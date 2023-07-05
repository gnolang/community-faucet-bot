#!/bin/bash

# a pacakge deployed to the realm, with package the prefixe gno.land/r/ is a smart contract code with capability to persist state on chain
# --deposit flag: a minumn fee of 100gnot is needed to deploy the smart contract in the package path gno.land/r/
source config.sh

default=" $1  --gas-fee 10000000ugnot --gas-wanted 1000000 --broadcast true --remote  $rpc --chainid $chainid"
package=' --pkgdir examples/gno.land/r/gnoland/faucet --pkgpath gno.land/r/gnoland/faucet '

#faucet admin  contract holds the faucet fund. let's deposit 200K gnot first
deposit=' --deposit 200000000000ugnot'

# call set faucet account
callAddController=' --pkgpath gno.land/r/gnoland/faucet --func AdminAddController --args '
addController="$default $callAddController $4"


# call transfer
callTransfer=' --pkgpath gno.land/r/gnoland/faucet --func Transfer --args '
transfer="$default $callTransfer $4 --args $5"


gnotx='gnokey maketx'

if [ "$2" = "call" ]
then
  if [ "$3" = "controller" ]
  then
    if [ "$4" = "" ]
    then
    echo 'need the third parameter: controller_address'

    else
    echo $gnotx call $addController
    $gnotx call $addController
    fi

  elif [ "$3" = "transfer" ]
  then
    if [ "$4" = "" ]
    then
    echo 'need the third parameter: to_address'

    else
    echo $gnotx call $transfer
    $gnotx call $transfer

    fi
  else
  echo 'the second parameter must be controller or transfer'
  fi

elif [ "$2" = "deploy" ]
then
  # deploy contracts
  echo "$gnotx addpkg $defaut $package"
  $gnotx addpkg $default $package $deposit

else
  echo 'the fist parameter must be either call or deploy'
  echo 'usage: deploy | call transfer [toaddress] | call controller [address]'
fi
