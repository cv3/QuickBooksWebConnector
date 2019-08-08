package qbwcGo

import (
	"bytes"
	"encoding/xml"

	"github.com/TeamFairmont/gabs"

	"github.com/Sirupsen/logrus"
)

//CustomerModQB will modify a customer record
func CustomerModQB(order, shipTo *gabs.Container, customerRecord CustomerRet) {
	var updateRecord = false
	var fieldMap = ReadFieldMapping("./fieldMaps/customerModMapping.json")
	var customerMod = CustomerModRq{}
	var escapedQBXML = bytes.Buffer{}
	var templateBuff = bytes.Buffer{}
	var err error
	var tPath = "./templates/QBXMLMsgsRq.t"
	var qbxmlWork = QBXMLWork{}

	customerMod.CustomerMod.ListID = customerRecord.ListID
	customerMod.CustomerMod.EditSequence = customerRecord.EditSequence
	customerMod.CustomerMod.Name = customerRecord.Name

	customerMod.CustomerMod.IsActive = fieldMap["IsActive"].Display(order)
	var accountRef = AccountRef{}
	if fieldMap["ClassRef.ListID"].Display(order)+fieldMap["ClassRef.FullName"].Display(order) != "" {
		accountRef.ListID = fieldMap["ClassRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["ClassRef.FullName"].Display(order)
		if customerRecord.ClassRef.FullName != accountRef.FullName {
			customerMod.CustomerMod.ClassRef = &accountRef
			updateRecord = true
		}
	}
	if fieldMap["ParentRef.ListID"].Display(order)+fieldMap["ParentRef.FullName"].Display(order) != "" {
		accountRef.ListID = fieldMap["ParentRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["ParentRef.FullName"].Display(order)
		if customerRecord.ParentRef.FullName != accountRef.FullName {
			customerMod.CustomerMod.ParentRef = &accountRef
			updateRecord = true
		}
	}
	if fieldMap["CompanyName"].Display(order) != "" && customerRecord.CompanyName != fieldMap["CompanyName"].Display(order) {
		customerMod.CustomerMod.CompanyName = fieldMap["CompanyName"].Display(order)
		updateRecord = true
	}
	if fieldMap["Salutation"].Display(order) != "" && customerRecord.Salutation != fieldMap["Salutation"].Display(order) {
		customerMod.CustomerMod.Salutation = fieldMap["Salutation"].Display(order)
		updateRecord = true
	}
	if fieldMap["FirstName"].Display(order) != "" && customerRecord.FirstName != fieldMap["FirstName"].Display(order) {
		customerMod.CustomerMod.FirstName = fieldMap["FirstName"].Display(order)
		updateRecord = true
	}
	if fieldMap["MiddleName"].Display(order) != "" && customerRecord.MiddleName != fieldMap["MiddleName"].Display(order) {
		customerMod.CustomerMod.MiddleName = fieldMap["MiddleName"].Display(order)
		updateRecord = true
	}
	if fieldMap["LastName"].Display(order) != "" && customerRecord.LastName != fieldMap["LastName"].Display(order) {
		customerMod.CustomerMod.LastName = fieldMap["LastName"].Display(order)
		updateRecord = true
	}
	if fieldMap["JobTitle"].Display(order) != "" && customerRecord.JobTitle != fieldMap["JobTitle"].Display(order) {
		customerMod.CustomerMod.JobTitle = fieldMap["JobTitle"].Display(order)
		updateRecord = true
	}
	var addressMap = ReadFieldMapping("./fieldMaps/addressMapping.json")
	var address = Address{}
	if addressMap["BillAddress.Addr1"].Display(order) != customerRecord.BillAddress.Addr1 || addressMap["BillAddress.Addr2"].Display(order) != customerRecord.BillAddress.Addr2 || addressMap["BillAddress.Addr3"].Display(order) != customerRecord.BillAddress.Addr3 || addressMap["BillAddress.PostalCode"].Display(order) != customerRecord.BillAddress.PostalCode || addressMap["BillAddress.City"].Display(order) != customerRecord.BillAddress.City {
		address.Addr1 = addressMap["BillAddress.Addr1"].Display(order)
		address.Addr2 = addressMap["BillAddress.Addr2"].Display(order)
		address.Addr3 = addressMap["BillAddress.Addr3"].Display(order)
		address.Addr4 = addressMap["BillAddress.Addr4"].Display(order)
		address.Addr5 = addressMap["BillAddress.Addr5"].Display(order)
		address.City = addressMap["BillAddress.City"].Display(order)
		address.State = addressMap["BillAddress.State"].Display(order)
		address.PostalCode = addressMap["BillAddress.PostalCode"].Display(order)
		address.Country = addressMap["BillAddress.Country"].Display(order)
		address.Note = addressMap["BillAddress.Note"].Display(order)
		customerMod.CustomerMod.BillAddress = &address
		updateRecord = true
	}
	if addressMap["ShipAddress.Addr1"].Display(shipTo) != customerRecord.ShipAddress.Addr1 || addressMap["ShipAddress.Addr2"].Display(shipTo) != customerRecord.ShipAddress.Addr2 || addressMap["ShipAddress.Addr3"].Display(shipTo) != customerRecord.ShipAddress.Addr3 || addressMap["ShipAddress.PostalCode"].Display(shipTo) != customerRecord.ShipAddress.PostalCode || addressMap["ShipAddress.City"].Display(shipTo) != customerRecord.ShipAddress.City {
		address = Address{}
		address.Addr1 = addressMap["ShipAddress.Addr1"].Display(shipTo)
		address.Addr2 = addressMap["ShipAddress.Addr2"].Display(shipTo)
		address.Addr3 = addressMap["ShipAddress.Addr3"].Display(shipTo)
		address.Addr4 = addressMap["ShipAddress.Addr4"].Display(shipTo)
		address.Addr5 = addressMap["ShipAddress.Addr5"].Display(shipTo)
		address.City = addressMap["ShipAddress.City"].Display(shipTo)
		address.State = addressMap["ShipAddress.State"].Display(shipTo)
		address.PostalCode = addressMap["ShipAddress.PostalCode"].Display(shipTo)
		address.Country = addressMap["ShipAddress.Country"].Display(shipTo)
		address.Note = addressMap["ShipAddress.Note"].Display(shipTo)
		customerMod.CustomerMod.ShipAddress = &address
		updateRecord = true
	}
	//re add all existing shipTos and add the new address from the order, if it is new
	var shipToAddresses = make([]Address, 0)
	var addressMatch = false
	for _, address := range customerRecord.ShipToAddress {
		var sToAddress = Address{}
		sToAddress.Name = address.Name
		sToAddress.Addr1 = address.Addr1
		sToAddress.Addr2 = address.Addr2
		sToAddress.Addr3 = address.Addr3
		sToAddress.Addr4 = address.Addr4
		sToAddress.Addr5 = address.Addr5
		sToAddress.City = address.City
		sToAddress.State = address.State
		sToAddress.PostalCode = address.PostalCode
		sToAddress.Country = address.Country
		sToAddress.Note = address.Note
		sToAddress.DefaultShipTo = address.DefaultShipTo
		//shipToAddresses = append(shipToAddresses, sToAddress)
		if address.Addr1 == addressMap["ShipToAddress.Addr1"].Display(shipTo) {
			if address.Addr2 == addressMap["ShipToAddress.Addr2"].Display(shipTo) {
				if address.Addr3 == addressMap["ShipToAddress.Addr3"].Display(shipTo) {
					if address.City == addressMap["ShipToAddress.City"].Display(shipTo) {
						if address.PostalCode == addressMap["ShipToAddress.PostalCode"].Display(shipTo) {
							addressMatch = true
							continue
						}
					}
				}
			}
		}
	} //end address match
	if !addressMatch { //if the order address did not match any of the shipTos

		var sToAddress = Address{}
		sToAddress.Name = addressMap["ShipToAddress.Name"].Display(shipTo)
		sToAddress.Addr1 = addressMap["ShipToAddress.Addr1"].Display(shipTo)
		sToAddress.Addr2 = addressMap["ShipToAddress.Addr2"].Display(shipTo)
		sToAddress.Addr3 = addressMap["ShipToAddress.Addr3"].Display(shipTo)
		sToAddress.Addr4 = addressMap["ShipToAddress.Addr4"].Display(shipTo)
		sToAddress.Addr5 = addressMap["ShipToAddress.Addr5"].Display(shipTo)
		sToAddress.City = addressMap["ShipToAddress.City"].Display(shipTo)
		sToAddress.State = addressMap["ShipToAddress.State"].Display(shipTo)
		sToAddress.PostalCode = addressMap["ShipToAddress.PostalCode"].Display(shipTo)
		sToAddress.Country = addressMap["ShipToAddress.Country"].Display(shipTo)
		sToAddress.Note = addressMap["ShipToAddress.Note"].Display(shipTo)
		sToAddress.DefaultShipTo = addressMap["ShipToAddress.DefaultShipTo"].Display(shipTo)
		shipToAddresses = append(shipToAddresses, sToAddress)
		customerMod.CustomerMod.ShipToAddress = &shipToAddresses
		updateRecord = true
	}
	if fieldMap["Phone"].Display(order) != "" && customerRecord.Phone != fieldMap["Phone"].Display(order) {
		customerMod.CustomerMod.Phone = fieldMap["Phone"].Display(order)
		updateRecord = true
	}
	if fieldMap["AltPhone"].Display(order) != "" && customerRecord.AltPhone != fieldMap["AltPhone"].Display(order) {
		customerMod.CustomerMod.AltPhone = fieldMap["AltPhone"].Display(order)
		updateRecord = true
	}
	if fieldMap["Fax"].Display(order) != "" && customerRecord.Fax != fieldMap["Fax"].Display(order) {
		customerMod.CustomerMod.Fax = fieldMap["Fax"].Display(order)
		updateRecord = true
	}
	/* Do not update emails, as this is what is used to match customer records
	if fieldMap["Email"].Display(order) != "" && customerRecord.Email != fieldMap["Email"].Display(order) {
		customerMod.CustomerMod.Email = fieldMap["Email"].Display(order)
		updateRecord = true
	}*/
	if fieldMap["Cc"].Display(order) != "" && customerRecord.CC != fieldMap["Cc"].Display(order) {
		customerMod.CustomerMod.Cc = fieldMap["Cc"].Display(order)
		updateRecord = true
	}
	if fieldMap["Contact"].Display(order) != "" && customerRecord.Contact != fieldMap["Contact"].Display(order) {
		customerMod.CustomerMod.Contact = fieldMap["Contact"].Display(order)
		updateRecord = true
	}
	if fieldMap["AltContact"].Display(order) != "" && customerRecord.AltContact != fieldMap["AltContact"].Display(order) {
		customerMod.CustomerMod.AltContact = fieldMap["AltContact"].Display(order)
		updateRecord = true
	}

	var additionalContactRef = make([]AdditionalContactRef, 0)
	var aContactRef = AdditionalContactRef{}
	if fieldMap["AdditionalContactRef.ContactName"].Display(order) != "" {
		aContactRef.ContactName = fieldMap["AdditionalContactRef.ContactName"].Display(order)
		aContactRef.ContactValue = fieldMap["AdditionalContactRef.ContactValue"].Display(order)
		additionalContactRef = append(additionalContactRef, aContactRef)
		customerMod.CustomerMod.AdditionalContactRef = &additionalContactRef
		updateRecord = true
	}
	var contactsMod = make([]ContactsRet, 0)
	var cMod = ContactsRet{}
	if fieldMap["ContactsMod.ListID"].Display(order) != "" {
		cMod.ListID = fieldMap["ContactsMod.ListID"].Display(order)
		cMod.EditSequence = fieldMap["ContactsMod.EditSequence"].Display(order)
		cMod.Salutation = fieldMap["ContactsMod.Salutation"].Display(order)
		cMod.FirstName = fieldMap["ContactsMod.FirstName"].Display(order)
		cMod.MiddleName = fieldMap["ContactsMod.MiddleName"].Display(order)
		cMod.LastName = fieldMap["ContactsMod.LastName"].Display(order)
		cMod.JobTitle = fieldMap["ContactsMod.JobTitle"].Display(order)

		additionalContactRef = make([]AdditionalContactRef, 0)
		aContactRef = AdditionalContactRef{}
		aContactRef.ContactName = fieldMap["ContactsMod.AdditionalContactRef.ContactName"].Display(order)
		aContactRef.ContactValue = fieldMap["ContactsMod.AdditionalContactRef.ContactValue"].Display(order)
		additionalContactRef = append(additionalContactRef, aContactRef)
		cMod.AdditionalContactRef = &additionalContactRef

		contactsMod = append(contactsMod, cMod)
		customerMod.CustomerMod.ContactsMod = &contactsMod
	}
	if fieldMap["CustomerTypeRef.ListID"].Display(order)+fieldMap["CustomerTypeRef.FullName"].Display(order) != "" && customerRecord.CustomerTypeRef.FullName != fieldMap["CustomerTypeRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["CustomerTypeRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["CustomerTypeRef.FullName"].Display(order)
		customerMod.CustomerMod.CustomerTypeRef = &accountRef
		updateRecord = true
	}
	if fieldMap["TermsRef.ListID"].Display(order)+fieldMap["TermsRef.FullName"].Display(order) != "" && customerRecord.TermsRef.FullName != fieldMap["TermsRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["TermsRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["TermsRef.FullName"].Display(order)
		customerMod.CustomerMod.TermsRef = &accountRef
		updateRecord = true
	}
	if fieldMap["SalesRepRef.ListID"].Display(order)+fieldMap["SalesRepRef.ListID"].Display(order) != "" && customerRecord.SalesRepRef.FullName != fieldMap["SalesRepRef.ListID"].Display(order) {
		accountRef.ListID = fieldMap["SalesRepRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["SalesRepRef.FullName"].Display(order)
		customerMod.CustomerMod.SalesRepRef = &accountRef
		updateRecord = true
	}
	if fieldMap["SalesTaxCodeRef.ListID"].Display(order)+fieldMap["SalesTaxCodeRef.FullName"].Display(order) != "" && customerRecord.SalesTaxCodeRef.FullName != fieldMap["SalesTaxCodeRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["SalesTaxCodeRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["SalesTaxCodeRef.FullName"].Display(order)
		customerMod.CustomerMod.SalesTaxCodeRef = &accountRef
		updateRecord = true
	}
	if fieldMap["ItemSalesTaxRef.ListID"].Display(order)+fieldMap["ItemSalesTaxRef.FullName"].Display(order) != "" && customerRecord.ItemSalesTaxRef.FullName != fieldMap["ItemSalesTaxRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["ItemSalesTaxRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["ItemSalesTaxRef.FullName"].Display(order)
		customerMod.CustomerMod.ItemSalesTaxRef = &accountRef
		updateRecord = true
	}
	if fieldMap["ResaleNumber"].Display(order) != "" && customerRecord.ResaleNumber != fieldMap["ResaleNumber"].Display(order) {
		customerMod.CustomerMod.ResaleNumber = fieldMap["ResaleNumber"].Display(order)
		updateRecord = true
	}
	if fieldMap["AccountNumber"].Display(order) != "" && customerRecord.AccountNumber != fieldMap["AccountNumber"].Display(order) {
		customerMod.CustomerMod.AccountNumber = fieldMap["AccountNumber"].Display(order)
		updateRecord = true
	}
	if fieldMap["CreditLimit"].Display(order) != "" && customerRecord.CreditLimit != fieldMap["CreditLimit"].Display(order) {
		customerMod.CustomerMod.CreditLimit = fieldMap["CreditLimit"].Display(order)
		updateRecord = true
	}
	if fieldMap["PreferredPaymentMethodRef.ListID"].Display(order)+fieldMap["PreferredPaymentMethodRef.FullNameD"].Display(order) != "" && customerRecord.PreferredPaymentMethodRef.FullName != fieldMap["PreferredPaymentMethodRef.FullNameD"].Display(order) {
		accountRef.ListID = fieldMap["PreferredPaymentMethodRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["PreferredPaymentMethodRef.FullName"].Display(order)
		customerMod.CustomerMod.PreferredPaymentMethodRef = &accountRef
		updateRecord = true
	}
	if fieldMap["CreditCardInfo.CreditCardNumber"].Display(order) != "" {
		var ccInfo = CreditCardInfo{}
		ccInfo.CreditCardNumber = fieldMap["CreditCardInfo.CreditCardNumber"].Display(order)
		ccInfo.ExpirationMonth = fieldMap["CreditCardInfo.ExpirationMonth"].Display(order)
		ccInfo.ExpirationYear = fieldMap["CreditCardInfo.ExpirationYear"].Display(order)
		ccInfo.NameOnCard = fieldMap["CreditCardInfo.NameOnCard"].Display(order)
		ccInfo.CreditCardAddress = fieldMap["CreditCardInfo.CreditCardAddress"].Display(order)
		ccInfo.CreditCardPostalCode = fieldMap["CreditCardInfo.CreditCardPostalCode"].Display(order)
		customerMod.CustomerMod.CreditCardInfo = &ccInfo
		updateRecord = true
	}
	if fieldMap["JobStatus"].Display(order) != "" && customerRecord.JobStatus != fieldMap["JobStatus"].Display(order) {
		customerMod.CustomerMod.JobStatus = fieldMap["JobStatus"].Display(order)
		updateRecord = true
	}
	if fieldMap["JobStartDate"].Display(order) != "" && customerRecord.JobStartDate != fieldMap["JobStartDate"].Display(order) {
		customerMod.CustomerMod.JobStartDate = fieldMap["JobStartDate"].Display(order)
		updateRecord = true
	}
	if fieldMap["JobProjectedEndDate"].Display(order) != "" && customerRecord.JobProjectedEndDate != fieldMap["JobProjectedEndDate"].Display(order) {
		customerMod.CustomerMod.JobProjectedEndDate = fieldMap["JobProjectedEndDate"].Display(order)
		updateRecord = true
	}
	if fieldMap["JobEndDate"].Display(order) != "" && customerRecord.JobEndDate != fieldMap["JobEndDate"].Display(order) {
		customerMod.CustomerMod.JobEndDate = fieldMap["JobEndDate"].Display(order)
		updateRecord = true
	}
	if fieldMap["JobDesc"].Display(order) != "" && customerRecord.JobDesc != fieldMap["JobDesc"].Display(order) {
		customerMod.CustomerMod.JobDesc = fieldMap["JobDesc"].Display(order)
		updateRecord = true
	}
	if fieldMap["CustomerMod.JobTypeRef.ListID"].Display(order)+fieldMap["CustomerMod.JobTypeRef.FullName"].Display(order) != "" && customerRecord.JobTypeRef.FullName != fieldMap["CustomerMod.JobTypeRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["CustomerMod.JobTypeRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["CustomerMod.JobTypeRef.FullName"].Display(order)
		customerMod.CustomerMod.JobTypeRef = &accountRef
		updateRecord = true
	}
	if fieldMap["Notes"].Display(order) != "" && customerRecord.Notes != fieldMap["Notes"].Display(order) {
		customerMod.CustomerMod.Notes = fieldMap["Notes"].Display(order)
		updateRecord = true
	}
	if fieldMap["AdditionalNotesMod.NoteID"].Display(order)+fieldMap["AdditionalNotesMod.Note"].Display(order) != "" {
		var additionalNotesMod = AdditionalNotes{}
		additionalNotesMod.NoteID = fieldMap["AdditionalNotesMod.NoteID"].Display(order)
		additionalNotesMod.Note = fieldMap["AdditionalNotesMod.Note"].Display(order)
		customerMod.CustomerMod.AdditionalNotesMod = &additionalNotesMod
	}
	if fieldMap["PreferredDeliveryMethod"].Display(order) != "" && customerRecord.PreferredDeliveryMethod != fieldMap["PreferredDeliveryMethod"].Display(order) {
		customerMod.CustomerMod.PreferredDeliveryMethod = fieldMap["PreferredDeliveryMethod"].Display(order)
		updateRecord = true
	}
	if fieldMap["PriceLevelRef.ListID"].Display(order)+fieldMap["PriceLevelRef.FullName"].Display(order) != "" && customerRecord.PriceLevelRef.FullName != fieldMap["PriceLevelRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["PriceLevelRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["PriceLevelRef.FullName"].Display(order)
		customerMod.CustomerMod.PriceLevelRef = &accountRef
		updateRecord = true
	}
	if fieldMap["CurrencyRef.ListID"].Display(order)+fieldMap["CurrencyRef.FullName"].Display(order) != "" && customerRecord.CurrencyRef.FullName != fieldMap["CurrencyRef.FullName"].Display(order) {
		accountRef.ListID = fieldMap["CurrencyRef.ListID"].Display(order)
		accountRef.FullName = fieldMap["CurrencyRef.FullName"].Display(order)
		customerMod.CustomerMod.CurrencyRef = &accountRef
		updateRecord = true
	}

	b, err := xml.MarshalIndent(&customerMod, "", "  ")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer mod")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error marshalling customer mod")
		return
	}
	qbxmlWork.AppendWork(string(b))
	LoadTemplate(&tPath, qbxmlWork, &templateBuff)
	err = xml.EscapeText(&escapedQBXML, templateBuff.Bytes())
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerModQB")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error Escaping template xml in CustomerModQB")
		return
	}
	templateBuff.Reset()
	if updateRecord {
		Log.WithFields(logrus.Fields{"Nme": customerRecord.Name, "ListID": customerRecord.ListID, "Email": customerRecord.Email}).Error("Updating customer record")
		workChan <- WorkCTX{Work: escapedQBXML.String()}
		var customerModResponse = <-customerModResponseChan
		b, err = xml.MarshalIndent(&customerModResponse, "", "  ")
		if err != nil {
			Log.WithFields(logrus.Fields{"Error": err}).Error("Error marshalling customermod response")
			ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error marshalling customermod response")
		}
	} else {
		ErrLog.Error("not updating customer from customerMod")
	}
}
