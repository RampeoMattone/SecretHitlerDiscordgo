package Game

import (
	"SecretHitlerDiscordgo/Utils"
	"math/rand"
)

// GovernmentCastVote will cast a vote for the election
func (G *Game) governmentCastVote(usr string, v bool) ElectionFeedback {
	G.Lock.Lock()
	defer G.Lock.Unlock()
	var (
		g     = G.game
		p     = g.PlayersMap[usr] // pointer to the player
		_, ok = g.Votes[p]
	)
	if !ok {
		g.Votes[p] = v
		if len(g.Votes) == len(g.Players) {
			var r int8 = 0
			for _, v := range g.Votes {
				switch v {
				case true:
					r += 1
				case false:
					r -= 1
				}
			}
			if r > 0 { // election successful
				g.lastElected.Clear()
				g.lastElected.AddAll(g.President, g.Chancellor)
				g.policyChoice = g.deck.draw(3)
				return Pass
			} else { // election failed
				// election tracker forces policies
				if g.ElectionTracker == 2 {
					return RejectAndForce
				}
				g.ElectionTracker++
				return Reject
			}
		}
		return VoteAck
	}
	return VoteError
}

// policyChoice will remove the chosen card from the policyChoice deck
func (G *Game) policyChoice(c *string, s uint8) PolicyFeedback {
	G.Lock.Lock()
	defer G.Lock.Unlock()
	var g = G.game
	switch g.turnStage {
	case PresidentPolicies:
		if s > 2 || *c != g.President.Id {
			return PolicyError
		}
		g.turnStage = ChancellorPolicies
	case ChancellorPolicies:
		if s > 1 || *c != g.Chancellor.Id {
			return PolicyError
		}
		if g.FascistTracker >= 5 {
			g.turnStage = VetoVote
		} else {
			g.turnStage = VoteEnaction
		}
	default:
		return PolicyError
	}
	g.policyChoice = append(g.policyChoice[:s], g.policyChoice[s+1:]...)
	return PolicyAck
}

// NewPresident at the start of a new turn, it will elect a new President and signal its userid via a channel
func (g game) newPresident() { // we have to use a channel because the function is ran In a goroutine
	g.turnNum++
	var p = g.Players[len(g.Players)%int(g.turnNum)]
	g.President = &p
	g.turnStage = ChancellorNeeded
}

// enactPolicy enacts a policy and echoes Out the winning party if there is one
func (G *Game) enactPolicy() (DidAnyoneWin, SpecialPowers) {
	G.Lock.Lock()
	defer G.Lock.Unlock()
	var g = G.game
	return g.enactPolicyUnsafe()
}

// enactPolicyUnsafe enacts a policy and echoes Out the winning party if there is one
func (g game) enactPolicyUnsafe() (DidAnyoneWin, SpecialPowers) {
	switch g.policyChoice[0] {
	case FascistPolicy:
		g.FascistTracker++
		if g.FascistTracker == 6 {
			return FascistsWon, 0
		}
	case LiberalPolicy:
		g.LiberalTracker++
		if g.LiberalTracker == 6 {
			return LiberalsWon, 0
		}
	}
	var p SpecialPowers
	switch len(g.Players) {
	case 5, 6:
		p = powersTable[0][g.FascistTracker]
	case 7, 8:
		p = powersTable[1][g.FascistTracker]
	case 9, 10:
		p = powersTable[2][g.FascistTracker]
	}
	return NobodyWon, p
}

// setRoles sets one player as HitlerRole and f Players as FascistRole
func (g game) setRoles(f int) {
	var rSet = make(Utils.Set) // make a set to store numbers we extract
	var r int                  // stores a random number
	// set #f Players as FascistRole
	for i := 0; i < f; i++ {
		r = rand.Intn(10) // extract a random number
		// if the number has already been extracted, reroll
		for rSet.Has(r) {
			r = rand.Intn(10)
		}
		rSet.Add(r)                     // add the extracted number to the set
		g.Players[r].Role = FascistRole // set the player Role
	}
	// set a player as HitlerRole
	for rSet.Has(r) {
		r = rand.Intn(10) // extract a random number
	}
	g.Players[r].Role = HitlerRole // set the player Role
}
