// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package saturnbot

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocolv1 "github.com/wndhydrnt/saturn-bot-go/protocol/v1"
)

type testPlugin struct {
	configReceived       map[string]string
	filterReturn         bool
	onPrClosedCallCount  int
	onPrCreatedCallCount int
	onPrMergedCallCount  int
}

func (tt *testPlugin) Apply(ctx Context) error {
	ctx.TemplateVars["tplSource"] = "apply"
	ctx.PluginData["pdSource"] = "apply"
	f, err := os.Create("unittest.txt")
	if err != nil {
		return err
	}

	return f.Close()
}

func (tt *testPlugin) Filter(ctx Context) (bool, error) {
	ctx.TemplateVars["tplSource"] = "filter"
	ctx.PluginData["pdSource"] = "filter"
	return tt.filterReturn, nil
}

func (tt *testPlugin) Init(config map[string]string) error {
	tt.configReceived = config
	return nil
}

func (tt *testPlugin) Name() string {
	return "testPlugin"
}

func (tt *testPlugin) OnPrClosed(ctx Context) error {
	tt.onPrClosedCallCount++
	return nil
}

func (tt *testPlugin) OnPrCreated(ctx Context) error {
	tt.onPrCreatedCallCount++
	return nil
}

func (tt *testPlugin) OnPrMerged(ctx Context) error {
	tt.onPrMergedCallCount++
	return nil
}

func (tt *testPlugin) Priority() int32 {
	return 10
}

func TestProvider_ExecuteActions_ApplySucceeds(t *testing.T) {
	p1 := &testPlugin{}
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	req := &protocolv1.ExecuteActionsRequest{
		Context: &protocolv1.Context{},
		Path:    dir,
	}

	p := &provider{plugin: p1}
	resp, err := p.ExecuteActions(req)

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, map[string]string{"pdSource": "apply"}, resp.GetPluginData())
	assert.Equal(t, map[string]string{"tplSource": "apply"}, resp.GetTemplateVars())
	_, err = os.Stat(path.Join(dir, "unittest.txt"))
	require.NoError(t, err)
}

func TestProvider_ExecuteFilters_Succeed(t *testing.T) {
	p1 := &testPlugin{
		filterReturn: true,
	}
	req := &protocolv1.ExecuteFiltersRequest{
		Context: &protocolv1.Context{},
	}

	p := &provider{plugin: p1}
	resp, err := p.ExecuteFilters(req)

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, map[string]string{"pdSource": "filter"}, resp.GetPluginData())
	assert.Equal(t, map[string]string{"tplSource": "filter"}, resp.GetTemplateVars())
	assert.True(t, resp.GetMatch())
}

func TestProvider_GetPlugin(t *testing.T) {
	p1 := &testPlugin{}
	req := &protocolv1.GetPluginRequest{Config: map[string]string{"custom": "config"}}

	p := &provider{plugin: p1}
	resp, err := p.GetPlugin(req)

	require.NoError(t, err)
	assert.Equal(t, p1.Name(), resp.GetName())
	assert.Equal(t, p1.Priority(), resp.GetPriority())
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, map[string]string{"custom": "config"}, p1.configReceived)
}

func TestProvider_OnPrClosed_Succeed(t *testing.T) {
	p1 := &testPlugin{}

	p := &provider{plugin: p1}
	resp, err := p.OnPrClosed(&protocolv1.OnPrClosedRequest{})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, p1.onPrClosedCallCount)
}

func TestProvider_OnPrCreated_Succeed(t *testing.T) {
	p1 := &testPlugin{}

	p := &provider{plugin: p1}
	resp, err := p.OnPrCreated(&protocolv1.OnPrCreatedRequest{})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, p1.onPrCreatedCallCount)
}

func TestProvider_OnPrMerged_Succeed(t *testing.T) {
	p1 := &testPlugin{}

	p := &provider{plugin: p1}
	resp, err := p.OnPrMerged(&protocolv1.OnPrMergedRequest{})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, p1.onPrMergedCallCount)
}
