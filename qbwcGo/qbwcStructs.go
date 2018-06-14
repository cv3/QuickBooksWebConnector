package qbwcGo

import (
	"encoding/xml"

	"github.com/amazingfly/cv3go"
)

//WorkCTX is a struct that will hold both the work to be done, and the data used to create it
type WorkCTX struct {
	Work        string          //holds the excaped qbxml
	Data        interface{}     //holds the struct that created the qbxml
	CV3Products []cv3go.Product //holds the cv3 products used to make the qbxml
	Type        string          //type of qbxml request
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
	ListID   string `xml:"ListID"`
	FullName string `xml:"FullName"`
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
	Addr1      string `xml:"Addr1"`
	Addr2      string `xml:"Addr2,omitempty"`
	Addr3      string `xml:"Addr3,omitempty"`
	Addr4      string `xml:"Addr4,omitempty"`
	Addr5      string `xml:"Addr5,omitempty"`
	City       string `xml:"City"`
	State      string `xml:"State"`
	PostalCode string `xml:"PostalCode"`
	Country    string `xml:"Country"`
	Note       string `xml:"Note,omitempty"`
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

//SalesReceiptAdd is the struct for added sales receipts from cv3 to qb
type SalesReceiptAdd struct {
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
	SalesReceiptLineAdd []SalesReceiptLineAdd `xml:"SalesReceiptLineAdd,ommitempty"`
	//1<!-- OR -->
	SalesReceiptLineGroupAdd []SalesReceiptLineGroupAdd `xml:"SalesReceiptLineGroupAdd,ommitempty"`
	//1<!-- END OR -->
	IncludeRetElement string `xml:"IncludeRetElement,ommitempty"`
}

//SalesOrderAdd is the struct to hold the variables for the qbxml call
type SalesOrderAdd struct {
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
	ItemRef                    AccountRef `xml:"ItemRef"`
	Desc                       string     `xml:"Desc"`
	Quantity                   string     `xml:"Quantity"`
	UnitOfMeasure              string     `xml:"UnitOfMeasure"`
	Rate                       string     `xml:"Rate"`
	RatePercent                string     `xml:"RatePercent"`
	PriceLevelRef              AccountRef `xml:"PriceLevelRef"`
	ClassRef                   AccountRef `xml:"ClassRef"`
	Amount                     string     `xml:"Amount"`
	OptionForPriceRuleConflict string     `xml:"OptionForPriceRuleConflict"` //OptionForPriceRuleConflict may have one of the following values: Zero, BasePrice -->
	InventorySiteRef           AccountRef `xml:"InventorySiteRef"`
	InventorySiteLocationRef   AccountRef `xml:"InventorySiteLocationRef"`
	SerialNumber               string     `xml:"SerialNumber"`
	LotNumber                  string     `xml:"LotNumber"`
	SalesTaxCodeRef            AccountRef `xml:"SalesTaxCodeRef"`
	IsManuallyClosed           string     `xml:"IsManuallyClosed"`
	Other1                     string     `xml:"Other1"`
	Other2                     string     `xml:"Other2"`
	DataExt                    []DataExt  `xml:"DataExt"`
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
