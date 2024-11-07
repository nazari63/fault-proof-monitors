# Fault Proof Monitors for Base

Fault proofs are a crucial component in ensuring the security and decentralization of **Base**. They enable a permissionless and decentralized approach to proposing and validating the state of the chain, allowing anyone to participate in securing the network. With fault proofs now live on Base Mainnet, monitoring these systems is essential to maintain the integrity and trustlessness of the network.

This repository contains a collection of monitors written in **gate**, a descriptive language developed by Hexagate for defining invariants for real-time monitoring of every block. These monitors are designed to observe and detect anomalies within the fault proof system on Base, ensuring that any potential issues are identified and addressed promptly.

By continuously monitoring fault proofs on Base, we aim to enhance the security and reliability of the chain, fostering a more open and decentralized onchain economy for everyone.

## Monitors

Each monitor in this repository targets specific aspects of the fault proof system to detect potential safety and liveness issues. Detailed explanations, including the purpose, technical overview, and importance of each monitor, are provided in the `docs` folder.

### List of Monitors

- `challenged_proposal.gate`
- `challenger_loses.gate`
- `credit_and_bond_discrepancy.gate`
- `duplicate_dispute_game.gate`
- `eth_deficit.gate`
- `eth_withdrawn_early.gate`
- `fault_proof_detection_child.gate`
- `fault_proof_detection_parent.gate`
- `incorrect_bond_balance.gate`
- `unresolvable_dispute_game.gate`

Please refer to the `docs` folder for comprehensive documentation on each monitor.

