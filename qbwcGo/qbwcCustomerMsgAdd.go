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
		customerMsg.Name = workCTX.Data.(SalesReceiptAdd).CustomerMsgRef.FullName
		//customerMsg.IsActive =
	case SalesOrderAdd:
		customerMsg.Name = workCTX.Data.(SalesOrderAdd).CustomerMsgRef.FullName
		//customerMsg.IsActive =
	}

	LoadTemplate(&tPath, &customerMsg, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template in CustomerMsgAddQB")
	}
	//Send prepaired QBXML to the workInsertChan
	workInsertChan <- WorkCTX{Work: escapedQBXML.String(), Data: workCTX.Data, Order: workCTX.Order, Type: "customerAddRq"}
	// TODO it seems that the workChan picks up before the workInsertChan,  Perhaps a better way can be found
	//time.Sleep(10)
	workChan <- workCTX
}
