package qbwcGo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

//DataExtAddQB will add a custom field value to a quickbooks object
func DataExtAddQB(TxnID string) string {
	var dataExts = LoadDataExtFile()
	var templateBuff = bytes.Buffer{}
	for _, dataExtAdd := range *dataExts {
		dataExtAdd.OwnerID = "0"
		dataExtAdd.TxnID = TxnID
		dataExtAdd.UseMacro = TxnID

		var tPath = `./templates/qbDataExtAdd.t`

		LoadTemplate(&tPath, dataExtAdd, &templateBuff)
	}
	return templateBuff.String()
}

//LoadDataExtFile will load the data ext set in the dataExtConfig file
func LoadDataExtFile() *[]DataExtAddRq {
	var dataExts = make([]DataExtAddRq, 0)
	var dataExtPath = "./config/qbwcDataExtConfig.json"
	dataExtFile, err := ioutil.ReadFile(dataExtPath)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error":       err,
			"dataExtPath": dataExtPath,
		}).Error("Error loading dataExt config file")
	}
	err = json.Unmarshal([]byte(dataExtFile), &dataExts)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error unmarshalling DataExt JSON")
	}
	return &dataExts
}
