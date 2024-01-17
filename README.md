# Axelmon

Axelmon is monitoring program for Axelar chain validator operator.

## Monitor List

- Heartbeat
- Maintainer
- EVM Vote

## Quick Guide

1. Build

```bash
go build
```

2. Configure config.toml file

```bash
# You can get a example of config.toml file by below command
cp config.toml.example config.toml
```

3. Execute

```bash
./axelmon -config config.toml
`
```
