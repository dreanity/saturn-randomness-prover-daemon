package daemon

import (
	"fmt"
	"time"

	"github.com/dreanity/saturn-daemon/internal/drand"
	"github.com/dreanity/saturn-daemon/internal/saturn"
	saturntypes "github.com/dreanity/saturn/x/randomness/types"
	"google.golang.org/grpc"
)

type RandomnessChan struct {
	Randomness    *[]saturntypes.UnprovenRandomness
	PaginationKey []byte
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
			fmt.Printf("The round №%d request time has expired", randomness.Round)
			continue
		}
	}

	return rounds
}

func getRound(c chan *drand.Round, urls []string, rRound uint64) {
	round, err := drand.GetRound(urls, rRound)
	if err != nil {
		fmt.Printf("Get round №%d error: %s", rRound, err)
		c <- nil
	}
	c <- round
}

func getUnprovenRandomnessAll(grpcConn *grpc.ClientConn, paginationKey []byte) (*[]saturntypes.UnprovenRandomness, []byte) {
	randomnessesChan := make(chan *RandomnessChan)

	go func() {
		randomnesses, pgk, err := saturn.GetUnprovenRandomnessAll(grpcConn, paginationKey)
		if err != nil {
			fmt.Printf("Get unproven Randomness all error: %s", err)
			randomnessesChan <- nil
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
		fmt.Printf("The unproven randomness all request time has expired")
		return nil, nil
	}

}
