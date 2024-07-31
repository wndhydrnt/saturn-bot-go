PROTOCOL_VERSION?=v0.10.0

clean:
	rm protocol/v1/saturnbot.proto

generate: protocol/v1/saturnbot.proto
	buf generate

protocol/v1/saturnbot.proto:
	curl -L --silent --fail -o ./protocol/v1/saturnbot.proto https://raw.githubusercontent.com/wndhydrnt/saturn-bot-protocol/$(PROTOCOL_VERSION)/protocol/v1/saturnbot.proto

test_cover:
	go test -covermode=set -coverpkg=./... -coverprofile cover.out -v ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out
