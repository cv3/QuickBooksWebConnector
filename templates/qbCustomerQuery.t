{{define "qbCustomerQuery.t"}}  <QBXMLMsgsRq onError="stopOnError">
        <CustomerQueryRq metaData="ENUMTYPE" iterator="ENUMTYPE" iteratorID="UUIDTYPE">
            <!-- BEGIN OR -->
            <ListID>IDTYPE</ListID>
            <!-- optional, may repeat -->
            <!-- OR -->
            <FullName>STRTYPE</FullName>
            <!-- optional, may repeat -->
            <!-- OR -->
            <MaxReturned>INTTYPE</MaxReturned>
            <!-- optional -->
            <!-- ActiveStatus may have one of the following values: ActiveOnly [DEFAULT], InactiveOnly, All -->
            <ActiveStatus>ENUMTYPE</ActiveStatus>
            <!-- optional -->
            <FromModifiedDate>DATETIMETYPE</FromModifiedDate>
            <!-- optional -->
            <ToModifiedDate>DATETIMETYPE</ToModifiedDate>
            <!-- optional -->
            <!-- BEGIN OR -->
            <NameFilter>
                <!-- optional -->
                <!-- MatchCriterion may have one of the following values: StartsWith, Contains, EndsWith -->
                <MatchCriterion>ENUMTYPE</MatchCriterion>
                <!-- required -->
                <Name>STRTYPE</Name>
                <!-- required -->
            </NameFilter>
            <!-- OR -->
            <NameRangeFilter>
                <!-- optional -->
                <FromName>STRTYPE</FromName>
                <!-- optional -->
                <ToName>STRTYPE</ToName>
                <!-- optional -->
            </NameRangeFilter>
            <!-- END OR -->
            <TotalBalanceFilter>
                <!-- optional -->
                <!-- Operator may have one of the following values: LessThan, LessThanEqual, Equal, GreaterThan, GreaterThanEqual -->
                <Operator>ENUMTYPE</Operator>
                <!-- required -->
                <Amount>AMTTYPE</Amount>
                <!-- required -->
            </TotalBalanceFilter>
            <CurrencyFilter>
                <!-- optional -->
                <!-- BEGIN OR -->
                <ListID>IDTYPE</ListID>
                <!-- optional, may repeat -->
                <!-- OR -->
                <FullName>STRTYPE</FullName>
                <!-- optional, may repeat -->
                <!-- END OR -->
            </CurrencyFilter>
            <ClassFilter>
                <!-- optional -->
                <!-- BEGIN OR -->
                <ListID>IDTYPE</ListID>
                <!-- optional, may repeat -->
                <!-- OR -->
                <FullName>STRTYPE</FullName>
                <!-- optional, may repeat -->
                <!-- OR -->
                <ListIDWithChildren>IDTYPE</ListIDWithChildren>
                <!-- optional -->
                <!-- OR -->
                <FullNameWithChildren>STRTYPE</FullNameWithChildren>
                <!-- optional -->
                <!-- END OR -->
            </ClassFilter>
            <!-- END OR -->
            <IncludeRetElement>STRTYPE</IncludeRetElement>
            <!-- optional, may repeat -->
            <OwnerID>GUIDTYPE</OwnerID>
            <!-- optional, may repeat -->
        </CustomerQueryRq>


{{end}}