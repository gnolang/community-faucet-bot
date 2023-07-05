#!/bin/bash
source config.sh

# uncomment the deployment section if the faucet contract needs to be published on the chain.
# We assume the faucet contract is loaded to the chain from local directory examples/gno.land/r
#
# echo "deploy contract..."
# ./faucet.sh $fundingkey deploy



# The faucet contract holds the fund, not the controller, but  we need to fund the controller and admin accounts so that it can pay the gas fee
# when it sends instructions to the faucet contract

# uncomment the fund admin account section when you run it on the production server.

# echo "fund admin account from test1 ..."
# ./send.sh $fundingkey g1u7y667z64x2h7vc6fmpcprgey4ck233jaww9zq 80000000ugnot


echo "fund controller1 account from $fundingkey and you will be asked to enter $fundingkey pass code"
./send.sh $fundingkey g1q0pjk6dd6lehd5q3gcpghvcf3rd6mqy7tge4va 8000000000ugnot

echo

echo "add controller1 from $adminkey and you will be asked to enter $adminkey pass code"
./faucet.sh $adminkey call controller g1q0pjk6dd6lehd5q3gcpghvcf3rd6mqy7tge4va
