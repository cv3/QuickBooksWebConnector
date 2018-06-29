{{define "qbCustomerMsgAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
    <QBXMLMsgsRq onError="stopOnError">
        <CustomerMsgAddRq>
            <CustomerMsgAdd>{{if .Name}}
                <Name >{{.Name}}</Name>{{end}}{{if .IsActive}}
                <IsActive >{{.IsActive}}</IsActive>{{end}}
            </CustomerMsgAdd>
        </CustomerMsgAddRq>
    </QBXMLMsgsRq>
</QBXML>
{{end}}