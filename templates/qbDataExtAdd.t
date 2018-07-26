{{define "qbDataExtAdd.t"}}        <DataExtAddRq>
            <DataExtAdd>{{if .OwnerID}}
                <OwnerID>{{.OwnerID}}</OwnerID>{{end}}{{if .DataExtName}}
                <DataExtName>{{.DataExtName}}</DataExtName>{{end}}{{if .ListDataExtType}}
                <ListDataExtType>{{.ListDataExtType}}</ListDataExtType>{{end}}{{if or .ListObjRef.FullName .ListObjRef.ListID}}
                <ListObjRef>{{if .ListObjRef.ListID}}
                    <ListID>{{.ListObjRef.ListID}}</ListID>{{end}}{{if .ListObjRef.FullName}}
                    <FullName>{{.ListObjRef.FullName}}</FullName>{{end}}
                </ListObjRef>{{end}}{{if .TxnDataExtType}}
                <TxnDataExtType>{{.TxnDataExtType}}</TxnDataExtType>{{end}}{{if or .TxnID .UseMacro}}
                <TxnID {{if .UseMacro}}useMacro="{{.UseMacro}}"{{end}}>{{if .TxnID}}{{.TxnID}}{{end}}</TxnID>{{end}}{{if .TxnLineID}}
                <TxnLineID>{{.TxnLineID}}</TxnLineID>{{end}}{{if .OtherDataExtType}}
                <OtherDataExtType>{{.OtherDataExtType}}</OtherDataExtType>{{end}}{{if .DataExtValue}}
                <DataExtValue>{{.DataExtValue}}</DataExtValue>{{end}}
            </DataExtAdd>
        </DataExtAddRq>{{end}}