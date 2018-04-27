#!/bin/bash

echo "set protocol address..."
read protocol

echo "set name"
read name

go run cmd/lrc/* nameRegistry registerName --name=$name --sender="0x4bad3053d574cd54513babe21db3f09bea1d387d" --config=/Users/fukun/projects/gohome/src/github.com/Loopring/relay/config/debug.toml --gasPrice="100000000000" --protocolAddress=$protocol
