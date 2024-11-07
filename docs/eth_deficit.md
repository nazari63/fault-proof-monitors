## Purpose

The `eth_deficit.gate` monitor ensures the **safety** and integrity of the fault proof system by detecting any deficits of ETH in the `DelayedWETH` contract associated with the specified dispute game. The `DelayedWETH` contract holds the ETH bonds deposited by participants in dispute games. A deficit indicates that more ETH has been withdrawn than should be allowed, potentially due to bugs in bond accounting or dispute game resolution, which can result in financial losses for honest participants. Note this monitor tracks deficits in relation to the Challenger (`cbChallenger`), who is assumed to be operating honestly and participating in every dispute game as necessary.

## Technical Overview

### How It Works

1. **Retrieving Key Balances and Credits**:

   - **DelayedWETH Address**: Retrieves the `DelayedWETH` contract address associated with the specific `disputeGame`.

   - **Challenger's Claim Credit (`claimCredit`)**: The amount of ETH that the Challenger (`cbChallenger`) can currently claim from the `disputeGame`.

   - **Challenger's Total Credit (`totalCredit`)**: The total amount of ETH that has been unlocked for the `cbChallenger` in the `DelayedWETH` contract.

   - **Dispute Game's ETH Balance (`ethBalanceDisputeGame`)**: The total amount of ETH held in the `DelayedWETH` contract for the `disputeGame`.

2. **Validating Balances**:

   - **Credit Consistency**:
     - Ensures that the `claimCredit` (what the challenger can claim) does not exceed the `totalCredit` (what has been unlocked for them).
     - Verifies that the amount the challenger is trying to claim is consistent with what is available.

   - **Total Credit vs. Dispute Game Balance**:
     - Ensures that the `totalCredit` does not exceed the `ethBalanceDisputeGame`.
     - Checks that the total credits unlocked for participants do not exceed the ETH actually held for the dispute game.

   - **Synchronization Check**:
     - Ensures that if `claimCredit` is zero, then `totalCredit` should also be zero.
     - Detects discrepancies that might indicate desynchronization between the dispute game and the `DelayedWETH` contract.

3. **Triggering Alerts**:

   - If any of the above conditions fail, the monitor raises an alert indicating a potential deficit or inconsistency in the ETH balances related to the dispute game.

### Importance of the Monitor

- **Preventing Financial Loss**: A deficit in the `DelayedWETH` contract can lead to losses for honest participants expecting to receive their bonds back upon dispute resolution.

- **Ensuring Correct Bond Accounting**: Accurate tracking of bonds is crucial for the incentivization mechanism of dispute games. Over or under-accounting undermines trust and the proper functioning of the system.

- **Detecting Critical Issues Early**: Prompt identification of discrepancies allows for immediate investigation and correction of potential bugs in bond accounting or resolution logic.

- **Maintaining System Integrity**: Ensures that the dispute game mechanism operates securely, preserving the safety and reliability of the network.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.
- `cbChallenger`: Address of the Challenger.
