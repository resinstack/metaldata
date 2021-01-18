package main

import (
	"os"

	"github.com/hashicorp/go-hclog"

	"github.com/resinstack/metaldata/pkg/http"
	"github.com/resinstack/metaldata/pkg/source/fs"
)

func main() {
	rootLog := hclog.New(&hclog.LoggerOptions{
		Name:  "metaldata",
		Level: hclog.LevelFromString("TRACE"),
	})

	s := http.New(rootLog)
	s.SetSource(fs.New(os.Getenv("MD_FSBASE")))
	s.Start(os.Getenv("MD_BIND"))
}
