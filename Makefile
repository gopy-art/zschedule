ifeq ($(OS),Windows_NT)
  EXECUTABLE_EXTENSION := .exe
else
  EXECUTABLE_EXTENSION :=
endif

.PHONY: build clean build-all gofmt

gofmt:
	goimports -w -l $(GO_FILES)

build:
	cd ./bin && go build -o zschedule
	cd ../

clean:
	cd ./bin/
	rm -f zschedule
	go clean
	cd ../