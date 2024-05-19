// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package saturnbot

import (
	"fmt"
	"os"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/wndhydrnt/saturn-bot-go/plugin"
	protocolv1 "github.com/wndhydrnt/saturn-bot-go/protocol/v1"
)

type Plugin interface {
	Apply(ctx Context) error
	Filter(ctx Context) (bool, error)
	Init(config map[string]string) error
	Name() string
	OnPrClosed(ctx Context) error
	OnPrCreated(ctx Context) error
	OnPrMerged(ctx Context) error
	Priority() int32
}

type BasePlugin struct{}

func (p BasePlugin) Apply(ctx Context) error {
	return nil
}

func (p BasePlugin) Filter(ctx Context) (bool, error) {
	return true, nil
}

func (p BasePlugin) Init(config map[string]string) error {
	return nil
}

func (p BasePlugin) OnPrClosed(ctx Context) error {
	return nil
}

func (p BasePlugin) OnPrCreated(ctx Context) error {
	return nil
}

func (p BasePlugin) OnPrMerged(ctx Context) error {
	return nil
}

func (p BasePlugin) Priority() int32 {
	return 0
}

type Context struct {
	*protocolv1.Context
	TemplateVars map[string]string
}

func newContext(c *protocolv1.Context) Context {
	return Context{
		Context:      c,
		TemplateVars: make(map[string]string),
	}
}

type provider struct {
	plugin Plugin
}

func (p *provider) ExecuteActions(req *protocolv1.ExecuteActionsRequest) (*protocolv1.ExecuteActionsResponse, error) {
	ctx := newContext(req.GetContext())
	err := inDirectory(req.Path, func() error {
		return p.plugin.Apply(ctx)
	})
	if err != nil {
		return &protocolv1.ExecuteActionsResponse{Error: fmtErr("apply of task failed: %s", err.Error())}, nil
	}

	return &protocolv1.ExecuteActionsResponse{PluginData: ctx.PluginData, TemplateVars: ctx.TemplateVars}, nil
}

func (p *provider) ExecuteFilters(req *protocolv1.ExecuteFiltersRequest) (*protocolv1.ExecuteFiltersResponse, error) {
	resp := &protocolv1.ExecuteFiltersResponse{}
	ctx := newContext(req.GetContext())
	match, err := p.plugin.Filter(ctx)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	resp.Match = match
	resp.PluginData = ctx.PluginData
	resp.TemplateVars = ctx.TemplateVars
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
	ctx := newContext(req.GetContext())
	err := p.plugin.OnPrClosed(ctx)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrCreated(req *protocolv1.OnPrCreatedRequest) (*protocolv1.OnPrCreatedResponse, error) {
	resp := &protocolv1.OnPrCreatedResponse{}
	ctx := newContext(req.GetContext())
	err := p.plugin.OnPrCreated(ctx)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrMerged(req *protocolv1.OnPrMergedRequest) (*protocolv1.OnPrMergedResponse, error) {
	resp := &protocolv1.OnPrMergedResponse{}
	ctx := newContext(req.GetContext())
	err := p.plugin.OnPrMerged(ctx)
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
