#!/bin/bash

./build/gnobot faucet $1 \
--pkgpath gno.land/r/faucet \
--token 'ABC' \
--bot-name "faucet" \
--channel "1029999999999999999" \
--guild "1019999999999999996" \
--chain-id test3 \
--remote test3.gno.land:36657 \
--limit 300000000ugnot \
--send   5000000ugnot
