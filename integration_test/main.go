// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	saturnbot "github.com/wndhydrnt/saturn-bot-go"
)

type IntegrationTest struct {
	saturnbot.BasePlugin

	eventOutTmpFilePath string
	staticContent       string
}

func (p *IntegrationTest) Apply(ctx saturnbot.Context) error {
	content := fmt.Sprintf("%s\n%s", p.staticContent, ctx.RunData["dynamic"])
	return os.WriteFile("integration-test.txt", []byte(content), 0600)
}

func (p *IntegrationTest) Filter(ctx saturnbot.Context) (bool, error) {
	if ctx.Repository.FullName == "git.localhost/integration/test" {
		return true, nil
	}

	if ctx.Repository.FullName == "git.localhost/integration/rundata" {
		ctx.RunData["plugin"] = "set by plugin"
		return true, nil
	}

	return false, nil
}

func (p *IntegrationTest) Init(config map[string]string) error {
	p.staticContent = config["content"]
	p.eventOutTmpFilePath = filepath.Join(os.TempDir(), config["event_out_tmp_file_path"])
	return nil
}

func (p *IntegrationTest) OnPrClosed(ctx saturnbot.Context) error {
	return os.WriteFile(p.eventOutTmpFilePath, []byte("Integration Test OnPrClosed"), 0644)
}

func (p *IntegrationTest) OnPrCreated(ctx saturnbot.Context) error {
	return os.WriteFile(p.eventOutTmpFilePath, []byte("Integration Test OnPrCreated"), 0644)
}

func (p *IntegrationTest) OnPrMerged(ctx saturnbot.Context) error {
	return os.WriteFile(p.eventOutTmpFilePath, []byte("Integration Test OnPrMerged"), 0644)
}

func (p *IntegrationTest) Name() string {
	return "integration-test"
}

func main() {
	saturnbot.ServePlugin(&IntegrationTest{})
}
