#!/bin/bash

echo "deploy contract..."
./faucet.sh deploy

echo "fund admin account..."
./send.sh test1 g14ykc8d4n2sr9lmlv80cgp2qu74emujyujvhe5w 80000000ugnot
echo "fund controller1 account..."
./send.sh test1 g1q0pjk6dd6lehd5q3gcpghvcf3rd6mqy7tge4va 80000000ugnot

echo "add controller1"
./faucet.sh call controller g1q0pjk6dd6lehd5q3gcpghvcf3rd6mqy7tge4va
