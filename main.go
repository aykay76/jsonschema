package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

type EntityRegistry struct {
	Entities []map[string]interface{}
	Events   []map[string]interface{}
}

type Entity struct {
	Properties map[string]interface{} `json:"properties"`
}

// example JSON string that contains an entity with properties
var jsonData = []byte(`{
	"name": "Country",
	"code": "US",
	"population": 331002651
}`)

func main() {
	entityRegistry, err := LoadSimulation()
	if err != nil {
		fmt.Println("Error loading simulation:", err)
		return
	}

	for _, entity := range entityRegistry.Entities {
		fmt.Println("Entity:", entity["type"], entity["id"], entity["name"])
	}
}

func LoadAndValidate(filename string, schemaPath string) ([]map[string]interface{}, error) {
	// 1. Read JSON file
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 2. Validate against schema
	schema := gojsonschema.NewReferenceLoader("file://./" + schemaPath)
	document := gojsonschema.NewStringLoader(string(jsonData))
	result, err := gojsonschema.Validate(schema, document)
	if err != nil {
		return nil, fmt.Errorf("Error validating JSON: %v", err)
	}

	if !result.Valid() {
		return nil, fmt.Errorf("Validation failed: %v", result.Errors())
	}

	// 3. Unmarshal into array of entities
	var entities []map[string]interface{}
	if err := json.Unmarshal(jsonData, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}

func LoadSimulation() (*EntityRegistry, error) {
	registry := &EntityRegistry{}

	// Load countries
	countries, err := LoadAndValidate("countries.json", "countries.schema.json")
	if err != nil {
		return nil, err
	}
	registry.Entities = countries

	// Load events
	events, err := LoadAndValidate("events.json", "events.schema.json")
	if err != nil {
		return nil, err
	}
	registry.Events = events

	return registry, nil
}
