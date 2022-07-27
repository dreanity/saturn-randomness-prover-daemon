package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

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

	var paginationKey []byte

	baseAccount := getBaseAccount(grpcConn, configs.Address.String())

	if baseAccount == nil {
		log.Fatal("Base account is nil")
	}

	ticker := time.NewTicker(time.Second)
	stop := make(chan bool)

	go func() {
		defer func() { stop <- true }()
		for {
			select {
			case <-ticker.C:
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
			case <-stop:
				log.Info("Stopping the deamon")
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ticker.Stop()

	stop <- true

	<-stop

	return nil
}
