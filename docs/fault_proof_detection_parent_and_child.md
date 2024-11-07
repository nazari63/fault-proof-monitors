# Fault Proof Detection Parent

## Purpose

The `fault_proof_detection_parent.gate` monitor checks the **correctness** of L2 output proposals submitted to **Ethereum**. With fault proofs now allowing anyone to propose output roots permissionlessly, it is crucial to verify that these proposals are valid. This monitor detects when an invalid L2 output proposal is submitted, enabling timely action to challenge and prevent incorrect state roots from affecting the network.

Additionally, this monitor is unique in that it requires inputs from both Ethereum and the specified L2 in order to verify output proposals.

## Technical Overview

### How It Works

1. **Monitoring Dispute Game Creation**:

   - The monitor listens for `DisputeGameCreated` events emitted by the `DisputeGameFactoryProxy` contract.
   - Each event contains information about a newly created dispute game, including the `disputeProxy` address, `gameType`, and the submitted `l2OutputProposal` (root claim).

2. **Fetching L2 Block Data**:

   - **Block Number**: Retrieves the L2 block number associated with the dispute game by calling `l2BlockNumber()` on the `disputeProxy` contract.
   - **State Root**: Fetches the state root of the L2 block using the `StateRoot` function with the block number and L2 chain ID.
   - **Block Hash**: Obtains the L2 block hash using the `BlockHash` function.
   - **Message Passer Storage Hash**: Retrieves the storage root hash of the L2 Cross Domain Message Passer contract at a specific block number.

3. **Computing Expected L2 Output Proposal**:

   - The monitor computes the expected L2 output proposal hash by concatenating:

     - A zero bytes32 value.
     - The L2 state root.
     - The message passer storage hash.
     - The L2 block hash.

   - The concatenated data is then hashed using `Keccak256` to produce the expected output root.

4. **Comparing Proposals**:

   - The monitor compares the computed expected L2 output root with the `l2OutputProposal` submitted in the `DisputeGameCreated` event.
   - If they do not match, it indicates that an incorrect output root has been submitted.

5. **Triggering Alerts**:

   - If an incorrect output proposal is detected, the monitor raises an alert for immediate action.
   - Additionally, it ensures that only one `DisputeGameCreated` event appears per block.

### Importance of the Monitor

- **Ensuring Correctness**: Verifies that only valid L2 output root proposals are submitted, maintaining the integrity of the L2 state rolled up to Ethereum.
- **Preventing Security Risks**: Detects invalid proposals that could lead to incorrect state being propagated to Ethereum, potentially causing severe security issues, including loss of funds or compromised state integrity.
- **Supporting Permissionless Proposals**: With the network allowing anyone to submit output roots, this monitor is essential to help safeguard against fraudulent or faulty proposals.

## Parameters

- `disputeGameFactoryProxy`: Address of the `DisputeGameFactoryProxy` contract used to monitor dispute game creations.
- `l2ChainId`: The chain ID of the L2 network to perform cross-chain calls and fetch L2 block data.

# Fault Proof Detection Child

## Purpose

The `fault_proof_detection_child.gate` monitor helps ensure that incorrect L2 output proposals detected by the `fault_proof_detection_parent.gate` monitor are challenged in time by the Challenger (`cbChallenger`). Once an invalid output root is identified, this monitor is deployed to track the corresponding dispute game, verifying that the challenger is actively contesting the invalid proposal to prevent it from affecting the network.

## Technical Overview

### How It Works

1. **Monitoring Dispute Game Moves**:

   - The monitor listens for `Move` events emitted by the specified `disputeGame` contract.
   - Each `Move` event represents a claim (move) made in the dispute game, and includes the `parentIndex`, `claim`, and `claimant` (the address making the move).

2. **Retrieving Claims**:

   - **Claim Count**: Retrieves the total number of claims (`claimDataLen`) in the dispute game.
   - **Parent Indices**:

     - **Even Parent Indices**: Represent challenge moves (attacks) against the parent claim.
     - **Odd Parent Indices**: Represent defense moves supporting the parent claim.

3. **Analyzing Challenger's Actions**:

   - **Challenge Moves by `cbChallenger`**:

     - Filters `Move` events where the `claimant` is the `cbChallenger` and the `parentIndex` is even, indicating that the challenger is attacking (challenging) the invalid output root.

   - **Defense Moves**:

     - Filters `Move` events where the `parentIndex` is odd, indicating defense moves. If an attacker is defending an invalid output root, it is a concern.

4. **Triggering Alerts**:

   - **Challenger Is Challenging**:

     - As the `cbChallenger` attacks the invalid output proposal, the monitor will fire alerts for each attack move, enabling us to track each move made.

   - **Attacker Defending**:

     - If any participant is making defense moves to support the invalid output root, an alert is raised.

### Importance of the Monitor

- **Ensuring Timely Challenges**: Verifies that the `cbChallenger` is actively challenging invalid output roots, preventing them from being accepted by the network.
- **Monitoring Adversarial Behavior**: Detects if external participants are defending invalid output roots, which could indicate coordinated attempts to compromise the network.
- **Complementing Parent Monitor**: Works in tandem with the `fault_proof_detection_parent.gate` monitor to provide a comprehensive detection and alert mechanism for invalid output proposals.

## Parameters

- `cbChallenger`: Address of the Challenger responsible for contesting invalid output roots.
- `disputeGame`: Address of the dispute game contract associated with the invalid output proposal as detected by the parent monitor.
