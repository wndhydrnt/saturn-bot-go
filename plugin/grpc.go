package plugin

import (
	"context"

	proto "github.com/wndhydrnt/saturn-sync-go/protocol/v1"
)

type ProviderGrpcServer struct {
	proto.UnimplementedTaskServiceServer
	Impl Provider
}

func (s *ProviderGrpcServer) ExecuteActions(ctx context.Context, request *proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error) {
	return s.Impl.ExecuteActions(request)
}

func (s *ProviderGrpcServer) ExecuteFilters(ctx context.Context, request *proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error) {
	return s.Impl.ExecuteFilters(request)
}

func (s *ProviderGrpcServer) ListTasks(ctx context.Context, request *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	return s.Impl.ListTasks(request)
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

type ProviderGrpcClient struct {
	client proto.TaskServiceClient
}

func (c *ProviderGrpcClient) ExecuteActions(req *proto.ExecuteActionsRequest) (*proto.ExecuteActionsResponse, error) {
	return c.client.ExecuteActions(context.Background(), req)
}

func (c *ProviderGrpcClient) ExecuteFilters(req *proto.ExecuteFiltersRequest) (*proto.ExecuteFiltersResponse, error) {
	return c.client.ExecuteFilters(context.Background(), req)
}

func (c *ProviderGrpcClient) ListTasks(req *proto.ListTasksRequest) (*proto.ListTasksResponse, error) {
	return c.client.ListTasks(context.Background(), req)
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
