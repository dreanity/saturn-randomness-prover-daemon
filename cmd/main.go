package main

import (
	"github.com/dreanity/saturn-randomness-prover-daemon/cmd/cobra"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	cobra.InitCmd()
}
