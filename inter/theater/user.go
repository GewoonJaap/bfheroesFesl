package theater

import (
	"fmt"

	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"
	"github.com/Synaxis/bfheroesFesl/storage/level"

	"github.com/sirupsen/logrus"
)

type answerUSER struct {
	TheaterID string `fesl:"TID"`
	Name      string `fesl:"NAME"` // ServerName / ClientName
	ClientID  string `fesl:"CID"`  // ?
}

func (tm *Theater) NewState(ident string) *level.State {
	return tm.level.NewState(ident)
}

// USER - SHARED Called to get user data about client? No idea
func (tm *Theater) USER(event network.ClientCommand) {
	if !event.Client.IsActive {
		logrus.Println("Client left")
		return
	}

	lkeyRedis := tm.level.NewObject("lkeys", Command.Message["LKEY"])

	redisState := tm.NewState(fmt.Sprintf(
		"%s:%s",
		"mm",
		Command.Message["LKEY"],
	))
	Client.HashState = redisState

	redisState.Set("id", lkeyRedis.Get("id"))
	redisState.Set("userID", lkeyRedis.Get("userID"))
	redisState.Set("name", lkeyRedis.Get("name"))

	Client.WriteEncode(&codec.Packet{
		Type: thtrUSER,
		Payload: answerUSER{
			TheaterID: Command.Message["TID"],
			Name:      lkeyRedis.Get("name"),
		},
	})
}
