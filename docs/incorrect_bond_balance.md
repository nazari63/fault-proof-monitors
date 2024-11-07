# incorrect_bond_balance.gate

## Purpose

The `incorrect_bond_balance.gate` monitor ensures the **safety** and integrity of the fault proof system on the **Base** by verifying that the total ETH bonds associated with a dispute game are correctly accounted for during the resolution process. Specifically, it checks that the sum of the ETH bonds for resolved claims and unresolved claims equals the maximum amount of ETH that the dispute game has ever held. Any imbalance indicates potential issues such as skipped subgames, multiple resolutions, or under-accounted claims, which can compromise the incentivization mechanism and result in financial losses for participants.

## Technical Overview

### How It Works

1. **Monitoring Resolved Claims**:

   - Retrieves all historical `resolveClaim` function calls on the `disputeGame` contract to identify resolved claims (subgames).
   - Extracts the claim indices of resolved claims.

2. **Calculating Resolved Bonds**:

   - Retrieves the bond amounts associated with each resolved claim.
   - Sums up the bond amounts for all resolved claims.

3. **Identifying Unresolved Claims**:

   - Determines the range of claim indices for unresolved claims based on the minimum resolved claim index.
   - Calculates the expected bond amounts for unresolved claims.

4. **Calculating Total ETH Held by the Dispute Game**:

   - Retrieves the current ETH balance of the `disputeGame` in the `DelayedWETH` contract.
   - Sums the current ETH balance with any past ETH withdrawals (unlocks) to determine the maximum amount of ETH that the dispute game has ever held.

5. **Validating Bond Accounting**:

   - Verifies that the sum of the bond amounts for resolved and unresolved claims equals the total ETH that has ever been held by the dispute game.
   - Checks for any discrepancies, which would indicate an imbalance.

6. **Triggering Alerts**:

   - If an imbalance is detected, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Ensuring Correct Incentivization**: Accurate bond accounting is essential for the incentivization mechanism of dispute games. Under-accounting or over-accounting of bonds undermines trust and the proper functioning of the system.

- **Preventing Financial Loss**: Imbalances can lead to financial losses for participants who are entitled to receive bond amounts upon dispute resolution.

- **Detecting Critical Issues**: Discrepancies may indicate that subgames have been skipped, resolved out of order, or resolved multiple times, pointing to broader issues in the dispute game logic.

- **Maintaining System Integrity**: Ensures that the dispute resolution process operates securely and reliably, preserving the safety of the Base network.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.

-