PROTOCOL_VERSION?=v0.11.4
INTEGRATION_TEST_BIN=integration-test-$(PROTOCOL_VERSION).$(shell uname -s)-$(shell uname -m)
INTEGRATION_TEST_BIN_PATH?=integration_test/$(INTEGRATION_TEST_BIN)
INTEGRATION_TEST_PLUGIN_PATH=integration_test/plugin-integration-test
SATURN_BOT_BIN_PATH?=saturn-bot

clean:
	rm protocol/v1/saturnbot.proto || true
	rm $(INTEGRATION_TEST_PLUGIN_PATH) || true
	rm $(INTEGRATION_TEST_BIN_PATH) || true

generate: protocol/v1/saturnbot.proto
	buf generate

protocol/v1/saturnbot.proto:
	curl -L --silent --fail -o ./protocol/v1/saturnbot.proto https://raw.githubusercontent.com/wndhydrnt/saturn-bot-protocol/$(PROTOCOL_VERSION)/protocol/v1/saturnbot.proto

test_cover:
	go test -covermode=set -coverpkg=./... -coverprofile cover.out -v ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out

$(INTEGRATION_TEST_PLUGIN_PATH):
	cd integration_test && go build

$(INTEGRATION_TEST_BIN_PATH):
	curl -fsSL -o $(INTEGRATION_TEST_BIN_PATH) "https://github.com/wndhydrnt/saturn-bot-protocol/releases/download/$(PROTOCOL_VERSION)/$(INTEGRATION_TEST_BIN)"
	chmod +x $(INTEGRATION_TEST_BIN_PATH)

.PHONY: test_integration
test_integration: $(INTEGRATION_TEST_BIN_PATH) $(INTEGRATION_TEST_PLUGIN_PATH)
	$(INTEGRATION_TEST_BIN_PATH) -plugin-path $(INTEGRATION_TEST_PLUGIN_PATH) -saturn-bot-path $(SATURN_BOT_BIN_PATH)
