{{define "qbCustomerAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
    <QBXML>
        <QBXMLMsgsRq onError="stopOnError">
            <CustomerAddRq>
                <CustomerAdd>{{if .Name}}
                    <Name>{{.Name}}</Name>{{end}}{{if .IsActive}}
                    <IsActive>{{.IsActive}}</IsActive>{{end}}{{if or .ClassRef.FullName .ClassRef.ListID}}
                    <ClassRef>{{if .ClassRef.ListID}}
                        <ListID>{{.CLassRef.ListID}}</ListID>{{end}}{{if .CustomerRef.FullName}}
                        <FullName>{{.CustomerRef.FullName}}</FullName>{{end}}
                    </ClassRef>{{end}}{{ if or .ParentRef.FullName .ParentRef.ListID}}
                    <ParentRef>{{if .ParentRef.ListID}}
                        <ListID>{{.ParentRef.ListID}}</ListID>{{end}}{{ if .ParentRef.FullName}}
                        <FullName>{{.ParentRef.FullName}}</FullName>{{end}}
                    </ParentRef>{{end}}{{if .CompanyName}}
                    <CompanyName>{{.CompanyName}}</CompanyName>{{end}}{{if .Salutation}}
                    <Salutation>{{.Salutation}}Salutation>{{end}}{{if .FirstName}}
                    <FirstName>{{.FirstName}}</FirstName>{{end}}{{if .MiddleName}}
                    <MiddleName>{{.MiddleName}}</MiddleName>{{end}}{{if .LastName}}
                    <LastName>{{.LastName}}</LastName>{{end}}{{if .JobTitle}}
                    <JobTitle>{{.JobTitle}}</JobTitle>{{end}}{{if or .BillAddress.Addr1 .BillAddress.State}}
                    <BillAddress>{{if .BillAddress.Addr1}}
                        <Addr1>{{.BillAddress.Addr1}}</Addr1>{{end}}{{if .BillAddress.Addr2}}
                        <Addr2>{{.BillAddress.Addr2}}</Addr2>{{end}}{{if .BillAddress.Addr3}}
                        <Addr3>{{.BillAddress.Addr3}}</Addr3>{{end}}{{if .BillAddress.Addr4}}
                        <Addr4>{{.BillAddress.Addr4}}</Addr4>{{end}}{{if .BillAddress.Addr5}}
                        <Addr5>{{.BillAddress.Addr5}}</Addr5>{{end}}{{if .BillAddress.City}}
                        <City>{{.BillAddress.City}}</City>{{end}}{{ if .BillAddress.State}}
                        <State>{{.BillAddress.State}}</State>{{end}}{{if .BillAddress.PostalCode}}
                        <PostalCode>{{.BillAddress.PostalCode}}</PostalCode>{{end}}{{if .BillAddress.Country}}
                        <Country>{{.BillAddress.Country}}</Country>{{end}}{{if .BillAddress.Note}}
                        <Note>{{.BillAddress.Note}}</Note>{{end}}
                    </BillAddress>{{end}}{{if or .ShipAddress.Addr1 .ShipAddress.State}}
                    <ShipAddress>{{if .ShipAddress.Addr1}}
                        <Addr1>{{.ShipAddress.Addr1}}</Addr1>{{end}}{{if .ShipAddress.Addr2}}
                        <Addr2>{{.ShipAddress.Addr2}}</Addr2>{{end}}{{if .ShipAddress.Addr3}}
                        <Addr3>{{.ShipAddress.Addr3}}</Addr3>{{end}}{{if .ShipAddress.Addr4}}
                        <Addr4>{{.ShipAddress.Addr4}}</Addr4>{{end}}{{if .ShipAddress.Addr5}}
                        <Addr5>{{.ShipAddress.Addr5}}</Addr5>{{end}}{{if .ShipAddress.City}}
                        <City>{{.ShipAddress.City}}</City>{{end}}{{if .ShipAddress.State}}
                        <State>{{.ShipAddress.State}}</State>{{end}}{{if .ShipAddress.PostalCode}}
                        <PostalCode>{{.ShipAddress.PostalCode}}</PostalCode>{{end}}{{if .ShipAddress.Country}}
                        <Country>{{.ShipAddress.Country}}</Country>{{end}}{{if .ShipAddress.Note}}
                        <Note>{{.ShipAddress.Note}}</Note>{{end}}
                    </ShipAddress>{{end}}{{range $index, $shipTo := .ShipToAddress}}
                    <ShipToAddress>{{if $shipTo.Name}}
                        <Name>{{$shipTo.Name}}</Name>{{end}}{{if $shipTo.Addr1}}
                        <Addr1>{{$shipTo.Addr1}}</Addr1>{{end}}{{if $shipTo.Addr2}}
                        <Addr2>{{$shipTo.Addr2}}</Addr2>{{end}}{{if $shipTo.Addr3}}
                        <Addr3>{{$shipTo.Addr3}}</Addr3>{{end}}{{if $shipTo.Addr4}}
                        <Addr4>{{$shipTo.Addr4}}</Addr4>{{end}}{{if $shipTo.Addr5}}
                        <Addr5>{{$shipTo.Addr5}}</Addr5>{{end}}{{if $shipTo.City}}
                        <City>{{$shipTo.City}}</City>{{end}}{{if $shipTo.State}}
                        <State>{{$shipTo.State}}</State>{{end}}{{if $shipTo.PostalCode}}
                        <PostalCode>{{$shipTo.PostalCode}}</PostalCode>{{end}}{{if $shipTo.Country}}
                        <Country>{{$shipTo.Country}}</Country>{{end}}{{if $shipTo.Note}}
                        <Note>{{$shipTo.Note}}</Note>{{end}}
                    </ShipToAddress>{{end}}{{if .Phone}}
                    <Phone>{{.Phone}}</Phone>{{end}}{{if .AltPhone}}
                    <AltPhone>{{.AltPhone}}</AltPhone>{{end}}{{if .Fax}}
                    <Fax>{{.Fax}}</Fax>{{end}}{{if .Email}}
                    <Email>{{.Email}}</Email>{{end}}{{if .Cc}}
                    <Cc>{{.Cc}}</Cc>{{end}}{{if .Contact}}
                    <Contact>{{.Contact}}</Contact>{{end}}{{if .AltContact}}
                    <AltContact>{{.AltContact}}</AltContact>{{end}}{{ range $index, $ref := .AdditionalContactRef}}
                    <AdditionalContactRef>{{if $ref.ContactName}}
                        <ContactName>{{$ref.ContactName}}</ContactName>{{end}}{{ if $ref.ContactValue}}
                        <ContactValue>{{$ref.ContactValue}}</ContactValue>{{end}}
                    </AdditionalContactRef>{{end}}{{range $index, $contact := .Contacts}}
                    <Contacts>{{if $contact.Salutation}}
                        <Salutation>{{$contact.Salutation}}</Salutation>{{end}}{{if $contact.FirstName}}
                        <FirstName>{{$contact.FirstName}}</FirstName>{{end}}{{if $contact.MiddleName}}
                        <MiddleName>{{$contact.MiddleName}}</MiddleName>{{end}}{{if $contact.LastName}}
                        <LastName>{{$contact.LastName}}</LastName>{{end}}{{if $contact.JobTitle}}
                        <JobTitle>{{$contact.JobTitle}}</JobTitle>{{end}}{{range $j, $ref := $contact.AdditionalContactRef}}
                        <AdditionalContactRef>{{if $ref.ContactName}}
                            <ContactName>{{$ref.ContactName}}</ContactName>{{end}}{{if $ref.ContactValue}}
                            <ContactValue>{{$ref.ContactValue}}</ContactValue>{{end}}
                        </AdditionalContactRef>{{end}}
                    </Contacts>{{end}}{{if or .CustomerTypeRef.FullName .CustomerTypeRef.ListID}}
                    <CustomerTypeRef>{{if .CustomerTypeRef.ListID}}
                        <ListID>{{.CustomerTypeRef.ListID}}</ListID>{{end}}{{if .CustomerTypeRef.FullName}}
                        <FullName>{{.CustomerTypeRef.FullName}}</FullName>{{end}}
                    </CustomerTypeRef>{{end}}{{if or .TermsRef.FullName .TermsRef.ListID}}
                    <TermsRef>{{if .TermsRef.ListID}}
                        <ListID>{{.TermsRef.ListID}}</ListID>{{end}}{{if .TermsRef.FullName}}
                        <FullName>{{.TermsRef.FullName}}</FullName>{{end}}
                    </TermsRef>{{end}}{{if or .SalesRepRef.FullName .SalesRepRef.ListID}}
                    <SalesRepRef>{{if .SalesRepRef.ListID}}
                        <ListID>{{.SalesRepRef.ListID}}</ListID>{{end}}{{if .SalesRepRef.FullName}}
                        <FullName>{{.SalesRepREf.FullName}}</FullName>{{end}}
                    </SalesRepRef>{{end}}{{if .OpenBalance}}
                    <OpenBalance>{{.OpenBalance}}</OpenBalance>{{end}}{{if .OpenBalanceDate}}
                    <OpenBalanceDate>{{.OpenBalanceDate}}</OpenBalanceDate>{{end}}{{if or .SalesTaxCodeRef.FullName .SalesTaxCodeRef.ListID}}
                    <SalesTaxCodeRef>{{if .SalesTaxCodeRef.ListID}}
                        <ListID>{{.SalesTaxCodeRef.ListID}}</ListID>{{end}}{{if .SalesTaxCodeRef.FullName}}
                        <FullName>{{.SalesTaxCodeRef.FullName}}</FullName>{{end}}
                    </SalesTaxCodeRef>{{end}}{{if or .ItemSalesTaxRef.FullName .ItemSalesTaxRef.ListID}}
                    <ItemSalesTaxRef>{{if .ItemSalesTaxRef.ListID}}
                        <ListID>{{.ItemSalesTaxRef.ListID}}</ListID>{{end}}{{if .ItemSalesTaxRef.FullName}}
                        <FullName>{{.ItemSalesTaxRef.FullName}}</FullName>{{end}}
                    </ItemSalesTaxRef>{{end}}{{if .ResaleNumber}}
                    <ResaleNumber>{{.ResaleNumber}}</ResaleNumber>{{end}}{{if .AccountNumber}}
                    <AccountNumber>{{.AccountNumber}}</AccountNumber>{{end}}{{if .CreditLimit}}
                    <CreditLimit>{{.CreditLimit}}</CreditLimit>{{end}}{{ if or .PreferredPaymentMethodRef.FullName .PreferredPaymentMethodRef.ListID}}
                    <PreferredPaymentMethodRef>{{if .PreferredPaymentMethodRef.ListID}}
                        <ListID>{{.PreferredPaymentMethodRef.ListID}}</ListID>{{end}}{{if .PreferredPaymentMethodRef.FullName}}
                        <FullName>{{.PreferredPaymentMethodRef.FullName}}</FullName>{{end}}
                    </PreferredPaymentMethodRef>{{end}}{{if .CreditCardInfo.CreditCardNumber}}
                    <CreditCardInfo>
                        <CreditCardNumber>{{.CeditCardInfo.CreditCardNumber}}</CreditCardNumber>{{if .CeditCardInfo.ExpirationMonth}}
                        <ExpirationMonth>{{.CeditCardInfo.ExpirationMonth}}</ExpirationMonth>{{end}}{{if .CeditCardInfo.ExpirationYear}}
                        <ExpirationYear>{{.CeditCardInfo.ExpirationYear}}</ExpirationYear>{{end}}{{if .CeditCardInfo.NameOnCard}}
                        <NameOnCard>{{.CeditCardInfo.NameOnCard}}NameOnCard>{{end}}{{if .CreditCardInfo.CreditCardAddress}}
                        <CreditCardAddress>{{.CreditCardInfo.CreditCardAddress}}</CreditCardAddress>{{end}}{{if .CreditCardInfo.CreditCardPostalCode}}
                        <CreditCardPostalCode>{{.CreditCardInfo.CreditCardPostalCode}}</CreditCardPostalCode>{{end}}
                    </CreditCardInfo>{{end}}{{if .JobStatus}}
                    <JobStatus>{{.JobStatus}}</JobStatus>{{end}}{{if .JobStartDate}}
                    <JobStartDate>{{.JobStartDate}}</JobStartDate>{{end}}{{if .JobProjectedEndDate}}
                    <JobProjectedEndDate>{{.JobProjectedEndDate}}</JobProjectedEndDate>{{end}}{{if .JobEndDate}}
                    <JobEndDate>{{.JobEndDate}}</JobEndDate>{{end}}{{if .JobDesc}}
                    <JobDesc>{{.JobDesc}}</JobDesc>{{end}}{{if or .JobTypeRef.FullName .JobTypeRef.ListID}}
                    <JobTypeRef>{{if .JobTypeRef.ListID}}
                        <ListID>{{.JobTypeRef.ListID}}</ListID>{{end}}{{if .JobTypeRef.FullName}}
                        <FullName>{{.JobTypeRef.FullName}}</FullName>{{end}}
                    </JobTypeRef>{{end}}{{if .Notes}}
                    <Notes>{{.Notes}}</Notes>{{end}}{{range $i, $note := .AdditionalNotes}}
                    <AdditionalNotes>{{if $note.Note}}
                        <Note>{{$note.Note}}</Note>{{end}}
                    </AdditionalNotes>{{end}}{{if .PreferredDeliveryMethod}}
                    <PreferredDeliveryMethod>{{.PreferredDeliveryMethod}}</PreferredDeliveryMethod>{{end}}{{if or .PriceLevelRef.FullName .PriceLevelRef.ListID}}
                    <PriceLevelRef>{{if .PriceLevelRef.ListID}}
                        <ListID>{{.PriceLevelRef.ListID}}</ListID>{{end}}{{if .PriceLevelRef.FullName}}
                        <FullName>{{.PriceLevelRef.FullName}}</FullName>{{end}}
                    </PriceLevelRef>{{end}}{{if .ExternalGUID}}
                    <ExternalGUID>{{.ExternalGUID}}</ExternalGUID>{{end}}{{if or .CurrencyRef.FullName .CurrencyRef.ListID}}
                    <CurrencyRef>{{if .CurrencyRef.ListID}}
                        <ListID>{{.CurrencyRef.ListID}}</ListID>{{end}}{{if .CurrencyRef.FullName}}
                        <FullName>{{.CurrencyRef.FullName}}</FullName>{{end}}
                    </CurrencyRef>{{end}}
                </CustomerAdd>
            </CustomerAddRq>
        </QBXMLMsgsRq>
    </QBXML>{{end}}