list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs
deps:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go
	go get golang.org/x/mobile/cmd/gomobile
	gx install
	gomobile clean
	gomobile init
ios:
	make deps
	gomobile bind -target ios -o build/panthalassa.framework
android:
	make deps
	gomobile bind -target android -o build/panthalassa.aar
build:
	make deps
	go build -o build/panthalassa
tests:
	go test ./...