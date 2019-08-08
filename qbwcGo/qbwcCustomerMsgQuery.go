package qbwcGo

import (
	"bytes"
	"encoding/xml"
	"strconv"

	"github.com/TeamFairmont/gabs"

	"github.com/Sirupsen/logrus"
)

//CustomerMsgQueryQB will query for customer messages, then return the list id, or add a new messages
func CustomerMsgQueryQB(order, shipTo *gabs.Container, name string) (string, int, error) {
	var listID = ""
	var statusCode int
	var customerMsgQuery = CustomerMsgQueryRq{}
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var err error
	var tPath = "./templates/QBXMLMsgsRq.t"
	var qbxmlWork = QBXMLWork{}
	var nameFilter = NameFilter{}
	nameFilter.Name = name
	nameFilter.MatchCriterion = "Contains" //StartsWith, Contains, EndsWith
	customerMsgQuery.NameFilter = &nameFilter
	b, err := xml.MarshalIndent(&customerMsgQuery, "", "  ")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customerMsgQuery")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customerMsgQuery")
		return "", 3, err
	}

	qbxmlWork.AppendWork(string(b))
	LoadTemplate(&tPath, qbxmlWork, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerMsgQuery")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerMsgQuery")
		return "", 3, err
	}
	workChan <- WorkCTX{Work: escapedQBXML.String()}
	var customerMsgQueryResponse = <-customerMsgQueryResponseChan
	statusCode, err = strconv.Atoi(customerMsgQueryResponse.StatusCode)
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error converting status code to int in CustomerMsgQueryQB")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error converting status code to int in CustomerMsgQueryQB")
	}
	switch statusCode {
	case 0:
		for _, msg := range customerMsgQueryResponse.CustomerMsgRet {
			if msg.Name == name {
				listID = msg.ListID
			}
		}
	default:
		b, err := xml.MarshalIndent(&customerMsgQueryResponse, "", "  ")
		if err != nil {
			ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Unknown erro in customerMsgQuery")
			Log.WithFields(logrus.Fields{"Error": err}).Error("Unknown erro in customerMsgQuery")
		}
		Log.WithFields(logrus.Fields{"statusCode": customerMsgQueryResponse.StatusCode, "CustomerMsgQueryRs": string(b)}).Error("Error with customer query")
		ErrLog.WithFields(logrus.Fields{"statusCode": customerMsgQueryResponse.StatusCode, "CustomerMsgQueryRs": string(b)}).Error("Error with customer query")
	}
	return listID, statusCode, err
}
