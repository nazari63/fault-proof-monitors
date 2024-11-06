package tests

import (
	"fmt"
	"testing"
)

var (
	monitorFiveFile = "duplicate_dispute_game.gate"
)

func TestDuplicateDisputeGameCreated(t *testing.T) {
	// We expect an alert to be fired when a dispute game is created in the current block that
	// has the same UUID as a previous dispute game

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim := "0x17bdb49e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c2143d99bf"
	extraData := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid := "0x4f73e8da3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128e32ece"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames": [][]interface{}{
			{0, rootClaim, extraData},
		},
		"createdDisputeGames": [][]interface{}{
			{99, []interface{}{"0x0000000000000000000000000000000000000000", 0, rootClaim}},
		},
		"createdDisputeGamesExtraData": []string{extraData},
		// technically it's the UUIDs that will trigger the alert, but we still want to mock the other data
		// as close as possible to the actual data
		"newDisputeGameUUIDs":      []string{uuid},
		"previousDisputeGameUUIDs": []string{uuid},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// check if the validate request failed
	// in an inverse way, this indicates that the monitor successfully fired an alert as the invariant was breached
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorFiveFile)
	}
}

func TestMultipleDuplicateDisputeGamesCreated(t *testing.T) {
	// We expect an alert to be fired when multiple dispute games are created in the current block
	// that have the same UUID as previous dispute game(s)

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim1 := "0xbbbbbb9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214bbbbbb"
	extraData1 := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid1 := "0xbbbbbbda3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128bbbbbb"

	rootClaim2 := "0xaaaaaa9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214aaaaaa"
	extraData2 := "0x0000000000000000000000000000000000000000000000000000000000aaaaaa"
	uuid2 := "0xaaaaaada3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128aaaaaa"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames": [][]interface{}{
			{0, rootClaim1, extraData1},
			{0, rootClaim2, extraData2},
		},
		"createdDisputeGames": [][]interface{}{
			{98, []interface{}{"0x0000000000000000000000000000000000000000", 0, rootClaim1}},
			{99, []interface{}{"0x0000000000000000000000000000000000000000", 0, rootClaim2}},
		},
		"createdDisputeGamesExtraData": []string{extraData1, extraData2},
		"newDisputeGameUUIDs":          []string{uuid1, uuid2},
		"previousDisputeGameUUIDs":     []string{uuid1, uuid2},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// check if the validate request failed
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorFiveFile)
	}

	// check to make sure two duplicate dispute game instances were identified
	duplicateGames := trace.(map[string]interface{})["foundDuplicateGameInfo"]
	if len(duplicateGames.([]interface{})) != 2 || duplicateGames.([]interface{})[0].(bool) != true || duplicateGames.([]interface{})[1].(bool) != true {
		fmt.Println(trace)
		t.Errorf("Monitor did not identify the correct number of duplicate dispute games")
	}
}

func TestDuplicateDisputeGameCreatedInSameBlock(t *testing.T) {
	// We expect an alert to be fired if more than one dispute game is created in the current block
	// and more than one of the newly-created dispute games have the same UUID

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim1 := "0xbbbbbb9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214bbbbbb"
	extraData1 := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid1 := "0xbbbbbbda3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128bbbbbb"

	rootClaim2 := "0xaaaaaa9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214aaaaaa"
	extraData2 := "0x0000000000000000000000000000000000000000000000000000000000aaaaaa"
	uuid2 := "0xaaaaaada3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128aaaaaa"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames": [][]interface{}{
			{0, rootClaim1, extraData1},
			{0, rootClaim1, extraData1},
		},
		"createdDisputeGames": [][]interface{}{
			{98, []interface{}{"0x0000000000000000000000000000000000000000", 0, rootClaim2}},
		},
		"createdDisputeGamesExtraData": []string{extraData2},
		"newDisputeGameUUIDs":          []string{uuid1, uuid1},
		"previousDisputeGameUUIDs":     []string{uuid2},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// check if the validate request failed
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorFiveFile)
	}
}

func TestDuplicateDisputeGameCreatedDifferentGameTypes(t *testing.T) {
	// We DO NOT expect an alert to be fired if a dispute game is created in the current block
	// that has the same UUID but a different game type as a previous dispute game

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim1 := "0xbbbbbb9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214bbbbbb"
	extraData1 := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid1 := "0xbbbbbbda3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128bbbbbb"

	rootClaim2 := "0xaaaaaa9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214aaaaaa"
	extraData2 := "0x0000000000000000000000000000000000000000000000000000000000aaaaaa"
	uuid2 := "0xaaaaaada3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128aaaaaa"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames": [][]interface{}{
			{0, rootClaim1, extraData1},
			{0, rootClaim2, extraData2},
		},
		"createdDisputeGames": [][]interface{}{
			{98, []interface{}{"0x0000000000000000000000000000000000000000", 2, rootClaim1}},
			{99, []interface{}{"0x0000000000000000000000000000000000000000", 2, rootClaim2}},
		},
		"createdDisputeGamesExtraData": []string{extraData1, extraData2},
		"newDisputeGameUUIDs":          []string{uuid1, uuid2},
		"previousDisputeGameUUIDs":     []string{},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// in this case we DO NOT expect the monitor to fire an alert
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorFiveFile)
	}
}

func TestNoDuplicateDisputeGameCreatedNoHistory(t *testing.T) {
	// We DO NOT expect an alert to be fired if a dispute game is created in the current block
	// and there is no history of a dispute game being created with the same UUID

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim1 := "0xbbbbbb9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214bbbbbb"
	extraData1 := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid1 := "0xbbbbbbda3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128bbbbbb"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames": [][]interface{}{
			{0, rootClaim1, extraData1},
		},
		"createdDisputeGames":          [][]interface{}{},
		"createdDisputeGamesExtraData": []string{},
		"newDisputeGameUUIDs":          []string{uuid1},
		"previousDisputeGameUUIDs":     []string{},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// in this case we DO NOT expect the monitor to fire an alert
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorFiveFile)
	}
}

func TestNoDuplicateDisputeGameCreatedNoGameCreatedInCurrBlock(t *testing.T) {
	// We DO NOT expect an alert to be fired if no dispute games are created in the current block
	// regardless of whether there are historical instances of duplciate dispute games being created

	// setup the param value - this doesn't matter as much as we'll be mocking source data
	params := map[string]any{
		"optimismPortalProxy": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorFiveFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorFiveFile, err)
	}

	// setup the mock data
	rootClaim1 := "0xbbbbbb9e89561f18e1dc284c1955238d2b942e0fa3b755279fce78c214bbbbbb"
	extraData1 := "0x0000000000000000000000000000000000000000000000000000000000bbbbbb"
	uuid1 := "0xbbbbbbda3b9d2fa9933b09187ee8b678b03fc2255e67975017d3462128bbbbbb"

	mocks := map[string]any{
		"disputeGameFactory": "0x0000000000000000000000000000000000000000",
		"respectedGameType":  0,
		"currBlock":          100,
		"newDisputeGames":    [][]interface{}{},
		"createdDisputeGames": [][]interface{}{
			{98, []interface{}{"0x0000000000000000000000000000000000000000", 0, rootClaim1}},
		},
		"createdDisputeGamesExtraData": []string{extraData1},
		"newDisputeGameUUIDs":          []string{},
		"previousDisputeGameUUIDs":     []string{uuid1},
	}

	// call out to hexagate API to run the gate file with params and mocks
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorFiveFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorFiveFile, exceptions)
	}

	// in this case we DO NOT expect the monitor to fire an alert
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorFiveFile)
	}
}
