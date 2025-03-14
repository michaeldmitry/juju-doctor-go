package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/canonical/juju-doctor/internal/engine"
	"github.com/canonical/juju-doctor/internal/utils"
)

// StringSlice is a custom type for handling multiple `-s` flags
type StringSlice []string

func (s *StringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type Flags struct {
	Model string
	Paths StringSlice
}

func parseFlags() *Flags {
	flags := &Flags{}

	flag.Var(&flags.Paths, "scriplet", "path of a scriplet containing probes to execute. Can be used multiple times.")
	flag.StringVar(&flags.Model, "model", "", "Specify the model name")

	// Parse the flags
	flag.Parse()
	return flags
}

func main() {
	log.SetFlags(0)
	ctx := context.Background()
	flags := parseFlags()
	if len(flags.Paths) == 0 {
		log.Fatalf("Error: at least one --scriplet flag is required")
	}

	scripletsFetcher, err := utils.NewFetcher()
	if err != nil {
		log.Fatalf("error creating scriplets fetcher %s", err.Error())
	}
	// Ensure cleanup
	defer os.RemoveAll(scripletsFetcher.Destination)
	if err := scripletsFetcher.CopyScriplets(flags.Paths); err != nil {
		log.Fatalf("error copying scriplets to filesystem %s", err.Error())
	}

	starlarkEngine, err := engine.NewStarlarkEngine(ctx)
	if err != nil {
		log.Fatalf("error creating a starlark engine %s", err.Error())
	}

	// load available scriplets
	if err := starlarkEngine.LoadScriplets(ctx, scripletsFetcher.Destination); err != nil {
		log.Fatalf("error loading scriplets %s", err.Error())
	}

	// Fire pre-defined events
	if err := starlarkEngine.FireStatusReadyEvent(ctx, flags.Model); err != nil {
		log.Fatalf("error firing status ready event %s", err.Error())
	}
	if err := starlarkEngine.FireBundleReadyEvent(ctx, flags.Model); err != nil {
		log.Fatalf("error firing bundle ready event %s", err.Error())
	}
	if err := starlarkEngine.FireShowUnitReadyEvent(ctx, flags.Model); err != nil {
		log.Fatalf("error firing show-unit ready event %s", err.Error())
	}

}
