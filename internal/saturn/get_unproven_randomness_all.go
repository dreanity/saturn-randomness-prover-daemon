package saturn

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dreanity/saturn/x/randomness/types"
	"google.golang.org/grpc"
)

func GetUnprovenRandomnessAll(grpcConn *grpc.ClientConn, paginationKey []byte) (*types.QueryAllUnprovenRandomnessResponse, error) {
	queryClient := types.NewQueryClient(grpcConn)
	res, err := queryClient.UnprovenRandomnessAll(context.Background(), &types.QueryAllUnprovenRandomnessRequest{
		Pagination: &query.PageRequest{
			Key: paginationKey,
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
