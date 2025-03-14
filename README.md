# scriplets-based go juju-doctor

Run juju-doctor on a set of event-based scriplets that will act as probes to perform validations.

juju-doctor will run a set of `juju` commands:
- `juju status ...`
- `juju show-unit ...` TODO
- `juju export-bundle ...` TODO 

Then, it will trigger events with the outputs as payloads, allowing scriplets to observe and perform necessary validations.

See example scriplets in `examples/`.

## Build juju-doctor
```
export GOPRIVATE=github.com/canonical/*
make build
```
This will generate a GO binary `juju-doctor` inside `./bin`


## Run on a live model
```
./bin/juju-doctor --scriplet ./examples/grafana_agent.star --model test
```
This will run a scriplet (e.g: a grafana-agent specific scriplet) to validate against a Juju model `test`

## Run with multiple artifacts
TODO