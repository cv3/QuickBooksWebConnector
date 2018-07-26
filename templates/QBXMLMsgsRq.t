{{define "QBXMLMsgsRq.t"}}<?xml version="1.0" encoding="utf-8"?>
<?qbxml version="13.0"?>
<QBXML>
    <QBXMLMsgsRq onError="stopOnError">{{range .}}
        {{.}}{{end}}
    </QBXMLMsgsRq>
</QBXML>{{end}}