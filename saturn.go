// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package saturn_sync

import (
	"errors"
	"fmt"
	"os"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/wndhydrnt/saturn-sync-go/plugin"
	protocolv1 "github.com/wndhydrnt/saturn-sync-go/protocol/v1"
)

var (
	errTaskNotFound = errors.New("task not found")
)

type Task struct {
	Name                  string
	AutoMerge             bool
	AutoMergeAfterSeconds int
	BranchName            string
	ChangeLimit           int
	CommitMessage         string
	CreateOnly            bool
	Disabled              bool
	Filters               *protocolv1.Filters
	KeepBranchAfterMerge  bool
	Labels                []string
	MergeOnce             bool
	PrBody                string
	PrTitle               string
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetAutoMerge() bool {
	return t.AutoMerge
}

func (t *Task) GetAutoMergeAfterSeconds() int {
	return t.AutoMergeAfterSeconds
}

func (t *Task) GetBranchName() string {
	return t.BranchName
}

func (t *Task) GetChangeLimit() int {
	return t.ChangeLimit
}

func (t *Task) GetCreateOnly() bool {
	return t.CreateOnly
}

func (t *Task) GetCommitMessage() string {
	return t.CommitMessage
}

func (t *Task) GetDisabled() bool {
	return t.Disabled
}

func (t *Task) GetFilters() *protocolv1.Filters {
	return t.Filters
}

func (t *Task) GetKeepBranchAfterMerge() bool {
	return t.KeepBranchAfterMerge
}

func (t *Task) GetLabels() []string {
	return t.Labels
}

func (t *Task) GetMergeOnce() bool {
	return t.MergeOnce
}

func (t *Task) GetPrBody() string {
	return t.PrBody
}

func (t *Task) GetPrTitle() string {
	return t.PrTitle
}

func (t *Task) Init(customConfig []byte) error {
	return nil
}

func (t *Task) OnPrClosed(ctx *protocolv1.Context) error {
	return nil
}

func (t *Task) OnPrCreated(ctx *protocolv1.Context) error {
	return nil
}

func (t *Task) OnPrMerged(ctx *protocolv1.Context) error {
	return nil
}

type Tasker interface {
	Apply(ctx *protocolv1.Context) error
	Filter(ctx *protocolv1.Context) (bool, error)
	GetName() string
	GetAutoMerge() bool
	GetAutoMergeAfterSeconds() int
	GetBranchName() string
	GetChangeLimit() int
	GetCommitMessage() string
	GetCreateOnly() bool
	GetDisabled() bool
	GetFilters() *protocolv1.Filters
	GetKeepBranchAfterMerge() bool
	GetLabels() []string
	GetMergeOnce() bool
	GetPrBody() string
	GetPrTitle() string
	Init(customConfig []byte) error
	OnPrClosed(ctx *protocolv1.Context) error
	OnPrCreated(ctx *protocolv1.Context) error
	OnPrMerged(ctx *protocolv1.Context) error
}

type provider struct {
	tasks []Tasker
}

func (p *provider) ExecuteActions(req *protocolv1.ExecuteActionsRequest) (*protocolv1.ExecuteActionsResponse, error) {
	task, err := p.findTask(req.TaskName)
	if err != nil {
		return &protocolv1.ExecuteActionsResponse{Error: Ptr(err.Error())}, nil
	}

	err = inDirectory(req.Path, func() error {
		return task.Apply(req.GetContext())
	})
	if err != nil {
		return &protocolv1.ExecuteActionsResponse{Error: fmtErr("apply of task failed: %s", err.Error())}, nil
	}

	return &protocolv1.ExecuteActionsResponse{}, nil
}

func (p *provider) ExecuteFilters(req *protocolv1.ExecuteFiltersRequest) (*protocolv1.ExecuteFiltersResponse, error) {
	resp := &protocolv1.ExecuteFiltersResponse{}
	task, err := p.findTask(req.TaskName)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	match, err := task.Filter(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	resp.Match = match
	return resp, nil
}

func (p *provider) ListTasks(req *protocolv1.ListTasksRequest) (*protocolv1.ListTasksResponse, error) {
	resp := &protocolv1.ListTasksResponse{}
	for _, t := range p.tasks {
		err := t.Init(req.CustomConfig)
		if err != nil {
			resp.Error = fmtErr("init of task %s failed: %s", t.GetName(), err.Error())
			return resp, nil
		}

		resp.Tasks = append(resp.Tasks, toProtoTask(t))
	}

	return resp, nil
}

func (p *provider) OnPrClosed(req *protocolv1.OnPrClosedRequest) (*protocolv1.OnPrClosedResponse, error) {
	resp := &protocolv1.OnPrClosedResponse{}
	task, err := p.findTask(req.TaskName)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	err = task.OnPrClosed(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrCreated(req *protocolv1.OnPrCreatedRequest) (*protocolv1.OnPrCreatedResponse, error) {
	resp := &protocolv1.OnPrCreatedResponse{}
	task, err := p.findTask(req.TaskName)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	err = task.OnPrCreated(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) OnPrMerged(req *protocolv1.OnPrMergedRequest) (*protocolv1.OnPrMergedResponse, error) {
	resp := &protocolv1.OnPrMergedResponse{}
	task, err := p.findTask(req.TaskName)
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	err = task.OnPrMerged(req.GetContext())
	if err != nil {
		resp.Error = Ptr(err.Error())
		return resp, nil
	}

	return resp, nil
}

func (p *provider) findTask(name string) (Tasker, error) {
	for _, t := range p.tasks {
		if t.GetName() == name {
			return t, nil
		}
	}

	return nil, errTaskNotFound
}

func ServeTask(task ...Tasker) {
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: goPlugin.PluginSet{
			"tasks": &plugin.ProviderPlugin{Impl: &provider{tasks: task}},
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}

func toProtoTask(t Tasker) *protocolv1.Task {
	return &protocolv1.Task{
		Name:                  t.GetName(),
		Actions:               []*protocolv1.Action{},
		AutoMerge:             Ptr(t.GetAutoMerge()),
		AutoMergeAfterSeconds: Ptr(int32(t.GetAutoMergeAfterSeconds())),
		BranchName:            Ptr(t.GetBranchName()),
		ChangeLimit:           Ptr(int32(t.GetChangeLimit())),
		CommitMessage:         Ptr(t.GetCommitMessage()),
		CreateOnly:            Ptr(t.GetCreateOnly()),
		Disabled:              Ptr(t.GetDisabled()),
		Filters:               t.GetFilters(),
		KeepBranchAfterMerge:  Ptr(t.GetKeepBranchAfterMerge()),
		Labels:                t.GetLabels(),
		MergeOnce:             Ptr(t.GetMergeOnce()),
		PrBody:                Ptr(t.GetPrBody()),
		PrTitle:               Ptr(t.GetPrTitle()),
	}
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
