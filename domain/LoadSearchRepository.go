package domain

import (
	"fmt"
	"golang/database"
	"golang/dto"
	"golang/errs"
	"golang/logger"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

type GetDbClient func() *sqlx.DB

type RepositoryDb struct {
	client GetDbClient
}

//AP_PC_07
//GetCPDBLoadResultRepo repostiory function is declared in here inside this function the Postgres (CP DB) is connected and the Load query is executed to get the
//load records based on the request
func (d RepositoryDb) GetCPDBLoadResultRepo(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]LoadSearchResDTO, int, *errs.AppErrorvalidation) {

	//Postgres CP DB connection is made in here
	client, DbError := database.GetCPDbClient()

	loadResponse := make([]LoadSearchResDTO, 0)
	if DbError != nil {
		return nil, 0, errs.ValidateResponse(nil, DbError.Code, DbError.Message)
	}

	//select query to get the load records
	loadQuery := `select ls.source_load_shipment_id as "loadId", ls.invoice_date as "invoiceDate", lst.status as "loadStatus", ls.shipper_city as "shipperCity", ls.shipper_state as "shipperState",  
                         ls.consignee_city as "consigneeCity", ls.consignee_state as "consigneeState", ls.total as "custTotal", ls.customer_purchase_order as "poNumber",
                         ls.shipper_bill_of_ladding as "shipBlNumber", lsc.pro_number as "proNumber", lssPICK.earliest_date as "shipDateStart", lssDROP.earliest_date as "shipDateEnd",
                         lssPICK.po_number as "pickupNumber", lssDROP.po_number as "deliveryNumber", ls.last_modified_date as "lastModified", 'otherTMS' as "loadOrigin", lsd.description as "loadMethod"
                   from cpdb.load_shipment ls
inner join cpdb.load_status lst on ls.load_status_id = lst.load_status_id
inner join cpdb.load_method lsd on ls.load_method_id = lsd.load_method_id
inner join cpdb.load_carrier lsc on lsc.load_shipment_id = ls.load_shipment_id
left join cpdb.load_shipment_stops lssPICK on lssPICK.load_carrier_id = lsc.load_carrier_id and lssPICK.load_shipment_stops_type_id = 1
left join cpdb.load_shipment_stops lssDROP on lssDROP.load_carrier_id = lsc.load_carrier_id and lssDROP.load_shipment_stops_type_id = 2
inner join cpdb.load_customer ldc on ldc.load_shipment_id = ls.load_shipment_id `

	//Wherey conditions are starts to form based on the request came
	whereQuery := `WHERE  (lst.status NOT IN ('QUOTE', 'DEAD', 'TEMPLATE')) `

	gcidValues := `AND ldc.global_customer_id in (`

	for i := 0; i < len(LoadSearchReqDTO.Gcid); i++ {
		if i-len(LoadSearchReqDTO.Gcid) == -1 {
			gcidValues += `'` + LoadSearchReqDTO.Gcid[i] + `')`
		} else {
			gcidValues += `'` + LoadSearchReqDTO.Gcid[i] + `',`
		}
	}

	whereQuery = whereQuery + gcidValues

	if LoadSearchReqDTO.LoadMethod != "" {
		whereQuery += `AND lsd.description ='` + LoadSearchReqDTO.LoadMethod + `'`
	}
	if LoadSearchReqDTO.LoadStatus != "" {
		whereQuery += `AND lst.status ='` + LoadSearchReqDTO.LoadStatus + `'`
	}
	if LoadSearchReqDTO.ShipDateStart != "" {
		whereQuery += `AND lssPICK.earliest_date >='` + LoadSearchReqDTO.ShipDateStart + `'`
	}
	if LoadSearchReqDTO.ShipDateEnd != "" {
		whereQuery += `AND lssPICK.earliest_date <='` + LoadSearchReqDTO.ShipDateEnd + `'`
	}

	LoadSearchReqDTO.PickupNumber = strings.ReplaceAll(LoadSearchReqDTO.PickupNumber, "'", "")
	LoadSearchReqDTO.DeliveryNumber = strings.ReplaceAll(LoadSearchReqDTO.DeliveryNumber, "'", "")
	LoadSearchReqDTO.CustomerblNumber = strings.ReplaceAll(LoadSearchReqDTO.CustomerblNumber, "'", "")
	LoadSearchReqDTO.PoNumber = strings.ReplaceAll(LoadSearchReqDTO.PoNumber, "'", "")

	if LoadSearchReqDTO.LoadId != "" {
		whereQuery = whereQuery + `AND ls.source_load_shipment_id like '%` + LoadSearchReqDTO.LoadId + `%'`
	} else if LoadSearchReqDTO.PickupNumber != "" {
		whereQuery = whereQuery + `AND lssPICK.po_number like '%` + LoadSearchReqDTO.PickupNumber + `%'`
	} else if LoadSearchReqDTO.DeliveryNumber != "" {
		whereQuery = whereQuery + `AND lssDROP.po_number like '%` + LoadSearchReqDTO.DeliveryNumber + `%'`
	} else if LoadSearchReqDTO.CustomerblNumber != "" {
		whereQuery = whereQuery + `AND ls.shipper_bill_of_ladding like '%` + LoadSearchReqDTO.CustomerblNumber + `%'`
	} else if LoadSearchReqDTO.PoNumber != "" {
		whereQuery = whereQuery + `AND ls.customer_purchase_order like '%` + LoadSearchReqDTO.PoNumber + `%'`
	}

	//orderBy clause is created based on the request
	orderBy := ""

	switch strings.ToUpper(LoadSearchReqDTO.SortColumn) {
	case "LOADID":
		orderBy += `ORDER BY ls.source_load_shipment_id`
	case "LOADSTATUS":
		orderBy += `ORDER BY lst.status`
	case "MOVETYPE":
		orderBy += `ORDER BY lsd.description`
	case "ORIGIN":
		orderBy += `ORDER BY ls.shipper_city`
	case "DESTINATION":
		orderBy += `ORDER BY ls.consignee_city`
	case "TOTAL":
		orderBy += `ORDER BY ls.total`
	default:
		orderBy += `ORDER BY ls.last_modified_date`
	}

	if strings.ToUpper(LoadSearchReqDTO.SortOrder) == "DESCENDING" {
		orderBy += ` DESC`
	} else if strings.ToUpper(LoadSearchReqDTO.SortOrder) == "ASCENDING" {
		orderBy += ` ASC`
	} else {
		orderBy += ` DESC`
	}

	//limit is created based on the request
	limit := ""
	if LoadSearchReqDTO.OtmsLoadRecords != "" {
		limit += ` LIMIT 10 OFFSET ` + LoadSearchReqDTO.OtmsLoadRecords
	} else {
		limit += ` LIMIT 10 OFFSET 0`
	}
	loadQuery += whereQuery + orderBy + limit

	// fmt.Print(loadQuery)

	loadCountQuery := `select count(ls.source_load_shipment_id) as "loadCount"
                   from cpdb.load_shipment ls
inner join cpdb.load_status lst on ls.load_status_id = lst.load_status_id
inner join cpdb.load_method lsd on ls.load_method_id = lsd.load_method_id
inner join cpdb.load_carrier lsc on lsc.load_shipment_id = ls.load_shipment_id
left join cpdb.load_shipment_stops lssPICK on lssPICK.load_carrier_id = lsc.load_carrier_id and lssPICK.load_shipment_stops_type_id = 1
left join cpdb.load_shipment_stops lssDROP on lssDROP.load_carrier_id = lsc.load_carrier_id and lssDROP.load_shipment_stops_type_id = 2
inner join cpdb.load_customer ldc on ldc.load_shipment_id = ls.load_shipment_id `

	loadCountQuery += whereQuery

	var count LoadCountRes
	//loadCount query is executed in here
	err1 := client.Get(&count, loadCountQuery)

	if err1 != nil {
		logger.Error("Error selecting values in the database:" + err1.Error())
		return nil, 0, errs.ValidateResponse(nil, 500, "Error selecting values in the database:"+err1.Error())
	}

	//loadQuery is executed in here
	err := client.Select(&loadResponse, loadQuery)

	client.Close()

	if err != nil {
		logger.Error("Error selecting values in the database:" + err.Error())
		return nil, 0, errs.ValidateResponse(nil, 500, "Error selecting values in the database:"+err.Error())
	}

	return loadResponse, count.LoadCount, nil

}

func NewLoadRepositoryDb() RepositoryDb {
	return RepositoryDb{}
}

//AP_PC_10
//GetBTMSLoadResultRepo repostiory function is declared in here inside this function the MySQL (BTMS DB) is connected and the Load query is executed to get the
//load records based on the request
func (d RepositoryDb) GetBTMSLoadResultRepo(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]LoadSearchResDTO, int, *errs.AppErrorvalidation) {

	dbclient, Dberror := database.BTMSDBConnection()

	if Dberror != nil {
		log.Fatal(Dberror)
		return nil, 0, errs.ValidateResponse(nil, 500, "Unexpected database error")
	}

	loadResponse := make([]LoadSearchResDTO, 0)

	//select query to get the load records
	loadQuery := `SELECT cast(loadsh.id as char) as loadId, invoice_date as invoiceDate, loadsh.status AS loadStatus, ship_date as shipDateStart, ship_city as shipperCity, ship_state as shipperState, loadsh.cust_total as custTotal,
      cust_po as poNumber,	ship_bl as shipBlNumber, cons_city as consigneeCity, cons_state as consigneeState, cons_date as shipDateEnd, lsstops.po as pickupNumber, last_stop.po as deliveryNumber,
      first_carr.pro_number as proNumber, loadsh.load_method as loadMethod, loadsh.lastmod as lastModified,"BTMS" as loadOrigin
	FROM loadsh
	LEFT JOIN lscarr first_carr ON first_carr.loadsh_id=loadsh.id AND first_carr.carr_order=1
	LEFT JOIN lscarr last_carr ON last_carr.loadsh_id=loadsh.id AND last_carr.carr_order=(SELECT MAX(carr_order) FROM lscarr WHERE loadsh_id=loadsh.id AND type='carrier')
	LEFT JOIN lsstops ON lsstops.lscarr_id=first_carr.id AND lsstops.stop_order=1
	LEFT JOIN lsstops last_stop ON last_stop.lscarr_id=last_carr.id AND last_stop.stop_order=(SELECT MAX(stop_order) FROM lsstops WHERE lscarr_id=last_carr.id) `

	//Where conditions are starts to form based on the request came
	whereQuery := `WHERE  (loadsh.status NOT IN ('QUOTE', 'DEAD', 'TEMPLATE')) `

	customerIdValues := `AND custm_id in (`

	for i := 0; i < len(LoadSearchReqDTO.CustomerId); i++ {
		if i-len(LoadSearchReqDTO.CustomerId) == -1 {
			customerIdValues += strconv.Itoa(LoadSearchReqDTO.CustomerId[i]) + `)`
		} else {
			customerIdValues += strconv.Itoa(LoadSearchReqDTO.CustomerId[i]) + `,`
		}
	}

	officeCodeValues := `AND loadsh.office in (`

	for i := 0; i < len(LoadSearchReqDTO.OfficeCode); i++ {
		if i-len(LoadSearchReqDTO.OfficeCode) == -1 {
			officeCodeValues += `'` + LoadSearchReqDTO.OfficeCode[i] + `')`
		} else {
			officeCodeValues += `'` + LoadSearchReqDTO.OfficeCode[i] + `',`
		}
	}

	whereQuery = whereQuery + customerIdValues + officeCodeValues

	if LoadSearchReqDTO.LoadMethod != "" {
		whereQuery += `AND loadsh.load_method ='` + LoadSearchReqDTO.LoadMethod + `'`
	}
	if LoadSearchReqDTO.LoadStatus != "" {
		whereQuery += `AND loadsh.status ='` + LoadSearchReqDTO.LoadStatus + `'`
	}
	if LoadSearchReqDTO.ShipDateStart != "" {
		whereQuery += `AND loadsh.ship_date >='` + LoadSearchReqDTO.ShipDateStart + `'`
	}
	if LoadSearchReqDTO.ShipDateEnd != "" {
		whereQuery += `AND loadsh.ship_date <='` + LoadSearchReqDTO.ShipDateEnd + `'`
	}

	LoadSearchReqDTO.PickupNumber = strings.ReplaceAll(LoadSearchReqDTO.PickupNumber, "'", "")
	LoadSearchReqDTO.DeliveryNumber = strings.ReplaceAll(LoadSearchReqDTO.DeliveryNumber, "'", "")
	LoadSearchReqDTO.CustomerblNumber = strings.ReplaceAll(LoadSearchReqDTO.CustomerblNumber, "'", "")
	LoadSearchReqDTO.PoNumber = strings.ReplaceAll(LoadSearchReqDTO.PoNumber, "'", "")
	re := regexp.MustCompile(`^\d+$`)

	if LoadSearchReqDTO.LoadId != "" && re.MatchString(LoadSearchReqDTO.LoadId) {
		whereQuery = whereQuery + `AND loadsh.id like '%` + LoadSearchReqDTO.LoadId + `%'` //` + LoadSearchReqDTO.LoadId
	} else if LoadSearchReqDTO.PickupNumber != "" {
		whereQuery = whereQuery + `AND lsstops.po like '%` + LoadSearchReqDTO.PickupNumber + `%'`
	} else if LoadSearchReqDTO.DeliveryNumber != "" {
		whereQuery = whereQuery + `AND last_stop.po like '%` + LoadSearchReqDTO.DeliveryNumber + `%'`
	} else if LoadSearchReqDTO.CustomerblNumber != "" {
		whereQuery = whereQuery + `AND loadsh.ship_bl like '%` + LoadSearchReqDTO.CustomerblNumber + `%'`
	} else if LoadSearchReqDTO.PoNumber != "" {
		whereQuery = whereQuery + `AND loadsh.cust_po like '%` + LoadSearchReqDTO.PoNumber + `%'`
	}

	orderBy := ""

	switch strings.ToUpper(LoadSearchReqDTO.SortColumn) {
	case "LOADID":
		orderBy += ` ORDER BY loadsh.id`
	case "LOADSTATUS":
		orderBy += ` ORDER BY loadsh.status`
	case "MOVETYPE":
		orderBy += ` ORDER BY loadsh.load_method`
	case "ORIGIN":
		orderBy += ` ORDER BY loadsh.ship_city`
	case "DESTINATION":
		orderBy += ` ORDER BY loadsh.cons_city`
	case "TOTAL":
		orderBy += ` ORDER BY loadsh.cust_total`
	default:
		orderBy += ` ORDER BY loadsh.lastmod`
	}

	if strings.ToUpper(LoadSearchReqDTO.SortOrder) == "DESCENDING" {
		orderBy += ` DESC`
	} else if strings.ToUpper(LoadSearchReqDTO.SortOrder) == "ASCENDING" {
		orderBy += ` ASC`
	} else {
		orderBy += ` DESC`
	}

	limit := ""
	if LoadSearchReqDTO.BtmsLoadRecords != "" {
		limit += ` LIMIT ` + LoadSearchReqDTO.BtmsLoadRecords + `, 10`
	} else {
		limit += ` LIMIT 0, 10`
	}

	loadQuery += whereQuery + orderBy + limit

	loadCountQuery := `SELECT count(loadsh.id) as loadCount FROM loadsh
	LEFT JOIN lscarr first_carr ON first_carr.loadsh_id=loadsh.id AND first_carr.carr_order=1
	LEFT JOIN lscarr last_carr ON last_carr.loadsh_id=loadsh.id AND last_carr.carr_order=(SELECT MAX(carr_order) FROM lscarr WHERE loadsh_id=loadsh.id AND type='carrier')
	LEFT JOIN lsstops ON lsstops.lscarr_id=first_carr.id AND lsstops.stop_order=1
	LEFT JOIN lsstops last_stop ON last_stop.lscarr_id=last_carr.id AND last_stop.stop_order=(SELECT MAX(stop_order) FROM lsstops WHERE lscarr_id=last_carr.id)`

	loadCountQuery += whereQuery

	var count LoadCountRes
	fmt.Printf("loadQuery: %+v\n", loadQuery)

	fmt.Printf("loadCountQuery: %+v\n", loadCountQuery)
	//loadCount query is executed in here
	err1 := dbclient.Get(&count, loadCountQuery)

	if err1 != nil {
		logger.Error("Error while selecting values in the database:" + err1.Error())
		return nil, 0, errs.ValidateResponse(nil, 500, "Unexpected database error")
	}
	//loadQuery is executed in here
	err2 := dbclient.Select(&loadResponse, loadQuery)

	//fmt.Printf("load response: %+v\n", loadResponse)

	dbclient.Close()

	if err2 != nil {
		logger.Error("Error selecting values in the database: load" + err2.Error())
		return nil, 0, errs.ValidateResponse(nil, 500, "Unexpected database error")
	}

	return loadResponse, count.LoadCount, nil

}

// AP_PC_07
//GetCPDBLoadDocsRepo repostiory function is declared in here inside this function the Postgres (CP DB) is connected and the Load Doc query is executed to get the
//document URL's based on the loadId's passed
func (d RepositoryDb) GetCPDBLoadDocsRepo(LoadDocsReqDTO *dto.LoadDocumentsReqDTO) ([]LoadCPDBResDTO, *errs.AppErrorvalidation) {
	//Postgres CP DB connection is made in here

	loadDocsRes := make([]LoadCPDBResDTO, 0)

	if len(LoadDocsReqDTO.OthertmsLoads) == 0 {
		return loadDocsRes, nil
	}

	client, DbError := database.GetCPDbClient()

	if DbError != nil {
		return nil, errs.ValidateResponse(nil, DbError.Code, DbError.Message)
	}

	//select query to get the document URL
	loadDocsQuery := `select ld.document_id as "id", ld.document_type_id as "typeId", ls.source_load_shipment_id as "loadId" , ld.document_name as "docName", ld.document_url as "docUrl"
from cpdb.document ld
inner join cpdb.document_type dt on dt.document_type_id = ld.document_type_id
inner join cpdb.load_shipment ls on ls.load_shipment_id = ld.load_shipment_id `

	//whereQuery is formed based on the request parameter
	whereQuery := `WHERE ls.source_load_shipment_id in (`

	for i := 0; i < len(LoadDocsReqDTO.OthertmsLoads); i++ {
		if i-len(LoadDocsReqDTO.OthertmsLoads) == -1 {
			whereQuery += `'` + LoadDocsReqDTO.OthertmsLoads[i] + `')`
		} else {
			whereQuery += `'` + LoadDocsReqDTO.OthertmsLoads[i] + `',`
		}
		// 	whereQuery += strconv.Itoa(LoadDocsReqDTO.OthertmsLoads[i]) + `)`
		// } else {
		// 	whereQuery += strconv.Itoa(LoadDocsReqDTO.OthertmsLoads[i]) + `,`
		// }
	}

	loadDocsQuery += whereQuery

	fmt.Print(loadDocsQuery)

	//loadDocsQuery is executed in here
	err := client.Select(&loadDocsRes, loadDocsQuery)

	// fmt.Print(loadDocsRes)

	client.Close()

	if err != nil {
		logger.Error("Error selecting values in the database:" + err.Error())
		return nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err.Error())
	}

	return loadDocsRes, nil

}

//AP_PC_10
//GetBTMSLoadDocsRepo repostiory function is declared in here inside this function the BTMS (MySQL DB) is connected and the query is executed to get the
// load status and load method for the loadId's passed
func (d RepositoryDb) GetBTMSLoadDocsRepo(LoadDocsReqDTO *dto.LoadDocumentsReqDTO) ([]LoadBTMSDBResDTO, *errs.AppErrorvalidation) {

	loadDocsResponse := make([]LoadBTMSDBResDTO, 0)

	if len(LoadDocsReqDTO.BtmsLoads) == 0 {

		return loadDocsResponse, nil
	}

	dbclient, Dberror := database.BTMSDBConnection()

	if Dberror != nil {
		log.Fatal(Dberror)
		return nil, errs.ValidateResponse(nil, 500, "Unexpected database error")
	}

	//select query to get the load status and load methods
	loadDocsQuery := `SELECT cast(loadsh.id as char) as loadId, loadsh.status AS loadStatus,loadsh.load_method as loadMethod FROM loadsh `

	whereQuery := `WHERE loadsh.id in (`

	for i := 0; i < len(LoadDocsReqDTO.BtmsLoads); i++ {
		if i-len(LoadDocsReqDTO.BtmsLoads) == -1 {
			whereQuery += LoadDocsReqDTO.BtmsLoads[i] + `)`
		} else {
			whereQuery += LoadDocsReqDTO.BtmsLoads[i] + `,`
		}
	}

	loadDocsQuery += whereQuery

	//loadDocQuery is executed in here
	err1 := dbclient.Select(&loadDocsResponse, loadDocsQuery)

	// fmt.Printf("load response: %+v\n", loadDocsResponse)

	dbclient.Close()

	if err1 != nil {
		logger.Error("Error selecting values in the database: load" + err1.Error())
		return nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err1.Error())
	}

	return loadDocsResponse, nil

}
