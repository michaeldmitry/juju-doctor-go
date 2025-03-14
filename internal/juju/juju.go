package juju

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
	"os/exec"
	"strings"
)

// GetJujuBundleOutput fetches the Juju bundle for the specified model.
func GetJujuBundleOutput(model string) (any, error) {

	args := []string{"export-bundle"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("juju", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error executing [juju %s]: %w", strings.Join(args, " "), err)
	}

	var singleDoc map[string]any
	var multiDoc []map[string]any
	// TODO Ask Juju if merging multi-doc into 1 YAML is fine, what is the purpose of multi docs vs. 1 YAML?
	if err := yaml.Unmarshal(output, &singleDoc); err == nil {
		return singleDoc, nil // It's a single YAML document
	} else if err := yaml.Unmarshal(output, &multiDoc); err == nil {
		return multiDoc, nil // It's a multi-document YAML
	} else {
		return nil, fmt.Errorf("error parsing Juju bundle YAML: %w", err)
	}
}

// GetJujuShowUnitOutput fetches Juju show-unit for all application units in the specified model.
func GetJujuShowUnitOutput(model string) (map[string]any, error) {

	// Get all the unit names from the juju status output
	jujuStatusObj, err := GetJujuStatusOutput(model)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(jujuStatusObj)
	if err != nil {
		return nil, fmt.Errorf("error converting Juju status to JSON %s", err)
	}
	unitNames := []string{}
	gjson.Get(string(jsonData), "applications").ForEach(func(appKey, appValue gjson.Result) bool {
		appValue.Get("units").ForEach(func(unitKey, _ gjson.Result) bool {
			unitNames = append(unitNames, unitKey.String())
			return true
		})
		return true
	})

	// Create a map of each unit's show-unit content
	results := make(map[string]any)
	for _, unitName := range unitNames {
		args := []string{"show-unit", unitName, "--format=json"}
		if model != "" {
			args = append(args, "--model", model)
		}
		cmd := exec.Command("juju", args...)
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error executing [juju %s]: %w", strings.Join(args, " "), err)
		}
		var data map[string]any
		if err := json.Unmarshal(output, &data); err != nil {
			return nil, fmt.Errorf("error parsing Juju status JSON for unit %s: %w", unitName, err)
		}
		results[unitName] = data
	}

	return results, nil
}

// GetJujuStatusOutput fetches the Juju status for the specified model.
func GetJujuStatusOutput(model string) (map[string]any, error) {

	args := []string{"status", "--format=json"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("juju", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing [juju %s]: %w", strings.Join(args, " "), err)
	}

	var data map[string]any
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("error parsing Juju status JSON: %w", err)
	}

	return data, nil
}
