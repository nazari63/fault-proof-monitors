package tests

import (
	"fmt"
	"testing"
)

var (
	monitorSeventeenFile = "credit_and_bond_discrepancy.gate"
)

func TestCreditAndBondDiscrepancyWrongBondAmount(t *testing.T) {
	// We expect an alert to be fired when the bond amount does not match the credit amount

	// set the param, which doesn't matter for this test suite
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveCalls": [][]interface{}{
			{0, 111111},
			{1, 111111},
		},
		"delayedWeth": "0x0000000000000000000000000000000000000000", // doesn't matter
		"unlocks": [][]interface{}{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000001},
		},
		"claimData": [][]interface{}{
			// wrong bond amount
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 9999, "0x00", 1, 123456},
			// reward goes to the counter address vs. the claimant
			{0, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001, "0x00", 2, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyWrongCreditAmount(t *testing.T) {
	// We expect an alert to be fired when the credit amount does not match the bond amount

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveCalls": [][]interface{}{
			{0, 111111},
			{1, 111111},
		},
		"delayedWeth": "0x0000000000000000000000000000000000000000", // doesn't matter
		"unlocks": [][]interface{}{
			// wrong credit amounts
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 9999},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000000},
		},
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000, "0x00", 1, 123456},
			// reward goes to the counter address vs. the claimant
			{0, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001, "0x00", 2, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyWrongClaimantAddress(t *testing.T) {
	// We expect an alert to be fired when the claimant address does not match the credited address

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveCalls": [][]interface{}{
			{0, 111111},
			{1, 111111},
		},
		"delayedWeth": "0x0000000000000000000000000000000000000000", // doesn't matter
		"unlocks": [][]interface{}{
			// right credit amounts, wrong claimant address
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001},
		},
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000, "0x00", 1, 123456},
			// reward goes to the counter address vs. the claimant
			{0, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001, "0x00", 2, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyCorrectAmountsAndAddresses(t *testing.T) {
	// We DO NOT expect an alert to be fired if the bond and credit amounts match
	// and the claimant address matches the credited address

	// set the param
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveCalls": [][]interface{}{
			{0, 111111},
			{1, 111111},
		},
		"delayedWeth": "0x0000000000000000000000000000000000000000", // doesn't matter
		"unlocks": [][]interface{}{
			// right credit amounts, right claimant addresses
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000001},
		},
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000, "0x00", 1, 123456},
			// reward goes to the counter address vs. the claimant
			{0, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001, "0x00", 2, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSeventeenFile)
	}
}

func TestCreditAndBondDiscrepancyNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when there is no address filtered in the current block trace

	// set the param, which doesn't matter for this test suite
	params := map[string]any{
		"disputeGame": "0x0000000000000000000000000000000000000000",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSeventeenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSeventeenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"resolveCalls": [][]interface{}{
			{0, 111111},
			{1, 111111},
		},
		"delayedWeth": "0x0000000000000000000000000000000000000000", // doesn't matter
		"unlocks": [][]interface{}{
			{"0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000000},
			{"0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1000001},
		},
		"claimData": [][]interface{}{
			// wrong bond amount
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 9999, "0x00", 1, 123456},
			// reward goes to the counter address vs. the claimant
			{0, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1000001, "0x00", 2, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSeventeenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSeventeenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSeventeenFile)
	}
}
