package engine

import (
	"context"
	"fmt"
	"os"

	"github.com/canonical/juju-doctor/internal/juju"
	"github.com/canonical/juju-doctor/internal/utils"

	"github.com/canonical/starform/starform"
	"github.com/canonical/starlark/starlark"
)

var AppName = "jujudoctor"
var StatusReadyEvent = "status"
var BundleReadyEvent = "bundle"
var ShowUnitReadyEvent = "show_unit"

type StarlarkEngine struct {
	ScriptSet *starform.ScriptSet
}

// NewStarlarkEngine generate a new Starlark engine to call starlark scriplets
func NewStarlarkEngine(ctx context.Context) (*StarlarkEngine, error) {
	cache, err := starform.NewDefaultCache(&starform.DefaultCacheOptions{MaxSize: 100})
	if err != nil {
		return nil, err
	}

	scriptSet, err := starform.NewScriptSet(&starform.ScriptSetOptions{
		App:   &starform.AppObject{Name: AppName},
		Cache: cache,
		Logger: &utils.WriterLogger{
			Writer:       os.Stdout,
			MinimumLevel: starform.PrintLevel,
		},
		RequiredSafety: starlark.CPUSafe | starlark.IOSafe | starlark.MemSafe,
		MaxAllocs:      4 * 1024 * 1024,
		MaxSteps:       100_000,
	})
	if err != nil {
		return nil, err
	}

	return &StarlarkEngine{
		ScriptSet: scriptSet,
	}, nil
}

// LoadScripts initializes and loads all Starlark scripts from the passed paths.
func (engine *StarlarkEngine) LoadScriplets(ctx context.Context, directory string) error {

	fs := os.DirFS(directory)
	sources, err := starform.LoadDirSources(ctx, &starform.LoadDirSourcesOptions{
		FS: fs,
		// Root: directory,
	})
	if err != nil {
		return err
	}

	if err := engine.ScriptSet.LoadSources(ctx, sources); err != nil {
		return err
	}

	return nil
}

func (engine *StarlarkEngine) FireStatusReadyEvent(ctx context.Context, model string) error {
	jujuStatusObj, err := juju.GetJujuStatusOutput(model)
	if err != nil {
		fmt.Printf("error getting juju status output %s\n", err.Error())
		return err
	}
	starlarkStatusObj, err := utils.ToStarlarkDict(jujuStatusObj)
	if err != nil {
		fmt.Printf("error converting juju status output to Starlark Dict %s\n", err.Error())
		return err
	}

	// Fire a juju status ready event
	if err := engine.ScriptSet.Handle(ctx, &starform.EventObject{
		Name:  StatusReadyEvent,
		Attrs: starlark.StringDict{"input": starlarkStatusObj},
	}); err != nil {
		return err
	}
	return nil
}

func (engine *StarlarkEngine) FireBundleReadyEvent(ctx context.Context, model string) error {
	jujuBundleOutput, err := juju.GetJujuBundleOutput(model)
	if err != nil {
		fmt.Printf("error getting juju bundle output %s\n", err.Error())
		return err
	}
	starlarkBundleObj, err := utils.ToStarlarkValue(jujuBundleOutput)
	if err != nil {
		fmt.Printf("error converting juju bundle output to Starlark Dict %s\n", err.Error())
		return err
	}

	// Fire a juju bundle ready event
	if err := engine.ScriptSet.Handle(ctx, &starform.EventObject{
		Name:  BundleReadyEvent,
		Attrs: starlark.StringDict{"input": starlarkBundleObj},
	}); err != nil {
		return err
	}
	return nil
}

func (engine *StarlarkEngine) FireShowUnitReadyEvent(ctx context.Context, model string) error {
	jujuShowUnitObj, err := juju.GetJujuShowUnitOutput(model)
	if err != nil {
		fmt.Printf("error getting juju show-unit output %s\n", err.Error())
		return err
	}
	starlarkShowUnitObj, err := utils.ToStarlarkDict(jujuShowUnitObj)
	if err != nil {
		fmt.Printf("error converting juju show-unit output to Starlark Dict %s\n", err.Error())
		return err
	}

	// Fire a juju show-unit ready event
	if err := engine.ScriptSet.Handle(ctx, &starform.EventObject{
		Name:  ShowUnitReadyEvent,
		Attrs: starlark.StringDict{"input": starlarkShowUnitObj},
	}); err != nil {
		return err
	}
	return nil
}
