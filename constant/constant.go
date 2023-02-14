package constant

import "os"

var (

	// Base URLS

	IMAGING_END_POINT         = os.Getenv("IMAGING_END_POINT")
	BROKERAGE_URL             = os.Getenv("BROKERAGE_URL")
	CUSTOMER_PORTAL_END_POINT = os.Getenv("CUSTOMER_PORTAL_END_POINT")
	BASE_URL                  = os.Getenv("BASE_URL")
	BANYAN_URL                = os.Getenv("BANYAN_URL")

	//Customer portal database (CPDB) config details

	PORT       = os.Getenv("PORT")
	CP_DB_USER = os.Getenv("CP_DB_USER")
	CP_DB_ADDR = os.Getenv("CP_DB_ADDR")
	CP_DB_PASS = os.Getenv("CP_DB_PASS")
	CP_DB_NAME = os.Getenv("CP_DB_NAME")

	//BTMS Connection Details

	BC_DB_USER = os.Getenv("BC_DB_USER")
	BC_DB_PASS = os.Getenv("BC_DB_PASS")
	BC_DB_HOST = os.Getenv("BC_DB_HOST")
	BC_DB_NAME = os.Getenv("BC_DB_NAME")

	//Shared database SDB)config details

	SH_DB_USER = os.Getenv("SH_DB_USER")
	SH_DB_ADDR = os.Getenv("SH_DB_ADDR")
	SH_DB_PASS = os.Getenv("SH_DB_PASS")
	SH_DB_NAME = os.Getenv("SH_DB_NAME")

	//Access token issue

	ACCESS_TOKEN_URL          = os.Getenv("ACCESS_TOKEN_URL")
	ACCESS_TOKEN_GRANTTYPE    = os.Getenv("ACCESS_TOKEN_GRANTTYPE")
	ACCESS_TOKEN_CLIENTID     = os.Getenv("ACCESS_TOKEN_CLIENTID")
	ACCESS_TOKEN_CLIENTSECRET = os.Getenv("ACCESS_TOKEN_CLIENTSECRET")

	InvoiceKey = os.Getenv("INVOICE_KEY")
)
