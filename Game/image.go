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

var(
	liberalBoard    image.Image
	fascistBoard56  image.Image
	fascistBoard78  image.Image
	fascistBoard910 image.Image
	liberalPolicy   image.Image
	fascistPolicy   image.Image
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
	boardX = 1523
	boardY = 567

	electionTrackerX = 541
	electionTrackerY = 461
	electionTrackerOffset = 138
	electionTrackerRadius = 19

)

func init() {
	var err error

	liberalBoard, err = gg.LoadPNG("./Game/assets/liberalBoard.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}

	liberalPolicy, err = gg.LoadPNG("./Game/assets/liberalPolicy.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}

	fascistBoard56, err = gg.LoadPNG("./Game/assets/fascistBoard_5-6.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}

	fascistBoard78, err = gg.LoadPNG("./Game/assets/fascistBoard_7-8.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}

	fascistBoard910, err = gg.LoadPNG("./Game/assets/fascistBoard_9-10.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}

	fascistPolicy, err = gg.LoadPNG("./Game/assets/fascistPolicy.png")
	if err != nil {
		lit.Error("Error while loading file: ", err)
	}
}

// Downloads avatar for a given user if it doesn't exist
func downloadAvatar(u *discordgo.User) {
	path := "./avatars/"+u.ID+".png"
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return
	}

	resp, err := http.Get(u.Avatar)
	if err != nil {
		lit.Error("Error while downloading file: ", err)
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		lit.Error("Error while creating file: ", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		lit.Error("Error while copy data to file: ", err)
	}
}

// Draws the fascist board
func (g Game) drawFascistBoard() *gg.Context {
	img := gg.NewContext(boardX, boardY)

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

	var i uint8
	for i = 0; i < g.fascistBoard; i++ {
		img.DrawImage(fascistPolicy, fascistX + int(fascistOffset*i), fascistY)
	}

	return img
}

// Draws the liberal board
func (g Game) drawLiberalBoard() *gg.Context {
	img := gg.NewContext(boardX, boardY)

	img.DrawImage(liberalBoard, 0, 0)

	var i uint8
	for i = 0; i < g.liberalBoard; i++ {
		img.DrawImage(liberalPolicy, liberalX + int(liberalOffset*i), liberalY)
	}

	for i = 0; i < g.electionTracker; i++ {
		img.DrawCircle(electionTrackerX + float64(electionTrackerOffset*i), electionTrackerY, electionTrackerRadius)
	}

	return img
}
