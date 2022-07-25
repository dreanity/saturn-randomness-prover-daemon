package daemon

import (
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
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

	return nil
}
