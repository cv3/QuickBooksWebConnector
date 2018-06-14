package qbwcGo

import "github.com/amazingfly/cv3go"

//InitQBItems adds certian items to Quick Books, to prepare it for the qbwc/cv3 intigration
func InitQBItems() {
	var shipping = cv3go.Product{}
	shipping.Description = "Shipping charges"
	shipping.Inactive = "false"
	shipping.Name = "Shipping"
	shipping.InventoryControl.InventoryControlExempt = "true"
	ImportCV3ItemsToQB(shipping)
}
