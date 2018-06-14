{{define "qbSalesOrderAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
   <QBXMLMsgsRq onError="stopOnError">
      <SalesOrderAddRq>
         <SalesOrderAdd >{{if or .CustomerRef.FullName .CustomerRef.ListID}}
            <CustomerRef>{{if .CustomerRef.ListID}}
               <ListID >{{.CustomerRef.ListID}}</ListID>{{end}}{{if .CustomerRef.FullName}}
               <FullName >{{.CustomerRef.FullName}}</FullName>{{end}}
            </CustomerRef>{{end}}{{if or .ClassRef.ListID .ClassRef.FullName}}
            <ClassRef>{{if .ClassRef.ListID}}
               <ListID >{{.ClassRef.ListID}}</ListID>{{end}}{{if .ClassRef.FullName}}
               <FullName >{{.ClassRef.FullName}}</FullName>{{end}}
            </ClassRef>{{end}}{{if or .TemplateRef.ListID .TemplateRef.FullName}}
            <TemplateRef>{{if .TemplateRef.ListID}}
               <ListID >{{.TemplateRef.ListID}}</ListID>{{end}}{{if .TemplateRef.FullName}}
               <FullName >{{.TemplateRef.FullName}}</FullName>{{end}}
            </TemplateRef>{{end}}{{if .TxnDate}}
            <TxnDate >{{.TxnDate}}</TxnDate>{{end}}{{if .RefNumber}}
            <RefNumber >{{.RefNumber}}</RefNumber>{{end}}{{if .BillAddress}}
            <BillAddress>{{if .BillAddress.Addr1}}
               <Addr1 >{{.BillAddress.Addr1}}</Addr1>{{end}}{{if .BillAddress.Addr2}}
               <Addr2 >{{ .BillAddress.Addr2}}</Addr2>{{end}}{{if .BillAddress.Addr3}}
               <Addr3 >{{ .BillAddress.Addr3}}</Addr3>{{end}}{{if .BillAddress.Addr4}}
               <Addr4 >{{ .BillAddress.Addr4}}</Addr4>{{end}}{{if .BillAddress.Addr5}}
               <Addr5 >{{ .BillAddress.Addr5}}</Addr5>{{end}}{{if .BillAddress.City}}
               <City >{{ .BillAddress.City}}</City>{{end}}{{if .BillAddress.State}}
               <State >{{ .BillAddress.State}}</State>{{end}}{{if .BillAddress.PostalCode}}
               <PostalCode >{{ .BillAddress.PostalCode}}</PostalCode>{{end}}{{if .BillAddress.Country}}
               <Country >{{ .BillAddress.Country}}</Country>{{end}}{{if .BillAddress.Note}}
               <Note >{{ .BillAddress.Note}}</Note>{{end}}
            </BillAddress>{{end}}{{if .ShipAddress}}
            <ShipAddress>{{if .ShipAddress.Addr1}}
               <Addr1 >{{ .ShipAddress.Addr1}}</Addr1>{{end}}{{if .ShipAddress.Addr2}}
               <Addr2 >{{ .ShipAddress.Addr2}}</Addr2>{{end}}{{if .ShipAddress.Addr3}}
               <Addr3 >{{ .ShipAddress.Addr3}}</Addr3>{{end}}{{if .ShipAddress.Addr4}}
               <Addr4 >{{ .ShipAddress.Addr4}}</Addr4>{{end}}{{if .ShipAddress.Addr5}}
               <Addr5 >{{ .ShipAddress.Addr5}}</Addr5>{{end}}{{if .ShipAddress.City}}
               <City >{{ .ShipAddress.City}}</City>{{end}}{{if .ShipAddress.State}}
               <State >{{ .ShipAddress.State}}</State>{{end}}{{if .ShipAddress.PostalCode}}
               <PostalCode >{{ .ShipAddress.PostalCode}}</PostalCode>{{end}}{{if .ShipAddress.Country}}
               <Country >{{ .ShipAddress.Country}}</Country>{{end}}{{if .ShipAddress.Note}}
               <Note >{{ .ShipAddress.Note}}</Note>{{end}}
            </ShipAddress>{{end}}{{if .PONumber}}
            <PONumber>{{ .PONumber}}</PONumber>{{end}}{{if or .TermsRef.ListID .TermsRef.FullName}}
            <TermsRef>{{if .TermsRef.ListID}}
                <ListID>{{ .TermsRef.ListID}}</ListID>{{end}}{{if .TermsRef.FullName}}
                <FullName>{{ .TermsRef.FullName}}</FullName>{{end}}
            </TermsRef>{{end}}{{if .DueDate}}
            <DueDate >{{.DueDate}}</DueDate>{{end}}{{if or .SalesRepRef.ListID .SalesRepRef.FullName}}
            <SalesRepRef>{{if .SalesRepRef.ListID}}
               <ListID >{{ .SalesRepRef.ListID}}</ListID>{{end}}{{if .SalesRepRef.FullName}}
               <FullName >{{ .SalesRepRef.FullName}}</FullName>{{end}}
            </SalesRepRef>{{end}}{{if .FOB}}
            <FOB >{{.FOB}}</FOB>{{end}}{{if .ShipDate}}
            <ShipDate >{{.ShipDate}}</ShipDate>{{end}}{{if or .ShipMethodRef.ListID .ShipMethodRef.FullName}}
            <ShipMethodRef>{{if .ShipMethodRef.ListID}}
               <ListID >{{.ShipMethodRef.ListID}}</ListID>{{end}}{{if .ShipMethodRef.FullName}}
               <FullName >{{.ShipMethodRef.FullName}}</FullName>{{end}}
            </ShipMethodRef>{{end}}{{if or .ItemSalesTaxRef.ListID .ItemSalesTaxRef.FullName}}
            <ItemSalesTaxRef>{{if .ItemSalesTaxRef.ListID}}
               <ListID >{{ .ItemSalesTaxRef.ListID}}</ListID>{{end}}{{if .ItemSalesTaxRef.FullName}}
               <FullName >{{ .ItemSalesTaxRef.FullName}}</FullName>{{end}}
            </ItemSalesTaxRef>{{end}}{{if .IsManuallyClosed}}
            <IsManuallyClosed>{{ .IsManuallyClosed}}</IsManuallyClosed>{{end}}{{if .Memo}}
            <Memo >{{.Memo}}</Memo>{{end}}{{if or .CustomerMsgRef.ListID .CustomerMsgRef.FullName}}
            <CustomerMsgRef>{{if .CustomerMsgRef.ListID}}
               <ListID >{{ .CustomerMsgRef.ListID}}</ListID>{{end}}{{if .CustomerMsgRef.FullName}}
               <FullName >{{ .CustomerMsgRef.FullName}}</FullName>{{end}}
            </CustomerMsgRef>{{end}}{{if .IsToBePrinted}}
            <IsToBePrinted >{{.IsToBePrinted}}</IsToBePrinted>{{end}}{{if .IsToBeEmailed}}
            <IsToBeEmailed >{{.IsToBeEmailed}}</IsToBeEmailed>{{end}}{{if or .CustomerSalesTaxCodeRef.ListID .CustomerSalesTaxCodeRef.FullName}}
            <CustomerSalesTaxCodeRef>{{if .CustomerSalesTaxCodeRef.ListID}}
               <ListID >{{.CustomerSalesTaxCodeRef.ListID}}</ListID>{{end}}{{if .CustomerSalesTaxCodeRef.FullName}}
               <FullName >{{ .CustomerSalesTaxCodeRef.FullName}}</FullName>{{end}}
            </CustomerSalesTaxCodeRef>{{end}}{{if .Other}}
            <Other >{{.Other}}</Other>{{end}}{{if .ExchangeRate}}
            <ExchangeRate >{{.ExchangeRate}}</ExchangeRate>{{end}}{{if .ExternalGUID}}
            <ExternalGUID >{{.ExternalGUID}}{{/*regex "0|(\{[0-9a-fA-F]{8}(\-([0-9a-fA-F]{4})){3}\-[0-9a-fA-F]{12}\})"*/}}</ExternalGUID>{{end}}{{if .SalesOrderLineAdds}}{{range $index, $lineAdd := .SalesOrderLineAdds}}
            <SalesOrderLineAdd >{{if or $lineAdd.ItemRef.ListID $lineAdd.ItemRef.FullName}}
               <ItemRef>{{if $lineAdd.ItemRef.ListID}}
                  <ListID >{{ $lineAdd.ItemRef.ListID}}</ListID>{{end}}{{if $lineAdd.ItemRef.FullName}}
                  <FullName >{{ $lineAdd.ItemRef.FullName}}</FullName>{{end}}
               </ItemRef>{{end}}{{if $lineAdd.Desc}}
               <Desc >{{ $lineAdd.Desc}}</Desc>{{end}}{{if $lineAdd.Quantity}}
               <Quantity >{{ $lineAdd.Quantity}}</Quantity>{{end}}{{if $lineAdd.UnitOfMeasure}}
               <UnitOfMeasure >{{ $lineAdd.UnitOfMeasure}}</UnitOfMeasure>{{end}}{{if $lineAdd.Rate}}
               <Rate >{{ $lineAdd.Rate}}</Rate>{{else if $lineAdd.RatePercent}}
               <RatePercent >{{$lineAdd.RatePercent}}</RatePercent>{{else if or $lineAdd.PriceLevelRef.FullName $lineAdd.PriceLevelRef.ListID}}
               <PriceLevelRef>{{if $lineAdd.PriceLevelRef.ListID}}
                  <ListID >{{$lineAdd.PriceLevelRef.ListID}}</ListID>{{end}}{{if $lineAdd.PriceLevelRef.FullName}}
                  <FullName >{{ $lineAdd.PriceLevelRef.FullName}}</FullName>{{end}}
               </PriceLevelRef>{{end}}{{if or $lineAdd.ClassRef.FullName $lineAdd.ClassRef.ListID}}
               <ClassRef>{{if $lineAdd.ClassRef.ListID}}
                  <ListID >{{ $lineAdd.ClassRef.ListID}}</ListID>{{end}}{{if $lineAdd.ClassRef.FullName}}
                  <FullName >{{ $lineAdd.ClassRef.FullName}}</FullName>{{end}}
               </ClassRef>{{end}}{{if $lineAdd.Amount}}
               <Amount >{{ $lineAdd.Amount}}</Amount>{{end}}{{if $lineAdd.OptionForPriceRuleConflict}}
               <OptionForPriceRuleConflict >{{ $lineAdd.OptionForPriceRuleConflict}}</OptionForPriceRuleConflict>{{/*OptionForPriceRuleConflict may have one of the following values: Zero, BasePrice*/}}{{end}}{{if or $lineAdd.InventorySiteRef.FullName $lineAdd.InventorySiteRef.ListID}}
               <InventorySiteRef>{{if $lineAdd.InventorySiteRef.ListID}}
                  <ListID >{{ $lineAdd.InventorySiteRef.ListID}}</ListID>{{end}}{{if $lineAdd.InventorySiteRef.FullName}}
                  <FullName >{{ $lineAdd.InventorySiteRef.FullName}}</FullName>{{end}}
               </InventorySiteRef>{{end}}{{if or $lineAdd.InventorySiteLocationRef.FullName $lineAdd.InventorySiteLocationRef.ListID}}
               <InventorySiteLocationRef>{{if $lineAdd.InventorySiteLocationRef.ListID}}
                  <ListID >{{ $lineAdd.InventorySiteLocationRef.ListID}}</ListID>{{end}}{{if $lineAdd.InventorySiteLocationRef.FullName}}
                  <FullName >{{ $lineAdd.InventorySiteLocationRef.FullName}}</FullName>{{end}}
               </InventorySiteLocationRef>{{end}}{{if $lineAdd.SerialNumber}}
               <SerialNumber >{{ $lineAdd.SerialNumber}}</SerialNumber>{{else if $lineAdd.LotNumber}}
               <LotNumber >{{ $lineAdd.LotNumber}}</LotNumber>{{end}}{{if or $lineAdd.SalesTaxCodeRef.FullName $lineAdd.SalesTaxCodeRef.ListID}}
               <SalesTaxCodeRef>{{if $lineAdd.SalesTaxCodeRef.ListID}}
                  <ListID >{{ $lineAdd.SalesTaxCodeRef.ListID}}</ListID>{{end}}{{if $lineAdd.SalesTaxCodeRef.FullName}}
                  <FullName >{{ $lineAdd.SalesTaxCodeRef.FullName}}</FullName>{{end}}
               </SalesTaxCodeRef>{{end}}{{if $lineAdd.IsManuallyClosed}}
               <IsManuallyClosed>{{ $lineAdd.IsManuallyClosed}}</IsManuallyClosed>{{end}}{{if $lineAdd.Other1}}
               <Other1 >{{ $lineAdd.Other1}}</Other1>{{end}}{{if $lineAdd.Other2}}
               <Other2 >{{ $lineAdd.Other2}}</Other2>{{end}}{{if $lineAdd.DataExt}}{{ range $j, $dataExt := $lineAdd.DataExt}}
               <DataExt>
                  <OwnerID >{{if $dataExt.OwnerID}}{{ $dataExt.OwnerID}}{{end}}</OwnerID>{{/*Required*/}}
                  <DataExtName >{{if $dataExt.DataExtName}}{{ $dataExt.DataExtName}}{{end}}</DataExtName>{{/*Required*/}}
                  <DataExtValue >{{if $dataExt.DataExtValue}}{{ $dataExt.DataExtValue}}{{end}}</DataExtValue>{{/*Required*/}}
               </DataExt>{{end}}{{end}}
            </SalesOrderLineAdd>{{end}}{{end}}{{if .SalesOrderLineGroupAdds}}{{ range $index, $groupAdd := .SalesOrderLineGroupAdds}}
            <SalesOrderLineGroupAdd>
               <ItemGroupRef>{{/*Required*/}}{{if $groupAdd.ItemGroupRef.ListID}}
                  <ListID >{{ $groupAdd.ItemGroupRef.ListID}}</ListID>{{end}}{{if $groupAdd.ItemGroupRef.FullName}}
                  <FullName >{{ $groupAdd.ItemGroupRef.FullName}}</FullName>{{end}}
               </ItemGroupRef>{{if $groupAdd.Quantity}}
               <Quantity >{{ $groupAdd.Quantity}}</Quantity>{{end}}{{if $groupAdd.UnitOfMeasure}}
               <UnitOfMeasure >{{ $groupAdd.UnitOfMeasure}}</UnitOfMeasure>{{end}}{{if or $groupAdd.InventorySiteRef.FullName $groupAdd.InventorySiteRef.ListID}}
               <InventorySiteRef>{{if $groupAdd.InventorySiteRef.ListID}}
                  <ListID >{{ $groupAdd.InventorySiteRef.ListID}}</ListID>{{end}}{{if $groupAdd.InventorySiteRef.FullName}}
                  <FullName >{{ $groupAdd.InventorySiteRef.FullName}}</FullName>{{end}}
               </InventorySiteRef>{{end}}{{if or $groupAdd.InventorySiteLocationRef.FullName $groupAdd.InventorySiteLocationRef.ListID}}
               <InventorySiteLocationRef>{{if $groupAdd.InventorySiteLocationRef.ListID}}
                  <ListID >{{ $groupAdd.InventorySiteLocationRef.ListID}}</ListID>{{end}}{{if $groupAdd.InventorySiteLocationRef.FullName}}
                  <FullName >{{ $groupAdd.InventorySiteLocationRef.FullName}}</FullName>{{end}}
               </InventorySiteLocationRef>{{end}}{{if $groupAdd.DataExt}}{{ range $j, $dataExt := $groupAdd.DataExt}}
               <DataExt>
                  <OwnerID >{{if $dataExt.OwnerID}}{{ $dataExt.OwnerID}}{{end}}</OwnerID>{{/*Required*/}}
                  <DataExtName >{{if $dataExt.DataExtName}}{{ $dataExt.DataExtName}}{{end}}</DataExtName>{{/*Required*/}}
                  <DataExtValue >{{if $dataExt.DataExtValue}}{{ $dataExt.DataExtValue}}{{end}}</DataExtValue>{{/*Required*/}}
               </DataExt>{{end}}{{end}}
            </SalesOrderLineGroupAdd>{{end}}{{end}}
         </SalesOrderAdd>{{/*<IncludeRetElement >STRTYPE</IncludeRetElement>*/}}
      </SalesOrderAddRq>
   </QBXMLMsgsRq>
</QBXML>{{end}}