package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"

	"github.com/gnolang/gno/pkgs/amino"

	"github.com/gnolang/gno/pkgs/crypto"

	"github.com/gnolang/gno/pkgs/crypto/keys/client"
	"github.com/gnolang/gno/pkgs/errors"
	"github.com/gnolang/gno/pkgs/std"
)

// Start a discord session
func (df *DiscordFaucet) Start() error {
	err := df.discordFaucet()
	if err != nil {
		return err
	}

	// Open a websocket connection to Discord and begin listening.

	return df.session.Open()
}

// Close discord session smoothly
func (df *DiscordFaucet) Close() {
	df.session.Close()
}

// it takes discord server API token and rpc client.
// We use rpc client to validate addresses and check abuses.
// and we cap the balance holding to 400 GNOT before issue new tokens from faucet

func (df *DiscordFaucet) discordFaucet() error {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + df.opts.BotToken)
	if err != nil {
		fmt.Println("failed to create discord bot session.", err)
		return err
	}

	df.session = dg

	// we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This is a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}
		// ignore message from other channels
		if m.ChannelID != df.opts.Channel {
			return
		}

		// ignore message from other discord guild/server

		if m.GuildID != df.opts.Guild {
			return
		}

		// ignore message does not mention anyone
		if len(m.Mentions) == 0 {
			return
		}

		user := m.Mentions[0]

		// ingore messages that doe not mentiont this bot

		if user.Bot != true || user.Username != df.opts.BotName {
			return
		}

		var res string
		// retrive toAddress from received disocrd message
		toAddr, bal, err := df.process(m)
		if err != nil {

			res = fmt.Sprintf("%s", err)
			dg.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> "+res)

			fmt.Printf("channel %s:<@%s> %s", m.ChannelID, m.Author.Username, res)

			return

		}
		var send std.Coins
		// if the account has no balance we give it upto the full.
		if bal.IsZero() {
			send = std.NewCoins(perAccountLimit)
		} else {
			send, _ = std.ParseCoins(df.opts.Send)
		}

		err = df.sendAmountTo(toAddr, send)
		if err != nil {

			dg.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> "+"faucet failed")

			fmt.Printf("channel %s:<@%s>%v", m.ChannelID, m.Author.Username, err)

			return

		}

		var amount string
		for _, v := range send {
			amount += v.String() + " "
		}

		res = fmt.Sprintf("Cha-Ching! %s +%s", toAddr, amount)

		dg.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+"> "+res)
	})

	return nil
}

// This function will be called every time a new
// message is created on any channel that the authenticated bot has access to.
// It returns true and valid response if the message from discord is valid.
// Other we returns false with an empty string,which we should ingore.
func (df *DiscordFaucet) process(m *discordgo.MessageCreate) (string, std.Coin, error) {
	validAddr, err := retrieveAddr(m.Content)

	zero := std.Coin{}

	if err != nil {
		return "", zero, errors.New("No valid addresse: %s", err)
	}

	bal, err := df.checkBalance(validAddr)
	if err != nil {
		return "", zero, errors.New("It seems our faucet is not working properly %s", err)
	}

	if bal.IsGTE(perAccountLimit) {
		return "", zero, errors.New("Your account %s still has %d%s, no need to accumulate more testing tokens", validAddr, bal.Amount, bal.Denom)
	}

	return validAddr, bal, nil
}

// retrive the first string with valid addresse length
func retrieveAddr(message string) (string, error) {
	var addr string
	words := strings.Fields(message)

	for _, w := range words {

		s := strings.TrimFunc(w, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if len(s) == addressStringLength {
			addr = s
			break
		}
	}

	ok, err := isValid(addr)

	if !ok {
		return "", err
	}

	return addr, nil
}

// A valid address format

func isValid(addr string) (bool, error) {
	// validate prefix

	if strings.HasPrefix(addr, crypto.Bech32AddrPrefix) == false {
		return false, errors.New("The address does not have correct prefix %s", crypto.Bech32AddrPrefix)
	}

	if addressStringLength != len(addr) {
		return false, errors.New("The address does not have correct length %d", addressStringLength)
	}

	_, err := crypto.AddressFromBech32(addr)
	if err != nil {
		return false, errors.Wrap(err, "parsing address")
	}
	return true, nil
}

// return true if the account balance is within limit
// return amount of token from
func (df *DiscordFaucet) checkBalance(addr string) (std.Coin, error) {
	qopts := client.QueryOptions{
		Path: fmt.Sprintf("auth/accounts/%s", addr),
	}
	qopts.Remote = df.opts.Remote
	qres, err := client.QueryHandler(qopts)
	if err != nil {
		return std.Coin{}, errors.Wrap(err, "query account")
	}
	var acc struct{ BaseAccount std.BaseAccount }
	err = amino.UnmarshalJSON(qres.Response.Data, &acc)
	if err != nil {
		return std.Coin{}, err
	}

	balances := acc.BaseAccount.GetCoins()

	bal := std.Coin{Denom: "ugnot", Amount: 0}

	if len(balances) > 0 {
		for i, v := range balances {
			if v.Denom == "ugnot" {
				bal = balances[i]
			}
		}
	}

	return bal, nil
}
