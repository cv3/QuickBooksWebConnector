	Most fields are set by mapping the CV3 field to the corresponding Quickbooks field.  Though several fields have hard coded values.
	
	//////////////////////////
	"CustomerRef.FullName":"",
	//If the billing name is not paypal, use it as firstName lastName
	qbOrderAdd.CustomerRef.FullName = CheckPath("billing.firstName", o) + " " + CheckPath("billing.lastName", o)
	//billing firstName is paypal, so just add paypal as a CustomerRef is required for a SalesOrderAdd
	qbOrderAdd.CustomerRef.FullName = CheckPath("billing.firstName", o)
	///////////////////////////

	Addresses
	///////////////////////////
	if title exists, add it in the front of the first address line followed by a space
	If firstName exists, add it to the first address line followed by a space
	If lastNamme exists, add it to the first address line followed by a space
	If company exists, add it to the end of the first address line
	Make sure it does not exceed the Quickbooks 41 char limit
	////////////////////////////

	///////////////////////////
	"TermsRef.FullName":"",
	"TermsRef.ListID":"",
	if TermsRef.FullName's mapped value = "creditcard", set it to "Credit Card"
	if TermsRef.FullName's mapped value = "paypal", set it to "PayPal"
	if TermsRef.FullName's mapped value = "ccpaypal", set it to "CCPaypal"
	else set it to its mapped value
	////////////////////////////
	
	Shipping
	////////////////////////////
	ItemRef.FullName and Quantity fields have hard coded values in the shippingOrderMapping.json file
	////////////////////////////

	QuickBooks TxnID
	////////////////////////////
	This is an internal quickbooks value, it get set by the cv3 order number appended with the index of the cv3 shipTo object that corresponds to the quickbook's order.
	////////////////////////////

	DataExt for adding custom fields
	////////////////////////////
	All hard coded in the binary to match Mac's Tie Downs custom fields
	Source = eCommerce
	Package = No
	/////////////////////////////

	Customer Adds
	/////////////////////////////
	New Customers have a few hard coded values set int customerAddMapping.json
	CustomerTypeRef.FullName
	SalesRepRef.FullName
	PreferredPaymentMethodRef.FullName
	TermsRef.FullName
	PriceLevelRef.FullName
	//////////////////////////////



	

Quickbooks and CV3 fields are matched up using various fieldMapper.json files.  A single field mapping will look like the following:

```"SKU":"ListID",```

Here we can see that the SKU field is being set to ListID, the assignment in the code looks as follows:
```itemTemp.Sku = CheckPath(fieldMap["SKU"], qbItem)```

This will look up SKU in the map and see it is mapped to ListID, then assigns the value located in ListID to itemTemp.Sku

Some of the data structures have multiple levels, to access sub levels we will use dot notation.
```"Retail.Price.StandardPrice":"SalesOrPurchase.Price",```

The above would map from CV3:
```
<Retail active="true">
    <Price >
        <StandardPrice>13.9500</StandardPrice>
    </Price>
</Retail>
```
to QuickBooks:
```
<SalesOrPurchase>
    <Price>13.95</Price>
</SalesOrPurchase>
```
