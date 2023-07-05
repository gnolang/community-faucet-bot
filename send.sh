#!/bin/bash
source config.sh

echo send from $1 to $2 for $3
gnokey maketx send $1 --to $2 --send $3  \
--gas-wanted 50000 --gas-fee 100000ugnot \
--remote $rpc \
--chainid $chainid \
--broadcast true
