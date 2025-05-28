package main

import (
	"fmt"
	"github.com/luislve17/amauta/linter"
)

func main() {
	fmt.Println("Hello world!")
	linter.LintFromRoot("foo", false)
}
