package Game

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func (g Game) join(usr *discordgo.User) bool {
	var p = Player{
		id:   usr.ID,
		role: LIBERAL_ROLE,
		name: usr.Username,
	}
	if len(g.players) == 0 {
		g.players = make([]Player, 10)
	}
	if len(g.players) > 10 {
		return false
	}
	g.players = append(g.players, p)
	return true
}

// sets one player as HITLER_ROLE and f players as FASCIST_ROLE
func (g Game) setRoles(f int) {
	var rSet = make(map[int]bool) // make a set to store numbers we extract
	var r int // stores a random number
	// set #f players as FASCIST_ROLE
	for i := 0; i < f; i++ {
		r = rand.Intn(10) // extract a random number
		// if the number has already been extracted, reroll
		for rSet[r] {
			r = rand.Intn(10)
		}
		rSet[r] = true                   // add the extracted number to the set
		g.players[r].role = FASCIST_ROLE // set the player role
	}
	// set a player as HITLER_ROLE
	for rSet[r] {
		r = rand.Intn(10) // extract a random number
	}
	g.players[r].role = HITLER_ROLE // set the player role
}

// init function for game setup
func (g Game) init() bool {
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
	return true
}