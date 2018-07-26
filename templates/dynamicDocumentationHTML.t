{{define "dynamicDocumentationHTML.t"}}<!DOCTYPE html>
<html>
<head>
<style>
table {
    width:100%;
}
td{
    width:50%;
}
table, th, tr, td {
    border: 1px solid black;
    border-collapse: collapse;
}
th, td {
    padding: 5px;
}
#mapData {
    padding: 0;
    border: 0;
}
#mapData table,
#mapData tr,
#mapData th,
#mapData td {
    border: 0;
    border-collapse: collapse;
}
#mapData table td {
  border: 1px solid black; 
}
#mapData table tr th:first-child {
  border-right: 1px solid black;
}
#mapData table tr:first-child td {
  border-top: 0;
}
#mapData table tr td:first-child {
  border-left: 0;
}
#mapData table tr:last-child td {
  border-bottom: 0;
}
#mapData table tr td:last-child {
  border-right: 0;
}
</style>
</head>
<body>
<h1>Quick Books Web conneector field mapping documentation</h1>
<br>
<p>The following tables show the data in their respective mapping files.  To change a mapping field you must open the file that corresponds to the table. Then add, remove, or modify a block of data</p>
<br>
<p>e.g.<br>
"Name":[<br>
    &emsp;{<br>
        &emsp;&emsp;"data":"firstname",<br>
        &emsp;&emsp;"mappedField":true<br>
    &emsp;},<br>
    &emsp;{<br>
        &emsp;&emsp;"data":" ",<br>
        &emsp;&emsp;"mappedField":false<br>
    &emsp;},<br>
    &emsp;{<br>
        &emsp;&emsp;"data":"lastname",<br>
        &emsp;&emsp;"mappedField":true<br>
    &emsp;}<br>
]<br>
<br>
Indicates the quickbooks Name field will concatenate the fields in top down order, so it would start with the cv3 field firstName and then be concatenate with a hardcoded value of " " and end with the cv3 field lastName
<br>
{{range $fileName, $mappingObjects := .}}
    <table >
    <caption><h1>{{$fileName}}</h1></caption>
        <tr>
            <th>QuickBooks Field</th>
            <th id="mapData">
                <table>
                    <tr>
                        <td colspan="2" >Mapped To</td>
                    </tr>
                    <tr>
                        <td>Data</td>
                        <td>Type</td>
                    </tr>
                </table>
            </th> 
        </tr>{{range $key, $mObj := $mappingObjects}}
        <tr>
            <td>{{$key}}</td>
            <td id="mapData">
                <table>{{range $type, $data := $mObj}}
                    <tr>
                        <td>{{if not $data.MappedField}}"{{$data.Data}}"{{else if $data.Data}}{{$data.Data}}{{end}}</td>
                        <td>{{if $data.MappedField}}Mapping{{else}}Hardcoded{{end}}</td>
                    </tr>{{end}}
                </table>
            </td>
        </tr>{{end}}
    </table>
    <br>{{end}}
</body>
</html>
{{end}}