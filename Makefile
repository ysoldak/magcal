.PHONY: test

test:
	go clean -testcache
#	go test -test.v .
	go test -test.v -run TestSearch .

bench:
	go test -benchtime=5s -run=XXX -bench=.
