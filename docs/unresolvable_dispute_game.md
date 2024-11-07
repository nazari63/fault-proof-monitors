## Purpose

The `unresolvable_dispute_game.gate` monitor ensures the **liveness** and integrity of the fault proof system by detecting dispute games that remain unresolved beyond a defined timeframe. In a properly functioning system, dispute games should resolve within a specific period after their creation. If a dispute game remains unresolved past this period, it may indicate potential issues in the dispute resolution process, which can prevent participants from reclaiming their ETH bonds and disrupt the processing of withdrawals associated with the disputed state root.

## Technical Overview

### How It Works

1. **Fetching Key Timestamps**:

   - **Creation Timestamp (`createdAt`)**: Retrieves the timestamp when the dispute game was created.

   - **Game Duration (`maxClockDuration`)**: Retrieves the maximum allowed duration for each participant's clock in the dispute game. This duration typically reflects the maximum time a participant has to make a move in the dispute game.

   - **Resolved Timestamp (`resolvedAt`)**: Retrieves the timestamp when the dispute game was resolved. If the game is unresolved, this value is zero.

2. **Calculating Expected Resolution Time**:

   - **Expected Resolution Timestamp**: Calculates an expected allowed time for the dispute game to be resolved by adding twice the game duration (to account for both the challenger and defender's clocks) and an additional buffer time (`extraTimeInSeconds`) to the creation timestamp.

  - **Formula**:
    ```
    expectedResolutionTimestamp = createdAt + (2 * maxClockDuration) + extraTimeInSeconds
    ```

  - **Explanation**:

    - **`2 * maxClockDuration`**: Accounts for the maximum time both the challenger and defender could take, assuming each uses the full duration allowed for their moves.

    - **`extraTimeInSeconds`**: An additional buffer time set to **172,800 seconds (2 days)**. This buffer accounts for any potential clock extensions or delays that might extend the dispute game beyond the standard duration. Extensions can occur due to the specific game mechanics that allow participants extra time under certain circumstances.

3. **Current Time Check**:

   - **Current Block Timestamp**: Retrieves the current block timestamp to compare against the expected resolution timestamp.

4. **Validating Resolution Status**:

   - The monitor checks whether:

     - The dispute game has been resolved (`resolvedAt` is not zero), **or**

     - The current time is less than or equal to the expected resolution timestamp.

   - **Alert Condition**: If neither condition is met, meaning the dispute game remains unresolved beyond the expected resolution time (including the extra buffer of 2 days), the monitor raises an alert.

### Importance of `extraTimeInSeconds`

- **Purpose**: The `extraTimeInSeconds` parameter is set to **172,800 seconds (2 days)** to provide an additional buffer on top of the calculated maximum game duration. This accounts for:

  - **Clock Extensions**: Certain dispute games may include mechanisms that extend the game duration under specific conditions, such as participants requesting more time or network delays.

  - **Ensuring Accuracy**: The buffer helps prevent false positives in the monitoring system by allowing for legitimate extensions beyond the standard maximum duration.

- **Effect on Monitoring**: By including this buffer, the monitor waits an additional 2 days beyond the expected duration before raising an alert. This ensures that only genuinely unresolvable dispute games (those that have exceeded all expected timeframes) or dispute games that are extending past the typical resolution window trigger an alert, reducing unnecessary investigations into games that are still within a reasonable extended duration.

## Importance of the Monitor

- **Ensuring Liveness**: Timely resolution of dispute games is crucial for the network to progress and the fault proof system's reliability.

- **Preventing Financial Loss**: Unresolved dispute games prevent participants from receiving their ETH bonds and rewards, leading to potential financial losses.

- **Maintaining Withdrawal Functionality**: The resolution status is required for processing withdrawals in systems like `OptimismPortal2`. An unresolved dispute game can block withdrawals associated with its output root claim.

- **Detecting Critical Issues**: Prolonged unresolved dispute games may indicate bugs or malfunctions in the dispute resolution contracts or offchain components, requiring immediate attention.

## Parameters

- `disputeGame`: Address of the dispute game contract being monitored.

- `extraTimeInSeconds`: An integer representing additional time added to the expected resolution date to account for potential clock extensions or network delays.
