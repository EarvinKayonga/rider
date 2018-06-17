package main

import (
	"log"
	"os"

	"github.com/EarvinKayonga/rider/application"
)

func main() {
	err := application.RunTrip(
		os.Args,
		application.Metadata{
			Branch:     branch,
			Compiler:   compiler,
			CompiledAt: compiledAt,
			Sha:        sha,
		},
	)
	if err != nil {
		log.Fatalf("%v occured with %v", err, os.Args)
	}
}

var (
	branch     string
	sha        string
	compiledAt string
	compiler   string
)
