[<img src="./docs/resources/bloxstaking_header_image.png" >](https://www.bloxstaking.com/)

<br>
<br>

# SSV - Secret Shared Validator

[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://pkg.go.dev/github.com/ethereum/eth2-ssv?tab=doc)
![Github Actions](https://github.com/ethereum/eth2-ssv/actions/workflows/full-test.yml/badge.svg?branch=stage)
![Github Actions](https://github.com/ethereum/eth2-ssv/actions/workflows/lint.yml/badge.svg?branch=stage)
![Test Coverage](./docs/resources/cov-badge.svg)
[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/eDXSP9R)

[comment]: <> ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/ethereum/eth2-ssv&#41;]&#40;https://goreportcard.com/report/github.com/ethereum/eth2-ssv&#41;)

[comment]: <> ([![Travis]&#40;https://travis-ci.com/ethereum/eth2-ssv.svg?branch=stage&#41;]&#40;https://travis-ci.com/ethereum/eth2-ssv&#41;)

## Introduction

Secret Shared Validator ('SSV') is a unique technology that enables the distributed control and operation of an Ethereum validator.

SSV uses an MPC threshold scheme with a consensus layer on top ([Istanbul BFT](https://arxiv.org/pdf/2002.03613.pdf)), 
that governs the network. \
Its core strength is in its robustness and fault tolerance which leads the way for an open network of staking operators 
to run validators in a decentralized and trustless way.

## SSV Spec
This repo contains the spec for SSV.Network node.

### SSV messages
SSV network message is called SSVMessage, it includes a MessageID and MsgType to route messages within the SSV node code, and, data for the actual message (QBFT/ Post consensus messages for example).

Any message data struct must be signed and nested within a signed message struct which follows the MessageSignature interface. 
A signed message structure includes the signature over the data structure, the signed root and signer list.

### Signing messages
The KeyManager interface has a function to sign roots, a slice of bytes. 
The root is computed over the original data structure (which follows the MessageRoot interface), domain and signature type.

**Use ComputeSigningRoot and ComputeSignatureDomain functions for signing**

Domain Constants:

| Domain         | Value                         | Description                       |
|----------------|-------------------------------|-----------------------------------|
| Primus Testnet | DomainType ("primus_testnet") | Domain for the the Primus testnet |

Signature type Constants:

| Signature Type       | Value                | Description                              |
|----------------------|----------------------|------------------------------------------|
| QBFT Signature       | [] byte {1, 0, 0, 0} | SignedMessage specific signatures        |
| PostConsensusSigType | [] byte {2, 0, 0, 0} | PostConsensusMessage specific signatures |

## Validator and DutyRunner instances
A validator instance is created for each validator independently, each validator will have multiple DutyRunners for each beacon chain duty type (Attestations, Blocks, etc.)
Duty runners are responsible for processing incoming messages and act upon them, completing a full beacon chain duty cycle.

CanStartNewDuty returns true if a new QBFT instance can start (meaning a new duty can get processed). 
As a general rule, new duties can't start until a full duty cycle (see below) is completed.  
One exception of the above is if a QBFT consensus decided, not all post consensus signatures were collected but 'PostConsensusSigCollectionSlotTimeout' slots passed.\
CanStartNewDuty Constants:

| Constant                              | Value | Description                                                                                                                       |
|---------------------------------------|-------|-----------------------------------------------------------------------------------------------------------------------------------|
| PostConsensusSigCollectionSlotTimeout | 32    | How many slots pass until a new QBFT instance can start without waiting for all post consensus partial signatures to be collected |

New Duty Full Cycle:

-> Received new beacon chain duty\
&nbsp;&nbsp;&nbsp;-> Check can start a new consensus instance\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Come to consensus on Duty + Duty data (AttestationData, etc.)\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Broadcast and collect partial signature to reconstruct signature\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-> Reconstruct signature, broadcast to BN

A duty runner holds a QBFT controller for processing QBFT messages and a dutyExecutionState which keeps progress for post consensus messages.
Partial signatures are collected and reconstructed (when threshold reached) to be broadcasted to the BN network.

## Validator Share
A share is generated and broadcasted publicly when a new SSV validator is registered to its operators.
Shares include: 
- Node ID: The Operator ID the share belongs to
- Validator Public Key
- Committee: An array of Nodes that constitute the SSV validator committee. A node must include it's NodeID and share public key.
- Domain

## Node
A node represents a registered SSV operator, each node has a unique ID and encryption key which is used to encrypt assigned shares.
NodeIDs are extremely important as they are used when splitting a validator key via Shamir-Secret-Sharing, later on they are used to verify messages and reconstruct signatures.

Shares use the Node data (for committee) to verify that incoming messages were signed by a committee member

## TODO
- [ ] Message Encoding - chose an encoding protocol and implement\
