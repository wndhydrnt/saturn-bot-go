// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package saturn_sync

import (
	"fmt"
	"os"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/wndhydrnt/saturn-sync-go/plugin"
	protocolv1 "github.com/wndhydrnt/saturn-sync-go/protocol/v1"
)

type Plugin interface {
	Apply(ctx *protocolv1.Context) error
	Filter(ctx *protocolv1.Context) (bool, error)
	Init(config map[string]string) error
	Name() string
	OnPrClosed(ctx *protocolv1.Context) error
	OnPrCreated(ctx *protocolv1.Context) error
	OnPrMerged(ctx *protocolv1.Context) error
	Priority() int32
}

type BasePlugin struct{}

func (p BasePlugin) Apply(ctx *protocolv1.Context) error {
	return nil
}

func (p BasePlugin) Filter(ctx *protocolv1.Context) (bool, error) {
	return true, nil
}

func (p BasePlugin) Init(config map[string]string) error {
	return nil
}

func (p BasePlugin) OnPrClosed(ctx *protocolv1.Context) error {
	return nil
}

func (p BasePlugin) OnPrCreated(ctx *protocolv1.Context) error {
	return nil
}

func (p BasePlugin) OnPrMerged(ctx *protocolv1.Context) error {
	return nil
}

func (p BasePlugin) Priority() int32 {
	return 0
}

type provider struct {
	plugin Plugin
}

func (p *provider) ExecuteActions(req *protocolv1.ExecuteActionsRequest) (*protocolv1.ExecuteActionsResponse, error) {
	err := inDirectory(req.Path, func() error {
		return p.plugin.Apply(req.GetContext())
	})
	if err != nil {
		return &protocolv1.ExecuteActionsResponse{Error: fmtErr("apply of task failed: %s", err.Error())}, nil
	}

	return &protocolv1.ExecuteActionsResponse{}, nil
}

func (p *provider) ExecuteFilters(req *protocolv1.ExecuteFiltersRequest) (*protocolv1.ExecuteFiltersResponse, error) {
	resp := &protocolv1.ExecuteFiltersResponse{}
	match, err := p.plugin.Filter(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	resp.Match = match
	return resp, nil
}

func (p *provider) GetPlugin(req *protocolv1.GetPluginRequest) (*protocolv1.GetPluginResponse, error) {
	resp := &protocolv1.GetPluginResponse{
		Name:     p.plugin.Name(),
		Priority: int32(p.plugin.Priority()),
	}
	cfg := req.GetConfig()
	if cfg == nil {
		cfg = map[string]string{}
	}

	err := p.plugin.Init(cfg)
	if err != nil {
		resp.Error = fmtErr("init of plugin %s failed: %s", p.plugin.Name(), err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrClosed(req *protocolv1.OnPrClosedRequest) (*protocolv1.OnPrClosedResponse, error) {
	resp := &protocolv1.OnPrClosedResponse{}
	err := p.plugin.OnPrClosed(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrCreated(req *protocolv1.OnPrCreatedRequest) (*protocolv1.OnPrCreatedResponse, error) {
	resp := &protocolv1.OnPrCreatedResponse{}
	err := p.plugin.OnPrCreated(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrMerged(req *protocolv1.OnPrMergedRequest) (*protocolv1.OnPrMergedResponse, error) {
	resp := &protocolv1.OnPrMergedResponse{}
	err := p.plugin.OnPrMerged(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func ServePlugin(p Plugin) {
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: goPlugin.PluginSet{
			plugin.ID: &plugin.ProviderPlugin{Impl: &provider{plugin: p}},
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}

func Ptr[T any](t T) *T {
	return &t
}

func fmtErr(format string, a ...any) *string {
	s := fmt.Sprintf(format, a...)
	return &s
}

func inDirectory(dir string, f func() error) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		return fmt.Errorf("change to work directory: %w", err)
	}

	funcErr := f()
	err = os.Chdir(currentDir)
	if err != nil {
		return fmt.Errorf("changing directory to %s at the end of applying actions: %w", currentDir, err)
	}

	return funcErr
}
