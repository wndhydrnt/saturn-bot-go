// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"

	proto "github.com/wndhydrnt/saturn-bot-go/protocol/v1"
)

type ProviderGrpcServer struct {
	proto.UnimplementedPluginServiceServer
	Impl Provider
}

func (s *ProviderGrpcServer) ExecuteActions(ctx context.Context, request *proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error) {
	return s.Impl.ExecuteActions(request)
}

func (s *ProviderGrpcServer) ExecuteFilters(ctx context.Context, request *proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error) {
	return s.Impl.ExecuteFilters(request)
}

func (s *ProviderGrpcServer) GetPlugin(ctx context.Context, request *proto.GetPluginRequest) (*proto.GetPluginResponse, error) {
	return s.Impl.GetPlugin(request)
}

func (s *ProviderGrpcServer) OnPrClosed(_ context.Context, request *proto.OnPrClosedRequest) (*proto.OnPrClosedResponse, error) {
	return s.Impl.OnPrClosed(request)
}

func (s *ProviderGrpcServer) OnPrCreated(_ context.Context, request *proto.OnPrCreatedRequest) (*proto.OnPrCreatedResponse, error) {
	return s.Impl.OnPrCreated(request)
}

func (s *ProviderGrpcServer) OnPrMerged(_ context.Context, request *proto.OnPrMergedRequest) (*proto.OnPrMergedResponse, error) {
	return s.Impl.OnPrMerged(request)
}

func (s *ProviderGrpcServer) Shutdown(_ context.Context, request *proto.ShutdownRequest) (*proto.ShutdownResponse, error) {
	return s.Impl.Shutdown(request)
}

type ProviderGrpcClient struct {
	client proto.PluginServiceClient
}

func (c *ProviderGrpcClient) ExecuteActions(req *proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error) {
	return c.client.ExecuteActions(context.Background(), req)
}

func (c *ProviderGrpcClient) ExecuteFilters(req *proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error) {
	return c.client.ExecuteFilters(context.Background(), req)
}

func (c *ProviderGrpcClient) GetPlugin(req *proto.GetPluginRequest) (*proto.GetPluginResponse, error) {
	return c.client.GetPlugin(context.Background(), req)
}

func (c *ProviderGrpcClient) OnPrClosed(req *proto.OnPrClosedRequest) (*proto.OnPrClosedResponse, error) {
	return c.client.OnPrClosed(context.Background(), req)
}

func (c *ProviderGrpcClient) OnPrCreated(req *proto.OnPrCreatedRequest) (*proto.OnPrCreatedResponse, error) {
	return c.client.OnPrCreated(context.Background(), req)
}

func (c *ProviderGrpcClient) OnPrMerged(req *proto.OnPrMergedRequest) (*proto.OnPrMergedResponse, error) {
	return c.client.OnPrMerged(context.Background(), req)
}

func (c *ProviderGrpcClient) Shutdown(req *proto.ShutdownRequest) (*proto.ShutdownResponse, error) {
	return c.client.Shutdown(context.Background(), req)
}
