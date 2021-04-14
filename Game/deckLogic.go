package Game

import (
	"math/rand"
)

// newDeck generates a new deck for a game and shuffles it ahead of time
func newDeck() Deck {
	var d = Deck{
		pos: 0,
		arr: [17]Policy{
			LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy,
		},
	}
	d.shuffle()
	return d
}

// shuffle shuffles the elements of the deck pseudorandomically
func (d Deck) shuffle() {
	d.pos = 0
	rand.Shuffle(17, func(i, j int) {
		d.arr[i], d.arr[j] = d.arr[j], d.arr[i]
	})
}

// draw draws the top n cards from the deck, making sure to move them away and not draw them again
func (d Deck) draw(n uint8) []Policy {
	var r = make([]Policy, n)
	for i := uint8(0); i < n; i++ {
		r[i] = d.arr[d.pos]
		d.pos++
	}
	// "If there are fewer than three tiles remaining in the policy deck at the end of a Legislative Session,
	// they are shuffled with the Discard pile to create a new policy deck. Unused policy tiles are not revealed."
	if d.pos > 14 {
		d.shuffle()
	}
	return r
}

// peek reveals the top 3 cards from the deck, making sure to leave the deck unaltered
func (d Deck) peek() [3]Policy {
	var r = [3]Policy{}
	for i := uint8(0); i < 3; i++ {
		r[i] = d.arr[d.pos+i]
	}
	return r
}
