package Game

import (
	"SecretHitlerDiscordgo/Utils"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// NewGame creates the basic structure for the game to let others join using the Join command
func NewGame() Game {
	return Game{
		game: game{
			deck:        newDeck(),
			players:     make([]Player, 10),
			lastElected: Utils.Set{},
			executed:    Utils.Set{},
			turnStage:   UNINITIALIZED,
		},
		lock: sync.RWMutex{},
	}
}

func (G *Game) Join(usr *discordgo.User) bool {
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	var p = Player{
		id:   usr.ID,
		role: LIBERAL_ROLE,
		name: usr.Username,
	}
	if len(g.players) > 10 {
		return false
	}
	g.playersMap[usr.ID] = &p
	g.players = append(g.players, p)
	go DownloadAvatar(usr)
	return true
}

// setRoles sets one player as HITLER_ROLE and f players as FASCIST_ROLE
func (g game) setRoles(f int) {
	var rSet = make(Utils.Set) // make a set to store numbers we extract
	var r int                  // stores a random number
	// set #f players as FASCIST_ROLE
	for i := 0; i < f; i++ {
		r = rand.Intn(10) // extract a random number
		// if the number has already been extracted, reroll
		for rSet.Has(r) {
			r = rand.Intn(10)
		}
		rSet.Add(r)                      // add the extracted number to the set
		g.players[r].role = FASCIST_ROLE // set the player role
	}
	// set a player as HITLER_ROLE
	for rSet.Has(r) {
		r = rand.Intn(10) // extract a random number
	}
	g.players[r].role = HITLER_ROLE // set the player role
}

// Start is ran at the end of the join phase, once the admin has agreed to start
// the function sets up the game's parts that can't be set up ahead of time
func (G *Game) Start(out chan<- string) bool {
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	if g.turnStage == UNINITIALIZED {
		g.turnStage = PRESIDENT_NEEDED
		switch len(g.players) {
		case 5, 6:
			g.setRoles(1)
		case 7, 8:
			g.setRoles(2)
		case 9, 10:
			g.setRoles(3)
		default:
			return false
		}
		go G.NewPresident(out)
	} else {
		return false
	}
	return true
}

// NewPresident at the start of a new turn, it will elect a new president and signal its userid via a channel
func (G *Game) NewPresident(out chan<- string) { // we have to use a channel because the function is ran in a goroutine
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	if g.turnStage == PRESIDENT_NEEDED {
		g.turnStage = CHANCELLOR_NEEDED
		g.turnNum++
		var p = g.players[len(g.players)%int(g.turnNum)]
		g.president = &p
		out <- p.id
	} else {
		out <- "" // empty string means error
	}
}

// NewChancellor will try to propose a chancellor
func (G *Game) NewChancellor(usr string, c string) bool {
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	if g.turnStage == CHANCELLOR_NEEDED && g.president.id == usr {
		g.turnStage = GOVERNMENT_ELECTION
		g.chancellor = g.playersMap[c]
		g.votes = make(map[*Player]bool, 10)
		return true
	} else {
		return false // empty string means error
	}
}

// GovernmentCastVote will cast a vote for the election
func (G *Game) GovernmentCastVote(usr string, v bool) ElectionFeedback {
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	var p = g.playersMap[usr] // pointer to the player
	var _, ok = g.votes[p]
	if g.turnStage == GOVERNMENT_ELECTION && !ok {
		g.votes[p] = v
		if len(g.votes) == len(g.players) {
			var r int8 = 0
			for _, v := range g.votes {
				switch v {
				case true:
					r += 1
				case false:
					r -= 1
				}
			}
			if r > 0 { // election successful
				g.lastElected.Clear()
				g.lastElected.AddAll(g.president, g.chancellor)
				g.turnStage = PRESIDENT_POLICIES
				g.policyChoice = g.deck.draw(3)
				return PASS
			} else { // election failed
				// election tracker forces policies
				if g.electionTracker == 2 {
					g.electionTracker = 0         // reset the tracker
					var policies = g.deck.draw(1) // draw the policy
					var policy = policies[0]
					if policy == LIBERAL_POLICY {
						g.liberalBoard++
					} else {
						g.fascistBoard++
					}
					return REJECT_AND_FORCE
				}
				g.electionTracker++
				return REJECT
			}
		}
		return ACK
	}
	return ElectionFeedback(ERROR)
}

// Choice will remove the chosen card from the policyChoice deck
func (G *Game) Choice(c uint8, s Stage) bool {
	G.lock.Lock()
	defer G.lock.Unlock()
	var g = G.game
	if g.turnStage == s {
		switch s {
		case PRESIDENT_POLICIES:
			if c > 2 {
				return false
			}
			g.turnStage = CHANCELLOR_POLICIES
		case CHANCELLOR_POLICIES:
			if c > 1 {
				return false
			}
			if g.fascistBoard >= 5 {
				g.turnStage = VETO_VOTE
			} else {
				g.turnStage = SPECIAL_POWER
			}
		default:
			return false
		}
		g.policyChoice = append(g.policyChoice[:c], g.policyChoice[c+1:]...)
		return true
	}
	return false
}
