package theater

import (
	"strconv"

	"github.com/Synaxis/bfheroesFesl/inter/network"

	"github.com/sirupsen/logrus"
)

func (tM *Theater) UPLA(event network.EventClientCommand) {
	if !event.Client.IsActive {
		return
	}

	var args []interface{}

	keys := 0

	pid := event.Command.Message["PID"]
	gid := event.Command.Message["GID"]

	for index, value := range event.Command.Message {
		if index == "TID" || index == "PID" || index == "GID" {
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

		args = append(args, gid)
		args = append(args, pid)
		args = append(args, index)
		args = append(args, value)
	}

	var err error
	_, err = tM.db.setServerPlayerStatsStatement(keys).Exec(args...)
	if err != nil {
		logrus.Errorln("Failed to update stats for player "+pid, err.Error())
	}

	gdata := tM.level.NewObject("gdata", event.Command.Message["GID"])

	num, _ := strconv.Atoi(gdata.Get("AP"))

	num++

	gdata.Set("AP", strconv.Itoa(num))
}
