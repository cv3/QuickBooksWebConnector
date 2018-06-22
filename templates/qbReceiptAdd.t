{{define "qbReceiptAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
   <QBXMLMsgsRq onError="stopOnError">
      <SalesReceiptAddRq>
         <SalesReceiptAdd >{{if or .CustomerRef.ListID .CustomerRef.FullName}}
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
            </ShipAddress>{{end}}{{if .IsPending}}
            <IsPending >{{ .IsPending}}</IsPending>{{end}}{{if .CheckNumber}}
            <CheckNumber >{{ .CheckNumber}}CheckNumber>{{end}}{{if or .PaymentMethodRef.ListID .PaymentMethodRef.FullName}}
            <PaymentMethodRef>{{if .PaymentMethodRef.ListID}}
               <ListID >{{ .PaymentMethodRef.ListID}}</ListID>{{end}}{{if .PaymentMethodRef.FullName}}
               <FullName >{{ .PaymentMethodRef.FullName}}</FullName>{{end}}
            </PaymentMethodRef>{{end}}{{if .DueDate}}
            <DueDate >{{.DueDate}}</DueDate>{{end}}{{if or .SalesRepRef.ListID .SalesRepRef.FullName}}
            <SalesRepRef>{{if .SalesRepRef.ListID}}
               <ListID >{{ .SalesRepRef.ListID}}</ListID>{{end}}{{if .SalesRepRef.FullName}}
               <FullName >{{ .SalesRepRef.FullName}}</FullName>{{end}}
            </SalesRepRef>{{end}}{{if .ShipDate}}
            <ShipDate >{{.ShipDate}}</ShipDate>{{end}}{{if or .ShipMethodRef.ListID .ShipMethodRef.FullName}}
            <ShipMethodRef>{{if .ShipMethodRef.ListID}}
               <ListID >{{.ShipMethodRef.ListID}}</ListID>{{end}}{{if .ShipMethodRef.FullName}}
               <FullName >{{.ShipMethodRef.FullName}}</FullName>{{end}}
            </ShipMethodRef>{{end}}{{if .FOB}}
            <FOB >{{.FOB}}</FOB>{{end}}{{if or .ItemSalesTaxRef.ListID .ItemSalesTaxRef.FullName}}
            <ItemSalesTaxRef>{{if .ItemSalesTaxRef.ListID}}
               <ListID >{{ .ItemSalesTaxRef.ListID}}</ListID>{{end}}{{if .ItemSalesTaxRef.FullName}}
               <FullName >{{ .ItemSalesTaxRef.FullName}}</FullName>{{end}}
            </ItemSalesTaxRef>{{end}}{{if .Memo}}
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
            </CustomerSalesTaxCodeRef>{{end}}{{if or .DepositToAccountRef.ListID .DepositToAccountRef.FullName}}
            <DepositToAccountRef>{{if .DepositToAccountRef.ListID}}
               <ListID >{{ .DepositToAccountRef.ListID}}</ListID>{{end}}{{if .DepositToAccountRef.FullName}}
               <FullName >{{ .DepositToAccountRef.FullName}}</FullName>{{end}}
            </DepositToAccountRef>{{end}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardNumber}}
            <CreditCardTxnInfo>
               <CreditCardTxnInputInfo>{{/*required*/}}
                  <CreditCardNumber >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardNumber}}</CreditCardNumber>{{/*required*/}}
                  <ExpirationMonth >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationMonth}}</ExpirationMonth>{{/*required*/}}
                  <ExpirationYear >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationYear}}</ExpirationYear>{{/*required*/}}
                  <NameOnCard >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.NameOnCard}}</NameOnCard>{{/*required*/}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardAddress}}
                  <CreditCardAddress >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardAddress}}</CreditCardAddress>{{end}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardPostalCode}}
                  <CreditCardPostalCode >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardPostalCode}}</CreditCardPostalCode>{{end}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.CommercialCardCode}}
                  <CommercialCardCode >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.CommercialCardCode}}</CommercialCardCode>{{end}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.TransactionMode}}
                  <TransactionMode >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.TransactionMode}}</TransactionMode>{{/*TransactionMode may have one of the following values: CardNotPresent [DEFAULT], CardPresent*/}}{{end}}{{if .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardTxnType}}
                  <CreditCardTxnType >{{ .CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardTxnType}}</CreditCardTxnType>{{/*CreditCardTxnType may have one of the following values: Authorization, Capture, Charge, Refund, VoiceAuthorization*/}}{{end}}
               </CreditCardTxnInputInfo>{{if .CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}
               <CreditCardTxnResultInfo>{{/*required*/}}
                  <ResultCode >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}{{end}}</ResultCode>{{/*required*/}}
                  <ResultMessage >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultMessage}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultMessage}}{{end}}</ResultMessage>{{/*required*/}}
                  <CreditCardTransID >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CreditCardTransID}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CreditCardTransID}}{{end}}</CreditCardTransID>{{/*required*/}}
                  <MerchantAccountNumber >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.MerchantAccountNumber}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.MerchantAccountNumber}}{{end}}</MerchantAccountNumber>{{/*required*/}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AuthorizationCode}}
                  <AuthorizationCode >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AuthorizationCode}}</AuthorizationCode>{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSStreet}}
                  <AVSStreet >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSStreet}}</AVSStreet>{{/*AVSStreet may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSZip}}
                  <AVSZip >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSZip}}</AVSZip>{{/*AVSZip may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CardSecurityCodeMatch}}
                  <CardSecurityCodeMatch >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CardSecurityCodeMatch}}</CardSecurityCodeMatch>{{/*CardSecurityCodeMatch may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ReconBatchID}}
                  <ReconBatchID >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ReconBatchID}}</ReconBatchID>{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentGroupingCode}}
                  <PaymentGroupingCode >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentGroupingCode}}</PaymentGroupingCode>{{end}}
                  <PaymentStatus >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentStatus}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentStatus}}{{end}}</PaymentStatus>{{/*Required: PaymentStatus may have one of the following values: Unknown, Completed*/}}
                  <TxnAuthorizationTime >{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationTime}}{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationTime}}{{end}}</TxnAuthorizationTime>{{/*Required*/}}{{if .CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationStamp}}
                  <TxnAuthorizationStamp >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationStamp}}</TxnAuthorizationStamp>{{end}}{{if .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ClientTransID}}
                  <ClientTransID >{{ .SalesReceiptAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ClientTransID}}</ClientTransID>{{end}}
               </CreditCardTxnResultInfo>{{end}}
            </CreditCardTxnInfo>{{end}}{{if .Other}}
            <Other >{{.Other}}</Other>{{end}}{{if .ExchangeRate}}
            <ExchangeRate >{{.ExchangeRate}}</ExchangeRate>{{end}}{{if .ExternalGUID}}
            <ExternalGUID >{{.ExternalGUID}}{{/*regex "0|(\{[0-9a-fA-F]{8}(\-([0-9a-fA-F]{4})){3}\-[0-9a-fA-F]{12}\})"*/}}</ExternalGUID>{{end}}{{if .SalesReceiptLineAdds}}{{range $index, $lineAdd := .SalesReceiptLineAdds}}
            <SalesReceiptLineAdd >{{if or $lineAdd.ItemRef.ListID $lineAdd.ItemRef.FullName}}
               <ItemRef>{{if $lineAdd.ItemRef.ListID}}
                  <ListID >{{ $lineAdd.ItemRef.ListID}}</ListID>{{end}}{{if $lineAdd.ItemRef.FullName}}
                  <FullName >{{ $lineAdd.ItemRef.FullName}}</FullName>{{end}}
               </ItemRef>{{end}}{{if $lineAdd.Desc}}
               <Desc >{{ $lineAdd.Desc}}</Desc>{{end}}{{if $lineAdd.Quantity}}
               <Quantity >{{ $lineAdd.Quantity}}</Quantity>{{end}}{{if $lineAdd.UnitOfMeasure}}
               <UnitOfMeasure >{{ $lineAdd.UnitOfMeasure}}</UnitOfMeasure>{{end}}{{if $lineAdd.Rate}}
               <Rate >{{ $lineAdd.Rate}}</Rate>{{else if $lineAdd.RatePercent}}
               <RatePercent >{{$lineAdd.RatePercent}}</RatePercent>{{else if or $lineAdd.PriceLevelRef.ListID $lineAdd.PriceLevelRef.FullName}}
               <PriceLevelRef>{{if $lineAdd.PriceLevelRef.ListID}}
                  <ListID >{{$lineAdd.PriceLevelRef.ListID}}</ListID>{{end}}{{if $lineAdd.PriceLevelRef.FullName}}
                  <FullName >{{ $lineAdd.PriceLevelRef.FullName}}</FullName>{{end}}
               </PriceLevelRef>{{end}}{{if or $lineAdd.ClassRef.ListID $lineAdd.ClassRef.FullName}}
               <ClassRef>{{if $lineAdd.ClassRef.ListID}}
                  <ListID >{{ $lineAdd.ClassRef.ListID}}</ListID>{{end}}{{if $lineAdd.ClassRef.FullName}}
                  <FullName >{{ $lineAdd.ClassRef.FullName}}</FullName>{{end}}
               </ClassRef>{{end}}{{if $lineAdd.Amount}}
               <Amount >{{ $lineAdd.Amount}}</Amount>{{end}}{{if $lineAdd.OptionForPriceRuleConflict}}
               <OptionForPriceRuleConflict >{{ $lineAdd.OptionForPriceRuleConflict}}</OptionForPriceRuleConflict>{{/*OptionForPriceRuleConflict may have one of the following values: Zero, BasePrice*/}}{{end}}{{if or $lineAdd.InventorySiteRef.ListID $lineAdd.InventorySiteRef.FullName}}
               <InventorySiteRef>{{if $lineAdd.InventorySiteRef.ListID}}
                  <ListID >{{ $lineAdd.InventorySiteRef.ListID}}</ListID>{{end}}{{if $lineAdd.InventorySiteRef.FullName}}
                  <FullName >{{ $lineAdd.InventorySiteRef.FullName}}</FullName>{{end}}
               </InventorySiteRef>{{end}}{{if or $lineAdd.InventorySiteLocationRef.ListID $lineAdd.InventorySiteLocationRef.FullName}}
               <InventorySiteLocationRef>{{if $lineAdd.InventorySiteLocationRef.ListID}}
                  <ListID >{{ $lineAdd.InventorySiteLocationRef.ListID}}</ListID>{{end}}{{if $lineAdd.InventorySiteLocationRef.FullName}}
                  <FullName >{{ $lineAdd.InventorySiteLocationRef.FullName}}</FullName>{{end}}
               </InventorySiteLocationRef>{{end}}{{if $lineAdd.SerialNumber}}
               <SerialNumber >{{ $lineAdd.SerialNumber}}</SerialNumber>{{else if $lineAdd.LotNumber}}
               <LotNumber >{{ $lineAdd.LotNumber}}</LotNumber>{{end}}{{if $lineAdd.ServiceDate}}
               <ServiceDate >{{ $lineAdd.ServiceDate}}</ServiceDate>{{end}}{{if or $lineAdd.SalesTaxCodeRef.ListID $lineAdd.SalesTaxCodeRef.FullName}}
               <SalesTaxCodeRef>{{if $lineAdd.SalesTaxCodeRef.ListID}}
                  <ListID >{{ $lineAdd.SalesTaxCodeRef.ListID}}</ListID>{{end}}{{if $lineAdd.SalesTaxCodeRef.FullName}}
                  <FullName >{{ $lineAdd.SalesTaxCodeRef.FullName}}</FullName>{{end}}
               </SalesTaxCodeRef>{{end}}{{if or $lineAdd.OverrideItemAccountRef.ListID $lineAdd.OverrideItemAccountRef.FullName}}
               <OverrideItemAccountRef>{{if $lineAdd.OverrideItemAccountRef.ListID}}
                  <ListID >{{ $lineAdd.OverrideItemAccountRef.ListID}}</ListID>{{end}}{{if $lineAdd.OverrideItemAccountRef.FullName}}
                  <FullName >{{ $lineAdd.OverrideItemAccountRef.FullName}}</FullName>{{end}}
               </OverrideItemAccountRef>{{end}}{{if $lineAdd.Other1}}
               <Other1 >{{ $lineAdd.Other1}}</Other1>{{end}}{{if $lineAdd.Other2}}
               <Other2 >{{ $lineAdd.Other2}}</Other2>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardNumber}}
               <CreditCardTxnInfo>
                  <CreditCardTxnInputInfo>{{/*required*/}}
                     <CreditCardNumber >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardNumber}}</CreditCardNumber>{{/*required*/}}
                     <ExpirationMonth >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationMonth}}</ExpirationMonth>{{/*required*/}}
                     <ExpirationYear >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.ExpirationYear}}</ExpirationYear>{{/*required*/}}
                     <NameOnCard >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.NameOnCard}}</NameOnCard>{{/*required*/}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardAddress}}
                     <CreditCardAddress >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardAddress}}</CreditCardAddress>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardPostalCode}}
                     <CreditCardPostalCode >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardPostalCode}}</CreditCardPostalCode>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CommercialCardCode}}
                     <CommercialCardCode >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CommercialCardCode}}</CommercialCardCode>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.TransactionMode}}
                     <TransactionMode >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.TransactionMode}}</TransactionMode>{{/*TransactionMode may have one of the following values: CardNotPresent [DEFAULT], CardPresent*/}}{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardTxnType}}
                  <CreditCardTxnType >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnInputInfo.CreditCardTxnType}}</CreditCardTxnType>{{/*CreditCardTxnType may have one of the following values: Authorization, Capture, Charge, Refund, VoiceAuthorization*/}}{{end}}
               </CreditCardTxnInputInfo>{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}
               <CreditCardTxnResultInfo>{{/*required*/}}
                  <ResultCode >{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultCode}}{{end}}</ResultCode>{{/*required*/}}
                     <ResultMessage >{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultMessage}}{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ResultMessage}}{{end}}</ResultMessage>{{/*required*/}}
                     <CreditCardTransID >{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CreditCardTransID}}{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CreditCardTransID}}{{end}}</CreditCardTransID>{{/*required*/}}
                     <MerchantAccountNumber >{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.MerchantAccountNumber}}{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.MerchantAccountNumber}}{{end}}</MerchantAccountNumber>{{/*required*/}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AuthorizationCode}}
                     <AuthorizationCode >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AuthorizationCode}}</AuthorizationCode>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSStreet}}
                     <AVSStreet >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSStreet}}</AVSStreet>{{/*AVSStreet may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSZip}}
                     <AVSZip >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.AVSZip}}</AVSZip>{{/*AVSZip may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CardSecurityCodeMatch}}
                     <CardSecurityCodeMatch >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.CardSecurityCodeMatch}}</CardSecurityCodeMatch>{{/*CardSecurityCodeMatch may have one of the following values: Pass, Fail, NotAvailable*/}}{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ReconBatchID}}
                     <ReconBatchID >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ReconBatchID}}</ReconBatchID>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentGroupingCode}}
                     <PaymentGroupingCode >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentGroupingCode}}</PaymentGroupingCode>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentStatus}}
                     <PaymentStatus >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.PaymentStatus}}</PaymentStatus>{{end}}{{/*Required: PaymentStatus may have one of the following values: Unknown, Completed*/}}
                     <TxnAuthorizationTime >{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationTime}}{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationTime}}{{end}}</TxnAuthorizationTime>{{/*Required*/}}{{if .CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationStamp}}
                     <TxnAuthorizationStamp >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.TxnAuthorizationStamp}}</TxnAuthorizationStamp>{{end}}{{if $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ClientTransID}}
                     <ClientTransID >{{ $lineAdd.CreditCardTxnInfo.CreditCardTxnResultInfo.ClientTransID}}</ClientTransID>{{end}}
                  </CreditCardTxnResultInfo>{{end}}
               </CreditCardTxnInfo>{{end}}{{if $lineAdd.DataExt}}{{ range $j, $dataExt := $lineAdd.DataExt}}
               <DataExt>
                  <OwnerID >{{if $dataExt.OwnerID}}{{ $dataExt.OwnerID}}{{end}}</OwnerID>{{/*Required*/}}
                  <DataExtName >{{if $dataExt.DataExtName}}{{ $dataExt.DataExtName}}{{end}}</DataExtName>{{/*Required*/}}
                  <DataExtValue >{{if $dataExt.DataExtValue}}{{ $dataExt.DataExtValue}}{{end}}</DataExtValue>{{/*Required*/}}
               </DataExt>{{end}}{{end}}
            </SalesReceiptLineAdd>{{end}}{{end}}{{if .SalesReceiptLineGroupAdd}}{{ range $index, $groupAdd := .SalesReceiptLineGroupAdd}}
            <SalesReceiptLineGroupAdd>
               <ItemGroupRef>{{/*Required*/}}{{if $groupAdd.ItemGroupRef.ListID}}
                  <ListID >{{ $groupAdd.ItemGroupRef.ListID}}</ListID>{{end}}{{if $groupAdd.ItemGroupRef.FullName}}
                  <FullName >{{ $groupAdd.ItemGroupRef.FullName}}</FullName>{{end}}
               </ItemGroupRef>{{if $groupAdd.Quantity}}
               <Quantity >{{ $groupAdd.Quantity}}</Quantity>{{end}}{{if $groupAdd.UnitOfMeasure}}
               <UnitOfMeasure >{{ $groupAdd.UnitOfMeasure}}</UnitOfMeasure>{{end}}{{if or $groupAdd.InventorySiteRef.ListID $groupAdd.InventorySiteRef.FullName}}
               <InventorySiteRef>{{if $groupAdd.InventorySiteRef.ListID}}
                  <ListID >{{ $groupAdd.InventorySiteRef.ListID}}</ListID>{{end}}{{if $groupAdd.InventorySiteRef.FullName}}
                  <FullName >{{ $groupAdd.InventorySiteRef.FullName}}</FullName>{{end}}
               </InventorySiteRef>{{end}}{{if or $groupAdd.InventorySiteLocationRef.ListID $groupAdd.InventorySiteLocationRef.FullName}}
               <InventorySiteLocationRef>{{if $groupAdd.InventorySiteLocationRef.ListID}}
                  <ListID >{{ $groupAdd.InventorySiteLocationRef.ListID}}</ListID>{{end}}{{if $groupAdd.InventorySiteLocationRef.FullName}}
                  <FullName >{{ $groupAdd.InventorySiteLocationRef.FullName}}</FullName>{{end}}
               </InventorySiteLocationRef>{{end}}{{if $groupAdd.DataExt}}{{ range $j, $dataExt := $groupAdd.DataExt}}
               <DataExt>
                  <OwnerID >{{if $dataExt.OwnerID}}{{ $dataExt.OwnerID}}{{end}}</OwnerID>{{/*Required*/}}
                  <DataExtName >{{if $dataExt.DataExtName}}{{ $dataExt.DataExtName}}{{end}}</DataExtName>{{/*Required*/}}
                  <DataExtValue >{{if $dataExt.DataExtValue}}{{ $dataExt.DataExtValue}}{{end}}</DataExtValue>{{/*Required*/}}
               </DataExt>{{end}}{{end}}
            </SalesReceiptLineGroupAdd>{{end}}{{end}}
         </SalesReceiptAdd>{{/*<IncludeRetElement >STRTYPE</IncludeRetElement>*/}}
      </SalesReceiptAddRq>
   </QBXMLMsgsRq>
</QBXML>{{end}}