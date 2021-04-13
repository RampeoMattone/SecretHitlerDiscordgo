package Game

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"testing"
)

func init() {
	_ = os.Mkdir("./temp", 0755)
	_ = os.Mkdir("./avatars", 0755)
}

func TestDrawFascistBoard(t *testing.T) {
	g := Game{
		fascistBoard: 6,
	}

	g.DrawFascistBoard().SavePNG("./temp/fascist-3.png")
}

func TestDrawLiberalBoard(t *testing.T) {
	g := Game{
		electionTracker: 3,
		liberalBoard:    5,
	}

	g.DrawLiberalBoard().SavePNG("./temp/liberal-3-2.png")
}

func TestDownloadAvatar(t *testing.T) {
	var (
		users = []discordgo.User{
			{
				ID:       "145618075452964864",
				Username: "TheTipo01",
				Avatar:   "93d255afb6f8d89fab55360edad0a9ef",
			},
			{
				ID:       "143060848091463680",
				Username: "dany_ev3",
				Avatar:   "a0527abc7a7a3674529c6271bcc15f16",
			},
			{
				ID:       "409711680633700373",
				Username: "techmccat",
				Avatar:   "13518517e19c32a9bf8e2cc740c2015e",
			},
			{
				ID:       "148395955962511360",
				Username: "Hexa",
				Avatar:   "f5452a1008bf89035c1661ba748a94f8",
			},
			{
				ID:       "322756205024116739",
				Username: "\U0001F9FF👄\U0001F9FF",
				Avatar:   "0fd0609328558c514d4edc2574f79691",
			},
			{
				ID:       "783071008164282439",
				Username: "Michele Bolla",
				Avatar:   "0656e7420082e5adaabe2afe7afb4244",
			},
		}

		players = make([]Player, 6)
	)

	for i, u := range users {
		DownloadAvatar(&u)
		players[i] = Player{
			id:   u.ID,
			role: 0,
			name: u.Username,
		}
	}

	players[0].role = HITLER_ROLE
	players[1].role = FASCIST_ROLE
	players[2].role = FASCIST_ROLE
	players[3].role = LIBERAL_ROLE
	players[4].role = LIBERAL_ROLE
	players[5].role = LIBERAL_ROLE

	g := Game{
		players: players,
	}

	g.DrawStatus(&g.players[0]).SavePNG("./temp/status.png")
}
