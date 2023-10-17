#!/bin/bash
source config.sh

./build/gnobot faucet $1 \
--pkgpath gno.land/r/gnoland/faucet \
--token $2 \
--bot-name "faucet" \
--channel "1029999999999999999" \
--guild "1019999999999999996" \
--chain-id $chainid \
--remote $rpc \
--limit 300000000ugnot \
--send   5000000ugnot
