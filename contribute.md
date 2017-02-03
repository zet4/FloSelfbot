# Go Contribution Guidelines

Please check these guidelines before contributing go code to a repository.

## Requirements

-   Run `go fmt` before pushing any code.
-   Run `golint` and `go vet` too -- some code may fail.

## A Short Introduction

If you are new to our Go development workflow:

-   Ensure you have [Go installed on your system](https://golang.org/doc/install).
-   Make sure that you have the environment variable `GOPATH` set somewhere, e.g. `$HOME/gopkg`
-   Clone FloSelfbot into the path `$GOPATH/src/github.com/Moonlington/FloSelfbot`
    -   NOTE: This is true even if you have forked FloSelfbot, dependencies in go are path based and must be in the right locations.
-   You are now free to make changes to the codebase as you please.
-   You can build the binary by running `go build bot.go` from the FloSelfbot directory.
    -   NOTE: when making changes remember to restart your daemon to ensure its running your new code.

## Imports

We strive to use the following convention when it comes to imports. First list stdlib imports, then local repository imports and then all other external imports. Separate them using one empty new line so they are not reordered by `go fmt` or `goimports`.

Example:
```go
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Moonlington/FloSelfbot/commands"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)
```

If a package name isn't the same as its directory, use the explicit name.
