{{define "qbItemInventoryAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
   <QBXMLMsgsRq onError="stopOnError">
      <ItemInventoryAddRq>
         <ItemInventoryAdd>
            <Name >{{if .Name}}{{ .Name}}{{end}}</Name>{{if .BarCode.BarCodeValue}}
            <BarCode>{{if .BarCode.BarCodeValue}}
               <BarCodeValue >{{ .BarCode.BarCodeValue}}</BarCodeValue>{{end}}{{if .BarCode.AssignEvenIfUsed}}
               <AssignEvenIfUsed >{{ .BarCode.AssignEvenIfUsed}}</AssignEvenIfUsed>{{end}}{{if .BarCode.AllowOverride}}
               <AllowOverride >{{ .BarCode.AllowOverride}}</AllowOverride>{{end}}
            </BarCode>{{end}}{{if .IsActive}}
            <IsActive >{{ .IsActive}}</IsActive>{{end}}{{if or .ClassRef.ListID .ClassRef.FullName}}
            <ClassRef>{{if .ClassRef.ListID}}
               <ListID >{{ .ClassRef.ListID}}</ListID>{{end}}{{if .ClassRef.FullName}}
               <FullName >{{ .ClassRef.FullName}}</FullName>{{end}}
            </ClassRef>{{end}}{{if or .ParentRef.ListID .ParentRef.FullName}}
            <ParentRef>{{if .ParentRef.ListID}}
               <ListID >{{ .ParentRef.ListID}}</ListID>{{end}}{{if .ParentRef.FullName}}
               <FullName >{{ .ParentRef.FullName}}</FullName>{{end}}
            </ParentRef>{{end}}{{if .ManufacturerPartNumber}}
            <ManufacturerPartNumber >{{ .ManufacturerPartNumber}}</ManufacturerPartNumber>{{end}}{{if or .UnitOfMeasureSetRef.ListID .UnitOfMeasureSetRef.FullName}}
            <UnitOfMeasureSetRef>{{if .UnitOfMeasureSetRef.ListID}}
               <ListID >{{ .UnitOfMeasureSetRef.ListID}}</ListID>{{end}}{{if .UnitOfMeasureSetRef.FullName}}
               <FullName >{{ .UnitOfMeasureSetRef.FullName}}</FullName>{{end}}
            </UnitOfMeasureSetRef>{{end}}{{if or .SalesTaxCodeRef.ListID .SalesTaxCodeRef.FullName}}
            <SalesTaxCodeRef>{{if .SalesTaxCodeRef.ListID}}
               <ListID >{{ .SalesTaxCodeRef.ListID}}</ListID>{{end}}{{if .SalesTaxCodeRef.FullName}}
               <FullName >{{ .SalesTaxCodeRef.FullName}}</FullName>{{end}}
            </SalesTaxCodeRef>{{end}}{{if .SalesDesc}}
            <SalesDesc >{{ .SalesDesc}}</SalesDesc>{{end}}{{if .SalesPrice}}
            <SalesPrice >{{.SalesPrice}}</SalesPrice>{{end}}{{if or .IncomeAccountRef.ListID .IncomeAccountRef.FullName}}
            <IncomeAccountRef>{{if .IncomeAccountRef.ListID}}
               <ListID >{{ .IncomeAccountRef.ListID}}</ListID>{{end}}{{if .IncomeAccountRef.FullName}}
               <FullName >{{ .IncomeAccountRef.FullName}}</FullName>{{end}}
            </IncomeAccountRef>{{end}}{{if .PurchaseDesc}}
            <PurchaseDesc >{{ .PurchaseDesc}}</PurchaseDesc>{{end}}{{if .PurchaseCost}}
            <PurchaseCost >{{ .PurchaseCost}}</PurchaseCost>{{end}}{{if or .COGSAccountRef.ListID .COGSAccountRef.FullName}}
            <COGSAccountRef>{{if .COGSAccountRef.ListID}}
               <ListID >{{ .COGSAccountRef.ListID}}</ListID>{{end}}{{if .COGSAccountRef.FullName}}
               <FullName >{{ .COGSAccountRef.FullName}}</FullName>{{end}}
            </COGSAccountRef>{{end}}{{if or .PrefVendorRef.ListID .PrefVendorRef.FullName}}
            <PrefVendorRef>{{if .PrefVendorRef.ListID}}
               <ListID >{{ .PrefVendorRef.ListID}}</ListID>{{end}}{{if .PrefVendorRef.FullName}}
               <FullName >{{ .PrefVendorRef.FullName}}</FullName>{{end}}
            </PrefVendorRef>{{end}}{{if or .AssetAccountRef.ListID .AssetAccountRef.FullName}}
            <AssetAccountRef>{{if .AssetAccountRef.ListID}}
               <ListID >{{ .AssetAccountRef.ListID}}</ListID>{{end}}{{if .AssetAccountRef.FullName}}
               <FullName >{{ .AssetAccountRef.FullName}}</FullName>{{end}}
            </AssetAccountRef>{{end}}{{if .ReorderPoint}}
            <ReorderPoint >{{ .ReorderPoint}}</ReorderPoint>{{end}}{{if .Max}}
            <Max >{{ .Max}}</Max>{{end}}{{if .QuantityOnHand}}
            <QuantityOnHand >{{ .QuantityOnHand}}</QuantityOnHand>{{end}}{{if .TotalValue}}
            <TotalValue >{{ .TotalValue}}</TotalValue>{{end}}{{if .InventoryDate}}
            <InventoryDate >{{ .InventoryDate}}</InventoryDate>{{end}}{{if .ExternalGUID}}
            <ExternalGUID >{{ .ExternalGUID}}</ExternalGUID>{{end}}
         </ItemInventoryAdd>
      </ItemInventoryAddRq>
   </QBXMLMsgsRq>
</QBXML>{{end}}
