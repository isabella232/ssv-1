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

## TODO
- [ ] Message Encoding - chose an encoding protocol and implement\
