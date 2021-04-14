// Package Game handles a game session
package Game

import (
	"SecretHitlerDiscordgo/Utils"
	"sync"
)

/*
Rules recap:
	Liberals win if one of the following happens:
		- five liberal laws are enacted
		- Hitler is assassinated
	Liberals win if one of the following happens:
		- six fascist laws are enacted
		- Hitler is elected Chancellor after the third fascist policy is enacted

Start:
			# PLAYERS			5 		6 		7 		8	 	9	 	10
			Liberals 			3 		4 		4 		5 		5 		6
			Fascists 			1+H 	1+H 	2+H 	2+H 	3+H 	3+H
	Fascists known by Hitler? 	Yes 	Yes 	No	 	No 		No 		No

Game Cycle:
	- the next player in queue gets the presidential title
	- the President chooses a candidate for the role of Chancellor ( out of the non term-limited Players )
	- all Players vote to elect the President + Chancellor pair ( Yes / No )
	* if the vote is a tie or a majority of Votes are No:
		- advance the election tracker by one step
		- check if the election tracker is on its last slot
		* if it's on the last slot:
			- the next Policy is revealed and enacted ( any power granted by this Policy is ignored )
			- the election tracker resets
			- any existing term-limits are forgotten
	* if the majority of Votes are Yes:
		- term-limits are updated for the new President and Chancellor
		* If three or more Fascist Policies have been enacted and the Chancellor is Hitler
			- fascists win
		- deafen and mute both the Chancellor and the President
		- reveal three policies to the President
		- let the President discard one of the three
		- let the Chancellor choose which of the two remaining policies to enact
		* if five fascists policies have been elected ( veto power )
			- ask both the Chancellor and the President whether they want to veto the policy election
			* if both agree
				- discard the remaining policy instead of enacting it
				- advance the election tracker by one step
				- check if the election tracker is on its last slot
						* if it's on the last slot:
							- the next Policy is revealed and enacted ( any power granted by this Policy is ignored )
							- the election tracker resets
							- any existing term-limits are forgotten
		- undeafen and unmute both the Chancellor and the President
		* if the enacted policy grants a presidential power
			+ investigate loyalty
				- the President chooses who to investigate ( a player may only be investigated once per game )
			+ call special election
				- the President chooses any other player to be the next President ( even those that are term-limited )
				- the round starts again without altering the queue for next President
				- the next round will pick up from the next player in queue
			+ policy peek
				- the President will see the top three cards in the deck
			+ execution
				- the President chooses a player to kill
				* if the player was Hitler
					- liberals win
				- the player is removed from the active members of the game and may only spectate
*/

type Role uint8

const (
	LiberalRole Role = 0
	FascistRole Role = 1
	HitlerRole  Role = 2
)

type Policy bool

const (
	LiberalPolicy Policy = true
	FascistPolicy Policy = false
)

type Player struct {
	id   string
	role Role
}

type Deck struct {
	arr [17]Policy
	pos uint8
}

type Game struct {
	in  chan interface{}
	out chan interface{}
	game
	lock sync.Mutex
}

type game struct {
	// MUTEXED - Public
	ElectionTracker uint8            // cycles from 0 to 3
	FascistTracker  uint8            // starts at 0 ( no cards ), ends at 6 ( 6 slots )
	LiberalTracker  uint8            // starts at 0 ( no cards ), ends at 5 ( 5 slots )
	President       *Player          // current President (elected or candidate)
	Chancellor      *Player          // current President (elected or candidate)
	Votes           map[*Player]bool // Votes for the government

	// UNMUTEXED - Public (they remain static since the start of the game)
	Id         int                // Id of the game
	Players    []Player           // collecion of the Players and roles
	PlayersMap map[string]*Player // maps discord ids to Players

	// UNMUTEXED - private
	deck         Deck      // deck for the game
	policyChoice []Policy  // deck to hold the policies that need to be enacted
	lastElected  Utils.Set // term limits for last Chancellor and last President
	executed     Utils.Set // pointer to Players who died
	turnNum      uint8     // used to calculate next President
	turnStage    Stage     // used to track the the turnNum's development
}

type Stage int8

const (
	Uninitialized Stage = iota
	ChancellorNeeded
	GovernmentElection
	PresidentPolicies
	ChancellorPolicies
	VetoVote
	VoteEnaction
	SpecialPower
)

type ElectionFeedback int8

const (
	VoteError ElectionFeedback = iota
	VoteAck
	Reject         // the vote was registered and the election was rejected
	RejectAndForce // the vote was registered, the election was rejected and a policy was drawn
	Pass           // the vote was registered and the election was approved
)

type PolicyFeedback int8

const (
	PolicyError PolicyFeedback = iota
	PolicyAck
)

type DidAnyoneWin int8

const (
	NobodyWon DidAnyoneWin = iota
	FascistsWon
	LiberalsWon
)

type SpecialPowers int8

const (
	Nothing SpecialPowers = iota
	Peek
	Investigate
	Election
	Execution
)

var powersTable = [3][6]SpecialPowers{
	{Nothing	, Nothing	 , Peek	   , Execution, Execution, Nothing},
	{Nothing	, Investigate, Election, Execution, Execution, Nothing},
	{Investigate, Investigate, Election, Execution, Execution, Nothing},
}
