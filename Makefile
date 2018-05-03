list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs
deps:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go
	go get github.com/mattn/goveralls
	go get -u github.com/kardianos/govendor
install:
	gx install
	go get ./...
deps_mobile:
	go get golang.org/x/mobile/cmd/gomobile
	gomobile clean
	gomobile init
deps_hack:
	gx-go rw
deps_hack_revert:
	gx-go uw
ios:
	gomobile bind -target ios -o build/panthalassa.framework -v github.com/Bit-Nation/panthalassa/mobile
android:
	gomobile bind -target android -o build/panthalassa.aar -v github.com/Bit-Nation/panthalassa/mobile
build:
	go build -o build/panthalassa
test:
	go fmt ./...
	go test ./...
test_coverage:
	go fmt ./...
	go test ./... -coverprofile=c.out && go tool cover -html=c.out
coveralls:
	goveralls -repotoken ${COVERALS_TOKEN}