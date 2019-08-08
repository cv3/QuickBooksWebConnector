package qbwcGo

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
	"github.com/amazingfly/cv3go"
)

//MappingObject is the struct to hold all the information for a single field mapping
type MappingObject []MapData

//MapData is one piece of a mapping object
type MapData struct {
	Data        string                   `json:"data"`        //the cv3 field name or the data to be hardcoded
	MappedField bool                     `json:"mappedField"` //indicates whether or not the data field is a hardcoded value or a cv3 field name
	SubMappings map[string]MappingObject `json:"subMappings"` //used to map the value found from the top level mapping to yet another value
}

//Display will display the data in the MappingObject in its desired format
func (mapObj MappingObject) Display(cv3Data ...*gabs.Container) string {
	return RecursiveMapCheck(mapObj, cv3Data...) //strings.TrimSpace(RecursiveMapCheck(mapObj, cv3Data...))
}

//RecursiveMapCheck will traverse the recursive mapping objects and return data as a string
func RecursiveMapCheck(mapObj MappingObject, cv3Data ...*gabs.Container) string {
	var displayBuf = bytes.Buffer{}
	//range over the map data elements to construct the desired string
	for _, mData := range mapObj {
		if mData.MappedField { //if this is a mapping
			//range over all gabs.Containers passed in, to check field mappings
			defaultFound := false
			for cDataIndex, cData := range cv3Data {
				data := strings.TrimSpace(CheckPath(mData.Data, cData)) //get the mapped data from the gabs.Container

				//if data != "" { //if the data is not nill
				if len(mData.SubMappings) > 0 { //if the data has subMappings to be resolved
					mObj, ok := mData.SubMappings[data] //get the desired mappingObject
					if ok {
						defaultFound = true
						displayBuf.WriteString(RecursiveMapCheck(mObj, cv3Data...))
					}
					if data == "" && !defaultFound && cDataIndex+1 == len(cv3Data) { //displayBuf.String() == "" { //Check for a default mapping as the matching mapData not found or empty
						for key, mObj := range mData.SubMappings {
							if strings.ToLower(key) == "default" {
								defaultFound = true
								displayBuf.WriteString(RecursiveMapCheck(mObj, cv3Data...))
							} // end if default
						} //end range loop on subMappings
					} //end else of OK
					if displayBuf.String() == "" && data != "" { //if non of the submappings mapped to anything, add the origional data
						defaultFound = true
						displayBuf.WriteString(data)
					}
				} else { //if there is no subMapping
					displayBuf.WriteString(data) //CheckPath(mData.Data, cData))
				}
				//} //end if data != ""
			} //end range loop on gabs.Containers
		} else { //hardcoded value
			if len(mData.SubMappings) > 0 { //if the data has subMappings to be resolved
				mObj, ok := mData.SubMappings[mData.Data] //get the desired mappingObject
				if ok {
					displayBuf.WriteString(RecursiveMapCheck(mObj, cv3Data...))
				}
				if mData.Data == "" { //displayBuf.String() == "" { //Check for a default mapping as the matching mapData not found or empty
					for key, mObj := range mData.SubMappings {
						if strings.ToLower(key) == "default" {
							displayBuf.WriteString(RecursiveMapCheck(mObj, cv3Data...))
						} // end if default
					} //end range loop on subMappings
				} //end else of OK
				if displayBuf.String() == "" { //if non of the submappings mapped to anything, add the origional data
					displayBuf.WriteString(mData.Data)
				}
			} else { // if there are no subMappings
				displayBuf.WriteString(mData.Data)
			}
		}
	}
	return displayBuf.String()
}

/*
//MappingObject is the struct to hold all the information for a single field mapping
type MappingObject []MapData

//MapData is one piece of a mapping object
type MapData struct {
	Data        string `json:"data"`        //the cv3 field name or the data to be hardcoded
	MappedField bool   `json:"mappedField"` //indicates whether or not the data field is a hardcoded value or a cv3 field name
}

//Display will display the data in the MappingObject in its desired format
func (mapObj MappingObject) Display(cv3Data ...*gabs.Container) string {
	var displayBuf = bytes.Buffer{}

	for _, mData := range mapObj {
		if mData.MappedField {
			for _, cData := range cv3Data {
				if CheckPath(mData.Data, cData) != "" {
					displayBuf.WriteString(CheckPath(mData.Data, cData))
					continue
				}
			}
		} else { //hardcoded value
			displayBuf.WriteString(mData.Data)
		}
	}
	return displayBuf.String()
}
*/

//QBXMLWork is the slice of finished templates to be added to the top level template
type QBXMLWork []string

//AppendWork checks if this is the first piece of work, as every piece of work after the first should have an extra newline character
func (qbxmlWork *QBXMLWork) AppendWork(w string) {
	if len(*qbxmlWork) > 0 {
		*qbxmlWork = append(*qbxmlWork, "\n")
	}
	*qbxmlWork = append(*qbxmlWork, w)
}

//WorkCTX is a struct that will hold both the work to be done, and the data used to create it
type WorkCTX struct {
	Attempted int    //keeps track of how many times this work was attempted
	Work      string //holds the excaped qbxml
	DataExts  []DataExtAddRq
	//TODO refactor Data to not be confused with DataExt
	Data          interface{}     //holds the struct that created the qbxml
	Order         *gabs.Container //holds the origional order information
	CV3Products   []cv3go.Product //holds the cv3 products used to make the qbxml
	Type          string          //type of qbxml request
	NoCustomer    bool            //set to true if there are problems with adding a customer name
	NoResendOrder bool            //used to teell CustomerAdd to not resent the order
}

//ItemGroupLine is a piece of the ItemGroupRet
type ItemGroupLine struct {
	ItemRef       AccountRef `xml:"ItemRef"`
	Quantity      string     `xml:"Quantity"`
	UnitOfMeasure string     `xml:"UnitOfMeasure"`
}

//ItemGroupRet is for the QBXML return type of ItemGroupRet
type ItemGroupRet struct {
	ItemBaseRet
	ItemDesc            string          `xml:"ItemDesc"`
	IsPrintItemsInGroup string          `xml:"IsPrintItemsInGroup"`
	ItemGroupLInes      []ItemGroupLine `xml:"ItemGroupLine"`
	BarCodeValue        string          `xml:"BarCodeValue"`
	UnitOfMeasureRef    AccountRef      `xml:"UnitOfMeasureRef"`
	SpecialItemType     string          `xml:"SpecialItemType"` //may have one of the following values: FinanceCharge, ReimbursableExpenseGroup, ReimbursableExpenseSubtotal
	ExternalGUID        string          `xml:"ExternalGUID"`
}

//ItemDiscountRet is for the QBXML return type of ItemDiscountRet
type ItemDiscountRet struct {
	ItemBaseRet
	ItemDesc        string     `xml:"ItemDesc"`
	SalesTaxCodeRef AccountRef `xml:"SalesTaxCodeRef"`
	DiscountRate    string     `xml:"DiscountRate"`
	AccountRef      AccountRef `xml:"AccountRef"`
}

//ItemSubtotalRet is for the QBXML return type of ItemSubtotalRet
type ItemSubtotalRet struct {
	ItemBaseRet
	ItemDesc string `xml:"ItemDesc"`
}

//ItemBaseRet has fields commont to all QBXML item return types
type ItemBaseRet struct {
	Name         string `xml:"Name"`
	FullName     string `xml:"FullName"`
	ListID       string `xml:"ListID"`
	TimeCreated  string `xml:"TimeCreated"`
	TimeModified string `xml:"TimeModified"`
	EditSequence string `xml:"EditSequence"`
	IsActive     string `xml:"IsActive"`
	Sublevel     string `xml:"Sublevel"`
}

//ItemOtherChargeRet is for the QBXML return type ItemOtherCharge
type ItemOtherChargeRet struct {
	ItemBaseRet
	SalesOrPurchase SalesOrPurchase `xml:"SalesOrPurchase"`
}

//BarCode holds BarCode information for items
type BarCode struct {
	BarCodeValue     string `xml:"BarCodeValue"`
	AssignEvenIfUsed string `xml:"AssignEvenIfUsed"`
	AllowOverride    string `xml:"AllowOverride"`
}

//ItemInventoryAdd is for the QBXML request to add an item to QB
type ItemInventoryAdd struct {
	Name                   string     `xml:"Name"`
	BarCode                BarCode    `xml:"BarCode"`
	IsActive               string     `xml:"IsActive"`
	ClassRef               AccountRef `xml:"ClassRef"`
	ParentRef              AccountRef `xml:"ParentRef"`
	ManufacturerPartNumber string     `xml:"ManufacturerPartNumber"`
	UnitOfMeasureSetRef    AccountRef `xml:"UnitOfMeasureSetRef"`
	SalesTaxCodeRef        AccountRef `xml:"SalesTaxCodeRef"`
	SalesDesc              string     `xml:"SalesDesc"`
	SalesPrice             string     `xml:"SalesPrice"`
	IncomeAccountRef       AccountRef `xml:"IncomeAccountRef"`
	PurchaseDesc           string     `xml:"PurchaseDesc"`
	PurchaseCost           string     `xml:"PurchaseCost"`
	COGSAccountRef         AccountRef `xml:"COGSAccountRef"`
	PrefVendorRef          AccountRef `xml:"PrefVendorRef"`
	AssetAccountRef        AccountRef `xml:"AssetAccountRef"`
	ReorderPoint           string     `xml:"ReorderPoint"`
	Max                    string     `xml:"Max"`
	QuantityOnHand         string     `xml:"QuantityOnHand"`
	TotalValue             string     `xml:"TotalValue"`
	InventoryDate          string     `xml:"InventoryDate"`
	ExternalGUID           string     `xml:"ExternalGUID"`
}

//ItemInventoryAddRs holds the information returned by Quick Books from and itemInventoryAdd request
type ItemInventoryAddRs struct {
	ResponseStatus
	ItemIventoryRet ItemIventoryRet `xml:"ItemInventoryRet"`
	ErrorRecovery   ErrorRecovery   `xml:"ErrorRecovery"`
}

//ItemIventoryRet is for the QBXML return type of ItemInventoryRet
type ItemIventoryRet struct {
	ItemBaseRet
	SalesDesc            string     `xml:"SalesDesc"`
	SalesPrice           string     `xml:"SalesPrice"`
	IncomeAccountRef     AccountRef `xml:"IncomeAccountRef"`
	PurchaseCost         string     `xml:"PurchaseCost"`
	COGSAccountRef       AccountRef `xml:"COGSAccountRef"`
	AssetAccountRef      AccountRef `xml:"AssetAccountRef"`
	QuantityOnHand       int        `xml:"QuantityOnHand"`
	AverageCost          string     `xml:"AverageCost"`
	QuantityOnOrder      int        `xml:"QuantityOnOrder"`
	QuantityOnSalesOrder int        `xml:"QuantityOnSalesOrder"`
}

//ItemGroupAdd holds the information for adding a group item
type ItemGroupAdd struct {
	Name                string          `xml:"Name"`
	BarCode             BarCode         `xml:"Barcode"`
	IsActive            string          `xml:"IsActive"`
	ItemDesc            string          `xml:"ItemDesc"`
	UnitOfMeasureSetRef AccountRef      `xml:"UnitOfMeasureSetRef"`
	IsPrintItemsInGroup string          `xml:"IsPrintItemsInGroup"`
	ExternalGUID        string          `xml:"ExternalGUID"`
	ItemGroupLine       []ItemGroupLine `xml:"ItemGroupLine"`
}

//ItemGroupAddRs is the struct to hold the information from ItemGroupAdd responses
type ItemGroupAddRs struct {
	ResponseStatus
	ItemGroupRet  ItemGroupRet  `xml:"ItemGroupRet"`
	ErrorRecovery ErrorRecovery `xml:"ErrorRecovery"`
}

//ItemNonInventoryRet is for the QBXML return type of ItemNonInventory
type ItemNonInventoryRet struct {
	ItemBaseRet
	ManufacturerPartNumber string          `xml:"ManufacturerPartNumber"`
	SalesOrPurchase        SalesOrPurchase `xml:"SalesOrPurchase"`
}

//SalesOrPurchase is the struct for the QBXML return item node SalesOrPurchase
type SalesOrPurchase struct {
	Desc         string     `xml:"Desc"`
	Price        string     `xml:"Price"`
	PricePercent string     `xml:"PricePercent"`
	AccountRef   AccountRef `xml:"AccountRef"`
}

//AccountRef Income, Expense, child of QBXML Items
type AccountRef struct {
	ListID   string `xml:"ListID,omitempty"`
	FullName string `xml:"FullName,omitempty"`
}

//SalesAndPurchase is the struct for the QBXML return item node SalesAndPurchase
type SalesAndPurchase struct {
	SalesDesc         string     `xml:"SalesDesc"`
	SalesPrice        string     `xml:"SalesPrice"`
	IncomeAccountRef  AccountRef `xml:"IncomeAccountRef"`
	PurchaseDesc      string     `xml:"PurchaseDesc"`
	PurchaseCost      string     `xml:"PurchaseCost"`
	ExpenseAccountRef AccountRef `xml:"ExpenseAccountRef"`
	PrefVendorRef     AccountRef `xml:"PrefVendorRef"`
}

//ItemServiceAdd is the struct that holds the item date when adding a service item to quickbooks
type ItemServiceAdd struct {
	Name                string           `xml:"Name"`
	BarCode             BarCode          `xml:"Barcode"`
	IsActive            string           `xml:"IsActive"`
	ClassRef            AccountRef       `xml:"ClassRef"`
	ParentRef           AccountRef       `xml:"ParentRef"`
	UnitOfMeasureSetRef AccountRef       `xml:"UnitOfMeasureSetRef"`
	SalesTaxCodeRef     AccountRef       `xml:"SalesTaxCodeRef"`
	SalesOrPurchase     SalesOrPurchase  `xml:"SalesOrPurchase"`
	SalesAndPurchase    SalesAndPurchase `xml:"SalesAndPurchase"`
	ExternalGUID        string           `xml:"ExternalGUID"`
}

//ItemServiceAddRs holds the information returned by Quick Books from and itemServiceAdd request
type ItemServiceAddRs struct {
	ResponseStatus
	ItemServiceRet ItemServiceRet `xml:"ItemServiceRet"`
	ErrorRecovery  ErrorRecovery  `xml:"ErrorRecovery"`
}

//ItemServiceRet is for the QBXML return type of ItemServiceRet
type ItemServiceRet struct {
	ItemBaseRet
	/*
		ListID               string     `xml:"ListID"`
		TimeCreated          string     `xml:"TimeCreated"`
		TimeModified         string     `xml:"TimeModified"`
		EditSequence         string     `xml:"EditSequence"`
		Name                 string     `xml:"Name"`
		FullName             string     `xml:"FullName"`
		IsActive             string     `xml:"IsActive"`
		Sublevel             string     `xml:"Sublevel"`
	*/
	SalesOrPurchase SalesOrPurchase `xml:"SalesOrPurchase"`
}

//ResponseStatus is a struct shared across multiple responses, holding status information
type ResponseStatus struct {
	StatusCode     string `xml:"statusCode,attr"`
	StatusSeverity string `xml:"statusSeverity,attr"`
	StatusMessage  string `xml:"statusMessage,attr"`
}

//ItemQueryCTX is the struct that holds the data for the qbxml ItemQuery request
type ItemQueryCTX struct {
	//BEGIN OR
	ListID []string `xml:"ListID"`
	//OR
	FullName []string `xml:"FullName"`
	//OR
	MaxReturned      int    `xml:"MaxReturned"`
	ActiveStatus     string `xml:"ActiveStatus"` //ActiveStatus may have one of the following values: ActiveOnly [DEFAULT], InactiveOnly, All
	FromModifiedDate string `xml:"FromModifiedDate"`
	ToModifiedDate   string `xml:"ToModifiedDate"`
	//BEGIN OR 2
	NameFilter NameFilter `xml:"NameFilter"`
	//OR 2
	NameRangeFilter NameRangeFilter `xml:"NameRangeFilter"`
	//END OR & END OR 2
	IncludeRetElement []string `xml:"IncludeRetElement"`
	OwnerID           []string `xml:"OwnderID"`
}

//NameFilter holds name filter data for ItemQueryCTX
type NameFilter struct {
	MatchCriterion string `xml:"MatchCriterion"` //REQUIRED MatchCriterion may have one of the following values: StartsWith, Contains, EndsWith
	Name           string `xml:"Name"`
}

//NameRangeFilter holds name range data for ItemQueryCTX
type NameRangeFilter struct {
	FromName string `xml:"FromName"`
	ToName   string `xml:"ToName"`
}

//ItemQueryRs is the struct for Quickbooks Items
type ItemQueryRs struct {
	ResponseStatus
	RequestID              string                `xml:"requestID,attr"`
	IteratorRemainingCount int                   `xml:"iteratorRemainingCount,attr,omitempty"`
	IteratorID             string                `xml:"iteratorID,attr,omitempty"`
	ItemServiceRets        []ItemServiceRet      `xml:"ItemServiceRet"`
	ItemInventoryRets      []ItemIventoryRet     `xml:"ItemInventoryRet"`
	ItemNonInventoryRets   []ItemNonInventoryRet `xml:"ItemNonInventoryRet"`
	ItemOtherChargeRets    []ItemOtherChargeRet  `xml:"ItemOtherChargeRet"`
	ItemSubtotalRets       []ItemSubtotalRet     `xml:"ItemSubtotalRet"`
	ItemDiscountRets       []ItemDiscountRet     `xml:"ItemDiscountRet"`
	ItemGroupRets          []ItemGroupRet        `xml:"ItemGroupRet"`
}

//ReceiveResponseXMLResponseCTX holds the context to fill the receiveResponseXML response
type ReceiveResponseXMLResponseCTX struct {
	Complete int //the % of completed work.  100 means all work is done
}

//ConnectionErrorResponseCTX holds the context to fill the connectionError Response
type ConnectionErrorResponseCTX struct {
	Result string //empty string //TODO what are valid values
}

//ServerVersionResponseCTX holts the context to fill the serverVersion response
type ServerVersionResponseCTX struct {
	Version string //a string representing the version of this server
}

//ClientVersionResponseCTX holds the context to fill the ClientVersion response
type ClientVersionResponseCTX struct {
	Result string //blank string, or string witha preceeding E:, W: or O: to indicate Error, Warning or OK
}

//ClientVersionCTX holds the client version information sent by QBWC
type ClientVersionCTX struct {
	StrVersion string `xml:"strVersion"`
}

//CloseConnectionResponseCTX holds the context to fill the closeConnection response
type CloseConnectionResponseCTX struct {
	CloseConnectionResult string //message to send when connection is closed
}

//CloseConnectionCTX holds the information sent by QBWC in the closeConnection soap request
type CloseConnectionCTX struct {
	Ticket string `xml:"ticket"`
}

//GetLastErrorResponseCTX holds the context to fill the getLastError response
type GetLastErrorResponseCTX struct {
	LastError string //blank if no error, string to indicate the error that occoured
}

//GetLastErrorCTX holds the information sent by QBWC in the GetLastError soap request
type GetLastErrorCTX struct {
	Ticket string `xml:"ticket"`
}

//SendRequestXMLResponseCTX holds the context to fill the sendRequestXML response
type SendRequestXMLResponseCTX struct {
	QBXML string //an escaped string of xml, that acts as insructions to Quick Books
}

//AuthenticateResponseCTX holds the data to fill the authenticate response
type AuthenticateResponseCTX struct {
	Ticket                  string //Session string: a string, none, nvu, or busy //TODO EXPEREMENT: QBWC talks about company file name, yet I don't believe it is correct. Origionally used 123a
	DelayUpdate             string //Can be a blank string. QBWC will pause for this many seconds
	EveryMinLowerLimit      string //Optional: Set the lower limit, in seconds, that the client can shoose to run updates
	MinimumRunEveryNSeconds string //Optional: TODO I'm not sure the difference
}

//ConnectionErrorCTX holds the data send by QBWC's connectionError call
type ConnectionErrorCTX struct {
	Ticket  string `xml:"ticket,omitempty"`
	HResult string `xml:"hresult,omitempty"`
	Message string `xml:"message,omitempty"`
}

//ReceiveResponseXMLCTX will hold the data send from QBWC's receiveResponseXML call
type ReceiveResponseXMLCTX struct {
	Ticket   string `xml:"ticket,omitempty"`
	Response []byte `xml:"response,omitempty"`
	HResult  string `xml:"hresult,omitempty"`
	Message  string `xml:"message,omitempty"`
}

//SessionCTX holds the session information
type SessionCTX struct {
	Ticket            string      `xml:"ticket,omitempty"` //session id returned in the authenticate response
	HCPResponse       []byte      `xml:"strHCPResponse,omitempty"`
	CompanyFileName   string      `xml:"strCompanyFileName,omitempty"`
	QBXMLCountry      string      `xml:"qbXMLCountry,omitempty"`
	QBXMLMajorVersion string      `xml:"qbXMLMajorVers,omitempty"`
	QBXMLMinorVersion string      `xml:"qbXMLMinorVers,omitempty"`
	QBXMLWorkQueue    []string    `xml:",omitempty"`
	QBXMLWorkChan     chan string `xml:",omitempty"`
}

//Node is the struct used to recursively read the xml
type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

//AuthenticateCTX is the struct used when qbwc sends authentication information
type AuthenticateCTX struct {
	XMLName  xml.Name `xml:"http://developer.intuit.com/ authenticate"`
	UserName string   `xml:"strUserName,omitempty"`
	Password string   `xml:"strPassword,omitempty"`
}

//AuthReturn is the return struct for authentication calls
type AuthReturn struct {
	Auth []string
}

//Address holds the billing addressss
type Address struct {
	Name          string `xml:"Name,omitempty"`
	Addr1         string `xml:"Addr1"`
	Addr2         string `xml:"Addr2,omitempty"`
	Addr3         string `xml:"Addr3,omitempty"`
	Addr4         string `xml:"Addr4,omitempty"`
	Addr5         string `xml:"Addr5,omitempty"`
	City          string `xml:"City"`
	State         string `xml:"State"`
	PostalCode    string `xml:"PostalCode"`
	Country       string `xml:"Country"`
	Note          string `xml:"Note,omitempty"`
	DefaultShipTo string `xml:"DefaultShipTo,omitempty"`
}

//CreditCardTxnInputInfo holds the credit card input information for CreditCardTxnInfo
type CreditCardTxnInputInfo struct {
	CreditCardNumber     string `xml:"CreditCardNumber"`
	ExpirationMonth      string `xml:"ExpirationMonth"`
	ExpirationYear       string `xml:"ExpirationYear"`
	NameOnCard           string `xml:"NameOnCard"`
	CreditCardAddress    string `xml:"CreditCardAddress"`
	CreditCardPostalCode string `xml:"CreditCardPostalCode"`
	CommercialCardCode   string `xml:"CommercialCardCode"`
	TransactionMode      string `xml:"TransactionMode"`   //TransactionMode may have one of the following values: CardNotPresent [DEFAULT], CardPresent
	CreditCardTxnType    string `xml:"CreditCardTxnType"` //CreditCardTxnType may have one of the following values: Authorization, Capture, Charge, Refund, VoiceAuthorization
}

//CreditCardTxnResultInfo holds the results from the credit card txn
type CreditCardTxnResultInfo struct {
	ResultCode            int/*required*/ `xml:"ResultCode"`
	ResultMessage         string/*required*/ `xml:"ResultMessage"`
	CreditCardTransID     string/*required*/ `xml:"CreditCardTransID"`
	MerchantAccountNumber string/*required*/ `xml:"MerchantAccountNumber"`
	AuthorizationCode     string                                   `xml:"AuthorizationCode"`
	AVSStreet             string                                   `xml:"AVSStreet"`             //AVSStreet may have one of the following values: Pass, Fail, NotAvailable
	AVSZip                string                                   `xml:"AVSZip"`                //AVSZip may have one of the following values: Pass, Fail, NotAvailable
	CardSecurityCodeMatch string                                   `xml:"CardSecurityCodeMatch"` //CardSecurityCodeMatch may have one of the following values: Pass, Fail, NotAvailable
	ReconBatchID          string                                   `xml:"ReconBatchID"`
	PaymentGroupingCode   int                                      `xml:"PaymentGroupingCode"`
	PaymentStatus         string/*required*/ `xml:"PaymentStatus"` //PaymentStatus may have one of the following values: Unknown, Completed
	TxnAuthorizationTime  string/*required*/ `xml:"TxnAuthorizationTime"`
	TxnAuthorizationStamp int    `xml:"TxnAuthorizationStamp"`
	ClientTransID         string `xml:"ClientTransID"`
}

//CreditCardTxnInfo holds credit card txn info
type CreditCardTxnInfo struct {
	CreditCardTxnInputInfo  CreditCardTxnInputInfo  `xml:"CreditCardTxnInputInfo"`
	CreditCardTxnResultInfo CreditCardTxnResultInfo `xml:"CreditCardTxnResultInfo"`
}

//DataExt Is part of the SalesReceiptLine items
type DataExt struct {
	OwnerID      string `xml:"OwnerID"`
	DataExtName  string `xml:"DataExtName"`
	DataExtValue string `xml:"DataExtValue"`
}

//DataExtRet is part of some receipt line item returns
type DataExtRet struct {
	DataExt
	DataExtType string `xml:"DataExtType"`
}

//SalesReceiptLineAdd one of the line item types for customerReceiptAdd
type SalesReceiptLineAdd struct {
	ItemRef       AccountRef `xml:"ItemRef"`
	Desc          string     `xml:"Desc,ommitempty"`
	Quantity      string     `xml:"Quantity,ommitempty"`
	UnitOfMeasure string     `xml:"UnitOfMeasure,ommitempty"`
	//2<!-- BEGIN OR -->
	Rate string `xml:"Rate,ommitempty"`
	//2<!-- OR -->
	RatePercent string `xml:"RatePercent,ommitempty"`
	//2<!-- OR -->
	PriceLevelRef AccountRef `xml:"PriceLevelRef,ommitempty"`
	//2<!-- END OR -->
	ClassRef                   AccountRef `xml:"ClassRef,ommitempty"`
	Amount                     string     `xml:"Amount,ommitempty"`
	OptionForPriceRuleConflict string     `xml:"OptionForPriceRuleConflict,ommitempty"`
	InventorySiteRef           AccountRef `xml:"InventorySiteRef,ommitempty"`
	InventorySiteLocationRef   AccountRef `xml:"InventorySiteLocationRef,ommitempty"`
	//2a<!-- BEGIN OR -->
	SerialNumber string `xml:"SerialNumber,ommitempty"`
	//2a<!-- OR -->
	LotNumber string `xml:"LotNumber,ommitempty"`
	//2a<!-- END OR -->
	ServiceDate            string            `xml:"ServiceDate,ommitempty"`
	SalesTaxCodeRef        AccountRef        `xml:"SalesTaxCodeRef,ommitempty"`
	OverrideItemAccountRef AccountRef        `xml:"OverrideItemAccount,ommitempty"`
	Other1                 string            `xml:"Other1,ommitempty"`
	Other2                 string            `xml:"Other2,ommitempty"`
	CreditCardTxnInfo      CreditCardTxnInfo `xml:"CreditCardTxnInfo,ommitempty"`
	DataExt                []DataExt         `xml:"DataExt,ommitempty"`
	Attributes             map[string]string //holds the cv3 attribute descriptions
}

//SalesReceiptPart is an itermediate struct used before we know if the cv3 item is a group or single
type SalesReceiptPart struct {
	Amount   string
	Quantity string
}

//SalesReceiptLineGroupAdd one of the line item types for customerReceiptAdd
type SalesReceiptLineGroupAdd struct {
	ItemGroupRef             AccountRef `xml:"ItemGroupRef"`
	Quantity                 string     `xml:"Quantity"`
	UnitOfMeasure            string     `xml:"UnitOfMeasure"`
	InventorySiteRef         AccountRef `xml:"InventorySiteRef"`
	InventorySiteLocationRef AccountRef `xml:"InventorySiteLocationRef"`
	DataExt                  []DataExt  `xml:"DataExt"`
}

/*
// BuildLineItems will take a *gabs.Container of CV3 items, itemFieldMap, skus a map of items to check, and mMap a map for keeping track of the used skus
func (receipt *SalesReceiptAdd) BuildLineItems(item *gabs.Container, itemFieldMap map[string]MappingObject, skus map[string]interface{}, workCTX *WorkCTX) { //*SalesReceiptAdd {
	var prod *SalesReceiptLineAdd
	var sku string
	//Check if this is an attribute product
	attr, err := item.Path("Attributes").Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute children in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute children in SalesReceiptAdd")
	}
	//Check if this is asubproduct
	subProds, err := item.Path("SubProducts").Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error making SubProducts children in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error making SubProducts children in SalesReceiptAdd")
	}
	if len(attr) > 1 { //if its an attribute product
		//Range over all attribute options and find the one that matches the SKU in the receipt
		for _, at := range attr {
			sku = CheckPath("SKU", at)
			pInterface, ok := skus[sku]
			if ok { //if the sku exists in the order data
				var p = pInterface.(*SalesReceiptLineAdd)
				if MatchAttributeCombinations(at, p.Attributes) == len(p.Attributes) {
					//Get the product data from the top level of the cv3 return
					prodInterface := AddReceiptItem(sku, prod, item, skus, &WorkCTX{}, itemFieldMap)
					if prodInterface != nil {
						prod = prodInterface.(*SalesReceiptLineAdd)
						//Get the product data from the specific attribute product
						prodInterface = AddReceiptItem(sku, prod, at, skus, workCTX, itemFieldMap)
						if prodInterface != nil {
							prod = prodInterface.(*SalesReceiptLineAdd)
							//append the attributes to the description
							var attrBuf = bytes.NewBufferString(prod.Desc)
							for _, attribute := range p.Attributes {
								attrBuf.WriteString(" ")
								attrBuf.WriteString(attribute)
							}
							prod.Desc = "" //attrBuf.String()
							//Add the product data to the quickbooks response data.
							receipt.SalesReceiptLineAdds = append(receipt.SalesReceiptLineAdds, *prod)
						}
					}
				}
			}
		}
	} else if len(subProds) > 1 {
		//Range over all attribute options and find the one that matches the SKU in the receipt
		for _, sProd := range subProds {
			sku = CheckPath("SKU", sProd)
			_, ok := skus[sku]
			if ok { //if the sku exists in the order data
				//Get the product data from the top level of the cv3 return
				prodInterface := AddReceiptItem(sku, prod, item, skus, &WorkCTX{}, itemFieldMap)
				if prodInterface != nil {
					prod = prodInterface.(*SalesReceiptLineAdd)
					//Get the product data from the specific sub product
					prodInterface = AddReceiptItem(sku, prod, sProd, skus, workCTX, itemFieldMap)
					if prodInterface != nil {
						prod = prodInterface.(*SalesReceiptLineAdd)
						//Add the product data to the quickbooks response data.
						receipt.SalesReceiptLineAdds = append(receipt.SalesReceiptLineAdds, *prod)
					}
				}
			}
		}
	} else { //not an attribute, or subproduct product
		sku = CheckPath("SKU", item)
		_, ok := skus[sku]
		if ok { //if the sku exists in the order data
			//Add all product data
			prodInterface := AddReceiptItem(sku, prod, item, skus, workCTX, itemFieldMap)
			prod = prodInterface.(*SalesReceiptLineAdd)
			//Add the product data to the quickbooks response data.
			receipt.SalesReceiptLineAdds = append(receipt.SalesReceiptLineAdds, *prod)
		}
	} //end else of attribute check
}
*/

//MatchAttributeCombinations loops through all the returned attribute combinations, and matches them against what was sent in the order information.
//This will return an int value to be compared with the length of attributes from the origional order information
func MatchAttributeCombinations(at *gabs.Container, pAttributes map[string]string) int {
	var attributeCombinationMatches = 0                       //count the attribute matches
	attrCombination, err := at.Path("Combination").Children() //get slice of all attributes
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute combination children in MatchAttributeCombinations")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute combination children in MatchAttributeCombinations")
	}
	//range over all attributes from the origional order information
	//then range over all attribute possibilities from the returned product
	//check if every attribute combination matches, otherwise it is te wrong attribute combination
	//then return the number of matches, to be compared with the length of the attribute combinations from the origional order information
	for _, orderAttributeDescription := range pAttributes {
		for _, atComb := range attrCombination {
			if orderAttributeDescription == CheckPath("content", atComb) {
				attributeCombinationMatches++
			}
		}
	}
	return attributeCombinationMatches
}

//DiscountCTX holds an orders discount contenxt
type DiscountCTX struct {
	Type              string
	TotalDiscount     float64
	Discount          float64
	RemainingDiscount float64
	SubTotal          float64
}

//SalesReceiptAdd is the struct for added sales receipts from cv3 to qb
type SalesReceiptAdd struct {
	DiscountCTX             *DiscountCTX
	DefMacro                string            `xml:"defmacro,attr"`
	CV3OrderID              string            //to keep track of order success
	ShipToIndex             int               //to keep track of orders with multiple shiptos
	CustomerRef             AccountRef        `xml:"CustomerRef,omitempty"`
	ClassRef                AccountRef        `xml:"ClassRef,omitempty"`
	TemplateRef             AccountRef        `xml:"TemplateRef,omitempty"`
	TxnDate                 string            `xml:"TxnDate,ommitempty"`
	RefNumber               string            `xml:"RefNumber,ommitempty"`
	BillAddress             Address           `xml:"BillAddress"`
	ShipAddress             Address           `xml:"ShipAddress"`
	IsPending               string            `xml:"IsPending,ommitempty"`
	CheckNumber             string            `xml:"CheckNumber,ommitempty"`
	PaymentMethodRef        AccountRef        `xml:"PaymentMethodRef,ommitempty"`
	DueDate                 string            `xml:"DueDate,ommitempty"`
	SalesRepRef             AccountRef        `xml:"SalesRepRef,ommitempty"`
	ShipDate                string            `xml:"ShipDate,ommitempty"`
	ShipMethodRef           AccountRef        `xml:"ShipMethodRef,ommitempty"`
	FOB                     string            `xml:"FOB,ommitempty"`
	ItemSalesTaxRef         AccountRef        `xml:"ItemSalesTaxRef,ommitempty"`
	Memo                    string            `xml:"Memo,ommitempty"`
	CustomerMsgRef          AccountRef        `xml:"CustomerMsgRef,ommitempty"`
	IsToBePrinted           bool              `xml:"IsToBePrinted,ommitempty"`
	IsToBeEmailed           bool              `xml:"IsToBeEmailed,ommitempty"`
	CustomerSalesTaxCodeRef AccountRef        `xml:"CustomerSalesTaxCodeRef,ommitempty"`
	DepositToAccountRef     AccountRef        `xml:"DepositToAccountRef,ommitempty"`
	CreditCardTxnInfo       CreditCardTxnInfo `xml:"CreditCardTxnInfo,ommitempty"`
	Other                   string            `xml:"Other,ommitempty"`
	ExchangeRate            float64           `xml:"ExchangeRate,ommitempty"`
	ExternalGUID            string            `xml:"ExternalGUID"` //ExternalGUID ragex = "0|(\{[0-9a-fA-F]{8}(\-([0-9a-fA-F]{4})){3}\-[0-9a-fA-F]{12}\})"
	//1<!-- BEGIN OR --> <<<probably
	SalesReceiptLineAdds []SalesReceiptLineAdd `xml:"SalesReceiptLineAdd,ommitempty"`
	//1<!-- OR -->
	SalesReceiptLineGroupAdd []SalesReceiptLineGroupAdd `xml:"SalesReceiptLineGroupAdd,ommitempty"`
	//1<!-- END OR -->
	IncludeRetElement string `xml:"IncludeRetElement,ommitempty"`
}

//OrderOrReceiptAdd Interface for SalesOrders or SalesReceiptsAdds
type OrderOrReceiptAdd interface {
	BuildLineItems()
}

/*
// BuildLineItems will take a *gabs.Container of CV3 items, itemFieldMap, skus a map of items, checks for sub or attribute products and adds them to the SalesOrderAdd object
func (order *SalesOrderAdd) BuildLineItems(item *gabs.Container, itemFieldMap map[string]MappingObject, skus map[string]interface{}, workCTX *WorkCTX) { //*SalesOrderAdd {
	var prod *SalesOrderLineAdd
	var sku string
	//Check if this is an attribute product
	attr, err := item.Path("Attributes").Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute children in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error making attribute children in SalesReceiptAdd")
	}
	//Check if this is asubproduct
	subProds, err := item.Path("SubProducts").Children()
	if err != nil {
		Log.WithFields(logrus.Fields{"Error": err}).Error("Error making SubProducts children in SalesReceiptAdd")
		ErrLog.WithFields(logrus.Fields{"Error": err}).Error("Error making SubProducts children in SalesReceiptAdd")
	}
	if len(attr) > 1 { //if its an attribute product
		//Range over all attribute options and find the one that matches the SKU in the order
		for _, at := range attr {
			sku = CheckPath("SKU", at)
			pInterface, ok := skus[sku]
			if ok { //if the sku exited in the order info
				var p = pInterface.(*SalesOrderLineAdd)
				if MatchAttributeCombinations(at, p.Attributes) == len(p.Attributes) {
					//Get the product data from the top level of the cv3 return
					prodInterface := AddOrderItem(sku, prod, item, skus, &WorkCTX{}, itemFieldMap)
					if prodInterface != nil {
						prod = prodInterface.(*SalesOrderLineAdd)
						//Get the product data from the specific attribute product
						prodInterface = AddOrderItem(sku, prod, at, skus, workCTX, itemFieldMap)
						if prodInterface != nil {
							prod = prodInterface.(*SalesOrderLineAdd)
							//append the attributes to the description

							//Add the product data to the quickbooks response data.
							order.SalesOrderLineAdds = append(order.SalesOrderLineAdds, *prod)
						}
					}
				}
			}
		}
	} else if len(subProds) > 1 { //if its a sub product
		//Range over all sub product options and find the one that matches the SKU in the order
		for _, sProd := range subProds {
			sku = CheckPath("SKU", sProd)
			_, ok := skus[sku]
			if ok { //if the sku matches what was in the order data
				//Get the product data from the top level of the cv3 product return
				prodInterface := AddOrderItem(sku, prod, item, skus, &WorkCTX{}, itemFieldMap)
				if prodInterface != nil {
					prod = prodInterface.(*SalesOrderLineAdd)
					//get the product data from the specific sub product
					prodInterface = AddOrderItem(sku, prod, sProd, skus, workCTX, itemFieldMap)
					if prodInterface != nil {
						prod = prodInterface.(*SalesOrderLineAdd)
						//Add the product data to the quickbooks response data.
						order.SalesOrderLineAdds = append(order.SalesOrderLineAdds, *prod)
					}
				}
			}
		}
	} else { //not an attribute, or subproduct product
		sku = CheckPath("SKU", item)
		_, ok := skus[sku]
		if ok { //if the sku was in the order data
			//Add all product data
			prodInterface := AddOrderItem(sku, prod, item, skus, workCTX, itemFieldMap)
			prod = prodInterface.(*SalesOrderLineAdd)
			//Add the product data to the quickbooks response data.
			order.SalesOrderLineAdds = append(order.SalesOrderLineAdds, *prod)
		}
	} //end else of attribute check
}
*/

//SalesOrderAdd is the struct to hold the variables for the qbxml call
type SalesOrderAdd struct {
	DiscountCTX             *DiscountCTX
	DefMacro                string `xml:"defmacro,attr"`
	CV3OrderID              string //to keep track of order success
	ShipToIndex             int
	CustomerRef             AccountRef               `xml:"CustomerRef"`
	ClassRef                AccountRef               `xml:"ClassRef"`
	TemplateRef             AccountRef               `xml:"TemplateRef"`
	TxnDate                 string                   `xml:"TxnDate"`
	RefNumber               string                   `xml:"RefNumber"`
	BillAddress             Address                  `xml:"BillAddress"`
	ShipAddress             Address                  `xml:"ShipAddress"`
	PONumber                string                   `xml:"PONumber"`
	TermsRef                AccountRef               `xml:"TermsRef"`
	DueDate                 string                   `xml:"DueDate"`
	SalesRepRef             AccountRef               `xml:"SalesRepRef"`
	FOB                     string                   `xml:"FOB"`
	ShipDate                string                   `xml:"ShipDate"`
	ShipMethodRef           AccountRef               `xml:"ShipMethodRef"`
	ItemSalesTaxRef         AccountRef               `xml:"ItemSalesTaxRef"`
	IsManuallyClosed        string                   `xml:"IsManuallyClosed"`
	Memo                    string                   `xml:"Memo"`
	CustomerMsgRef          AccountRef               `xml:"CustomerMsgRef"`
	IsToBePrinted           string                   `xml:"IsToBePrinted"`
	IsToBeEmailed           string                   `xml:"IsToBeEmailed"`
	CustomerSalesTaxCodeRef AccountRef               `xml:"CustomerSalesTaxCodeRef"`
	Other                   string                   `xml:"Other"`
	ExchangeRate            string                   `xml:"ExchangeRate"`
	ExternalGUID            string                   `xml:"ExternalGUID"`
	SalesOrderLineAdds      []SalesOrderLineAdd      `xml:"SalesOrderLineAdds"`
	SalesOrderLineGroupAdds []SalesOrderLineGroupAdd `xml:"SalesOrderLineGroupAdd"`
}

//SalesOrderLineAdd holds the information about the items in the sales order
type SalesOrderLineAdd struct {
	ItemRef                    AccountRef        `xml:"ItemRef"`
	Desc                       string            `xml:"Desc"`
	Quantity                   string            `xml:"Quantity"`
	UnitOfMeasure              string            `xml:"UnitOfMeasure"`
	Rate                       string            `xml:"Rate"`
	RatePercent                string            `xml:"RatePercent"`
	PriceLevelRef              AccountRef        `xml:"PriceLevelRef"`
	ClassRef                   AccountRef        `xml:"ClassRef"`
	Amount                     string            `xml:"Amount"`
	OptionForPriceRuleConflict string            `xml:"OptionForPriceRuleConflict"` //OptionForPriceRuleConflict may have one of the following values: Zero, BasePrice -->
	InventorySiteRef           AccountRef        `xml:"InventorySiteRef"`
	InventorySiteLocationRef   AccountRef        `xml:"InventorySiteLocationRef"`
	SerialNumber               string            `xml:"SerialNumber"`
	LotNumber                  string            `xml:"LotNumber"`
	SalesTaxCodeRef            AccountRef        `xml:"SalesTaxCodeRef"`
	IsManuallyClosed           string            `xml:"IsManuallyClosed"`
	Other1                     string            `xml:"Other1"`
	Other2                     string            `xml:"Other2"`
	DataExt                    []DataExt         `xml:"DataExt"`
	Attributes                 map[string]string //holds the cv3 attribute descriptions
}

//SalesOrderLineGroupAdd hold the information for the group items in a salesOrder
type SalesOrderLineGroupAdd struct {
	ItemGroupRef             AccountRef `xml:"ItemGroupRef"`
	Quantity                 string     `xml:"Quantity"`
	UnitOfMeasure            string     `xml:"UnitOfMeasure"`
	InventorySiteRef         AccountRef `xml:"InventorySiteRef"`
	InventorySiteLocationRef AccountRef `xml:"InventorySiteLocationRef"`
	DataExt                  []DataExt  `xml:"DataExt"`
}

//ErrorRecovery hold information on tthe error that occoured
type ErrorRecovery struct {
	ListID       string `xml:"ListID"`
	OwnerID      string `xml:"OwnerID"`
	TxnID        string `xml:"TxnID"`
	TxnNumber    string `xml:"TxnNumber"`
	EditSequence string `xml:"EditSequence"`
	ExternalGUID string `xml:"ExternalGUID"`
}

//SalesOrderAddRs is the top level struct for holding the information returned for a salesOrderAddRq
type SalesOrderAddRs struct {
	ResponseStatus
	SalesOrderRet SalesReceiptRet `xml:"SalesOrderRet"`
	ErrorRecovery ErrorRecovery   `xml:"ErrorRecovery"`
}

//SalesOrderRet holds the information from the sales Order add
type SalesOrderRet struct {
	TxnID                     string                 `xml:"TxnID"`
	TimeCreated               string                 `xml:"TimeCreated"`
	TimeModified              string                 `xml:"TimeModified"`
	EditSequence              string                 `xml:"EditSequence"`
	TxnNumber                 string                 `xml:"TxnNumber"`
	CustomerRef               AccountRef             `xml:"CustomerRef"`
	ClassRef                  AccountRef             `xml:"ClassRef"`
	TemplateRef               AccountRef             `xml:"TemplateRef"`
	TxnDate                   string                 `xml:"TxnDate"`
	RefNumber                 string                 `xml:"RefNumber"`
	BillAddress               Address                `xml:"BillAddress"`
	BillAddressBlock          Address                `xml:"BillAddressBlock"`
	ShipAddress               Address                `xml:"ShipAddress"`
	ShipAddressBlock          Address                `xml:"ShipAddressBlock"`
	PONumber                  string                 `xml:"PONumber"`
	TermsRef                  AccountRef             `xml:"TermsRef"`
	DueDate                   string                 `xml:"DueDate"`
	SalesRepRef               AccountRef             `xml:"SalesRepRef"`
	FOB                       string                 `xml:"FOB"`
	ShipDate                  string                 `xml:"ShipDate"`
	ShipMethodRef             AccountRef             `xml:"ShipMethodRef"`
	SubTotal                  string                 `xml:"SubTotal"`
	ItemSalesTaxRef           AccountRef             `xml:"ItemSalesTaxRef"`
	SalesTaxPercentage        string                 `xml:"SalesTaxPercentage"`
	SalesTaxTotal             string                 `xml:"SalesTaxTotal"`
	TotalAmount               string                 `xml:"TotalAmount"`
	CurrencyRef               AccountRef             `xml:"CurrencyRef"`
	ExchangeRate              string                 `xml:"ExchangeRate"`
	TotalAmountInHomeCurrency string                 `xml:"TotalAmountInHomeCurrency"`
	IsManuallyClosed          string                 `xml:"IsManUallyClosed"`
	IsFullyInvoiced           string                 `xml:"IsFullyInvoiced"`
	Memo                      string                 `xml:"Memo"`
	CustomerMsgRef            AccountRef             `xml:"CustomerMsgRef"`
	IsToBePrinted             string                 `xml:"IsToBePrinted"`
	IsToBeEmailed             string                 `xml:"IsToBeEmailed"`
	CustomerSalesTaxCodeRef   AccountRef             `xml:"CustomerSalesTaxCodeRef"`
	Other                     string                 `xml:"Other"`
	ExternalGUID              string                 `xml:"ExternalGUID"`
	LinkedTxn                 []LinkedTxn            `xml:"LinkedTxn"`
	SalesOrderLineRet         SalesOrderLineRet      `xml:"SalesReceiptLineRet"`
	SalesOrderLineGroupRet    SalesOrderLineGroupRet `xml:"SalesReceiptLineGroupRet"`
	DataExtRet                DataExtRet             `xml:"DataExtRet"`
}

//LinkedTxn is the struct to hold the data returned in a salesOrderAdd's LinkedTxnRet
type LinkedTxn struct {
	TxnID     string `xml:"TxnID"`
	TxnType   string `xml:"TxnType"` //TxnType may have one of the following values: ARRefundCreditCard, Bill, BillPaymentCheck, BillPaymentCreditCard, BuildAssembly, Charge, Check, CreditCardCharge, CreditCardCredit, CreditMemo, Deposit, Estimate, InventoryAdjustment, Invoice, ItemReceipt, JournalEntry, LiabilityAdjustment, Paycheck, PayrollLiabilityCheck, PurchaseOrder, ReceivePayment, SalesOrder, SalesReceipt, SalesTaxPaymentCheck, Transfer, VendorCredit, YTDAdjustment
	TxnDate   string `xml:"TxnDate"`
	RefNumber string `xml:"RefNumber"`
	LinkType  string `xml:"LinkType"` //LinkType may have one of the following values: AMTTYPE, QUANTYPE
	Amount    string `xml:"Amount"`
}

//SalesOrderLineRet hold the returned information from a sales order line item
type SalesOrderLineRet struct {
	TxnLineID                string     `xml:"TxnLineID"`
	ItemRef                  AccountRef `xml:"ItemFef"`
	Desc                     string     `xml:"Desc"`
	Quantity                 string     `xml:"Quantity"`
	UnitOfMeasure            string     `xml:"UnitOfMeasure"`
	OverrideUOMSetRef        AccountRef `xml:"OverrideUOMSetRef"`
	Rate                     string     `xml:"Rate"`
	RatePercent              string     `xml:"RatePercent"`
	ClassRef                 AccountRef `xml:"ClassRef"`
	Amount                   string     `xml:"Amount"`
	InventorySiteRef         AccountRef `xml:"InventorySiteRef"`
	InventorySiteLocationRef AccountRef `xml:"InventorySiteLocationRef"`
	SerialNumber             string     `xml:"SerialNumber"`
	LotNumber                string     `xml:"LotNumber"`
	SalesTaxCodeRef          AccountRef `xml:"SalesTaxCodeRef"`
	Invoiced                 string     `xml:"Invoiced"`
	IsManuallyClosed         string     `xml:"IsManuallyClosed"`
	Other1                   string     `xml:"Other1"`
	Other2                   string     `xml:"Other2"`
	DataExtRet               DataExtRet `xml:"DataExtRet"`
}

//SalesOrderLineGroupRet hold the returned information from a sales order line group item
type SalesOrderLineGroupRet struct {
	TxnLineID           string                `xml:"TxnLineID"`
	ItemGroupRef        AccountRef            `xml:"ItemGroupRef"`
	Desc                string                `xml:"Desc"`
	Quantity            string                `xml:"Quantity"`
	UnitOfMeasure       string                `xml:"UnitOfMeasure"`
	OverrideUOMSetRef   AccountRef            `xml:"OverrideUOMSetRef"`
	IsPrintItemsInGroup string                `xml:"IsPrintItemsInGroup"`
	TotalAmount         string                `xml:"TotalAmount"`
	SalesOrderLineRet   []SalesReceiptLineRet `xml:"SalesReceiptLineRet"`
	DataExtRet          DataExtRet            `xml:"DataExtRet"`
}

//SalesReceiptAddRs is the top level struct for holding the information returned for a salesReceiptAddRq
type SalesReceiptAddRs struct {
	ResponseStatus
	SalesReceiptRet SalesReceiptRet `xml:"SalesReciptRet"`
	ErrorRecovery   ErrorRecovery   `xml:"ErrorRecovery"`
}

//SalesReceiptRet holds the information from the sales receipt add
type SalesReceiptRet struct {
	TxnID                     string                   `xml:"TxnID"`
	TimeCreated               string                   `xml:"TimeCreated"`
	TimeModified              string                   `xml:"TimeModified"`
	EditSequence              string                   `xml:"EditSequence"`
	TxnNumber                 string                   `xml:"TxnNumber"`
	CustomerRef               AccountRef               `xml:"CustomerRef"`
	ClassRef                  AccountRef               `xml:"ClassRef"`
	TemplateRef               AccountRef               `xml:"TemplateRef"`
	TxnDate                   string                   `xml:"TxnDate"`
	RefNumber                 string                   `xml:"RefNumber"`
	BillAddress               Address                  `xml:"BillAddress"`
	BillAddressBlock          Address                  `xml:"BillAddressBlock"`
	ShipAddress               Address                  `xml:"ShipAddress"`
	ShipAddressBlock          Address                  `xml:"ShipAddressBlock"`
	IsPending                 string                   `xml:"IsPending"`
	CheckNumber               string                   `xml:"CheckNumber"`
	PaymentMethodRef          AccountRef               `xml:"PaymentMethodRef"`
	DueDate                   string                   `xml:"DueDate"`
	SalesRepRef               AccountRef               `xml:"SalesRepRef"`
	ShipDate                  string                   `xml:"ShipDate"`
	ShipMethodRef             AccountRef               `xml:"ShipMethodRef"`
	FOB                       string                   `xml:"FOB"`
	SubTotal                  string                   `xml:"SubTotal"`
	ItemSalesTaxRef           AccountRef               `xml:"ItemSalesTaxRef"`
	SalesTaxPercentage        string                   `xml:"SalesTaxPercentage"`
	SalesTaxTotal             string                   `xml:"SalesTaxTotal"`
	TotalAmount               string                   `xml:"TotalAmount"`
	CurrencyRef               AccountRef               `xml:"CurrencyRef"`
	ExchangeRate              string                   `xml:"ExchangeRate"`
	TotalAmountInHomeCurrency string                   `xml:"TotalAmountInHomeCurrency"`
	Memo                      string                   `xml:"Memo"`
	CustomerMsgRef            AccountRef               `xml:"CustomerMsgRef"`
	IsToBePrinted             string                   `xml:"IsToBePrinted"`
	IsToBeEmailed             string                   `xml:"IsToBeEmailed"`
	CustomerSalesTaxCodeRef   AccountRef               `xml:"CustomerSalesTaxCodeRef"`
	DepositToAccountRef       AccountRef               `xml:"DepositToAccountRef"`
	CreditCardTxnInfo         CreditCardTxnInfo        `xml:"CreditCardTxnInfo"`
	Other                     string                   `xml:"Other"`
	ExternalGUID              string                   `xml:"ExternalGUID"`
	SalesReceiptLineRet       SalesReceiptLineRet      `xml:"SalesReceiptLineRet"`
	SalesReceiptLineGroupRet  SalesReceiptLineGroupRet `xml:"SalesReceiptLineGroupRet"`
	DataExtRet                DataExtRet               `xml:"DataExtRet"`
}

//SalesReceiptLineRet hold the returned information from a sales receipt line item
type SalesReceiptLineRet struct {
	TxnLineID                string            `xml:"TxnLineID"`
	ItemRef                  AccountRef        `xml:"ItemFef"`
	Desc                     string            `xml:"Desc"`
	Quantity                 string            `xml:"Quantity"`
	UnitOfMeasure            string            `xml:"UnitOfMeasure"`
	OverrideUOMSetRef        AccountRef        `xml:"OverrideUOMSetRef"`
	Rate                     string            `xml:"Rate"`
	RatePercent              string            `xml:"RatePercent"`
	ClassRef                 AccountRef        `xml:"ClassRef"`
	Amount                   string            `xml:"Amount"`
	InventorySiteRef         AccountRef        `xml:"InventorySiteRef"`
	InventorySiteLocationRef AccountRef        `xml:"InventorySiteLocationRef"`
	SerialNumber             string            `xml:"SerialNumber"`
	LotNumber                string            `xml:"LotNumber"`
	ServiceDate              string            `xml:"ServiceDate"`
	SalesTaxCodeRef          AccountRef        `xml:"SalesTaxCodeRef"`
	Other1                   string            `xml:"Other1"`
	Other2                   string            `xml:"Other2"`
	CreditCardTxnInfo        CreditCardTxnInfo `xml:"CreditCardTxnInfo"`
	DataExtRet               DataExtRet        `xml:"DataExtRet"`
}

//SalesReceiptLineGroupRet hold the returned information from a sales receipt line group item
type SalesReceiptLineGroupRet struct {
	TxnLineID           string                `xml:"TxnLineID"`
	ItemGroupRef        AccountRef            `xml:"ItemGroupRef"`
	Desc                string                `xml:"Desc"`
	Quantity            string                `xml:"Quantity"`
	UnitOfMeasure       string                `xml:"UnitOfMeasure"`
	OverrideUOMSetRef   AccountRef            `xml:"OverrideUOMSetRef"`
	IsPrintItemsInGroup string                `xml:"IsPrintItemsInGroup"`
	TotalAmount         string                `xml:"TotalAmount"`
	SalesReceiptLineRet []SalesReceiptLineRet `xml:"SalesReceiptLineRet"`
	DataExtRet          DataExtRet            `xml:"DataExtRet"`
}

//Config is the struct to hold the config data
type Config struct {
	OrderType              string         `json:"orderType"`   //SalesReceipt, or SalesOrder are current options, case insensitive
	InitQBItems            bool           `json:"initQBItems"` //adds the Shipping item to QuickBooks
	AutoImportCV3Items     bool           `json:"autoImportCV3Items"`
	ItemUpdates            ItemUpdates    `json:"itemUpdates"`
	ListenPort             string         `json:"listenPort"`
	Logging                Logging        `json:"logging"`
	CV3Credentials         CV3Credentials `json:"cv3Credentials"`
	QBWCCredentials        Credentials    `json:"qbwcCredentials"`
	ServerVersion          string         `json:"serverVersion"`
	QBWCVersion            string         `json:"qbwcVersion"`
	CloseConnectionMessage string         `json:"closeConnectionMessage"`
	MaxWorkAttempts        int            `json:"maxWorkAttempts"`
	ConfirmOrders          bool           `json:"confirmOrders"`
	DataExtActive          bool           `json:"dataExtActive"`
	NameArrangement        struct {
		First           string `json:"first"`
		SeperatorString string `json:"seperatorString"`
		Last            string `json:"last"`
	} `json:"nameArrangement"`
}

//ItemUpdates is part of the config to keep track of item updates
type ItemUpdates struct {
	UpdateCV3Items bool   `json:"updateCV3Items"` //if set to true, QB items will be uploaded to V3
	LastUpdate     string `json:"lastUpdate"`
	UpdateNewOnly  bool   `json:"updateNewOnly"`
}

//Logging holds the logging data for the config
type Logging struct {
	Level      string `json:"Level"` //debug info warn error
	OutputPath string `json:"outputPath"`
}

//Credentials holds the credential information for cv3
type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

//CV3Credentials holds the credentials for cv3 authentication
type CV3Credentials struct {
	Credentials
	ServiceID string `json:"serviceID"`
}

//OrderSuccessTracker is used as a map to keep track of what orders were impored successfully
type OrderSuccessTracker struct {
	OrderID               string //cv3 orderID
	ShipToLength          int    //how many shipTos in the order
	SuccessCount          int    //how many successful shipTos in the order
	ShipToSuccessTrackers map[int]ShipToSuccessTracker
}

//ShipToSuccessTracker is used as a slice in of OrderSuccessTracker, to keep track of which shipTos have been successfully added to Quick Bookis
type ShipToSuccessTracker struct {
	OrderID        string
	Index          int
	Success        bool
	QBErrorCode    string
	QBErrorMessage string
}

//CustomerAddRq is the top level struct for holding the CustomerAddRq data
type CustomerAddRq struct {
	Name                      string                 `xml:"Name"`
	IsActive                  string                 `xml:"IsActive"`
	ClassRef                  AccountRef             `xml:"ClassRef"`
	ParentRef                 AccountRef             `xml:"ParentRef"`
	CompanyName               string                 `xml:"CompanyName"`
	Salutation                string                 `xml:"Salutation"`
	FirstName                 string                 `xml:"FirstName"`
	MiddleName                string                 `xml:"MiddleName"`
	LastName                  string                 `xml:"LastName"`
	JobTitle                  string                 `xml:"JobTitle"`
	BillAddress               Address                `xml:"BillAddress"`
	ShipAddress               Address                `xml:"ShipAddress"`
	ShipToAddress             []QBShipToAddress      `xml:"ShipToAddress"`
	Phone                     string                 `xml:"Phone"`
	AltPhone                  string                 `xml:"AltPhone"`
	Fax                       string                 `xml:"Fax"`
	Email                     string                 `xml:"Email"`
	Cc                        string                 `xml:"Cc"`
	Contact                   string                 `xml:"Contact"`
	AltContact                string                 `xml:"AltContact"`
	AdditionalContactRef      []AdditionalContactRef `xml:"AdditionalContactRef"`
	Contacts                  []Contacts             `xml:"Contacts"`
	CustomerTypeRef           AccountRef             `xml:"CustomerTypeRef"`
	TermsRef                  AccountRef             `xml:"TermsRef"`
	SalesRepRef               AccountRef             `xml:"SalesRepRef"`
	OpenBalance               string                 `xml:"OpenBalance"`
	OpenBalanceDate           string                 `xml:"OpenBalanceDate"`
	SalesTaxCodeRef           AccountRef             `xml:"SalesTaxCodeRef"`
	ItemSalesTaxRef           AccountRef             `xml:"ItemSalesTaxRef"`
	ResaleNumber              string                 `xml:"ResaleNumber"`
	AccountNumber             string                 `xml:"AccountNumber"`
	CreditLimit               string                 `xml:"CreditLimit"`
	PreferredPaymentMethodRef AccountRef             `xml:"PreferredPaymentMethodRef"`
	CreditCardInfo            CreditCardInfo         `xml:"CreditCardInfo"`
	JobStatus                 string                 `xml:"JobStatus"`
	JobStartDate              string                 `xml:"JobStartDate"`
	JobProjectedEndDate       string                 `xml:"JobProjectedEndDate"`
	JobEndDate                string                 `xml:"JobEndDate"`
	JobDesc                   string                 `xml:"JobDesc"`
	JobTypeRef                AccountRef             `xml:"JobTypeRef"`
	Notes                     string                 `xml:"Notes"`
	AdditionalNotes           []struct {
		Note string `xml:"Note"`
	} `xml:"AdditionalNotes"`
	PreferredDeliveryMethod string     `xml:"PreferredDeliveryMethod"`
	PriceLevelRef           AccountRef `xml:"PriceLevelRef"`
	ExternalGUID            string     `xml:"ExternalGUID"`
	CurrencyRef             AccountRef `xml:"CurrencyRef"`
}

//CustomerAddRs hods the data from the response to CustomerAddRq
type CustomerAddRs struct {
	ResponseStatus
}

//CustomerMsgAddRq is the struct to hold the information to add a customer message reference
type CustomerMsgAddRq struct {
	Name     string `xml:"Name"`
	IsActive string `xml:"IsActive"`
}

//CustomerMsgQueryRq is the struct to hold data for querying for customer messages
type CustomerMsgQueryRq struct {
	XMLName           xml.Name         `xml:"CustomerMsgQueryRq,omitempty"`
	MetaData          string           `xml:"metaData,attr,omitempty"`
	ListID            string           `xml:"ListID,omitempty"`
	FullName          string           `xml:"FullName,omitempty"`
	MaxReturned       string           `xml:"MaxReturned,omitempty"`
	ActiveStatus      string           `xml:"ActiveStatus,omitempty"` //ActiveStatus may have one of the following values: ActiveOnly [DEFAULT], InactiveOnly, All
	FromModifiedDate  string           `xml:"FromModifiedDate,omitempty"`
	ToModifiedDate    string           `xml:"ToModifiedDate,omitempty"`
	NameFilter        *NameFilter      `xml:"NameFilter,omitempty"`
	NameRangeFilter   *NameRangeFilter `xml:"NameRangeFilter,omitempty"`
	IncludeRetElement []string         `xml:"IncludeRetElement,omitempty"`
}

//CustomerMsgQueryRs is the struct to hold data for the response to the customerMsgQuery call
type CustomerMsgQueryRs struct {
	XMLName xml.Name `xml:"CustomerMsgQueryRs"`
	ResponseStatus
	CustomerMsgRet []struct {
		ListID       string `xml:"ListID"`
		TimeCreated  string `xml:"TimeCreated"`
		TimeModified string `xml:"TimeModified"`
		EditSequence string `xml:"EditSequence"`
		Name         string `xml:"Name"`
		IsActive     string `xml:"IsActive"`
	} `xml:"CustomerMsgRet"`
}

//QBShipToAddress holds a quickbooks customer's shipTo addresses
type QBShipToAddress struct {
	Address
	Name          string `xml:"Name"`
	DefaultShipTo string `xml:"DefaultShipTo"`
}

//AdditionalContactRef is a unique ref struct to hold contact information
type AdditionalContactRef struct {
	ContactName  string `xml:"ContactName"`
	ContactValue string `xml:"ContactValue"`
}

//Contacts is a struct to hold contact information when adding a customer
type Contacts struct {
	Salutation           string                 `xml:"Salutation"`
	FirstName            string                 `xml:"FirstName"`
	MiddleName           string                 `xml:"MiddleName"`
	LastName             string                 `xml:"LastName"`
	JobTitle             string                 `xml:"JobTitle"`
	AdditionalContactRef []AdditionalContactRef `xml:"AdditionalContactRef"`
}

//CreditCardInfo is to hold the credit card info whwen adding a customer
type CreditCardInfo struct {
	CreditCardNumber     string `xml:"CreditCardNumber"`
	ExpirationMonth      string `xml:"ExpirationMonth"`
	ExpirationYear       string `xml:"ExpirationYear"`
	NameOnCard           string `xml:"NameOnCard"`
	CreditCardAddress    string `xml:"CreditCardAddress"`
	CreditCardPostalCode string `xml:"CreditCardPostalCode"`
}

//DataExtAddRq is the sruct to add custom data fields to a quick books object
type DataExtAddRq struct {
	OwnerID     string `xml:"OwnerID"`
	DataExtName string `xml:"DataExtName"`
	//<!-- BEGIN OR -->
	ListDataExtType string     `xml:"ListDataExtType"` //may have one of the following values: Account, Customer, Employee, Item, OtherName, Vendor
	ListObjRef      AccountRef `xml:"ListObjRef"`
	//<!-- OR -->
	TxnDataExtType string `xml:"TxnDataExtType"` //may have one of the following values: ARRefundCreditCard, Bill, BillPaymentCheck, BillPaymentCreditCard, BuildAssembly, Charge, Check, CreditCardCharge, CreditCardCredit, CreditMemo, Deposit, Estimate, InventoryAdjustment, Invoice, ItemReceipt, JournalEntry, PurchaseOrder, ReceivePayment, SalesOrder, SalesReceipt, SalesTaxPaymentCheck, VendorCredit -->
	TxnID          string `xml:"TxnID"`
	UseMacro       string `xml:"UserMacro,attr"`
	TxnLineID      string `xml:"TxnLineID"`
	//<!-- OR -->
	OtherDataExtType string `xml:"OtherDataExtType"` //may have one of the following values: Company -->
	//<!-- END OR -->
	DataExtValue string `xml:"DataExtValue"`
}

//CustomerModRs is the struct to hold data from the customer mod response
type CustomerModRs struct {
	XMLName xml.Name `xml:"CustomerModRs"`
	ResponseStatus
	CustomerRet struct {
		ListID                    string               `xml:"ListID"`
		TimeCreated               string               `xml:"TimeCreated"`
		TimeModified              string               `xml:"TimeModified"`
		EditSequence              string               `xml:"EditSequence"`
		Name                      string               `xml:"Name"`
		FullName                  string               `xml:"FullName"`
		IsActive                  string               `xml:"IsActive"`
		ClassRef                  AccountRef           `xml:"ClassRef"`
		ParentRef                 AccountRef           `xml:"ParentRef"`
		Sublevel                  string               `xml:"Sublevel"`
		CompanyName               string               `xml:"CompanyName"`
		Salutation                string               `xml:"Salutation"`
		FirstName                 string               `xml:"FirstName"`
		MiddleName                string               `xml:"MiddleName"`
		LastName                  string               `xml:"LastName"`
		JobTitle                  string               `xml:"JobTitle"`
		BillAddress               Address              `xml:"BillAddress"`
		BillAddressBlock          Address              `xml:"BillAddressBlock"`
		ShipAddress               Address              `xml:"ShipAddress"`
		ShipAddressBlock          Address              `xml:"ShipAddressBlock"`
		ShipToAddress             []Address            `xml:"ShipToAddress"`
		Phone                     string               `xml:"Phone"`
		AltPhone                  string               `xml:"AltPhone"`
		Fax                       string               `xml:"Fax"`
		Email                     string               `xml:"Email"`
		Cc                        string               `xml:"Cc"`
		Contact                   string               `xml:"Contact"`
		AltContact                string               `xml:"AltContact"`
		AdditionalContactRef      AdditionalContactRef `xml:"AdditionalContactRef"`
		ContactsRet               ContactsRet          `xml:"ContactsRet"`
		CustomerTypeRef           AccountRef           `xml:"CustomerTypeRef"`
		TermsRef                  AccountRef           `xml:"TermsRef"`
		SalesRepRef               AccountRef           `xml:"SalesRepRef"`
		Balance                   string               `xml:"Balance"`
		TotalBalance              string               `xml:"TotalBalance"`
		SalesTaxCodeRef           AccountRef           `xml:"SalesTaxCodeRef"`
		ItemSalesTaxRef           AccountRef           `xml:"ItemSalesTaxRef"`
		ResaleNumber              string               `xml:"ResaleNumber"`
		AccountNumber             string               `xml:"AccountNumber"`
		CreditLimit               string               `xml:"CreditLimit"`
		PreferredPaymentMethodRef AccountRef           `xml:"PreferredPaymentMethodRef"`
		CreditCardInfo            CreditCardInfo       `xml:"CreditCardInfo"`
		JobStatus                 string               `xml:"JobStatus"`
		JobStartDate              string               `xml:"JobStartDate"`
		JobProjectedEndDate       string               `xml:"JobProjectedEndDate"`
		JobEndDate                string               `xml:"JobEndDate"`
		JobDesc                   string               `xml:"JobDesc"`
		JobTypeRef                AccountRef           `xml:"JobTypeRef"`
		Notes                     string               `xml:"Notes"`
		AdditionalNotesRet        AdditionalNotes      `xml:"AdditionalNotesRet"`
		PreferredDeliveryMethod   string               `xml:"PreferredDeliveryMethod"`
		PriceLevelRef             AccountRef           `xml:"PriceLevelRef"`
		ExternalGUID              string               `xml:"ExternalGUID"`
		CurrencyRef               AccountRef           `xml:"CurrencyRef"`
		DataExtRet                DataExtRet           `xml:"DataExtRet"`
	} `xml:"CustomerRet"`
	ErrorRecovery *ErrorRecovery `xml:"ErrorRecovery"`
}

//CustomerModRq isis the struct to hold data for the customerModRq call
type CustomerModRq struct {
	XMLName     xml.Name `xml:"CustomerModRq"`
	CustomerMod struct {
		ListID                    string                  `xml:"ListID,omitempty"`
		EditSequence              string                  `xml:"EditSequence,omitempty"`
		Name                      string                  `xml:"Name,omitempty"`
		IsActive                  string                  `xml:"IsActive,omitempty"`
		ClassRef                  *AccountRef             `xml:"ClassRef,omitempty"`
		ParentRef                 *AccountRef             `xml:"ParentRef,omitempty"`
		CompanyName               string                  `xml:"CompanyName,omitempty"`
		Salutation                string                  `xml:"Salutation,omitempty"`
		FirstName                 string                  `xml:"FirstName,omitempty"`
		MiddleName                string                  `xml:"MiddleName,omitempty"`
		LastName                  string                  `xml:"LastName,omitempty"`
		JobTitle                  string                  `xml:"JobTitle,omitempty"`
		BillAddress               *Address                `xml:"BillAddress,omitempty"`
		ShipAddress               *Address                `xml:"ShipAddress,omitempty"`
		ShipToAddress             *[]Address              `xml:"ShipToAddress,omitempty"`
		Phone                     string                  `xml:"Phone,omitempty"`
		AltPhone                  string                  `xml:"AltPhone,omitempty"`
		Fax                       string                  `xml:"Fax,omitempty"`
		Email                     string                  `xml:"Email,omitempty"`
		Cc                        string                  `xml:"Cc,omitempty"`
		Contact                   string                  `xml:"Contact,omitempty"`
		AltContact                string                  `xml:"AltContact,omitempty"`
		AdditionalContactRef      *[]AdditionalContactRef `xml:"AdditionalContactRef,omitempty"`
		ContactsMod               *[]ContactsRet          `xml:"ContactsMod,omitempty"`
		CustomerTypeRef           *AccountRef             `xml:"CustomerTypeRef,omitempty"`
		TermsRef                  *AccountRef             `xml:"TermsRef,omitempty"`
		SalesRepRef               *AccountRef             `xml:"SalesRepRef,omitempty"`
		SalesTaxCodeRef           *AccountRef             `xml:"SalesTaxCodeRef,omitempty"`
		ItemSalesTaxRef           *AccountRef             `xml:"ItemSalesTaxRef,omitempty"`
		ResaleNumber              string                  `xml:"ResaleNumber,omitempty"`
		AccountNumber             string                  `xml:"AccountNumber,omitempty"`
		CreditLimit               string                  `xml:"CreditLimit,omitempty"`
		PreferredPaymentMethodRef *AccountRef             `xml:"PreferredPaymentMethodRef,omitempty"`
		CreditCardInfo            *CreditCardInfo         `xml:"CreditCardInfo,omitempty"`
		JobStatus                 string                  `xml:"JobStatus,omitempty"` //may have one of the following values: Awarded, Closed, InProgress, None [DEFAULT], NotAwarded, Pending
		JobStartDate              string                  `xml:"JobStartDate,omitempty"`
		JobProjectedEndDate       string                  `xml:"JobProjectedEndDate,omitempty"`
		JobEndDate                string                  `xml:"JobEndDate,omitempty"`
		JobDesc                   string                  `xml:"JobDesc,omitempty"`
		JobTypeRef                *AccountRef             `xml:"JobTypeRef,omitempty"`
		Notes                     string                  `xml:"Notes,omitempty"`
		AdditionalNotesMod        *AdditionalNotes        `xml:"AdditionalNotesMod,omitempty"`
		PreferredDeliveryMethod   string                  `xml:"PreferredDeliveryMethod,omitempty"` //may have one of the following values: None [Default], Email, Fax
		PriceLevelRef             *AccountRef             `xml:"PriceLevelRef,omitempty"`
		CurrencyRef               *AccountRef             `xml:"CurrencyRef,omitempty"`
	} `xml:"CustomerMod,omitempty"`
	IncludeRetElement []string `xml:"IncludeRetElement,omitempty"`
}

//AdditionalNotes is the struct to hold data for the additional notes node
type AdditionalNotes struct {
	NoteID string `xml:"NoteID,omitempty"`
	Note   string `xml:"Note,omitempty"`
}

//CustomerQueryRq is the struct to hold data for the customerQueryRq call
type CustomerQueryRq struct {
	XMLName            xml.Name            `xml:"CustomerQueryRq"`
	MetaData           string              `xml:"metaData,attr,omitempty"`
	Iterator           string              `xml:"iterator,attr,omitempty"`
	IteratorID         string              `xml:"iteratorID,attr,omitempty"`
	ListID             []string            `xml:"ListID,omitempty"`
	FullName           []string            `xml:"FullName,omitempty"`
	MaxReturned        int                 `xml:"MaxReturned,omitempty"`
	ActiveStatus       string              `xml:"ActiveStatus,omitempty"` //ActiveStatus may have one of the following values: ActiveOnly [DEFAULT], InactiveOnly, All
	FromModifiedDate   string              `xml:"FromModifiedDate,omitempty"`
	ToModifiedDate     string              `xml:"ToModifiedDate,omitempty"`
	NameFilter         *NameFilter         `xml:"NameFilter,omitempty"`
	NameRangeFilter    *NameRangeFilter    `xml:"NameRangeFilter,omitempty"`
	TotalBalanceFilter *TotalBalanceFilter `xml:"TotalBalanceFilter,omitempty"`
	CurrencyFilter     *CurrencyFilter     `xml:"CurrencyFilter,omitempty"`
	ClassFilter        *ClassFilter        `xml:"ClassFilte,omitempty"`
	IncludeRetElement  []string            `xml:"IncludeRetElement,omitempty"`
	OwnerID            []string            `xml:"OwnerID,omitempty"`
}

//TotalBalanceFilter is the struct to hold data for customer queries
type TotalBalanceFilter struct {
	Operator string `xml:"Operator,omitempty"` //Operator may have one of the following values: LessThan, LessThanEqual, Equal, GreaterThan, GreaterThanEqual
	Amount   string `xml:"Amount,omitempty"`
}

//CurrencyFilter is the struct to hold data for customer queries
type CurrencyFilter struct {
	ListID   []string `xml:"ListID,omitempty"`
	FullName []string `xml:"FullName,omitempty"`
}

//ClassFilter is the struct to hold data for customer quries
type ClassFilter struct {
	ListID               []string `xml:"ListID,omitempty"`
	FullName             []string `xml:"FullName,omitempty"`
	ListIDWithChildren   string   `xml:"ListIDWithChildren,omitempty"`
	FullNameWithChildren string   `xml:"FullNameWithChildren,omitempty"`
}

//CustomerQueryRs aaa
type CustomerQueryRs struct {
	ResponseStatus
	RetCount               int           `xml:"retCount,attr"`
	IteratorRemainingCount int           `xml:"iteratorRemainingCount,attr"`
	IteratorID             string        `xml:"iteratorID,attr"`
	CustomerRet            []CustomerRet `xml:"CustomerRet"`
}

//CustomerRet is the struct to hold data for the customer query response
type CustomerRet struct {
	ListID                    string                 `xml:"ListID"`
	TimeCreated               string                 `xml:"TimeCreated"`
	TimeModified              string                 `xml:"TimeModified"`
	EditSequence              string                 `xml:"EditSequence"`
	Name                      string                 `xml:"Name"`
	FullName                  string                 `xml:"FullName"`
	IsActive                  string                 `xml:"IsActive"`
	ClassRef                  AccountRef             `xml:"ClassRef"`
	ParentRef                 AccountRef             `xml:"ParentRef"`
	Sublevel                  string                 `xml:"Sublevel"`
	CompanyName               string                 `xml:"CompanyName"`
	Salutation                string                 `xml:"Salutation"`
	FirstName                 string                 `xml:"FirstName"`
	MiddleName                string                 `xml:"MiddleName"`
	LastName                  string                 `xml:"LastName"`
	JobTitle                  string                 `xml:"JobTitle"`
	BillAddress               Address                `xml:"BillAddress"`
	BillAddressBlock          Address                `xml:"BillAddressBlock"`
	ShipAddress               Address                `xml:"ShipAddress"`
	ShipAddressBlock          Address                `xml:"ShipAddressBlock"`
	ShipToAddress             []QBShipToAddress      `xml:"ShipToAddress"`
	Phone                     string                 `xml:"Phone"`
	AltPhone                  string                 `xml:"AltPhone"`
	Fax                       string                 `xml:"Fax"`
	Email                     string                 `xml:"Email"`
	CC                        string                 `xml:"Cc"`
	Contact                   string                 `xml:"Contact"`
	AltContact                string                 `xml:"AltContact"`
	AdditionalContactRef      []AdditionalContactRef `xml:"AdditionalContactRef"`
	ContactsRet               []ContactsRet          `xml:"ContactsRet"`
	CustomerTypeRef           AccountRef             `xml:"CustomerTypeRef"`
	TermsRef                  AccountRef             `xml:"TermsRef"`
	SalesRepRef               AccountRef             `xml:"SalesRepRef"`
	Balance                   string                 `xml:"Balance"`
	TotalBalance              string                 `xml:"TotalBalance"`
	SalesTaxCodeRef           AccountRef             `xml:"SalesTaxCodeRef"`
	ItemSalesTaxRef           AccountRef             `xml:"ItemSalesTaxRef"`
	ResaleNumber              string                 `xml:"ResaleNumber"`
	AccountNumber             string                 `xml:"AccountNumber"`
	CreditLimit               string                 `xml:"CreditLimit"`
	PreferredPaymentMethodRef AccountRef             `xml:"PreferredPaymentMethodRef"`
	CreditCardInfo            CreditCardInfo         `xml:"CreditCardInfo"`
	JobStatus                 string                 `xml:"JobStatus"`
	JobStartDate              string                 `xml:"JobStartDate"`
	JobProjectedEndDate       string                 `xml:"JobProjectedEndDate"`
	JobEndDate                string                 `xml:"JobEndDate"`
	JobDesc                   string                 `xml:"JobDesc"`
	JobTypeRef                AccountRef             `xml:"JobTypeRef"`
	Notes                     string                 `xml:"Notes"`
	AdditionalNotesRet        []AdditionalNotes      `xml:"AdditionalNotesRet"`
	PreferredDeliveryMethod   string                 `xml:"PreferredDeliveryMethod"`
	PriceLevelRef             AccountRef             `xml:"PriceLevelRef"`
	CurrencyRef               AccountRef             `xml:"CurrencyRef"`
}

//ContactsRet is the struct to hold data for customer query responses
type ContactsRet struct {
	ListID               string                  `xml:"ListID,omitempty"`
	TimeCreated          string                  `xml:"TimeCreated,omitempty"`
	TimeModified         string                  `xml:"TimeModified,omitempty"`
	EditSequence         string                  `xml:"EditSequence,omitempty"`
	Contact              string                  `xml:"Contact,omitempty"`
	Salutation           string                  `xml:"Salutation,omitempty"`
	FirstName            string                  `xml:"FirstName,omitempty"`
	MiddleName           string                  `xml:"MiddleName,omitempty"`
	LastName             string                  `xml:"LastName,omitempty"`
	JobTitle             string                  `xml:"JobTitle,omitempty"`
	AdditionalContactRef *[]AdditionalContactRef `xml:"AdditionalContactRef,omitempty"`
}

/*
        <CreditCardInfo>
            <!-- optional -->
            <CreditCardNumber>STRTYPE</CreditCardNumber>
            <!-- optional -->
            <ExpirationMonth>INTTYPE</ExpirationMonth>
            <!-- optional -->
            <ExpirationYear>INTTYPE</ExpirationYear>
            <!-- optional -->
            <NameOnCard>STRTYPE</NameOnCard>
            <!-- optional -->
            <CreditCardAddress>STRTYPE</CreditCardAddress>
            <!-- optional -->
            <CreditCardPostalCode>STRTYPE</CreditCardPostalCode>
            <!-- optional -->
        </CreditCardInfo>
        <!-- JobStatus may have one of the following values: Awarded, Closed, InProgress, None [DEFAULT], NotAwarded, Pending -->
        <JobStatus>ENUMTYPE</JobStatus>
        <!-- optional -->
        <JobStartDate>DATETYPE</JobStartDate>
        <!-- optional -->
        <JobProjectedEndDate>DATETYPE</JobProjectedEndDate>
        <!-- optional -->
        <JobEndDate>DATETYPE</JobEndDate>
        <!-- optional -->
        <JobDesc>STRTYPE</JobDesc>
        <!-- optional -->
        <JobTypeRef>
            <!-- optional -->
            <ListID>IDTYPE</ListID>
            <!-- optional -->
            <FullName>STRTYPE</FullName>
            <!-- optional -->
        </JobTypeRef>
        <Notes>STRTYPE</Notes>
        <!-- optional -->
        <AdditionalNotesRet>
            <!-- optional, may repeat -->
            <NoteID>INTTYPE</NoteID>
            <!-- required -->
            <Date>DATETYPE</Date>
            <!-- required -->
            <Note>STRTYPE</Note>
            <!-- required -->
        </AdditionalNotesRet>
        <!-- PreferredDeliveryMethod may have one of the following values: None [Default], Email, Fax -->
        <PreferredDeliveryMethod>ENUMTYPE</PreferredDeliveryMethod>
        <!-- optional -->
        <PriceLevelRef>
            <!-- optional -->
            <ListID>IDTYPE</ListID>
            <!-- optional -->
            <FullName>STRTYPE</FullName>
            <!-- optional -->
        </PriceLevelRef>
        <ExternalGUID>GUIDTYPE</ExternalGUID>
        <!-- optional -->
        <CurrencyRef>
            <!-- optional -->
            <ListID>IDTYPE</ListID>
            <!-- optional -->
            <FullName>STRTYPE</FullName>
            <!-- optional -->
        </CurrencyRef>
        <DataExtRet>
            <!-- optional, may repeat -->
            <OwnerID>GUIDTYPE</OwnerID>
            <!-- optional -->
            <DataExtName>STRTYPE</DataExtName>
            <!-- required -->
            <!-- DataExtType may have one of the following values: AMTTYPE, DATETIMETYPE, INTTYPE, PERCENTTYPE, PRICETYPE, QUANTYPE, STR1024TYPE, STR255TYPE -->
            <DataExtType>ENUMTYPE</DataExtType>
            <!-- required -->
            <DataExtValue>STRTYPE</DataExtValue>
            <!-- required -->
        </DataExtRet>
    </CustomerRet>
</CustomerQueryRs>
*/
