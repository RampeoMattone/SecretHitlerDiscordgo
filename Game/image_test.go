package Game

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"sync"
	"testing"
)

func init() {
	_ = os.Mkdir("./temp", 0755)
	_ = os.Mkdir("./avatars", 0755)
}

func TestDrawFascistBoard(t *testing.T) {
	g := Game{
		game: game{
			fascistBoard: 6,
		},
		lock: sync.RWMutex{},
	}

	g.DrawFascistBoard().SavePNG("./temp/fascist-3.png")
}

func TestDrawLiberalBoard(t *testing.T) {
	g := Game{
		game: game{
			electionTracker: 3,
			liberalBoard:    5,
		},
		lock: sync.RWMutex{},
	}

	g.DrawLiberalBoard().SavePNG("./temp/liberal-3-2.png")
}

func TestDrawStatus(t *testing.T) {
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
			{
				ID:       "271001798473416704",
				Username: "Xx_DNS_xX",
				Avatar:   "e77e1dbc885595545f47c74bdda6dec0",
			},
			{
				ID:       "145874051032678400",
				Username: "slashtube",
				Avatar:   "ee2bc862adc078bd5814ba4bbb2d96f5",
			},
		}

		players = make([]Player, 8)
	)

	for i, u := range users {
		DownloadAvatar(&u)
		players[i] = Player{
			id:   u.ID,
			role: 0,
			name: u.Username,
		}
	}

	players[0].role = FASCIST_ROLE
	players[1].role = FASCIST_ROLE
	players[2].role = HITLER_ROLE
	players[3].role = LIBERAL_ROLE
	players[4].role = LIBERAL_ROLE
	players[5].role = FASCIST_ROLE
	players[6].role = LIBERAL_ROLE
	players[7].role = FASCIST_ROLE

	g := Game{
		game: game{
			players:    players,
			chancellor: &players[0],
			president:  &players[5],
		},
		lock: sync.RWMutex{},
	}

	g.DrawStatus(&g.players[2]).SavePNG("./temp/statusHitler.png")
	g.DrawStatus(&g.players[0]).SavePNG("./temp/statusFascist.png")
	g.DrawStatus(&g.players[3]).SavePNG("./temp/statusLiberal.png")
}
