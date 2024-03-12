// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package saturn_sync

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocolv1 "github.com/wndhydrnt/saturn-sync-go/protocol/v1"
)

type testTask struct {
	Task
	customConfigReceived []byte
	filterReturn         bool
	onPrClosedCallCount  int
	onPrCreatedCallCount int
	onPrMergedCallCount  int
}

func (tt *testTask) Apply(ctx *protocolv1.Context) error {
	f, err := os.Create("unittest.txt")
	if err != nil {
		return err
	}

	return f.Close()
}

func (tt *testTask) Filter(ctx *protocolv1.Context) (bool, error) {
	return tt.filterReturn, nil
}

func (tt *testTask) Init(customConfig []byte) error {
	tt.customConfigReceived = customConfig
	return nil
}

func (tt *testTask) OnPrClosed(ctx *protocolv1.Context) error {
	tt.onPrClosedCallCount++
	return nil
}

func (tt *testTask) OnPrCreated(ctx *protocolv1.Context) error {
	tt.onPrCreatedCallCount++
	return nil
}

func (tt *testTask) OnPrMerged(ctx *protocolv1.Context) error {
	tt.onPrMergedCallCount++
	return nil
}

func TestProvider_ExecuteActions_ApplySucceeds(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
	}
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	req := &protocolv1.ExecuteActionsRequest{Path: dir, TaskName: "t1"}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.ExecuteActions(req)

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	_, err = os.Stat(path.Join(dir, "unittest.txt"))
	require.NoError(t, err)
}

func TestProvider_ExecuteActions_TaskNotFound(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	req := &protocolv1.ExecuteActionsRequest{Path: dir, TaskName: "t1"}

	p := &provider{tasks: []Tasker{}}
	resp, err := p.ExecuteActions(req)

	require.NoError(t, err)
	assert.Equal(t, "task not found", resp.GetError())
}

func TestProvider_ExecuteFilters_Succeed(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
		filterReturn: true,
	}
	req := &protocolv1.ExecuteFiltersRequest{TaskName: "t1"}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.ExecuteFilters(req)

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.True(t, resp.GetMatch())
}

func TestProvider_ExecuteFilters_TaskNotFound(t *testing.T) {
	req := &protocolv1.ExecuteFiltersRequest{TaskName: "t1"}

	p := &provider{tasks: []Tasker{}}
	resp, err := p.ExecuteFilters(req)

	require.NoError(t, err)
	assert.Equal(t, "task not found", resp.GetError())
}

func TestProvider_ListTasks(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
	}
	req := &protocolv1.ListTasksRequest{CustomConfig: []byte("unittest")}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.ListTasks(req)

	require.NoError(t, err)
	assert.Len(t, resp.GetTasks(), 1)
	assert.Equal(t, t1.GetName(), resp.GetTasks()[0].GetName())
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, []byte("unittest"), t1.customConfigReceived)
}

func TestProvider_OnPrClosed_Succeed(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
	}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.OnPrClosed(&protocolv1.OnPrClosedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, t1.onPrClosedCallCount)
}

func TestProvider_OnPrClosed_TaskNotFound(t *testing.T) {
	p := &provider{tasks: []Tasker{}}
	resp, err := p.OnPrClosed(&protocolv1.OnPrClosedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "task not found", resp.GetError())
}

func TestProvider_OnPrCreated_Succeed(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
	}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.OnPrCreated(&protocolv1.OnPrCreatedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, t1.onPrCreatedCallCount)
}

func TestProvider_OnPrCreated_TaskNotFound(t *testing.T) {
	p := &provider{tasks: []Tasker{}}
	resp, err := p.OnPrCreated(&protocolv1.OnPrCreatedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "task not found", resp.GetError())
}

func TestProvider_OnPrMerged_Succeed(t *testing.T) {
	t1 := &testTask{
		Task: Task{
			Name: "t1",
		},
	}

	p := &provider{tasks: []Tasker{t1}}
	resp, err := p.OnPrMerged(&protocolv1.OnPrMergedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "", resp.GetError())
	assert.Equal(t, 1, t1.onPrMergedCallCount)
}

func TestProvider_OnPrMerged_TaskNotFound(t *testing.T) {
	p := &provider{tasks: []Tasker{}}
	resp, err := p.OnPrMerged(&protocolv1.OnPrMergedRequest{TaskName: "t1"})

	require.NoError(t, err)
	assert.Equal(t, "task not found", resp.GetError())
}
