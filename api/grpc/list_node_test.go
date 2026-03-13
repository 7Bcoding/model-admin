//  add unit test for list nodes

package grpc

import (
	"context"
	nexus_api_server_v1 "llm-ops/api/nexus/api/nexus-api-server/v1"
	"llm-ops/config"
	"testing"
)

func TestListNodes(t *testing.T) {
	//  use us-ca-01 cluster
	cluster := config.NexusClusters["us-ca-01"]
	conn := NewApiServerClient(&RemoteServerOptions{
		Address:  cluster.GrpcURL,
		Token:    cluster.Token,
		Insecure: false,
	})
	client, err := conn.Client()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.ListNodes(context.Background(), &nexus_api_server_v1.ListNodesRequest{})
	if err != nil {
		t.Fatalf("failed to list nodes: %v", err)
	}

	t.Logf("list nodes response: %v", resp)
}
