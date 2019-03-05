package qbwcGo

import (
	"github.com/amazingfly/cv3go"
)

//ConfirmCV3Order will let CV3 know the order information has been received
func ConfirmCV3Order(orderID string) {
	if cfg.ConfirmOrders {
		//Call CV3 for the desired order confirmation
		var api = cv3go.NewApi()
		//api.Debug = true
		api.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID)***REMOVED***
		api.OrderConfirm(orderID)
		api.Execute()
	}
}
