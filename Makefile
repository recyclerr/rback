rback_version := 0.4.0

.PHONY: build clean

build :
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./release/linux_rback .
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./release/macos_rback .

clean :
	@rm ./release/*