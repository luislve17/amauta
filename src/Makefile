run:
	go run main.go --lint --render -i ./dist/manifest.amauta

run-black:
	go run main.go --lint --render -i ./dist/manifest.amauta -theme dark

test:
	@go test -v -cover ./... | awk '\
	/^=== RUN/ { print "\033[1;34m" $$0 "\033[0m"; next } \
	/^--- PASS/ { print "\033[1;32m" $$0 "\033[0m"; next } \
	/^--- FAIL/ { print "\033[1;31m" $$0 "\033[0m"; next } \
	/^PASS$$/ { print "\033[1;32m" $$0 "\033[0m"; next } \
	/^FAIL$$/ { print "\033[1;31m" $$0 "\033[0m"; next } \
	{ print }'

dlv-test:
	@dlv test $(PKG) -- -test.run ^$(TEST)$$

