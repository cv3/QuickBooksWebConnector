package qbwcGo

import (
	"encoding/json"
	"encoding/xml"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
	"github.com/amazingfly/cv3go"
)

//UpdateCV3Items will recieve a qb itemQueryRs with items to be converted and sent to CV3
func UpdateCV3Items(items ItemQueryRs) error {
	Log.Debug("Starting qbwcUpdateCV3Items")
	var itemInventoryMapPath = "./fieldMaps/itemInventoryMapping.json"
	var itemServiceMapPath = "./fieldMaps/itemServiceMapping.json"
	var itemNonInventoryMapPath = "./fieldMaps/itemNonInventoryMapping.json"
	var checkSKU []string
	var b []byte
	var err error
	var cv3Items = cv3go.Products{}
	//Check every case, using fallthrough
	switch {
	case len(items.ItemInventoryRets) > 0:
		for _, inventory := range items.ItemInventoryRets {
			ItemMapping(&inventory, &checkSKU, &cv3Items, itemInventoryMapPath)
		}
		fallthrough
	case len(items.ItemNonInventoryRets) > 0:
		for _, nonInventory := range items.ItemNonInventoryRets {
			ItemMapping(&nonInventory, &checkSKU, &cv3Items, itemNonInventoryMapPath)
		}
		fallthrough
	case len(items.ItemServiceRets) > 0:
		for _, service := range items.ItemServiceRets {
			ItemMapping(&service, &checkSKU, &cv3Items, itemServiceMapPath)
		}
		fallthrough
	default:
	}

	var total = len(cv3Items.Products) //total number of items
	var block = 10                     //size of block to send
	var max = 0                        //maximum index
	//split product add calls up into blocks
	for min := 0; max != total; min = max {
		if min == 0 { //first pass
			if block <= total { //total is greater than the block size
				max = block - 1 // set max index to size of block -1 to keep min and max from overlapping
			} else { //total is smaller then the block size
				max = total //set max to total
			}
		} else { //min != 0
			if max+block <= total { //total is greater than the new max index
				max += block //add a block size to the maximum index
			} else { //total is less then the old max plus a new block
				max = total //set the maximum index to the total
			}
		}
		b, err = xml.Marshal(cv3Items.Products[min:max])
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("error marshaling cv3 products inUpdateCV3Items ")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("error marshaling cv3 products inUpdateCV3Items ")
		}
		api := cv3go.NewApi()
		api.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID)***REMOVED***
		//api.Debug = true
		api.PushInventory(string(b), false)
		b = api.Execute()
	}
	Log.WithFields(logrus.Fields{
		"ItemInventoryRets":         len(items.ItemInventoryRets),
		"ItemNonInventoryRets":      len(items.ItemNonInventoryRets),
		"ItemServiceRets":           len(items.ItemServiceRets),
		"ItemGroupRets":             len(items.ItemGroupRets),
		"ItemDiscountRets":          len(items.ItemDiscountRets),
		"ItemOtherChargeRets":       len(items.ItemOtherChargeRets),
		"ItemSubtotalRets":          len(items.ItemSubtotalRets),
		"cv3Items.Products created": len(cv3Items.Products),
	}).Debug("returns received")
	//set config's LastUpdate to now
	cfg.ItemUpdates.LastUpdate = time.Now().Format(time.RFC3339)
	SaveConfig()
	return nil
}

//ItemMapping set the cv3 product fields with the QB data
//Load the field mapping, convert itemInterface into a gabs container, then use the fieldMap to set the desired fields
func ItemMapping(itemInterface interface{}, checkSKU *[]string, cv3Items *cv3go.Products, mapPath string) {
	//Load the field mapping json file into a map[string]string
	var fieldMap = ReadFieldMapping(mapPath)
	//Marshal itemInterface into json to convert to gabs container
	jsonBytes, err := json.Marshal(itemInterface)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "mapPath": mapPath}).Error("Error marshaling itemInterface in ItemMapping")
		ErrLog.WithFields(logrus.Fields{"Error": err, "mapPath": mapPath}).Error("Error marshaling itemInterface in ItemMapping")
	} //convert json to gabs
	qbItem, err := gabs.ParseJSON(jsonBytes)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "mapPath": mapPath}).Error("Error parsing json to gabs container in ItemMapping")
		ErrLog.WithFields(logrus.Fields{"Error": err, "mapPath": mapPath}).Error("Error parsing json to gabs container in ItemMapping")
	}
	*checkSKU = append(*checkSKU, CheckPath("ListID", qbItem))
	var itemTemp = cv3go.Product{}

	itemTemp.Sku = fieldMap["SKU"].Display(qbItem)
	itemTemp.Name = fieldMap["Name"].Display(qbItem)
	//Check if this is an inventory item
	if strings.Contains(mapPath, "itemInventory") {
		itemTemp.InventoryControl.InventoryOnHand = fieldMap["InventoryControl.InventoryOnHand"].Display(qbItem)
		onOrder, err := strconv.ParseInt(fieldMap["InventoryControl.OnOrder"].Display(qbItem), 0, 32)
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error parsing int for onOrder")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error parsing int for onOrder")
		}
		itemTemp.InventoryControl.OnOrder = int(onOrder)
		qtyOnHand, err := strconv.ParseInt(CheckPath("QuantityOnHand", qbItem), 0, 32)
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error parsing int for QuantityOnHand")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error parsing int for QuantityOnHand")
		}
		if qtyOnHand > 0 {
			itemTemp.InventoryControl.InventoryStatus = "In Stock"
		}
	} else { //nonInventory or service items
		itemTemp.InventoryControl.InventoryControlExempt = "true"
	} //end inventory if
	if strings.ToUpper(CheckPath("IsActive", qbItem)) == "TRUE" {
		itemTemp.Inactive = "false"
	}
	itemTemp.Description = fieldMap["Description"].Display(qbItem)
	itemTemp.Retail.Price.StandardPrice = fieldMap["Retail.Price.StandardPrice"].Display(qbItem)
	itemTemp.Retail.Price.PriceCategory = "Retail"
	itemTemp.Retail.Active = "true"

	cv3Items.Products = append(cv3Items.Products, itemTemp)
}

/*
	if len(items.ItemGroupRets) > 0 {
		for i, group := range items.ItemGroupRets {
			checkSKU = append(checkSKU, group.ListID)
			var itemTemp = cv3go.Product{}
			itemTemp.Sku = group.ListID
			if group.FullName != "" {
				itemTemp.Name = group.FullName
			} else {
				itemTemp.Name = group.Name
			}

			if strings.ToUpper(group.IsActive) == "TRUE" {
				itemTemp.Inactive = "false"
			} else {

			}
			itemTemp.Description = group.ItemDesc
			for _, item := range group.ItemGroupLInes {
				var itemSubTemp = cv3go.SubProduct{}
				itemSubTemp.Name = item.ItemRef.FullName
				itemSubTemp.Sku = item.ItemRef.ListID

				itemSubTemp.InventoryControl.InventoryOnHand = item.Quantity
				itemTemp.SubProducts.SubProducts = append(itemTemp.SubProducts.SubProducts, itemSubTemp)
			}
			itemTemp.SubProducts.Active = "true"
			itemTemp.SubProducts.SubProducts[i].Sku = group.ItemGroupLInes[i].ItemRef.ListID
			itemTemp.Retail.Price.PriceCategory = "Retail"
			itemTemp.InventoryControl.InventoryControlExempt = "true"
			itemTemp.Retail.Active = "true"
			itemTemp.Categories.IDs = append(itemTemp.Categories.IDs, "15")
			cv3Items.Products = append(cv3Items.Products, itemTemp)
		}
	}
*/
