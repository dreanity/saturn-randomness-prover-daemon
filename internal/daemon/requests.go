package daemon

import (
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/drand"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/saturn"
	randtypes "github.com/dreanity/saturn/x/randomness/types"
	saturntypes "github.com/dreanity/saturn/x/randomness/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RandomnessChan struct {
	Randomness    *[]saturntypes.UnprovenRandomness
	PaginationKey []byte
}

func getAccount(grpcConn *grpc.ClientConn, address string) *authtypes.GenesisAccount {
	accountChan := make(chan *authtypes.GenesisAccount)

	go func() {
		account, err := saturn.GetAccount(grpcConn, address)

		if err != nil {
			log.Error(err)
			accountChan <- nil
			return
		}

		accountChan <- account
	}()

	select {
	case account := <-accountChan:
		return account
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
	randomnessesChan := make(chan *randtypes.QueryAllUnprovenRandomnessResponse)

	go func() {
		res, err := saturn.GetUnprovenRandomnessAll(grpcConn, paginationKey)
		log.Infoln(res, "ALL")
		if err != nil {
			log.Errorf("Get unproven Randomness all error: %s", err)
			randomnessesChan <- nil
			return
		}
		randomnessesChan <- res
	}()

	select {
	case randomness := <-randomnessesChan:
		if randomness != nil {
			return &randomness.UnprovenRandomness, randomness.Pagination.NextKey
		}

		return nil, nil
	case <-time.After(2 * time.Second):
		log.Warn("The unproven randomness all request time has expired")
		return nil, nil
	}
}
