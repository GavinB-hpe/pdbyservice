package main

import (
	"reflect"
	"testing"

	"github.com/gavinB-hpe/pdbyservice/dbtalker"
	"github.com/gavinB-hpe/pdbyservice/globals"
	"github.com/gavinB-hpe/pdbyservice/model"
)

func TestPrintIncidents(t *testing.T) {
	// Create a mock DBTalker
	dbt := dbtalker.NewDBTalker(model.ConnectDatabase(globals.DEFAULTDBTYPE, "test.db"))

	// Set up test data
	keys := []string{"service1", "service2"}
	mxc := 10

	// Call the function being tested
	printIncidents(true, keys, mxc, dbt)

	// Add your assertions here
}

func TestGetPerServiceIncidents(t *testing.T) {
	// Create a mock DBTalker
	dbt := dbtalker.NewDBTalker(model.ConnectDatabase(globals.DEFAULTDBTYPE, "test.db"))

	// Set up test data
	sr := true
	key := "service1"

	// Call the function being tested
	result := getPerServiceIncidents(sr, key, dbt)

	// Add your assertions here
	if result == nil {
		t.Errorf("Got unexpected nil %v", result)
	}
}

func TestCheckSeenServices(t *testing.T) {
	sn := map[string]string{
		"service1": "Service 1",
		"service2": "Service 2",
		"service3": "Service 3",
	}

	sd := map[string]map[string]string{
		"service1": {
			"name": "Service 1",
		},
		"service3": {
			"name": "Service 3",
		},
	}

	expected := map[string]string{
		"service2": "Service 2",
	}

	result := checkSeenServices(sn, &sd)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
func TestReadServiceData(t *testing.T) {
	// Set up test data
	fn := "services-data.json"

	// Call the function being tested
	result := readServiceData(fn)

	// Add your assertions here
	if len(*result) == 0 {
		t.Errorf("Got unexpected empty map %v", result)
	}
}
