package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gnolang/gno/pkgs/amino"
	"github.com/gnolang/gno/pkgs/command"

	"github.com/gnolang/gno/pkgs/crypto/keys"
	"github.com/gnolang/gno/pkgs/crypto/keys/client"
	"github.com/gnolang/gno/pkgs/errors"

	"github.com/gnolang/gno/pkgs/sdk/vm"
	"github.com/gnolang/gno/pkgs/std"
)

var perAccountLimit = std.NewCoin("ugnot", 350000000) // 350 gnot

// A valid gno address string length is 40
const addressStringLength = 40

type faucetOptions struct {
	client.BaseOptions        // home, ...
	ChainID            string `flag:"chain-id" help:"chain id must provide"`

	GasWanted int64  `flag:"gas-wanted" help:"gas requested for tx"`
	GasFee    string `flag:"gas-fee" help:"gas payment fee"`
	Send      string `flag:"send" help:"send coins per request"`
	Memo      string `flag:"memo" help:"any descriptive text"`
	PkgPath   string `flag:"pkgpath" help:"pacakge holding the faucet fund"`
	Limit     string `flag:"limit" help:"per transfer and per user account limit"`
	BotToken  string `flag:"token" help:"discord bot token - must provide"`
	BotName   string `flag:"bot-name" help:"discord bot name - must provide"`
	Channel   string `flag:"channel" help:"discord bot channel id - must provide"`
	Guild     string `flag:"guild" help:"discord bot guild/server id - must provide"`
}

var DefaultFaucetOptions = faucetOptions{
	BaseOptions: client.DefaultBaseOptions,
	ChainID:     "", // must override

	GasWanted: 800000,
	GasFee:    "1000000ugnot",
	Send:      "5000000ugnot",
	Memo:      "disccord faucet",
	PkgPath:   "gno.land/r/faucet",
	Limit:     "350000000ugnot",

	BotToken: "", // must override
	BotName:  "", // must override
	Channel:  "", // must override
	Guild:    "", // must override

}

// DiscordFaucet access local keybase, remote chain endpoint discord session
type DiscordFaucet struct {
	// local wallet key store
	keybase keys.Keybase

	keyinfo keys.Info

	keyname string

	keypass string

	sequence uint64

	// discord session
	session *discordgo.Session

	// other faucetOptions
	opts faucetOptions
}

// NewDiscordFaucet construct an instance
func NewDiscordFaucet(name, pass string, opts faucetOptions) (*DiscordFaucet, error) {
	kb, err := keys.NewKeyBaseFromDir(opts.Home)
	if err != nil {
		return nil, err
	}
	info, err := kb.GetByName(name)
	if err != nil {
		return nil, err
	}

	// sign a dummy value to validate key pass
	const dummy = "test"
	_, _, err = kb.Sign(name, pass, []byte(dummy))
	if err != nil {
		return nil, err
	}

	return &DiscordFaucet{
		keybase: kb,
		keyinfo: info,
		keyname: name,
		keypass: pass,

		opts: opts,
	}, nil
}

func faucetApp(cmd *command.Command, args []string, iopts interface{}) error {
	opts := iopts.(faucetOptions)

	if len(args) != 1 {
		cmd.ErrPrintfln("Usage: facuet <keyname>")
		return errors.New("invalid args")
	}

	if opts.ChainID == "" {
		return errors.New("chain-id not specified")
	}
	if opts.BotToken == "" {
		return errors.New("discord bot token not specified")
	}

	if opts.Channel == "" {
		return errors.New("discord bot channel id not specified")
	}

	if opts.BotName == "" {
		return errors.New("discord bot name not specified")
	}

	if opts.Guild == "" {
		return errors.New("discord bot guild/server id not specified")
	}

	remote := opts.Remote
	if remote == "" || remote == "y" {
		return errors.New("missing remote url")
	}

	var err error
	perAccountLimit, err = std.ParseCoin(opts.Limit)
	if err != nil {
		panic(err)
	}

	// XXX XXX
	// Read supply account pubkey.

	name := args[0]
	var pass string
	if opts.Quiet {
		pass, err = cmd.GetPassword("", false)
	} else {
		pass, err = cmd.GetPassword("Enter password.", false)
	}
	if err != nil {
		return err
	}

	// validate password

	// start a discord session

	df, err := NewDiscordFaucet(name, pass, opts)
	if err != nil {
		return err
	}
	// Start the bot

	err = df.Start()
	if err != nil {
		return err
	}

	// Cleanly close down the Discord session.
	defer df.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}

// call faucet contract to send tokens to requester
func (df *DiscordFaucet) sendAmountTo(to string, send std.Coins) error {
	function := "Transfer"
	// Read faucet account pubkey.

	faucet := df.keyinfo.GetAddress()

	// parse gas wanted & fee.
	gaswanted := df.opts.GasWanted

	gasfee, err := std.ParseCoin(df.opts.GasFee)
	if err != nil {
		panic(err)
	}
	sendAmount := strconv.FormatInt(send.AmountOf("ugnot"), 10)
	// construct msg & tx and marshal.
	msg := vm.MsgCall{
		Caller:  faucet,
		PkgPath: df.opts.PkgPath,
		Func:    function,
		Args:    []string{to, sendAmount},
	}
	tx := std.Tx{
		Msgs:       []std.Msg{msg},
		Fee:        std.NewFee(gaswanted, gasfee),
		Signatures: nil,
		Memo:       df.opts.Memo,
	}

	return df.signAndBroadcast(tx)
}

func (df *DiscordFaucet) signAndBroadcast(tx std.Tx) error {
	// query account to get account number and sequece

	accountAddr := df.keyinfo.GetAddress().String()

	qopts := client.QueryOptions{
		Path: fmt.Sprintf("auth/accounts/%s", accountAddr),
	}
	qopts.Remote = df.opts.Remote
	qres, err := client.QueryHandler(qopts)
	if err != nil {
		return errors.Wrap(err, "query account")
	}
	var qret struct{ BaseAccount std.BaseAccount }
	err = amino.UnmarshalJSON(qres.Response.Data, &qret)
	if err != nil {
		return err
	}

	// sign tx
	accountNumber := qret.BaseAccount.AccountNumber
	sequence := qret.BaseAccount.Sequence
	sopts := client.SignOptions{
		Sequence:      &sequence,
		AccountNumber: &accountNumber,
		ChainID:       df.opts.ChainID,
		NameOrBech32:  df.keyname,
		TxJson:        amino.MustMarshalJSON(tx),
	}
	sopts.Home = df.opts.Home
	sopts.Pass = df.keypass

	signedTx, err := client.SignHandler(sopts)
	if err != nil {
		return errors.Wrap(err, "sign tx")
	}

	// broadcast tx bytes.

	// broadcast signed tx
	bopts := client.BroadcastOptions{
		Tx: signedTx,
	}
	bopts.Remote = df.opts.Remote
	bres, err := client.BroadcastHandler(bopts)
	if err != nil {
		return errors.Wrap(err, "broadcast tx")
	}
	if bres.CheckTx.IsErr() {
		return errors.Wrap(bres.CheckTx.Error, "check transaction failed: log:%s", bres.CheckTx.Log)
	}
	if bres.DeliverTx.IsErr() {
		return errors.Wrap(bres.DeliverTx.Error, "deliver transaction failed: log:%s", bres.DeliverTx.Log)
	}

	fmt.Println("Message Delivered!")
	fmt.Println("GAS WANTED:", bres.DeliverTx.GasWanted)
	fmt.Println("GAS USED:  ", bres.DeliverTx.GasUsed)

	if string(bres.DeliverTx.Data) != "(\"\" string)" {
		return errors.New(string(bres.DeliverTx.Data))
	}

	return nil
}
