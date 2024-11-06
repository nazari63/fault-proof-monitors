package tests

import (
	"fmt"
	"testing"
)

var (
	monitorTwentyFile = "unresolvable_dispute_game.gate"
)

func TestUnresolvableDisputeGame(t *testing.T) {
	// We expect an alert to be fired when a dispute game has not resolved within the time limit

	// set the params
	params := map[string]any{
		"disputeGame":        "0x0000000000000000000000000000000000000000",
		"extraTimeInSeconds": 172800,
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTwentyFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTwentyFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"creationTimestamp": 555555,
		"gameDuration":      100,
		"resolvedAt":        0,      // game hasn't resolved yet
		"currentTimestamp":  728556, // creationTimestamp + (2 * gameDuration) + extraTime + 1
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTwentyFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTwentyFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorTwentyFile)
	}
}

func TestUnresolvableDisputeGameNoAlertGameUnderTimeLimit(t *testing.T) {
	// We DO NOT expect an alert to be fired when a dispute game has not resolved but
	// the time limit has not been reached

	// set the params
	params := map[string]any{
		"disputeGame":        "0x0000000000000000000000000000000000000000",
		"extraTimeInSeconds": 172800,
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTwentyFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTwentyFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"creationTimestamp": 555555,
		"gameDuration":      100,
		"resolvedAt":        0,      // game hasn't resolved yet
		"currentTimestamp":  728554, // creationTimestamp + (2 * gameDuration) + extraTime - 1
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTwentyFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTwentyFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorTwentyFile)
	}
}

func TestUnresolvableDisputeGameNoAlertGameResolved(t *testing.T) {
	// We DO NOT expect an alert to be fired when a dispute game has resolved

	// set the params
	params := map[string]any{
		"disputeGame":        "0x0000000000000000000000000000000000000000",
		"extraTimeInSeconds": 172800,
	}

	// read in the gate file
	data, err := ReadGateFile(monitorTwentyFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorTwentyFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"creationTimestamp": 555555,
		"gameDuration":      100,
		"resolvedAt":        555655, // game has resolved already
		"currentTimestamp":  728554, // doesn't matter
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorTwentyFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorTwentyFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorTwentyFile)
	}
}
