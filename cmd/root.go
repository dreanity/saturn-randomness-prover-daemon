package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/dreanity/saturn-daemon/internal/daemon"
	"github.com/spf13/cobra"
)

const (
	PrivateKey  = "private-key"
	NodeGrpcUrl = "node-grpc-url"
	DrandUrls   = "drand-urls"
	ChainID     = "chain-id"
)

func InitCmd() {
	rootCmd := &cobra.Command{
		Use:   "start",
		Short: "Start saturn daemon and set configs",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			pk, err := cmd.Flags().GetString(PrivateKey)
			if err != nil {
				return err
			}

			ngu, err := cmd.Flags().GetString(NodeGrpcUrl)
			if err != nil {
				return err
			}

			du, err := cmd.Flags().GetStringArray(DrandUrls)
			if err != nil {
				return err
			}

			cid, err := cmd.Flags().GetString(ChainID)
			if err != nil {
				return err
			}

			pkBytes, err := hex.DecodeString(pk)
			if err != nil {
				return err
			}

			privateKey := secp256k1.PrivKey{Key: pkBytes}
			pubKey := privateKey.PubKey()
			accAddress, err := types.AccAddressFromHex(pubKey.Address().String())
			if err != nil {
				return err
			}

			cfg := daemon.Configs{
				PrivateKey:  privateKey,
				PublicKey:   pubKey,
				NodeGrpcUrl: ngu,
				DrandUrls:   du,
				ChainID:     cid,
				Address:     accAddress,
			}

			if err = daemon.StartDaemon(&cfg); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().StringP(PrivateKey, "pk", "", "The private key from which the transaction will be sent (required)")
	rootCmd.MarkFlagRequired(PrivateKey)
	rootCmd.Flags().StringP(NodeGrpcUrl, "ngu", "", "A grpc url to the node to which the transaction will be sent (required)")
	rootCmd.MarkFlagRequired(NodeGrpcUrl)
	rootCmd.Flags().StringP(DrandUrls, "du", "", "Urls to drand nodes (required)")
	rootCmd.MarkFlagRequired(DrandUrls)
	rootCmd.Flags().StringP(ChainID, "cid", "", "Chain identifier (required)")
	rootCmd.MarkFlagRequired(ChainID)

	execute(rootCmd)
}

func execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
