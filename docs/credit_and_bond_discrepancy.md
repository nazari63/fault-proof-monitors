## Purpose

The `credit_and_bond_discrepancy.gate` monitor ensures the **safety** and integrity of the fault proof system by verifying that the ETH bond amounts credited during the resolution of dispute games match the original bond values. In a properly functioning system, when a dispute game or its subgames are resolved, the bond amounts unlocked and credited to participants should exactly equal the original bond amounts that were deposited. Discrepancies in these amounts indicate potential bugs in the dispute game contracts or the contracts that hold the bonds (`DelayedWETH`), which could compromise the incentivization mechanism and put user funds at risk.

## Technical Overview

### How It Works

1. **Monitoring Resolved Subgames**: The monitor listens for `resolveClaim` function calls on the dispute game contract, which indicate that a subgame has been resolved.

2. **Retrieving Claim Data**:

   - For each resolved subgame (identified by its claim index), the monitor retrieves the corresponding claim data from the dispute game contract.
   - The claim data includes information such as the claimant, any counterparty, and the original bond amount deposited for that claim.

3. **Determining Expected Recipient and Bond Amount**:

   - **Expected Recipient**:
     - If the claim was countered (i.e., there is a `counteredBy` address), the recipient of the bond should be the counterparty.
     - If the claim was not countered, the recipient is the original claimant.
   - **Expected Bond Amount**: The original bond amount associated with the claim, which should be unlocked upon resolution.

4. **Retrieving Bond Unlock Events**:

   - The monitor retrieves all `unlock` function calls on the `DelayedWETH` contract (the contract holding the bonds), which indicate bonds being unlocked and credited to participants.

5. **Cross-Referencing Values**:

   - The monitor cross-references the expected recipients and bond amounts with the actual bond unlocks recorded.
   - It checks that for each expected recipient and bond amount, there is a matching unlock event.

6. **Detecting Discrepancies**:

   - If any expected bond unlock is not found in the actual unlock events, it indicates a discrepancy.
   - Such a discrepancy means that the bond amount credited does not match the original bond value, signaling a potential issue.

7. **Triggering Alerts**:

   - If discrepancies are detected, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Ensuring Correct Incentivization**: The dispute game mechanism relies on proper financial incentives. Participants are motivated to act honestly because they risk losing their bonds if they behave maliciously. Over or under-accounting of ETH bonds undermines this mechanism.

- **Preventing Financial Loss**: Discrepancies in bond amounts can lead to financial losses for participants who do not receive the correct bond value upon resolution.

- **Maintaining Trust in the System**: Accurate bond accounting is essential for participants to trust the dispute resolution process and continue participating, which is crucial for the fault proof system's functionality.

- **Detecting Critical Bugs**: Discrepancies may indicate bugs or logic flaws in the dispute game contracts or related components, requiring prompt attention to prevent further issues.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.