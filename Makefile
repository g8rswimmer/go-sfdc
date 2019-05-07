.PHONY: tests
tests:
	go test ./... -cover

.PHONY: docs
docs:
	godoc -http=:6060