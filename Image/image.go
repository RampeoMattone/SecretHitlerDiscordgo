package image

import (
	"SecretHitlerDiscordgo/database"
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/RampeoMattone/SecretGopher"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/lit"
	"github.com/fogleman/gg"
	"github.com/spf13/viper"
	"image"
	"io/ioutil"
	"net/http"
)

var (
	// Images of the various boards
	images = make(map[string]image.Image)
	// Cache for the image of the avatars
	avatarImageCache = make(map[string]image.Image)
	// Map for storing the status image with only avatars
	baseImages = make(map[*SecretGopher.GameState]image.Image)
)

const (
	// Coordinate for the upper-left corner for where to put policy cards
	// Fascist board
	fascistX = 168
	fascistY = 165
	// Liberal board
	liberalX = 268
	liberalY = 160

	// Offset from one card to the others
	// Fascist board
	fascistOffset = 204
	// Liberal board
	liberalOffset = 205

	// Dimension of the boards
	boardHeight = 1523
	boardWidth  = 567

	// Coordinate for the upper-left corner for where to put the election tracker
	electionTrackerX = 540.5
	electionTrackerY = 461.5
	// Offset from one tracker to the others
	electionTrackerOffset = 139.5
	// Radius of the polygon to draw
	electionTrackerRadius = 21.0

	// Dimension of the status
	statusHeight = 1000
	statusWidth  = 1555
)

func init() {
	lit.LogLevel = lit.LogError

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			lit.Error("Config file not found! See example_config.yml")
			return
		}
	} else {
		// DB setup
		database.InitializeDatabase(viper.GetString("drivername"), viper.GetString("datasourcename"))

		database.ExecQuery(database.TblUsers, database.DB)
		database.ExecQuery(database.TblGames, database.DB)
		database.ExecQuery(database.TblPlayers, database.DB)
		database.ExecQuery(database.TblRounds, database.DB)
		database.ExecQuery(database.TblActions, database.DB)

		// Loads all of the image used
		images["liberalBoard"], err = gg.LoadPNG("./assets/liberalBoard.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["liberalPolicy"], err = gg.LoadPNG("./assets/liberalPolicy.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["fascistBoard56"], err = gg.LoadPNG("./assets/fascistBoard_5-6.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["fascistBoard78"], err = gg.LoadPNG("./assets/fascistBoard_7-8.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["fascistBoard910"], err = gg.LoadPNG("./assets/fascistBoard_9-10.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["fascistPolicy"], err = gg.LoadPNG("./assets/fascistPolicy.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["fascistRole"], err = gg.LoadPNG("./assets/fascistRole.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["hitlerRole"], err = gg.LoadPNG("./assets/hitlerRole.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}

		images["liberalRole"], err = gg.LoadPNG("./assets/liberalRole.png")
		if err != nil {
			lit.Error("Error while loading file: %v", err)
		}
	}
}

// DownloadAvatar downloads avatar for a given user if it doesn't exist
func DownloadAvatar(u *discordgo.User) {
	// Check stored avatar
	var hash string
	err := database.DB.QueryRow("SELECT avatarHash FROM users WHERE id=?", u.ID).Scan(&hash)

	if hash == u.Avatar {
		// Avatar is already updated
		return
	}

	// If the user doesn't exist, add it
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		stm, _ := database.DB.Prepare("INSERT INTO users(id, avatarHash, name) VALUES(?, ?, ?)")
		_, err = stm.Exec(u.ID, u.Avatar, base64.StdEncoding.EncodeToString([]byte(u.Username)))
		if err != nil {
			lit.Error("Error while inserting user into database: %v", err)
			return
		}
	}

	// Start HTTP request
	resp, err := http.Get(u.AvatarURL("256"))
	if err != nil {
		lit.Error("Error while downloading file: %v", err)
		return
	}

	img, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	_, err = database.DB.Exec("UPDATE users SET avatarImage=? WHERE id=?", img, u.ID)
	if err != nil {
		lit.Error("Error inserting image into database: %v", err)
	}
}

// DrawFascistBoard draws the fascist board
func DrawFascistBoard(G *SecretGopher.GameState) *gg.Context {
	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Use the appropriate board type
	switch len(G.Roles) {
	case 7, 8:
		img.DrawImage(images["fascistBoard78"], 0, 0)

	case 9, 10:
		img.DrawImage(images["fascistBoard910"], 0, 0)

	default:
		img.DrawImage(images["fascistBoard56"], 0, 0)
	}

	// Draw the policy card
	var i int8
	for i = 0; i < G.FascistTracker; i++ {
		img.DrawImage(images["fascistPolicy"], fascistX+fascistOffset*int(i), fascistY)
	}

	return img
}

// DrawLiberalBoard draws the liberal board
func DrawLiberalBoard(G *SecretGopher.GameState) *gg.Context {
	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Draw the board
	img.DrawImage(images["liberalBoard"], 0, 0)

	// Draw the policy cards
	var i int8
	for i = 0; i < G.LiberalTracker; i++ {
		img.DrawImage(images["liberalPolicy"], liberalX+liberalOffset*int(i), liberalY)
	}

	// And the election tracker
	if G.ElectionTracker > 0 {
		// The circles aren't centered by a few pixels...
		var offset float64
		switch G.ElectionTracker {
		case 1:
			offset = 1.5

		case 2:
			offset = 1.0
		}

		// Draw the polygon
		img.DrawRegularPolygon(6, electionTrackerX+electionTrackerOffset*float64(G.ElectionTracker-1)+offset, electionTrackerY, electionTrackerRadius, 0.0)
		// Set color
		img.SetRGB(1, 0, 0)
		// File the polygon
		img.Fill()
	}

	return img
}

// Draws the base image for a given game with all of the avatars
func drawBase(players []*discordgo.User) *gg.Context {
	img := gg.NewContext(statusWidth, statusHeight)

	for i, p := range players {
		// If it's not loaded, load the image
		if _, ok := avatarImageCache[p.ID]; !ok {
			loadAvatar(p.ID)
		}

		if i < 5 {
			img.DrawImageAnchored(avatarImageCache[p.ID], 200+(i*300), 180, 0.5, 0.5)
		} else {
			img.DrawImageAnchored(avatarImageCache[p.ID], 200+((i-5)*300), 600, 0.5, 0.5)
		}
	}

	return img
}

// DrawStatus draws the status image for a given player
func DrawStatus(G *SecretGopher.GameState, forP int8, players []*discordgo.User) *gg.Context {

	if _, ok := baseImages[G]; !ok {
		baseImages[G] = drawBase(players).Image()
	}

	// Create new image
	img := gg.NewContext(statusWidth, statusHeight)
	// Draw the base on top of it
	img.DrawImage(baseImages[G], 0, 0)

	// Loads the font
	if err := img.LoadFontFace("./fonts/Karantina-Regular.ttf", 96); err != nil {
		lit.Error("Error while loading font: %v", err)
	}

	// Draws the president
	img.SetRGB(0, 0, 0)
	if G.President < 5 {
		img.DrawStringAnchored("President", 200+(float64(G.President)*300), 350, 0.5, 0.5)
	} else {
		img.DrawStringAnchored("President", 200+(float64(G.President-5)*300), 770, 0.5, 0.5)
	}

	// Draws the chancellor
	img.SetRGB(0, 0, 0)
	if G.Chancellor < 5 {
		img.DrawStringAnchored("Chancellor", 200+(float64(G.Chancellor)*300), 350, 0.5, 0.5)
	} else {
		img.DrawStringAnchored("Chancellor", 200+(float64(G.Chancellor-5)*300), 770, 0.5, 0.5)
	}

	// Draw fascist and Hitler if the user is a fascist or is Hitler with less then 6 players
	if G.Roles[forP] == SecretGopher.FascistParty || (G.Roles[forP] == SecretGopher.Hitler && len(G.Roles) <= 6) {
		for i, p := range G.Roles {
			switch p {
			case SecretGopher.FascistParty:
				if i < 5 {
					img.DrawImage(images["fascistRole"], 250+(i*300), 0)
				} else {
					img.DrawImage(images["fascistRole"], 250+((i-5)*300), 420)
				}

			case SecretGopher.Hitler:
				if i < 5 {
					img.DrawImage(images["hitlerRole"], 250+(i*300), 0)
				} else {
					img.DrawImage(images["hitlerRole"], 250+((i-5)*300), 420)
				}
			}
		}
	} else {
		// Else check if the user is Hitler to draw his card
		if G.Roles[forP] == SecretGopher.Hitler {
			if forP < 5 {
				img.DrawImage(images["hitlerRole"], 250+(int(forP)*300), 0)
			} else {
				img.DrawImage(images["hitlerRole"], 250+((int(forP)-5)*300), 420)
			}
		} else {
			// Else check if the user is a liberal to draw his card
			if G.Roles[forP] == SecretGopher.LiberalParty {
				if forP < 5 {
					img.DrawImage(images["liberalRole"], 250+(int(forP)*300), 0)
				} else {
					img.DrawImage(images["liberalRole"], 250+((int(forP)-5)*300), 420)
				}
			}
		}
	}

	return img
}

func loadAvatar(id string) {
	var b []byte
	_ = database.DB.QueryRow("SELECT avatarImage FROM users WHERE id=?", id).Scan(&b)

	avatarImageCache[id], _, _ = image.Decode(bytes.NewBuffer(b))
}
