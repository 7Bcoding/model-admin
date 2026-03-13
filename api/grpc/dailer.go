package grpc

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	apiV1 "llm-ops/api/nexus/api/nexus-api-server/v1"

	"github.com/go-kratos/kratos/v2/middleware"
	kratosgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

type RemoteServerOptions struct {
	Address  string
	Insecure bool
	Token    string
}

type ApiServerClient struct {
	option    *RemoteServerOptions
	conn      *grpc.ClientConn
	apiClient apiV1.NexusApiClient
	mutex     sync.RWMutex
}

func NewApiServerClient(option *RemoteServerOptions) *ApiServerClient {
	return &ApiServerClient{
		option: option,
	}
}

func (client *ApiServerClient) Client() (apiV1.NexusApiClient, error) {
	if client.apiClient != nil {
		return client.apiClient, nil
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if client.apiClient != nil {
		return client.apiClient, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*40)
	defer cancel()
	conn, err := dial(ctx, client.option)
	if err != nil {
		return nil, err
	}

	client.apiClient = apiV1.NewNexusApiClient(conn)
	return client.apiClient, nil
}

func dial(ctx context.Context, option *RemoteServerOptions, m ...middleware.Middleware) (*grpc.ClientConn, error) {
	middleWares := make([]middleware.Middleware, 0)
	if option.Token != "" {
		middleWares = append(middleWares, Client(option.Token))
	}
	ctx, _ = context.WithTimeout(ctx, time.Second*30)
	middleWares = append(middleWares, m...)
	if !option.Insecure {
		return kratosgrpc.Dial(
			ctx,
			kratosgrpc.WithTLSConfig(&tls.Config{}),
			kratosgrpc.WithEndpoint(option.Address),
			kratosgrpc.WithTimeout(time.Second*30),
			kratosgrpc.WithOptions(grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1.0 * time.Second,
					Multiplier: 1.5,
					Jitter:     0.2,
					MaxDelay:   60 * time.Second,
				},
				MinConnectTimeout: time.Second * 30,
			})),
			kratosgrpc.WithMiddleware(
				middleWares...,
			),
		)
	} else {
		return kratosgrpc.DialInsecure(
			ctx,
			kratosgrpc.WithEndpoint(option.Address),
			kratosgrpc.WithTimeout(time.Second*30),
			kratosgrpc.WithOptions(grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1.0 * time.Second,
					Multiplier: 1.5,
					Jitter:     0.2,
					MaxDelay:   60 * time.Second,
				},
				MinConnectTimeout: time.Second * 30,
			})),
			kratosgrpc.WithMiddleware(
				middleWares...,
			),
		)
	}
}
