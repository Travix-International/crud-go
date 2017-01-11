GITHUB_API_TOKEN := ""
VERSION :=""

build: install run-tests

install:
	go get -u github.com/FiloSottile/gvt
	gvt restore

run-tests:
	go test -cover -v

cover:
	go test -coverprofile=cover.tmp && go tool cover -html=cover.tmp
