{{define "qbItemServiceAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
   <QBXMLMsgsRq onError="stopOnError">
      <ItemServiceAddRq>
         <ItemServiceAdd>
            <Name >{{ if .Name}}{{ .Name}}{{end}}</Name>{{ if .BarCode.BarCodeValue}}
            <BarCode>{{ if .BarCode.BarCodeValue}}
               <BarCodeValue >{{ .BarCode.BarCodeValue}}</BarCodeValue>{{end}}{{ if .BarCode.AssignEvenIfUsed}}
               <AssignEvenIfUsed >{{ .BarCode.AssignEvenIfUsed}}</AssignEvenIfUsed>{{end}}{{ if .BarCode.AllowOverride}}
               <AllowOverride >{{ .BarCode.AllowOverride}}</AllowOverride>{{end}}
            </BarCode>{{end}}{{ if .IsActive}}
            <IsActive >{{ .IsActive}}</IsActive>{{end}}{{ if or .ClassRef.ListID .ClassRef.FullName}}
            <ClassRef>{{ if .ClassRef.ListID}}
               <ListID >{{ .ClassRef.ListID}}</ListID>{{end}}{{ if .ClassRef.FullName}}
               <FullName >{{ .ClassRef.FullName}}</FullName>{{end}}
            </ClassRef>{{end}}{{ if or .ParentRef.ListID .ParentRef.FullName}}
            <ParentRef>{{ if .ParentRef.ListID}}
               <ListID >{{ .ParentRef.ListID}}</ListID>{{end}}{{ if .ParentRef.FullName}}
               <FullName >{{ .ParentRef.FullName}}</FullName>{{end}}
            </ParentRef>{{end}}{{ if or .UnitOfMeasureSetRef.ListID .UnitOfMeasureSetRef.FullName}}
            <UnitOfMeasureSetRef>{{ if .UnitOfMeasureSetRef.ListID}}
               <ListID >{{ .UnitOfMeasureSetRef.ListID}}</ListID>{{end}}{{ if .UnitOfMeasureSetRef.FullName}}
               <FullName >{{ .UnitOfMeasureSetRef.FullName}}</FullName>{{end}}
            </UnitOfMeasureSetRef>{{end}}{{ if or .SalesTaxCodeRef.ListID .SalesTaxCodeRef.FullName}}
            <SalesTaxCodeRef>{{ if .SalesTaxCodeRef.ListID}}
               <ListID >{{ .SalesTaxCodeRef.ListID}}</ListID>{{end}}{{ if .SalesTaxCodeRef.FullName}}
               <FullName >{{ .SalesTaxCodeRef.FullName}}</FullName>{{end}}
            </SalesTaxCodeRef>{{end}}{{ if or .SalesOrPurchase.Price .SalesOrPurchase.Desc}}
            <SalesOrPurchase>{{ if .SalesOrPurchase.Desc}}
                <Desc>{{ .SalesOrPurchase.Desc}}</Desc>{{end}}{{ if .SalesOrPurchase.Price}}
                <Price>{{.SalesOrPurchase.Price}}</Price>{{else if .SalesOrPurchase.PricePercent}}
                <PricePercent>{{.SalesOrPurchase.PricePercent}}</PricePercent>{{end}}{{ if or .SalesOrPurchase.AccountRef.FullName .SalesOrPurchase.AccountRef.ListID }}
                <AccountRef>{{ if .SalesOrPurchase.AccountRef.ListID}}
                    <ListID>{{.SalesOrPurchase.AccountRef.ListID}}</ListID>{{end}}{{ if .SalesOrPurchase.AccountRef.FullName}}
                    <FullName>{{ .SalesOrPurchase.AccountRef.FullName}}</FullName>{{end}}
                </AccountRef>{{end}}
            </SalesOrPurchase>{{else if or .SalesAndPurchase.SalesPrice .SalesAndPurchase.SalesDesc}}
            <SalesAndPurchase>{{ if .SalesAndPurchase.SalesDesc}}
                <SalesDesc>{{ .SalesAndPurchase.SalesDesc}}</SalesDesc>{{end}}{{ if .SalesAndPurchase.SalesPrice}}
                <SalesPrice>{{.SalesAndPurchase.SalesPrice}}</SalesPrice>{{end}}{{ if or .SalesAndPurchase.IncomeAccountRef.FullName .SalesAndPurchase.IncomeAccountRef.ListID}} 
                <IncomeAccountRef>{{ if .SalesAndPurchase.IncomeAccountRef.ListID}}
                    <ListID>{{.SalesAndPurchase.IncomeAccountRef.ListID}}</ListID>{{end}}{{if .SalesAndPurchase.IncomeAccountRef.FullName}}
                    <FullName>{{.SalesAndPurchase.IncomeAccountRef.FullName}}</FullName>{{end}}
                </IncomeAccountRef>{{end}}{{ if .SalesAndPurchase.PurchaseDesc}}
                <PurchaseDesc>{{ .SalesAndPurchase.PurchaseDesc}}</PurchaseDesc>{{end}}{{ if .SalesAndPurchase.PurchaseCost}}
                <PurchaseCost>{{ .SalesAndPurchase.PurchaseCost}}</PurchaseCost>{{end}}{{ if or .SalesAndPurchase.ExpenseAccountRef.FullName .SalesAndPurchase.ExpenseAccountRef.ListID}}
                <ExpenseAccountRef>{{ if .SalesAndPurchase.ExpenseAccountRef.ListID}}
                    <ListID>{{ .SalesAndPurchase.ExpenseAccountRef.ListID}}</ListID>{{end}}{{ if .SalesAndPurchase.ExpenseAccountRef.FullName}}
                    <FullName>{{.SalesAndPurchase.ExpenseAccountRef.FullName}}</FullName>{{end}}
                </ExpenseAccountRef>{{end}}{{ if or .SalesAndPurchase.PrefVendorRef.FullName .SalesAndPurchase.PrefVendorRef.ListID}}
                <PrefVendorRef>{{ if .SalesAndPurchase.PrefVendorRef.ListID}}
                    <ListID>{{ .SalesAndPurchase.PrefVendorRef.ListID}}</ListID>{{end}}{{ if .SalesAndPurchase.PrefVendorRef.FullName}}
                    <FullName>{{ .SalesAndPurchase.PrefVendorRef.FullName}}</FullName>{{end}}
                </PrefVendorRef>{{end}}
            </SalesAndPurchase>{{end}}{{ if .ExternalGUID}}
            <ExternalGUID >{{ .ExternalGUID}}</ExternalGUID>{{end}}
         </ItemServiceAdd>
      </ItemServiceAddRq>
   </QBXMLMsgsRq>
</QBXML>{{end}}
