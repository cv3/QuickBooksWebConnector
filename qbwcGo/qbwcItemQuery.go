package qbwcGo

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/Sirupsen/logrus"
)

//QBItemQuery readies a qbwc ItemQuery QBXML
func QBItemQuery() {
	Log.Debug("sending itemQueryRq")
	var tp = `./templates/itemQueryRequest.t`
	var x = bytes.Buffer{}
	var xx = bytes.Buffer{}
	var queryCTX = ItemQueryCTX{}
	//update items modified after the LastUpdate set in the config then update the configf
	if cfg.ItemUpdates.UpdateNewOnly {
		//parse time to make sure format is correct
		_, err := time.Parse(time.RFC3339, cfg.ItemUpdates.LastUpdate)
		if err != nil { //parse error
			Log.WithFields(logrus.Fields{
				"Error":      err,
				"LastUpdate": cfg.ItemUpdates.LastUpdate,
			}).Error("Error parsing time from config's ItemUpdates.LastUpdate")
			ErrLog.WithFields(logrus.Fields{
				"Error":      err,
				"LastUpdate": cfg.ItemUpdates.LastUpdate,
			}).Error("Error parsing time from config's ItemUpdates.LastUpdate")
		} else {
			queryCTX.FromModifiedDate = cfg.ItemUpdates.LastUpdate
		}
	}
	//Update the LastUpdate config var
	//cfg.ItemUpdates.LastUpdate = time.Now().Format(time.RFC3339)
	//SaveConfig()

	LoadTemplate(&tp, &queryCTX, &x)
	xml.Escape(&xx, x.Bytes())
	workChan <- WorkCTX{Work: xx.String(), Type: "ItemQueryRs"}
}
