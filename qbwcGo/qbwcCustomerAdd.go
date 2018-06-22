package qbwcGo

import (
	"bytes"
	"encoding/xml"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
)

//CustomerAddQB will add a customer to the quickbooks database
func CustomerAddQB(workCTX WorkCTX) {
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var tPath = `./templates/qbCustomerAdd.t`
	var customer = CustomerAddRq{}
	var fieldMap = ReadFieldMapping("./fieldMaps/customerAddMapping.json")
	//var shipToIndex string
	//var isReceipt bool
	switch workCTX.Data.(type) {
	case SalesReceiptAdd:
		//isReceipt = true
		//shipToIndex = strconv.Itoa(workCTX.Data.(SalesReceiptAdd).ShipToIndex)
		customer.BillAddress = workCTX.Data.(SalesReceiptAdd).BillAddress
		customer.ShipAddress = workCTX.Data.(SalesReceiptAdd).ShipAddress

	case SalesOrderAdd:
		//isReceipt = false
		//shipToIndex = strconv.Itoa(workCTX.Data.(SalesOrderAdd).ShipToIndex)
		customer.BillAddress = workCTX.Data.(SalesOrderAdd).BillAddress
		customer.ShipAddress = workCTX.Data.(SalesOrderAdd).ShipAddress
	}

	//prepare the shipTo gabs container for a range loop
	shipToMapper, err := workCTX.Order.Path("shipTos").Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error getting shipTosMapper Children in MakeSalesReceipt")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error getting shipTosMapper Children in MakeSalesReceipt")
	}
	var qbShipTos = make([]QBShipToAddress, len(shipToMapper))
	for i, shipTo := range shipToMapper {
		var qbShipTo = QBShipToAddress{}
		//qbShipTo.

		//Start shipping address
		var addr = make([]string, 0) // For adding shipTo address info
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
			qbShipTo.Addr1 = FieldCharLimit(addr[0], addrCharLimit)
		}
		if len(addr) > 1 {
			qbShipTo.Addr2 = FieldCharLimit(addr[1], addrCharLimit)
		}
		if len(addr) > 2 {
			qbShipTo.Addr3 = FieldCharLimit(addr[2], addrCharLimit)
		}
		qbShipTo.City = FieldCharLimit(CheckPath("city", shipTo), cityCharLimit)
		qbShipTo.State = FieldCharLimit(CheckPath("state", shipTo), stateCharLimit)
		qbShipTo.PostalCode = FieldCharLimit(CheckPath("zip", shipTo), zipCharLimit)
		qbShipTo.Country = FieldCharLimit(CheckPath("country", shipTo), countryCharLimit)
		qbShipTo.Name = CheckPath("name", shipTo)
		//qbShipTo.Note =
		//qbShipTo.DefaultShipTo =
		qbShipTos[i] = qbShipTo
	}

	customer.Name = CheckPath("billing.lastName", workCTX.Order) + ", " + CheckPath("billing.firstName", workCTX.Order)
	//customer.AccountNumber = CheckPath(fieldMap["AccountNumber"], workCTX.Order)
	customer.Email = CheckPath(fieldMap["Email"], workCTX.Order)
	customer.Phone = CheckPath("billing.phone", workCTX.Order)
	customer.FirstName = CheckPath("billing.firstName", workCTX.Order)
	customer.LastName = CheckPath("billing.lastName", workCTX.Order)

	customer.Cc = CheckPath(fieldMap["Cc"], workCTX.Order)
	customer.ClassRef.FullName = CheckPath(fieldMap["ClassRef.FullName"], workCTX.Order)
	customer.ClassRef.ListID = CheckPath(fieldMap["ClassRef.ListID"], workCTX.Order)
	customer.CompanyName = CheckPath(fieldMap["CompanyName"], workCTX.Order)
	//customer. = CheckPath(fieldMap[""], workCTX.Order)

	LoadTemplate(&tPath, &customer, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")
	}
	//Send prepaired QBXML to the workInsertChan
	workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: workCTX.Data, Order: workCTX.Order, Type: "customerAddRq"}
	workChan <- workCTX
}

//BuildCustomerFromCV3Order takes a cv3 order *gabs.Container and builds a customerAdd object
func BuildCustomerFromCV3Order(shipTo *gabs.Container) {

}
