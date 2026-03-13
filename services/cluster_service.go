package services

import (
	"context"
	pb "llm-ops/api/nexus/api/nexus-api-server/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClusterService struct {
	url    string
	token  string
	conn   *grpc.ClientConn
	client pb.NexusApiClient
}

func NewClusterService(url, token string) (*ClusterService, error) {
	// 创建 gRPC 连接
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ClusterService{
		url:    url,
		token:  token,
		conn:   conn,
		client: pb.NewNexusApiClient(conn),
	}, nil
}

func (s *ClusterService) ListNodes(ctx context.Context) ([]*pb.Node, error) {
	// 添加认证头
	md := metadata.Pairs("authorization", "Bearer "+s.token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// 调用 gRPC 方法
	resp, err := s.client.ListNodes(ctx, &pb.ListNodesRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}

func (s *ClusterService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
