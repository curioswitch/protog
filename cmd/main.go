package main

import (
	"os"
	"strings"

	"github.com/curioswitch/protog/internal/cmd"
)

func main() {
	env := map[string]string{}
	for _, e := range os.Environ() {
		if key, value, ok := strings.Cut(e, "="); ok {
			env[key] = value
		}
	}
	_ = cmd.Run(os.Args[1:], env)
}
