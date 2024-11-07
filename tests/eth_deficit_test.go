package tests

import (
	"fmt"
	"testing"
)

var (
	monitorElevenFile = "eth_deficit.gate"
)

func TestETHDeficitTotalCreditDeficit(t *testing.T) {
	// We expect an alert to be fired when totalCredit is less than claimCredit

	// set the params, which don't really matter for these tests
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"claimCredit":           100,
		"totalCredit":           []interface{}{50, 123456},
		"ethBalanceDisputeGame": 500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitTotalETHBalanceDeficit(t *testing.T) {
	// We expect an alert to be fired when ethBalanceDisputeGame is less than totalCredit

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"claimCredit":           50,
		"totalCredit":           []interface{}{150, 123456},
		"ethBalanceDisputeGame": 100,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitCurrCreditZeroTotalCreditNonZero(t *testing.T) {
	// We expect an alert to be fired when claimCredit is zero and totalCredit is non-zero

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"claimCredit":           0,
		"totalCredit":           []interface{}{150, 123456},
		"ethBalanceDisputeGame": 1500,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorElevenFile)
	}
}

func TestETHDeficitNoDeficit(t *testing.T) {
	// We DO NOT expect an alert to be fired if there is no deficit

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorElevenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorElevenFile, err)
	}

	// set the mock data
	mocks := map[string]any{
		"delayedWETH":           "0x0000000000000000000000000000000000000000",
		"claimCredit":           100,
		"totalCredit":           []interface{}{200, 123456},
		"ethBalanceDisputeGame": 300,
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorElevenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorElevenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorElevenFile)
	}
}
