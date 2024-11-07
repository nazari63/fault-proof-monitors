## Purpose

The `duplicate_dispute_game.gate` monitor ensures the **integrity** of the fault proof system on the **Base** network by detecting the creation of duplicate dispute games for the same state root claim. In the system, only one dispute game should exist per state root claim. Creating multiple dispute games for the same claim can lead to inconsistent dispute outcomes and ambiguity in which dispute game the system should reference, potentially compromising the dispute resolution process.

## Technical Overview

### How It Works

1. **Monitoring Dispute Game Creation**: The monitor observes the `DisputeGameFactory` contract for any `create` function calls that initiate new dispute games.

2. **Collecting Dispute Game Data**:

   - **DisputeGameFactory Address**: Retrieved from the `OptimismPortalProxy` contract.
   - **Current Respected Game Type**: Obtained from the `OptimismPortalProxy` contract to focus on the relevant game type.
   - **DisputeGameCreated Events**: All historical `DisputeGameCreated` events are collected, including associated block numbers.
   - **Extra Data**: For each created dispute game, the `extraData` field is retrieved, which is necessary for identifying unique games.

3. **Calculating Unique Identifiers (UUIDs)**:

   - For each dispute game, a UUID is calculated using the `getGameUUID` function of the `DisputeGameFactory` contract.
   - The UUID is derived from the combination of `gameType`, `rootClaim`, and `extraData`.

4. **Detecting Duplicates**:

   - **Existing UUIDs**: A mapping of UUIDs from previously created dispute games (excluding those in the current block) is constructed.
   - **New UUIDs**: UUIDs of dispute games created in the current block are calculated.
   - **Comparison**: The monitor checks if any new UUIDs match existing ones, indicating a duplicate dispute game.
   - Additionally, it checks for duplicates within the new dispute games themselves.

5. **Triggering Alerts**:

   - If a duplicate UUID is found—meaning a dispute game with the same `gameType`, `rootClaim`, and `extraData` already exists—the monitor raises an alert for immediate investigation.

### Importance of the Monitor

- **Preventing Inconsistencies**: Multiple dispute games for the same state root claim can lead to conflicting dispute outcomes, undermining the reliability of the dispute resolution process.
- **System Integrity**: Ensures that the system references the correct dispute game, avoiding ambiguity and potential security vulnerabilities.
- **Maintaining Protocol Rules**: Enforces the protocol's rule that only one dispute game can exist per state root claim.

## Parameters

- `optimismPortalProxy`: Address of the `OptimismPortalProxy` contract used to retrieve the `DisputeGameFactory` contract address and the respected game type.