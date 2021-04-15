package database

import (
	"database/sql"
	"github.com/bwmarrin/lit"
	// Possible drivers
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

// DB is the database connection
var DB *sql.DB

// Various table used
const (
	TblUsers   = "CREATE TABLE IF NOT EXISTS `users`( `id` varchar(18) NOT NULL, `avatarHash` varchar(32) NOT NULL, `name` varchar(200) NOT NULL, `avatarImage` mediumblob, PRIMARY KEY (`id`))"
	TblGames   = "CREATE TABLE IF NOT EXISTS `games`( `id` int(11) unsigned NOT NULL AUTO_INCREMENT, `startedAt` datetime NOT NULL, `finishedAt` datetime DEFAULT NULL, `result` enum('FASCIST','LIBERAL') DEFAULT NULL, PRIMARY KEY (`id`))"
	TblPlayers = "CREATE TABLE IF NOT EXISTS `players`( `id` int(10) unsigned NOT NULL, `userID` varchar(18) NOT NULL, `gameID` int(10) unsigned NOT NULL, `role` enum('FASCIST','LIBERAL') NOT NULL, PRIMARY KEY (`id`), KEY `FK__users` (`userID`), KEY `FK_players_games` (`gameID`), CONSTRAINT `FK__users` FOREIGN KEY (`userID`) REFERENCES `users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION, CONSTRAINT `FK_players_games` FOREIGN KEY (`gameID`) REFERENCES `games` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION)"
	TblRounds  = "CREATE TABLE IF NOT EXISTS `rounds`( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `gameID` int(10) unsigned NOT NULL, `chancellor` varchar(18) NOT NULL, `president` varchar(18) NOT NULL, `policy` enum('FASCIST','LIBERAL') DEFAULT NULL, PRIMARY KEY (`id`), KEY `FK__players` (`chancellor`), KEY `FK__players_2` (`president`), KEY `FK_rounds_games` (`gameID`), CONSTRAINT `FK__players` FOREIGN KEY (`chancellor`) REFERENCES `players` (`userID`) ON DELETE NO ACTION ON UPDATE NO ACTION, CONSTRAINT `FK__players_2` FOREIGN KEY (`president`) REFERENCES `players` (`userID`) ON DELETE NO ACTION ON UPDATE NO ACTION, CONSTRAINT `FK_rounds_games` FOREIGN KEY (`gameID`) REFERENCES `games` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION)"
	TblActions = "CREATE TABLE IF NOT EXISTS `actions`( `id` int(10) unsigned NOT NULL AUTO_INCREMENT, `roundID` int(10) unsigned NOT NULL, `action` enum('PEEK','KILL') NOT NULL, `killed` int(10) unsigned DEFAULT NULL, PRIMARY KEY (`id`), KEY `FK__rounds` (`roundID`), KEY `FK_actions_players` (`killed`), CONSTRAINT `FK__rounds` FOREIGN KEY (`roundID`) REFERENCES `rounds` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION, CONSTRAINT `FK_actions_players` FOREIGN KEY (`killed`) REFERENCES `players` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION)"
)

// InitializeDatabase initializes db connection given the driver and datasource name.
// Supported driver are sqlite and mysql.
func InitializeDatabase(driver, dsn string) {
	// Open database connection
	if DB == nil {
		var err error
		DB, err = sql.Open(driver, dsn)
		if err != nil {
			lit.Error("Error opening db connection: %v", err)
			return
		}
	}
}

// ExecQuery Executes a simple query given a DB
func ExecQuery(query string, db *sql.DB) {
	_, err := db.Exec(query)
	if err != nil {
		lit.Error("Error creating table: %v: %v", err)
		return
	}
}
