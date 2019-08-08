package qbwcGo

import (
	"bytes"
	"encoding/xml"
	"errors"
	"strings"

	"github.com/TeamFairmont/gabs"

	"github.com/Sirupsen/logrus"
)

/*customerStatusCodes
0: customerQuery returned records with matching email.  Return the listID to be used in the order
1: No customer records returned, the customer name does not exist in quickbooks so add the customer.
2: Records that matched the customer name were returned, but nomatching emails were found.  So Add a new customer with the email appened to the customer nam
3: other error
*/

//CustomerQueryQB will query for customers, attempt to match the email from the order
//if the email matches, the listID will be returned to be used in the addOrder, and the customers info will be updated, excluding the the email
//if the name exists but the email does not match, a new customer record will be created with the new email appended to the customer name
//if there is no matching name a new customer will be created.
func CustomerQueryQB(order, shipTo *gabs.Container, name, email string) (string, int, error) {
	var listID = ""
	var customerQueryStatusCode int
	var customerQuery = CustomerQueryRq{}
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var err error
	var tPath = "./templates/QBXMLMsgsRq.t"
	var qbxmlWork = QBXMLWork{}
	var nameFilter = NameFilter{}
	nameFilter.Name = name
	nameFilter.MatchCriterion = "Contains" //StartsWith, Contains, EndsWith
	customerQuery.NameFilter = &nameFilter
	b, err := xml.MarshalIndent(&customerQuery, "", "  ")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer query")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer query")
		return "", 3, err
	}

	qbxmlWork.AppendWork(string(b))
	LoadTemplate(&tPath, qbxmlWork, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in GetCV3Orders")
		return "", 3, err
	}
	workChan <- WorkCTX{Work: escapedQBXML.String()}
	var customerQueryResponse = <-customerQueryResponseChan
	if customerQueryResponse.StatusCode == "0" {
		for _, cust := range customerQueryResponse.CustomerRet {
			if strings.ToLower(cust.Email) == strings.ToLower(email) {
				listID = cust.ListID
				customerQueryStatusCode = 0
				CustomerModQB(order, shipTo, cust)
				return listID, customerQueryStatusCode, nil
			}
		}
	} else {
		Log.WithFields(logrus.Fields{"statusCode": customerQueryResponse.StatusCode, "statusMessage": customerQueryResponse.StatusMessage}).Error("Error with customer query")
		ErrLog.WithFields(logrus.Fields{"statusCode": customerQueryResponse.StatusCode, "statusMessage": customerQueryResponse.StatusMessage}).Error("Error with customer query")
	}
	switch {
	case len(customerQueryResponse.CustomerRet) < 1:
		customerQueryStatusCode = 1
		return "", customerQueryStatusCode, errors.New("No customer records match the customer name")
	case listID == "":
		customerQueryStatusCode = 2
		return "", customerQueryStatusCode, errors.New("No customer records that match the customers email")
	}
	customerQueryStatusCode = 3
	return "", customerQueryStatusCode, errors.New(customerQueryResponse.StatusMessage)
}
