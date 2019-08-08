package qbwcGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
	"github.com/amazingfly/cv3go"
)

var customerNameCharLimit = 41
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

	if cfg.DataExtActive {
		rand.Seed(time.Now().UTC().UnixNano())
	}

	//Call CV3 for the desired orders
	var api = cv3go.NewApi()
	//api.Debug = true
	api.SetCredentials(cfg.CV3Credentials.User, cfg.CV3Credentials.Pass, cfg.CV3Credentials.ServiceID) //("connector", "#CV3C0nn3ct0r!a", "6c383bc896")
	api.GetOrdersNew()
	//api.GetOrdersRange("50564", "50564") //("25678", "25678") //("7152", "7152") //"7142")

	var d = api.Execute(true)
	Log.Debug(string(d))
	ord, err := gabs.ParseJSON(d)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err, "json": string(d)}).Error("Error parsing json into gabs container in GetCV3Order")
		ErrLog.WithFields(logrus.Fields{"Error": err, "json": string(d)}).Error("Error parsing json into gabs container in GetCV3Order")
	}
	ordTrim := ord.Path("CV3Data.orders")
	/*
		ordTrim, err := gabs.ParseJSONFile("./orderDiscount.json")
		if err != nil {
			fmt.Println(err)
		}
	*/
	switch strings.ToLower(cfg.OrderType) {
	case "salesreceipt":
		MakeSalesReceipt(&workCount, &workCTX, ordTrim)
		break
	case "salesorder":
		MakeSalesOrder(&workCount, &workCTX, ordTrim)
		break
	default:
		Log.WithFields(logrus.Fields{"OrderType": cfg.OrderType}).Error("Error in GetCV3Orders, invalid order type in config")
		ErrLog.WithFields(logrus.Fields{"OrderType": cfg.OrderType}).Error("Error in GetCV3Orders, invalid order type in config")
	}
	if workCount < 1 {
		//workChan <- WorkCTX{Work: "", Type: "NoOp"}
		if CheckPath("CV3Data.error", ord) != "" {
			getLastErrChan <- CheckPath("CV3Data.error", ord)
			Log.WithFields(logrus.Fields{"Error": CheckPath("CV3Data.error", ord), "Json": ord.String()}).Error("Error in CV3 order return")
			ErrLog.WithFields(logrus.Fields{"Error": CheckPath("CV3Data.error", ord), "Json": ord.String()}).Error("Error in CV3 order return")
		} else {
			getLastErrChan <- "No new Orders"
			Log.WithFields(logrus.Fields{"Json": ord.String()}).Info("No new orders in CV3 order return")
		}
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
	var addrFieldMap = ReadFieldMapping("./fieldMaps/addressMapping.json")

	sort.SliceStable(oMapper, func(i, j int) bool { return CheckPath("orderID", oMapper[i]) < CheckPath("orderID", oMapper[j]) })

	//iterate over the orders, then the shiptos, as each shipto needs to be handled as a seperate sales receipt in QB
	for _, o := range oMapper {

		var orderDiscount = new(DiscountCTX)
		orderDiscount.TotalDiscount = ParseFloat(CheckPath("totalOrderDiscount.totalDiscount", o))
		orderDiscount.RemainingDiscount = orderDiscount.TotalDiscount

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
			var customerListID string
			var customerQueryStatusCode int
			//If the billing name is not paypal, so use it as the customers name
			if !strings.Contains(strings.ToLower(fieldMap["CustomerRef.FullName"].Display(o)), "paypal") {
				customerListID, customerQueryStatusCode, err = CustomerQueryQB(o, shipTo, fieldMap["CustomerRef.FullName"].Display(o), strings.ToLower(fieldMap["Email"].Display(o)))
				if err != nil && customerQueryStatusCode == 3 {
					Log.WithFields(logrus.Fields{"Error": err}).Error("Error in customerQuery")
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error in customerQuery")
				}
			}
			var qbReceiptAdd = SalesReceiptAdd{} //object to hold current shipto information
			qbReceiptAdd.DiscountCTX = orderDiscount
			qbReceiptAdd.DiscountCTX.SubTotal = 0.0
			qbReceiptAdd.CV3OrderID = CheckPath("orderID", o)
			*workCount++
			//checkPayMethod will return the transactionID of the passed in PayMethod
			CheckPayMethod(o, "payMethod")
			CheckPayMethod(o, "additionalPayMethod")
			qbReceiptAdd.ClassRef.FullName = fieldMap["ClassRef.FullName"].Display(o)
			qbReceiptAdd.ClassRef.ListID = fieldMap["ClassRef.ListID"].Display(o)
			qbReceiptAdd.Other = fieldMap["Other"].Display(o)
			if fieldMap["ExchangeRate"].Display(o) != "" {
				exchRate, err := strconv.ParseFloat(fieldMap["ExchangeRate"].Display(o), 64)
				if err != nil {
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing string to float for exchange rate in salesReceiptAdd")
				} else {
					qbReceiptAdd.ExchangeRate = exchRate
				}
			}
			if strings.ToLower(fieldMap["IsToPeEmailed"].Display(o)) == "true" {
				qbReceiptAdd.IsToBeEmailed = true
			}
			if strings.ToLower(fieldMap["IsToPePrinted"].Display(o)) == "true" {
				qbReceiptAdd.IsToBePrinted = true
			}
			if strings.ToLower(fieldMap["IsPending"].Display(o)) == "true" {
				qbReceiptAdd.IsPending = "true"
			}
			qbReceiptAdd.FOB = fieldMap["FOB"].Display(o)

			var customerMsgQueryStatusCode int
			qbReceiptAdd.CustomerMsgRef.FullName = ConvertCustomerMsgRef(fieldMap["CustomerMsgRef.FullName"].Display(shipTo))
			if qbReceiptAdd.CustomerMsgRef.FullName != "" {
				var customerMsgListID string
				customerMsgListID, customerMsgQueryStatusCode, err = CustomerMsgQueryQB(o, shipTo, qbReceiptAdd.CustomerMsgRef.FullName)
				if err != nil {
					Log.WithFields(logrus.Fields{"Error": err}).Error("Error in customerMsgQueryQB")
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error in customerMsgQueryQB")
				}
				if customerMsgListID != "" {
					qbReceiptAdd.CustomerMsgRef.FullName = ""
					qbReceiptAdd.CustomerMsgRef.ListID = customerMsgListID
				}
			}

			//qbReceiptAdd.CustomerMsgRef.ListID = fieldMap["CustomerMsgRef.ListID"].Display(shipTo)
			qbReceiptAdd.CustomerSalesTaxCodeRef.FullName = fieldMap["CustomerSalesTaxCodeRef.FullName"].Display(o)
			qbReceiptAdd.CustomerSalesTaxCodeRef.ListID = fieldMap["CustomerSalesTaxCodeRef.ListID"].Display(o)
			qbReceiptAdd.ItemSalesTaxRef.FullName = fieldMap["ItemSalesTaxRef.FullName"].Display(o)
			qbReceiptAdd.ItemSalesTaxRef.ListID = fieldMap["ItemSalesTaxRef.ListID"].Display(o)
			qbReceiptAdd.DepositToAccountRef.FullName = fieldMap["DepositToAccountRef.FullName"].Display(o)
			qbReceiptAdd.DepositToAccountRef.ListID = fieldMap["DepositToAccountRef.ListID"].Display(o)

			//Direct mappingFor Beatrice Bakery "W"
			qbReceiptAdd.SalesRepRef.FullName = fieldMap["SalesRepRef.FullName"].Display() //fieldMap["SalesRepRef.FullName"].Display(o)

			qbReceiptAdd.SalesRepRef.ListID = fieldMap["SalesRepRef.ListID"].Display(o)
			qbReceiptAdd.TemplateRef.FullName = fieldMap["TemplateRef.FullName"].Display(o)
			qbReceiptAdd.TemplateRef.ListID = fieldMap["TemplateRef.ListID"].Display(o)
			qbReceiptAdd.RefNumber = fieldMap["RefNumber"].Display(o)
			qbReceiptAdd.ShipToIndex = shipToIndex

			//start billing information assignment
			//QB will either accept addr 1-5 or addr 1-2 and city state zip country

			qbReceiptAdd.BillAddress.Addr1 = FieldCharLimit(addrFieldMap["BillAddress.Addr1"].Display(o), addrCharLimit)
			qbReceiptAdd.BillAddress.Addr2 = FieldCharLimit(addrFieldMap["BillAddress.Addr2"].Display(o), addrCharLimit)
			qbReceiptAdd.BillAddress.Addr3 = FieldCharLimit(addrFieldMap["BillAddress.Addr3"].Display(o), addrCharLimit)
			qbReceiptAdd.BillAddress.City = FieldCharLimit(addrFieldMap["BillAddress.City"].Display(o), cityCharLimit)
			qbReceiptAdd.BillAddress.Country = FieldCharLimit(addrFieldMap["BillAddress.Country"].Display(o), countryCharLimit)
			qbReceiptAdd.BillAddress.PostalCode = FieldCharLimit(addrFieldMap["BillAddress.PostalCode"].Display(o), zipCharLimit)
			qbReceiptAdd.BillAddress.State = FieldCharLimit(addrFieldMap["BillAddress.State"].Display(o), stateCharLimit)
			//qbReceiptAdd.BillAddress.Note not used
			//end billing information

			qbReceiptAdd.ShipMethodRef.FullName = fieldMap["ShipMethodRef.FullName"].Display(shipTo)
			qbReceiptAdd.ShipMethodRef.ListID = fieldMap["ShipMethodRef.ListID"].Display(shipTo)
			qbReceiptAdd.Memo = fieldMap["Memo"].Display(o)
			qbReceiptAdd.PaymentMethodRef.FullName = fieldMap["PaymentMethodRef.FullName"].Display(o)
			qbReceiptAdd.PaymentMethodRef.ListID = fieldMap["PaymentMethodRef.ListID"].Display(o)

			//If the billing name is not paypal, so use it as the customers name
			if !strings.Contains(strings.ToLower(CheckPath("billing.firstName", o)), "paypal") {
				qbReceiptAdd.CustomerRef.FullName = fieldMap["CustomerRef.FullName"].Display(o) //BuildName(CheckPath("billing.firstName", o), CheckPath("billing.lastName", o))
				qbReceiptAdd.CustomerRef.ListID = customerListID
			} //else bliiling firstname is paypal, so do not add any customer info
			qbReceiptAdd.ShipDate = fieldMap["ShipDate"].Display(shipTo)

			//Start shipping address

			qbReceiptAdd.ShipAddress.Addr1 = FieldCharLimit(addrFieldMap["ShipAddress.Addr1"].Display(shipTo), addrCharLimit)
			qbReceiptAdd.ShipAddress.Addr2 = FieldCharLimit(addrFieldMap["ShipAddress.Addr2"].Display(shipTo), addrCharLimit)
			qbReceiptAdd.ShipAddress.Addr3 = FieldCharLimit(addrFieldMap["ShipAddress.Addr3"].Display(shipTo), addrCharLimit)
			qbReceiptAdd.ShipAddress.City = FieldCharLimit(addrFieldMap["ShipAddress.City"].Display(shipTo), cityCharLimit)
			qbReceiptAdd.ShipAddress.State = FieldCharLimit(addrFieldMap["ShipAddress.State"].Display(shipTo), stateCharLimit)
			qbReceiptAdd.ShipAddress.PostalCode = FieldCharLimit(addrFieldMap["ShipAddress.PostalCode"].Display(shipTo), zipCharLimit)
			qbReceiptAdd.ShipAddress.Country = FieldCharLimit(addrFieldMap["ShipAddress.Country"].Display(shipTo), countryCharLimit)
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
			var s = make([]string, 0) //will hold product skus
			var sMap = make(map[string]bool, 0)
			var skus = make(map[string]interface{}, 0) //will hold a salesLineAdd or salesGroupLineAdd
			//Prepare the shipToProducts gabs container for range loop
			shipToProductsChildren, err := shipTo.Path("shipToProducts").Children()
			if err != nil {
				Log.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesReceipt")
				ErrLog.WithFields(logrus.Fields{"Error": err, "ShipToProductsMapper": shipTo.Path("shipToProducts")}).Error("Error getting shipToProductsMapper Children in MakeSalesReceipt")
			}
			var shipToProductFieldMap = ReadFieldMapping("./fieldMaps/salesReceiptLineAddMapping.json")
			//iterate over shipToProducts, save their skus, and start building their line add objects
			for _, prod := range shipToProductsChildren {
				//check for duplicates exist?
				if sMap[CheckPath("SKU", prod)] == true {
					//sku already exists in slice
					Log.WithFields(logrus.Fields{"sku": CheckPath("SKU", prod)}).Debug("duplicate sku found in shipTo.ShipToProducts range loop, in GetCV3Orders")
				} else {
					sMap[CheckPath("SKU", prod)] = true //set to true to find duplicates
					s = append(s, CheckPath("SKU", prod))
					var temp = &SalesReceiptLineAdd{}
					//these variables must be set from the shipToProducts
					tempInterface := qbReceiptAdd.AddReceiptItem("sku", temp, prod, skus, &WorkCTX{}, shipToProductFieldMap)
					temp = tempInterface.(*SalesReceiptLineAdd)

					skus[CheckPath("SKU", prod)] = temp
					qbReceiptAdd.SalesReceiptLineAdds = append(qbReceiptAdd.SalesReceiptLineAdds, *temp)
				}
			}
			qbReceiptAdd.TxnDate = fieldMap["TxnDate"].Display(o)

			//Calculate discounts
			qbReceiptAdd.AddDiscount(o, shipToIndex)
			qbReceiptAdd.AddShipping(shipTo)

			qbReceiptAdd.AddTax(o, shipTo)

			if cfg.DataExtActive {
				//Add a defMacr aka TxnID
				var txnBuff = bytes.NewBufferString("TxnID:")
				txnBuff.WriteString(CheckPath("orderID", o))
				txnBuff.WriteString(strconv.Itoa(qbReceiptAdd.ShipToIndex))
				qbReceiptAdd.DefMacro = txnBuff.String()
			}

			switch customerQueryStatusCode {
			case 0: // customer record was matched, we have the lisID, so remove the customerRef.FullName to make sure the ListID is used by QB
				qbReceiptAdd.CustomerRef.FullName = ""
				break
			case 1: //no customer records matched the given customer name, create a new customer
				//Leave the customer name alone, to be added the way it is
				Log.Debug("Completely new customer is the answer here")
				break
			case 2: //customer records matched the given customer name but no emails matched the customers email, so create a new customer with the email appended to the end of the name
				qbReceiptAdd.CustomerRef.FullName = FieldCharLimit(qbReceiptAdd.CustomerRef.FullName+" "+strings.ToLower(fieldMap["Email"].Display(o)), customerNameCharLimit)
				Log.WithFields(logrus.Fields{"FullName": qbReceiptAdd.CustomerRef.FullName, "Email": strings.ToLower(fieldMap["Email"].Display(o))}).Debug("Resetting customerRef.FullName")
				break
			}

			var templateBuff = bytes.Buffer{}
			var escapedQBXML = bytes.Buffer{}
			var tPath = `./templates/qbReceiptAdd.t`
			var qbxmlWork = QBXMLWork{} //make([]string, 0)

			LoadTemplate(&tPath, qbReceiptAdd, &templateBuff)
			//Add the SalesReceiptQBXML to the slice for use in QBXMLMsgsRq, the top level template
			qbxmlWork.AppendWork(templateBuff.String())

			if cfg.DataExtActive {
				//Prepare the DataExtAdds //TODO make universal rith now its only for macsTieDowns
				var dExts = DataExtAddQB(qbReceiptAdd.DefMacro)
				qbxmlWork.AppendWork(dExts)
			}

			//Reset the template buffer, and build and execute the toplevel QBXML template with the preceeding templates as data
			templateBuff.Reset()
			tPath = "./templates/QBXMLMsgsRq.t"
			LoadTemplate(&tPath, qbxmlWork, &templateBuff)
			err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
			if err != nil {
				Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
				ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
			}
			//add the QBXML to the work slice
			workCTX.Work = escapedQBXML.String()
			workCTX.Data = qbReceiptAdd
			workCTX.Order = o
			workCTX.Type = "SalesReceiptAdd"
			switch customerQueryStatusCode {
			case 0:
				//do nothing as a customer record was matched
				break
			case 1: //no customer records matched the given customer name, create a new customer
				//Do not resend the order in CustomerAddQB
				workCTX.NoResendOrder = true
				Log.Debug("Completly new customer Add")
				CustomerAddQB(*workCTX)
				break
			case 2: //customer records matched the given customer name but no emails matched the customers email, so create a new customer with the email appended to the end of the name
				//Do not resend the order in CustomerAddQB
				workCTX.NoResendOrder = true
				Log.Debug("Adding customer with appended email")
				CustomerAddQB(*workCTX)
				break
			}

			switch customerMsgQueryStatusCode {
			case 0:
				//Do nothing, the customerMsg was found
			case 1: //Add customerMsg
				CustomerMsgAddQB(*workCTX)
			}
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
	sort.SliceStable(oMapper, func(i, j int) bool { return CheckPath("orderID", oMapper[i]) < CheckPath("orderID", oMapper[j]) })

	var fieldMap = ReadFieldMapping("./fieldMaps/orderMapping.json")
	var addrFieldMap = ReadFieldMapping("./fieldMaps/addressMapping.json")
	//iterate over the orders, then the shiptos, as each shipto needs to be handled as a seperate sales receipt in QB
	for _, o := range oMapper {
		var orderDiscount = new(DiscountCTX)
		orderDiscount.TotalDiscount = ParseFloat(CheckPath("totalOrderDiscount.totalDiscount", o))
		orderDiscount.RemainingDiscount = orderDiscount.TotalDiscount
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
			qbOrderAdd.DiscountCTX = orderDiscount
			qbOrderAdd.DiscountCTX.SubTotal = 0.0
			qbOrderAdd.CV3OrderID = CheckPath("orderID", o)
			qbOrderAdd.RefNumber = fieldMap["RefNumber"].Display(o)
			qbOrderAdd.ShipToIndex = shipToIndex

			//checkPayMethod returns the transactionID of the passed in payMethod
			CheckPayMethod(o, "payMethod")
			CheckPayMethod(o, "additionalPayMethod")
			qbOrderAdd.ClassRef.FullName = fieldMap["ClassRef.FullName"].Display(o)
			qbOrderAdd.ClassRef.ListID = fieldMap["ClassRef.ListID"].Display(o)
			qbOrderAdd.Other = fieldMap["Other"].Display(o)
			qbOrderAdd.ExchangeRate = fieldMap["ExchangeRate"].Display(o)

			if strings.ToLower(fieldMap["IsToPeEmailed"].Display(o)) == "true" {
				qbOrderAdd.IsToBeEmailed = "true"
			}
			if strings.ToLower(fieldMap["IsToPePrinted"].Display(o)) == "true" {
				qbOrderAdd.IsToBePrinted = "true"
			}
			qbOrderAdd.FOB = fieldMap["FOB"].Display(o)

			qbOrderAdd.CustomerMsgRef.FullName = fieldMap["CustomerMsgRef.FullName"].Display(shipTo)
			var customerMsgQueryStatusCode int
			if qbOrderAdd.CustomerMsgRef.FullName != "" {
				var customerMsgListID string
				customerMsgListID, customerMsgQueryStatusCode, err = CustomerMsgQueryQB(o, shipTo, qbOrderAdd.CustomerMsgRef.FullName)
				if err != nil {
					Log.WithFields(logrus.Fields{"Error": err}).Error("Error in customerMsgQueryQB")
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error in customerMsgQueryQB")
				}
				if customerMsgListID != "" {
					qbOrderAdd.CustomerMsgRef.FullName = ""
					qbOrderAdd.CustomerMsgRef.ListID = customerMsgListID
				}
			}

			qbOrderAdd.CustomerSalesTaxCodeRef.FullName = ConvertCustomerMsgRef(fieldMap["CustomerMsgRef.FullName"].Display(shipTo))

			//qbOrderAdd.CustomerMsgRef.ListID = fieldMap["CustomerMsgRef.ListID"].Display(shipTo)

			qbOrderAdd.CustomerSalesTaxCodeRef.ListID = fieldMap["CustomerSalesTaxCodeRef.ListID"].Display(o)
			qbOrderAdd.ItemSalesTaxRef.FullName = fieldMap["ItemSalesTaxRef.FullName"].Display(o)
			qbOrderAdd.ItemSalesTaxRef.ListID = fieldMap["ItemSalesTaxRef.ListID"].Display(o)
			qbOrderAdd.SalesRepRef.FullName = fieldMap["SalesRepRef.FullName"].Display()
			qbOrderAdd.SalesRepRef.ListID = fieldMap["SalesRepRef.ListID"].Display(o)
			qbOrderAdd.TemplateRef.FullName = fieldMap["TemplateRef.FullName"].Display(o)
			qbOrderAdd.TemplateRef.ListID = fieldMap["TemplateRef.ListID"].Display(o)
			//
			//TODO make dynamicly laded from file
			//
			/*was for Macs Tie Downs
			if fieldMap["TermsRef.FullName"].Display(o) == "creditcard" {
				qbOrderAdd.TermsRef.FullName = "Credit Card"
			} else if fieldMap["TermsRef.FullName"].Display(o) == "paypal" {
				qbOrderAdd.TermsRef.FullName = "PayPal"
			} else if fieldMap["TermsRef.FullName"].Display(o) == "ccpaypal" {
				qbOrderAdd.TermsRef.FullName = "CCPaypal"
			} else {
				qbOrderAdd.TermsRef.FullName = fieldMap["TermsRef.FullName"].Display(o)
			}
			*/
			qbOrderAdd.TermsRef.FullName = fieldMap["TermsRef.FullName"].Display(o)
			qbOrderAdd.TermsRef.ListID = fieldMap["TermsRef.ListID"].Display(o)
			qbOrderAdd.IsManuallyClosed = fieldMap["IsManuallyClosed"].Display(o)

			//start billing information assignment
			qbOrderAdd.BillAddress.Addr1 = FieldCharLimit(addrFieldMap["BillAddress.Addr1"].Display(o), addrCharLimit)
			qbOrderAdd.BillAddress.Addr2 = FieldCharLimit(addrFieldMap["BillAddress.Addr2"].Display(o), addrCharLimit)
			qbOrderAdd.BillAddress.Addr3 = FieldCharLimit(addrFieldMap["BillAddress.Addr3"].Display(o), addrCharLimit)
			qbOrderAdd.BillAddress.City = FieldCharLimit(addrFieldMap["BillAddress.City"].Display(o), cityCharLimit)
			qbOrderAdd.BillAddress.Country = FieldCharLimit(addrFieldMap["BillAddress.Country"].Display(o), countryCharLimit)
			qbOrderAdd.BillAddress.PostalCode = FieldCharLimit(addrFieldMap["BillAddress.PostalCode"].Display(o), zipCharLimit)
			qbOrderAdd.BillAddress.State = FieldCharLimit(addrFieldMap["BillAddress.State"].Display(o), stateCharLimit)
			//end billing information

			qbOrderAdd.ShipMethodRef.FullName = fieldMap["ShipMethodRef.FullName"].Display(shipTo)
			qbOrderAdd.ShipMethodRef.ListID = fieldMap["ShipMethodRef.ListID"].Display(shipTo)
			qbOrderAdd.Memo = fieldMap["Memo"].Display(o)

			qbOrderAdd.ShipDate = fieldMap["ShipDate"].Display(shipTo)

			//Start shipping address
			qbOrderAdd.ShipAddress.Addr1 = FieldCharLimit(addrFieldMap["ShipAddress.Addr1"].Display(shipTo), addrCharLimit)
			qbOrderAdd.ShipAddress.Addr2 = FieldCharLimit(addrFieldMap["ShipAddress.Addr2"].Display(shipTo), addrCharLimit)
			qbOrderAdd.ShipAddress.Addr3 = FieldCharLimit(addrFieldMap["ShipAddress.Addr3"].Display(shipTo), addrCharLimit)
			qbOrderAdd.ShipAddress.City = FieldCharLimit(addrFieldMap["ShipAddress.City"].Display(shipTo), cityCharLimit)
			qbOrderAdd.ShipAddress.State = FieldCharLimit(addrFieldMap["ShipAddress.State"].Display(shipTo), stateCharLimit)
			qbOrderAdd.ShipAddress.PostalCode = FieldCharLimit(addrFieldMap["ShipAddress.PostalCode"].Display(shipTo), zipCharLimit)
			qbOrderAdd.ShipAddress.Country = FieldCharLimit(addrFieldMap["ShipAddress.Country"].Display(shipTo), countryCharLimit)
			qbOrderAdd.PONumber = fieldMap["PONumber"].Display(o)
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
			var shipToProductFieldMap = ReadFieldMapping("./fieldMaps/orderLineAddMapping.json")
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
					tempInterface := qbOrderAdd.AddOrderItem("sku", temp, prod, o, skus, &WorkCTX{}, shipToProductFieldMap)
					temp = tempInterface.(*SalesOrderLineAdd)

					//temp.SalesTaxCodeRef.FullName = "Tax"
					skus[CheckPath("SKU", prod)] = temp

					qbOrderAdd.SalesOrderLineAdds = append(qbOrderAdd.SalesOrderLineAdds, *temp)
				}
			}
			qbOrderAdd.TxnDate = fieldMap["TxnDate"].Display(o)

			qbOrderAdd.AddDiscount(o, shipToIndex)
			qbOrderAdd.AddShipping(shipTo)

			//Add the tax if applicable
			qbOrderAdd.AddTax(o, shipTo)

			if cfg.DataExtActive {
				//Add a defMacr aka TxnID
				var txnBuff = bytes.NewBufferString("TxnID:")
				txnBuff.WriteString(CheckPath("orderID", o))
				txnBuff.WriteString(strconv.Itoa(qbOrderAdd.ShipToIndex))
				txnBuff.WriteString(strconv.Itoa(rand.Int()))
				qbOrderAdd.DefMacro = txnBuff.String()
			}

			var customerListID string
			var customerQueryStatusCode int
			//If the billing name is not paypal, so use it as the customers name
			if !strings.Contains(strings.ToLower(fieldMap["CustomerRef.FullName"].Display(o)), "paypal") {
				customerListID, customerQueryStatusCode, err = CustomerQueryQB(o, shipTo, fieldMap["CustomerRef.FullName"].Display(o), strings.ToLower(fieldMap["Email"].Display(o)))
				if err != nil && customerQueryStatusCode == 3 {
					Log.WithFields(logrus.Fields{"Error": err}).Error("Error in customerQuery")
					ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error in customerQuery")
				}
				ErrLog.WithFields(logrus.Fields{"ListID": customerListID, "customer query status code": customerQueryStatusCode}).Error("Customer query returned in getCV3Orders")

				qbOrderAdd.CustomerRef.FullName = fieldMap["CustomerRef.FullName"].Display(o)
				qbOrderAdd.CustomerRef.ListID = customerListID
			}

			switch customerQueryStatusCode {
			case 0: // customer record was matched, we have the lisID, so remove the customerRef.FullName to make sure the ListID is used by QB
				qbOrderAdd.CustomerRef.FullName = ""
				break
			case 1: //no customer records matched the given customer name, create a new customer
				//Leave the customer name alone, to be added the way it is
				ErrLog.Error("Completely new customer is the answer here")
				break
			case 2: //customer records matched the given customer name but no emails matched the customers email, so create a new customer with the email appended to the end of the name
				qbOrderAdd.CustomerRef.FullName = FieldCharLimit(qbOrderAdd.CustomerRef.FullName+" "+strings.ToLower(fieldMap["Email"].Display(o)), customerNameCharLimit)
				Log.WithFields(logrus.Fields{"FullName": qbOrderAdd.CustomerRef.FullName, "Email": strings.ToLower(fieldMap["Email"].Display(o))}).Debug("Resetting customerRef.FullName")
				break
			}

			//Build the templates, add them to the top level template and then xml.Escape them
			var templateBuff = bytes.Buffer{}
			var escapedQBXML = bytes.Buffer{}
			var qbxmlWork = QBXMLWork{} //make([]string, 0)
			var tPath = `./templates/qbSalesOrderAdd.t`
			LoadTemplate(&tPath, qbOrderAdd, &templateBuff)
			//Add the SalesOrderQBXML to th slice for use in QBXMLMsgsRq, the top level template
			qbxmlWork.AppendWork(templateBuff.String())

			if cfg.DataExtActive {
				//Prepare the DataExtAdds //TODO make universal rith now its only for macsTieDowns
				var dExts = DataExtAddQB(qbOrderAdd.DefMacro)
				qbxmlWork.AppendWork(dExts)
			}

			//Reset the template buffer, and build and execute the toplevel QBXML template with the preceeding templates as data
			templateBuff.Reset()
			tPath = "./templates/QBXMLMsgsRq.t"
			LoadTemplate(&tPath, qbxmlWork, &templateBuff)
			err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
			if err != nil {
				Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
				ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
			}

			workCTX.Work = escapedQBXML.String()
			workCTX.Data = qbOrderAdd
			workCTX.Order = o
			workCTX.Type = "SalesOrderAdd"

			switch customerQueryStatusCode {
			case 0:
				//do nothing as a customer record was matched
				break
			case 1: //no customer records matched the given customer name, create a new customer
				//Do not resend the order in CustomerAddQB
				workCTX.NoResendOrder = true
				Log.Debug("Completly new customer Add")
				CustomerAddQB(*workCTX)
				break
			case 2: //customer records matched the given customer name but no emails matched the customers email, so create a new customer with the email appended to the end of the name
				//Do not resend the order in CustomerAddQB
				workCTX.NoResendOrder = true
				Log.Debug("Adding customer with appended email")
				CustomerAddQB(*workCTX)
				break
			}

			switch customerMsgQueryStatusCode {
			case 0:
				//Do nothing, the customerMsg was found
			case 1: //Add customerMsg
				CustomerMsgAddQB(*workCTX)
			}
			workChan <- *workCTX

			shipToSuccessChan <- ShipToSuccessTracker{
				Index:   shipToIndex,
				OrderID: qbOrderAdd.CV3OrderID,
			}
			Log.WithFields(logrus.Fields{"non escaped xml request": templateBuff.String()}).Debug("logging non escaped xml")
			Log.WithFields(logrus.Fields{"xml request": escapedQBXML.String()}).Debug("Logging full xml request")
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

//ReadDictFile reads the json dictionary file and returns a map[string]string
func ReadDictFile(mapPath string) map[string]string {
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

//ReadFieldMapping will load a map[string]MappingObject from a json file
func ReadFieldMapping(path string) map[string]MappingObject {
	var mapObj = make(map[string]MappingObject, 0)
	mapFile, err := ioutil.ReadFile(path)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error":   err,
			"mapPath": path,
		}).Error("Error loading field map file")
		ErrLog.WithFields(logrus.Fields{
			"error":   err,
			"mapPath": path,
		}).Error("Error loading field map file")
	}
	err = json.Unmarshal(mapFile, &mapObj)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error unmarshalling field map JSON")
		ErrLog.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error unmarshalling field map JSON")
	}
	Log.WithFields(logrus.Fields{
		"mapPath": path,
	}).Debug("Loaded field map")
	return mapObj
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
	return ""
}

//AddOrderItem will add the cv3 product data to the quickbooks return object
func (orderAdd *SalesOrderAdd) AddOrderItem(sku string, prod interface{}, item, o *gabs.Container, skus map[string]interface{}, workCTX *WorkCTX, itemFieldMap map[string]MappingObject) interface{} {
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
	if itemFieldMap["ItemRef.FullName"].Display(item) != "" {
		if len(itemFieldMap["ItemRef.FullName"].Display(item)) > 31 {
			skus[sku].(*SalesOrderLineAdd).ItemRef.FullName = html.EscapeString(itemFieldMap["ItemRef.FullName"].Display(item)[:31])
		} else {
			skus[sku].(*SalesOrderLineAdd).ItemRef.FullName = html.EscapeString(itemFieldMap["ItemRef.FullName"].Display(item))
		}
	}
	if itemFieldMap["Quantity"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).Quantity = itemFieldMap["Quantity"].Display(item)
	}
	if itemFieldMap["ItemRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).ItemRef.ListID = itemFieldMap["ItemRef.ListID"].Display(item)
	}
	if itemFieldMap["ClassRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).ClassRef.FullName = itemFieldMap["ClassRef.FullName"].Display(item)
	}
	if itemFieldMap["ClassRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).ClassRef.ListID = itemFieldMap["ClassRef.ListID"].Display(item)
	}
	if itemFieldMap["Desc"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).Desc = html.EscapeString(itemFieldMap["Desc"].Display(item))
	}
	if itemFieldMap["Other1"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).Other1 = itemFieldMap["Other1"].Display(item)
	}
	if itemFieldMap["Other2"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).Other2 = itemFieldMap["Other2"].Display(item)
	}
	if itemFieldMap["UnitOfMeasure"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).UnitOfMeasure = itemFieldMap["UnitOfMeasure"].Display(item)
	}
	if itemFieldMap["Rate"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).Rate = itemFieldMap["Rate"].Display(item)
		//Calculate subtotal for use in discount calculations
		orderAdd.DiscountCTX.SubTotal += ParseFloat(skus[sku].(*SalesOrderLineAdd).Rate) * ParseFloat(skus[sku].(*SalesOrderLineAdd).Quantity)
	}
	if itemFieldMap["RatePercent"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).RatePercent = itemFieldMap["RatePercent"].Display(item)
	}
	if itemFieldMap["PriceLevelRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).PriceLevelRef.FullName = itemFieldMap["PriceLevelRef.FullName"].Display(item)
	}
	if itemFieldMap["PriceLevelRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).PriceLevelRef.ListID = itemFieldMap["PriceLevelRef.ListID"].Display(item)
	}
	if itemFieldMap["OptionForPriceRuleConflict"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).OptionForPriceRuleConflict = itemFieldMap["OptionForPriceRuleConflict"].Display(item)
	}
	if itemFieldMap["InventorySiteRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteRef.FullName = itemFieldMap["InventorySiteRef.FullName"].Display(item)
	}
	if itemFieldMap["InventorySiteRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteRef.ListID = itemFieldMap["InventorySiteRef.ListID"].Display(item)
	}
	if itemFieldMap["InventorySiteLocationRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteLocationRef.FullName = itemFieldMap["InventorySiteLocationRef.FullName"].Display(item)
	}
	if itemFieldMap["InventorySiteLocationRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).InventorySiteLocationRef.ListID = itemFieldMap["InventorySiteLocationRef.ListID"].Display(item)
	}
	if itemFieldMap["SerialNumber"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).SerialNumber = itemFieldMap["SerialNumber"].Display(item)
	}
	if itemFieldMap["LotNumber"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).LotNumber = itemFieldMap["LotNumber"].Display(item)
	}
	if itemFieldMap["SalesTaxCodeRef.FullName"].Display(item, o) != "" {
		Log.WithFields(logrus.Fields{"SalesTaxCodeRef.FullName": itemFieldMap["SalesTaxCodeRef.FullName"].Display(item, o)}).Debug("logging the salesTaxCodeRef.FullName from mapping")
		skus[sku].(*SalesOrderLineAdd).SalesTaxCodeRef.FullName = itemFieldMap["SalesTaxCodeRef.FullName"].Display(item, o)
	}
	if itemFieldMap["SalesTaxCodeRef.ListID"].Display(item, o) != "" {
		skus[sku].(*SalesOrderLineAdd).SalesTaxCodeRef.ListID = itemFieldMap["SalesTaxCodeRef.ListID"].Display(item, o)
	}
	if itemFieldMap["IsManuallyClosed"].Display(item) != "" {
		skus[sku].(*SalesOrderLineAdd).IsManuallyClosed = itemFieldMap["IsManuallyClosed"].Display(item)
	}
	return skus[sku].(*SalesOrderLineAdd)
}

//AddReceiptItem will add the cv3 product data to the quickbooks return object
func (receiptAdd *SalesReceiptAdd) AddReceiptItem(sku string, prod interface{}, item *gabs.Container, skus map[string]interface{}, workCTX *WorkCTX, itemFieldMap map[string]MappingObject) interface{} {
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
	if itemFieldMap["ItemRef.FullName"].Display(item) != "" {
		if len(itemFieldMap["ItemRef.FullName"].Display(item)) > 31 {
			skus[sku].(*SalesReceiptLineAdd).ItemRef.FullName = html.EscapeString(itemFieldMap["ItemRef.FullName"].Display(item)[:31])
		} else {
			skus[sku].(*SalesReceiptLineAdd).ItemRef.FullName = html.EscapeString(itemFieldMap["ItemRef.FullName"].Display(item))
		}
	}
	if itemFieldMap["Quantity"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Quantity = itemFieldMap["Quantity"].Display(item)
	}
	if itemFieldMap["ClassRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).ClassRef.FullName = itemFieldMap["ClassRef.FullName"].Display(item)
	}
	if itemFieldMap["Desc"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Desc = html.EscapeString(itemFieldMap["Desc"].Display(item))
	}
	if itemFieldMap["Other1"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Other1 = itemFieldMap["Other1"].Display(item)
	}
	if itemFieldMap["Other2"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Other2 = itemFieldMap["Other2"].Display(item)
	}
	if itemFieldMap["UnitOfMeasure"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).UnitOfMeasure = itemFieldMap["UnitOfMeasure"].Display(item)
	}
	if itemFieldMap["Rate"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).Rate = itemFieldMap["Rate"].Display(item)
		receiptAdd.DiscountCTX.SubTotal += ParseFloat(skus[sku].(*SalesReceiptLineAdd).Rate) * ParseFloat(skus[sku].(*SalesReceiptLineAdd).Quantity)
	}
	if itemFieldMap["RatePercent"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).RatePercent = itemFieldMap["RatePercent"].Display(item)
	}
	if itemFieldMap["PriceLevelRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).PriceLevelRef.FullName = itemFieldMap["PriceLevelRef.FullName"].Display(item)
	}
	if itemFieldMap["PriceLevelRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).PriceLevelRef.ListID = itemFieldMap["PriceLevelRef.ListID"].Display(item)
	}
	if itemFieldMap["OptionForPriceRuleConflict"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OptionForPriceRuleConflict = itemFieldMap["OptionForPriceRuleConflict"].Display(item)
	}
	if itemFieldMap["InventorySiteRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteRef.FullName = itemFieldMap["InventorySiteRef.FullName"].Display(item)
	}
	if itemFieldMap["InventorySiteRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteRef.ListID = itemFieldMap["InventorySiteRef.ListID"].Display(item)
	}
	if itemFieldMap["InventorySiteLocationRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteLocationRef.FullName = itemFieldMap["InventorySiteLocationRef.FullName"].Display(item)
	}
	if itemFieldMap["InventorySiteLocationRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).InventorySiteLocationRef.ListID = itemFieldMap["InventorySiteLocationRef.ListID"].Display(item)
	}
	if itemFieldMap["SerialNumber"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SerialNumber = itemFieldMap["SerialNumber"].Display(item)
	}
	if itemFieldMap["LotNumber"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).LotNumber = itemFieldMap["LotNumber"].Display(item)
	}
	if itemFieldMap["ServiceDate"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).ServiceDate = itemFieldMap["ServiceDate"].Display(item)
	}
	if itemFieldMap["SalesTaxCodeRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SalesTaxCodeRef.FullName = itemFieldMap["SalesTaxCodeRef.FullName"].Display(item)
	}
	if itemFieldMap["SalesTaxCodeRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).SalesTaxCodeRef.ListID = itemFieldMap["SalesTaxCodeRef.ListID"].Display(item)
	}
	if itemFieldMap["OverrideItemAccountRef.FullName"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OverrideItemAccountRef.FullName = itemFieldMap["OverrideItemAccountRef.FullName"].Display(item)
	}
	if itemFieldMap["OverrideItemAccountRef.ListID"].Display(item) != "" {
		skus[sku].(*SalesReceiptLineAdd).OverrideItemAccountRef.ListID = itemFieldMap["OverrideItemAccountRef.ListID"].Display(item)
	}
	//Add item to the sales receipt
	//receiptAdd.SalesReceiptLineAdds = append(receiptAdd.SalesReceiptLineAdds, *skus[sku].(*SalesReceiptLineAdd))
	//return salesReceiptLine item
	return skus[sku].(*SalesReceiptLineAdd)
}

//AddDiscount Adds a discount line item when a totalOrderDiscount exists
func (orderAdd *SalesOrderAdd) AddDiscount(o *gabs.Container, shipToIndex int) {
	var orderDiscount = SalesOrderLineAdd{}
	discountMap := ReadFieldMapping("./fieldMaps/discountOrderMapping.json")
	if CheckPath("totalOrderDiscount.totalDiscount", o) != "" {
		if orderAdd.DiscountCTX.RemainingDiscount != 0.0 {
			if orderAdd.DiscountCTX.RemainingDiscount <= orderAdd.DiscountCTX.SubTotal {
				orderDiscount.ItemRef.FullName = discountMap["ItemRef.FullName"].Display()
				orderDiscount.Quantity = discountMap["Quantity"].Display()
				orderDiscount.SalesTaxCodeRef.FullName = discountMap["SalesTaxCodeRef"].Display()
				orderDiscount.Rate = strconv.FormatFloat(0.0-orderAdd.DiscountCTX.RemainingDiscount, 'f', 2, 64)
				orderAdd.DiscountCTX.RemainingDiscount = 0.0
				orderAdd.SalesOrderLineAdds = append(orderAdd.SalesOrderLineAdds, orderDiscount)
			} else {
				//subtract this cv3 shipTo's subtotal from the remaining discount for use in next shipTo
				orderAdd.DiscountCTX.RemainingDiscount = orderAdd.DiscountCTX.RemainingDiscount - orderAdd.DiscountCTX.SubTotal
				orderDiscount.ItemRef.FullName = discountMap["ItemRef.FullName"].Display()
				orderDiscount.Quantity = discountMap["Quantity"].Display()
				orderDiscount.SalesTaxCodeRef.FullName = discountMap["SalesTaxCodeRef"].Display()
				orderDiscount.Rate = strconv.FormatFloat(0.0-orderAdd.DiscountCTX.SubTotal, 'f', 2, 64)
				orderAdd.SalesOrderLineAdds = append(orderAdd.SalesOrderLineAdds, orderDiscount)
			}
		}
	}
}

//AddDiscount Adds a discount line item when a totalOrderDiscount exists
func (receiptAdd *SalesReceiptAdd) AddDiscount(o *gabs.Container, shipToIndex int) {
	var receiptDiscount = SalesReceiptLineAdd{}
	discountMap := ReadFieldMapping("./fieldMaps/discountReceiptMapping.json")
	//Check if there is a discount, then split it between cv3 shipTos if needed
	if CheckPath("totalOrderDiscount.totalDiscount", o) != "" {
		if receiptAdd.DiscountCTX.RemainingDiscount != 0.0 {
			if receiptAdd.DiscountCTX.RemainingDiscount <= receiptAdd.DiscountCTX.SubTotal {
				receiptDiscount.ClassRef.FullName = discountMap["ClassRef.FullName"].Display()
				receiptDiscount.Desc = discountMap["Desc"].Display(o)
				receiptDiscount.ItemRef.FullName = discountMap["ItemRef.FullName"].Display()
				receiptDiscount.Quantity = discountMap["Quantity"].Display()
				receiptDiscount.SalesTaxCodeRef.FullName = discountMap["SalesTaxCodeRef.FullName"].Display()
				receiptDiscount.Rate = strconv.FormatFloat(0.0-receiptAdd.DiscountCTX.RemainingDiscount, 'f', 2, 64)
				receiptAdd.DiscountCTX.RemainingDiscount = 0.0 //discount is used up, set it to zero
				receiptAdd.SalesReceiptLineAdds = append(receiptAdd.SalesReceiptLineAdds, receiptDiscount)
			} else { //remainingDiscount is greater than the subtotal of this cv3 shipTo / qbwc salesReceipt
				//subtract this cv3 shipTo's subtotal from the remaining discount for use in next shipTo
				receiptAdd.DiscountCTX.RemainingDiscount = receiptAdd.DiscountCTX.RemainingDiscount - receiptAdd.DiscountCTX.SubTotal
				receiptDiscount.ClassRef.FullName = discountMap["ClassRef.FullName"].Display()
				receiptDiscount.Desc = discountMap["Desc"].Display(o)
				receiptDiscount.ItemRef.FullName = discountMap["ItemRef.FullName"].Display()
				receiptDiscount.Quantity = discountMap["Quantity"].Display()
				receiptDiscount.SalesTaxCodeRef.FullName = discountMap["SalesTaxCodeRef.FullName"].Display()
				receiptDiscount.Rate = strconv.FormatFloat(0.0-receiptAdd.DiscountCTX.SubTotal, 'f', 2, 64)
				receiptAdd.SalesReceiptLineAdds = append(receiptAdd.SalesReceiptLineAdds, receiptDiscount)
			}
		}
	}
}

//AddShipping will add a shipping line item
func (receiptAdd *SalesReceiptAdd) AddShipping(shipTo *gabs.Container) {
	shipMap := ReadFieldMapping("./fieldMaps/shippingReceiptMapping.json")
	//Add shipping
	var p = SalesReceiptLineAdd{}
	//HardCode fields
	p.ItemRef.FullName = shipMap["ItemRef.FullName"].Display() //cfg.HardCodedFields["shipping"]["ItemRef.FullName"] //"*SHIPPING CHARGES-retail" //shipMap["ItemRef.FullName"].Display(shipTo)
	p.Quantity = shipMap["Quantity"].Display()                 //cfg.HardCodedFields["shipping"]["Quantity"]                 //"1"
	//mapped fields
	p.Amount = shipMap["Amount"].Display(shipTo)
	p.Desc = shipMap["Desc"].Display(shipTo)
	p.ClassRef.FullName = shipMap["ClassRef.FullName"].Display(shipTo)
	p.ClassRef.ListID = shipMap["ClassRef.ListID"].Display(shipTo)
	p.InventorySiteRef.FullName = shipMap["InventorySiteRef.FullName"].Display(shipTo)
	p.InventorySiteRef.ListID = shipMap["InventorySiteRef.ListID"].Display(shipTo)
	p.InventorySiteLocationRef.FullName = shipMap["InventorySiteLocationRef.FullName"].Display(shipTo)
	p.InventorySiteLocationRef.ListID = shipMap["InventorySiteLocationRef.ListID"].Display(shipTo)
	p.SalesTaxCodeRef.FullName = shipMap["SalesTaxCodeRef.FullName"].Display(shipTo)
	p.SalesTaxCodeRef.ListID = shipMap["SalesTaxCodeRed.ListID"].Display(shipTo)
	p.PriceLevelRef.FullName = shipMap["PriceLevelRef.FullName"].Display(shipTo)
	p.PriceLevelRef.ListID = shipMap["PriceLevelRef.ListID"].Display(shipTo)
	p.OverrideItemAccountRef.FullName = shipMap["OverrideItemAccountRef.FullName"].Display(shipTo)
	p.OverrideItemAccountRef.ListID = shipMap["OverrideItemAccountRef.ListID"].Display(shipTo)
	p.OptionForPriceRuleConflict = shipMap["OptionForPriceRuleConflict"].Display(shipTo)
	p.SerialNumber = shipMap["SerialNumber"].Display(shipTo)
	p.LotNumber = shipMap["LotNumber"].Display(shipTo)
	p.OptionForPriceRuleConflict = shipMap["OptionForPriceRuleConflict"].Display(shipTo)
	p.Rate = shipMap["Rate"].Display(shipTo)
	p.RatePercent = shipMap["RatePercent"].Display(shipTo)
	p.UnitOfMeasure = shipMap["UnitOfMeasure"].Display(shipTo)

	p.Other1 = shipMap["Other1"].Display(shipTo)
	p.Other2 = shipMap["Other2"].Display(shipTo)

	p.ServiceDate = shipMap["ServiceDate"].Display(shipTo)
	receiptAdd.SalesReceiptLineAdds = append(receiptAdd.SalesReceiptLineAdds, p)
}

//AddShipping will add a shipping line item
func (orderAdd *SalesOrderAdd) AddShipping(shipTo *gabs.Container) {
	shipMap := ReadFieldMapping("./fieldMaps/shippingOrderMapping.json")
	//Add shipping
	var p = SalesOrderLineAdd{}
	//Set hard coded values
	p.ItemRef.FullName = shipMap["ItemRef.FullName"].Display() //"Shipping" //shipMap["ItemRef.FullName"].Display(shipTo)
	p.Quantity = shipMap["Quantity"].Display()                 //"1"
	//Set mapped fields
	p.Amount = shipMap["Amount"].Display(shipTo)

	//maybe add via: for macs tie downs
	p.Desc = shipMap["Desc"].Display(shipTo)
	p.ClassRef.FullName = shipMap["ClassRef.FullName"].Display(shipTo)
	p.ClassRef.ListID = shipMap["ClassRef.ListID"].Display(shipTo)
	p.InventorySiteRef.FullName = shipMap["InventorySiteRef.FullName"].Display(shipTo)
	p.InventorySiteRef.ListID = shipMap["InventorySiteRef.ListID"].Display(shipTo)
	p.InventorySiteLocationRef.FullName = shipMap["InventorySiteLocationRef.FullName"].Display(shipTo)
	p.InventorySiteLocationRef.ListID = shipMap["InventorySiteLocationRef.ListID"].Display(shipTo)
	p.IsManuallyClosed = shipMap["IsManuallyClosed"].Display(shipTo)
	p.LotNumber = shipMap["LotNumber"].Display(shipTo)
	p.OptionForPriceRuleConflict = shipMap["OptionForPriceRuleConflict"].Display(shipTo)

	p.Other1 = shipMap["Other1"].Display(shipTo)
	p.Other2 = shipMap["Other2"].Display(shipTo)
	orderAdd.SalesOrderLineAdds = append(orderAdd.SalesOrderLineAdds, p)
}

//AddTax will add a tax item
func (orderAdd *SalesOrderAdd) AddTax(o, shipTo *gabs.Container) {
	//Create Tax Item if tax > 0
	taxFloat, err := strconv.ParseFloat(CheckPath("tax", shipTo), 64)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
	}
	ErrLog.WithFields(logrus.Fields{"orderAdd.ItemSalesTaxRef.FullName": orderAdd.ItemSalesTaxRef.FullName}).Error("Logging for troubleshoting")
	if taxFloat > 0 && orderAdd.ItemSalesTaxRef.FullName == "" {
		stateTaxDict := ReadDictFile("./fieldMaps/stateTaxMapping.json")
		orderAdd.ItemSalesTaxRef.FullName = stateTaxDict[CheckPath("billing.state", o)]
	}
}

//AddTax will add a tax item
func (receiptAdd *SalesReceiptAdd) AddTax(o, shipTo *gabs.Container) {
	//Create Tax Item if tax > 0
	taxFloat, err := strconv.ParseFloat(CheckPath("tax", shipTo), 64)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
	}
	if taxFloat > 0 {
		stateTaxMap := ReadDictFile("./fieldMaps/stateTaxMapping.json")
		receiptAdd.ItemSalesTaxRef.FullName = stateTaxMap[CheckPath("billing.state", o)]
	}
	/*
		//The way it is for Beatrice Bakery at the moment
		//Create Tax Item if tax > 0
		taxFloat, err := strconv.ParseFloat(CheckPath("tax", shipTo), 64)
		if err != nil {
			Log.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
			ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error parsing tax in SalesReceiptAdd")
		}
		if taxFloat > 0 {
			var taxMap = ReadFieldMapping("./fieldMaps/taxReceiptMapping.json")
			var taxItem = SalesReceiptLineAdd{}
			taxItem.ItemRef.FullName = taxMap["ItemRef.FullName"].Display() //"Tax"
			taxItem.Desc = taxMap["Desc"].Display()                         //"Tax"
			taxItem.Quantity = taxMap["Quantity"].Display()                 //"1"
			taxItem.Amount = taxMap["Amount"].Display(shipTo)
			receiptAdd.SalesReceiptLineAdds = append(receiptAdd.SalesReceiptLineAdds, taxItem)
		}*/
}

//ParseFloat is a wrapper for strconv.ParseFloat, and will parse a string into a float
func ParseFloat(s string) float64 {
	var f float64
	var err error
	if s != "" {
		f, err = strconv.ParseFloat(s, 64)
		if err != nil {
			Log.WithFields(logrus.Fields{"Error": err, "value": s}).Error("Error parsing float in add order")
			ErrLog.WithFields(logrus.Fields{"Error": err, "value": s}).Error("Error parsing float in add order")
		}
	}
	return f
}
