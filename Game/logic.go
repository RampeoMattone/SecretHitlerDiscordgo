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

// newDeck generates a new deck for a game and shuffles it ahead of time
func newDeck() Deck {
	var d = Deck{
		pos: 0,
		arr: [17]Policy{
			LIBERAL_POLICY, LIBERAL_POLICY, LIBERAL_POLICY, LIBERAL_POLICY, LIBERAL_POLICY, LIBERAL_POLICY,
		},
	}
	d.shuffle()
	return d
}

// shuffle shuffles the elements of the deck pseudorandomically
func (d Deck) shuffle() {
	rand.Shuffle(17, func(i, j int) {
		d.arr[i], d.arr[j] = d.arr[j], d.arr[i]
	})
}

// NewGame creates the basic structure for the game to let others join using the Join command
func NewGame() Game {
	return Game{
		deck:            newDeck(),
		players:         make([]Player, 10),
		playersMap:      make(map[string]*Player, 10),
		lastElected:     Utils.Set{},
		executed:        Utils.Set{},
		turnNum:         0,
		turnStage:       UNINITIALIZED,
		electionTracker: 0,
		fascistBoard:    0,
		liberalBoard:    0,
		mut:             sync.Mutex{},
	}
}

func (g *Game) Join(usr *discordgo.User) bool {
	g.mut.Lock()
	defer g.mut.Unlock()
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
func (g Game) setRoles(f int) {
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
func (g *Game) Start(out chan<- string) bool {
	g.mut.Lock()
	defer g.mut.Unlock()
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
		go g.NewPresident(out)
	} else {
		return false
	}
	return true
}

// NewPresident at the start of a new turn, it will elect a new president and signal its userid via a channel
func (g *Game) NewPresident(out chan<- string) { // we have to use a channel because the function is ran in a goroutine
	g.mut.Lock()
	defer g.mut.Unlock()
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
func (g *Game) NewChancellor(usr string, c string) bool {
	g.mut.Lock()
	defer g.mut.Unlock()
	if g.turnStage == CHANCELLOR_NEEDED && g.president.id == usr {
		g.turnStage = GOVERNMENT_ELECTION
		g.chancellor = g.playersMap[c]
		return true
	} else {
		return false // empty string means error
	}
}

// GovernmentCastVote will cast a vote for the election
func (g *Game) GovernmentCastVote(usr string, v bool) ElectionFeedback {
	g.mut.Lock()
	defer g.mut.Unlock()
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
				// TODO add deck interaction for the president
				return PASS
			} else { // election failed
				return REJECT
			}
		}
		return ACK
	}
	return ElectionFeedback(ERROR)
}
