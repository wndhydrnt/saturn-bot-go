test_cover:
	go test -coverprofile cover.out -v ./...
	go tool cover -html cover.out -o cover.html
	rm cover.out
