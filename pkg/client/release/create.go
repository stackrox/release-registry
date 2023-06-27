package release

import (
	"context"

	v1 "github.com/stackrox/release-registry/gen/go/proto/api/v1"
	"github.com/stackrox/release-registry/pkg/client/common"
	"google.golang.org/grpc"
)

func Execute(input string) (string, error) {
	return common.WithGRPCHandler(run)(input)
}

func run(ctx context.Context, conn *grpc.ClientConn, input string) (string, error) {
	resp, err := v1.NewReleaseServiceClient(conn).
	if err != nil {
		return "", err
	}

	return resp.Commit, nil
}
