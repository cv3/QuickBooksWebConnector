package qbwcGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
	"github.com/amazingfly/cv3go"
)

var addrCharLimit = 41
var cityCharLimit = 31
var stateCharLimit = 21
var zipCharLimit = 13
var countryCharLimit = 31
var noteCharLimit = 41

//GetCV3Orders will recieve a qb itemQueryRs with items to be converted and sent to CV3
func GetCV3Orders() { //(workChan chan string, doneChan chan bool) {
	var workCount = 0
	var workCTX = WorkCTX{}
	Log.Debug("staring getcv3orders")

	//Call CV3 for the desired orders
	var api = cv3go.NewApi()
	//api.Debug = true
	api.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID)***REMOVED***
	api.GetOrdersNew()
	//api.GetOrdersRange("7152", "7152") //"7142")
	var d = api.Execute(true)
	Log.Debug(string(d))
	ord, err := gabs.ParseJSON(d)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "json": string(d)}).Error("Error parsing json into gabs container in GetCV3Order")
		ErrLog.WithFields(logrus.Fields{"Error": err, "json": string(d)}).Error("Error parsing json into gabs container in GetCV3Order")
	}
	ord = ord.Path("CV3Data.orders")
	Log.Debug(ord.String())

	switch strings.ToLower(cfg.OrderType) {
	case "salesreceipt":
		MakeSalesReceipt(&workCount, &workCTX, ord)
		break
	case "salesorder":
		MakeSalesOrder(&workCount, &workCTX, ord)
		break
	default:
		Log.WithFields(logrus.Fields{"OrderType": cfg.OrderType}).Error("Error in GetCV3Orders, invalid order type in config")
		ErrLog.WithFields(logrus.Fields{"OrderType": cfg.OrderType}).Error("Error in GetCV3Orders, invalid order type in config")
	}

	if workCount < 1 {
		workChan <- WorkCTX{Work: ""}
		getLastErrChan <- "No new Orders"
	}
}

//MakeSalesReceipt takes the cv3 order and turns it into a qbxml salesReceiptAdd
func MakeSalesReceipt(workCount *int, workCTX *WorkCTX, ordersMapper *gabs.Container) {
	//Prepare gabs container for range loop
	oMapper, err := ordersMapper.Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "OrdersMapper": ordersMapper}).Error("Error getting ordersMapper Children in MakeSalesReceipt")
		ErrLog.WithFields(logrus.Fields{"Error": err, "OrdersMapper": ordersMapper}).Error("Error getting ordersMapper Children in MakeSalesReceipt")
	} //Load the dynamic field mappings froma  file
	var fieldMap = ReadFieldMapping("./fieldMaps/receiptMapping.json")
	//iterate over the orders, then the shiptos, as each shipto needs to be handled as a seperate sales receipt in QB
	for _, o := range oMapper {
		//prepare the shipTo gabs container for a range loop
		shipToMapper, err := o.Path("shipTos").Children()
		if err != nil {
			Log.WithFields(logrus.Fields{"Error": err, "ShipTosMapper": o.Path("shipTos")}).Error("Error getting shipTosMapper Children in MakeSalesReceipt")
			ErrLog.WithFields(logrus.Fields{"Error": err, "ShipTosMapper": o.Path("shipTos")}).Error("Error getting shipTosMapper Children in MakeSalesReceipt")
		}
		//start racking the orders and shipTos
		orderSuccessChan <- OrderSuccessTracker{
			OrderID:      CheckPath("orderID", o),
			ShipToLength: len(shipToMapper),
		}
		for shipToIndex, shipTo := range shipToMapper {
			var qbReceiptAdd = SalesReceiptAdd{} //object to hold current shipto information
			qbReceiptAdd.CV3OrderID = CheckPath("orderID", o)
			*workCount++

			//checkPayMethod will return the transactionID of the passed in PayMethod
			CheckPayMethod(o, "payMethod")
			CheckPayMethod(o, "additionalPayMethod")
			qbReceiptAdd.ClassRef.FullName = CheckPath(fieldMap["ClassRef.FullName"], o)
			qbReceiptAdd.ClassRef.ListID = CheckPath(fieldMap["ClassRef.ListID"], o)
			qbReceiptAdd.Other = CheckPath(fieldMap["Other"], o)
			if CheckPath(fieldMap["ExchangeRate"], o) != "" {
				exchRate, err := strconv.ParseFloat(CheckPath(fieldMap["ExchangeRate"], o), 64)
				if err != nil {
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing string to float for exchange rate in salesReceiptAdd")
				} else {
					qbReceiptAdd.ExchangeRate = exchRate
				}
			}
			if strings.ToLower(CheckPath(fieldMap["IsToPeEmailed"], o)) == "true" {
				qbReceiptAdd.IsToBeEmailed = true
			}
			if strings.ToLower(CheckPath(fieldMap["IsToPePrinted"], o)) == "true" {
				qbReceiptAdd.IsToBePrinted = true
			}
			if strings.ToLower(CheckPath(fieldMap["IsPending"], o)) == "true" {
				qbReceiptAdd.IsPending = "true"
			}
			qbReceiptAdd.FOB = CheckPath(fieldMap["FOB"], o)
			qbReceiptAdd.CustomerMsgRef.FullName = CheckPath(fieldMap["CustomerMsgRef"], shipTo)
			//qbReceiptAdd.CustomerMsgRef.ListID = CheckPath(fieldMap["CustomerMsgRef.ListID"], shipTo)
			qbReceiptAdd.CustomerSalesTaxCodeRef.FullName = CheckPath(fieldMap["CustomerSalesTaxCodeRef.FullName"], o)
			qbReceiptAdd.CustomerSalesTaxCodeRef.ListID = CheckPath(fieldMap["CustomerSalesTaxCodeRef.ListID"], o)
			qbReceiptAdd.ItemSalesTaxRef.FullName = CheckPath(fieldMap["ItemSalesTaxRef.FullName"], o)
			qbReceiptAdd.ItemSalesTaxRef.ListID = CheckPath(fieldMap["ItemSalesTaxRef.ListID"], o)
			qbReceiptAdd.DepositToAccountRef.FullName = CheckPath(fieldMap["DepositToAccountRef.FullName"], o)
			qbReceiptAdd.DepositToAccountRef.ListID = CheckPath(fieldMap["DepositToAccountRef.ListID"], o)

			//Direct mappingFor Beatrice Bakery "W"
			qbReceiptAdd.SalesRepRef.FullName = fieldMap["SalesRepRef.FullName"] //CheckPath(fieldMap["SalesRepRef.FullName"], o)

			qbReceiptAdd.SalesRepRef.ListID = CheckPath(fieldMap["SalesRepRef.ListID"], o)
			qbReceiptAdd.TemplateRef.FullName = CheckPath(fieldMap["TemplateRef.FullName"], o)
			qbReceiptAdd.TemplateRef.ListID = CheckPath(fieldMap["TemplateRef.ListID"], o)
			qbReceiptAdd.RefNumber = CheckPath(fieldMap["RefNumber"], o)
			qbReceiptAdd.ShipToIndex = shipToIndex

			//start billing information assignment
			//QB will either accept addr 1-5 or addr 1-2 and city state zip country
			var addr = make([]string, 0) // For adding Billing address info
			//If first or last name is not empty, add them as the first line
			if CheckPath("billing.firstName", o) != "" || CheckPath("billing.lastName", o) != "" {
				//if title is not empty add it before the name
				if CheckPath("billing.title", o) != "" {
					addr = append(addr, CheckPath("billing.title", o)+" "+CheckPath("billing.firstName", o)+" "+CheckPath("billing.lastName", o))
				} else {
					addr = append(addr, CheckPath("billing.firstName", o)+" "+CheckPath("billing.lastName", o))
				}
			} //if billing company is not empty ad it as the next available address slot
			if CheckPath("billing.company", o) != "" {
				addr = append(addr, CheckPath("billing.company", o))
			} //add billing address line 1 as the next available address slot
			addr = append(addr, CheckPath("billing.address1", o))
			//if billing address line 2 is not empty add it as the next available address slot
			if CheckPath("billing.address2", o) != "" {
				addr = append(addr, CheckPath("billing.address2", o))
			}
			//add the billing address info to the QB struct
			if len(addr) > 0 {
				qbReceiptAdd.BillAddress.Addr1 = FieldCharLimit(addr[0], addrCharLimit)
			}
			if len(addr) > 1 {
				qbReceiptAdd.BillAddress.Addr2 = FieldCharLimit(addr[1], addrCharLimit)
			}
			if len(addr) > 2 {
				qbReceiptAdd.BillAddress.Addr3 = FieldCharLimit(addr[2], addrCharLimit)
			}
			qbReceiptAdd.BillAddress.City = FieldCharLimit(CheckPath("billing.city", o), cityCharLimit)
			qbReceiptAdd.BillAddress.Country = FieldCharLimit(CheckPath("billing.country", o), countryCharLimit)
			qbReceiptAdd.BillAddress.PostalCode = FieldCharLimit(CheckPath("billing.zip", o), zipCharLimit)
			qbReceiptAdd.BillAddress.State = FieldCharLimit(CheckPath("billing.state", o), stateCharLimit)
			//qbReceiptAdd.BillAddress.Note not used
			//end billing information

			qbReceiptAdd.ShipMethodRef.FullName = CheckPath(fieldMap["ShipMethodRef.FullName"], shipTo)
			qbReceiptAdd.ShipMethodRef.ListID = CheckPath(fieldMap["ShipMethodRef.ListID"], shipTo)
			//hard code prefix for Beatrirce bakery
			qbReceiptAdd.Memo = "WEB# " + CheckPath(fieldMap["Memo"], o)
			qbReceiptAdd.PaymentMethodRef.FullName = CheckPath(fieldMap["PaymentMethodRef.FullName"], o)
			qbReceiptAdd.PaymentMethodRef.ListID = CheckPath(fieldMap["PaymentMethodRef.ListID"], o)

			//If the billing name is not paypal, use it as the customers name
			if !strings.Contains(strings.ToLower(CheckPath("billing.firstName", o)), "paypal") {
				qbReceiptAdd.CustomerRef.FullName = CheckPath("billing.lastName", o) + ", " + CheckPath("billing.firstName", o)
			} //else bliiling firstname is paypal, so do not add any customer info
			qbReceiptAdd.ShipDate = CheckPath(fieldMap["ShipDate"], shipTo)

			//Start shipping address
			addr = make([]string, 0) // For adding Billing address info
			//If first or last name is not empty, add them as the first line
			if CheckPath("firstName", shipTo) != "" || CheckPath("lastName", shipTo) != "" {
				//if title is not empty add it before the name
				if CheckPath("title", shipTo) != "" {
					addr = append(addr, CheckPath("title", shipTo)+" "+CheckPath("firstName", shipTo)+" "+CheckPath("lastName", shipTo))
				} else {
					addr = append(addr, CheckPath("firstName", shipTo)+" "+CheckPath("lastName", shipTo))
				}
			} //if shiping company is not empty ad it as the next available address slot
			if CheckPath("company", shipTo) != "" {
				addr = append(addr, CheckPath("company", shipTo))
			} //add Shiping address line 1 as the next available address slot
			addr = append(addr, CheckPath("address1", shipTo))
			//if shiping address line 2 is not empty add it as the next available address slot
			if CheckPath("address2", shipTo) != "" {
				addr = append(addr, CheckPath("address2", shipTo))
			}
			//add the shiping address info to the QB struct
			if len(addr) > 0 {
				qbReceiptAdd.ShipAddress.Addr1 = FieldCharLimit(addr[0], addrCharLimit)
			}
			if len(addr) > 1 {
				qbReceiptAdd.ShipAddress.Addr2 = FieldCharLimit(addr[1], addrCharLimit)
			}
			if len(addr) > 2 {
				qbReceiptAdd.ShipAddress.Addr3 = FieldCharLimit(addr[2], addrCharLimit)
			}
			qbReceiptAdd.ShipAddress.City = FieldCharLimit(CheckPath("city", shipTo), cityCharLimit)
			qbReceiptAdd.ShipAddress.State = FieldCharLimit(CheckPath("state", shipTo), stateCharLimit)
			qbReceiptAdd.ShipAddress.PostalCode = FieldCharLimit(CheckPath("zip", shipTo), zipCharLimit)
			qbReceiptAdd.ShipAddress.Country = FieldCharLimit(CheckPath("country", shipTo), countryCharLimit)

			//qbReceiptAdd.ShipAddress.Note = FieldCharLimit(CheckPath("message", shipTo), noteCharLimit)
			//end shipping address
			/*
				//NO CCINFO WILL BE PASSED AT THIS TIME
				//QB REQUIRES MORE TRANSACTION RESULT DATA THAN CV3 SENDS
				//start Credit Cart transaction input information assignment
				//start require
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.NameOnCard = o.CCName
				//qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardNumber = o.CCNum
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationMonth = o.CCExpM
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationYear = o.CCExpY
				//end require
				//end credit card transaction input information assignment
				//start credit card transaction result information assignment
				//require and not found in cv3 orders
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode = 1
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultMessage = "a"
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CreditCardTransID = "a"
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.MerchantAccountNumber = "a"
				//PaymentStatus may have one of the following values: Unknown, Completed
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentStatus = "Unknown"
				qbReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationTime = "2017-10-31T15:07:38-05:00"
			*/
			/*
				//generate GUID
				uuid, err := uuid.NewRandom()
				if err != nil {
					Log.WithFields(logrus.Fields{"error": err}).Error("error generating guid in getCV3Orders ")
					ErrLog.WithFields(logrus.Fields{"error": err}).Error("error generating guid in getCV3Orders ")
				}
				//ExternalGUID ragex = "0|(\{[0-9a-fA-F]{8}(\-([0-9a-fA-F]{4})){3}\-[0-9a-fA-F]{12}\})"
				qbReceiptAdd.ExternalGUID = `{` + uuid.String() + `}` //"{1904A826-7368-11DC-8317-F7AD55D89593}"
			*/
			var s = make([]string, 0) //will hold product skus
			var sMap = make(map[string]bool, 0)
			var skus = make(map[string]interface{}, 0) //will hold a salesLineAdd or salesGroupLineAdd
			//Prepare the shipToProducts gabs container for range loop
			shipToProductsChildren, err := shipTo.Path("shipToProducts").Children()
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesReceipt")
				ErrLog.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesReceipt")
			}
			var shipToProductFieldMap = ReadFieldMapping("./fieldMaps/cv3ShipToProductReceiptMapping.json")
			//iterate over shipToProducts, save their skus, and start building their line add objects
			for _, prod := range shipToProductsChildren {
				//TODO no duplicates exist?
				if sMap[CheckPath("SKU", prod)] == true {
					//sku already exists in slice
					Log.WithFields(logrus.Fields{"sku": CheckPath("SKU", prod)}).Debug("duplicate sku found in shipTo.ShipToProducts range loop, in GetCV3Orders")
				} else {
					sMap[CheckPath("SKU", prod)] = true //set to true to find duplicates
					s = append(s, CheckPath("SKU", prod))
					var temp = &SalesReceiptLineAdd{}
					//these variables must be set from the shipToProducts
					tempInterface := AddReceiptItem("sku", temp, prod, skus, &WorkCTX{}, shipToProductFieldMap)
					temp = tempInterface.(*SalesReceiptLineAdd)

					//store attribute descriptions, so they can be matched later
					attributes, err := prod.Path("attributes").ChildrenMap()
					if err != nil {
						Log.WithFields(logrus.Fields{"Error": err}).Debug("Error making product attribute children in SalesReceiptAdd")
					} else { //attributes found
						temp.Attributes = make(map[string]string, len(attributes))
						for key, attribute := range attributes {
							temp.Attributes[key] = attribute.Data().(string)
						}
					}
					//temp.Attributes =
					skus[CheckPath("SKU", prod)] = temp
				}
			}
			//Using the slice of skus, get cv3 product information
			var api2 = cv3go.NewApi()
			//api2.Debug = true
			api2.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID)***REMOVED***
			api2.GetProductSKUs(s, false)
			var d2 = api2.Execute(true)
			//Strip the <![CDATA[...]]> tags
			d2 = bytes.Replace(d2, []byte("<![CDATA["), []byte(""), -1)
			d2 = bytes.Replace(d2, []byte("]]>"), []byte(""), -1)
			//Parse return into gabs container
			products, err := gabs.ParseJSON(d2)
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "json": string(d2)}).Error("Error parsing cv3 product json into gabs container in MakeSalesReceipt")
				ErrLog.WithFields(logrus.Fields{"Error": err, "json": string(d2)}).Error("Error parsing cv3 product json into gabs container in MakeSalesReceipt")
			} //ready the gabs container for range loop
			prodChildren, err := products.Path("CV3Data.products").Children()
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "CV3ProductsMapper": products.Path("CV3Data.products")}).Error("Error getting CV3ProductsMapper Children in MakeSalesReceipt")
				ErrLog.WithFields(logrus.Fields{"Error": err, "CV3ProductsMapper": products.Path("CV3Data.products")}).Error("Error getting CV3ProductsMapper Children in MakeSalesReceipt")
			}
			var itemFieldMap = ReadFieldMapping("./fieldMaps/salesReceiptLineAddMapping.json")
			//Build the line items for the order
			for _, item := range prodChildren {
				qbReceiptAdd.BuildLineItems(item, itemFieldMap, skus, workCTX)
			}

			qbReceiptAdd.TxnDate = CheckPath(fieldMap["TxnDate"], o)
			shipMap := ReadFieldMapping("./fieldMaps/shippingReceiptMapping.json")
			//Add shipping
			var p = SalesReceiptLineAdd{}
			//HardCode fields
			p.ItemRef.FullName = shipMap["ItemRef.FullName"] //cfg.HardCodedFields["shipping"]["ItemRef.FullName"] //"*SHIPPING CHARGES-retail" //CheckPath(shipMap["ItemRef.FullName"], shipTo)
			p.Quantity = shipMap["Quantity"]                 //cfg.HardCodedFields["shipping"]["Quantity"]                 //"1"
			//mapped fields
			p.Amount = CheckPath(shipMap["Amount"], shipTo)
			p.Desc = CheckPath(shipMap["Desc"], shipTo)
			p.ClassRef.FullName = CheckPath(shipMap["ClassRef.FullName"], shipTo)
			p.ClassRef.ListID = CheckPath(shipMap["ClassRef.ListID"], shipTo)
			p.InventorySiteRef.FullName = CheckPath(shipMap["InventorySiteRef.FullName"], shipTo)
			p.InventorySiteRef.ListID = CheckPath(shipMap["InventorySiteRef.ListID"], shipTo)
			p.InventorySiteLocationRef.FullName = CheckPath(shipMap["InventorySiteLocationRef.FullName"], shipTo)
			p.InventorySiteLocationRef.ListID = CheckPath(shipMap["InventorySiteLocationRef.ListID"], shipTo)
			p.SalesTaxCodeRef.FullName = CheckPath(shipMap["SalesTaxCodeRef.FullName"], shipTo)
			p.SalesTaxCodeRef.ListID = CheckPath(shipMap["SalesTaxCodeRed.ListID"], shipTo)
			p.PriceLevelRef.FullName = CheckPath(shipMap["PriceLevelRef.FullName"], shipTo)
			p.PriceLevelRef.ListID = CheckPath(shipMap["PriceLevelRef.ListID"], shipTo)
			p.OverrideItemAccountRef.FullName = CheckPath(shipMap["OverrideItemAccountRef.FullName"], shipTo)
			p.OverrideItemAccountRef.ListID = CheckPath(shipMap["OverrideItemAccountRef.ListID"], shipTo)
			p.OptionForPriceRuleConflict = CheckPath(shipMap["OptionForPriceRuleConflict"], shipTo)
			p.SerialNumber = CheckPath(shipMap["SerialNumber"], shipTo)
			p.LotNumber = CheckPath(shipMap["LotNumber"], shipTo)
			p.OptionForPriceRuleConflict = CheckPath(shipMap["OptionForPriceRuleConflict"], shipTo)
			p.Rate = CheckPath(shipMap["Rate"], shipTo)
			p.RatePercent = CheckPath(shipMap["RatePercent"], shipTo)
			p.UnitOfMeasure = CheckPath(shipMap["UnitOfMeasure"], shipTo)

			p.Other1 = CheckPath(shipMap["Other1"], shipTo)
			p.Other2 = CheckPath(shipMap["Other2"], shipTo)

			p.ServiceDate = CheckPath(shipMap["ServiceDate"], shipTo)
			qbReceiptAdd.SalesReceiptLineAdds = append(qbReceiptAdd.SalesReceiptLineAdds, p)

			//Create Tax Item if tax > 0
			taxFloat, err := strconv.ParseFloat(CheckPath("tax", shipTo), 64)
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
				ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
			}
			if taxFloat > 0 {
				var taxMap = ReadFieldMapping("./fieldMaps/taxReceiptMapping.json")
				var taxItem = SalesReceiptLineAdd{}
				taxItem.ItemRef.FullName = taxMap["ItemRef.FullName"] //"Tax"
				taxItem.Desc = taxMap["Desc"]                         //"Tax"
				taxItem.Quantity = taxMap["Quantity"]                 //"1"
				taxItem.Amount = CheckPath(taxMap["Amount"], shipTo)
				qbReceiptAdd.SalesReceiptLineAdds = append(qbReceiptAdd.SalesReceiptLineAdds, taxItem)
			}
			var templateBuff = bytes.Buffer{}
			var escapedQBXML = bytes.Buffer{}
			var tPath = `./templates/qbReceiptAdd.t`

			LoadTemplate(&tPath, qbReceiptAdd, &templateBuff)
			err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
			if err != nil {
				Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
				ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
			}
			//add the QBXML to the work slice
			workCTX.Work = escapedQBXML.String()
			workCTX.Data = qbReceiptAdd
			workCTX.Order = o
			workChan <- *workCTX

			shipToSuccessChan <- ShipToSuccessTracker{
				Index:   shipToIndex,
				OrderID: qbReceiptAdd.CV3OrderID,
			}
			templateBuff.Reset()
			escapedQBXML.Reset()
		}
	}
}

//MakeSalesOrder takes the cv3 order and turns it into a qbxml salesOrderAdd
func MakeSalesOrder(workCount *int, workCTX *WorkCTX, ordersMapper *gabs.Container) {
	//Prepare gabs container for range loop
	oMapper, err := ordersMapper.Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "OrdersMapper": ordersMapper}).Error("Error getting ordersMapper Children in MakeSalesOrder")
		ErrLog.WithFields(logrus.Fields{"Error": err, "OrdersMapper": ordersMapper}).Error("Error getting ordersMapper Children in MakeSalesOrder")
	} //load dynamic field map
	var fieldMap = ReadFieldMapping("./fieldMaps/orderMapping.json")
	//iterate over the orders, then the shiptos, as each shipto needs to be handled as a seperate sales receipt in QB
	for _, o := range oMapper {
		//prepare shipTo gabs container for range loop
		shipToMapper, err := o.Path("shipTos").Children()
		if err != nil {
			Log.WithFields(logrus.Fields{"Error": err, "ShipTosMapper": o.Path("shipTos")}).Error("Error getting shipTosMapper Children in MakeSalesOrder")
			ErrLog.WithFields(logrus.Fields{"Error": err, "ShipTosMapper": o.Path("shipTos")}).Error("Error getting shipTosMapper Children in MakeSalesOrder")
		}
		//start racking the orders and shipTos
		orderSuccessChan <- OrderSuccessTracker{
			OrderID:      CheckPath("orderID", o),
			ShipToLength: len(shipToMapper),
		}
		for shipToIndex, shipTo := range shipToMapper {
			*workCount++
			var qbOrderAdd = SalesOrderAdd{} //object to hold current shipto information
			qbOrderAdd.CV3OrderID = CheckPath("orderID", o)
			qbOrderAdd.RefNumber = CheckPath(fieldMap["RefNumber"], o)
			qbOrderAdd.ShipToIndex = shipToIndex

			//checkPayMethod returns the transactionID of the passed in payMethod
			CheckPayMethod(o, "payMethod")
			CheckPayMethod(o, "additionalPayMethod")
			qbOrderAdd.ClassRef.FullName = CheckPath(fieldMap["ClassRef.FullName"], o)
			qbOrderAdd.ClassRef.ListID = CheckPath(fieldMap["ClassRef.ListID"], o)
			qbOrderAdd.Other = CheckPath(fieldMap["Other"], o)
			qbOrderAdd.ExchangeRate = CheckPath(fieldMap["ExchangeRate"], o)

			if strings.ToLower(CheckPath(fieldMap["IsToPeEmailed"], o)) == "true" {
				qbOrderAdd.IsToBeEmailed = "true"
			}
			if strings.ToLower(CheckPath(fieldMap["IsToPePrinted"], o)) == "true" {
				qbOrderAdd.IsToBePrinted = "true"
			}
			qbOrderAdd.FOB = CheckPath(fieldMap["FOB"], o)

			qbOrderAdd.CustomerMsgRef.FullName = CheckPath(fieldMap["CustomerMsgRef"], shipTo)
			qbOrderAdd.CustomerMsgRef.ListID = CheckPath(fieldMap["CustomerMsgRef.ListID"], shipTo)
			qbOrderAdd.CustomerSalesTaxCodeRef.FullName = CheckPath(fieldMap["CustomerSalesTaxCodeRef.FullName"], o)
			qbOrderAdd.CustomerSalesTaxCodeRef.ListID = CheckPath(fieldMap["CustomerSalesTaxCodeRef.ListID"], o)
			qbOrderAdd.ItemSalesTaxRef.FullName = CheckPath(fieldMap["ItemSalesTaxRef.FullName"], o)
			qbOrderAdd.ItemSalesTaxRef.ListID = CheckPath(fieldMap["ItemSalesTaxRef.ListID"], o)
			//edit for mac's tie downs
			qbOrderAdd.SalesRepRef.FullName = fieldMap["SalesRepRef.FullName"] //CheckPath(fieldMap["SalesRepRef.FullName"], o)
			qbOrderAdd.SalesRepRef.ListID = CheckPath(fieldMap["SalesRepRef.ListID"], o)
			qbOrderAdd.TemplateRef.FullName = CheckPath(fieldMap["TemplateRef.FullName"], o)
			qbOrderAdd.TemplateRef.ListID = CheckPath(fieldMap["TemplateRef.ListID"], o)
			if CheckPath(fieldMap["TermsRef.FullName"], o) == "creditcard" {
				qbOrderAdd.TermsRef.FullName = "Credit Card"
			} else if CheckPath(fieldMap["TermsRef.FullName"], o) == "paypal" {
				qbOrderAdd.TermsRef.FullName = "PayPal"
			} else if CheckPath(fieldMap["TermsRef.FullName"], o) == "ccpaypal" {
				qbOrderAdd.TermsRef.FullName = "CCPaypal"
			} else {
				qbOrderAdd.TermsRef.FullName = CheckPath(fieldMap["TermsRef.FullName"], o)
			}

			qbOrderAdd.TermsRef.ListID = CheckPath(fieldMap["TermsRef.ListID"], o)
			qbOrderAdd.IsManuallyClosed = CheckPath(fieldMap["IsManuallyClosed"], o)

			//start billing information assignment
			//QB will either accept addr 1-5 or addr 1-2 and city state zip country
			var addr = make([]string, 0) // For adding Billing address info
			//If first or last name is not empty, add them as the first line
			if CheckPath("billing.firstName", o) != "" || CheckPath("billing.lastName", o) != "" {
				//if title is not empty add it before the name
				if CheckPath("billing.title", o) != "" {
					addr = append(addr, CheckPath("billing.title", o)+" "+CheckPath("billing.firstName", o)+" "+CheckPath("billing.lastName", o))
				} else {
					addr = append(addr, CheckPath("billing.firstName", o)+" "+CheckPath("billing.lastName", o))
				}
			} //if billing company is not empty ad it as the next available address slot
			if CheckPath("billing.company", o) != "" {
				addr = append(addr, CheckPath("billing.company", o))
			} //add billing address line 1 as the next available address slot
			addr = append(addr, CheckPath("billing.address1", o))
			//if billing address line 2 is not empty add it as the next available address slot
			if CheckPath("billing.address2", o) != "" {
				addr = append(addr, CheckPath("billing.address2", o))
			}
			//add the billing address info to the QB struct
			if len(addr) > 0 {
				qbOrderAdd.BillAddress.Addr1 = FieldCharLimit(addr[0], addrCharLimit)
			}
			if len(addr) > 1 {
				qbOrderAdd.BillAddress.Addr2 = FieldCharLimit(addr[1], addrCharLimit)
			}
			if len(addr) > 2 {
				qbOrderAdd.BillAddress.Addr3 = FieldCharLimit(addr[2], addrCharLimit)
			}
			qbOrderAdd.BillAddress.City = FieldCharLimit(CheckPath("billing.city", o), cityCharLimit)
			qbOrderAdd.BillAddress.Country = FieldCharLimit(CheckPath("billing.country", o), countryCharLimit)
			qbOrderAdd.BillAddress.PostalCode = FieldCharLimit(CheckPath("billing.zip", o), zipCharLimit)
			qbOrderAdd.BillAddress.State = FieldCharLimit(CheckPath("billing.state", o), stateCharLimit)
			//end billing information

			qbOrderAdd.ShipMethodRef.FullName = CheckPath(fieldMap["ShipMethodRef.FullName"], shipTo)
			qbOrderAdd.ShipMethodRef.ListID = CheckPath(fieldMap["ShipMethodRef.ListID"], shipTo)
			qbOrderAdd.Memo = CheckPath(fieldMap["Memo"], o)

			//If the billing name is not paypal, use it as the customers name
			if !strings.Contains(strings.ToLower(CheckPath("billing.firstName", o)), "paypal") {
				//No Comma for Mac's Tie downs
				qbOrderAdd.CustomerRef.FullName = CheckPath("billing.lastName", o) + " " + CheckPath("billing.firstName", o)
			} else { //billing firstName is paypal, so just add paypal as a CustomerRef is required for a SalesOrderAdd
				qbOrderAdd.CustomerRef.FullName = CheckPath("billing.firstName", o)
			}
			qbOrderAdd.ShipDate = CheckPath(fieldMap["ShipDate"], shipTo)

			//Start shipping address//TODO see if QB can handle multiple shipping addresses in the xml
			addr = make([]string, 0) // For adding Billing address info
			//If first or last name is not empty, add them as the first line
			if CheckPath("firstName", shipTo) != "" || CheckPath("lastName", shipTo) != "" {
				//if title is not empty add it before the name
				if CheckPath("title", shipTo) != "" {
					addr = append(addr, CheckPath("title", shipTo)+" "+CheckPath("firstName", shipTo)+" "+CheckPath("lastName", shipTo))
				} else {
					addr = append(addr, CheckPath("firstName", shipTo)+" "+CheckPath("lastName", shipTo))
				}
			} //if shiping company is not empty ad it as the next available address slot
			if CheckPath("company", shipTo) != "" {
				addr = append(addr, CheckPath("company", shipTo))
			} //add Shiping address line 1 as the next available address slot
			addr = append(addr, CheckPath("address1", shipTo))
			//if shiping address line 2 is not empty add it as the next available address slot
			if CheckPath("address2", shipTo) != "" {
				addr = append(addr, CheckPath("address2", shipTo))
			}
			//add the shiping address info to the QB struct
			if len(addr) > 0 {
				qbOrderAdd.ShipAddress.Addr1 = FieldCharLimit(addr[0], addrCharLimit)
			}
			if len(addr) > 1 {
				qbOrderAdd.ShipAddress.Addr2 = FieldCharLimit(addr[1], addrCharLimit)
			}
			if len(addr) > 2 {
				qbOrderAdd.ShipAddress.Addr3 = FieldCharLimit(addr[2], addrCharLimit)
			}
			qbOrderAdd.ShipAddress.City = FieldCharLimit(CheckPath("city", shipTo), cityCharLimit)
			qbOrderAdd.ShipAddress.State = FieldCharLimit(CheckPath("state", shipTo), stateCharLimit)
			qbOrderAdd.ShipAddress.PostalCode = FieldCharLimit(CheckPath("zip", shipTo), zipCharLimit)
			qbOrderAdd.ShipAddress.Country = FieldCharLimit(CheckPath("country", shipTo), countryCharLimit)

			//PONUMBER FOR MAC TIE DOWN
			qbOrderAdd.PONumber = CheckPath(fieldMap["PONumber"], o)

			qbOrderAdd.ShipAddress.Note = FieldCharLimit(CheckPath("message", shipTo), noteCharLimit)
			//end shipping address

			var s = make([]string, 0) //will hold product skus
			var sMap = make(map[string]bool, 0)
			var skus = make(map[string]interface{}, 0) //will hold a salesLineAdd or salesGroupLineAdd
			//Prepare shipToProducts Gabs container for range loop
			shipToProductsChildren, err := shipTo.Path("shipToProducts").Children()
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesOrder")
				ErrLog.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesOrder")
			}
			var shipToProductFieldMap = ReadFieldMapping("./fieldMaps/cv3ShipToProductOrderMapping.json")
			//iterate over shipToProducts, save their skus, and start building their line add objects
			for _, prod := range shipToProductsChildren {
				if sMap[CheckPath("SKU", prod)] == true {
					//sku already exists in slice
					Log.WithFields(logrus.Fields{"sku": CheckPath("SKU", prod)}).Debug("duplicate sku found in shipTo.ShipToProducts range loop, in GetCV3Orders")
				} else {
					sMap[CheckPath("SKU", prod)] = true //set to true to find duplicates
					s = append(s, CheckPath("SKU", prod))
					var temp = &SalesOrderLineAdd{} //SalesReceiptPart{}
					//these variables must be set from the shipToProducts
					tempInterface := AddOrderItem("sku", temp, prod, skus, &WorkCTX{}, shipToProductFieldMap)
					temp = tempInterface.(*SalesOrderLineAdd)
					//store attribute descriptions, so they can be matched later
					attributes, err := prod.Path("attributes").ChildrenMap()
					if err != nil {
						Log.WithFields(logrus.Fields{"Error": err}).Debug("Error making product attribute children in SalesOrderAdd")
					} else { //attributes found
						temp.Attributes = make(map[string]string, len(attributes))
						for key, attribute := range attributes {
							temp.Attributes[key] = attribute.Data().(string)
						}
					}
					skus[CheckPath("SKU", prod)] = temp
				}
			}
			//Using the slice of skus, get cv3 product information
			var api2 = cv3go.NewApi()
			//api2.Debug = true
			api2.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID)***REMOVED***
			api2.GetProductSKUs(s, false)
			var d2 = api2.Execute(true)
			d2 = bytes.Replace(d2, []byte("<![CDATA["), []byte(""), -1)
			d2 = bytes.Replace(d2, []byte("]]>"), []byte(""), -1)

			//Parse return into gabs container
			products, err := gabs.ParseJSON(d2)
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "json": string(d2)}).Error("Error parsing cv3 product json into gabs container in MakeSalesOrder")
				ErrLog.WithFields(logrus.Fields{"Error": err, "json": string(d2)}).Error("Error parsing cv3 product json into gabs container in MakeSalesOrder")
			} //ready the gabs container for range loop
			prodChildren, err := products.Path("CV3Data.products").Children()
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "CV3ProductsMapper": products.Path("CV3Data.products")}).Error("Error getting CV3ProductsMapper Children in MakeSalesOrder")
				ErrLog.WithFields(logrus.Fields{"Error": err, "CV3ProductsMapper": products.Path("CV3Data.products")}).Error("Error getting CV3ProductsMapper Children in MakeSalesOrder")
			}
			var itemFieldMap = ReadFieldMapping("./fieldMaps/salesOrderLineAddMapping.json")
			//Build the line items for this order
			for _, item := range prodChildren {
				qbOrderAdd.BuildLineItems(item, itemFieldMap, skus, workCTX)
			}
			qbOrderAdd.TxnDate = CheckPath(fieldMap["TxnDate"], o)
			shipMap := ReadFieldMapping("./fieldMaps/shippingOrderMapping.json")
			//Add shipping
			var p = SalesOrderLineAdd{}
			//Set hard coded values
			p.ItemRef.FullName = shipMap["ItemRef.FullName"] //"Shipping" //CheckPath(shipMap["ItemRef.FullName"], shipTo)
			p.Quantity = shipMap["Quantity"]                 //"1"
			//Set mapped fields
			p.Amount = CheckPath(shipMap["Amount"], shipTo)

			//maybe add via: for macs tie downs
			p.Desc = CheckPath(shipMap["Desc"], shipTo)
			p.ClassRef.FullName = CheckPath(shipMap["ClassRef.FullName"], shipTo)
			p.ClassRef.ListID = CheckPath(shipMap["ClassRef.ListID"], shipTo)
			p.InventorySiteRef.FullName = CheckPath(shipMap["InventorySiteRef.FullName"], shipTo)
			p.InventorySiteRef.ListID = CheckPath(shipMap["InventorySiteRef.ListID"], shipTo)
			p.InventorySiteLocationRef.FullName = CheckPath(shipMap["InventorySiteLocationRef.FullName"], shipTo)
			p.InventorySiteLocationRef.ListID = CheckPath(shipMap["InventorySiteLocationRef.ListID"], shipTo)
			p.IsManuallyClosed = CheckPath(shipMap["IsManuallyClosed"], shipTo)
			p.LotNumber = CheckPath(shipMap["LotNumber"], shipTo)
			p.OptionForPriceRuleConflict = CheckPath(shipMap["OptionForPriceRuleConflict"], shipTo)

			p.Other1 = CheckPath(shipMap["Other1"], shipTo)
			p.Other2 = CheckPath(shipMap["Other2"], shipTo)
			qbOrderAdd.SalesOrderLineAdds = append(qbOrderAdd.SalesOrderLineAdds, p)

			/*
				//Create Tax Item if tax > 0
				taxFloat, err := strconv.ParseFloat(CheckPath("tax", shipTo), 64)
				if err != nil {
					Log.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
				}
					if taxFloat > 0 {
						var taxItem = SalesOrderLineAdd{}
						taxItem.ItemRef.FullName = "Tax"
						taxItem.Desc = "Tax"
						taxItem.Quantity = "1"
						taxItem.Amount = CheckPath("tax", shipTo)
						qbOrderAdd.SalesOrderLineAdds = append(qbOrderAdd.SalesOrderLineAdds, taxItem)
					}*/

			var templateBuff = bytes.Buffer{}
			var escapedQBXML = bytes.Buffer{}
			var tPath = `./templates/qbSalesOrderAdd.t`

			LoadTemplate(&tPath, qbOrderAdd, &templateBuff)
			err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
			if err != nil {
				Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
				ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
			}

			//add the QBXML to the work slice
			workCTX.Work = escapedQBXML.String()
			workCTX.Data = qbOrderAdd
			workCTX.Order = o
			workChan <- *workCTX

			shipToSuccessChan <- ShipToSuccessTracker{
				Index:   shipToIndex,
				OrderID: qbOrderAdd.CV3OrderID,
			}
			templateBuff.Reset()
			escapedQBXML.Reset()
		}
	}
}

//FieldCharLimit wil truncate characters past the passed in limit
func FieldCharLimit(s string, limit int) string {
	if len(s) > limit {
		return EscapeField(s[:31])
	}
	return EscapeField(s)
}

//CheckPayMethod checks the payMethod and additionalPayMethod fields and maaps the proper fields
func CheckPayMethod(o *gabs.Container, s string) {
	switch CheckPath(s, o) {
	case "paypal_express":
		o.SetP(CheckPath("payPalTransactionID", o), s+"Txn")
	case "paypal":
		o.SetP(CheckPath("payPalTransactionID", o), s+"Txn")
	case "amazon_pay":
		o.SetP(CheckPath("amazonOrderIDs", o), s+"Txn")
	case "purchaseorder":
		o.SetP(CheckPath("purchaseOrder", o), s+"Txn")
	case "giftcertificate_internal":
		o.SetP(CheckPath("giftCertificate", o), s+"Txn")
	case "giftcertificate":
		o.SetP(CheckPath("giftCertificate", o), s+"Txn")
	case "echeck":
		o.SetP(CheckPath("ECAccountName", o), s+"Txn")
	case "visa_checkout":
		o.SetP(CheckPath("visaCheckoutInfo.TransactionID", o), s+"Txn")
	case "creditcard":
		var b = bytes.Buffer{}
		b.WriteString(CheckPath("billing.CCInfo.CCType", o))
		b.WriteString(" ")
		b.WriteString(CheckPath("billing.CCInfo.CCName", o))
		o.SetP(b.String(), s+"Txn")
	case "authorize":
		o.SetP(CheckPath("CCTransactionID", o), s+"Txn")
	case "call":
		//return ""
	case "check":
		//return ""
	case "onfile":
		//return ""
	case "invoice":
		//return ""
	default:
		//return ""
	}
}

//ReadFieldMapping reads the json field map file and returns a map
func ReadFieldMapping(mapPath string) map[string]string {
	var fieldMap map[string]string
	mapFile, err := ioutil.ReadFile(mapPath)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error":   err,
			"mapPath": mapPath,
		}).Error("Error loading field map file")
		ErrLog.WithFields(logrus.Fields{
			"error":   err,
			"mapPath": mapPath,
		}).Error("Error loading field map file")
	}
	err = json.Unmarshal(mapFile, &fieldMap)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error unmarshalling field map JSON")
		ErrLog.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error unmarshalling field map JSON")
	}
	Log.WithFields(logrus.Fields{
		"mapPath": mapPath,
	}).Debug("Loaded field map")
	return fieldMap
}

//CheckPath checks if a gabs path exists and then returns the value as a string
func CheckPath(path string, o *gabs.Container) string {
	if o.ExistsP(path) {
		switch t := o.Path(path).Data().(type) {
		case float64:
			return strconv.FormatFloat(o.Path(path).Data().(float64), 'f', -1, 64)
		case string:
			return o.Path(path).Data().(string)
		case int:
			return strconv.Itoa(o.Path(path).Data().(int))
		default:
			Log.WithFields(logrus.Fields{"type": t, "path": path, "gabs object": o.Path(path)}).Debug("Gabs object is not one of the set types")
			return ""
		}
	}
	Log.WithFields(logrus.Fields{"path": path}).Debug("Path not found in CheckPath")
	return ""
}

//AddOrderItem will add the cv3 product data to the quickbooks return object
func AddOrderItem(sku string, prod interface{}, item *gabs.Container, skus map[string]interface{}, workCTX *WorkCTX, itemFieldMap map[string]string) interface{} {
	//unmarshal into cv3 product, to add to workCTX list of products
	var m = cv3go.Product{}
	err := json.Unmarshal(item.Bytes(), &m)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "product json": item.String()}).Error("Error Unmarshalling cv3 products in MakeSalesOrder")
		ErrLog.WithFields(logrus.Fields{"Error": err, "product json": item.String()}).Error("Error Unmarshalling cv3 products in MakeSalesOrder")
	}
	workCTX.CV3Products = append(workCTX.CV3Products, m)
	if skus[sku] == nil {
		skus[sku] = prod
	}
	if CheckPath(itemFieldMap["ItemRef.FullName"], item) != "" {
		if len(CheckPath(itemFieldMap["ItemRef.FullName"], item)) > 31 {
			skus[sku].(*SalesOrderLineAdd).ItemRef.FullName = html.EscapeString(CheckPath(itemFieldMap["ItemRef.FullName"], item)[:31])
		} else {
			skus[sku].(*SalesOrderLineAdd).ItemRef.FullName = html.EscapeString(CheckPath(itemFieldMap["ItemRef.FullName"], item))
		}
	}
	if CheckPath(itemFieldMap["Quantity"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).Quantity = CheckPath(itemFieldMap["Quantity"], item)
	}
	if CheckPath(itemFieldMap["ItemRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).ItemRef.ListID = CheckPath(itemFieldMap["ItemRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["ClassRef.FullName"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).ClassRef.FullName = CheckPath(itemFieldMap["ClassRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["ClassRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).ClassRef.ListID = CheckPath(itemFieldMap["ClassRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["Desc"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).Desc = html.EscapeString(CheckPath(itemFieldMap["Desc"], item))
	}
	if CheckPath(itemFieldMap["Other1"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).Other1 = CheckPath(itemFieldMap["Other1"], item)
	}
	if CheckPath(itemFieldMap["Other2"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).Other2 = CheckPath(itemFieldMap["Other2"], item)
	}
	if CheckPath(itemFieldMap["UnitOfMeasure"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).UnitOfMeasure = CheckPath(itemFieldMap["UnitOfMeasure"], item)
	}
	if CheckPath(itemFieldMap["Rate"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).Rate = CheckPath(itemFieldMap["Rate"], item)
	}
	if CheckPath(itemFieldMap["RatePercent"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).RatePercent = CheckPath(itemFieldMap["RatePercent"], item)
	}
	if CheckPath(itemFieldMap["PriceLevelRef.FullName"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).PriceLevelRef.FullName = CheckPath(itemFieldMap["PriceLevelRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["PriceLevelRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).PriceLevelRef.ListID = CheckPath(itemFieldMap["PriceLevelRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["OptionForPriceRuleConflict"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).OptionForPriceRuleConflict = CheckPath(itemFieldMap["OptionForPriceRuleConflict"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteRef.FullName"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteRef.FullName = CheckPath(itemFieldMap["InventorySiteRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteRef.ListID = CheckPath(itemFieldMap["InventorySiteRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteLocationRef.FullName"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteLocationRef.FullName = CheckPath(itemFieldMap["InventorySiteLocationRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteLocationRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteLocationRef.ListID = CheckPath(itemFieldMap["InventorySiteLocationRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["SerialNumber"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).SerialNumber = CheckPath(itemFieldMap["SerialNumber"], item)
	}
	if CheckPath(itemFieldMap["LotNumber"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).LotNumber = CheckPath(itemFieldMap["LotNumber"], item)
	}
	if CheckPath(itemFieldMap["SalesTaxCodeRef.FullName"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).SalesTaxCodeRef.FullName = CheckPath(itemFieldMap["SalesTaxCodeRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["SalesTaxCodeRef.ListID"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).SalesTaxCodeRef.ListID = CheckPath(itemFieldMap["SalesTaxCodeRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["IsManuallyClosed"], item) != "" {
		skus[sku].(*SalesOrderLineAdd).IsManuallyClosed = CheckPath(itemFieldMap["IsManuallyClosed"], item)
	}
	return skus[sku].(*SalesOrderLineAdd)
}

//AddReceiptItem will add the cv3 product data to the quickbooks return object
func AddReceiptItem(sku string, prod interface{}, item *gabs.Container, skus map[string]interface{}, workCTX *WorkCTX, itemFieldMap map[string]string) interface{} {
	//unmarshal into cv3 product, to add to workCTX list of products
	var m = cv3go.Product{}
	err := json.Unmarshal(item.Bytes(), &m)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "product json": item.String()}).Error("Error Unmarshalling cv3 products in MakeSalesReceipt")
		ErrLog.WithFields(logrus.Fields{"Error": err, "product json": item.String()}).Error("Error Unmarshalling cv3 products in MakeSalesReceipt")
	}
	workCTX.CV3Products = append(workCTX.CV3Products, m)
	if skus[sku] == nil {
		skus[sku] = prod
	}
	if CheckPath(itemFieldMap["ItemRef.FullName"], item) != "" {
		if len(CheckPath(itemFieldMap["ItemRef.FullName"], item)) > 31 {
			skus[sku].(*SalesReceiptLineAdd).ItemRef.FullName = html.EscapeString(CheckPath(itemFieldMap["ItemRef.FullName"], item)[:31])
		} else {
			skus[sku].(*SalesReceiptLineAdd).ItemRef.FullName = html.EscapeString(CheckPath(itemFieldMap["ItemRef.FullName"], item))
		}
	}
	if CheckPath(itemFieldMap["Quantity"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Quantity = CheckPath(itemFieldMap["Quantity"], item)
	}
	if CheckPath(itemFieldMap["ClassRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).ClassRef.FullName = CheckPath(itemFieldMap["ClassRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["Desc"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Desc = html.EscapeString(CheckPath(itemFieldMap["Desc"], item))
	}
	if CheckPath(itemFieldMap["Other1"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Other1 = CheckPath(itemFieldMap["Other1"], item)
	}
	if CheckPath(itemFieldMap["Other2"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Other2 = CheckPath(itemFieldMap["Other2"], item)
	}
	if CheckPath(itemFieldMap["UnitOfMeasure"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).UnitOfMeasure = CheckPath(itemFieldMap["UnitOfMeasure"], item)
	}
	if CheckPath(itemFieldMap["Rate"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Rate = CheckPath(itemFieldMap["Rate"], item)
	}
	if CheckPath(itemFieldMap["RatePercent"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).RatePercent = CheckPath(itemFieldMap["RatePercent"], item)
	}
	if CheckPath(itemFieldMap["PriceLevelRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).PriceLevelRef.FullName = CheckPath(itemFieldMap["PriceLevelRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["PriceLevelRef.ListID"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).PriceLevelRef.ListID = CheckPath(itemFieldMap["PriceLevelRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["OptionForPriceRuleConflict"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OptionForPriceRuleConflict = CheckPath(itemFieldMap["OptionForPriceRuleConflict"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteLocationRef.FullName = CheckPath(itemFieldMap["InventorySiteRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteRef.ListID"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteRef.ListID = CheckPath(itemFieldMap["InventorySiteRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteLocationRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteLocationRef.FullName = CheckPath(itemFieldMap["InventorySiteLocationRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["InventorySiteLocationRef.ListID"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteLocationRef.ListID = CheckPath(itemFieldMap["InventorySiteLocationRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["SerialNumber"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SerialNumber = CheckPath(itemFieldMap["SerialNumber"], item)
	}
	if CheckPath(itemFieldMap["LotNumber"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).LotNumber = CheckPath(itemFieldMap["LotNumber"], item)
	}
	if CheckPath(itemFieldMap["ServiceDate"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).ServiceDate = CheckPath(itemFieldMap["ServiceDate"], item)
	}
	if CheckPath(itemFieldMap["SalesTaxCodeRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SalesTaxCodeRef.FullName = CheckPath(itemFieldMap["SalesTaxCodeRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["SalesTaxCodeRef.ListID"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SalesTaxCodeRef.ListID = CheckPath(itemFieldMap["SalesTaxCodeRef.ListID"], item)
	}
	if CheckPath(itemFieldMap["OverrideItemAccountRef.FullName"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OverrideItemAccountRef.FullName = CheckPath(itemFieldMap["OverrideItemAccountRef.FullName"], item)
	}
	if CheckPath(itemFieldMap["OverrideItemAccountRef.ListID"], item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OverrideItemAccountRef.ListID = CheckPath(itemFieldMap["OverrideItemAccountRef.ListID"], item)
	}
	//return salesReceiptLine item
	return skus[sku].(*SalesReceiptLineAdd)
}
