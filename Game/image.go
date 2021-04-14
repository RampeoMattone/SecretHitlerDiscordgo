package Game

import (
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/lit"
	"github.com/fogleman/gg"
	"image"
	"io"
	"net/http"
	"os"
)

var (
	// Images of the various boards
	images = make(map[string]image.Image)
	// Cache for the image of the avatars
	avatarImageCache = make(map[string]image.Image)
	// Map for storing the status image with only avatars
	baseImages = make(map[int]image.Image)
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
	var err error

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
}

// DownloadAvatar downloads avatar for a given user if it doesn't exist
func DownloadAvatar(u *discordgo.User) {
	// If the avatar already exist, just return
	path := "./avatars/" + u.ID + ".png"
	_, err := os.Stat(path)
	if os.IsExist(err) {
		return
	}

	// Start HTTP request
	resp, err := http.Get(u.AvatarURL("256"))
	if err != nil {
		lit.Error("Error while downloading file: %v", err)
		return
	}
	defer resp.Body.Close()

	// Creates file
	f, err := os.Create(path)
	if err != nil {
		lit.Error("Error while creating file: %v", err)
		return
	}
	defer f.Close()

	// And writes the data
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		lit.Error("Error while copy data to file: %v", err)
	}
}

// DrawFascistBoard draws the fascist board
func (G *Game) DrawFascistBoard() *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()

	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Use the appropriate board type
	switch len(G.game.players) {
	case 7, 8:

		img.DrawImage(images["fascistBoard78"], 0, 0)
		break
	case 9, 10:

		img.DrawImage(images["fascistBoard910"], 0, 0)
		break
	default:

		img.DrawImage(images["fascistBoard56"], 0, 0)
	}

	// Draw the policy card
	var i uint8
	for i = 0; i < G.game.fascistBoard; i++ {
		img.DrawImage(images["fascistPolicy"], fascistX+fascistOffset*int(i), fascistY)
	}

	return img
}

// DrawLiberalBoard draws the liberal board
func (G *Game) DrawLiberalBoard() *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()

	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Draw the board
	img.DrawImage(images["liberalBoard"], 0, 0)

	// Draw the policy cards
	var i uint8
	for i = 0; i < G.game.liberalBoard; i++ {
		img.DrawImage(images["liberalPolicy"], liberalX+liberalOffset*int(i), liberalY)
	}

	// And the election tracker
	if G.game.electionTracker > 0 {
		// The circles aren't centered by a few pixels...
		var offset float64
		switch G.game.electionTracker {
		case 1:
			offset = 1.5
			break
		case 2:
			offset = 1.0
		}

		// Draw the polygon
		img.DrawRegularPolygon(6, electionTrackerX+electionTrackerOffset*float64(G.game.electionTracker-1)+offset, electionTrackerY, electionTrackerRadius, 0.0)
		// Set color
		img.SetRGB(1, 0, 0)
		// File the polygon
		img.Fill()
	}

	return img
}

// Draws the base image for a given game with all of the avatars
func (G *Game) drawBase() *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()

	img := gg.NewContext(statusWidth, statusHeight)

	for i, p := range G.game.players {
		// If it's not loaded, load the image
		if _, ok := avatarImageCache[p.id]; !ok {
			loadAvatar(p.id)
		}

		if i < 5 {
			img.DrawImageAnchored(avatarImageCache[p.id], 200+(i*300), 180, 0.5, 0.5)
		} else {
			img.DrawImageAnchored(avatarImageCache[p.id], 200+((i-5)*300), 600, 0.5, 0.5)
		}
	}

	return img
}

// DrawStatus draws the status image for a given player
func (G *Game) DrawStatus(forP *Player) *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()

	if _, ok := baseImages[G.game.id]; !ok {
		baseImages[G.game.id] = G.drawBase().Image()
	}

	// Create new image
	img := gg.NewContext(statusWidth, statusHeight)
	// Draw the base on top of it
	img.DrawImage(baseImages[G.game.id], 0, 0)

	// Loads the font
	if err := img.LoadFontFace("./fonts/Karantina-Regular.ttf", 96); err != nil {
		lit.Error("Error while loading font: %v", err)
	}

	var president, chancellor bool
	for i, p := range G.game.players {
		// Draws the president
		if G.game.president.id == p.id {
			img.SetRGB(0, 0, 0)

			if i < 5 {
				img.DrawStringAnchored("President", 200+(float64(i)*300), 350, 0.5, 0.5)
			} else {
				img.DrawStringAnchored("President", 200+(float64(i-5)*300), 770, 0.5, 0.5)
			}
			president = true
		} else {
			// Draws the chancellor
			if G.game.chancellor.id == p.id {
				img.SetRGB(0, 0, 0)

				if i < 5 {
					img.DrawStringAnchored("Chancellor", 200+(float64(i)*300), 350, 0.5, 0.5)
				} else {
					img.DrawStringAnchored("Chancellor", 200+(float64(i-5)*300), 770, 0.5, 0.5)
				}
				chancellor = true
			}
		}

		// If we drew both text
		if president && chancellor {
			break
		}
	}


	if forP.role == FASCIST_ROLE || (forP.role == HITLER_ROLE && len(G.game.players) <= 6) {
		for i, p := range G.game.players {
			switch p.role {
			case FASCIST_ROLE:
				if i < 5 {
					img.DrawImage(images["fascistRole"], 250+(i*300), 0)
				} else {
					img.DrawImage(images["fascistRole"], 250+((i-5)*300), 420)
				}
				break

			case HITLER_ROLE:
				if i < 5 {
					img.DrawImage(images["hitlerRole"], 250+(i*300), 0)
				} else {
					img.DrawImage(images["hitlerRole"], 250+((i-5)*300), 420)
				}
				break
			}
		}
	}

	return img
}

func loadAvatar(id string) {
	var err error
	avatarImageCache[id], err = gg.LoadPNG("./avatars/" + id + ".png")
	if err != nil {
		lit.Error("Error loading avatar: %v", err)
	}
}
