package theater

import (
	"fmt"
	"net"

	"github.com/Synaxis/bfheroesFesl/inter/matchmaking"
	"github.com/Synaxis/bfheroesFesl/inter/network"
	"github.com/Synaxis/bfheroesFesl/inter/network/codec"

	"github.com/sirupsen/logrus"
)

type ansCGAM struct {
	TheaterID  string `fesl:"TID"`
	LobbyID    string `fesl:"LID"`
	MaxPlayers string `fesl:"MAX-PLAYERS"`
	EKEY       string `fesl:"EKEY"`
	UGID       string `fesl:"UGID"`
	Secret     string `fesl:"SECRET"`
	JOIN       string `fesl:"JOIN"`
	J          string `fesl:"J"`
	GameID     string `fesl:"GID"`
}

// CGAM - SERVER called to create a game
func (tm *Theater) CGAM(event network.EventClientCommand) {
	addr, ok := event.Client.IpAddr.(*net.TCPAddr)
	if !ok {
		logrus.Errorln("Failed turning IpAddr to net.TCPAddr")
		return
	}

	res, err := tm.db.stmtCreateServer.Exec(
		event.Command.Message["NAME"],
		event.Command.Message["B-U-community_name"],
		event.Command.Message["INT-IP"],
		event.Command.Message["INT-PORT"],
		event.Command.Message["B-version"],
	)
	if err != nil {
		logrus.Error("Cannot create New server", err)
		return
	}

	id, _ := res.LastInsertId()
	gameID := fmt.Sprintf("%d", id)

	// Store our server for easy access later
	matchmaking.Games[gameID] = event.Client

	var args []interface{}

	// Setup a new key for our game
	gameServer := tm.level.NewObject("gdata", gameID)

	keys := 0

	// Stores what we know about this game in the redis db
	for index, value := range event.Command.Message {
		if index == "TID" {
			continue
		}

		keys++

		// Strip quotes
		if len(value) > 0 && value[0] == '"' {
			value = value[1:]
		}
		if len(value) > 0 && value[len(value)-1] == '"' {
			value = value[:len(value)-1]
		}
		gameServer.Set(index, value)

		args = append(args, gameID)
		args = append(args, index)
		args = append(args, value)
	}

	gameServer.Set("LID", "1")
	gameServer.Set("GID", gameID)
	gameServer.Set("IP", addr.IP.String())
	gameServer.Set("AP", "0")
	gameServer.Set("QUEUE-LENGTH", "0")

	event.Client.HashState.Set("gdata:GID", gameID)

	_, err = tm.db.setServerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Error("Failed setting stats for game server "+gameID, err.Error())
		return
	}

	event.Client.WriteEncode(&codec.Packet{
		Type: thtrCGAM,
		Payload: ansCGAM{
			TheaterID:  event.Command.Message["TID"],
			LobbyID:    "1",
			UGID:       event.Command.Message["UGID"],
			MaxPlayers: event.Command.Message["MAX-PLAYERS"],
			EKEY:       `O65zZ2D2A58mNrZw1hmuJw%3d%3d`,
			Secret:     `2587913`,
			JOIN:       event.Command.Message["JOIN"],
			J:          event.Command.Message["JOIN"],
			GameID:     gameID,
		},
	})

	// Create game in database
	_, err = tm.db.stmtAddGame.Exec(gameID, addr.IP.String(), event.Command.Message["PORT"], event.Command.Message["B-version"], event.Command.Message["JOIN"], event.Command.Message["B-U-map"], 0, 0, event.Command.Message["MAX-PLAYERS"], 0, 0, "")
	if err != nil {
		logrus.Errorf("Failed to add game: %v", err)
	}
}
