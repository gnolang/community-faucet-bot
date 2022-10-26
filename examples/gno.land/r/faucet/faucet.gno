package faucet

import (
	"errors"
	"std"
	"strings"

	"gno.land/p/ufmt"
)

var (
	// configurable by admin.
	gAdminAddr   std.Address    = "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5"
	gControllers [3]std.Address = [3]std.Address{} // limit it to 3

	gInPause = false
	gMessage = "# Community Faucet.\n\n"

	// internal vars, for stats.
	gTotalTransferred std.Coins
	gTotalTransfers   = uint(0)

	// per account request limit 350 gnot
	limit std.Coins = std.Coins{std.Coin{"ugnot", 350000000}}
)

func Transfer(to std.Address, send int64) string {
	if err := assertIsController(); err != nil {
		return err.Error()
	}

	if gInPause {
		return errors.New("faucet in pause").Error()
	}

	var sendCoins std.Coins

	if send <= limit.AmountOf("ugnot") {

		sendCoins = std.Coins{std.Coin{Denom: "ugnot", Amount: send}}

	} else {

		return errors.New("Per request limit " + limit.String() + " exceed").Error()
	}

	gTotalTransferred = gTotalTransferred.Add(sendCoins)
	gTotalTransfers++

	banker := std.GetBanker(std.BankerTypeRealmSend)
	pkgaddr := std.GetOrigPkgAddr()
	banker.SendCoins(pkgaddr, to, sendCoins)
	return ""
}
func GetPerTransferLimit() int64 {

	return limit.AmountOf("ugnot")

}
func Render(path string) string {
	banker := std.GetBanker(std.BankerTypeRealmSend)
	balance := banker.GetCoins(std.GetOrigPkgAddr())

	output := path + gMessage
	if gInPause {
		output += "Status: inactive.\n"
	} else {
		output += "Status: active.\n\n"
	}
	output += ufmt.Sprintf("Balance: %s.\n", balance.String())
	output += ufmt.Sprintf("Total transfer: %s (in %d times).\n\n", gTotalTransferred.String(), gTotalTransfers)

	output += "Package address: " + std.GetOrigPkgAddr().String() + "\n\n"
	output += ufmt.Sprintf("Admin: %s\n\n ", gAdminAddr.String())
	output += ufmt.Sprintf("Controller:\n\n ")

	for _, v := range gControllers {

		output += ufmt.Sprintf("%s  ", v.String())

	}
	output += ufmt.Sprintf("\n\n")
	output += ufmt.Sprintf("Per request limit: %s\n\n", limit.String())

	return output
}

func assertIsController() error {
	caller := std.GetOrigCaller()

	for _, v := range gControllers {
		if caller == v {
			return nil
		}
	}
	return errors.New(caller.String() + " is not on the controller list")
}
