package main

import (
	"github.com/saturn-randomness-prover-daemon/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	cmd.InitCmd()
}
