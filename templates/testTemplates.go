package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/TeamFairmont/QuickBooksWebConnector/qbwcGo"
	"github.com/amazingfly/cv3go"
)

var nn = bytes.Buffer{}
var url = "https://machws.cloud.machsoftware.com:8181/machws"

func main() {
	var buff = bytes.Buffer{}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		t, err := template.ParseFiles(f.Name())
		if err != nil {
			fmt.Println(err)
		}

		switch f.Name() {
		case "qbItemGroupAdd.t":
			var ctx = qbwcGo.ItemGroupAdd{}
			err = t.Execute(&buff, ctx)
			if err != nil {
				fmt.Println("Error executing template: ", err)
				fmt.Println(buff.String())
			} else {
				fmt.Println("template: ", f.Name())
			}
		case "qbItemInventoryAdd.t":
			var ctx = qbwcGo.ItemInventoryAdd{}
			err = t.Execute(&buff, ctx)
			if err != nil {
				fmt.Println("Error executing template: ", err)
				fmt.Println(buff.String())
			} else {
				fmt.Println("template: ", f.Name())
			}
		case "qbItemServiceAdd.t":
			var ctx = qbwcGo.ItemServiceAdd{}
			err = t.Execute(&buff, ctx)
			if err != nil {
				fmt.Println("Error executing template: ", err)
				fmt.Println(buff.String())
			} else {
				fmt.Println("template: ", f.Name())
			}
		case "qbReceiptAdd.t":
			var ctx = qbwcGo.SalesReceiptAdd{}
			err = t.Execute(&buff, ctx)
			if err != nil {
				fmt.Println("Error executing template: ", err)
				fmt.Println(buff.String())
			} else {
				fmt.Println("template: ", f.Name())
			}
		case "qbSalesOrderAdd.t":
			var ctx = qbwcGo.SalesOrderAdd{}
			err = t.Execute(&buff, ctx)
			if err != nil {
				fmt.Println("Error executing template: ", err)
				fmt.Println(buff.String())
			} else {
				fmt.Println("template: ", f.Name())
			}
		default:
			fmt.Println("did not match! ", f.Name())
		}
		cv3go.PrintToFile(buff.Bytes(), "test"+f.Name())
		buff.Reset()
	}
}
