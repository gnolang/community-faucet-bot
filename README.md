# community-faucet-bot
the community bot(s) interacting with gno.land/r/gnoland/faucet

See https://github.com/gnolang/gno/issues/364.

## GNO Discord Faucet

Community members can request testing tokens from the discord channel.

## Problem to solve

A centralized faucet endpoint could easily be an abusive target.
In addition, it adds operational complexity and cost to operating an endpoint.

### Details

There are four issues we try to solve in the web-based faucet.

1) It is centralized and can be an abuse target.

2) First-time users need to understand the fee structure to register and interact with the board and require multiple requests to get in fees

3) There need to be more tokens allocated to the faucet wallet to support many user registrations and board creation.

4) Faucet need to use fund from a contract instead of a user account.

## Solution

A decentralized discord bot can be a faucet to any discord server and channels.

Discord provides sophisticated verification to prevent abuse. We can leverage it to distribute test tokens in the community group.

It is friendly to community members. The moderator can easily manage and prevent people from abusing the faucet.

### Details

A discord faucet bot can solve 1, 2, and partially 3

1) A discord bot can run on any computer and be configured to use any wallet to support the faucet. There is no direct attacking point as a service. Instead, the testing token request is funneled through the discord app, which the admin moderates. Multiple bots can stand by and respond to the same channel as back to each other.

2) In the gno land test net, users must request 200gnot to register and cost 100gnot to interact with the board contract if not registered.

How do we allow people quickly get started and prevent people keep accumulating testing tokens?

We can set a limit for each, say 400gnot per account. We give max 400gnot for first-time users. So they can get started the right way and have enough tokens to try everything. For non-first-time users, we provide a regular amount of 1gnot to cover most of the gas fee usage.

3) A faucet can drain out with regular usage. We can run a backup bot with another faucet wallet monitor on a different channel as a backup.

## Other aspects

4) We can also recycle the tokens from the user and broad contract back to the faucet manually or automatically.

5) 200 gnot user registration fee is to prevent people from spamming the user board. We can add a limited number of users can register per day and lower the fee. It slows down the pace that the faucet drains out.

## Features

- Each account is limited to 350gnot max.
- Instead of incrementally issuing a small amount to users, we give the first-time user the max amount. It costs 200gnot to register and 100gnot to create a board in the current test net. No need to request multiple times to be able to register users and create boards.

## Instruction

#### 0) Install the gnokey

the discord bot will need to access local key store

     git clone https://github.com/gnolang/gno
     cd gno
     make install_gnokey

make sure you include $GOPATH/bin in $PATH

#### 0.1) create admin and controller accounts

We use test1 to fund admin and controllers. The admin account address is hard coded in the faucet.gno.
We could set the contract deployer as faucet contract admin. However, it is less secure.

     // only use test1 for testing
     gnokey add test1 --recover

     // import faucet admin key for testing, on production admin should have been already created for it to be included in the contract.

     gnokey add admin --recover

     gnokey add controller1



#### 1) build gnobot

    make

#### 2) Assign controller to the faucet

Please modify the controller address in ./provision.sh

On production, please use a separate admin account than test1 ( g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5 )
DO NOT use test1 as admin in production, and mnemonic code is published on github.
adminkey's address should match with the admin address in faucet.gno

     ./provision.sh

#### 3) start the bot

Review the flag in the script, and provide a gnokey account that can instruct the faucet contract to distribute funds.

      ./startbot.sh controller1 DISCORDBOT_TOKEN


the executed command and flags


    ./build/gnobot faucet controller1 --chain-id test3 --token DISCORDBOT_TOKEN --channel DISCORD_CHANNELID -bot-name DISCORD_BOTNAME --guild $DISCORD_GUILD --remote test3.gno.land:36657


The following flags are required when you run the gnobot. We do not recommend storing the discord token on your local machine, not even in the env file.

--chain-id

--token

--channel

--bot-name

--guild

--limit // default 350000000ugnot

--send  // default 5000000ugnot

DISCORDBOT_TOKEN:

Your discord application bot token. (It is NOT a OAuth2 token )

DISCORD_CHANNELID:

The id of a channel that you add the bot to

DISCORD_GUILD:

The server id of the discord server that you added the bot

DISCORD_BOTNAME:

The name of your bot


#### 4) Final step, fund the faucet contract!

Once the bot starts, you will see a contract package address printed.

    Please verify the address with the contract deployed on the chain for  gno.land/r/gnoland/faucet
    Please make sure to deposit funds to this faucet contract: g1ttrq7mp4zy6dssnmgyyktnn4hcj3ys8xhju0n7

./send.sh test1 g1ttrq7mp4zy6dssnmgyyktnn4hcj3ys8xhju0n7 200000000000ugnot
