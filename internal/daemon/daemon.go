package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/saturn"
	"google.golang.org/grpc"
)

type Configs struct {
	PrivateKey  secp256k1.PrivKey
	PublicKey   cryptotypes.PubKey
	NodeGrpcUrl string
	DrandUrls   []string
	ChainID     string
	Address     types.AccAddress
}

func StartDaemon(configs *Configs) (err error) {
	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		configs.NodeGrpcUrl, // Or your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	go func() {
		var paginationKey []byte

		baseAccount := getBaseAccount(grpcConn, configs.Address.String())

		if baseAccount == nil {
			log.Fatal()
		}

		for {
			randomnesses, pgk := getUnprovenRandomnessAll(grpcConn, paginationKey)

			if randomnesses == nil || pgk == nil {
				continue
			}
			paginationKey = pgk

			rounds := getRounds(randomnesses, configs.DrandUrls)

			for _, round := range rounds {
				err := saturn.SendProveRandomnessMsg(
					context.Background(),
					grpcConn,
					&round,
					configs.PrivateKey,
					configs.PublicKey,
					configs.Address.String(),
					baseAccount.AccountNumber,
					baseAccount.Sequence,
					configs.ChainID,
				)

				if err != nil {
					continue
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	return nil
}
