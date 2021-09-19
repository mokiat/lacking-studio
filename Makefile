.PHONY: run
run:
	go run ./cmd/lacking

.PHONY: install
install:
	go install github.com/mokiat/lacking-studio/cmd/lacking
