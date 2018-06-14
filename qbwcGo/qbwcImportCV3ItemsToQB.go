package qbwcGo

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/amazingfly/cv3go"
)

//ImportCV3ItemsToQB will recieve a cv3 product and add a qb itemQueryRs to the work queue
func ImportCV3ItemsToQB(prod cv3go.Product) {
	//Add to wait group, to let sendResponse know there is more work to be done
	//wg.Add(1)
	var templateBuff = bytes.Buffer{}
	var escapedQBXML = bytes.Buffer{}
	var err error
	var itemInventoryAdd = ItemInventoryAdd{}
	var itemServiceAdd = ItemServiceAdd{}
	//if inventory controll is disabled, use itemServiceAdd
	if strings.ToLower(prod.InventoryControl.InventoryControlExempt) == "true" {
		itemServiceAdd.Name = EscapeName(prod.Name)
		itemServiceAdd.SalesOrPurchase.Desc = prod.Description
		itemServiceAdd.SalesOrPurchase.Price = prod.Retail.Price.StandardPrice
		itemServiceAdd.IsActive = "true"
		itemServiceAdd.SalesOrPurchase.AccountRef.FullName = "Shipping and Delivery Income"

		//build and return the template
		var tPath = `./templates/qbItemServiceAdd.t`

		LoadTemplate(&tPath, itemServiceAdd, &templateBuff)
		err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")

		}
		//Send prepared QBXML to the workChan
		workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: itemServiceAdd, Type: "ItemServiceAddRs"}
		escapedQBXML.Reset()
		templateBuff.Reset()
	} else {
		itemInventoryAdd.Name = EscapeName(prod.Name)
		itemInventoryAdd.SalesDesc = EscapeField(prod.Description)
		itemInventoryAdd.SalesPrice = prod.Retail.Price.StandardPrice
		itemInventoryAdd.IsActive = "true"
		itemInventoryAdd.QuantityOnHand = prod.InventoryControl.InventoryOnHand

		//COGS must be set with a valid QB account
		itemInventoryAdd.COGSAccountRef.FullName = "Sales - Software"
		//IncomeAccount must be set with a valid QB account
		itemInventoryAdd.IncomeAccountRef.FullName = "Sales - Software"
		//AssetAccount must be set with a valid QB account
		itemInventoryAdd.AssetAccountRef.FullName = "Sales - Software"

		//My qb may not allow ClassRef
		//itemInventoryAdd.ClassRef.FullName = prod.Sku
		//itemInventoryAdd.BarCode
		//itemInventoryAdd.ManufacturerPartNumber
		//itemInventoryAdd.ExternalGUID
		//itemInventoryAdd.InventoryDate
		//itemInventoryAdd.Max
		//itemInventoryAdd.ParentRef
		//itemInventoryAdd.PrefVendorRef
		//itemInventoryAdd.PurchaseCost
		//itemInventoryAdd.PurchaseDesc
		//itemInventoryAdd.ReorderPoint
		//itemInventoryAdd.SalesTaxCodeRef
		//itemInventoryAdd.TotalValue
		//itemInventoryAdd.UnitOfMeasureSetRef

		//Add subproducts as their own product in QB
		for _, subProd := range prod.SubProducts.SubProducts {
			var temp = ItemInventoryAdd{}
			temp.Name = EscapeName(subProd.Name)
			temp.SalesDesc = EscapeField(prod.Description)
			temp.SalesPrice = subProd.Retail.Price.StandardPrice
			//temp.ClassRef.FullName = subProd.Sku

			//COGS must be set with a valid QB account
			temp.COGSAccountRef.FullName = "Sales - Software"
			//IncomeAccount must be set with a valid QB account
			temp.IncomeAccountRef.FullName = "Sales - Software"
			//AssetAccount must be set with a valid QB account
			temp.AssetAccountRef.FullName = "Sales - Software"
			temp.IsActive = "true"
			temp.QuantityOnHand = subProd.InventoryControl.InventoryOnHand

			//build and return the template
			var tPath = `./templates/qbItemInventoryAdd.t`

			LoadTemplate(&tPath, temp, &templateBuff)
			err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
			if err != nil {
				Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")
			}
			//Send prepared QBXML to the workChan
			workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: temp, Type: "ItemInventoryAddRs"}
			escapedQBXML.Reset()
			templateBuff.Reset()
		}
		//build and return the template
		var tPath = `./templates/qbItemInventoryAdd.t`

		LoadTemplate(&tPath, itemInventoryAdd, &templateBuff)
		err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")
		}
		//Send prepaired QBXML to the workChan
		workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: itemInventoryAdd, Type: "ItemInventoryAddRs"}
	}

}

//EscapeField takes a string and returns an escaped string
func EscapeField(s string) string {
	var escaped = bytes.Buffer{}
	//XML escape the passed in string
	err := xml.EscapeText(&escaped, []byte(s))
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error escaping text in EscapeField")
	}
	return escaped.String()
}

//EscapeName will check the name field for sizee requirements then escape it
func EscapeName(s string) string {
	//QBXML Name fields cannot be > 31, and must be XML escaped
	if len(s) > 31 {
		return EscapeField(s[:31])
	}
	return EscapeField(s)

}
