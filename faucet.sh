#!/bin/bash

# a pacakge deployed to the realm, with package the prefixe gno.land/r/ is a smart contract code with capability to persist state on chain
# --deposit flag: a minumn fee of 100gnot is needed to deploy the smart contract in the package path gno.land/r/

default=' test1  --gas-fee 10000000ugnot --gas-wanted 1000000 --broadcast true --remote test3.gno.land:36657 --chainid test3'
package=' --pkgdir examples/gno.land/r/faucet --pkgpath gno.land/r/faucet '

#faucet admin  contract holds the faucet fund. let's deposit 200K gnot first
deposit=' --deposit 200000000000ugnot'

# call set faucet account
callAddController=' --pkgpath gno.land/r/faucet --func AdminAddController --args '
addController="$default $callAddController $3"


# call transfer
callTransfer=' --pkgpath gno.land/r/faucet --func Transfer --args '
transfer="$default $callTransfer $3 --args $4"


gnotx='gnokey maketx'

if [ "$1" = "call" ]
then
  if [ "$2" = "controller" ]
  then
    if [ "$3" = "" ]
    then
    echo 'need the third parameter: controller_address'

    else
    echo $gnotx call $addController
    $gnotx call $addController
    fi

  elif [ "$2" = "transfer" ]
  then
    if [ "$3" = "" ]
    then
    echo 'need the third parameter: to_address'

    else
    echo $gnotx call $transfer
    $gnotx call $transfer

    fi
  else
  echo 'the second parameter must be controller or transfer'
  fi

elif [ "$1" = "deploy" ]
then
  # deploy contracts
  echo "$gnotx addpkg $defaut $package"
  $gnotx addpkg $default $package $deposit

else
  echo 'the fist parameter must be either call or deploy'
  echo 'usage: deploy | call transfer [toaddress] | call controller [address]'
fi
