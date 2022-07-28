package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/saturn"
	log "github.com/sirupsen/logrus"
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

	ticker := time.NewTicker(time.Second)
	stop := make(chan bool)

	go func() {
		defer func() { stop <- true }()
		for {
			select {
			case <-ticker.C:
				randomnesses, pgk := getUnprovenRandomnessAll(grpcConn, paginationKey)
				log.Infoln(randomnesses, pgk)
				if randomnesses == nil {
					log.Warnln("Randomnesses is nil or PaginationKey is nil")
					continue
				}
				paginationKey = pgk

				rounds := getRounds(randomnesses, configs.DrandUrls)
				for _, round := range rounds {
					account := getAccount(grpcConn, configs.Address.String())

					if account == nil {
						log.Warnln("Base account is nil")
						continue
					}

					err := saturn.SendProveRandomnessMsg(
						context.Background(),
						grpcConn,
						&round,
						configs.PrivateKey,
						configs.PublicKey,
						configs.Address.String(),
						(*account).GetAccountNumber(),
						(*account).GetSequence(),
						configs.ChainID,
					)

					if err != nil {
						log.Errorln(err)
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
