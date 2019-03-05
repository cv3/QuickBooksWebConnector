package qbwcGo

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/TeamFairmont/gabs"
	"github.com/amazingfly/cv3go"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

//CfgPath is the path to the config file, used first in qbwcServer.go then in SendAuthenticateResponse() to reload the config
var CfgPath = "./config/qbwcConfig.json"
var cfg = Config{} //holds config data
//Log is the normal logger
var Log *logrus.Logger

//ErrLog is specifically for logging error level events
var ErrLog *logrus.Logger
var globalSession = SessionCTX{}           //holds session data
var workChan = make(chan WorkCTX, 5)       //channel for work to be sent on
var workInsertChan = make(chan WorkCTX, 5) //channel for work to be inserted in the front of the line
//var insertWG = new(sync.WaitGroup)
var checkWorckChan = make(chan WorkCTX, 5)       //channel used in SendReceiveResponseXML, will hold the workCTX that has just been done.
var checkWorckInsertChan = make(chan WorkCTX, 5) //same as checkWorkChan, but for use when using workInsertChan
var doneChan = make(chan bool)                   //signals sent from CloseConnection so the order tracking go routine, it is done.
var getLastErrChan = make(chan string, 5)        //channel to send to the getLastError
var waiting bool                                 //used in the NoOp holding pattern

var shipToSuccessChan = make(chan ShipToSuccessTracker, 9999) //sent from GetCV3Orders, to start tracking a shipTo
var orderSuccessChan = make(chan OrderSuccessTracker, 9999)   //sent from GetCV3Orders, to start tracking an order
var confirmShipToChan = make(chan ShipToSuccessTracker, 9999) //sent from SendReceiveResponseXML to confirm a shipTo has been successfully added to QuickBooks

//QBWCHandler is the only handler, to handle qbwc soap requests by switching on xml node names
func QBWCHandler(w http.ResponseWriter, r *http.Request) {
	//ready the http header
	w.Header().Set("content-type", "text/xml")
	//read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error reading incoming message body")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error reading incoming message body")
		getLastErrChan <- err.Error()
	}
	var xmlNode = Node{}                //Struct to hold xml data,
	err = xml.Unmarshal(body, &xmlNode) // unmarshal xml into node struct
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling incoming message body")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling incoming message body")
		getLastErrChan <- err.Error()
	}
	//determine what the incoming call is by recursively check xml nodes
	CheckNodes([]Node{xmlNode}, Node{}, func(node, parentNode Node) bool {
		switch node.XMLName.Local {
		case "authenticate": //Authenticate the given credentials, and start the session
			//Initialize channels to make sure no old data remains
			InitChannels()
			SendAuthenticateResponse(parentNode, w)
			return true //end recursive CheckNode
		case "serverVersion": //Send this servers version information to QBWC
			SendServerVersionResponse(parentNode, w)
			return false //end recursive CheckNode
		case "clientVersion": //recieve QBWC's client version information
			SendClientVersionResponse(parentNode, w)
			return false //end recursive CheckNode
		case "closeConnection": //Close the connection to QBWC
			SendCloseConnectionResponse(parentNode, w)
			return false //end recursive CheckNode
		case "getLastError": //Send QBWC the last error this server has encountered
			SendGetLastErrorResponse(parentNode, w)
			return false //end recursive CheckNode
		case "sendRequestXML":
			SendSendRequestXMLResponse(parentNode, w)
			return false
		case "receiveResponseXML": //receive QBWC's response for the message sent in sendRequestXML
			SendReceiveResponseXMLResponse(parentNode, w)
			return false //end recursive CheckNode
		case "connectionError": //receive error message from QBWC if there was an error while it attempted to connect to QB
			SendConnectionErrorResponse(parentNode, w)
			return false //end recursive CheckNode
		default:
			return true //continue recursive CheckNode
		}
	})
}

//CheckNodes will recursively check xml nodes, and run the passed in funcion f
func CheckNodes(nodes []Node, parentNode Node, f func(Node, Node) bool) {
	for _, node := range nodes {
		if f(node, parentNode) { //if f returns true
			//pass current node's slice of child nodes, and itself as the new parent node
			CheckNodes(node.Nodes, node, f)
		}
	}
}

//SendReceiveResponseXMLResponse sends the response to QBWC's receiveResponseXML call
//this function receives the workCTX used in SendSendRequestXML on one of the checkWork channels.
//Then determines what type of response has been received by switching on the xml node name, then runs the corresponding function,
//Next determines if there is more work to be done, then sends its response
func SendReceiveResponseXMLResponse(parentNode Node, w http.ResponseWriter) {
	Log.Debug("Sending receiveResponseXML response")
	waiting = false       //we are no longer waiting
	var checkWork WorkCTX //workCTX to be used if errors occur
	select {              //check the insert channel first
	case cwi, ok := <-checkWorckInsertChan:
		if !ok {
			Log.WithFields(logrus.Fields{"ok": ok}).Error("check work insert chan error!")
			ErrLog.WithFields(logrus.Fields{"ok": ok}).Error("check work insert chan error!")
		}
		checkWork = cwi
		break //otherwise use the normal channel
	case cw, ok := <-checkWorckChan:
		if !ok {
			Log.WithFields(logrus.Fields{"ok": ok}).Error("check work chan error!")
			ErrLog.WithFields(logrus.Fields{"ok": ok}).Error("check work chan error!")
		}
		checkWork = cw
		break
	default: //never happens?
		Log.Error("checkWorkChan empty!")
		ErrLog.Error("checkWorkChan empty!")
	}
	//Unmarshal the content of the incoming xml into the proper struct
	var receiveResponseXMLCTX = ReceiveResponseXMLCTX{}
	err := xml.Unmarshal(parentNode.Content, &receiveResponseXMLCTX)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "parent node": parentNode.XMLName}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse()")
		ErrLog.WithFields(logrus.Fields{"error": err, "parent node": parentNode.XMLName}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse()")
		getLastErrChan <- err.Error()
	} else { //no error unmarshaling xml
		//Unmarshal the QBXML response, in the response field of the incoming xml
		var responseNode = Node{}
		err = xml.Unmarshal(receiveResponseXMLCTX.Response, &responseNode)
		if err != nil {
			Log.WithFields(logrus.Fields{"error": err, "node": string(receiveResponseXMLCTX.Response)}).Error("Error unmarshalling receiveResponseXMLCTX.Response")
			ErrLog.WithFields(logrus.Fields{"error": err, "node": string(receiveResponseXMLCTX.Response)}).Error("Error unmarshalling receiveResponseXMLCTX.Response")
		}
		//Find the ItemQueryRs node
		CheckNodes([]Node{responseNode}, Node{}, func(node, parentNode Node) bool {
			switch node.XMLName.Local {
			case "SalesOrderAddRs":
				SalesOrderAddRsHandler(parentNode, checkWork)
				return false // end recursive CheckNode
			case "SalesReceiptAddRs":
				SalesReceiptAddRsHandler(parentNode, checkWork)
				return false // end recursive CheckNode
			//IteemInventoryAddRs found
			case "ItemInventoryAddRs":
				ItemInventoryAddRsHandler(parentNode, checkWork)
				return false //end recursive checknode
			//not currently used
			//IteemGroupAddRs found
			case "ItemGroupAddRs":
				ItemGroupAddRsHandler(parentNode, checkWork)
				return false //end recursive checkNode
			//ItemQueryRs found, send the items to CV3
			case "ItemQueryRs":
				ItemQueryRsHandler(parentNode, checkWork)
				return false //end recursive CheckNode
			case "CustomerAddRs":
				CustomerAddRsHandler(parentNode, checkWork)
				return false
			} //end main switch
			return true //continue recursive checkNode
		}) //end check node
		//Load the response template, and write the response to w
		var b = bytes.Buffer{}
		var ctx = ReceiveResponseXMLResponseCTX{}
		var tPath = `./templates/receiveResponseXMLResponse.t`
		//if the session has more work to do
		if len(workChan) > 0 || len(workInsertChan) > 0 {
			ctx.Complete = 50
			Log.WithFields(logrus.Fields{"workInsertChan lenght": len(workInsertChan), "workChan length": len(workChan)}).Debug("More work remains")
		} else { //if no more work, wait and check agian, needed to make sure late work is handled
			time.Sleep(10 * time.Second)
			if len(workChan) > 0 || len(workInsertChan) > 0 {
				ctx.Complete = 50
				Log.WithFields(logrus.Fields{"workInsertChan lenght": len(workInsertChan), "workChan length": len(workChan)}).Debug("More work remains")
			} else { //done
				ctx.Complete = 100
				Log.WithFields(logrus.Fields{"workInsertChan lenght": len(workInsertChan), "workChan length": len(workChan)}).Debug("No more work, DONE")
			}
		}
		LoadTemplate(&tPath, ctx, &b)
		w.Write(b.Bytes())
	}
}

//SendSendRequestXMLResponse sends the response to SendRequestXLM call from QBWC.
//This function first checks if there is work ready to be done.  Adds this work to the queue and sends its workCTX on the checkWork channel to be used in SendReceiveResponseXMLResponse
//Then sends the desired work to be done
func SendSendRequestXMLResponse(parentNode Node, w http.ResponseWriter) {
	Log.Debug("sending sendRequestXML response")
	//Unmarshal the content of the incoming xml into the proper struct
	err := xml.Unmarshal(parentNode.Content, &globalSession)
	if err != nil { //Unmarshal error
		Log.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("error unmarshalling sendRequestXML")
		ErrLog.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("error unmarshalling sendRequestXML")
		getLastErrChan <- err.Error()
	} else { //no error unmarshaling xml
		var work = WorkCTX{}
		//if an insert job is still in progress, wait until it finishes.
		//insertWG.Wait()
		select {
		//check the workInsert channel first
		case work, ok := <-workInsertChan:
			if ok {
				if work.Attempted < cfg.MaxWorkAttempts {
					work.Attempted++
					//Need to track the work being
					checkWorckInsertChan <- work
					//add work to the Queue
					globalSession.QBXMLWorkQueue = append(globalSession.QBXMLWorkQueue, work.Work)
				} else { //too many retires on this workCTX
					Log.WithFields(logrus.Fields{"Attempts": work.Attempted}).Error("This work has been retried too many times")
					ErrLog.WithFields(logrus.Fields{"Attempts": work.Attempted}).Error("This work has been retried too many times")
					SendSendRequestXMLResponse(parentNode, w)
				}
			}
			break
		//check the normal work channel if the insert channel is empty
		case work, ok := <-workChan:
			if ok {
				if work.Attempted < cfg.MaxWorkAttempts {
					if work.Work == "" { //No work, ask QB to pause for 5 seconds
						getLastErrChan <- "NoOp"
						Log.Debug("sending NoOp on getLastErrChan")
					} else { //Need to track the work being done
						work.Attempted++
						checkWorckChan <- work
					}
					//add work to the work queue
					globalSession.QBXMLWorkQueue = append(globalSession.QBXMLWorkQueue, work.Work)
				} else { //too many retires on this workCTX
					Log.WithFields(logrus.Fields{"Attempts": work.Attempted}).Error("This work has been retried too many times")
					ErrLog.WithFields(logrus.Fields{"Attempts": work.Attempted}).Error("This work has been retried too many times")
					SendSendRequestXMLResponse(parentNode, w)
				}
			}
			break
		default: //happens during NoOp
			if waiting { // continue NoOp holding pattern
				Log.Debug("sending NoOp on getLastErrChan")
				globalSession.QBXMLWorkQueue = append(globalSession.QBXMLWorkQueue, "")
				getLastErrChan <- "NoOp"
			} else {
				Log.Error("not waiting and could not read from workChan, in sendRequestXML")
				ErrLog.Error("not waiting and could not read from workChan, in sendRequestXML")
			}
		}
		//Load the desired QBXML template, and fill its data to be added to the response field of the response
		var tPath = `./templates/sendRequestXMLResponse.t`
		var ctx = SendRequestXMLResponseCTX{}
		var b = bytes.Buffer{}

		//if there is work in the work queue, add it to the sendRequest context
		if len(globalSession.QBXMLWorkQueue) > 0 {
			ctx.QBXML = globalSession.QBXMLWorkQueue[0]
			globalSession.QBXMLWorkQueue = append(globalSession.QBXMLWorkQueue[:0], globalSession.QBXMLWorkQueue[1:]...)
			Log.WithFields(logrus.Fields{"request type": work.Type}).Info("Sending request to Quick Books")
		}
		//Load the soap response template, with the response field set, then write to w
		LoadTemplate(&tPath, ctx, &b)
		w.Write(b.Bytes())
	}
}

//SendConnectionErrorResponse sends the response to QBWC's connectionError call, when the QuickBooks Web Connector has an issue connecting to QuickBooks, and logs the error
func SendConnectionErrorResponse(parentNode Node, w http.ResponseWriter) {
	Log.Debug("sending connectionError response")
	//Unmarshal the content of the incoming xml into the proper struct
	var connectionErrorCTX = ConnectionErrorCTX{}
	err := xml.Unmarshal(parentNode.Content, &connectionErrorCTX)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("Error unmarshalling connectionError")
		ErrLog.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("Error unmarshalling connectionError")
		getLastErrChan <- err.Error()
	} else { //no error unmarshaling xml
		Log.WithFields(logrus.Fields{
			"Message": connectionErrorCTX.Message,
			"Ticket":  connectionErrorCTX.Ticket,
			"HResult": connectionErrorCTX.HResult,
		}).Error("ConnectionError, connection from QuickBooks Web Connector to QuickBooks Desktop")
		ErrLog.WithFields(logrus.Fields{
			"Message": connectionErrorCTX.Message,
			"Ticket":  connectionErrorCTX.Ticket,
			"HResult": connectionErrorCTX.HResult,
		}).Error("ConnectionError, connection from QuickBooks Web Connector to QuickBooks Desktop")
		//Load the soap response template and write to w
		var b = bytes.Buffer{}
		var ctx = ConnectionErrorResponseCTX{}
		var tPath = `./templates/connectionErrorResponse.t`
		LoadTemplate(&tPath, ctx, &b)
		w.Write(b.Bytes())
	}
}

//SendServerVersionResponse sends the server version information to QBWC
func SendServerVersionResponse(parentNode Node, w http.ResponseWriter) {
	Log.Debug("sending server version")
	//Load the soap response template and write to w
	var b = bytes.Buffer{}
	var ctx = ServerVersionResponseCTX{}
	var tPath = `./templates/serverVersionResponse.t`
	//add the server version information
	ctx.Version = cfg.ServerVersion
	LoadTemplate(&tPath, ctx, &b)
	w.Write(b.Bytes())
}

//SendClientVersionResponse sends the response to the client version call
func SendClientVersionResponse(parentNode Node, w http.ResponseWriter) {
	var clientVersionCTX = ClientVersionCTX{}
	err := xml.Unmarshal(parentNode.Content, &clientVersionCTX)
	if err != nil {
		Log.Error("Error unmarshalling ClientVersion soap request")
		ErrLog.Error("Error unmarshalling ClientVersion soap request")
		getLastErrChan <- err.Error()
	} else {
		Log.WithFields(logrus.Fields{"Version in config": cfg.QBWCVersion, "Version Received": clientVersionCTX.StrVersion}).Debug("Client version in")
	}
	//Load the soap response template and write to w
	var b = bytes.Buffer{}
	var ctx = ClientVersionResponseCTX{}
	ctx.Result = cfg.QBWCVersion
	var tPath = `./templates/clientVersionResponse.t`
	LoadTemplate(&tPath, ctx, &b)
	w.Write(b.Bytes())
}

//SendCloseConnectionResponse sends the response to the close connection call.
//This function will signal the order success tracker to end, this will be the last call from QBWC
func SendCloseConnectionResponse(parentNode Node, w http.ResponseWriter) {
	var closeConnectionCTX = CloseConnectionCTX{}
	err := xml.Unmarshal(parentNode.Content, &closeConnectionCTX)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling CloseConnection")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling CloseConnection")
		getLastErrChan <- err.Error()
	}
	go func() { //send done signal to order tracker
		doneChan <- true
	}()
	Log.WithFields(logrus.Fields{"ticket": closeConnectionCTX.Ticket}).Info("Connection closed")
	var b = bytes.Buffer{}
	var ctx = CloseConnectionResponseCTX{}
	ctx.CloseConnectionResult = cfg.CloseConnectionMessage
	var tPath = `./templates/closeConnectionResponse.t`
	LoadTemplate(&tPath, ctx, &b)
	w.Write(b.Bytes())
}

//SendGetLastErrorResponse sends the response to getLastError call.  This call is send when QBWC thinks this server has an error, or is in the NoOp holding pattern
func SendGetLastErrorResponse(parentNode Node, w http.ResponseWriter) {
	var getLastErrorCTX = GetLastErrorCTX{}
	var b = bytes.Buffer{}
	var ctx = GetLastErrorResponseCTX{}
	var tPath = `./templates/getLastErrorResponse.t`
	err := xml.Unmarshal(parentNode.Content, &getLastErrorCTX)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling getLastError")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling getLastError")
		getLastErrChan <- err.Error()
	} else {
		Log.WithFields(logrus.Fields{"ticket": getLastErrorCTX.Ticket}).Debug("send get last error response")

		select { //get the last error from getLastErrChan
		case lastErr, ok := <-getLastErrChan:
			if ok {
				Log.WithFields(logrus.Fields{"error": lastErr}).Debug("Getting last error from GetLastErrChan")
				ctx.LastError = lastErr
			}
		default:
			Log.Debug("nothing on getLastErrChan, sending blank")
			ctx.LastError = ""
		}
	}
	LoadTemplate(&tPath, ctx, &b)
	w.Write(b.Bytes())
}

//SendAuthenticateResponse will send the response to authenticate
//This function will authenticate the credentials sent by QBWC.
//Start storing session information, the NoOp holding pattern, the order success tracker, and start "getting" work
func SendAuthenticateResponse(parentNode Node, w http.ResponseWriter) {
	//reload config
	LoadConfig(CfgPath)
	//Seed random int for ticken / session id
	rand.Seed(time.Now().UTC().UnixNano())
	var tPath = `./templates/authenticateResponse.t`
	var b = bytes.Buffer{}
	var ctx = AuthenticateResponseCTX{}
	Log.Debug("sending authentication response")
	//unmarshal the parent node that contains the authentication information
	var authCTX = AuthenticateCTX{}
	err := xml.Unmarshal(parentNode.Content, &authCTX)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("Error unmarshalling authenticate")
		ErrLog.WithFields(logrus.Fields{"error": err, "node": parentNode.XMLName}).Error("Error unmarshalling authenticate")
		getLastErrChan <- err.Error()
	} else { //no error unmarshaling xml
		//user verification
		if authCTX.UserName == cfg.QBWCCredentials.User {
			if authCTX.Password == cfg.QBWCCredentials.Pass {
				//set session ID
				ctx.Ticket = strconv.Itoa(1000 + rand.Intn(1000))
				globalSession = SessionCTX{Ticket: ctx.Ticket}

				//Start holding patter until work is ready
				waiting = true
				go func() {
					workChan <- WorkCTX{Work: ""}
				}()

				//start keeping track of a new order or ship to, as each shipto is a seperate quickbooks salesReceiptAdd
				go StartOrderTracker()
				//go get the work
				go InitWork()
			} else { //invalid password
				Log.Error("Invalid user or password")
				ErrLog.Error("Invalid user or password")
				getLastErrChan <- "Invalid user or password"
			} //end password verification
		} else { //invalid username
			Log.Error("Invalid user or password")
			ErrLog.Error("Invalid user or password")
			getLastErrChan <- "Invalid user or password"
		} //end user verification
	}
	LoadTemplate(&tPath, ctx, &b)
	w.Write(b.Bytes())
}

// LoadTemplate accepts a path to a template, a struct for the temmplates context, and the request body to fill
func LoadTemplate(tPath *string, ctx interface{}, requestBody *bytes.Buffer) {
	t, err := template.ParseFiles(*tPath)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "filepath": *tPath}).Error("Error parsing template file")
		ErrLog.WithFields(logrus.Fields{"error": err, "filepath": *tPath}).Error("Error parsing template file")
		getLastErrChan <- err.Error()
	} // Populate requestBody with the executed template and context

	err = t.Execute(requestBody, ctx)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "filepath": *tPath}).Error("Error executing template")
		ErrLog.WithFields(logrus.Fields{"error": err, "filepath": *tPath}).Error("Error executing template")
		getLastErrChan <- err.Error()
	}
}

//InitWork will go get work to be done
func InitWork() {
	//Check config, if special Quick Books items need to be initialized
	if cfg.InitQBItems {
		go InitQBItems()
		cfg.InitQBItems = false
		SaveConfig()
		Log.Info("Initializing Quick Books Itesm")
	}
	if cfg.ItemUpdates.UpdateCV3Items {
		go QBItemQuery()
		Log.WithFields(logrus.Fields{"UpdateCV3Items": cfg.ItemUpdates.UpdateCV3Items}).Info("Updating CV3 items from Quick Books")
	} //else config set to not update cv3 items from QB
	go GetCV3Orders()
	//go AddTestItemsToQB()
}

//AddTestItemsToQB will add test items to the QB DB
func AddTestItemsToQB() {
	var products = cv3go.Products{}

	for i := 1; i < 40000; i++ {
		var temp = cv3go.Product{}
		temp.Name = "massive test item " + strconv.Itoa(i)
		temp.Description = "the best " + temp.Name
		temp.Retail.Price.StandardPrice = strconv.Itoa(i + i)
		temp.Retail.Price.PriceCategory = "Retail"
		temp.InventoryControl.InventoryOnHand = strconv.Itoa(i + i)
		temp.Categories.IDs = []string{"8"}

		products.Products = append(products.Products, temp)
	}

	for _, item := range products.Products {
		ImportCV3ItemsToQB(item)
	}
}

//InitLog initialized the log
func InitLog() *logrus.Logger {
	var l = logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.JSONFormatter), //&logrus.TextFormatter{DisableTimestamp: true}, //new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
	return &l
}

//SetLogFile sets the file location to store the log
func SetLogFile(path string, l *logrus.Logger) *logrus.Logger {
	lumberjackLogrotate := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    50,
		MaxBackups: 0,
		MaxAge:     0,
		Compress:   true,
	}
	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	l.Out = logMultiWriter
	/*file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		l.Out = file
		l.WithFields(logrus.Fields{"filepath": path}).Debug("logger set to use file")
	} else {
		Log.WithFields(logrus.Fields{"error": err, "filepath": path}).Error("Failed to log to file")
		ErrLog.WithFields(logrus.Fields{"error": err, "filepath": path}).Error("Failed to log to file")
		getLastErrChan <- err.Error()
	}*/
	return l
}

//SetLogLevel sets the logging level
func SetLogLevel(level string, l *logrus.Logger) *logrus.Logger {
	levelMap := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"warn":  logrus.WarnLevel,
		"panic": logrus.PanicLevel,
		"info":  logrus.InfoLevel,
	}
	_, ok := levelMap[strings.ToLower(level)]
	if ok {
		l.Level = levelMap[strings.ToLower(level)]
		l.WithFields(logrus.Fields{"level": strings.ToLower(level)}).Info("logging level set")
	} else {
		Log.WithFields(logrus.Fields{"level": strings.ToLower(level)}).Error("Error setting log level")
		ErrLog.WithFields(logrus.Fields{"level": strings.ToLower(level)}).Error("Error setting log level")
	}
	return l
}

//LoadConfig loads the configuration file
func LoadConfig(cfgPath string) *Config {
	// load the config file
	configFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error":   err,
			"cfgPath": cfgPath,
		}).Fatal("Error loading config file")
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(configFile), &cfg)
	if err != nil {
		getLastErrChan <- err.Error()
		Log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Error unmarshalling config JSON")
		os.Exit(1)
	}
	Log.WithFields(logrus.Fields{
		"cfgPath":         cfgPath,
		"listenPort":      cfg.ListenPort,
		"UpdateCV3Items":  cfg.ItemUpdates.UpdateCV3Items,
		"log level":       cfg.Logging.Level,
		"log output path": cfg.Logging.OutputPath,
	}).Debug("Loaded config")
	return &cfg
}

//SaveConfig saves the config
func SaveConfig() {
	var out bytes.Buffer
	//Create json from the loaded config
	b, err := json.MarshalIndent(&cfg, "", "	")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error marshaling config")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error marshaling config")
		getLastErrChan <- err.Error()
	} //indent the json, and store it in out
	err = json.Indent(&out, b, "", "\t")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error indenting config in SaveConfig")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error indenting config in SaveConfig")
		getLastErrChan <- err.Error()
	}
	err = ioutil.WriteFile(`./config/qbwcConfig.json`, out.Bytes(), 0755)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "filePath": `./config/qbwcConfig.json`}).Error("Error writing config file")
		ErrLog.WithFields(logrus.Fields{"error": err, "filePath": `./config/qbwcConfig.json`}).Error("Error writing config file")
		getLastErrChan <- err.Error()
	}
}

//StartOrderTracker starts the order confirmation tracking
func StartOrderTracker() {
	var orderTrackerMap = make(map[string]*OrderSuccessTracker)
	for {
		select { //receiving order, start tracking order
		case orderTracker := <-orderSuccessChan:
			if _, ok := orderTrackerMap[orderTracker.OrderID]; !ok {
				orderTracker.ShipToSuccessTrackers = make(map[int]ShipToSuccessTracker)
				Log.WithFields(logrus.Fields{"orderID": orderTracker.OrderID, "shipTo length": orderTracker.ShipToLength}).Debug("order tracker in")
				orderTrackerMap[orderTracker.OrderID] = &OrderSuccessTracker{OrderID: orderTracker.OrderID, ShipToLength: orderTracker.ShipToLength, ShipToSuccessTrackers: orderTracker.ShipToSuccessTrackers, SuccessCount: orderTracker.SuccessCount} //orderTracker
			} else {
				Log.WithFields(logrus.Fields{"orderID": orderTracker.OrderID}).Error("Order tracker already exists")
				ErrLog.WithFields(logrus.Fields{"orderID": orderTracker.OrderID}).Error("Order tracker already exists")
			}
			break
		//Receiving shipTo, add shipTo to the order
		case shipToTracker := <-shipToSuccessChan:
			Log.WithFields(logrus.Fields{"orderID": shipToTracker.OrderID, "shipToIndex": shipToTracker.Index}).Debug("shipToTracker in")
			if _, ok := orderTrackerMap[shipToTracker.OrderID]; ok {
				if _, ok := orderTrackerMap[shipToTracker.OrderID].ShipToSuccessTrackers[shipToTracker.Index]; !ok {
					orderTrackerMap[shipToTracker.OrderID].ShipToSuccessTrackers[shipToTracker.Index] = ShipToSuccessTracker{Success: shipToTracker.Success, QBErrorMessage: shipToTracker.QBErrorMessage, QBErrorCode: shipToTracker.QBErrorCode, OrderID: shipToTracker.OrderID, Index: shipToTracker.Index}
				} else {
					Log.WithFields(logrus.Fields{"orderID": shipToTracker.OrderID, "shipToIndex": shipToTracker.Index}).Error("ShipTo tracker already exists")
					ErrLog.WithFields(logrus.Fields{"orderID": shipToTracker.OrderID, "shipToIndex": shipToTracker.Index}).Error("ShipTo tracker already exists")
				}
			} else { //order tracker for this shipTo does not exist, resend
				Log.WithFields(logrus.Fields{"orderID": shipToTracker.OrderID, "shipToIndex": shipToTracker.Index}).Error("Order tracker for this shipTo does not exist")
				ErrLog.WithFields(logrus.Fields{"orderID": shipToTracker.OrderID, "shipToIndex": shipToTracker.Index}).Error("Order tracker for this shipTo does not exist")
				go func(shipToTracker ShipToSuccessTracker) {
					shipToSuccessChan <- shipToTracker
				}(shipToTracker)
			}
			break
		//receive shipTo success confirmation
		case confirm := <-confirmShipToChan:
			//Overwrite the existing shipTo with the confirmation shipTo
			if _, ok := orderTrackerMap[confirm.OrderID]; ok {
				if _, ok := orderTrackerMap[confirm.OrderID].ShipToSuccessTrackers[confirm.Index]; ok {
					orderTrackerMap[confirm.OrderID].ShipToSuccessTrackers[confirm.Index] = ShipToSuccessTracker{OrderID: confirm.OrderID, Index: confirm.Index, QBErrorCode: confirm.QBErrorCode, QBErrorMessage: confirm.QBErrorMessage, Success: confirm.Success} //confirm

					orderTrackerMap[confirm.OrderID].SuccessCount = 0
					for _, shipTo := range orderTrackerMap[confirm.OrderID].ShipToSuccessTrackers {
						if shipTo.Success { //if successful
							orderTrackerMap[confirm.OrderID].SuccessCount++
						}
					} //check each orderTrackers successCount compated to shipToLength
					//if the orderTrackers shipToLength is equal to the count, they have all been successful. Send order confirmation to CV3 and delete the orderTracker
					if orderTrackerMap[confirm.OrderID].ShipToLength == orderTrackerMap[confirm.OrderID].SuccessCount && orderTrackerMap[confirm.OrderID].SuccessCount > 0 {
						Log.WithFields(logrus.Fields{"orderID": confirm.OrderID}).Debug("Confirming order")
						go ConfirmCV3Order(confirm.OrderID)      //send confirmation to CV3
						delete(orderTrackerMap, confirm.OrderID) //remove order tracker from the map
					}
				} else {
					Log.WithFields(logrus.Fields{"orderID": confirm.OrderID, "shipToIndex": confirm.Index}).Error("ShipTo does not exist, it cannot be confirmed")
					ErrLog.WithFields(logrus.Fields{"orderID": confirm.OrderID, "shipToIndex": confirm.Index}).Error("ShipTo does not exist, it cannot be confirmed")
				}
			} else { //else shipTo not in map, resend
				Log.WithFields(logrus.Fields{"orderID": confirm.OrderID, "shipToIndex": confirm.Index}).Error("ShipTo confirmation does not exist")
				ErrLog.WithFields(logrus.Fields{"orderID": confirm.OrderID, "shipToIndex": confirm.Index}).Error("ShipTo confirmation does not exist")
				go func(confirm ShipToSuccessTracker) {
					confirmShipToChan <- confirm
				}(confirm)
			}
			break
		//Session is finished, log unsuccsessful orders.
		case <-doneChan:
			for id, orderSuccessTracker := range orderTrackerMap {
				var errBuff = bytes.Buffer{}
				//compile a error report for the order's shipTos
				for _, shipToSuccessTracker := range orderSuccessTracker.ShipToSuccessTrackers {
					if shipToSuccessTracker.QBErrorCode != "" {
						errBuff.WriteString("shipToIndex: ")
						errBuff.WriteString(strconv.Itoa(shipToSuccessTracker.Index))
						errBuff.WriteString(" QBErrorCode: ")
						errBuff.WriteString(shipToSuccessTracker.QBErrorCode)
						errBuff.WriteString(" QBErrorMessage: ")
						errBuff.WriteString(shipToSuccessTracker.QBErrorMessage)
						errBuff.WriteString(" || ")
					}
				}

				Log.WithFields(logrus.Fields{
					"cv3OrderID":         id,
					"successful shipTos": orderSuccessTracker.SuccessCount,
					"total shipTos":      orderSuccessTracker.ShipToLength,
					"order error report": errBuff.String(),
				}).Error("Quick Books did not import the order successfully")
				ErrLog.WithFields(logrus.Fields{
					"cv3OrderID":         id,
					"successful shipTos": orderSuccessTracker.SuccessCount,
					"total shipTos":      orderSuccessTracker.ShipToLength,
					"order error report": errBuff.String(),
				}).Error("Quick Books did not import the order successfully")
			}
			return //end go routine
		}
	}
}

//SalesOrderAddRsHandler is used in SendResponseXML to handle a SalesOrderAddRs
func SalesOrderAddRsHandler(parentNode Node, checkWork WorkCTX) {
	var salesOrderAddRs = SalesOrderAddRs{}
	Log.Debug("SalesOrderAddRs in")
	err := xml.Unmarshal(parentNode.Content, &salesOrderAddRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err, "node": string(parentNode.Content)}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse() switch case SalesOrderAddRs")
		ErrLog.WithFields(logrus.Fields{"error": err, "node": string(parentNode.Content)}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse() switch case SalesOrderAddRs")
		getLastErrChan <- err.Error()
	} else { //else salesOrderAddRs unmarshalled
		switch salesOrderAddRs.StatusCode {
		case "0": //no error
			Log.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Info("SalesOrderAdd successful, sending confirmation")
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesOrderAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesOrderAdd).CV3OrderID,
					Success:        true,
					QBErrorCode:    salesOrderAddRs.StatusCode,
					QBErrorMessage: salesOrderAddRs.StatusMessage,
				}
			}()
			break
		case "3140": //item not in QB
			Log.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("SalesOrderAddRs error 3140, items not in Quick Books")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("SalesOrderAddRs error 3140, items not in Quick Books")

			//Check error message to see what feild had errors
			switch { //check for CustomerMsg first to avoid a false positive with Customer
			case strings.Contains(salesOrderAddRs.StatusMessage, "CustomerMsg"):
				Log.Debug("CustomerMsg not found, attempting to add new customerMsg to Quickbooks")
				go CustomerMsgAddQB(checkWork)
				break
			case strings.Contains(salesOrderAddRs.StatusMessage, "Customer"):
				Log.Debug("Customer not found, attempting to add new customer to Quickbooks")
				if cfg.MaxWorkAttempts-checkWork.Attempted > 0 {
					go CustomerAddQB(checkWork)
				} else {
					Log.WithFields(logrus.Fields{"OrderID": CheckPath("orderID", checkWork.Order)}).Error("Maximum number of attempts to add this salesReceipt has been exceeded")
					ErrLog.WithFields(logrus.Fields{"OrderID": CheckPath("orderID", checkWork.Order)}).Error("Maximum number of attempts to add this salesReceipt has been exceeded")
				}
				break
			}
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesOrderAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesOrderAdd).CV3OrderID,
					Success:        false,
					QBErrorCode:    salesOrderAddRs.StatusCode,
					QBErrorMessage: salesOrderAddRs.StatusMessage,
				}
			}() //if set to true in the config, this will automatically import order items from an order that fails to import to QuickBooks because the items do not exist in quick books
			go AutoImportItemsToQB(checkWork)
			break
		case "3180": //occurs when quickbooks thinks the list is being accessed from another location, may only happen in Enterprise version.  The solution is to simply resent the call
			Log.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs, Resending")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs, Resending")
			workChan <- checkWork
			break
		case "3270": //occurs when using SalesOrderAdds with a quick books version that does not support it
			Log.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs, you may not have a feature you are trying to use. E.G.  If you are trying to use SalesOrders and do are using QuickBooks Pro")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs, you may not have a feature you are trying to use. E.G.  If you are trying to use SalesOrders and do are using QuickBooks Pro")
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesOrderAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesOrderAdd).CV3OrderID,
					Success:        false,
					QBErrorCode:    salesOrderAddRs.StatusCode,
					QBErrorMessage: salesOrderAddRs.StatusMessage,
				}
			}()
			break
		default:
			Log.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesOrderAddRs.StatusSeverity,
				"message":         salesOrderAddRs.StatusMessage,
				"Status Code":     salesOrderAddRs.StatusCode,
			}).Error("Error in SalesOrderAddRs")
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesOrderAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesOrderAdd).CV3OrderID,
					Success:        false,
					QBErrorCode:    salesOrderAddRs.StatusCode,
					QBErrorMessage: salesOrderAddRs.StatusMessage,
				}
			}()
		}
	}
}

//SalesReceiptAddRsHandler is used in SendResponseXML to handle a SalesReceiptAddRs
func SalesReceiptAddRsHandler(parentNode Node, checkWork WorkCTX) {
	var salesReceiptAddRs = SalesReceiptAddRs{}
	Log.Debug("SalesReceiptAddRs in")
	err := xml.Unmarshal(parentNode.Content, &salesReceiptAddRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse() switch case SalesReceiptAddRs")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling xml node in SendReceiveResponseXMLResponse() switch case SalesReceiptAddRs")
		getLastErrChan <- err.Error()
	} else { //else salesRecieptAddRs unmarshalled
		switch salesReceiptAddRs.StatusCode {
		case "0": //no error
			Log.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Info("SalesReceiptAdd successful, sending confirmation")
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesReceiptAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesReceiptAdd).CV3OrderID,
					Success:        true,
					QBErrorCode:    salesReceiptAddRs.StatusCode,
					QBErrorMessage: salesReceiptAddRs.StatusMessage,
				}
			}()
			break
		case "3140": //Account Ref not in QB
			Log.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("SalesReceiptAddRs error 3140, Account Ref not in Quick Books")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("SalesReceiptAddRs error 3140, Account Ref not in Quick Books")

			//Check error message to see what feild had errors
			switch { //check for CustomerMsg first to avoid a false positive with Customer
			case strings.Contains(salesReceiptAddRs.StatusMessage, "CustomerMsg"):
				Log.Debug("CustomerMsg not found, attempting to add new customerMsg to Quickbooks")
				go CustomerMsgAddQB(checkWork)
				break
			case strings.Contains(salesReceiptAddRs.StatusMessage, "Customer"):
				Log.Debug("Customer not found, attempting to add new customer to Quickbooks")
				if cfg.MaxWorkAttempts-checkWork.Attempted > 0 {
					go CustomerAddQB(checkWork)
				} else {
					Log.WithFields(logrus.Fields{"OrderID": CheckPath("orderID", checkWork.Order)}).Error("Maximum number of attempts to add this salesReceipt has been exceeded")
					ErrLog.WithFields(logrus.Fields{"OrderID": CheckPath("orderID", checkWork.Order)}).Error("Maximum number of attempts to add this salesReceipt has been exceeded")
				}
				break
			}
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesReceiptAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesReceiptAdd).CV3OrderID,
					Success:        false,
					QBErrorCode:    salesReceiptAddRs.StatusCode,
					QBErrorMessage: salesReceiptAddRs.StatusMessage,
				}
			}() //if set to true in the config, this will automatically import order items from an order that fails to import to QuickBooks because the items do not exist in quick books
			go AutoImportItemsToQB(checkWork)
			break
		case "3180": //occurs when quickbooks thinks the list is being accessed from another location, may only happen in Enterprise version.  The solution is to simply resent the call
			Log.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("Error in SalesReceiptAddRs, Resending")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("Error in SalesReceiptAddRs, Resending")
			workChan <- checkWork
			break
		default:
			Log.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("Error in SalesReceiptAddRs")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": salesReceiptAddRs.StatusSeverity,
				"message":         salesReceiptAddRs.StatusMessage,
				"Status Code":     salesReceiptAddRs.StatusCode,
			}).Error("Error in SalesReceiptAddRs")
			go func() {
				confirmShipToChan <- ShipToSuccessTracker{
					Index:          checkWork.Data.(SalesReceiptAdd).ShipToIndex,
					OrderID:        checkWork.Data.(SalesReceiptAdd).CV3OrderID,
					Success:        false,
					QBErrorCode:    salesReceiptAddRs.StatusCode,
					QBErrorMessage: salesReceiptAddRs.StatusMessage,
				}
			}()
		}
	}
}

//ItemInventoryAddRsHandler is used in SendResponseXML to handle a ItemInventoryAddRs
func ItemInventoryAddRsHandler(parentNode Node, checkWork WorkCTX) {
	var itemInventoryAddRs = ItemInventoryAddRs{}
	err := xml.Unmarshal(parentNode.Content, &itemInventoryAddRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("error unmarshalling ItemInventoryAddRs")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("error unmarshalling ItemInventoryAddRs")
		getLastErrChan <- err.Error()
	} else {
		switch itemInventoryAddRs.StatusCode {
		case "0":
			Log.WithFields(logrus.Fields{
				"Status Severity": itemInventoryAddRs.StatusSeverity,
				"message":         itemInventoryAddRs.StatusMessage,
				"Status Code":     itemInventoryAddRs.StatusCode,
			}).Info("itemInventoryAddRs in")
			break
		case "3100":
			Log.WithFields(logrus.Fields{
				"Status Severity": itemInventoryAddRs.StatusSeverity,
				"message":         itemInventoryAddRs.StatusMessage,
				"Status Code":     itemInventoryAddRs.StatusCode,
			}).Error("Erro in itemInventoryAddRs, item already in QuickBooks")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": itemInventoryAddRs.StatusSeverity,
				"message":         itemInventoryAddRs.StatusMessage,
				"Status Code":     itemInventoryAddRs.StatusCode,
			}).Error("Erro in itemInventoryAddRs, item already in QuickBooks")
			break
		default:
			Log.WithFields(logrus.Fields{
				"Status Severity": itemInventoryAddRs.StatusSeverity,
				"message":         itemInventoryAddRs.StatusMessage,
				"Status Code":     itemInventoryAddRs.StatusCode,
			}).Error("Error in itemInventoryAddRs")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": itemInventoryAddRs.StatusSeverity,
				"message":         itemInventoryAddRs.StatusMessage,
				"Status Code":     itemInventoryAddRs.StatusCode,
			}).Error("Error in itemInventoryAddRs")
			break
		}
	}
}

//ItemGroupAddRsHandler asdf
func ItemGroupAddRsHandler(parentNode Node, checkWork WorkCTX) {
	var itemGroupAddRs = ItemGroupAddRs{}
	err := xml.Unmarshal(parentNode.Content, &itemGroupAddRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling ItemGroupAddRs")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling ItemGroupAddRs")
		getLastErrChan <- err.Error()
	} else {
		if itemGroupAddRs.StatusCode != "0" && itemGroupAddRs.StatusCode != "3100" {
			Log.WithFields(logrus.Fields{
				"Status Severity": itemGroupAddRs.StatusSeverity,
				"message":         itemGroupAddRs.StatusMessage,
				"Status Code":     itemGroupAddRs.StatusCode,
			}).Debug(":::ItemGroupAddRs has an error!!!")
		} else {
			Log.WithFields(logrus.Fields{
				"Status Severity": itemGroupAddRs.StatusSeverity,
				"message":         itemGroupAddRs.StatusMessage,
				"Status Code":     itemGroupAddRs.StatusCode,
			}).Error("itemGroupAddRs in")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": itemGroupAddRs.StatusSeverity,
				"message":         itemGroupAddRs.StatusMessage,
				"Status Code":     itemGroupAddRs.StatusCode,
			}).Error("itemGroupAddRs in")
		}
	}
}

//CustomerAddRsHandler will handle the response from a CustomerAddRq
func CustomerAddRsHandler(parentNode Node, checkWork WorkCTX) {
	var customerAddRs = CustomerAddRs{}
	//Unmarshal the CustomerAddRs xml into the proper struct
	err := xml.Unmarshal(parentNode.Content, &customerAddRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling CustomerAddRs")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling CustomerAddRs")
		getLastErrChan <- err.Error()
	} else {
		//send items to CV3
		switch customerAddRs.StatusCode {
		case "0":
			Log.WithFields(logrus.Fields{
				"Status Severity": customerAddRs.StatusSeverity,
				"message":         customerAddRs.StatusMessage,
				"Status Code":     customerAddRs.StatusCode,
			}).Info("customerAddRs in")
			break
		case "3100":
			Log.WithFields(logrus.Fields{
				"Status Severity": customerAddRs.StatusSeverity,
				"message":         customerAddRs.StatusMessage,
				"Status Code":     customerAddRs.StatusCode,
			}).Info("customerAddRs in")
			//Check the message to see how to handle this error
			if checkWork.Attempted <= 2 {
				switch {
				//Customer already exists as a employee or vendor. Add "Cust" to the end of their name.
				case customerAddRs.StatusMessage[len(customerAddRs.StatusMessage)-39:] == " of the list element is already in use.":
					Log.WithFields(logrus.Fields{"OrderID": CheckPath("orderID", checkWork.Order)}).Info("Attempting to add a customer that already exists as a vendor or employee is not allowed")
					//Throw out the old order so this does not repeat
					<-workInsertChan
					var order = gabs.New()
					var workCTX = WorkCTX{Attempted: checkWork.Attempted}
					var workCount = 0
					var fieldMap = ReadFieldMapping("./fieldMaps/customerAddMapping.json")
					//append cust to the end of the last section of the customer name, as designated in customerAddMapping, and put that into the gabs order object
					checkWork.Order.SetP(CheckPath(fieldMap["Name"][len(fieldMap["Name"])-1].Data, checkWork.Order)+"Cust", fieldMap["Name"][len(fieldMap["Name"])-1].Data)
					/*
						//Check the data int the config nameArrangement's last field
						switch { //lowercase and check for the existance of first or last to allow for user error
						case strings.Contains(strings.ToLower(fieldMap["Name"][0]["data"], "first"):

							checkWork.Order.SetP(CheckPath("billing.firstName", checkWork.Order)+"Cust", "billing.firstName")
							break
						case strings.Contains(strings.ToLower(cfg.NameArrangement.Last), "last"):
							checkWork.Order.SetP(CheckPath("billing.lastName", checkWork.Order)+"Cust", "billing.lastName")
							break
						}*/
					//Put order info inside a gabs object to be compatible with the MakeOrder or MakeReceipt functions
					order.Set(checkWork.Order.Data(), "0")
					if err != nil {
						fmt.Println(err)
					}
					//Create new orders or receipts with the altered name
					switch strings.ToLower(cfg.OrderType) {
					case "salesreceipt":
						MakeSalesReceipt(&workCount, &workCTX, order)
						break
					case "salesorder":
						MakeSalesOrder(&workCount, &workCTX, order)
						break
					}
				}
				break
			}
		default:
			Log.WithFields(logrus.Fields{
				"Status Severity": customerAddRs.StatusSeverity,
				"message":         customerAddRs.StatusMessage,
				"Status Code":     customerAddRs.StatusCode,
			}).Error("Error in customerAddRs")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": customerAddRs.StatusSeverity,
				"message":         customerAddRs.StatusMessage,
				"Status Code":     customerAddRs.StatusCode,
			}).Error("Error in customerAddRs")
		}
	}
}

//ItemQueryRsHandler asdf
func ItemQueryRsHandler(parentNode Node, checkWork WorkCTX) {
	var itemQueryRs = ItemQueryRs{}
	//Unmarshal the ItemQueryRs xml into the proper struct
	err := xml.Unmarshal(parentNode.Content, &itemQueryRs)
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling ItemQueryRs")
		ErrLog.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling ItemQueryRs")
		getLastErrChan <- err.Error()
	} else {
		//send items to CV3
		switch itemQueryRs.StatusCode {
		case "0":
			Log.WithFields(logrus.Fields{
				"Status Severity": itemQueryRs.StatusSeverity,
				"message":         itemQueryRs.StatusMessage,
				"Status Code":     itemQueryRs.StatusCode,
			}).Info("itemQueryRs in")
			go UpdateCV3Items(itemQueryRs)
			break
		case "1":
			Log.WithFields(logrus.Fields{
				"Status Severity": itemQueryRs.StatusSeverity,
				"message":         itemQueryRs.StatusMessage,
				"Status Code":     itemQueryRs.StatusCode,
			}).Info("No QB item matches")
		default:
			Log.WithFields(logrus.Fields{
				"Status Severity": itemQueryRs.StatusSeverity,
				"message":         itemQueryRs.StatusMessage,
				"Status Code":     itemQueryRs.StatusCode,
			}).Error("Error in itemQueryRs")
			ErrLog.WithFields(logrus.Fields{
				"Status Severity": itemQueryRs.StatusSeverity,
				"message":         itemQueryRs.StatusMessage,
				"Status Code":     itemQueryRs.StatusCode,
			}).Error("Error in itemQueryRs")
		}
	}
}

//AutoImportItemsToQB will automatically import any items in a CV3 order, that do not exist in quickbooks
func AutoImportItemsToQB(checkWork WorkCTX) {
	if cfg.AutoImportCV3Items {
		for _, item := range checkWork.CV3Products {
			_ = item
			ImportCV3ItemsToQB(item)
			Log.WithFields(logrus.Fields{"Item Name": item.Name, "Item SKU": item.Sku}).Debug("Automatically importing order items from CV3 to QB")
		}
		workInsertChan <- checkWork
	}
	Log.Debug("automatic import of cv3 items is currently disabled")
}

// StrExtract Retrieves a string between two delimiters.
// sExper:  Specifies the expression to search.
// cAdelim: Specifies the character that delimits the beginning of sExper.
// cCdelim: Specifies the character that delimits the end of sExper.
// nOccur:  Specifies at which occurrence of cAdelim in sExper to start the extraction.
func StrExtract(sExper, sAdelim, sCdelim string, nOccur int) string {
	aExper := strings.Split(sExper, sAdelim)
	if len(aExper) <= nOccur {
		return ""
	}
	sMember := aExper[nOccur]
	aExper = strings.Split(sMember, sCdelim)

	return aExper[0]
}

//ConvertColonPairsToNewlines for use with stores that use :: as delimiters between newlines
func ConvertColonPairsToNewlines(incoming string) string {
	return strings.Replace(incoming, "::", "\n", -1)
}

//ConvertColonsToSpaces is to remove colons from qbxml as colons indicate parent/child relationships to quickbooks
func ConvertColonsToSpaces(incoming string) string {
	return strings.Replace(incoming, ":", " ", -1)
}

//StripBlanksAndNewlines returns an empty string "" if the incoming string contains nothing but spaces and newlines.
// Otherwise, return the string as it was received.
func StripBlanksAndNewlines(incoming string) string {
	reLeadcloseWhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reInsideWhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := reLeadcloseWhtsp.ReplaceAllString(incoming, "")
	final = reInsideWhtsp.ReplaceAllString(final, " ")
	if final == "" {
		return final
	}
	return incoming
}

//ConvertCustomerMsgRef will check for the cv3 elimiter :: and replace single colons with a space
func ConvertCustomerMsgRef(s string) string {
	converted := ConvertColonPairsToNewlines(s)
	converted = ConvertColonsToSpaces(converted)
	converted = StripBlanksAndNewlines(converted)
	return converted
}

//InitChannels Resets the go channels to prevent data being left in the buffers
func InitChannels() {
	shipToSuccessChan = make(chan ShipToSuccessTracker, 9999)
	orderSuccessChan = make(chan OrderSuccessTracker, 9999)
	confirmShipToChan = make(chan ShipToSuccessTracker, 9999)
	workChan = make(chan WorkCTX, 5)
	workInsertChan = make(chan WorkCTX, 5)
	checkWorckChan = make(chan WorkCTX, 5)
	checkWorckInsertChan = make(chan WorkCTX, 5)
	getLastErrChan = make(chan string, 5)
}
