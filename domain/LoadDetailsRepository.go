package domain

import (
	"fmt"
	"golang/database"
	"golang/dto"
	"golang/errs"
	"golang/logger"
	"log"
)

//AP_PC_15
//GetOTMSLoadDetailsRepo repostiory function is declared in here inside this function the postgres SQL (CP DB) is connected and the Load query is executed to get the
//load details based on the loadId
func (d RepositoryDb) GetOTMSLoadDetailsRepo(loadDetailsVar dto.GetLoadDetailsReqDTO) ([]LoadShipConsRefDTO, []LoadCommoditiesDTO, []LoadCPDBResDTO, *errs.AppErrorvalidation) {

	dbclient, Dberror := database.GetCPDbClient()

	if Dberror != nil {
		log.Fatal(Dberror)
		return nil, nil, nil, errs.ValidateResponse(nil, Dberror.Code, Dberror.Message)
	}

	LoadShipConsRefRes := make([]LoadShipConsRefDTO, 0)

	loadCommoditiesRes := make([]LoadCommoditiesDTO, 0)

	loadDocsRes := make([]LoadCPDBResDTO, 0)

	var shipConsRefQuery = LoadShipConsRefQuery + `'` + loadDetailsVar.LoadId + `'`

	var commodityQuery = LoadCommodityQuery + `'` + loadDetailsVar.LoadId + `'`

	var loadDocumentsQuery = LoadDocsQuery + `'` + loadDetailsVar.LoadId + `'`

	// shipConsRefQuery is executed in here
	err := dbclient.Select(&LoadShipConsRefRes, shipConsRefQuery)

	if err != nil {
		logger.Error("Error selecting values in the database:" + err.Error())
		return nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err.Error())
	}

	// returning empty response for invalid loadId
	if len(LoadShipConsRefRes) == 0 {
		return LoadShipConsRefRes, nil, nil, nil
	}

	// commodityQuery is executed in here
	err1 := dbclient.Select(&loadCommoditiesRes, commodityQuery)

	fmt.Println(commodityQuery)

	if err1 != nil {
		logger.Error("Error selecting values in the database:" + err1.Error())
		return nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err1.Error())
	}

	err2 := dbclient.Select(&loadDocsRes, loadDocumentsQuery)

	dbclient.Close()

	if err2 != nil {
		logger.Error("Error selecting values in the database:" + err2.Error())
		return nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err2.Error())
	}

	// fmt.Printf("LoadShipConsRefRes: %+v", LoadShipConsRefRes)

	return LoadShipConsRefRes, loadCommoditiesRes, loadDocsRes, nil

}

//AP_PC_07
//GetBTMSLoadDetailsRepo repostiory function is declared in here inside this function the mysql SQL (BTMS DB) is connected and the Load query is executed to get the
//load details based on the loadId
func (d RepositoryDb) GetBTMSLoadDetailsRepo(loadDetailsVar dto.GetLoadDetailsReqDTO) ([]LoadShipConsRefDTO, []LoadCommoditiesDTO, []EventTrackingUpdatesTL, []LocationBreadCrumbs, []EventTrackingUpdatesELTL, *errs.AppErrorvalidation) {

	dbclient, Dberror := database.BTMSDBConnection()

	if Dberror != nil {
		log.Fatal(Dberror)
		return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, Dberror.Code, Dberror.Message)
	}

	LoadShipConsRefRes := make([]LoadShipConsRefDTO, 0)

	loadCommoditiesRes := make([]LoadCommoditiesDTO, 0)

	var shipConsRefQuery = BtmsLoadShipConsRefQuery + loadDetailsVar.LoadId

	var commodityQuery = BtmsLoadCommodityQuery + loadDetailsVar.LoadId

	err := dbclient.Select(&LoadShipConsRefRes, shipConsRefQuery)

	if err != nil {
		logger.Error("Error selecting values in the database:" + err.Error())
		return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err.Error())
	}

	// returning empty response for invalid loadId
	if len(LoadShipConsRefRes) == 0 {
		return LoadShipConsRefRes, nil, nil, nil, nil, nil
	}

	err1 := dbclient.Select(&loadCommoditiesRes, commodityQuery)

	if err1 != nil {
		logger.Error("Error selecting values in the database:" + err1.Error())
		return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+err1.Error())
	}

	if LoadShipConsRefRes[0].LoadMethod == "TL" || LoadShipConsRefRes[0].LoadMethod == "Intermodal" {

		trackingEventsUpdateRes := make([]EventTrackingUpdatesTL, 0)

		var trackingEventsQuery = LoadEventsUpdateQueryTL + loadDetailsVar.LoadId + ` ORDER BY event_date DESC`

		eventErr := dbclient.Select(&trackingEventsUpdateRes, trackingEventsQuery)

		if eventErr != nil {
			logger.Error("Error selecting values in the database:" + eventErr.Error())
			return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+eventErr.Error())
		}

		locationBreadCrumbsRes := make([]LocationBreadCrumbs, 0)

		var locationUpdatesQuery = LocationBreadCrumbsQuery + loadDetailsVar.LoadId + ` ORDER BY event_date DESC`

		locationErr := dbclient.Select(&locationBreadCrumbsRes, locationUpdatesQuery)

		dbclient.Close()

		if locationErr != nil {
			logger.Error("Error selecting values in the database:" + locationErr.Error())
			return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+locationErr.Error())
		}

		return LoadShipConsRefRes, loadCommoditiesRes, trackingEventsUpdateRes, locationBreadCrumbsRes, nil, nil

	} else if LoadShipConsRefRes[0].LoadMethod == "ELTL" {
		trackingEventsUpdateRes := make([]EventTrackingUpdatesELTL, 0)

		var trackingEventsQuery = LoadEventsUpdateQueryELTL + ` ORDER BY 1 ASC`

		loadId := loadDetailsVar.LoadId

		eventErr := dbclient.Select(&trackingEventsUpdateRes, trackingEventsQuery, loadId, loadId, loadId, loadId, loadId)

		dbclient.Close()

		if eventErr != nil {
			logger.Error("Error selecting values in the database:" + eventErr.Error())
			return nil, nil, nil, nil, nil, errs.ValidateResponse(nil, 500, "Error selecting values in the database: "+eventErr.Error())
		}

		return LoadShipConsRefRes, loadCommoditiesRes, nil, nil, trackingEventsUpdateRes, nil

	}

	dbclient.Close()

	return LoadShipConsRefRes, loadCommoditiesRes, nil, nil, nil, nil
}

func NewLoadDetailsRepositoryDb() RepositoryDb {
	return RepositoryDb{}
}
