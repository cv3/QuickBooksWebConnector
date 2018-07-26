package qbwcGo

import (
	"bytes"
	"encoding/xml"

	"github.com/Sirupsen/logrus"
)

//CustomerAddQB will add a customer to the quickbooks database
func CustomerAddQB(workCTX WorkCTX) {
	//insertWG.Add(1)
	var err error
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var tPath = `./templates/qbCustomerAdd.t`
	var customer = CustomerAddRq{}
	var fieldMap = ReadFieldMapping("./fieldMaps/customerAddMapping.json")
	var addrFieldMap = ReadFieldMapping("./fieldMaps/addressMapping.json")
	//var shipToIndex string
	//var isReceipt bool
	switch workCTX.Data.(type) {
	case SalesReceiptAdd:
		var tempReceiptAdd = workCTX.Data.(SalesReceiptAdd)
		customer.BillAddress = tempReceiptAdd.BillAddress
		customer.ShipAddress = tempReceiptAdd.ShipAddress
		//if a macro was defined add a trailing _c so it will not be a duplicate
		if tempReceiptAdd.DefMacro != "" {
			tempReceiptAdd.DefMacro = tempReceiptAdd.DefMacro + "_C"
		}
		workCTX.Data = tempReceiptAdd

		var templateBuff = bytes.Buffer{}
		var escapedQBXML = bytes.Buffer{}
		var tPath = `./templates/qbReceiptAdd.t`
		var qbxmlWork = QBXMLWork{} //make([]string, 0)

		LoadTemplate(&tPath, tempReceiptAdd, &templateBuff)
		//Add the SalesReceiptQBXML to th slice for use in QBXMLMsgsRq, the top level template
		qbxmlWork.AppendWork(templateBuff.String())

		//Prepare the DataExtAdds //TODO make universal rith now its only for macsTieDowns
		if cfg.DataExtActive {
			var dExts = DataExtAddQB(tempReceiptAdd.DefMacro)
			qbxmlWork.AppendWork(dExts)
		}
		//Reset the template buffer, and build and execute the toplevel QBXML template with the preceeding templates as data
		templateBuff.Reset()
		tPath = "./templates/QBXMLMsgsRq.t"
		LoadTemplate(&tPath, qbxmlWork, &templateBuff)
		err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerAddQB")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerAddQB")
		}
		//add the QBXML to the work slice
		workCTX.Work = escapedQBXML.String()
		break
	case SalesOrderAdd:
		var tempOrderAdd = workCTX.Data.(SalesOrderAdd)
		customer.BillAddress = tempOrderAdd.BillAddress
		customer.ShipAddress = tempOrderAdd.ShipAddress
		//if a macro was defined add a trailing _c so it will not be a duplicate
		if tempOrderAdd.DefMacro != "" {
			tempOrderAdd.DefMacro = tempOrderAdd.DefMacro + "_C"
		}
		workCTX.Data = tempOrderAdd

		var templateBuff = bytes.Buffer{}
		var escapedQBXML = bytes.Buffer{}
		var tPath = `./templates/qbSalesOrderAdd.t`
		var qbxmlWork = QBXMLWork{} //make([]string, 0)

		LoadTemplate(&tPath, tempOrderAdd, &templateBuff)
		//Add the SalesOrderQBXML to th slice for use in QBXMLMsgsRq, the top level template
		qbxmlWork.AppendWork(templateBuff.String())

		//Prepare the DataExtAdds //TODO make universal rith now its only for macsTieDowns
		if cfg.DataExtActive {
			var dExts = DataExtAddQB(tempOrderAdd.DefMacro)
			qbxmlWork.AppendWork(dExts)
		}
		//Reset the template buffer, and build and execute the toplevel QBXML template with the preceeding templates as data
		templateBuff.Reset()
		tPath = "./templates/QBXMLMsgsRq.t"
		LoadTemplate(&tPath, qbxmlWork, &templateBuff)
		err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerAddQB")
			ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerAddQB")
		}
		//add the QBXML to the work slice
		workCTX.Work = escapedQBXML.String()
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

		/*TODO remove after Mac's Tie Down install is complete
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
		}
		*/

		//Start shipping address
		var addrLine1 = bytes.Buffer{}
		var addr = make([]string, 0) // For adding Billing address info
		//TODO make the same as billing?
		switch {
		case CheckPath("title", shipTo) != "":
			addrLine1.WriteString(CheckPath("title", shipTo))
			addrLine1.WriteString(" ")
			fallthrough
		case CheckPath("firstName", shipTo) != "":
			addrLine1.WriteString(CheckPath("firstName", shipTo))
			addrLine1.WriteString(" ")
			fallthrough
		case CheckPath("firstName", shipTo) != "":
			addrLine1.WriteString(CheckPath("lastName", shipTo))
			fallthrough
		case CheckPath("company", shipTo) != "":
			addrLine1.WriteString(" ")
			CheckPath("company", shipTo)
		}
		addr = append(addr, FieldCharLimit(addrLine1.String(), addrCharLimit))
		/*
			//add Shiping address line 1 as the next available address slot
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
		*/

		qbShipTo.Addr1 = FieldCharLimit(addrFieldMap["ShipAddress.Addr1"].Display(shipTo), addrCharLimit)
		qbShipTo.Addr2 = FieldCharLimit(addrFieldMap["ShipAddress.Addr2"].Display(shipTo), addrCharLimit)
		qbShipTo.Addr3 = FieldCharLimit(addrFieldMap["ShipAddress.Addr3"].Display(shipTo), addrCharLimit)
		qbShipTo.City = FieldCharLimit(addrFieldMap["ShipAddress.City"].Display(shipTo), cityCharLimit)
		qbShipTo.State = FieldCharLimit(addrFieldMap["ShipAddress.State"].Display(shipTo), stateCharLimit)
		qbShipTo.PostalCode = FieldCharLimit(addrFieldMap["ShipAddress.PostalCode"].Display(shipTo), zipCharLimit)
		qbShipTo.Country = FieldCharLimit(addrFieldMap["ShipAddress.Country"].Display(shipTo), countryCharLimit)

		qbShipTo.Name = qbShipTo.Addr1 //CheckPath("name", shipTo)
		//qbShipTo.Note =
		//qbShipTo.DefaultShipTo =
		qbShipTos[i] = qbShipTo
	}
	//Direct mapping for beatrice bakery
	customer.CustomerTypeRef.FullName = fieldMap["CustomerTypeRef.FullName"].Display()
	//Direct mapping for beatrice bakery
	customer.SalesRepRef.FullName = fieldMap["SalesRepRef.FullName"].Display()

	customer.Name = fieldMap["Name"].Display(workCTX.Order) //BuildName(CheckPath("billing.firstName", workCTX.Order), CheckPath("billing.lastName", workCTX.Order)) //CheckPath("billing.lastName", workCTX.Order) + ", " + CheckPath("billing.firstName", workCTX.Order)
	//customer.AccountNumber = fieldMap["AccountNumber"].Display(workCTX.Order)
	customer.Email = fieldMap["Email"].Display(workCTX.Order)
	customer.Phone = CheckPath("billing.phone", workCTX.Order)
	customer.FirstName = fieldMap["FirstName"].Display(workCTX.Order) //CheckPath("billing.firstName", workCTX.Order)
	customer.LastName = fieldMap["LastName"].Display(workCTX.Order)   //CheckPath("billing.lastName", workCTX.Order)

	customer.Cc = fieldMap["Cc"].Display(workCTX.Order)
	customer.ClassRef.FullName = fieldMap["ClassRef.FullName"].Display(workCTX.Order)
	customer.ClassRef.ListID = fieldMap["ClassRef.ListID"].Display(workCTX.Order)
	customer.CompanyName = fieldMap["CompanyName"].Display(workCTX.Order)

	//Direct mappings credit card for Beatrice,
	customer.PreferredPaymentMethodRef.FullName = fieldMap["PreferredPaymentMethodRef.FullName"].Display() //fieldMap[""].Display(workCTX.Order)
	//Direct mappings NET for beatrice bakery
	customer.TermsRef.FullName = fieldMap["TermsRef.FullName"].Display()
	//Direct mappings RETAIL for beatrice bakery
	customer.PriceLevelRef.FullName = fieldMap["PriceLevelRef.FullName"].Display()

	LoadTemplate(&tPath, &customer, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in ImportCV3ItemsToQB")
	}
	//Send prepaired QBXML to the workInsertChan
	workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: workCTX.Data, Order: workCTX.Order, Type: "customerAddRq", Attempted: workCTX.Attempted}
	//insertWG.Done()
	//workChan <- workCTX
	workInsertChan <- workCTX
}
