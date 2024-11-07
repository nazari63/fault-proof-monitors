package tests

import (
	"fmt"
	"strings"
	"testing"
)

var (
	monitorThirteenFile = "challenger_loses.gate"
)

func TestChallengerLostTopLevelChallenge(t *testing.T) {
	// We expect an alert to be fired if the honest challenger was challenging a root claim and the claim resolved in favor of the defenders

	// set the params
	params := map[string]any{
		"disputeGame":      "0x000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveEvents": [][]interface{}{
			{2}, // resolution status of the dispute game, 2 = DEFENDER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // challenger attacks root claim
			{1, "0x01", "0x00000000000000000000000000000000000000AA"}, // defender moves against challenger
		},
		"claimCount": 3, // 3 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was not countered
			{11111111, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 0, "0x00", 0, 123455},
			// cb challenger claim was countered successfully
			{0, "0x00000000000000000000000000000000000000AA", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// defender claim was also not countered
			{1, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 2, "0x00", 2, 123457},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorThirteenFile)
	}

	expectedAlert := "Challenger lost the dispute game while challenging a state root"
	foundAlert := false
	// we expect to see the specific alert fired
	for _, alert := range failed {
		alertStr := alert.([]interface{})[0]
		if strings.Contains(alertStr.(string), expectedAlert) {
			foundAlert = true
			break
		}
	}

	if !foundAlert {
		fmt.Println(failed)
		fmt.Println(trace)
		t.Errorf("Monitor did not fire the expected alert for %s", monitorThirteenFile)
	}
}

func TestChallengerLostTopLevelDefenseAndSubgame(t *testing.T) {
	// We expect an alert to be fired if the honest challenger was defending a root claim and the claim resolved in favor of the other challengers

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveEvents": [][]interface{}{
			{1}, // resolution status of the dispute game, 1 = CHALLENGER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x00000000000000000000000000000000000000AA"}, // attacker challenges root claim
			{1, "0x01", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // challenger defends root claim by challenging the attacker's claim
			{2, "0x02", "0x00000000000000000000000000000000000000AA"}, // attacker challenges honest challenger's claim
		},
		"claimCount": 4, // 4 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was countered successfully
			{11111111, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000BB", 0, "0x00", 0, 123455},
			// attacker claim was not countered
			{0, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 1, "0x11", 1, 123456},
			// honest challenger defense move was countered
			{1, "0x00000000000000000000000000000000000000AA", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 2, "0x22", 2, 123457},
			// attacker claim was not countered
			{2, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 3, "0x33", 3, 123458},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorThirteenFile)
	}

	expectedAlert := "Challenger lost the dispute game while defending a state root"
	foundAlert := false
	// we expect to see the specific alert fired
	for _, alert := range failed {
		alertStr := alert.([]interface{})[0]
		if strings.Contains(alertStr.(string), expectedAlert) {
			foundAlert = true
			break
		}
	}

	if !foundAlert {
		fmt.Println(failed)
		fmt.Println(trace)
		t.Errorf("Monitor did not fire the expected alert for %s", monitorThirteenFile)
	}
}

func TestChallengerLostSubgames(t *testing.T) {
	// We expect an alert to be fired when the honest challenger loses any subgame claim, even if the top-level game was won

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveEvents": [][]interface{}{
			{2}, // resolution status of the dispute game, 2 = DEFENDER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x00000000000000000000000000000000000000AA"}, // attacker challenges root claim
			{1, "0x1a", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // challenger defends root claim by challenging the attacker's claim
			{1, "0x1b", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // challenger (unrealistically) defends the root claim again on the same claim index
			{2, "0x02", "0x00000000000000000000000000000000000000AA"}, // attacker challenges one of the honest challenger's claim
		},
		"claimCount": 5, // 5 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim not countered
			{11111111, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000BB", 0, "0x00", 0, 123455},
			// attacker claim on root claim is countered
			{0, "0x49277EE36A024120Ee218127354c4a3591dc90A9", "0x00000000000000000000000000000000000000AA", 1, "0x33", 1, 123456},
			// challenger first defense move is uncountered
			{1, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 2, "0x11", 2, 123457},
			// challenger second defense move was countered
			{1, "0x00000000000000000000000000000000000000AA", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 2, "0x22", 2, 123458},
			// attacker claim on honest challenger's second defense move was not countered
			{2, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 3, "0x33", 3, 123459},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we expect to see the alert fired
	if len(failed) == 0 {
		fmt.Println(trace)
		t.Errorf("Monitor did not fire an alert for %s when it was supposed to", monitorThirteenFile)
	}

	expectedAlert := "Challenger lost one or more subgames"
	foundAlert := false
	// we expect to see the specific alert fired
	for _, alert := range failed {
		alertStr := alert.([]interface{})[0]
		if strings.Contains(alertStr.(string), expectedAlert) {
			foundAlert = true
			break
		}
	}

	if !foundAlert {
		fmt.Println(failed)
		fmt.Println(trace)
		t.Errorf("Monitor did not fire the expected alert for %s", monitorThirteenFile)
	}
}

func TestChallengerWins(t *testing.T) {
	// We DO NOT expect an alert to be fired when the honest challenger wins all the claims it makes

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveEvents": [][]interface{}{
			{1}, // resolution status of the dispute game, 1 = CHALLENGER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // honest hallenger attacks root claim
		},
		"claimCount": 2, // 2 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was countered by the honest challenger
			{11111111, "0x49277EE36A024120Ee218127354c4a3591dc90A9", "0x00000000000000000000000000000000000000AA", 0, "0x00", 0, 123455},
			// challenger claim was not countered
			{0, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorThirteenFile)
	}
}

func TestChallengerGameInProgress(t *testing.T) {
	// We DO NOT expect an alert to be fired when the dispute game is still in progress

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"addressesInTrace": []any{"0x000000000000000000000000000000000000000"},
		"resolveEvents": [][]interface{}{
			{0}, // resolution status of the dispute game, 0 = IN_PROGRESS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // honest challenger attacks root claim
		},
		"claimCount": 2, // 2 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// resolution of all claims has not occurred yet
			{11111111, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 0, "0x00", 0, 123455},
			{0, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorThirteenFile)
	}
}

func TestChallengerLostTopLevelChallengeNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when the honest challenger loses a top-level challenge and there is no filtered address

	// set the params
	params := map[string]any{
		"disputeGame":      "0x000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"resolveEvents": [][]interface{}{
			{2}, // resolution status of the dispute game, 2 = DEFENDER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // honest challenger attacks root claim
			{1, "0x01", "0x00000000000000000000000000000000000000AA"}, // defender moves against the honest challenger
		},
		"claimCount": 3, // 3 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was not countered
			{11111111, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 0, "0x00", 0, 123455},
			// challenger claim was countered successfully
			{0, "0x00000000000000000000000000000000000000AA", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
			// defender claim was also not countered
			{1, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 2, "0x00", 2, 123457},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorThirteenFile)
	}
}

func TestChallengerLostTopLevelDefenseAndSubgameNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when the honest challenger loses a top-level defense and subgame and there is no filtered address

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"resolveEvents": [][]interface{}{
			{1}, // resolution status of the dispute game, 1 = CHALLENGER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x00000000000000000000000000000000000000AA"}, // attacker challenges root claim
			{1, "0x01", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // honest challenger defends root claim by challenging the attacker's claim
			{2, "0x02", "0x00000000000000000000000000000000000000AA"}, // attacker challenges the honest challenger's claim
		},
		"claimCount": 4, // 4 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was countered successfully
			{11111111, "0x00000000000000000000000000000000000000AA", "0x00000000000000000000000000000000000000BB", 0, "0x00", 0, 123455},
			// attacker claim was not countered
			{0, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 1, "0x11", 1, 123456},
			// challenger defense move was countered
			{1, "0x00000000000000000000000000000000000000AA", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 2, "0x22", 2, 123457},
			// attacker claim was not countered
			{2, "0x0000000000000000000000000000000000000000", "0x00000000000000000000000000000000000000AA", 3, "0x33", 3, 123458},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorThirteenFile)
	}
}

func TestChallengerLostSubgamesNoFilterAddress(t *testing.T) {
	// We DO NOT expect an alert to be fired when the honest challenger loses any subgame claim and there is no filtered address

	// set the params
	params := map[string]any{
		"disputeGame":      "0x0000000000000000000000000000000000000000",
		"honestChallenger": "0x49277EE36A024120Ee218127354c4a3591dc90A9",
	}

	// read in the gate file
	data, err := ReadGateFile(monitorThirteenFile)
	if err != nil {
		t.Errorf("Error reading file %s: %v", monitorThirteenFile, err)
	}

	// set the mock data that we will pass along with the Gate file and params to the validate request endpoint
	mocks := map[string]any{
		"resolveEvents": [][]interface{}{
			{1}, // resolution status of the dispute game, 1 = CHALLENGER_WINS
		},
		"historicalMoveEvents": [][]interface{}{
			{0, "0x00", "0x49277EE36A024120Ee218127354c4a3591dc90A9"}, // honest challenger attacks root claim
		},
		"claimCount": 2, // 2 claims total, inclusive of the root claim which doesn't count as a Move
		"claimResults": [][]interface{}{
			// root claim was countered by cb challenger
			{11111111, "0x49277EE36A024120Ee218127354c4a3591dc90A9", "0x00000000000000000000000000000000000000AA", 0, "0x00", 0, 123455},
			// challenger claim was not countered
			{0, "0x0000000000000000000000000000000000000000", "0x49277EE36A024120Ee218127354c4a3591dc90A9", 1, "0x00", 1, 123456},
		},
	}

	// call the validate request endpoint and parse the results
	failed, exceptions, trace, err := HandleValidateRequest(data, params, mocks)
	if err != nil {
		t.Errorf("Error handling validate request for %s: %v", monitorThirteenFile, err)
	}

	// check if the validate request threw any exceptions
	if len(exceptions) > 0 {
		fmt.Println(trace)
		t.Errorf("Exceptions for %s: %v", monitorThirteenFile, exceptions)
	}

	// we DO NOT expect to see the alert fired
	if len(failed) > 0 {
		fmt.Println(trace)
		t.Errorf("Monitor fired an alert for %s when it was not supposed to", monitorThirteenFile)
	}
}
