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
