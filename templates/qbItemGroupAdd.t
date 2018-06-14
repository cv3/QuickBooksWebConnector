{{define "qbItemGroupAdd.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
    <QBXMLMsgsRq onError="stopOnError">
        <ItemGroupAddRq>
            <ItemGroupAdd>
                <Name>{{if .Name}}{{.Name}}{{end}}</Name>{{if .BarCode.BarCodeValue}}
                <BarCode>
                    <BarCodeValue>{{.BarCode.BarCodeValue}}</BarCodeValue>{{if .BarCode.AssignEvenIfUsed}}
                    <AssignEvenIfUsed>{{ .BarCode.AssignEvenIfUsed}}</AssignEvenIfUsed>{{end}}{{if .BarCode.AllowOverride}}
                    <AllowOverride>{{ .BarCode.AllowOverride}}</AllowOverride>{{end}}
                </BarCode>{{end}}{{if.IsActive}}
                <IsActive>{{ .IsActive}}</IsActive>{{end}}{{if .ItemDesc}}
                <ItemDesc>{{ .ItemDesc}}</ItemDesc>{{end}}{{if or .UnitOfMeasureSetRef.ListID .UnitOfMeasureSetRef.FullName}}
                <UnitOfMeasureSetRef>{{if .UnitOfMeasureSetRef.ListID}}
                    <ListID>{{ .UnitOfMeasureSetRef.ListID}}</ListID>{{end}}
                    <FullName>{{ .UnitOfMeasureSetRef.FullName}}</FullName>
                </UnitOfMeasureSetRef>{{end}}{{if .IsPrintItemsInGroup}}
                <IsPrintItemsInGroup>{{ .IsPrintItemsInGroup}}</IsPrintItemsInGroup>{{end}}{{if .ExternalGUID}}
                <ExternalGUID>{{ .ExternalGUID}}</ExternalGUID>{{end}}{{range $index, $groupLine := .ItemGroupLine}}
                <ItemGroupLine>{{if $groupLine.ItemRef.FullName}}
                    <ItemRef>{{if $groupLine.ItemRef.ListID}}
                        <ListID>{{ $groupLine.ItemRef.ListID}}</ListID>{{end}}
                        <FullName>{{ $groupLine.ItemRef.FullName}}</FullName>
                    </ItemRef>{{end}}{{if $groupLine.Quantity}}
                    <Quantity>{{ $groupLine.Quantity}}</Quantity>{{end}}{{if $groupLine.UnitOfMeasure}}
                    <UnitOfMeasure>{{ $groupLine.UnitOfMeasure}}</UnitOfMeasure>{{end}}
                </ItemGroupLine>{{end}}
            </ItemGroupAdd>
        </ItemGroupAddRq>
    </QBXMLMsgsRq>
</QBXML>{{end}}