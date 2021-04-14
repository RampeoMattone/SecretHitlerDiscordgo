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
		- Hitler is elected chancellor after the third fascist policy is enacted

Start:
			# PLAYERS			5 		6 		7 		8	 	9	 	10
			Liberals 			3 		4 		4 		5 		5 		6
			Fascists 			1+H 	1+H 	2+H 	2+H 	3+H 	3+H
	Fascists known by Hitler? 	Yes 	Yes 	No	 	No 		No 		No

Game Cycle:
	- the next player in queue gets the presidential title
	- the president chooses a candidate for the role of chancellor ( out of the non term-limited players )
	- all players vote to elect the president + chancellor pair ( Yes / No )
	* if the vote is a tie or a majority of votes are No:
		- advance the election tracker by one step
		- check if the election tracker is on its last slot
		* if it's on the last slot:
			- the next Policy is revealed and enacted ( any power granted by this Policy is ignored )
			- the election tracker resets
			- any existing term-limits are forgotten
	* if the majority of votes are Yes:
		- term-limits are updated for the new president and chancellor
		* If three or more Fascist Policies have been enacted and the chancellor is Hitler
			- fascists win
		- deafen and mute both the chancellor and the president
		- reveal three policies to the president
		- let the president discard one of the three
		- let the chancellor choose which of the two remaining policies to enact
		* if five fascists policies have been elected ( veto power )
			- ask both the chancellor and the president whether they want to veto the policy election
			* if both agree
				- discard the remaining policy instead of enacting it
				- advance the election tracker by one step
				- check if the election tracker is on its last slot
						* if it's on the last slot:
							- the next Policy is revealed and enacted ( any power granted by this Policy is ignored )
							- the election tracker resets
							- any existing term-limits are forgotten
		- undeafen and unmute both the chancellor and the president
		* if the enacted policy grants a presidential power
			+ investigate loyalty
				- the president chooses who to investigate ( a player may only be investigated once per game )
			+ call special election
				- the president chooses any other player to be the next president ( even those that are term-limited )
				- the round starts again without altering the queue for next president
				- the next round will pick up from the next player in queue
			+ policy peek
				- the president will see the top three cards in the deck
			+ execution
				- the president chooses a player to kill
				* if the player was Hitler
					- liberals win
				- the player is removed from the active members of the game and may only spectate
*/

const ERROR int8 = -1 // general error

type Role uint8

const (
	LIBERAL_ROLE Role = 0
	FASCIST_ROLE Role = 1
	HITLER_ROLE  Role = 2
)

type Policy bool

const (
	LIBERAL_POLICY Policy = true
	FASCIST_POLICY Policy = false
)

type Player struct {
	id   string
	role Role
	name string
}

type Deck struct {
	arr [17]Policy
	pos uint8
}

type Game struct {
	game
	lock			sync.RWMutex
}

type game struct {
	id          int
	players     []Player
	deck        Deck
	playersMap  map[string]*Player // maps discord ids to players
	president   *Player
	chancellor  *Player
	votes       map[*Player]bool
	lastElected Utils.Set // term limits for last chancellor and last president
	executed    Utils.Set // pointer to players who died
	turnNum     uint8     // used to calculate next president
	turnStage       Stage     // used to track the the turnNum's development
	electionTracker uint8     // cycles from 0 to 3
	fascistBoard    uint8     // starts at 0 ( no cards ), ends at 6 ( 6 slots )
	liberalBoard    uint8     // starts at 0 ( no cards ), ends at 5 ( 5 slots )
}

type Stage int8

const (
	UNINITIALIZED       Stage = -1
	PRESIDENT_NEEDED    Stage = 0
	CHANCELLOR_NEEDED   Stage = 1
	GOVERNMENT_ELECTION Stage = 2
	PRESIDENT_POLICIES  Stage = 3
	CHANCELLOR_POLICIES Stage = 4
	VETO_VOTE           Stage = 5
	SPECIAL_POWER       Stage = 6
)

type ElectionFeedback int8

const (
	ACK    ElectionFeedback = 0 // the vote was registered
	REJECT ElectionFeedback = 1 // the vote was registered and the election was rejected
	PASS   ElectionFeedback = 2 // the vote was registered and the election was approved
)
