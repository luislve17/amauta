test:
	@go test -v -cover ./... | awk '\
	/^=== RUN/ { print "\033[1;34m" $$0 "\033[0m"; next } \
	/^--- PASS/ { print "\033[1;32m" $$0 "\033[0m"; next } \
	/^--- FAIL/ { print "\033[1;31m" $$0 "\033[0m"; next } \
	/^PASS$$/ { print "\033[1;32m" $$0 "\033[0m"; next } \
	/^FAIL$$/ { print "\033[1;31m" $$0 "\033[0m"; next } \
	/coverage: [0-9.]+%/ { \
		split($$2, a, "%"); \
		cov = a[1] + 0; \
		color = cov < 50 ? "\033[1;31m" : (cov < 80 ? "\033[1;33m" : "\033[1;32m"); \
		print color $$0 "\033[0m"; next \
	} \
	{ print }'

dlv-test:
	@dlv test $(PKG) -- -test.run ^$(TEST)$$

