package tests

import (
	"fmt"
	"testing"
)

var (
	monitorSixteenFile = "challenged_proposal.gate"
)

func TestChallengedProposalChallengerAttacksRootClaim(t *testing.T) {
	// We expect an alert to be fired when the challenger attacks the root claim

	// set the params, which DO matter for these tests
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		// we only use move events for length, but we still need the shape of the data to be accurate
		"moveEvents": [][]interface{}{
			{2, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x00"},
		},
		"claimCount": 4,
		"claimData": [][]interface{}{
			// the root claim doesn't have a real parent index since it is the root, so the index is type(uint32).max
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// root claim is being attacked by the CB challenger
			{0, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 2, 123456},
			{1, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 3, 123456},
			{2, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 4, 1233456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorSixteenFile)
	}
}

func TestChallengedProposalChallengerDefendsRootClaim(t *testing.T) {
	// We DO NOT expect an alert to be fired when the challenger defends the root claim

	// set the params
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"moveEvents": [][]interface{}{
			{2, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x00"},
		},
		"claimCount": 4,
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			{0, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 2, 123456},
			// root claim is being defended by the CB challenger
			{1, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 3, 123456},
			{2, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 4, 1233456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSixteenFile)
	}
}

func TestChallengedProposalChallengerNotCB(t *testing.T) {
	// We DO NOT expect an alert to be fired when the challenger is not CB, regardless of
	// whether the root claim submitted by CB proposer is challenged or not

	// set the params
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"moveEvents": [][]interface{}{
			{2, "0x09dE888033b1e815419a3fb865f0DA5689332FdB", "0x00"},
		},
		"claimCount": 4,
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// root claim is challenged, but by a random address
			{0, "0x0000000000000000000000000000000000000000", "0x09dE888033b1e815419a3fb865f0DA5689332FdB", 1, "0x00", 2, 123456},
			{1, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 3, 123456},
			{2, "0x0000000000000000000000000000000000000000", "0x09dE888033b1e815419a3fb865f0DA5689332FdB", 1, "0x00", 4, 1233456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSixteenFile)
	}
}

func TestChallengedProposalProposerNotCB(t *testing.T) {
	// We DO NOT expect an alert to be fired when the proposer is not CB, regardless of
	// whether the challenger attacks the root claim or not

	// set the params
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"moveEvents": [][]interface{}{
			{2, "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", "0x00"},
		},
		"claimCount": 4,
		"claimData": [][]interface{}{
			// root claim is NOT proposed by the CB proposer
			{4294967295, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 1, 123456},
			// root claim is challenged by the CB challenger
			{0, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 2, 123456},
			{1, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 3, 123456},
			{2, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 4, 1233456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSixteenFile)
	}
}

func TestChallengedProposalOnlyRootClaim(t *testing.T) {
	// We DO NOT expect an alert to be fired when the only claim is the root claim

	// set the params
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"moveEvents": [][]interface{}{
			{0, "0x49277EE36A024120Ee218127354c4a3591dc90A9", "0x00"},
		},
		"claimCount": 4,
		"claimData": [][]interface{}{
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// root claim is unchallenged
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSixteenFile)
	}
}

func TestChallengedProposalNoMoveEventInBlock(t *testing.T) {
	// We DO NOT expect an alert to be fired when there is no move event in the block,
	// regardless of whether the claimData indicates the root claim is being challenged or not

	// set the params, which DO matter for these tests
	params := map[string]any{
		"disputeGame":  "0x0000000000000000000000000000000000000000",
		"cbProposer":   "0x49277EE36A024120Ee218127354c4a3591dc90A9",
		"cbChallenger": "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorSixteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorSixteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		// same setup as the first test, except no moveEvents have been emitted in the block
		"claimCount": 4,
		"claimData": [][]interface{}{
			// the root claim doesn't have a real parent index since it is the root, so the index is type(uint32).max
			{4294967295, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// root claim is being attacked by the CB challenger
			{0, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 2, 123456},
			{1, "0x0000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000", 1, "0x00", 3, 123456},
			{2, "0x0000000000000000000000000000000000000000", "0xc96775081bcA132B0E7cbECDd0B58d9Ec07Fdaa4", 1, "0x00", 4, 1233456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorSixteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorSixteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorSixteenFile)
	}
}
