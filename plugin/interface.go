package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	proto "github.com/wndhydrnt/saturn-sync-go/protocol/v1"
	"google.golang.org/grpc"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SATURN_SYNC_MAGIC_COOKIE",
	MagicCookieValue: "9P59IdZaEoZpENXXY2SHuvjczxUVHJaVhGG8RgeIVXfPx6c5wt34g6NLtRNehFT6",
}

var PluginMap = map[string]plugin.Plugin{
	"tasks": &ProviderPlugin{},
}

// Provider defines the methods to call remote code via go-plugin.
type Provider interface {
	ExecuteActions(*proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error)
	ExecuteFilters(*proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error)
	ListTasks(*proto.ListTasksRequest) (*proto.ListTasksResponse, error)
	OnPrClosed(*proto.OnPrClosedRequest) (*proto.OnPrClosedResponse, error)
	OnPrCreated(*proto.OnPrCreatedRequest) (*proto.OnPrCreatedResponse, error)
	OnPrMerged(*proto.OnPrMergedRequest) (*proto.OnPrMergedResponse, error)
}

// ProviderPlugin is the bridge between custom code and go-plugin.
type ProviderPlugin struct {
	plugin.Plugin
	Impl Provider
}

// GRPCServer implements GRPCPlugin.
// https://github.com/hashicorp/go-plugin/blob/8d2aaa458971cba97c3bfec1b0380322e024b514/plugin.go#L36C6-L36C16
func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterTaskServiceServer(s, &ProviderGrpcServer{Impl: p.Impl})
	return nil
}

// GRPCClient implements GRPCPlugin.
// https://github.com/hashicorp/go-plugin/blob/8d2aaa458971cba97c3bfec1b0380322e024b514/plugin.go#L36C6-L36C16
func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ProviderGrpcClient{client: proto.NewTaskServiceClient(c)}, nil
}
