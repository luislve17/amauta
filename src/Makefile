build:
	go build -o amauta

run:
	go run main.go --render -i ./dist/docs/main.amauta -o dist/doc.html -theme $(theme)

run-debug:
	go run main.go --lint --render -i ./dist/debug.amauta

test:
	@go test ./... -json | tparse -smallscreen -progress| awk '{ \
	  gsub("PASS", "\033[32mPASS\033[0m"); \
	  gsub("FAIL", "\033[31mFAIL\033[0m"); \
	  print \
	}'

dlv-test:
	@dlv test $(PKG) -- -test.run ^$(TEST)$$

