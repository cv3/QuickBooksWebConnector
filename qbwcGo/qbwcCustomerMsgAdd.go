package qbwcGo

import (
	"bytes"
	"encoding/xml"

	"github.com/Sirupsen/logrus"
)

//CustomerMsgAddQB will add a customerMsgRef to the quickbooks database
func CustomerMsgAddQB(workCTX WorkCTX) {
	var err error
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var tPath = `./templates/qbCustomerMsgAdd.t`
	var customerMsg = CustomerMsgAddRq{}
	//var fieldMap = ReadFieldMapping("./fieldMaps/customerMsgAddMapping.json")

	switch workCTX.Data.(type) {
	case SalesReceiptAdd:
		var tempReceiptAdd = workCTX.Data.(SalesReceiptAdd)
		customerMsg.Name = tempReceiptAdd.CustomerMsgRef.FullName

		tempReceiptAdd.DefMacro = tempReceiptAdd.DefMacro + "_M"
		workCTX.Data = tempReceiptAdd

		var templateBuff = bytes.Buffer{}
		var escapedQBXML = bytes.Buffer{}
		var tPath = `./templates/qbReceiptAdd.t`
		var qbxmlWork = make([]string, 0)

		LoadTemplate(&tPath, tempReceiptAdd, &templateBuff)
		//Add the SalesReceiptQBXML to th slice for use in QBXMLMsgsRq, the top level template
		qbxmlWork = append(qbxmlWork, templateBuff.String())

		//Prepare the DataExtAdds
		if cfg.DataExtActive {
			var dExts = DataExtAddQB(tempReceiptAdd.DefMacro)
			qbxmlWork = append(qbxmlWork, dExts)
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

		//customerMsg.IsActive =
	case SalesOrderAdd:
		var tempOrderAdd = workCTX.Data.(SalesOrderAdd)
		customerMsg.Name = tempOrderAdd.CustomerMsgRef.FullName

		tempOrderAdd.DefMacro = tempOrderAdd.DefMacro + "_M"
		workCTX.Data = tempOrderAdd

		var templateBuff = bytes.Buffer{}
		var escapedQBXML = bytes.Buffer{}
		var tPath = `./templates/qbOrderAdd.t`
		var qbxmlWork = make([]string, 0)

		LoadTemplate(&tPath, tempOrderAdd, &templateBuff)
		//Add the SalesOrderQBXML to the slice for use in QBXMLMsgsRq, the top level template
		qbxmlWork = append(qbxmlWork, templateBuff.String())

		//Prepare the DataExtAdds
		if cfg.DataExtActive {
			var dExts = DataExtAddQB(tempOrderAdd.DefMacro)
			qbxmlWork = append(qbxmlWork, dExts)
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

	LoadTemplate(&tPath, &customerMsg, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in CustomerMsgAddQB")
	}
	//Send prepaired QBXML to the workInsertChan
	workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: workCTX.Data, Order: workCTX.Order, Type: "customerAddRq"}
	//insertWG.Done()
	workInsertChan <- workCTX
}
