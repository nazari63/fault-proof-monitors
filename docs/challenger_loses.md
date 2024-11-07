## Purpose

The `challenger_loses.gate` monitor ensures the **safety** of the fault proof system by detecting situations where the Challenger (`honestChallenger`), assumed to be operating honestly, loses a dispute game or any subgame within the dispute game that they should have won. Since the `honestChallenger` acts in good faith, a game indicates a critical issue—possibly a logic flaw in the dispute game contracts or offchain components—that could compromise the integrity and security of the network.

## Technical Overview

### How It Works

1. **Monitoring Dispute Resolutions**: The monitor observes dispute games and its subgames involving the `honestChallenger` and checks for the resolution statuses.

2. **Identifying Challenger's Role**:
   
   - **As a Challenger**: When the `honestChallenger` initiates a challenge against an incorrect state output root.
   
   - **As a Defender**: When the `honestChallenger` defends a correct state output root against challenges.

3. **Analyzing Outcomes**:
   
   - The monitor retrieves the `Resolved` events from the dispute game contract to determine the outcome:
     - **Defender Wins (Status 1)**: Indicates the defender won the dispute.
     - **Challenger Wins (Status 2)**: Indicates the challenger won the dispute.

   - It also examines all moves (subgames) made in the dispute game to identify the actions taken by the `honestChallenger`.

4. **Detecting Anomalies**:
   
   - **Challenger Loses as Challenger**: If the `honestChallenger` initiated the challenge but the dispute resolves with the defender winning, this suggests an incorrect outcome.
   
   - **Challenger Loses as Defender**: If the `honestChallenger` defended a correct claim but the dispute resolves with the challenger winning, this indicates a problem.

   - **Lost Subgames**: The monitor also checks if the `honestChallenger` lost any subgames within the dispute game, which should not happen if the challenger is acting honestly.

5. **Triggering Alerts**:

   - If any of the above anomalies are detected, the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Ensuring Safety**: The fault proof system must correctly resolve disputes to maintain the integrity of the network. An honest challenger losing a dispute they should have won compromises the system's safety.

- **Detecting Critical Bugs**: Such an event may indicate serious logic flaws in the dispute game contracts or issues with offchain components such as the challenger software.

- **Protecting User Funds**: Incorrect dispute resolutions could potentially put user funds at risk by allowing invalid state transitions or withdrawals.

- **Maintaining Trust**: Ensures that the network operates securely and that the mechanisms in place to prevent fraud are functioning correctly.

## Parameters

- `honestChallenger`: Address of the Challenger.
- `disputeGame`: Address of the dispute game contract being monitored.

