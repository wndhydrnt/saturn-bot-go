// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"

	"github.com/hashicorp/go-plugin"
	proto "github.com/wndhydrnt/saturn-bot-go/protocol/v1"
	"google.golang.org/grpc"
)

const (
	ID = "saturn-bot-plugin"
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "SATURN_SYNC_MAGIC_COOKIE",
	MagicCookieValue: "9P59IdZaEoZpENXXY2SHuvjczxUVHJaVhGG8RgeIVXfPx6c5wt34g6NLtRNehFT6",
}

var PluginMap = map[string]plugin.Plugin{
	ID: &ProviderPlugin{},
}

// Provider defines the methods to call remote code via go-plugin.
type Provider interface {
	ExecuteActions(*proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error)
	ExecuteFilters(*proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error)
	GetPlugin(*proto.GetPluginRequest) (*proto.GetPluginResponse, error)
	OnPrClosed(*proto.OnPrClosedRequest) (*proto.OnPrClosedResponse, error)
	OnPrCreated(*proto.OnPrCreatedRequest) (*proto.OnPrCreatedResponse, error)
	OnPrMerged(*proto.OnPrMergedRequest) (*proto.OnPrMergedResponse, error)
	Shutdown(*proto.ShutdownRequest) (*proto.ShutdownResponse, error)
}

// ProviderPlugin is the bridge between custom code and go-plugin.
type ProviderPlugin struct {
	plugin.Plugin
	Impl Provider
}

// GRPCServer implements GRPCPlugin.
// https://github.com/hashicorp/go-plugin/blob/8d2aaa458971cba97c3bfec1b0380322e024b514/plugin.go#L36C6-L36C16
func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServiceServer(s, &ProviderGrpcServer{Impl: p.Impl})
	return nil
}

// GRPCClient implements GRPCPlugin.
// https://github.com/hashicorp/go-plugin/blob/8d2aaa458971cba97c3bfec1b0380322e024b514/plugin.go#L36C6-L36C16
func (p *ProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &ProviderGrpcClient{client: proto.NewPluginServiceClient(c)}, nil
}
