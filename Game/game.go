// Package Game handles a game session
package Game

/*
Rules recap:
	Liberals win if one of the following happens:
		- five liberal laws are enacted
		- Hitler is assassinated
	Liberals win if one of the following happens:
		- six fascist laws are enacted
		- Hitler is elected chancellor after the third fascist policy is enacted

Setup:
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
		- TODO ... legislative session ...
*/
