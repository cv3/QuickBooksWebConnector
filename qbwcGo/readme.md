QuickBooks Web Connector Programmers Guide:
https://developer-static.intuit.com/qbSDK-current/doc/PDF/QBWC_proguide.pdf

QuickBookds IDN Unified On Screen Reference:
https://developer-static.intuit.com/qbSDK-current/Common/newOSR/index.html

CommerceV3 Prezi describing the CV3 integration service - https://prezi.com/view/YvGSjnefni3pkLZaracf/

The CV3 integration service for QuickBooks Web Connector handles communication betwwen CommerceV3 stores and QuickBooks.  All messages to and from the QuickBooks Web Connector flow through the URL set in the <AppURL> field in wsdl/QBWebConnect.qwc

The Intuit Web Connector sessions consists of several calls to the CV3 integration server.  All roundtrip messages begin within Intuit's web connector, and every message sent from the CV3 integration service to Intuit's web connector is a response to a request:

1. ServerVersion
    1. Asks the server for its version information
2. ClientVersion
    1. Intuit Web Connector sends its versioning information to the server
3. Authenticate
    1. Performs basic authentication
    2. Starts the NoOp holding patter to tell the Intuit Web Connector to wait for work
    3. Starts the order success tracking system 
    4. Gets the work from CV3 via the InitWork() function which calls the GetCV3Orders() function
        1. A call to the CV3 webservice to retrieve orders
        2. Map the CV3 order information to the correct qbXML order type as set in the config
        3. Send the qbXML to the work channel
        4. Send order data to the order success tracking system
4. SendRequestXML
    1. Check for work in the work channel, otherwise continue the NoOp holding pattern
    2. Add work context to the sendRequestXML template
    3. Send the work back to the Intuit Web Connector
5. ReceiveResponseXML
    1. Receives the response from the preceeding SendRequestXML call
    2. Checks for what type of response via the CheckNode function and routs the the corresponding handler function
        1. Checks for success or error status codes 
        2. If success, the order success tracker will be updated to allow for an order confirmation to be sent to CV3
        3. If errors are found, the service will trigger the creation of more work to be sent to Intuit's Web Connector
    3. Intuit's Web Connector will then call sendRequestXML to start a new work loop
6. CloseConnection
    1. Ends the session when work is done


ReceiveResponseXML handlers:

1. SalesOrderAddRs
    1. Check StatusCode
        1. "0" OK, send the order tracking system a success signal
        2. "3140" Element not found in QuickBooks
            1. Check StatusMessage to potentially:
                1. Add new Customer
                2. Add new CustomerMSG
        3. "3180" occurs when quickbooks thinks the list is being accessed from another location
            1. May only happen in Enterprise version  
            2. Resend the work to the work channel
        4. "3270" occurs when using SalesOrderAdds with an unsupported version of QuickBooks
            1. If this occurs, the config may be set to use SalesOrders, and instead should use SalesReceipts
2. SalesReceiptAddRs
    1. Check StatusCode
        1. "0" OK, send the order tracking system a success signal
        2. "3140" Element not found in QuickBooks
            1. Check StatusMessage to potentially:
                1. Add new Customer
                2. Add new CustomerMSG
        3. "3180" occurs when quickbooks thinks the list is being accessed from another location
            1. May only happen in Enterprise version  
            2. Resend the work to the work channel
3. CustomerAddRs
    1. Check StatusCode
        1. "0" OK, Customer added successfully
        2. "3100" Element already exists in QuickBooks
            1. Check StatusMsg to potentially
                1. Customer name already exists as a vendor or employee
                    1. Append "Cust" to the end of the customer name and re attempt customer add
                    2. Update the order to match the new customer name
These calls are routed using a string match on the incoming XML.  A recursive CheckNodes function is used to get a match on any XML node desired, then the corresponding function will be called.
