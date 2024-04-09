# Axelmon

`Axelmon` is a monitoring tool designed specifically for validators on the Axelar network. With `Axelmon`, validators can effortlessly verify the health and performance of both `vald` and `tofnd` nodes. Additionally, this tool promptly alerts operators in case any issues arise with these nodes.

### Contents

- [Axelmon](#axelmon)
  - [Monitoring List](#monitoring-list)
  - [Features](#features)
  - [Quick Guide](#quick-guide)
- [Deep dive](#deep-dive)
  - [Background](#background)
  - [What does it do?](#what-does-it-do)

## Monitoring List

- Heartbeat
  - Monitors `vald` node sent properly `heartbeat tx` to the Axelar network when `heartbeat event` emitted.
- Maintainer
  - Monitors the registration status of maintainers on the Axelar network.
- EVM Vote
  - Monitors votes for registered external chains(Maintainers).

## Features

- JSON API
  - Access monitoring results via a JSON API, available at the following path: `/`.
- Prometheus
  - Also support Prometheus for providing monitoring results. You can access them at the following path: `/metrics`.
- Telegram Alerts
  - Receive alerts via Telegram if any issues arise with your node.

## Quick Guide

1. Build

```bash
go build
```

2. Configure `config.toml` file

```bash
cp config.toml.example config.toml
```

1. Run

```bash
./axelmon -config config.toml
```

# Deep dive

## Background

Axelar network is a chain that mainly services bridge, and it is important for validator teams to operate nodes.

Axelar network validators will run following nodes:

- `axelard`
  - This is the main component that communicates with other nodes via P2P and secures the network’s consensus via Tendermint BFT.
  - Axelar network itself is preparing sign batches that can be relayed to appropriate EVM compatible chains.
- `vald`
  - It listens to events from the Axelar network, such as signing requests, and verifying events on EVM chains. It’s connected to EVM compatible RPC nodes of external chains, where it queries for events such as deposit confirmations and message send.
  - On one side it is connected to the Axelar network, where it reads instructions what to verify. On the other side, it’s connected to EVM compatible RPC nodes of external chains, where it verifies that batches were executed correctly from Axelar network.
- `tofnd`
  - gRPC service which wraps rust implementation of multi-party signing protocols, such as, ECDSA multisig, GG20 threshold-ECDSA protocol (t of n), into RPC calls invoked by vald.Signs transaction batches for sending to smartcontracts on destination chains
  - It has to be executed in highly secured environment.

The bridge service of Axelar network relies on the `vald` and `tofnd` nodes. Any issues with these nodes can significantly impact quality and security. Swift problem recognition is crucial to avoid disruptions. Once an issue is identified, operators investigate by analyzing logs and collaborating with other validator teams or the Foundation on Discord.

This is why we developed `Axelmon`, which helps us for swiftly recognize issues with the `vald` and `tofnd` nodes.

## **What does it do?**

![diagram](docs/diagram.png)

`Axelmon` is a monitoring tool designed to oversee nodes related to bridge service of Axelar network. Let’s delve into what `Axelmon` monitors:

- Heartbeat
  - Heartbeat events emit at regular intervals, specifically when $CurrentHeight \mod 50 == 0$.
  - Within 1-2 blocks from this event emitted height, the validator’s **`vald`** node must send a `heartbeat transaction` to the Axelar network.
  - `Axelmon` monitors the miss count for `n` heartbeats.
- Maintainer
  - The Maintainer is the external chains that have chosen to support on the Axelar network.
  - It could be automatically deregistered, if a validator misses many votes.
  - `Axelmon` monitors whether the maintainer is registered or not.
- EVM Vote
  - Whenever a bridge service-related transaction occurs for a registered external chain, a vote is cast.
  - `Axelmon` monitors the miss count for n votes at each registered chains.

And `Axelmon` provide the result of monitoring:

- JSON API ( `localhost:${LISTENPORT}/` )

  ```json
  {
    "maintainers": {
      "status": true,
      "maintainer": {
        "Avalanche": true,
        "Ethereum": true
      }
    },
    "heartbeat": {
      "status": true,
      "missed": "0 / 3"
    },
    "EVMVotes": {
      "chain": {
        "Avalanche": {
          "status": true,
          "missed": "0 / 10"
        },
        "Ethereum": {
          "status": true,
          "missed": "0 / 10"
        }
      }
    }
  }
  ```

  - maintainers
    - status: If any of the mainainers are deregistered, the state will be `false`.
    - maintainer: `true` if registered, `false` if deregistered.
  - heartbeat
    - status: `false` if it is gte to the threshold filled in `miss_cnt` field of `config.toml`.
    - missed: Miss count for heartbeats.
      - (`miss_cnt` field) / (`check_n` field) in `config.toml`.
  - EVMVotes
    - status: `false` if it is gte to the threshold filled in `miss_cnt` field of `config.toml`.
    - missed: Miss count for votes.
      - (`miss_cnt` field) / (`check_n` field) in `config.toml`.

- Prometheus ( `localhost:${LISTENPORT}/metrics` )
  ```
  # HELP evm_votes_total Number of EVM votes
  # TYPE evm_votes_total counter
  evm_votes_total{network_name="Avalanche",status="missed"} 0
  evm_votes_total{network_name="Avalanche",status="success"} 10
  # HELP heartbeats_total Heartbeat of the application
  # TYPE heartbeats_total counter
  heartbeats_total{status="missed"} 0
  heartbeats_total{status="success"} 3
  # HELP maintainers_status_list Maintainer's status, 1 is active, 0 is not active
  # TYPE maintainers_status_list gauge
  maintainers_status_list{network_name="Avalanche"} 1
  ```
- Notification through Telegram
  - The alert condition corresponds to the condition for `false` status in the API.
