package router

import (
	"net/http"

	h "../handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{Name: "register", Method: "POST", Pattern: "/authorization/v1/register", HandlerFunc: h.Register},
	Route{Name: "payment", Method: "POST", Pattern: "/payment/v1/wallet", HandlerFunc: h.Payment},
	Route{Name: "payment", Method: "POST", Pattern: "/payment/v1/trade", HandlerFunc: h.Trade},
	Route{Name: "payment", Method: "POST", Pattern: "/payment/v1/transfer", HandlerFunc: h.Transfer},
	Route{Name: "payment", Method: "POST", Pattern: "/payment/v1/offline", HandlerFunc: h.Offline},
	Route{Name: "payment", Method: "POST", Pattern: "/payment/v1/nolimit", HandlerFunc: h.Nolimit},
	// Route{Name: "ETH", Method: "POST", Pattern: "/ETH", HandlerFunc: h.ETH},
	// Route{Name: "BTC", Method: "POST", Pattern: "/BTC", HandlerFunc: h.BTC},
	// Route{Name: "GPE", Method: "POST", Pattern: "/GPE", HandlerFunc: h.GPE},
	// Route{Name: "USDT", Method: "POST", Pattern: "/USDT", HandlerFunc: h.USDT},
	// Route{Name: "WALLET", Method: "POST", Pattern: "/WALLET", HandlerFunc: h.Wallet},
}
