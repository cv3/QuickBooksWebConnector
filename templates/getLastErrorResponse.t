{{define "getLastErrorResponse.t"}}<?xml version="1.0" encoding="utf-8"?>{{/*no extra whitespace is important*/}}
    <soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
        <soap12:Body>
            <getLastErrorResponse xmlns="http://developer.intuit.com/">
                <getLastErrorResult>{{if .LastError}}{{.LastError}}{{end}}</getLastErrorResult>
            </getLastErrorResponse>
        </soap12:Body>
    </soap12:Envelope>
{{end}}