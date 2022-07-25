package saturn

import (
	"context"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"google.golang.org/grpc"
)

func GetBaseAccount(grpcConn *grpc.ClientConn, address string) (*types.BaseAccount, error) {
	queryClient := types.NewQueryClient(grpcConn)
	res, err := queryClient.Account(context.Background(), &types.QueryAccountRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	acc := types.BaseAccount{}

	if err := types.ModuleCdc.UnpackAny(res.Account, &acc); err != nil {
		return nil, err
	}

	return &acc, nil
}
