package main

import (
	"github.com/hashicorp/go-hclog"

	"github.com/resinstack/metaldata/pkg/http"
)

func main() {
	rootLog := hclog.New(&hclog.LoggerOptions{
		Name:  "metaldata",
		Level: hclog.LevelFromString("TRACE"),
	})

	s := http.New(rootLog)
	s.Start(":1234")
}
