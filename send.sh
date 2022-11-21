#!/bin/bash
echo send - from to amount
gnokey maketx send $1 --to $2 --send $3  \
--gas-wanted 50000 --gas-fee 100000ugnot \
--remote test3.gno.land:36657 \
--chainid test3 \
--broadcast true
