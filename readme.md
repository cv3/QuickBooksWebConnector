# CommerceV3 QuickBooks Web Connector #
## Supported Versions of QuickBooks ##
The CommerceV3 QuickBooks Web Connector supports all versions of QuickBooks that Intuit's Web Connector software supports.  At this time, that includes the following:
  - U.S. editions of QuickBooks Financial Software products
      - QuickBooks Enterprise Solutions
      - QuickBooks Premier (2002 or later)
      - QuickBooks Pro (2002 or later)
      - QuickBooks Simple Start (2006 or later)
  - QuickBooks Point of Sale (v4.0 or later)
  - Canadian editions of QuickBooks Pro, Premier or Enterprise (2003 or later)
  - UK editions of QuickBooks Pro or Accountant Edition (2003 or later)

For more information, please visit https://developer.intuit.com/docs/0200_quickbooks_desktop/0400_tools/web_connector

The CV3 App (also referred to as the M.O.M. Connector) QuickBooks module has been deprecated as of QuickBooks 2015.  The software described in this document is the recommended method to integrate CommerceV3 with QuickBooks.


## Installation Guide ##

QuickBooks 2015 or later must be installed.  Installation of QuickBooks requires a valid license code and activation, provided by Intuit.  QuickBooks may be obtained from http://support.quickbooks.intuit.com/support/ProductUpdates.aspx

### Download & install the  QuickBooks Web Connector ###
Before continuing, QuickBooks 2015 or later must be installed, but not running.

1. The QuickBooks Web Connector may be obtained from https://developer.intuit.com/docs/0200_quickbooks_desktop/0100_essentials/quickbooks_web_connector
    1. Click **Download, unzip and install the QuickBooks Web Connector**
        1. Click the Download button
    
2. Once the QBWebConnector archive has downloaded, extract the archive, then explore into the directory.

3. Right click on QBWebConnectorInstaller.exe and select "Run as administrator" to install.
    1. If prompted, "Do you want to allow the following program to make changes to this computer?", select **Yes**.
    2. Proceed through the installer using the defaults until complete.
    3. The QBWebConnector archive & extracted directory may now be deleted.
    4. Run the QuickBooks Web Connector after installation by clicking Start > Web Connector
    5. Click the **Hide** button at the lower-right corner of the QuickBooks Web Connector to minimize to the system tray.
    
### Download & install the CV3 Integration Server ###
1. The latest CV3 Integration Server may be obtained from https://github.com/cv3/QuickBooksWebConnector/releases/
    1. Click **qbwcServer.zip**

2. Once qbwcServer.zip has downloaded, use [7-Zip](https://www.7-zip.org/) to extract the archive (Windows built-in unzip application may fail to extract the archive), then explore into the directory.

3. Double-click **createStartup.bat** to enable launching of the CV3 Integration Server when the current user logs in.  The batch file will quickly run then close without any prompts.

### Configure the QuickBooks Web Connector with the CV3 Integration Server ###

1. Using a text editor such as Notepad++, edit qbwcServer/wsdl/**QBWebConnect.qwc**
    1. Enter your desired QuickBooks Web Connector username in between `<UserName>`...`</UserName>`.  It is recommended that you not keep the default value, "user".
    2. Save the file.

2. Using a text editor, edit the config file located in qbwcServer/config/**qbwcConfig.json**
    1. Under **"cv3Credentials"**, enter your store's CV3 webservice information.  Please contact CommerceV3 support for assistance obtaining your webservice credentials.
    ```
    "cv3Credentials": {
		"user": "cv3_webservice_username", // Replace cv3_webservice_username with your credentials
		"pass": "cv3_webservice_password", // Replace cv3_webservice_password with your credentials
		"serviceID": "cv3_webservice_id"   // Replace cv3_webservice_id with your credentials
	}
    ```
    2. Under **"qbwcCredentials"**, replace "user": "**user**" as it was entered in the  qbwcServer/wsdl/QBWebConnect.qwc <UserName> field.  Replace "pass": "**pass**" with the password you intend to use.
    ```
    "qbwcCredentials": {
		"user": "qbwc_username", // Replace qbwc_username with your credentials
		"pass": "qbwc_password"  // Replace qbwc_password with your credentials
	}
    ```
    3. If you want to send all of your QuickBooks Inventory Items to CV3, set "updateCV3Items" to **true**
    4. Set "orderType" to **SalesOrder** or **SalesReceipt**.  Note that to use sales orders, QuickBooks Desktop must be the Premier or Enterprise version.  If you are running another version, such as QuickBooks Desktop Pro, set orderType to SalesReceipt.
    5. Save the file

3. Manually launch the CV3 Integration Server by double-clicking qbwcServer/**qbwcServer.exe**  (It will launch automatically whenever the current Windows user logs in, and may be launched manually at this time).

### Import the QuickBooks Web Connector configuration into QuickBooks ###

1. Start QuickBooks as a user with administrator privileges, then minimize it to the taskbar by clicking the _ at the top-right.

2. Double-click qbwcServer/wsdl/**QBWebConnect.qwc**
    1. Click **OK** to grant the web service access to QuickBooks.
    2. When prompted "Do you want to allow this application to read and modify this company file?", select **Yes, whenever this QuickBooks company file is open**.  Then click **Continue...**
    3. Click **Done**.

3. In the Quick Books Web Connector:
    1. Enter your password as it was entered in qbwcServer/wsdl/**QBWebConnect.qwc**.  When asked to save the password click **Yes**
    2. Check the box on the left if it is unchecked.
    3. Click "Update selected" to manually run the application.
    
4. Create a shortcut on the desktop to start the QBWebConnector
    1. Explore to C:\Program Files (x86)\Common Files\Intuit\QuickBooks\QBWebConnector\
    2. Right click on QBWebConnector.exe and select Send to > Desktop (Create shortcut)

---

## Confirming QuickBooks and CV3 are communicating ##

### Export items from QuickBooks to CV3 ###
1. In QuickBooks, select Edit > Preferences...
    1. Accounting > Company Preferences
        - Check "Use account numbers" and "Use class tracking for transactions".
    2. Items & Inventory
        - Check "Inventory and purchase orders are active".

2. Create a new item by selecting Lists > Item List
    1. Item > New
        - Type: Inventory Part  (Service part & non-inventory part will also work)
        - Item Name/Number = This will be the CV3 Product Name
        - Sales Price = CV3 Retail Standard Price
        - Income Account > Pick from the dropdown.
        - On Hand = CV3 Inventory Count
        - Click "OK"

3. Open QuickBooks Web Connector and click "Update Selected".  Once it completes, the item you created should appear in your CV3 store's All Products page.

### Import orders from CV3 to QuickBooks ###
1. Visit your store's website and place an order with the item(s) exported from QuickBooks.  After completing checkout, return to QuickBooks Web Connector and click "Update Selected".

2. To view the orders in QuickBooks, click Reports > Report Center
    1. Click Sales (Left Column) > Sales by Ship To Address (Run)
    2. Customize Report > Filters
    3. Current Filter Choices: Remove **Name** and **TransactionType**, then click **OK**
    - After clicking OK, the test order should be displayed.

#### Troubleshooting ####
Intuit saves logging in C:\ProgramData\Intuit\QBWebConnector\log\QWCLog.txt
This log file will show you both the QBXML request and response send between QuickBooks and the QBWC server.

To enable verbose logging:
  - Exit the Web Connector
  - Run regedit (Start > Run > regedit.exe)
  - Navigate to: \HKEY_CURRENT_USER\Software\Intuit\QBWebConnector
  - Change the 'Level' key to VERBOSE  (The default setting is DEBUG)
  - Start the QuickBooks Web Connector and attempt an update.

Note about permissions to Intuit's log file:
  - All users running QBWC will need permission to write to C:\ProgramData\Intuit\QBWebConnector\log\QWCLog.txt

---

## Additional Instructions ##

### Field Mapping ###
Custom field mapping is performed in the json files within qbwcServer/**fieldMaps**. For field mapping documentation, view the readme in [qbwcServer/fieldMaps/mapping docs](fieldMaps/mapping%20docs).
