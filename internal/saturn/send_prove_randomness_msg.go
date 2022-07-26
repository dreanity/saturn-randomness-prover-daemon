package saturn

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client/tx"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dreanity/saturn-randomness-prover-daemon/internal/drand"
	"github.com/dreanity/saturn/app"
	"github.com/dreanity/saturn/x/randomness/types"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"google.golang.org/grpc"
)

const (
	uhydrogen = "uhydrogen"
	gasLimit  = 150_000
)

func SendProveRandomnessMsg(
	ctx context.Context,
	grpcConn *grpc.ClientConn,
	round *drand.Round,
	privKey secp256k1.PrivKey,
	pubKey cryptotypes.PubKey,
	accAddress string,
	accNum uint64,
	accSeq uint64,
	chainId string,
) error {
	encoding := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	txBuilder := encoding.TxConfig.NewTxBuilder()

	msg := types.NewMsgProveRandomness(
		accAddress,
		round.Round,
		round.Randomness,
		round.Signature,
		round.PreviousSignature,
	)

	if err := txBuilder.SetMsgs(msg); err != nil {
		return err
	}

	gasPrice := sdktypes.NewCoin(uhydrogen, sdktypes.NewInt(1))

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdktypes.NewCoins(gasPrice))

	//-----------------------------------------------------------------------

	sigV2 := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode:  encoding.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
	}

	txBuilder.SetSignatures(sigV2)

	signerData := xauthsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accNum,
		Sequence:      accSeq,
	}

	sigV2, err := tx.SignWithPrivKey(
		encoding.TxConfig.SignModeHandler().DefaultMode(),
		signerData,
		txBuilder,
		&privKey,
		encoding.TxConfig,
		accSeq)

	if err != nil {
		return err
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return err
	}

	txSender := sdktx.NewServiceClient(grpcConn)
	txBytes, err := encoding.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	res, err := txSender.BroadcastTx(
		ctx,
		&sdktx.BroadcastTxRequest{
			Mode:    sdktx.BroadcastMode_BROADCAST_MODE_ASYNC,
			TxBytes: txBytes,
		})

	if err != nil {
		return err
	}

	_ = res
	return nil
}
