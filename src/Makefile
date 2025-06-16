test:
	cd core && go test -v -cover ./...

dlv-test:
	cd core && dlv test $(PKG) -- -test.run ^$(TEST)$$
