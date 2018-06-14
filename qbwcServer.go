package main

import (
	"log"
	"net/http"

	"github.com/TeamFairmont/QuickBooksWebConnector/qbwcGo"
)

func main() {
	qbwcGo.Log = qbwcGo.InitLog()
	var cfg = qbwcGo.LoadConfig(qbwcGo.CfgPath)
	qbwcGo.Log = qbwcGo.SetLogLevel(cfg.Logging.Level, qbwcGo.Log)
	qbwcGo.ErrLog = qbwcGo.InitLog()
	qbwcGo.ErrLog = qbwcGo.SetLogLevel("error", qbwcGo.ErrLog)
	qbwcGo.ErrLog = qbwcGo.SetLogFile("./qbwcErrLog.log", qbwcGo.ErrLog)
	qbwcGo.Log = qbwcGo.SetLogFile(cfg.Logging.OutputPath, qbwcGo.Log)
	//add handlers to the default mux
	http.HandleFunc("/qbwc", qbwcGo.QBWCHandler)

	s := &http.Server{
		Addr:    cfg.ListenPort,
		Handler: nil, //nil to use default mux
	}

	log.Fatal(s.ListenAndServe())
}
