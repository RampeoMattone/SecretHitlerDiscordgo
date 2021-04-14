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
	liberalBoard    image.Image
	fascistBoard56  image.Image
	fascistBoard78  image.Image
	fascistBoard910 image.Image
	liberalPolicy   image.Image
	fascistPolicy   image.Image
	// Cache for the image of the avatars
	avatarImageCache = make(map[string]image.Image)
	//
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
	statusHeight = 820
	statusWidth  = 1555
)

func init() {
	var err error

	// Loads all of the image used
	liberalBoard, err = gg.LoadPNG("./Game/assets/liberalBoard.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}

	liberalPolicy, err = gg.LoadPNG("./Game/assets/liberalPolicy.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}

	fascistBoard56, err = gg.LoadPNG("./Game/assets/fascistBoard_5-6.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}

	fascistBoard78, err = gg.LoadPNG("./Game/assets/fascistBoard_7-8.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}

	fascistBoard910, err = gg.LoadPNG("./Game/assets/fascistBoard_9-10.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}

	fascistPolicy, err = gg.LoadPNG("./Game/assets/fascistPolicy.png")
	if err != nil {
		lit.Error("Error while loading file: %v", err)
	}
}

// Downloads avatar for a given user if it doesn't exist
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

// Draws the fascist board
func (G *Game) DrawFascistBoard() *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()
	var g = G.game

	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Use the appropriate board type
	switch len(g.players) {
	case 7, 8:
		img.DrawImage(fascistBoard78, 0, 0)
		break
	case 9, 10:
		img.DrawImage(fascistBoard910, 0, 0)
		break
	default:
		img.DrawImage(fascistBoard56, 0, 0)
	}

	// Draw the policy card
	var i uint8
	for i = 0; i < g.fascistBoard; i++ {
		img.DrawImage(fascistPolicy, fascistX+fascistOffset*int(i), fascistY)
	}

	return img
}

// Draws the liberal board
func (G *Game) DrawLiberalBoard() *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()
	var g = G.game

	// Create new blank image with boardHeight and boardWidth dimensions
	img := gg.NewContext(boardHeight, boardWidth)

	// Draw the board
	img.DrawImage(liberalBoard, 0, 0)

	// Draw the policy cards
	var i uint8
	for i = 0; i < g.liberalBoard; i++ {
		img.DrawImage(liberalPolicy, liberalX+liberalOffset*int(i), liberalY)
	}

	// And the election tracker
	if g.electionTracker > 0 {
		// The circles aren't centered by a few pixels...
		var offset float64
		switch g.electionTracker {
		case 1:
			offset = 1.5
			break
		case 2:
			offset = 1.0
		}

		// Draw the polygon
		img.DrawRegularPolygon(6, electionTrackerX+electionTrackerOffset*float64(g.electionTracker-1)+offset, electionTrackerY, electionTrackerRadius, 0.0)
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
	var g = G.game

	img := gg.NewContext(statusWidth, statusHeight)

	for i, p := range g.players {
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

// Draws the status image for a given player
func (G *Game) DrawStatus(forP *Player) *gg.Context {
	G.lock.RLock()
	defer G.lock.RUnlock()
	var g = G.game

	if _, ok := baseImages[g.id]; !ok {
		baseImages[g.id] = G.drawBase().Image()
	}

	img := gg.NewContext(statusWidth, statusHeight)

	return img
}

func loadAvatar(id string) {
	var err error
	avatarImageCache[id], err = gg.LoadPNG("./avatars/" + id + ".png")
	if err != nil {
		lit.Error("Error loading avatar: %v", err)
	}
}
