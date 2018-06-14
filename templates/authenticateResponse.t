{{define "authenticateResponse.t"}}<?xml version="1.0" encoding="utf-8"?>
<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
    <soap12:Body>
        <authenticateResponse xmlns="http://developer.intuit.com/">
            <authenticateResult>
                <string>{{if .Ticket}}{{.Ticket}}{{end}}</string>{{/* having no extra white space is very important*/}}
                <string>{{if .DelayUpdate}}{{.DelayUpdate}}{{end}}</string>{{ if .EveryMinLowerLimit}}
                <string>{{.EveryMinLowerLimit}}</string>{{end}}{{/*optional, only add if the field is populated*/}}{{if .MinimumRunEveryNSeconds}}
                <string>{{.MinimumRunEveryNSeconds}}</string>{{end}}{{/*optional, only add if the field is populated*/}}
            </authenticateResult>
        </authenticateResponse>
    </soap12:Body>
</soap12:Envelope>{{end}}