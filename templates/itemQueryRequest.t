{{define "itemQueryRequest.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="7.0"?>
<QBXML>
  <QBXMLMsgsRq onError="stopOnError">
    <ItemQueryRq requestID="SXRlbVF1ZXJ5fDEyMA==" >{{  range $id := .ListID}}
      <ListID>{{$id}}</ListID>{{end}}{{  range $name := .FullName}}
      <FullName>{{$name}}</FullName>{{end}}{{ if .MaxReturned}}
      <MaxReturned>{{.MaxReturned}}</MaxReturned>{{end}}{{ if .ActiveStatus}}
      <ActiveStatus>{{.ActiveStatus}}</ActiveStatus>{{end}}{{ if .FromModifiedDate}}
      <FromModifiedDate>{{.FromModifiedDate}}</FromModifiedDate>{{end}}{{ if .ToModifiedDate}}
      <ToModifiedDate>{{.ToModifiedDate}}</ToModifiedDate>{{end}}{{ if .NameFilter.MatchCriterion }}
      <NameFilter>
        <MatchCriterion>{{.NameFilter.MatchCriterion}}</MatchCriterion>{{ if .NameFilter.Name}}
        <Name>{{.NameFilter.Name}}{{end}}
      </NameFilter>{{end}}{{if .NameRangeFilter.FromName}}
      <NameRangeFilter>
        <FromName>{{.NameRangeFilter.FromName}}</FromName>
        <ToName>{{if .NameRangeFilter.ToName}}{{.NameRangeFilter.ToName}}{{end}}</ToName>
      </NameRangeFilter>{{end}}{{range $element := .IncludeRetElement}}
      <IncludeRetElement>{{$element}}</IncludeRetElement>{{end}}{{range $id := .OwnerID}}
      <OwnerID>{{$id}}</OwnerID>{{end}}
    </ItemQueryRq>
  </QBXMLMsgsRq>
</QBXML>{{end}}