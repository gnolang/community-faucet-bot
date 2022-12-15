# community-faucet-bot
the community bot(s) interacting with gno.land/r/gnoland/faucet

See https://github.com/gnolang/gno/issues/364.

## GNO Discord Faucet

Community members can request testing tokens from the discord channel.

## Problem to solve

A centralized faucet endpoint could easily be an abusive target.
In addition, it adds operational complexity and cost to operating an endpoint.

### Details

There are three issues we try to solve in the web-based faucet.

1) It is centralized and can be an abuse target.

2) First-time users need to understand the fee structure to register and interact with the board and require multiple requests to get in fees

3) There need to be more tokens allocated to the faucet wallet to support many user registrations and board creation.



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

5) 200 gnot user registration fee is to prevent people from spamming the user board. We can add a limited number of users can register per day and lower the fee. It slows down the pace that the faucets drain out.

## Features

- Each account is limited to 350gnot max.
- Instead of incrementally issuing a small amount to users, we give the first-time user the max amount. It costs 200gnot to register and 100gnot to create a board in the current test net. No need to request multiple times to be able to register users and create boards.

## Instruction

#### 0) Install the gnokey

     git clone https://github.com/gnolang/gno
     cd gno
     make install_gnokey

make sure you include $GOPATH/bin in $PATH

#### 0.1) create admin and controller accounts

     gnokey add admin

     gnokey add controller1

     gnokey add controller2


#### 1) build gnobot

    make

#### 2) Deploy faucet contract and assign controller the faucet

 please modify the the address in the script.

     ./provison.sh

#### 3) start the bot

check out

      ./startbot.sh

The following flags are required when you run the gnobot. We do not recommend storing the discord token on your local machine, not even in the env file.

--chain-id

--token

--channel

--bot-name

--guild

--limit // default 350000000ugnot

--send  // default 5000000ugnot

DISCORDBOT_TOKEN:

Your discord application bot token. (NOT OAuth2 token )

DISCORD_CHANNELID:

The id of a channel that you add the bot to

DISCORD_GUILD:

The id of the discord server that you add the bot to

DISCORD_BOTNAME:

The name of your bot


      ./build/gnodiscord faucet test1 --chain-id test3 --token DISCORDBOT_TOKEN --channel DISCORD_CHANNELID -bot-name DISCORD_BOTNAME --guild $DISCORD_GUILD --remote test3.gno.land:36657
