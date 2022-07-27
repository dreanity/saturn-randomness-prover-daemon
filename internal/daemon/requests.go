package daemon

import (
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/drand"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/saturn"
	saturntypes "github.com/dreanity/saturn/x/randomness/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RandomnessChan struct {
	Randomness    *[]saturntypes.UnprovenRandomness
	PaginationKey []byte
}

func getBaseAccount(grpcConn *grpc.ClientConn, address string) *authtypes.BaseAccount {
	baseAccountChan := make(chan *authtypes.BaseAccount)

	go func() {
		baseAccount, err := saturn.GetBaseAccount(grpcConn, address)
		if err != nil {
			log.Error(err)
			baseAccountChan <- nil
			return
		}

		baseAccountChan <- baseAccount
	}()

	select {
	case baseAccount := <-baseAccountChan:
		return baseAccount
	case <-time.After(2 * time.Second):
		log.Warn("The base account request time has expired")
		return nil
	}
}

func getRounds(
	randomnesses *[]saturntypes.UnprovenRandomness,
	urls []string,
) []drand.Round {
	var rounds []drand.Round
	for _, randomness := range *randomnesses {
		roundChan := make(chan *drand.Round)

		go getRound(roundChan, urls, randomness.Round)

		select {
		case round := <-roundChan:
			if round != nil {
				rounds = append(rounds, *round)
			}
		case <-time.After(2 * time.Second):
			log.Warnf("The round №%d request time has expired", randomness.Round)
			continue
		}
	}

	return rounds
}

func getRound(c chan *drand.Round, urls []string, rRound uint64) {
	round, err := drand.GetRound(urls, rRound)
	if err != nil {
		log.Errorf("Get round №%d error: %s", rRound, err)
		c <- nil
		return
	}
	c <- round
}

func getUnprovenRandomnessAll(grpcConn *grpc.ClientConn, paginationKey []byte) (*[]saturntypes.UnprovenRandomness, []byte) {
	randomnessesChan := make(chan *RandomnessChan)

	go func() {
		randomnesses, pgk, err := saturn.GetUnprovenRandomnessAll(grpcConn, paginationKey)
		if err != nil {
			log.Errorf("Get unproven Randomness all error: %s", err)
			randomnessesChan <- nil
			return
		}
		randomnessesChan <- &RandomnessChan{
			Randomness:    randomnesses,
			PaginationKey: pgk,
		}
	}()

	select {
	case randomness := <-randomnessesChan:
		return randomness.Randomness, randomness.PaginationKey
	case <-time.After(2 * time.Second):
		log.Warn("The unproven randomness all request time has expired")
		return nil, nil
	}
}
