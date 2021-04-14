package Game

import (
	"SecretHitlerDiscordgo/Utils"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// collection of functions that handle the flow of the game

// NewGame creates the basic structure for the game to let others join using the Join command
func NewGame() *Game {
	G := Game{
		Out:  make(chan interface{}),
		In:   make(chan interface{}),
		// Lock: sync.Mutex{}, // todo test if needed
		game: game{
			Players:   make([]Player, 10),
			turnStage: Uninitialized,
			// deck:        newDeck(),   // TODO init later on
			// lastElected: Utils.Set{}, // TODO init later on
			// executed:    Utils.Set{}, // TODO init later on
		},
	}
	go G.Handler()
	return &G
}

func (G *Game) Handler() {
	g := G.game
	vetoVoted := Utils.Set{}
	vetoResult := uint8(0)
	defer close(G.Out)
	defer close(G.In)
	for {
		in := <-G.In
		switch in.(type) {
		// Join the game
		case Join:
			if G.game.turnStage == Uninitialized {
				id := in.(Join)
				p := Player{
					Id: *id,
				}
				if len(g.Players) > 10 {
					G.Out <- Error(LobbyFull)
					continue
				}
				g.PlayersMap[*id] = &p
				g.Players = append(g.Players, p)
				G.Out <- Ack{}
			} else {
				G.Out <- Error(WrongPhase)
			}
		// Start the game
		case Start:
			if g.turnStage == Uninitialized {
				switch len(g.Players) {
				case 5, 6:
					g.setRoles(1)
				case 7, 8:
					g.setRoles(2)
				case 9, 10:
					g.setRoles(3)
				default:
					G.Out <- Error(NotEnoughPlayers)
					continue
				}
				g.newPresident()
				G.Out <- Ack{}
			} else {
				G.Out <- Error(WrongPhase)
			}
		case ElectChancellor:
			e := in.(ElectChancellor)
			if g.turnStage == ChancellorNeeded && g.President.Id == *e.caller {
				if v, ok := g.PlayersMap[*e.proposal]; ok && !g.lastElected.Has(*e.proposal) {
					g.Chancellor = v
					g.Votes = make(map[*Player]bool, 10)
					g.turnStage = GovernmentElection
					G.Out <- Ack{}
				} else {
					G.Out <- Error(InvalidInput)
				}
			} else {
				G.Out <- Error(WrongPhase)
			}
		case VoteGovernment:
			if g.turnStage == GovernmentElection {
				e := in.(VoteGovernment)
				f := G.governmentCastVote(*e.caller, e.vote)
				switch f {
				case VoteAck:
					G.Out <- Ack{}
				case VoteError:
					G.Out <- Error(InvalidInput)
				case Reject:
					G.Out <- GovernmentElectionFailed{ ForcedPolicy: false }
					g.newPresident()
				case RejectAndForce:
					g.ElectionTracker = 0           // reset the tracker
					g.policyChoice = g.deck.draw(1) // draw the policy
					g.enactPolicyUnsafe()           // todo handle return types and win conditions
					G.Out <- GovernmentElectionFailed{ ForcedPolicy: true }
					g.newPresident()
				case Pass:
					g.turnStage = PresidentPolicies
					G.Out <- GovernmentElectionSuccess{}
				}
			} else {
				G.Out <- Error(WrongPhase)
			}
		case PolicyChoice:
			e := in.(PolicyChoice)
			switch G.policyChoice(e.caller, e.selected) {
			case PolicyAck:
				switch g.turnStage {
				case ChancellorPolicies:
					G.Out <- Ack{}
				case VetoVote:
					G.Out <- StartVeto{}
				case VoteEnaction:
					win, power := G.enactPolicy()
					G.Out <- PolicyEnacted{
						win: win,
						power: power,
					}
					// todo if won, close the loop
				}
			case PolicyError:
				G.Out <- Error(WrongPhase)
			}
		case PolicyVeto:
			if g.turnStage == VetoVote {
				e := in.(PolicyVeto)
				if !vetoVoted.Has(*e.caller) && (g.PlayersMap[*e.caller] == g.President || g.PlayersMap[*e.caller] == g.
					Chancellor) {
					vetoVoted.Add(*e.caller)
					vetoResult++
					if len(vetoVoted) == 2 {
						if vetoResult == 2 {
							// G.Lock.Lock() // todo test if needed
							g.newPresident()
							if g.ElectionTracker == 2 {
								g.ElectionTracker = 0
								g.policyChoice = g.deck.draw(1)
								win, power := g.enactPolicyUnsafe()
								G.Out <- VetoResult{
									success: true,
									force: true,
									enacted: PolicyEnacted{
										win: win,
										power: power,
									},
								}
							} else {
								g.ElectionTracker++
								G.Out <- VetoResult{
									success: true,
									force: false,
								}
							}
							// G.Lock.Unlock() // todo test if needed
						} else {
							win, power := G.enactPolicy()
							G.Out <- VetoResult{
								success: false,
								enacted: PolicyEnacted{
									win: win,
									power: power,
								},
							}
						}
					} else {
						G.Out <- Ack{}
					}
				} else {
					G.Out <- Error(WrongPhase)
				}
			}
		}
	}
}
