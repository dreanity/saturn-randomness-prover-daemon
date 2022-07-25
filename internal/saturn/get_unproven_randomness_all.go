package saturn

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dreanity/saturn/x/randomness/types"
	"google.golang.org/grpc"
)

func GetUnprovenRandomnessAll(grpcConn *grpc.ClientConn, paginationKey []byte) (*[]types.UnprovenRandomness, []byte, error) {
	queryClient := types.NewQueryClient(grpcConn)
	res, err := queryClient.UnprovenRandomnessAll(context.Background(), &types.QueryAllUnprovenRandomnessRequest{
		Pagination: &query.PageRequest{
			Key: paginationKey,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return &res.UnprovenRandomness, res.Pagination.NextKey, nil
}
