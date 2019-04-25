package qbwcGo

import (
	"bytes"
	"encoding/xml"

	"github.com/Sirupsen/logrus"
)

//CustomerQueryQB will query for customers
func CustomerQueryQB() {
	var customerQuery = CustomerQueryRq{}
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var err error
	var tPath = "./templates/QBXMLMsgsRq.t"
	var qbxmlWork = QBXMLWork{}
	_ = customerQuery
	var nameFilter = NameFilter{}
	nameFilter.Name = "Bright Man"
	nameFilter.MatchCriterion = "Contains"
	customerQuery.NameFilter = &nameFilter
	//customerQuery.FullName = append(customerQuery.FullName, "Bright Man")
	b, err := xml.MarshalIndent(&customerQuery, "", "  ")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer query")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer query")
	}

	qbxmlWork.AppendWork(string(b))
	//Reset the template buffer, and build and execute the toplevel QBXML template with the preceeding templates as data
	templateBuff.Reset()
	LoadTemplate(&tPath, qbxmlWork, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
	}

	workChan <- WorkCTX{Work: escapedQBXML.String()}
}
