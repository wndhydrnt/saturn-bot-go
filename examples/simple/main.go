// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"os"

	saturnbot "github.com/wndhydrnt/saturn-bot-go"
)

type Example struct {
	saturnbot.BasePlugin

	message string
}

func (e *Example) Init(config map[string]string) error {
	e.message = config["message"]
	return nil
}

func (e *Example) Filter(ctx saturnbot.Context) (bool, error) {
	// Match a single repository.
	// Implement more complex matching logic here by calling APIs.
	match := ctx.GetRepository().GetFullName() == "github.com/wndhydrnt/saturn-bot-example"
	return match, nil
}

func (e *Example) Apply(ctx saturnbot.Context) error {
	// Create a file in the root of the repository.
	f, err := os.Create("hello-go.txt")
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(e.message + "\n")
	if err != nil {
		return err
	}

	return nil
}

func (e *Example) Name() string {
	return "Example"
}

func main() {
	// Initialize and serve the plugin.
	saturnbot.ServePlugin(&Example{})
}
