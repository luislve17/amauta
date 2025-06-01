test:
	go test -v -cover ./...

dlv-test:
	dlv test $(PKG) -- -test.run ^$(TEST)$$
