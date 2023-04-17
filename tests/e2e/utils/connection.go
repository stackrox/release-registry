// Package utils contains end-to-end utils.
package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	minPort = 30000
	maxPort = 39999
)

// RemotePort is the port of the mock server for e2e tests.
//
//nolint:gochecknoglobals
var RemotePort int

func getRandomIntInRange(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max-min+1) + min
}

// bearerToken implements the credentials.PerRPCCredentials interface, and sets
// a bearer token on the connection metadata.
type bearerToken string

var _ credentials.PerRPCCredentials = (*bearerToken)(nil)

func (t bearerToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + string(t),
	}, nil
}

func (t bearerToken) RequireTransportSecurity() bool {
	return true
}

func getRemote(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}

// GetGRPCConnection creates an authenticated GRPC client connection to a remote.
func GetGRPCConnection(ctx context.Context, remotePort int, token string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, getRemote(remotePort),
		grpc.WithPerRPCCredentials(bearerToken(token)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})),
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create grpc connection")
	}

	return conn, nil
}
