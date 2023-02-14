package domain

import (
	"fmt"
	"golang/database"
	"golang/dto"
	"golang/errs"
	"golang/logger"
	"log"
	"strconv"
	"strings"
)

// BGLD_ps_1.3-1.7
// This Function will be used to verify the UserId, Login Id and whether they have permission to access bookLtl.
func (d RepositoryDb) VerifyAuthDataRepo(authDataReq *dto.AuthDataReq) (*VerifyAuthRes, *errs.AppErrorvalidation) {

	var Res VerifyAuthRes
	//errorsArr := make([]dto.ValidateResponse, 0)

	dbClient, DbError := database.GetCpdbClient()

	if DbError != nil {
		fmt.Println("entered error")
		return nil, DbError
	}

	verifyAuthDataResp := dbClient.Get(&Res, `SELECT users.user_id as "userId", user_permission.mysunteck_login_id as "sunteckLoginId",
                              user_permission.cust_master_id as "custmId", permission.permission_type as "customerPermission"    FROM cpdb.users 
                              INNER JOIN cpdb.user_permission ON users.user_id=user_permission.user_id
                              INNER JOIN cpdb.permission ON user_permission.permission_id=permission.permission_id
                              INNER JOIN cpdb.customer_permission ON customer_permission.global_customer_id=user_permission.global_customer_id AND customer_permission.permission_id=user_permission.permission_id 
                              WHERE user_permission.active = 1 AND customer_permission.active = 1 AND users.user_id = $1 AND user_permission.mysunteck_login_id = $2 AND permission.permission_type = 'Book LTL'`, authDataReq.UserId, authDataReq.LoginId)

	dbClient.Close()

	if verifyAuthDataResp != nil {
		logger.Error("Error while getting values from database:" + verifyAuthDataResp.Error())
		//errorsArr = append(errorsArr, dto.ValidateResponse{Code: "QLTL062", Message: "Unexpected database error"})
		//return nil, errs.UnexpectedErrorResponse(errorsArr)
		return nil, errs.ValidateResponse(nil, 500, "Error while trying to connect to Customer Portal database")

	}
	return &Res, nil
}

// BGLD_ps_1.3-1.11
// This function will get all the customer details like custmId and CustaId and Office code etc... from db
func (d RepositoryDb) GetCustDataRepo(loginId string) (*CustDataRes, *errs.AppError) {

	var Res CustDataRes
	errorsArr := make([]dto.ValidateResponse, 0)

	dbClient, DbError := database.BTMSDBWriterConnection()

	if DbError != nil {
		fmt.Println("entered error")
		return nil, DbError
	}

	custDataResp := dbClient.Get(&Res, `SELECT cml.custm_id, cml.custa_id, cml.office_code, 
                                       CASE WHEN (cm.use_banyan = 1 or (cm.use_mode = 1 and cm.use_banyan = 1)) THEN 'Banyan' ELSE 'Mode' END as "rating_engine" 
                                       from cust_mysunteck_logins cml
                                       inner join cust_master cm on cm.id = cml.custm_id
                                       where login = ?`, loginId)

	dbClient.Close()

	if custDataResp != nil {
		logger.Error("Error while getting values from database:" + custDataResp.Error())
		errorsArr = append(errorsArr, dto.ValidateResponse{Code: "QLTL062", Message: "Unexpected database error"})
		return nil, errs.NewUnexpectedError("BL029", "Unexpected database error")
	}

	return &Res, nil
}

// BGLD_ps_1.3-1.13
// The rating engine whether it is mode or banyan will be known by CheckRatingEngine(bookReq.CustmId)
func (d RepositoryDb) CheckRatingEngine(CustmId string) (*CheckRatingEngineResponse, *errs.AppError) {

	dbclnt, DbError := database.GetSharedDbClient()

	if DbError != nil {
		return nil, DbError
	}

	ratingEngine := "select rating_engine from btmsdb.quote_customer_details where customer_id =" + CustmId

	var ratingEngineVar CheckRatingEngineResponse
	custaIdResp := dbclnt.Get(&ratingEngineVar, ratingEngine)

	dbclnt.Close()

	if custaIdResp != nil {
		fmt.Println("entered custaID")
		return &ratingEngineVar, nil
	}
	return &ratingEngineVar, nil
}

// BGLD_ps_1.3-1.16
// The QuoteDetails() will be called when Rating engine is "Mode"
// In this function it gets the xml response of quotedetails to compare the rquest with xml Response for valid load creation
func (d RepositoryDb) QuoteDetails(loadId int) (*QuoteDetailsResponse, *errs.AppError) {

	dbclnt, DbError := database.GetSharedDbClient()

	if DbError != nil {
		return nil, DbError
	}

	QuoteDetails := "select xml_traffic from btmsdb.quote_xml where quote_load_details_id = (select quote_load_details_id from btmsdb.quote_load_details where banyan_load_id =" + strconv.Itoa(loadId) + ")"

	var QuoteDetailsVar QuoteDetailsResponse
	QuoteDetailsResp := dbclnt.Get(&QuoteDetailsVar, QuoteDetails)
	//fmt.Println("QuoteDetailsVarshshQuoteDetailsResp", QuoteDetailsResp)

	dbclnt.Close()

	if QuoteDetailsResp != nil {

		return &QuoteDetailsVar, nil
	}
	return &QuoteDetailsVar, nil
}

// BGLD_ps_1.20-1.27
// Valid load for banyan rating engine will be created with InsertBookDetailsModeRepo()
func (d RepositoryDb) InsertBookDetailsRepo(banres *dto.BanyanGetLoadDetailsDto, bookreq *dto.BookLTLRequestDTO, SelectedAccessorials string) (*dto.BookLTLResponseDTO, *errs.AppError) {
	var err error
	//var responseVar string
	DBClient, _ := database.BTMSDBConnection()
	tx, rollError := DBClient.Beginx()
	if rollError != nil {
		fmt.Println("error in dbconnection CreateUserRepo")
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	}
	fmt.Println("cncncnc", bookreq.CustaId)
	CustomerInfoQuery := `select  cust_agent.salesperson as "salesperson",cust_master.name as "name",cust_agent.id as "custaId", cust_agent.custm_id as "custmId", agents.tel as "agentPhoneNumber", agents.name as "agentName", agents.email as "agentEmail", agents.office as "officeCode", IFNULL(SUM(loadsh_balance.current_amount), 0) AS "daxBalance", IFNULL(cust_bal_temp.fats_balance, 0) AS "btmsBalance",agents.id as "salesId",cust_agent.contact as "contact",cust_agent.phone as "cusPhone",COALESCE(cust_agent.email, '')  as "cusEmail",COALESCE(cust_agent.acct_manager_id, '0')  as 'accountmanagerid',COALESCE(cust_agent.acct_manager, '') as 'accountmanager' ,cust_agent.fax as "cusFax" from agents inner join cust_agent on cust_agent.sales_id=agents.id inner join cust_master on cust_master.id=cust_agent.custm_id inner join cust_bal_temp on cust_master.id=cust_bal_temp.custm_id left JOIN loadsh_balance on cust_bal_temp.custm_id = loadsh_balance.custm_id where cust_agent.id ='` + bookreq.CustaId + `'`

	fmt.Println(CustomerInfoQuery, "scsc")
	var CustDetailsVar CustmerDetails
	CustDetailsResp := DBClient.Get(&CustDetailsVar, CustomerInfoQuery)

	if CustDetailsResp != nil {
		logger.Error("uhghjfbkdois" + CustDetailsResp.Error())
		fmt.Println("", CustDetailsVar.Name)
	}
	fmt.Println(CustDetailsVar, "asdfgCustDetailsVar")
	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}
	floor_plan_success := 1
	floor_plan_inches := 12
	floor_plan_inches_avg := 6
	floor_plan_inches_user := 0
	equip := "v"
	banyan_rating_complete := 1
	paytype := "0"
	load_ltl_quote_requestQuery := `INSERT INTO load_ltl_quote_request
	      SET cust_name='` + CustDetailsVar.Name + `',
	          custm_id='` + bookreq.CustmId + `',
	          custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	          salesperson='',
	          sales_id=` + strconv.Itoa(CustDetailsVar.SalesId) + `,
	          office='` + CustDetailsVar.OfficeCode + `',
	          created=NOW(),
	          ship_zip='` + bookreq.ShipperDetails.ShipZipCode + `',
	          ship_state='` + bookreq.ShipperDetails.ShipState + `',
	          ship_city='` + bookreq.ShipperDetails.ShipCity + `',
	          ship_limited='',
	          cons_zip='` + bookreq.ConsigneeDetails.ConsZipCode + `',
	          cons_state='` + bookreq.ConsigneeDetails.ConsState + `',
	          cons_city='` + bookreq.ConsigneeDetails.ConsCity + `',
	          cons_limited='',
	          floor_plan_success='` + strconv.Itoa(floor_plan_success) + `',
	          floor_plan_inches='` + strconv.Itoa(floor_plan_inches) + `',
	          floor_plan_inches_avg='` + strconv.Itoa(floor_plan_inches_avg) + `',
	          floor_plan_inches_user='` + strconv.Itoa(floor_plan_inches_user) + `',
	          banyan_rating_complete='` + strconv.Itoa(banyan_rating_complete) + `',
	          banyan_load_id='` + strconv.Itoa(bookreq.LoadId) + `',
	          banyan_rerate_reason='',
	          equip='` + equip + `',
	          paytype='` + paytype + `',
	          fee_codes='', notes='a:0:{}'`
	//'GUA,SEP,NOT'

	_, err = tx.Query(load_ltl_quote_requestQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	for i, _ := range bookreq.Commodities.Commodities {
		stackable := 0
		hazmat := 0
		if bookreq.Commodities.Commodities[i].Stackable == true {
			stackable = 1
		}
		if bookreq.Commodities.Commodities[i].Hazmat == true {
			stackable = 1
		}

		load_ltl_quote_request_items_Query := `INSERT INTO load_ltl_quote_request_items
		(load_ltl_quote_request_id, descrip, pallet_length, pallet_width,pallet_height,weight,nmfc,class,qty, load_item_id,stackable,
		  hazmat,density_item,locked,preapproved)
		VALUES( (select id from load_ltl_quote_request where banyan_load_id ='` + strconv.Itoa(bookreq.LoadId) + `' limit 1), '` + bookreq.Commodities.Commodities[i].Desc + `', '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `', '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Weight) + `','` + bookreq.Commodities.Commodities[i].NMFC + `','` + bookreq.Commodities.Commodities[i].Class + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Quantity) + `','` + strconv.Itoa(i) + `','` + strconv.Itoa(stackable) + `','` + strconv.Itoa(hazmat) + `', 0, '', '0')`

		_, err = tx.Query(load_ltl_quote_request_items_Query)

		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database 33", "500")
		} else {
			fmt.Println("success-Query33")
		}
	}
	banyanQuoteId := 0
	banyanQuoteNumber := ""
	CarrLtlId := ""
	CarrId := ""
	CarrName := ""
	ltlFlag := ""
	CustRate := ""
	CarrRate := ""
	CustTotal := ""
	CarrTotal := ""
	TransitTime := ""
	TotalWeight := 0
	CustFeulSurCharge := ""
	CarrFeulSurCharge := ""
	SCAC := ""
	ServiceId := ""
	CarrNotes := ""

	for i, _ := range bookreq.Commodities.Commodities {

		TotalWeight += bookreq.Commodities.Commodities[i].Weight
	}
	for i, _ := range banres.LoadDetails[0].Quotes {

		if bookreq.QuoteId == banres.LoadDetails[0].Quotes[i].QuoteID {
			banyanQuoteId = banres.LoadDetails[0].Quotes[i].QuoteID
			banyanQuoteNumber = banres.LoadDetails[0].Quotes[i].QuoteNumber
			CustRate = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CustomerPrice.FreightCharge, 'g', 5, 64)
			CarrRate = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CarrierPrice.FreightCharge, 'g', 5, 64)
			CustTotal = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CustomerPrice.NetPrice, 'g', 5, 64)
			CarrTotal = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CarrierPrice.NetPrice, 'g', 5, 64)
			TransitTime = strconv.Itoa(banres.LoadDetails[0].Quotes[i].TransitTime)
			CarrFeulSurCharge = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CarrierPrice.FuelSurcharge, 'g', 5, 64)
			CustFeulSurCharge = strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CustomerPrice.FuelSurcharge, 'g', 5, 64)
			ServiceId = banres.LoadDetails[0].Quotes[i].ServiceID
			CarrNotes = banres.LoadDetails[0].Quotes[i].CarrierNote
			fmt.Println(banyanQuoteId, banyanQuoteNumber, CustRate, CarrRate, CustTotal, CarrTotal, TransitTime, CarrFeulSurCharge, CustFeulSurCharge, ServiceId, CarrNotes, "djdjdjdkslzmxx")

			fmt.Println("banyanQuoteNumber", banyanQuoteId)
			CarrIdLoadshQuery := ` SELECT L.id, L.name, L.notes, L.pallet_length, L.pallet_width, L.pallet_height, L.pallet_weight, L.ltl_flag, L.ban_supports204,
                      CASE WHEN L.carr_deleted=1 THEN 0 ELSE 1 END as carr_activated, COALESCE(B.id, C.id) as carr_id,
                      CASE WHEN C.status IN ('ACTIVE','CAUTION') THEN 1 ELSE 0 END as carr_usable,
                      CASE WHEN P.blacklist=1 THEN 'BLACKLISTED' WHEN P.blacklist=0 THEN 'FAVORITE' ELSE NULL END as carr_pref,
                      CASE WHEN (office_exclusivity IS NULL OR office_exclusivity LIKE CONCAT('%,',A.office,',%')) THEN 0 ELSE 1 END AS office_prohibited
                    FROM sunteck_fats.carriers_LTL L
                    INNER JOIN cust_agent CA ON CA.id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    INNER JOIN agents A ON A.id=CA.sales_id
                    LEFT OUTER JOIN sunteck_fats.carriers C ON C.scac=L.scac
                    LEFT OUTER JOIN sunteck_fats.carriers B ON B.scac=L.broker_scac AND L.broker_scac<>''
                    LEFT OUTER JOIN cust_agent_carr_prefs P ON P.ltl_carr=1 AND P.carr_id=L.id AND P.custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    WHERE L.scac='` + banres.LoadDetails[0].Quotes[i].SCAC + `'AND L.broker_scac=''`
			SCAC = banres.LoadDetails[0].Quotes[i].SCAC

			var CarrIdLoadshVar CarrIdDetails
			CarrIdLoadshResp := DBClient.Get(&CarrIdLoadshVar, CarrIdLoadshQuery)

			if CarrIdLoadshResp != nil {
				//	logger.Error("uhghjfbkdois" + CarrIdResp.Error())
				fmt.Println("", CustDetailsVar.Name)
			}
			fmt.Println(CarrIdLoadshVar, "acncnmcgCustDetailsVar")
			CarrName = CarrIdLoadshVar.Name
			CarrId = strconv.Itoa(CarrIdLoadshVar.CarrId)
			ltlFlag = strconv.Itoa(CarrIdLoadshVar.LtlFlag)
			CarrLtlId = strconv.Itoa(CarrIdLoadshVar.Id)
			///////////////////////////////////

		}

		fmt.Println("banyanQuoteNumber", banyanQuoteId)
		CarrIdQuery := ` SELECT L.id, L.name, L.notes, L.pallet_length, L.pallet_width, L.pallet_height, L.pallet_weight, L.ltl_flag, L.ban_supports204,
                      CASE WHEN L.carr_deleted=1 THEN 0 ELSE 1 END as carr_activated, COALESCE(B.id, C.id) as carr_id,
                      CASE WHEN C.status IN ('ACTIVE','CAUTION') THEN 1 ELSE 0 END as carr_usable,
                      CASE WHEN P.blacklist=1 THEN 'BLACKLISTED' WHEN P.blacklist=0 THEN 'FAVORITE' ELSE NULL END as carr_pref,
                      CASE WHEN (office_exclusivity IS NULL OR office_exclusivity LIKE CONCAT('%,',A.office,',%')) THEN 0 ELSE 1 END AS office_prohibited
                    FROM sunteck_fats.carriers_LTL L
                    INNER JOIN cust_agent CA ON CA.id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    INNER JOIN agents A ON A.id=CA.sales_id
                    LEFT OUTER JOIN sunteck_fats.carriers C ON C.scac=L.scac
                    LEFT OUTER JOIN sunteck_fats.carriers B ON B.scac=L.broker_scac AND L.broker_scac<>''
                    LEFT OUTER JOIN cust_agent_carr_prefs P ON P.ltl_carr=1 AND P.carr_id=L.id AND P.custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    WHERE L.scac='` + banres.LoadDetails[0].Quotes[i].SCAC + `'AND L.broker_scac=''`

		//var CustDetailsVar CustmerDetails
		//_, err = DBClient.Get(CustomerInfoQuery, &CustDetailsVar)
		fmt.Println(CarrIdQuery)
		var CarrIdVar CarrIdDetails
		CarrIdResp := DBClient.Get(&CarrIdVar, CarrIdQuery)
		//CustomerInfoResp := DBClient.Get(&CustDetailsVar, CustomerInfoQuery)
		//fmt.Println(CustomerInfoResp, "asdassasfgCustDetailsVar")

		if CarrIdResp != nil {
			//	logger.Error("uhghjfbkdois" + CarrIdResp.Error())
			fmt.Println("", CustDetailsVar.Name)
		}
		fmt.Println(CarrIdVar, "acncnmcgCustDetailsVar")

		tx.Commit()
		tx, rollError = DBClient.Beginx()

		Query34 := `INSERT INTO load_ltl_quote_request_results
		(
		  load_ltl_quote_request_id,
		  carriers_LTL_id,
		  quote_unique_id,
		  load_ltl_quote_request_spot_result_id,
		  banyan_fee,
		  raw_soap_result
		  )
		VALUES
		(
		  (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' order by id desc limit 1),
		  '` + strconv.Itoa(CarrIdVar.Id) + `',
		  '` + banres.LoadDetails[0].Quotes[i].QuoteNumber + `',
		  NULL,
		  '0',
		  'NULL'
		)`

		_, err = tx.Query(Query34)

		if err != nil {
			tx.Rollback()
			//logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database 34", "500")
		} else {
			fmt.Println("success-Query34", "iteration", i)
		}
		tx.Commit()
		tx, rollError = DBClient.Beginx()
		fee_codes_banyan_ltl_quote_mappings := `INSERT INTO fee_codes_banyan_ltl_quote_mappings
		   (
		     ltl_quote_result_id,
		     fee_codes_banyan_mapping_id
		     )
		   VALUES
		   (
		     (select id from load_ltl_quote_request_results where quote_unique_id =  '` + banres.LoadDetails[0].Quotes[i].QuoteNumber + `'  and load_ltl_quote_request_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' order by id desc limit 1) order by id desc limit 1),'0')`
		fmt.Println(fee_codes_banyan_ltl_quote_mappings)
		_, err = tx.Query(fee_codes_banyan_ltl_quote_mappings)

		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query37", "iteration", i)
		}

	}

	loadshquery := `INSERT INTO loadsh
	  SET created      = NOW(),
	      created_by   = '',
	      created_id   = '',
	      date         = NOW(),
	      lastmod      =  NOW(),
	      never_edited = 1,
	      lastuser     = '',
	      carr_count   = 1,
	      stop_count   = 2,
	      miles        = '',
	      fk_auto_load_status_update = '1',
	      ltl_flag     = '` + ltlFlag + `',
	      ltl_indirect = 1,
	      ltl_foreign_quote_id = '` + banyanQuoteNumber + `',
	      ltl_foreign_quote_number = '` + strconv.Itoa(banyanQuoteId) + `',
	      ltl_quote_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' limit 1),
	      ltl_quote_result_id = (select id from load_ltl_quote_request_results where quote_unique_id =  '` + banyanQuoteNumber + `'  and load_ltl_quote_request_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' limit 1)limit 1),
	      ltl_quote_spot_result_id = NULL,
	      ltl_carr_contact = 'NULL',
	      edispatch_load_id = '` + strconv.Itoa(bookreq.LoadId) + `',
	      banyan_fee   = '0',
	      status       = 'QUOTE',
	      dispatcher   = '',
	      disp_id      = '',
	      salesperson  = '` + CustDetailsVar.Salesperson + `',
	      sales_id     = '` + strconv.Itoa(CustDetailsVar.SalesId) + `',
	      acct_manager = '` + CustDetailsVar.AccountManager + `',
	      acct_manager_id ='` + CustDetailsVar.AccountManagerId + `',
	      cust_name    = '` + CustDetailsVar.Name + `',
	      custa_id     ='` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	      custm_id     = '` + bookreq.CustmId + `',
	      office       = '` + CustDetailsVar.OfficeCode + `',
	      ship_city    = '` + bookreq.ShipperDetails.ShipCity + `',
	      ship_state   ='` + bookreq.ShipperDetails.ShipState + `',
	      ship_zip     ='` + bookreq.ShipperDetails.ShipZipCode + `',
	      ship_country = '` + bookreq.ShipperDetails.ShipCountry + `',
	      cons_city    = '` + bookreq.ConsigneeDetails.ConsCity + `',
	      cons_state   = '` + bookreq.ConsigneeDetails.ConsState + `',
	      cons_zip     = '` + bookreq.ConsigneeDetails.ConsZipCode + `',
	      carr_eq      = 'V',
	      cons_country =  '` + bookreq.ConsigneeDetails.ConsCountry + `',
	      carr_name    = '` + CarrName + `',
	      carr_id      = '` + CarrId + `',
	      cust_rate    = '` + CustRate + `',
	      carr_rate    ='` + CarrRate + `',
	      cust_total   = '` + CustTotal + `',
	      carr_total   = '` + CarrTotal + `',
	      linear_feet  = '',
	      ltl_volume_quote = 0,
	      ltl_transit_days =  '` + TransitTime + `',
	      load_qty     = '` + strconv.Itoa(bookreq.Commodities.Commodities[0].Quantity) + `',
	      load_type    = '` + TransitTime + `',
	      load_weight  = '` + strconv.Itoa(TotalWeight) + `',
	      load_method  = 'ELTL',
	      paytype      = '0',
	      share_amt    = '',
	      share_calc   = '',
	      share_agent_id = '',
		  ship_bl='` + bookreq.BlNumber + `',
          cust_po='` + bookreq.PoNumber + `',
          cust_shipid='` + bookreq.ShippingNumber + `'`

	////miles
	_, err = tx.Query(loadshquery)
	tx.Commit()
	fmt.Println(loadshquery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query61")
	}
	tx, rollError = DBClient.Beginx()
	commissionsQuery := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 post_id,
	 adjustment_type_id,
	 adjustment_payment_id,
	 parent_commission_id,
	 billing_adjustment_id,
	 note,
	 date,
	 total,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT
	 CASE
	     WHEN cg.internal=1 THEN 'INTERNAL_SHARE'
	     ELSE 'SHARE'
	 END                 AS 'type',
	 'OPEN'             AS 'status',
	 o.commission_basis  AS 'basis',
	 a.office            AS 'office',
	 l.id                AS 'loadsh_id',
	 cg.agent_id         AS 'agent_id',
	 cg.id               AS 'group_id',
	 NULL                AS 'post_id',
	 NULL                AS 'adjustment_type_id',
	 NULL                AS 'adjustment_payment_id',
	 NULL                AS 'parent_commission_id',
	 NULL                AS 'billing_adjustment_id',
	 NULL                AS 'note',
	 NULL                AS 'date',
	 NULL                AS 'total',
	 cg.amount           AS 'amount',
	 cg.calc             AS 'calc',
	 '510110'            AS 'gl_code'
	FROM loadsh l
	JOIN offices o
	 ON o.code = l.office
	LEFT JOIN commission_groups cg
	 ON cg.load_field_name = 'custa_id'
	 AND l.custa_id = cg.load_field_value
	LEFT JOIN agents a
	 ON a.id = cg.agent_id
	WHERE l.id = (select id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1)
	AND cg.id IS NOT NULL
	AND  a.id IS NOT NULL
	AND IFNULL(cg.active, 1) > 0
	AND IFNULL(cg.incentive_id, 0) < 1 -- exclude incentives
	)`

	log.Println("check 1 62")
	fmt.Println(commissionsQuery)
	_, err = tx.Query(commissionsQuery)

	if err != nil {
		//tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query62")
	}

	commissions_internalQuery := `INSERT INTO commissions_internal
	(
	 type,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 commission_id,
	 parent_internal_id,
	 amount,
	 calc,
	 deleted
	)
	(
	SELECT
	 'ORIGINAL'               AS 'type',
	 a.office            AS 'office',
	 l.id                AS 'loadsh_id',
	 cg.agent_id         AS 'agent_id',
	 cg.id               AS 'group_id',
	 (SELECT MAX(c2.id) FROM commissions c2 WHERE c2.loadsh_id = l.id AND c2.type = 'MAIN') AS 'commission_id',
	 NULL                AS 'parent_internal_id',
	 cg.amount           AS 'amount',
	 cg.calc             AS 'calc',
	 0                   AS 'deleted'
	FROM loadsh l
	JOIN offices o
	 ON o.code = l.office
	LEFT JOIN commission_groups cg
	 ON cg.load_field_name = 'custa_id'
	 AND l.custa_id = cg.load_field_value
	LEFT JOIN agents a
	 ON a.id = cg.agent_id
	WHERE l.id = (select id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1)
	AND cg.id IS NOT NULL
	AND  a.id IS NOT NULL
	AND IFNULL(cg.active, 1) > 0
	AND IFNULL(cg.internal, 0) > 0 -- internal only
	)`

	_, err = tx.Query(commissions_internalQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query63")
	}

	loadshIdQuery := `select id , ltl_quote_id ,carr_id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1`
	fmt.Println("chdhchc", loadshIdQuery)
	var loadshIdVar LoadShIdStruct
	loadshIdResp := DBClient.Get(&loadshIdVar, loadshIdQuery)

	if loadshIdResp != nil {

		fmt.Println("", CustDetailsVar.Name)
	}
	fmt.Println(loadshIdVar, "loadshIdVar")

	loadshUpdate := ` UPDATE loadsh
	LEFT OUTER JOIN (
	  SELECT loadsh_id, amount, calc, agent_id
	  FROM commissions
	  WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND status <> 'DELETED'
	  AND type = 'SHARE'
	  ORDER BY id ASC
	  LIMIT 1
	) TT ON TT.loadsh_id = loadsh.id
	SET loadsh.share_amt = IFNULL(TT.amount, 0),
	  loadsh.share_calc = IFNULL(TT.calc, 'USD'),
	  loadsh.share_agent_id = IFNULL(TT.agent_id, 0)
	WHERE loadsh.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(loadshUpdate)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query64")
	}
	log.Println("check 2 exected")

	tx.Commit()
	tx, rollError = DBClient.Beginx()
	commissionCalc := `INSERT INTO commissions (type, calc, amount, office, loadsh_id, agent_id, basis)
	SELECT
	 'MAIN' AS type,
	 'PCT'  AS calc,
	 CAST(o.commission_pct AS DECIMAL(10,2) ) AS amount,
	 l.office,
	 l.id AS loadsh_id,
	 l.sales_id AS agent_id,
	 o.commission_basis AS basis
	FROM loadsh AS l
	JOIN offices AS o ON l.office = o.code
	LEFT JOIN commissions c ON l.id = c.loadsh_id AND c.type = 'MAIN'
	/* ONLY 1 active MAIN record per load */
	WHERE l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND IFNULL((SELECT 1 FROM commissions c2 WHERE c2.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' AND c2.type = 'MAIN' AND c2.status <> 'DELETED' LIMIT 1), 0) < 1`

	_, err = tx.Query(commissionCalc)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query65")
	}

	CommissionCQuery := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	JOIN offices o ON l.office = o.code
	JOIN agents a ON c.agent_id = a.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	 SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	 FROM lscarr tc
	 JOIN lsstops ts ON ts.lscarr_id = tc.id
	 JOIN lsfees tf ON tf.lsstop_id = ts.id
	 WHERE tf.excluded_from_commissions = 1
	 AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	SET
	c.office = a.office,
	c.basis = o.commission_basis,
	c.date = NULLIF(
	     (CASE
	         WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'invoice_date'  THEN IF( NULLIF(l.invoice_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.invoice_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'delfinal_date' THEN IF( NULLIF(l.delfinal_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.delfinal_date), DATE(NOW())),  NULL )
	         ELSE NULL
	     END)
	 , '0000-00-00'),
	 c.total = CAST( IF( c.calc = 'PCT',
	     ROUND(LOAD_COMMISSION(
	         (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	         (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	         l.broker_liability_insurance,
	         0.00, /* charge GPS only to MAIN */
	         l.banyan_fee,
	         c.amount/100
	     ),2),
	     c.amount )
	 AS DECIMAL(12,2) ),
	 c.adjustment_payment_id = NULL,
	 c.cust_total = l.cust_total,
	 c.carr_total = l.carr_total,
	 c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	 c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND ( c.type = 'SHARE' OR c.type = 'INTERNAL_SHARE' )
	/* only update if commission date is in the future or TBD */
	AND (
	     NULLIF(c.date, '0000-00-00') IS NULL
	     OR c.date > DATE(NOW())
	 )`

	_, err = tx.Query(CommissionCQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query66")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	Query67 := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	 SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	 FROM lscarr tc
	 JOIN lsstops ts ON ts.lscarr_id = tc.id
	 JOIN lsfees tf ON tf.lsstop_id = ts.id
	 WHERE tf.excluded_from_commissions = 1
	 AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	JOIN offices o ON l.office = o.code
	/* FLAT USD SHAREs only joined to MAIN type */
	LEFT JOIN
	(
	 SELECT
	     IFNULL(SUM(amount), 0) amount,
	     loadsh_id
	 FROM commissions
	 WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND type = 'SHARE'
	 AND calc = 'USD'
	 AND status <> 'DELETED'
	) flat_shares ON l.id = flat_shares.loadsh_id
	/* PCT Based SHAREs 'shared' from MAIN type */
	LEFT JOIN
	(
	 SELECT
	     IFNULL(SUM(amount), 0) amount,
	     loadsh_id
	 FROM commissions
	 WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND type = 'SHARE'
	 AND calc = 'PCT'
	 AND status <> 'DELETED'
	) shares ON l.id = shares.loadsh_id
	SET
	 c.office = l.office,
	 c.agent_id = IF( l.office = ( SELECT office FROM agents WHERE id = l.sales_id ), l.sales_id, o.owner_id),
	 c.date = NULLIF(
	     (CASE
	         WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'invoice_date'  THEN l.invoice_date
	         WHEN o.commission_basis = 'delfinal_date' THEN l.delfinal_date
	         ELSE NULL
	     END)
	 , '0000-00-00'),
	 c.amount = o.commission_pct,
	 c.basis = o.commission_basis,
	 c.calc = 'PCT',
	 c.total = CAST(
	     ROUND(LOAD_COMMISSION(
	         (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	         (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	         l.broker_liability_insurance,
	         l.gps_total,
	         l.banyan_fee,
	         /* subtract percentage based SHAREs from MAIN */
	         ( o.commission_pct - IFNULL(shares.amount, 0) )/100
	     ),2)
	     /* subtract FLAT USD shares from MAIN */
	     - IFNULL(flat_shares.amount,0)
	 AS DECIMAL(12,2) ),
	 c.cust_total = l.cust_total,
	 c.carr_total = l.carr_total,
	 c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	 c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND c.status <> 'DELETED'
	AND c.type = 'MAIN'
	/* only update if commission date is in the future or TBD */
	AND (
	     NULLIF(c.date, '0000-00-00') IS NULL
	     OR c.date > DATE(NOW())
	 )`

	_, err = tx.Query(Query67)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query67")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	log.Println("check 3 ")

	commissionsInsert := `INSERT INTO commissions
	(type, status, basis, office, loadsh_id, agent_id, group_id, amount, calc, gl_code, date, total)
	(
	 SELECT * FROM (
	     SELECT
	         pivot.type,
	         'OPEN' status,
	         o.commission_basis,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             l.office,
	             IFNULL(share_agent.office, a.office)
	         ) office,
	         l.id loadsh_id,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             o.owner_id,
	             IFNULL(share.agent_id, ca.sales_id)
	         ) agent_id,
	         cg.id group_id,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             cg.office_amount,
	             cg.sales_amount * IFNULL(share.amount / all_shares.total, 1.00)
	         ) amount,
	         'PCT' calc,
	         cg.gl_code,
	         DATE(main.date) date,
	         LOAD_COMMISSION(
	             (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	             (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	             l.broker_liability_insurance,
	             0, /* charge GPS only to MAIN */
	             l.banyan_fee,
	             IF(pivot.type = 'OFFICE_INCENTIVE',
	                 cg.office_amount,
	                 cg.sales_amount * IFNULL(share.amount/all_shares.total, 1.00)
	             )/100
	         ) AS total
	     FROM commission_incentives ci
	     JOIN commission_groups cg ON ci.id = cg.incentive_id
	     JOIN loadsh l ON (
	     (cg.load_field_name = 'office' AND l.office = cg.load_field_value)
	         OR
	     (cg.load_field_name = 'custa_id' AND l.custa_id = cg.load_field_value)
	     )
	     JOIN offices o ON l.office = o.code
	     /* Fees to be removed from profit */
	     LEFT JOIN (
	         SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	         FROM lsfees tf
	         LEFT JOIN lsstops ts ON tf.lsstop_id = ts.id
	         LEFT JOIN lscarr tc ON ts.lscarr_id = tc.id
	         WHERE tf.excluded_from_commissions = 1
	         
	         
	         AND tc.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	         GROUP BY tc.loadsh_id
	     ) excluded_fees ON l.id = excluded_fees.loadsh_id
	     JOIN cust_agent ca ON l.custa_id = ca.id
	     -- join to MAIN in order to get commissions date of MAIN commission
	     JOIN (
	         SELECT date, loadsh_id
	         FROM commissions
	         WHERE type = 'MAIN'
	         AND NULLIF(date, '0000-00-00') IS NOT NULL
	         AND NULLIF(post_id, 0) IS NULL
	         AND status <> 'DELETED'
	         AND loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	     ) main ON l.id = main.loadsh_id
	     JOIN agents a ON ca.sales_id = a.id
	     JOIN (
	         SELECT 'OFFICE_INCENTIVE' AS type
				UNION ALL
				SELECT 'SALES_INCENTIVE'
			) pivot
	     -- get sum of all PCT Load Default Share amounts, so that we can figure out proportionate split of SALES_INCENTIVE percentage
	     LEFT JOIN (
	         SELECT loadsh.id as load_id, SUM(all_shares.amount) as total
	         FROM loadsh
	         INNER JOIN commission_groups as all_shares
	             ON IFNULL(all_shares.incentive_id, 0) = 0
	             AND all_shares.load_field_name = 'custa_id'
	             AND all_shares.load_field_value = loadsh.custa_id
	             AND all_shares.active = 1
	             AND all_shares.internal = 0
	             AND all_shares.calc = 'PCT'
	         WHERE loadsh.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	         ) all_shares ON all_shares.load_id = l.id
	     -- for sales incentives only, join to any PCT Load Default Share for this customer;
	     --  if none exist, larger SELECT will still return one SALES_INCENTIVE record because this is an outer join;
	     --  but if any _do_ exist, larger select will return one SALES_INCENTIVE record _per_ Load Default Share
	     LEFT JOIN commission_groups as share
	         ON IFNULL(share.incentive_id, 0) = 0
	         AND pivot.type = 'SALES_INCENTIVE'
	         AND share.load_field_name = cg.load_field_name
	         AND share.load_field_value = cg.load_field_value
	         AND share.active = 1
	         AND share.internal = 0
	         AND share.calc = 'PCT'
	     -- join to agent record of share agents (so we can get office)
	     LEFT JOIN agents as share_agent
	         ON share_agent.id = share.agent_id
	     -- join to any existing SALES_INCENTIVE for this load
	     LEFT JOIN commissions c ON c.group_id = cg.id
	         AND c.type = 'SALES_INCENTIVE'
	         AND pivot.type = 'SALES_INCENTIVE'
	         AND c.loadsh_id = l.id
	     -- join to any existing OFFICE_INCENTIVE for this load
	     LEFT JOIN commissions c2 ON c2.group_id = cg.id
	         AND c2.type = 'OFFICE_INCENTIVE'
	         AND pivot.type = 'OFFICE_INCENTIVE'
	         AND c2.loadsh_id = l.id
	     WHERE ci.active > 0
	     -- MAIN commission on load has a date value set
	     AND NULLIF(main.date, '0000-00-00') IS NOT NULL
	     -- commission date is after effective date (if set)
	     AND ( NULLIF(cg.effective_date, '0000-00-00') IS NULL
	         OR cg.effective_date <= (CASE
	             WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	             ELSE main.date
	         END)
	     )
	     -- commission date is before expiry date (if set)
	     AND ( NULLIF(cg.expiry_date, '0000-00-00') IS NULL
	         OR cg.expiry_date >= (CASE
	             WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	             ELSE main.date
	         END)
	     )
	     -- SALES_INCENTIVE and/or OFFICE_INCENTIVE are applicable but do not already exist for this load
	     AND ( (pivot.type = 'SALES_INCENTIVE' AND ci.sales_share > 0 AND cg.sales_amount <> 0 AND c.id IS NULL)
	         OR (pivot.type = 'OFFICE_INCENTIVE' AND ci.office_share > 0 AND cg.office_amount <> 0 AND c2.id IS NULL)
	     )
	     AND cg.active > 0
	     AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 ) temp_new_incentives
	 WHERE temp_new_incentives.amount <> 0
	)`

	_, err = tx.Query(commissionsInsert)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query68")
	}

	commissionsQuery2 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'INTERNAL_SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsQuery2)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query69")
	}
	log.Println("check 5 ")

	commissionsQuery3 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsQuery3)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query70")
	}

	commissionsQuery4 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsQuery4)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query71")
	}
	log.Println("check 6 62")

	commissionsQuery5 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SALES_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsQuery5)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query72")
	}

	commissionsQuery6 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'MAIN'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsQuery6)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query73")
	}

	commissionsQuery7 := `update commission_incentives ci
	JOIN commission_groups cg ON ci.id = cg.incentive_id
	JOIN cust_agent ca ON cg.load_field_value = ca.id AND cg.load_field_name = 'custa_id'
	JOIN cust_master cm ON ca.custm_id = cm.id
	JOIN (
	 SELECT
	     cust_master.id as custm_id,
	     IF(commission_groups.basis = 'ship_date',
	         loadsh.ship_date,
	         CURDATE()
	     ) as eff_date
	 FROM commissions
	 JOIN commission_groups ON commissions.group_id = commission_groups.id
	 JOIN cust_agent ON commission_groups.load_field_value = cust_agent.id AND commission_groups.load_field_name = 'custa_id'
	 JOIN cust_master ON cust_agent.custm_id = cust_master.id
	 JOIN loadsh ON commissions.loadsh_id = loadsh.id
	 WHERE commissions.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND NULLIF(commission_groups.incentive_id, 0) IS NOT NULL
	 AND commission_groups.auto_expiry_months IS NOT NULL
	 AND NULLIF(commission_groups.effective_date, '0000-00-00') IS NULL
	 AND NULLIF(commission_groups.expiry_date, '0000-00-00') IS NULL
	 GROUP BY cust_master.id
	) TT ON cm.id = TT.custm_id
	SET
		cg.effective_date = TT.eff_date,
		cg.expiry_date = DATE_ADD(TT.eff_date, INTERVAL cg.auto_expiry_months MONTH)
	WHERE ci.active > 0
	AND cg.auto_expiry_months IS NOT NULL
	AND NULLIF(cg.effective_date, '0000-00-00') IS NULL
	AND NULLIF(cg.expiry_date, '0000-00-00') IS NULL`

	_, err = tx.Query(commissionsQuery7)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query75")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsQuery8 := `INSERT IGNORE INTO loadsh_nocapreload VALUES ('` + strconv.Itoa(loadshIdVar.LoadshId) + `', 0)`

	_, err = tx.Query(commissionsQuery8)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	}

	commissionsQuery9 := `UPDATE edispatch_doc_images
	               SET loadsh_id     = '` + strconv.Itoa(loadshIdVar.LoadshId) + `',
	                   document_type = 'Packing-Worksheet-L` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	               WHERE ltl_quote_id = '` + strconv.Itoa(loadshIdVar.LtlQuoteId) + `'`
	fmt.Println(commissionsQuery9)
	_, err = tx.Query(commissionsQuery9)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query77")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lscarrQuery := `INSERT INTO lscarr
	     SET created       = NOW(),
	         lastmod       = NOW(),
	         loadsh_id     = '` + strconv.Itoa(loadshIdVar.LoadshId) + `',
	         carr_id       = '` + strconv.Itoa(loadshIdVar.CarrId) + `',
	         carr_id_LTL   =  '` + CarrLtlId + `',
	         ltl_special_instruct='Quote #` + banyanQuoteNumber + ` ` + SelectedAccessorials + `',
	         carr_order    = 1,
	         carr_name     =  '` + CarrName + `',
	         cust_rate     = '` + CustRate + `',
	         carr_rate     = '` + CarrRate + `',
	         carr_total    = '` + CarrTotal + `',
	         miles_agent   = '',
	         miles_total   = '',
	         miles_total_questionable = 0,
	         trailer_length = '53',
	         equip         = 'V',
	         cargo_value_opt_id  = '',
	         weight_total  = '` + strconv.Itoa(TotalWeight) + `'`
	fmt.Println(lscarrQuery, SelectedAccessorials, "ssmsm")
	_, err = tx.Query(lscarrQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query82")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	// warning
	lsstopsPickQuery := `INSERT INTO lsstops
	     SET created = NOW(),
	         lastmod      = NOW(),
	         lscarr_id    = (select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1),
	         stop_order   = 1,
	         type         = 'PICK',
	         city         = '` + bookreq.ShipperDetails.ShipCity + `',
	         state        = '` + bookreq.ShipperDetails.ShipState + `',
	         zip          = '` + bookreq.ShipperDetails.ShipZipCode + `',
	         contact      = '` + bookreq.ShipperDetails.ShipContactName + `',
	         email          = '` + bookreq.ShipperDetails.ShipEmail + `',
	         bl          = '` + bookreq.BlNumber + `',
	         time          = '` + bookreq.ShipperDetails.ShipEarliestTime + `',
	         time2          = '` + bookreq.ShipperDetails.ShipLatestTime + `',
	         date          = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	         date2          = '` + bookreq.ShipperDetails.ShipLatestDate + `',	       
	         tel          = '` + bookreq.ShipperDetails.ShipPhone + `',	       
	         fax          = '` + bookreq.ShipperDetails.ShipFax + `',	       	       
	         instruct          = '` + bookreq.ShipperDetails.ShipLoadNotes + `',
	         addr1          = '` + bookreq.ShipperDetails.ShipAddress1 + `',
	         addr2          = '` + bookreq.ShipperDetails.ShipAddress2 + `',
	         name          = '` + bookreq.ShipperDetails.ShipName + `',
	         country      = '` + bookreq.ShipperDetails.ShipCountry + `'`

	_, err = tx.Query(lsstopsPickQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query83")
	}
	//tx.Commit()
	//tx, rollError = DBClient.Beginx()
	//	// warning
	lsstopsDropQuery := `INSERT INTO lsstops
	     SET created = NOW(),
	         lastmod      = NOW(),
	         lscarr_id    = (select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1),
	         stop_order   = 2,
	         type         = 'DROP',
	         city         = '` + bookreq.ConsigneeDetails.ConsCity + `',
	         state        = '` + bookreq.ConsigneeDetails.ConsState + `',
	         zip          ='` + bookreq.ConsigneeDetails.ConsZipCode + `',
	         contact      = '` + bookreq.ConsigneeDetails.ConsContactName + `',
	         email          = '` + bookreq.ConsigneeDetails.ConsEmail + `',
	         bl          = '` + bookreq.BlNumber + `',
	         time          = '` + bookreq.ConsigneeDetails.ConsEarliestTime + `',
	         time2          = '` + bookreq.ConsigneeDetails.ConsLatestTime + `',
	         date          = '` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	         date2          = '` + bookreq.ConsigneeDetails.ConsLatestDate + `',	       
	         tel          = '` + bookreq.ConsigneeDetails.ConsPhone + `',	       
	         fax          = '` + bookreq.ConsigneeDetails.ConsFax + `',	          
	         instruct          = '` + bookreq.ConsigneeDetails.ConsLoadNotes + `',
	         addr1          = '` + bookreq.ConsigneeDetails.ConsAddress1 + `',
	         addr2          = '` + bookreq.ConsigneeDetails.ConsAddress2 + `',
	         name          = '` + bookreq.ConsigneeDetails.ConsName + `',
	         country      = '` + bookreq.ConsigneeDetails.ConsCountry + `'`

	_, err = tx.Query(lsstopsDropQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query84")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	for i, _ := range bookreq.Commodities.Commodities {
		stackable := 0
		hazmat := 0
		if bookreq.Commodities.Commodities[i].Stackable == true {
			stackable = 1
		}
		if bookreq.Commodities.Commodities[i].Hazmat == true {
			stackable = 1
		}

		lsitemsQuery := `INSERT INTO lsitems
	       SET created    = NOW(),
	           lastmod    = NOW(),
	           lsstop_id  =(select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' order by id limit 1)  order by id limit 1),
	           item_order ='` + strconv.Itoa(i+1) + `',
	           qty        = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Quantity) + `',
	           type       =  '` + bookreq.Commodities.Commodities[i].EquipmentType + `',
	           descrip    = '` + bookreq.Commodities.Commodities[i].Desc + `',
	           density    =  '` + bookreq.Commodities.Commodities[i].Density + `',
	           pallet_length = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `',
	           pallet_width  = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `',
	           pallet_height = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `',
	           weight     ='` + strconv.Itoa(bookreq.Commodities.Commodities[i].Weight) + `',
	           stackable  = '` + strconv.Itoa(stackable) + `',
	           hazmat     ='` + strconv.Itoa(hazmat) + `',
	           class      = '` + bookreq.Commodities.Commodities[i].Class + `',
	           nmfc       = '` + bookreq.Commodities.Commodities[i].NMFC + `'`

		_, err = tx.Query(lsitemsQuery)

		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query85")
		}

		load_ltl_quote_approved_itemsQuery := `INSERT INTO load_ltl_quote_approved_items
	         SET custa_id   = '` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	             created    = NOW(),
	             zip        = '` + bookreq.ShipperDetails.ShipZipCode + `',
	             descrip    = '` + bookreq.Commodities.Commodities[i].Desc + `',
	             nmfc       = '` + bookreq.Commodities.Commodities[i].NMFC + `',
	             class      = '` + bookreq.Commodities.Commodities[i].Class + `',
	             pallet_length = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `',
	             pallet_width  = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `',
	             pallet_height ='` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `',
	             stackable     = '` + strconv.Itoa(stackable) + `',
	             hazmat        ='` + strconv.Itoa(hazmat) + `',
	             density_item  = 0,
	             density       = '` + bookreq.Commodities.Commodities[i].Density + `'`

		_, err = tx.Query(load_ltl_quote_approved_itemsQuery)

		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query86")
		}
	}

	lsfeesQuery := `INSERT INTO lsfees
	                     SET created     = NOW(),
	                         lastmod     = NOW(),
	                         lsstop_id   = (select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' order by id limit 1)  order by id limit 1), fee_order = '1', qty = '1', type = 'FUEL SURCHARGE', descrip = '', cust_rate_type = 'FLAT', cust_rate = 'NULL', cust_charge ='` + CustFeulSurCharge + `', carr_charge = '` + CarrFeulSurCharge + `', carr_rate_type = 'FLAT', carr_rate = 'NULL', code = 'FUE'`
	fmt.Println(lsfeesQuery)
	_, err = tx.Query(lsfeesQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query88")
	}

	for i, _ := range banres.LoadDetails[0].Quotes {

		if bookreq.QuoteId == banres.LoadDetails[0].Quotes[i].QuoteID {
			//strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CustomerPrice.FuelSurcharge, 'g', 5, 64)
			for j, _ := range banres.LoadDetails[0].Quotes[i].CustomerPrice.Charges {

				lsfeesQuery2 := `INSERT INTO lsfees
			                     SET created     = NOW(),
			                         lastmod     = NOW(),
			                         lsstop_id   =  (select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' order by id limit 1)  order by id limit 1), fee_order = '` + strconv.Itoa(j+2) + `', qty = '1', type = '` + banres.LoadDetails[0].Quotes[i].CustomerPrice.Charges[j].Name + `', descrip = '', cust_rate_type = 'FLAT', cust_rate = 'NULL', cust_charge = '` + strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CustomerPrice.Charges[j].Amount, 'g', 5, 64) + `', carr_charge =  '` + strconv.FormatFloat(banres.LoadDetails[0].Quotes[i].CarrierPrice.Charges[j].Amount, 'g', 5, 64) + `', carr_rate_type = 'FLAT', carr_rate = 'NULL', code = '` + banres.LoadDetails[0].Quotes[i].CarrierPrice.Charges[j].Code + `'`

				_, err = tx.Query(lsfeesQuery2)
				//if rollError != nil {
				//	fmt.Println("error in dbconnection CreateUserRepo")
				//	return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
				//}
				if err != nil {
					tx.Rollback()
					logger.Error("Error while inserting values in the database:" + err.Error())
					return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
				} else {
					fmt.Println("success-Query89")
				}
			}
		}
	}

	load_ltl_quote_requestQuery2 := `UPDATE load_ltl_quote_request SET loadsh_id='` + strconv.Itoa(loadshIdVar.LoadshId) + `' WHERE id='` + strconv.Itoa(loadshIdVar.LtlQuoteId) + `'`

	_, err = tx.Query(load_ltl_quote_requestQuery2)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query94")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	loadshUpdateQuery := `UPDATE loadsh
	SET
	lastmod = NOW(),
	lastuser = 'Customer',
	created_by = 'Customer',
	status = 'ACTIVE',
	dispatcher = '` + CustDetailsVar.Salesperson + `',
	disp_id ='` + strconv.Itoa(CustDetailsVar.SalesId) + `',
	office = '` + CustDetailsVar.OfficeCode + `',
	ship_name = '` + bookreq.ShipperDetails.ShipName + `',
	ship_date = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	cons_name = '` + bookreq.ConsigneeDetails.ConsName + `',
	cons_date ='` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	carr_name = '` + CarrName + `',
	cust_total_primary = '` + CustTotal + `',
	load_length = '53',
	miles = 'null',
	id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	WHERE id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(loadshUpdateQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query161")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lscarridQuery := `select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1`

	var lscarridVar LscarrIdStruct
	lscarridResp := DBClient.Get(&lscarridVar, lscarridQuery)

	if lscarridResp != nil {
		logger.Error("uhghjfbkdois" + lscarridResp.Error())
		fmt.Println("", lscarridVar)
	}
	fmt.Println(lscarridVar, "asdfgCustDetailsVar")
	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsidQuery := `select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'   limit 1)  order by id limit 1`

	var LsstopsidVar LsstopIdStruct
	LsstopsidResp := DBClient.Get(&LsstopsidVar, lsstopsidQuery)

	if LsstopsidResp != nil {
		logger.Error("uhghjfbkdois" + LsstopsidResp.Error())
		fmt.Println("", LsstopsidVar)
	}
	fmt.Println(LsstopsidResp, "asdfgCustDetailsVar")

	if err != nil {
		//tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}

	lsstopsiddropQuery := `select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1)  order by id limit 1`

	var lsstopsiddropVar LsstopIdStruct
	lsstopsiddropResp := DBClient.Get(&lsstopsiddropVar, lsstopsiddropQuery)

	if lsstopsiddropResp != nil {
		logger.Error("uhghjfbkdois" + lsstopsiddropResp.Error())
		fmt.Println("", lsstopsiddropVar)
	}
	fmt.Println(lsstopsiddropResp, "asdfgCustDetailsVar")

	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}

	lscarrQuery3 := `UPDATE lscarr
	SET
	lastmod =  NOW(),
	carr_name ='` + CarrName + `',
	miles_total = '',
	miles_agent = '',
	ltl_special_instruct = 'Quote #` + banyanQuoteNumber + ` ` + SelectedAccessorials + `',
	id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'
	WHERE id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'`

	fmt.Println(lscarrQuery3)
	_, err = tx.Query(lscarrQuery3)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query162")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsPickUpdateQuery := `UPDATE lsstops
	SET
	lastmod =  NOW(),
	name = '` + bookreq.ShipperDetails.ShipName + `',
	addr1 = '` + bookreq.ShipperDetails.ShipAddress1 + `',
	date = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	blind = '',
	showph = '',
	id = '` + strconv.Itoa(LsstopsidVar.LsstopId) + `'
	WHERE id = '` + strconv.Itoa(LsstopsidVar.LsstopId) + `'`
	fmt.Println(lsstopsPickUpdateQuery)
	_, err = tx.Query(lsstopsPickUpdateQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query163")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsDropUpdateQuery := `UPDATE lsstops
	SET
	lastmod =  NOW(),
	name ='` + bookreq.ConsigneeDetails.ConsName + `',
	addr1 = '` + bookreq.ConsigneeDetails.ConsAddress1 + `',
	date = '` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	blind = '',
	showph = '',
	id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'
	WHERE id ='` + strconv.Itoa(lscarridVar.LscarrId) + `'`

	_, err = tx.Query(lsstopsDropUpdateQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query164")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	loadshUpdatequery2 := `UPDATE loadsh
	LEFT OUTER JOIN (
	  SELECT loadsh_id, amount, calc, agent_id
	  FROM commissions
	  WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND status <> 'DELETED'
	  AND type = 'SHARE'
	  ORDER BY id ASC
	  LIMIT 1
	) TT ON TT.loadsh_id = loadsh.id
	SET loadsh.share_amt = IFNULL(TT.amount, 0),
	  loadsh.share_calc = IFNULL(TT.calc, 'USD'),
	  loadsh.share_agent_id = IFNULL(TT.agent_id, 0)
	WHERE loadsh.id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(loadshUpdatequery2)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query166")
	}
	//

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsQuery11 := `INSERT INTO commissions (type, calc, amount, office, loadsh_id, agent_id, basis)
	SELECT
	  'MAIN' AS type,
	  'PCT'  AS calc,
	  CAST(o.commission_pct AS DECIMAL(10,2) ) AS amount,
	  l.office,
	  l.id AS loadsh_id,
	  l.sales_id AS agent_id,
	  o.commission_basis AS basis
	FROM loadsh AS l
	JOIN offices AS o ON l.office = o.code
	LEFT JOIN commissions c ON l.id = c.loadsh_id AND c.type = 'MAIN'
	/* ONLY 1 active MAIN record per load */
	WHERE l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND IFNULL((SELECT 1 FROM commissions c2 WHERE c2.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' AND c2.type = 'MAIN' AND c2.status <> 'DELETED' LIMIT 1), 0) < 1`

	_, err = tx.Query(commissionsQuery11)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query167")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsQuery12 := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	JOIN offices o ON l.office = o.code
	JOIN agents a ON c.agent_id = a.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	  SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	  FROM lscarr tc
	  JOIN lsstops ts ON ts.lscarr_id = tc.id
	  JOIN lsfees tf ON tf.lsstop_id = ts.id
	  WHERE tf.excluded_from_commissions = 1
	  AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	SET
	 c.office = a.office,
	 c.basis = o.commission_basis,
	 c.date = NULLIF(
	      (CASE
	          WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'invoice_date'  THEN IF( NULLIF(l.invoice_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.invoice_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'delfinal_date' THEN IF( NULLIF(l.delfinal_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.delfinal_date), DATE(NOW())),  NULL )
	          ELSE NULL
	      END)
	  , '0000-00-00'),
	  c.total = CAST( IF( c.calc = 'PCT',
	      ROUND(LOAD_COMMISSION(
	          (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	          (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	          l.broker_liability_insurance,
	          0.00, /* charge GPS only to MAIN */
	          l.banyan_fee,
	          c.amount/100
	      ),2),
	      c.amount )
	  AS DECIMAL(12,2) ),
	  c.adjustment_payment_id = NULL,
	  c.cust_total = l.cust_total,
	  c.carr_total = l.carr_total,
	  c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	  c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND ( c.type = 'SHARE' OR c.type = 'INTERNAL_SHARE' )
	/* only update if commission date is in the future or TBD */
	AND (
	      NULLIF(c.date, '0000-00-00') IS NULL
	      OR c.date > DATE(NOW())
	  )`

	_, err = tx.Query(commissionsQuery12)
	log.Println("check 7 ")

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query168")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsQuery13 := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	  SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	  FROM lscarr tc
	  JOIN lsstops ts ON ts.lscarr_id = tc.id
	  JOIN lsfees tf ON tf.lsstop_id = ts.id
	  WHERE tf.excluded_from_commissions = 1
	  AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	JOIN offices o ON l.office = o.code
	/* FLAT USD SHAREs only joined to MAIN type */
	LEFT JOIN
	(
	  SELECT
	      IFNULL(SUM(amount), 0) amount,
	      loadsh_id
	  FROM commissions
	  WHERE loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND type = 'SHARE'
	  AND calc = 'USD'
	  AND status <> 'DELETED'
	) flat_shares ON l.id = flat_shares.loadsh_id
	/* PCT Based SHAREs 'shared' from MAIN type */
	LEFT JOIN
	(
	  SELECT
	      IFNULL(SUM(amount), 0) amount,
	      loadsh_id
	  FROM commissions
	  WHERE loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND type = 'SHARE'
	  AND calc = 'PCT'
	  AND status <> 'DELETED'
	) shares ON l.id = shares.loadsh_id
	SET
	  c.office = l.office,
	  c.agent_id = IF( l.office = ( SELECT office FROM agents WHERE id = l.sales_id ), l.sales_id, o.owner_id),
	  c.date = NULLIF(
	      (CASE
	          WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'invoice_date'  THEN l.invoice_date
	          WHEN o.commission_basis = 'delfinal_date' THEN l.delfinal_date
	          ELSE NULL
	      END)
	  , '0000-00-00'),
	  c.amount = o.commission_pct,
	  c.basis = o.commission_basis,
	  c.calc = 'PCT',
	  c.total = CAST(
	      ROUND(LOAD_COMMISSION(
	          (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	          (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	          l.broker_liability_insurance,
	          l.gps_total,
	          l.banyan_fee,
	          /* subtract percentage based SHAREs from MAIN */
	          ( o.commission_pct - IFNULL(shares.amount, 0) )/100
	      ),2)
	      /* subtract FLAT USD shares from MAIN */
	      - IFNULL(flat_shares.amount,0)
	  AS DECIMAL(12,2) ),
	  c.cust_total = l.cust_total,
	  c.carr_total = l.carr_total,
	  c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	  c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND c.status <> 'DELETED'
	AND c.type = 'MAIN'
	/* only update if commission date is in the future or TBD */
	AND (
	      NULLIF(c.date, '0000-00-00') IS NULL
	      OR c.date > DATE(NOW())
	  )`

	_, err = tx.Query(commissionsQuery13)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query169")
	}
	commissionsQuery15 := `INSERT INTO commissions
	(type, status, basis, office, loadsh_id, agent_id, group_id, amount, calc, gl_code, date, total)
	(
	  SELECT * FROM (
	      SELECT
	          pivot.type,
	          'OPEN' status,
	          o.commission_basis,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              l.office,
	              IFNULL(share_agent.office, a.office)
	          ) office,
	          l.id loadsh_id,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              o.owner_id,
	              IFNULL(share.agent_id, ca.sales_id)
	          ) agent_id,
	          cg.id group_id,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              cg.office_amount,
	              cg.sales_amount * IFNULL(share.amount / all_shares.total, 1.00)
	          ) amount,
	          'PCT' calc,
	          cg.gl_code,
	          DATE(main.date) date,
	          LOAD_COMMISSION(
	              (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	              (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	              l.broker_liability_insurance,
	              0, /* charge GPS only to MAIN */
	              l.banyan_fee,
	              IF(pivot.type = 'OFFICE_INCENTIVE',
	                  cg.office_amount,
	                  cg.sales_amount * IFNULL(share.amount/all_shares.total, 1.00)
	              )/100
	          ) AS total
	      FROM commission_incentives ci
	      JOIN commission_groups cg ON ci.id = cg.incentive_id
	      JOIN loadsh l ON (
	      (cg.load_field_name = 'office' AND l.office = cg.load_field_value)
	          OR
	      (cg.load_field_name = 'custa_id' AND l.custa_id = cg.load_field_value)
	      )
	      JOIN offices o ON l.office = o.code
	      /* Fees to be removed from profit */
	      LEFT JOIN (
	          SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	          FROM lsfees tf
	          LEFT JOIN lsstops ts ON tf.lsstop_id = ts.id
	          LEFT JOIN lscarr tc ON ts.lscarr_id = tc.id
	          WHERE tf.excluded_from_commissions = 1
	          AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	          GROUP BY tc.loadsh_id
	      ) excluded_fees ON l.id = excluded_fees.loadsh_id
	      JOIN cust_agent ca ON l.custa_id = ca.id
	      -- join to MAIN in order to get commissions date of MAIN commission
	      JOIN (
	          SELECT date, loadsh_id
	          FROM commissions
	          WHERE type = 'MAIN'
	          AND NULLIF(date, '0000-00-00') IS NOT NULL
	          AND NULLIF(post_id, 0) IS NULL
	          AND status <> 'DELETED'
	          AND loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	      ) main ON l.id = main.loadsh_id
	      JOIN agents a ON ca.sales_id = a.id
	      JOIN (
	          SELECT 'OFFICE_INCENTIVE' AS type
				UNION ALL
				SELECT 'SALES_INCENTIVE'
			) pivot
	      -- get sum of all PCT Load Default Share amounts, so that we can figure out proportionate split of SALES_INCENTIVE percentage
	      LEFT JOIN (
	          SELECT loadsh.id as load_id, SUM(all_shares.amount) as total
	          FROM loadsh
	          INNER JOIN commission_groups as all_shares
	              ON IFNULL(all_shares.incentive_id, 0) = 0
	              AND all_shares.load_field_name = 'custa_id'
	              AND all_shares.load_field_value = loadsh.custa_id
	              AND all_shares.active = 1
	              AND all_shares.internal = 0
	              AND all_shares.calc = 'PCT'
	          WHERE loadsh.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	          ) all_shares ON all_shares.load_id = l.id
	      -- for sales incentives only, join to any PCT Load Default Share for this customer;
	      --  if none exist, larger SELECT will still return one SALES_INCENTIVE record because this is an outer join;
	      --  but if any _do_ exist, larger select will return one SALES_INCENTIVE record _per_ Load Default Share
	      LEFT JOIN commission_groups as share
	          ON IFNULL(share.incentive_id, 0) = 0
	          AND pivot.type = 'SALES_INCENTIVE'
	          AND share.load_field_name = cg.load_field_name
	          AND share.load_field_value = cg.load_field_value
	          AND share.active = 1
	          AND share.internal = 0
	          AND share.calc = 'PCT'
	      -- join to agent record of share agents (so we can get office)
	      LEFT JOIN agents as share_agent
	          ON share_agent.id = share.agent_id
	      -- join to any existing SALES_INCENTIVE for this load
	      LEFT JOIN commissions c ON c.group_id = cg.id
	          AND c.type = 'SALES_INCENTIVE'
	          AND pivot.type = 'SALES_INCENTIVE'
	          AND c.loadsh_id = l.id
	      -- join to any existing OFFICE_INCENTIVE for this load
	      LEFT JOIN commissions c2 ON c2.group_id = cg.id
	          AND c2.type = 'OFFICE_INCENTIVE'
	          AND pivot.type = 'OFFICE_INCENTIVE'
	          AND c2.loadsh_id = l.id
	      WHERE ci.active > 0
	      -- MAIN commission on load has a date value set
	      AND NULLIF(main.date, '0000-00-00') IS NOT NULL
	      -- commission date is after effective date (if set)
	      AND ( NULLIF(cg.effective_date, '0000-00-00') IS NULL
	          OR cg.effective_date <= (CASE
	              WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	              ELSE main.date
	          END)
	      )
	      -- commission date is before expiry date (if set)
	      AND ( NULLIF(cg.expiry_date, '0000-00-00') IS NULL
	          OR cg.expiry_date >= (CASE
	              WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	              ELSE main.date
	          END)
	      )
	      -- SALES_INCENTIVE and/or OFFICE_INCENTIVE are applicable but do not already exist for this load
	      AND ( (pivot.type = 'SALES_INCENTIVE' AND ci.sales_share > 0 AND cg.sales_amount <> 0 AND c.id IS NULL)
	          OR (pivot.type = 'OFFICE_INCENTIVE' AND ci.office_share > 0 AND cg.office_amount <> 0 AND c2.id IS NULL)
	      )
	      AND cg.active > 0
	      AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  ) temp_new_incentives
	  WHERE temp_new_incentives.amount <> 0
	)`

	_, err = tx.Query(commissionsQuery15)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query170")
	}
	commissions81 := `INSERT INTO commissions
	(
	  type,
	  status,
	  basis,
	  office,
	  loadsh_id,
	  agent_id,
	  group_id,
	  billing_adjustment_id,
	  parent_commission_id,
	  note,
	  date,
	  total,
	  cust_total,
	  carr_total,
	  profit,
	  amount,
	  calc,
	  gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	      CASE
	          WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	          WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	          WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	          ELSE 'LOAD_EDIT'
	      END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	      NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                  ROUND(LOAD_COMMISSION(
	                      (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                      (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                      l.broker_liability_insurance,
	                      IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                      l.banyan_fee,
	                      c.amount/100
	                  ),2)
	          END
	          /* ... and subtract this original commission total... */
	          - c.total
	          /* ... and subtract any other applicable commissions totals: */
	          - IFNULL((
	              CASE
	                  /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                  WHEN c.type = 'MAIN'
	                  THEN (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                          AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                          AND status <> 'DELETED'
	                  )
	                  /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                          with this commission's ID as parent commission ID. */
	                  ELSE (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          parent_commission_id = c.id
	                          AND status <> 'DELETED'
	                  )
	              END
	          ), 0)
	      AS DECIMAL(12,2) ) AS total
	      /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.cust_total
	             - c.cust_total
	             - IFNULL((
	              SELECT SUM(cust_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as cust_total
	      /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.carr_total
	             - c.carr_total
	             - IFNULL((
	              SELECT SUM(carr_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as carr_total
	      /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	             - c.profit
	             - IFNULL((
	              SELECT SUM(profit)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	  /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	  WHERE c.status <> 'DELETED'
	  /* only insert if previous commission date is NOT in the future or TBD */
	  AND NOT (
	          NULLIF(c.date, '0000-00-00') IS NULL
	          OR c.date > DATE(NOW())
	      )
	  /* Original commission records of specified type only */
	  AND c.type = 'INTERNAL_SHARE'
	  AND c.parent_commission_id IS NULL
	  AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissions81)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query171")
	}

	commissions18 := `INSERT INTO commissions
	(
	  type,
	  status,
	  basis,
	  office,
	  loadsh_id,
	  agent_id,
	  group_id,
	  billing_adjustment_id,
	  parent_commission_id,
	  note,
	  date,
	  total,
	  cust_total,
	  carr_total,
	  profit,
	  amount,
	  calc,
	  gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	      CASE
	          WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	          WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	          WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	          ELSE 'LOAD_EDIT'
	      END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	      NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                  ROUND(LOAD_COMMISSION(
	                      (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                      (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                      l.broker_liability_insurance,
	                      IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                      l.banyan_fee,
	                      c.amount/100
	                  ),2)
	          END
	          /* ... and subtract this original commission total... */
	          - c.total
	          /* ... and subtract any other applicable commissions totals: */
	          - IFNULL((
	              CASE
	                  /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                  WHEN c.type = 'MAIN'
	                  THEN (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                          AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                          AND status <> 'DELETED'
	                  )
	                  /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                          with this commission's ID as parent commission ID. */
	                  ELSE (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          parent_commission_id = c.id
	                          AND status <> 'DELETED'
	                  )
	              END
	          ), 0)
	      AS DECIMAL(12,2) ) AS total
	      /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.cust_total
	             - c.cust_total
	             - IFNULL((
	              SELECT SUM(cust_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as cust_total
	      /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.carr_total
	             - c.carr_total
	             - IFNULL((
	              SELECT SUM(carr_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as carr_total
	      /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	             - c.profit
	             - IFNULL((
	              SELECT SUM(profit)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	  /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	  WHERE c.status <> 'DELETED'
	  /* only insert if previous commission date is NOT in the future or TBD */
	  AND NOT (
	          NULLIF(c.date, '0000-00-00') IS NULL
	          OR c.date > DATE(NOW())
	      )
	  /* Original commission records of specified type only */
	  AND c.type = 'SHARE'
	  AND c.parent_commission_id IS NULL
	  AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissions18)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query172")
	}

	log.Println("check 8")

	commissions27 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'OFFICE_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissions27)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query173")
	}

	commissions72 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SALES_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissions72)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query174")
	}
	commissions54 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'MAIN'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissions54)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query175")
	}

	//	Error Code: 1175. You are using safe update mode and you tried to update a table without a WHERE that uses a KEY column To disable safe mode, toggle the option in Preferences -> SQL Editor and reconnect.	0.234 sec
	//not working
	commissions45 := `UPDATE commission_incentives ci
	JOIN commission_groups cg ON ci.id = cg.incentive_id
	JOIN cust_agent ca ON cg.load_field_value = ca.id AND cg.load_field_name = 'custa_id'
	JOIN cust_master cm ON ca.custm_id = cm.id
	JOIN (
	  SELECT
	      cust_master.id as custm_id,
	      IF(commission_groups.basis = 'ship_date',
	          loadsh.ship_date,
	          CURDATE()
	      ) as eff_date
	  FROM commissions
	  JOIN commission_groups ON commissions.group_id = commission_groups.id
	  JOIN cust_agent ON commission_groups.load_field_value = cust_agent.id AND commission_groups.load_field_name = 'custa_id'
	  JOIN cust_master ON cust_agent.custm_id = cust_master.id
	  JOIN loadsh ON commissions.loadsh_id = loadsh.id
	  WHERE commissions.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND NULLIF(commission_groups.incentive_id, 0) IS NOT NULL
	  AND commission_groups.auto_expiry_months IS NOT NULL
	  AND NULLIF(commission_groups.effective_date, '0000-00-00') IS NULL
	  AND NULLIF(commission_groups.expiry_date, '0000-00-00') IS NULL
	  GROUP BY cust_master.id
	) TT ON cm.id = TT.custm_id
	SET
		cg.effective_date = TT.eff_date,
		cg.expiry_date = DATE_ADD(TT.eff_date, INTERVAL cg.auto_expiry_months MONTH)
	WHERE ci.active > 0
	AND cg.auto_expiry_months IS NOT NULL
	AND NULLIF(cg.effective_date, '0000-00-00') IS NULL
	AND NULLIF(cg.expiry_date, '0000-00-00') IS NULL`
	fmt.Println(commissions45)
	_, err = tx.Query(commissions45)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error Query177", "500")
	} else {
		fmt.Println("success-Query177")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	Query180 := `INSERT INTO loadsh_post (carrier_email, comments, created_at, deleted_at, emailed_at, id, is_matched, last_error, last_error_payload, last_job_params, last_operation, last_payload, load_board_id, load_id, posted_at, posted_id, posted_load, updated_at)
	VALUES (NULL, '', '2022-08-02 07:28:03', NULL, NULL, '0', '0', NULL, NULL, NULL, NULL, NULL, '1','` + strconv.Itoa(loadshIdVar.LoadshId) + `', NULL, NULL, NULL, NULL)`
	log.Println("check 9 ")

	_, err = tx.Query(Query180)

	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query180")
	}
	tx.Commit()
	DBClient.Close()
	//	DBClient.Close()
	//	if err != nil {
	//		logger.Error("Error while inserting values in the database:" + err.Error())
	//		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	//	}

	BookResponse := dto.BookLTLResponseDTO{
		LoadNumber:  loadshIdVar.LoadshId,
		QuoteId:     banyanQuoteId,
		QuoteNumber: banyanQuoteNumber,
		AgentEmail:  CustDetailsVar.AgentEmail,
		CustEmail:   CustDetailsVar.CustAgentEmail,
		CustFax:     CustDetailsVar.CustAgentFax,
		CustPhone:   CustDetailsVar.CustAgentPhone,
		Contact:     CustDetailsVar.Contact,
		PriceDetails: dto.PriceDetails{
			Scac:               SCAC,
			Service:            ServiceId,
			CarrierName:        CarrName,
			CarrierNotes:       CarrNotes,
			TransitTime:        TransitTime,
			FlatPrice:          CustTotal,
			FuelSurchargePrice: CustFeulSurCharge,
		},
		TotalPrice: CustTotal,
	}

	return &BookResponse, nil

}

// BGLD_ps_1.3-1.8
// Valid load for banyan rating engine will be created with InsertBookDetailsModeRepo()
func (d RepositoryDb) InsertBookDetailsModeRepo(modeRes *dto.QuoteDetailsDto, bookreq *dto.BookLTLRequestDTO) (*dto.BookLTLResponseDTO, *errs.AppError) {
	var err error
	//var responseVar string
	DBClient, DbError := database.BTMSDBWriterConnection()
	tx, rollError := DBClient.Beginx()
	if rollError != nil {
		fmt.Println("error in dbconnection CreateUserRepo")
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	}
	if DbError != nil {
		fmt.Println("entered error")
		return nil, DbError
	}
	CustomerInfoQuery := `select  cust_agent.salesperson as "salesperson",cust_master.name as "name",cust_agent.id as "custaId", cust_agent.custm_id as "custmId", agents.tel as "agentPhoneNumber", agents.name as "agentName", agents.email as "agentEmail", agents.office as "officeCode", IFNULL(SUM(loadsh_balance.current_amount), 0) AS "daxBalance", IFNULL(cust_bal_temp.fats_balance, 0) AS "btmsBalance",agents.id as "salesId",cust_agent.contact as "contact",cust_agent.phone as "cusPhone",COALESCE(cust_agent.email, '')  as "cusEmail",COALESCE(cust_agent.acct_manager_id, '0')  as 'accountmanagerid',COALESCE(cust_agent.acct_manager, '') as 'accountmanager' ,cust_agent.fax as "cusFax" from agents inner join cust_agent on cust_agent.sales_id=agents.id inner join cust_master on cust_master.id=cust_agent.custm_id inner join cust_bal_temp on cust_master.id=cust_bal_temp.custm_id left JOIN loadsh_balance on cust_bal_temp.custm_id = loadsh_balance.custm_id where cust_agent.id ='` + bookreq.CustaId + `'`

	//var CustDetailsVar CustmerDetails
	//_, err = DBClient.Get(CustomerInfoQuery, &CustDetailsVar)

	var CustDetailsVar CustmerDetails
	fmt.Println(CustDetailsVar)
	CustDetailsResp := DBClient.Get(&CustDetailsVar, CustomerInfoQuery)

	if CustDetailsResp != nil {
		logger.Error(CustDetailsResp.Error())
		fmt.Println("", CustDetailsVar.Name)
	}
	fmt.Printf("asdfgCustDetailsVar %+v", CustDetailsVar)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-CustDetailsResp")
	}
	floor_plan_success := 1
	floor_plan_inches := 12
	floor_plan_inches_avg := 6
	floor_plan_inches_user := 0
	equip := "v"
	banyan_rating_complete := 1
	paytype := "0"
	load_ltl_quote_requestQuery := `INSERT INTO load_ltl_quote_request
	      SET cust_name='` + CustDetailsVar.Name + `',
	          custm_id='` + bookreq.CustmId + `',
	          custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	          salesperson='',
	          sales_id=` + strconv.Itoa(CustDetailsVar.SalesId) + `,
	          office='` + CustDetailsVar.OfficeCode + `',
	          created=NOW(),
	          ship_zip='` + bookreq.ShipperDetails.ShipZipCode + `',
	          ship_state='` + bookreq.ShipperDetails.ShipState + `',
	          ship_city='` + bookreq.ShipperDetails.ShipCity + `',
	          ship_limited='',
	          cons_zip='` + bookreq.ConsigneeDetails.ConsZipCode + `',
	          cons_state='` + bookreq.ConsigneeDetails.ConsState + `',
	          cons_city='` + bookreq.ConsigneeDetails.ConsCity + `',
	          cons_limited='',
	          floor_plan_success='` + strconv.Itoa(floor_plan_success) + `',
	          floor_plan_inches='` + strconv.Itoa(floor_plan_inches) + `',
	          floor_plan_inches_avg='` + strconv.Itoa(floor_plan_inches_avg) + `',
	          floor_plan_inches_user='` + strconv.Itoa(floor_plan_inches_user) + `',
	          banyan_rating_complete='` + strconv.Itoa(banyan_rating_complete) + `',
	          banyan_load_id='` + strconv.Itoa(bookreq.LoadId) + `',
	          banyan_rerate_reason='',
	          equip='` + equip + `',
	          paytype='` + paytype + `',
	          fee_codes='', notes='a:0:{}'`
	//'GUA,SEP,NOT'

	_, err = tx.Query(load_ltl_quote_requestQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	for i, _ := range bookreq.Commodities.Commodities {
		stackable := 0
		hazmat := 0
		if bookreq.Commodities.Commodities[i].Stackable == true {
			stackable = 1
		}
		if bookreq.Commodities.Commodities[i].Hazmat == true {
			stackable = 1
		}

		load_ltl_quote_request_items_Query := `INSERT INTO load_ltl_quote_request_items
		(load_ltl_quote_request_id, descrip, pallet_length, pallet_width,pallet_height,weight,nmfc,class,qty, load_item_id,stackable,
		  hazmat,density_item,locked,preapproved)
		VALUES( (select id from load_ltl_quote_request where banyan_load_id ='` + strconv.Itoa(bookreq.LoadId) + `' order by id limit 1), '` + bookreq.Commodities.Commodities[i].Desc + `', '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `', '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Weight) + `','` + bookreq.Commodities.Commodities[i].NMFC + `','` + bookreq.Commodities.Commodities[i].Class + `','` + strconv.Itoa(bookreq.Commodities.Commodities[i].Quantity) + `','` + strconv.Itoa(i) + `','` + strconv.Itoa(stackable) + `','` + strconv.Itoa(hazmat) + `', 0, '', '0')`

		fmt.Println("Querylalalala", load_ltl_quote_request_items_Query)
		_, err = tx.Query(load_ltl_quote_request_items_Query)
		if rollError != nil {
			fmt.Println("error in dbconnection CreateUserRepo")
			return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
		}
		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database 33", "500")
		} else {
			fmt.Println("success-Query33")
		}
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	banyanQuoteId := ""
	banyanQuoteNumber := ""
	CarrLtlId := ""
	CarrId := ""
	CarrName := ""
	ltlFlag := ""
	CustRate := ""
	CarrRate := ""
	CustTotal := ""
	CarrTotal := ""
	TransitTime := ""
	TotalWeight := 0
	SCAC := ""
	Distance := ""
	ServiceId := ""
	for i, _ := range bookreq.Commodities.Commodities {

		TotalWeight += bookreq.Commodities.Commodities[i].Weight
	}

	for i, _ := range modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet {

		if modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.Type == "Charge" {
			if strings.Contains(modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.QuoteInformation.QuoteNumber, strconv.Itoa(bookreq.QuoteId)) {
				str := modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.QuoteInformation.QuoteNumber

				tempArray := strings.Split(str, ";")

				banyanQuoteId = tempArray[1]
				banyanQuoteNumber = tempArray[0]

				CustRate = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].SubTotal
				CarrRate = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.SubTotal
				CustTotal = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Total
				CarrTotal = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.Total
				TransitTime = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.ServiceDays
				Distance = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.Distance
				ServiceId = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].ID
				SCAC = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].SCAC

				fmt.Println("banyanQuoteNumber", banyanQuoteId)
				CarrIdLoadshQuery := ` SELECT L.id, L.name, L.notes, L.pallet_length, L.pallet_width, L.pallet_height, L.pallet_weight, L.ltl_flag, L.ban_supports204,
                      CASE WHEN L.carr_deleted=1 THEN 0 ELSE 1 END as carr_activated, COALESCE(B.id, C.id) as carr_id,
                      CASE WHEN C.status IN ('ACTIVE','CAUTION') THEN 1 ELSE 0 END as carr_usable,
                      CASE WHEN P.blacklist=1 THEN 'BLACKLISTED' WHEN P.blacklist=0 THEN 'FAVORITE' ELSE NULL END as carr_pref,
                      CASE WHEN (office_exclusivity IS NULL OR office_exclusivity LIKE CONCAT('%,',A.office,',%')) THEN 0 ELSE 1 END AS office_prohibited
                    FROM sunteck_fats.carriers_LTL L
                    INNER JOIN cust_agent CA ON CA.id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    INNER JOIN agents A ON A.id=CA.sales_id
                    LEFT OUTER JOIN sunteck_fats.carriers C ON C.scac=L.scac
                    LEFT OUTER JOIN sunteck_fats.carriers B ON B.scac=L.broker_scac AND L.broker_scac<>''
                    LEFT OUTER JOIN cust_agent_carr_prefs P ON P.ltl_carr=1 AND P.carr_id=L.id AND P.custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    WHERE L.scac='` + SCAC + `'AND L.broker_scac=''`
				fmt.Println(CarrIdLoadshQuery, "acncnmcgCustDetailsVar")
				//var CustDetailsVar CustmerDetails
				//_, err = DBClient.Get(CustomerInfoQuery, &CustDetailsVar)
				fmt.Println(CarrIdLoadshQuery)
				var CarrIdLoadshVar CarrIdDetails
				CarrIdLoadshResp := DBClient.Get(&CarrIdLoadshVar, CarrIdLoadshQuery)
				//CustomerInfoResp := DBClient.Get(&CustDetailsVar, CustomerInfoQuery)
				//fmt.Println(CustomerInfoResp, "asdassasfgCustDetailsVar")

				if CarrIdLoadshResp != nil {
					//	logger.Error("uhghjfbkdois" + CarrIdResp.Error())
					fmt.Println("", CustDetailsVar.Name)
				}

				CarrName = CarrIdLoadshVar.Name
				CarrId = strconv.Itoa(CarrIdLoadshVar.CarrId)
				ltlFlag = strconv.Itoa(CarrIdLoadshVar.LtlFlag)
				CarrLtlId = strconv.Itoa(CarrIdLoadshVar.Id)
			}
		}

		fmt.Println("banyanQuoteNumber", banyanQuoteId)
		CarrIdQuery := ` SELECT L.id, L.name, L.notes, L.pallet_length, L.pallet_width, L.pallet_height, L.pallet_weight, L.ltl_flag, L.ban_supports204,
                      CASE WHEN L.carr_deleted=1 THEN 0 ELSE 1 END as carr_activated, COALESCE(B.id, C.id) as carr_id,
                      CASE WHEN C.status IN ('ACTIVE','CAUTION') THEN 1 ELSE 0 END as carr_usable,
                      CASE WHEN P.blacklist=1 THEN 'BLACKLISTED' WHEN P.blacklist=0 THEN 'FAVORITE' ELSE NULL END as carr_pref,
                      CASE WHEN (office_exclusivity IS NULL OR office_exclusivity LIKE CONCAT('%,',A.office,',%')) THEN 0 ELSE 1 END AS office_prohibited
                    FROM sunteck_fats.carriers_LTL L
                    INNER JOIN cust_agent CA ON CA.id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    INNER JOIN agents A ON A.id=CA.sales_id
                    LEFT OUTER JOIN sunteck_fats.carriers C ON C.scac=L.scac
                    LEFT OUTER JOIN sunteck_fats.carriers B ON B.scac=L.broker_scac AND L.broker_scac<>''
                    LEFT OUTER JOIN cust_agent_carr_prefs P ON P.ltl_carr=1 AND P.carr_id=L.id AND P.custa_id='` + strconv.Itoa(CustDetailsVar.CustaId) + `'
                    WHERE L.scac='` + SCAC + `'AND L.broker_scac=''`

		//var CustDetailsVar CustmerDetails
		//_, err = DBClient.Get(CustomerInfoQuery, &CustDetailsVar)

		var CarrIdVar CarrIdDetails
		CarrIdResp := DBClient.Get(&CarrIdVar, CarrIdQuery)
		//CustomerInfoResp := DBClient.Get(&CustDetailsVar, CustomerInfoQuery)
		//fmt.Println(CustomerInfoResp, "asdassasfgCustDetailsVar")

		if CarrIdResp != nil {
			//	logger.Error("uhghjfbkdois" + CarrIdResp.Error())
			fmt.Println("", CustDetailsVar.Name)
		}
		fmt.Println(CarrIdVar, "acncnmcgCustDetailsVar")
		//
		load_ltl_quote_request_results_Query := `INSERT INTO load_ltl_quote_request_results
		(
		  load_ltl_quote_request_id,
		  carriers_LTL_id,
		  quote_unique_id,
		  load_ltl_quote_request_spot_result_id,
		  banyan_fee,
		  raw_soap_result
		  )
		VALUES
		(
		  (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' order by id limit 1),
		  '` + strconv.Itoa(CarrIdVar.Id) + `',
		  '` + banyanQuoteNumber + `',
		  NULL,
		  '0',
		  'NULL'
		)`

		_, err = tx.Query(load_ltl_quote_request_results_Query)
		if rollError != nil {
			fmt.Println("error in dbconnection CreateUserRepo")
			return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
		}
		if err != nil {
			tx.Rollback()
			//logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database 34", "500")
		} else {
			fmt.Println("success-Query34", "iteration", i)
		}

		fee_codes_banyan_ltl_quote_mappingsQuery := `INSERT INTO fee_codes_banyan_ltl_quote_mappings
		   (
		     ltl_quote_result_id,
		     fee_codes_banyan_mapping_id
		     )
		   VALUES
		   (
		     (select id from load_ltl_quote_request_results where quote_unique_id =  '` + banyanQuoteNumber + `'  and load_ltl_quote_request_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' order by id limit 1)order by id limit 1),'0')`
		fmt.Println(fee_codes_banyan_ltl_quote_mappingsQuery)
		_, err = tx.Query(fee_codes_banyan_ltl_quote_mappingsQuery)
		if rollError != nil {
			fmt.Println("error in dbconnection CreateUserRepo")
			return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
		}
		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query37", "iteration", i)
		}
		//if strconv.Itoa(bookReq.QuoteId) == tempArray[1] {
		//
		//}
	}

	LoadShQuery := `INSERT INTO loadsh
	  SET created      = NOW(),
	      created_by   = '',
	      created_id   = '',
	      date         = NOW(),
	      lastmod      =  NOW(),
	      never_edited = 1,
	      lastuser     = '',
	      carr_count   = 1,
	      stop_count   = 2,
	      miles        = '` + Distance + `',
	      fk_auto_load_status_update = '1',
	      ltl_flag     = '` + ltlFlag + `',
	      ltl_indirect = 1,
	      ltl_foreign_quote_id = '` + banyanQuoteNumber + `',
	      ltl_foreign_quote_number = '` + banyanQuoteId + `',
	      ltl_quote_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' limit 1),
	      ltl_quote_result_id = (select id from load_ltl_quote_request_results where quote_unique_id =  '` + banyanQuoteNumber + `'  and load_ltl_quote_request_id = (select  id from load_ltl_quote_request where banyan_load_id = '` + strconv.Itoa(bookreq.LoadId) + `' limit 1)limit 1),
	      ltl_quote_spot_result_id = NULL,
	      ltl_carr_contact = 'NULL',
	      edispatch_load_id = '` + strconv.Itoa(bookreq.LoadId) + `',
	      banyan_fee   = '0',
	      status       = 'QUOTE',
	      dispatcher   = '',
	      disp_id      = '',
	      salesperson  = '` + CustDetailsVar.Salesperson + `',
	      sales_id     = '` + strconv.Itoa(CustDetailsVar.SalesId) + `',
	      acct_manager = '` + CustDetailsVar.AccountManager + `',
	      acct_manager_id ='` + CustDetailsVar.AccountManagerId + `',
	      cust_name    = '` + CustDetailsVar.Name + `',
	      custa_id     ='` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	      custm_id     = '` + bookreq.CustmId + `',
	      office       = '` + CustDetailsVar.OfficeCode + `',
	      ship_city    = '` + bookreq.ShipperDetails.ShipCity + `',
	      ship_state   ='` + bookreq.ShipperDetails.ShipState + `',
	      ship_zip     ='` + bookreq.ShipperDetails.ShipZipCode + `',
	      ship_country = '` + bookreq.ShipperDetails.ShipCountry + `',
	      cons_city    = '` + bookreq.ConsigneeDetails.ConsCity + `',
	      cons_state   = '` + bookreq.ConsigneeDetails.ConsState + `',
	      cons_zip     = '` + bookreq.ConsigneeDetails.ConsZipCode + `',
	      carr_eq      = 'V',
	      cons_country =  '` + bookreq.ConsigneeDetails.ConsCountry + `',
	      carr_name    = '` + CarrName + `',
	      carr_id      = '` + CarrId + `',
	      cust_rate    = '` + CustRate + `',
	      carr_rate    ='` + CarrRate + `',
	      cust_total   = '` + CustTotal + `',
	      carr_total   = '` + CarrTotal + `',
	      linear_feet  = '',
	      ltl_volume_quote = 0,
	      ltl_transit_days =  '` + TransitTime + `',
	      load_qty     = '` + strconv.Itoa(bookreq.Commodities.Commodities[0].Quantity) + `',
	      load_type    = '` + TransitTime + `',
	      load_weight  = '` + strconv.Itoa(TotalWeight) + `',
	      load_method  = 'ELTL',
	      paytype      = '0',
	      share_amt    = '',
	      share_calc   = '',
	      share_agent_id = '',
		  ship_bl='` + bookreq.BlNumber + `',
          cust_po='` + bookreq.PoNumber + `',
          cust_shipid='` + bookreq.ShippingNumber + `'`

	_, err = tx.Query(LoadShQuery)
	fmt.Println(LoadShQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query61")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsQuery := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 post_id,
	 adjustment_type_id,
	 adjustment_payment_id,
	 parent_commission_id,
	 billing_adjustment_id,
	 note,
	 date,
	 total,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT
	 CASE
	     WHEN cg.internal=1 THEN 'INTERNAL_SHARE'
	     ELSE 'SHARE'
	 END                 AS 'type',
	 'OPEN'             AS 'status',
	 o.commission_basis  AS 'basis',
	 a.office            AS 'office',
	 l.id                AS 'loadsh_id',
	 cg.agent_id         AS 'agent_id',
	 cg.id               AS 'group_id',
	 NULL                AS 'post_id',
	 NULL                AS 'adjustment_type_id',
	 NULL                AS 'adjustment_payment_id',
	 NULL                AS 'parent_commission_id',
	 NULL                AS 'billing_adjustment_id',
	 NULL                AS 'note',
	 NULL                AS 'date',
	 NULL                AS 'total',
	 cg.amount           AS 'amount',
	 cg.calc             AS 'calc',
	 '510110'            AS 'gl_code'
	FROM loadsh l
	JOIN offices o
	 ON o.code = l.office
	LEFT JOIN commission_groups cg
	 ON cg.load_field_name = 'custa_id'
	 AND l.custa_id = cg.load_field_value
	LEFT JOIN agents a
	 ON a.id = cg.agent_id
	WHERE l.id = (select id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1)
	AND cg.id IS NOT NULL
	AND  a.id IS NOT NULL
	AND IFNULL(cg.active, 1) > 0
	AND IFNULL(cg.incentive_id, 0) < 1 -- exclude incentives
	)`
	fmt.Println(commissionsQuery)
	_, err = tx.Query(commissionsQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query62")
	}

	commissions_internal := `INSERT INTO commissions_internal
	(
	 type,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 commission_id,
	 parent_internal_id,
	 amount,
	 calc,
	 deleted
	)
	(
	SELECT
	 'ORIGINAL'               AS 'type',
	 a.office            AS 'office',
	 l.id                AS 'loadsh_id',
	 cg.agent_id         AS 'agent_id',
	 cg.id               AS 'group_id',
	 (SELECT MAX(c2.id) FROM commissions c2 WHERE c2.loadsh_id = l.id AND c2.type = 'MAIN') AS 'commission_id',
	 NULL                AS 'parent_internal_id',
	 cg.amount           AS 'amount',
	 cg.calc             AS 'calc',
	 0                   AS 'deleted'
	FROM loadsh l
	JOIN offices o
	 ON o.code = l.office
	LEFT JOIN commission_groups cg
	 ON cg.load_field_name = 'custa_id'
	 AND l.custa_id = cg.load_field_value
	LEFT JOIN agents a
	 ON a.id = cg.agent_id
	WHERE l.id = (select id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1)
	AND cg.id IS NOT NULL
	AND  a.id IS NOT NULL
	AND IFNULL(cg.active, 1) > 0
	AND IFNULL(cg.internal, 0) > 0 -- internal only
	)`

	_, err = tx.Query(commissions_internal)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query63")
	}

	loadshIdQuery := `select id , ltl_quote_id ,carr_id from loadsh where ltl_foreign_quote_id = '` + banyanQuoteNumber + `' order by id desc limit 1`

	fmt.Println(loadshIdQuery, "dhcc")

	var loadshIdVar LoadShIdStruct
	loadshIdResp := DBClient.Get(&loadshIdVar, loadshIdQuery)

	if loadshIdResp != nil {
		//	logger.Error("uhghjfbkdois" + CarrIdResp.Error())
		fmt.Println("", CustDetailsVar.Name)
	}
	fmt.Println(loadshIdVar, "loadshIdVar")

	Query64 := ` UPDATE loadsh
	LEFT OUTER JOIN (
	  SELECT loadsh_id, amount, calc, agent_id
	  FROM commissions
	  WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND status <> 'DELETED'
	  AND type = 'SHARE'
	  ORDER BY id ASC
	  LIMIT 1
	) TT ON TT.loadsh_id = loadsh.id
	SET loadsh.share_amt = IFNULL(TT.amount, 0),
	  loadsh.share_calc = IFNULL(TT.calc, 'USD'),
	  loadsh.share_agent_id = IFNULL(TT.agent_id, 0)
	WHERE loadsh.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(Query64)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query64")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	commissionsQueryMain := `INSERT INTO commissions (type, calc, amount, office, loadsh_id, agent_id, basis)
	SELECT
	 'MAIN' AS type,
	 'PCT'  AS calc,
	 CAST(o.commission_pct AS DECIMAL(10,2) ) AS amount,
	 l.office,
	 l.id AS loadsh_id,
	 l.sales_id AS agent_id,
	 o.commission_basis AS basis
	FROM loadsh AS l
	JOIN offices AS o ON l.office = o.code
	LEFT JOIN commissions c ON l.id = c.loadsh_id AND c.type = 'MAIN'
	/* ONLY 1 active MAIN record per load */
	WHERE l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND IFNULL((SELECT 1 FROM commissions c2 WHERE c2.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' AND c2.type = 'MAIN' AND c2.status <> 'DELETED' LIMIT 1), 0) < 1`

	_, err = tx.Query(commissionsQueryMain)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query65")
	}

	commissionsC := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	JOIN offices o ON l.office = o.code
	JOIN agents a ON c.agent_id = a.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	 SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	 FROM lscarr tc
	 JOIN lsstops ts ON ts.lscarr_id = tc.id
	 JOIN lsfees tf ON tf.lsstop_id = ts.id
	 WHERE tf.excluded_from_commissions = 1
	 AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	SET
	c.office = a.office,
	c.basis = o.commission_basis,
	c.date = NULLIF(
	     (CASE
	         WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'invoice_date'  THEN IF( NULLIF(l.invoice_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.invoice_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'delfinal_date' THEN IF( NULLIF(l.delfinal_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.delfinal_date), DATE(NOW())),  NULL )
	         ELSE NULL
	     END)
	 , '0000-00-00'),
	 c.total = CAST( IF( c.calc = 'PCT',
	     ROUND(LOAD_COMMISSION(
	         (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	         (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	         l.broker_liability_insurance,
	         0.00, /* charge GPS only to MAIN */
	         l.banyan_fee,
	         c.amount/100
	     ),2),
	     c.amount )
	 AS DECIMAL(12,2) ),
	 c.adjustment_payment_id = NULL,
	 c.cust_total = l.cust_total,
	 c.carr_total = l.carr_total,
	 c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	 c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND ( c.type = 'SHARE' OR c.type = 'INTERNAL_SHARE' )
	/* only update if commission date is in the future or TBD */
	AND (
	     NULLIF(c.date, '0000-00-00') IS NULL
	     OR c.date > DATE(NOW())
	 )`

	_, err = tx.Query(commissionsC)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query66")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsUpdate := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	 SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	 FROM lscarr tc
	 JOIN lsstops ts ON ts.lscarr_id = tc.id
	 JOIN lsfees tf ON tf.lsstop_id = ts.id
	 WHERE tf.excluded_from_commissions = 1
	 AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	JOIN offices o ON l.office = o.code
	/* FLAT USD SHAREs only joined to MAIN type */
	LEFT JOIN
	(
	 SELECT
	     IFNULL(SUM(amount), 0) amount,
	     loadsh_id
	 FROM commissions
	 WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND type = 'SHARE'
	 AND calc = 'USD'
	 AND status <> 'DELETED'
	) flat_shares ON l.id = flat_shares.loadsh_id
	/* PCT Based SHAREs 'shared' from MAIN type */
	LEFT JOIN
	(
	 SELECT
	     IFNULL(SUM(amount), 0) amount,
	     loadsh_id
	 FROM commissions
	 WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND type = 'SHARE'
	 AND calc = 'PCT'
	 AND status <> 'DELETED'
	) shares ON l.id = shares.loadsh_id
	SET
	 c.office = l.office,
	 c.agent_id = IF( l.office = ( SELECT office FROM agents WHERE id = l.sales_id ), l.sales_id, o.owner_id),
	 c.date = NULLIF(
	     (CASE
	         WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	         WHEN o.commission_basis = 'invoice_date'  THEN l.invoice_date
	         WHEN o.commission_basis = 'delfinal_date' THEN l.delfinal_date
	         ELSE NULL
	     END)
	 , '0000-00-00'),
	 c.amount = o.commission_pct,
	 c.basis = o.commission_basis,
	 c.calc = 'PCT',
	 c.total = CAST(
	     ROUND(LOAD_COMMISSION(
	         (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	         (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	         l.broker_liability_insurance,
	         l.gps_total,
	         l.banyan_fee,
	         /* subtract percentage based SHAREs from MAIN */
	         ( o.commission_pct - IFNULL(shares.amount, 0) )/100
	     ),2)
	     /* subtract FLAT USD shares from MAIN */
	     - IFNULL(flat_shares.amount,0)
	 AS DECIMAL(12,2) ),
	 c.cust_total = l.cust_total,
	 c.carr_total = l.carr_total,
	 c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	 c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND c.status <> 'DELETED'
	AND c.type = 'MAIN'
	/* only update if commission date is in the future or TBD */
	AND (
	     NULLIF(c.date, '0000-00-00') IS NULL
	     OR c.date > DATE(NOW())
	 )`

	_, err = tx.Query(commissionsUpdate)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query67")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsInsert := `INSERT INTO commissions
	(type, status, basis, office, loadsh_id, agent_id, group_id, amount, calc, gl_code, date, total)
	(
	 SELECT * FROM (
	     SELECT
	         pivot.type,
	         'OPEN' status,
	         o.commission_basis,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             l.office,
	             IFNULL(share_agent.office, a.office)
	         ) office,
	         l.id loadsh_id,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             o.owner_id,
	             IFNULL(share.agent_id, ca.sales_id)
	         ) agent_id,
	         cg.id group_id,
	         IF(pivot.type = 'OFFICE_INCENTIVE',
	             cg.office_amount,
	             cg.sales_amount * IFNULL(share.amount / all_shares.total, 1.00)
	         ) amount,
	         'PCT' calc,
	         cg.gl_code,
	         DATE(main.date) date,
	         LOAD_COMMISSION(
	             (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	             (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	             l.broker_liability_insurance,
	             0, /* charge GPS only to MAIN */
	             l.banyan_fee,
	             IF(pivot.type = 'OFFICE_INCENTIVE',
	                 cg.office_amount,
	                 cg.sales_amount * IFNULL(share.amount/all_shares.total, 1.00)
	             )/100
	         ) AS total
	     FROM commission_incentives ci
	     JOIN commission_groups cg ON ci.id = cg.incentive_id
	     JOIN loadsh l ON (
	     (cg.load_field_name = 'office' AND l.office = cg.load_field_value)
	         OR
	     (cg.load_field_name = 'custa_id' AND l.custa_id = cg.load_field_value)
	     )
	     JOIN offices o ON l.office = o.code
	     /* Fees to be removed from profit */
	     LEFT JOIN (
	         SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	         FROM lsfees tf
	         LEFT JOIN lsstops ts ON tf.lsstop_id = ts.id
	         LEFT JOIN lscarr tc ON ts.lscarr_id = tc.id
	         WHERE tf.excluded_from_commissions = 1
	         AND tc.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	         GROUP BY tc.loadsh_id
	     ) excluded_fees ON l.id = excluded_fees.loadsh_id
	     JOIN cust_agent ca ON l.custa_id = ca.id
	     -- join to MAIN in order to get commissions date of MAIN commission
	     JOIN (
	         SELECT date, loadsh_id
	         FROM commissions
	         WHERE type = 'MAIN'
	         AND NULLIF(date, '0000-00-00') IS NOT NULL
	         AND NULLIF(post_id, 0) IS NULL
	         AND status <> 'DELETED'
	         AND loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	     ) main ON l.id = main.loadsh_id
	     JOIN agents a ON ca.sales_id = a.id
	     JOIN (
	         SELECT 'OFFICE_INCENTIVE' AS type
				UNION ALL
				SELECT 'SALES_INCENTIVE'
			) pivot
	     -- get sum of all PCT Load Default Share amounts, so that we can figure out proportionate split of SALES_INCENTIVE percentage
	     LEFT JOIN (
	         SELECT loadsh.id as load_id, SUM(all_shares.amount) as total
	         FROM loadsh
	         INNER JOIN commission_groups as all_shares
	             ON IFNULL(all_shares.incentive_id, 0) = 0
	             AND all_shares.load_field_name = 'custa_id'
	             AND all_shares.load_field_value = loadsh.custa_id
	             AND all_shares.active = 1
	             AND all_shares.internal = 0
	             AND all_shares.calc = 'PCT'
	         WHERE loadsh.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	         ) all_shares ON all_shares.load_id = l.id
	     -- for sales incentives only, join to any PCT Load Default Share for this customer;
	     --  if none exist, larger SELECT will still return one SALES_INCENTIVE record because this is an outer join;
	     --  but if any _do_ exist, larger select will return one SALES_INCENTIVE record _per_ Load Default Share
	     LEFT JOIN commission_groups as share
	         ON IFNULL(share.incentive_id, 0) = 0
	         AND pivot.type = 'SALES_INCENTIVE'
	         AND share.load_field_name = cg.load_field_name
	         AND share.load_field_value = cg.load_field_value
	         AND share.active = 1
	         AND share.internal = 0
	         AND share.calc = 'PCT'
	     -- join to agent record of share agents (so we can get office)
	     LEFT JOIN agents as share_agent
	         ON share_agent.id = share.agent_id
	     -- join to any existing SALES_INCENTIVE for this load
	     LEFT JOIN commissions c ON c.group_id = cg.id
	         AND c.type = 'SALES_INCENTIVE'
	         AND pivot.type = 'SALES_INCENTIVE'
	         AND c.loadsh_id = l.id
	     -- join to any existing OFFICE_INCENTIVE for this load
	     LEFT JOIN commissions c2 ON c2.group_id = cg.id
	         AND c2.type = 'OFFICE_INCENTIVE'
	         AND pivot.type = 'OFFICE_INCENTIVE'
	         AND c2.loadsh_id = l.id
	     WHERE ci.active > 0
	     -- MAIN commission on load has a date value set
	     AND NULLIF(main.date, '0000-00-00') IS NOT NULL
	     -- commission date is after effective date (if set)
	     AND ( NULLIF(cg.effective_date, '0000-00-00') IS NULL
	         OR cg.effective_date <= (CASE
	             WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	             ELSE main.date
	         END)
	     )
	     -- commission date is before expiry date (if set)
	     AND ( NULLIF(cg.expiry_date, '0000-00-00') IS NULL
	         OR cg.expiry_date >= (CASE
	             WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	             ELSE main.date
	         END)
	     )
	     -- SALES_INCENTIVE and/or OFFICE_INCENTIVE are applicable but do not already exist for this load
	     AND ( (pivot.type = 'SALES_INCENTIVE' AND ci.sales_share > 0 AND cg.sales_amount <> 0 AND c.id IS NULL)
	         OR (pivot.type = 'OFFICE_INCENTIVE' AND ci.office_share > 0 AND cg.office_amount <> 0 AND c2.id IS NULL)
	     )
	     AND cg.active > 0
	     AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 ) temp_new_incentives
	 WHERE temp_new_incentives.amount <> 0
	)`

	_, err = tx.Query(commissionsInsert)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query68")
	}

	commissionsInsert2 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'INTERNAL_SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsInsert2)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query69")
	}
	//
	////// repeat
	commissionsInsert3 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsInsert3)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query70")
	}
	////
	//// repeat
	commissionsInsert4 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SHARE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsInsert4)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-commissionsInsert4")
	}
	////
	////// repeat
	commissionsInsert5 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SALES_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsInsert5)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query72")
	}
	////// repeat
	commissionsInsert6 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'MAIN'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsInsert6)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query73")
	}

	commission_incentives := `update commission_incentives ci
	JOIN commission_groups cg ON ci.id = cg.incentive_id
	JOIN cust_agent ca ON cg.load_field_value = ca.id AND cg.load_field_name = 'custa_id'
	JOIN cust_master cm ON ca.custm_id = cm.id
	JOIN (
	 SELECT
	     cust_master.id as custm_id,
	     IF(commission_groups.basis = 'ship_date',
	         loadsh.ship_date,
	         CURDATE()
	     ) as eff_date
	 FROM commissions
	 JOIN commission_groups ON commissions.group_id = commission_groups.id
	 JOIN cust_agent ON commission_groups.load_field_value = cust_agent.id AND commission_groups.load_field_name = 'custa_id'
	 JOIN cust_master ON cust_agent.custm_id = cust_master.id
	 JOIN loadsh ON commissions.loadsh_id = loadsh.id
	 WHERE commissions.loadsh_id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 AND NULLIF(commission_groups.incentive_id, 0) IS NOT NULL
	 AND commission_groups.auto_expiry_months IS NOT NULL
	 AND NULLIF(commission_groups.effective_date, '0000-00-00') IS NULL
	 AND NULLIF(commission_groups.expiry_date, '0000-00-00') IS NULL
	 GROUP BY cust_master.id
	) TT ON cm.id = TT.custm_id
	SET
		cg.effective_date = TT.eff_date,
		cg.expiry_date = DATE_ADD(TT.eff_date, INTERVAL cg.auto_expiry_months MONTH)
	WHERE ci.active > 0
	AND cg.auto_expiry_months IS NOT NULL
	AND NULLIF(cg.effective_date, '0000-00-00') IS NULL
	AND NULLIF(cg.expiry_date, '0000-00-00') IS NULL`

	_, err = tx.Query(commission_incentives)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query75")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	loadsh_nocapreload := `INSERT IGNORE INTO loadsh_nocapreload VALUES ('` + strconv.Itoa(loadshIdVar.LoadshId) + `', 0)`

	_, err = tx.Query(loadsh_nocapreload)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	edispatch_doc_images := `UPDATE edispatch_doc_images
	               SET loadsh_id     = '` + strconv.Itoa(loadshIdVar.LoadshId) + `',
	                   document_type = 'Packing-Worksheet-L` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	               WHERE ltl_quote_id = '` + strconv.Itoa(loadshIdVar.LtlQuoteId) + `'`
	fmt.Println(edispatch_doc_images)
	_, err = tx.Query(edispatch_doc_images)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query77")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lscarr := `INSERT INTO lscarr
	     SET created       = NOW(),
	         lastmod       = NOW(),
	         loadsh_id     = '` + strconv.Itoa(loadshIdVar.LoadshId) + `',
	         carr_id       = '` + strconv.Itoa(loadshIdVar.CarrId) + `',
	         carr_id_LTL   =  '` + CarrLtlId + `',
	         ltl_special_instruct='Quote #` + banyanQuoteNumber + `',
	         carr_order    = 1,
	         carr_name     =  '` + CarrName + `',
	         cust_rate     = '` + CustRate + `',
	         carr_rate     = '` + CarrRate + `',
	         carr_total    = '` + CarrTotal + `',
	         miles_agent   = '',
	         miles_total   = '` + Distance + `',
	         miles_total_questionable = 0,
	         trailer_length = '53',
	         equip         = 'V',
	         cargo_value_opt_id  = '',
	         weight_total  = '` + strconv.Itoa(TotalWeight) + `'`
	//fmt.Println(Query82, SelectedAccessorials, "ssmsm")
	_, err = tx.Query(lscarr)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query82")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsPick := `INSERT INTO lsstops
	     SET created = NOW(),
	         lastmod      = NOW(),
	         lscarr_id    = (select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1),
	         stop_order   = 1,
	         type         = 'PICK',
	         city         = '` + bookreq.ShipperDetails.ShipCity + `',
	         state        = '` + bookreq.ShipperDetails.ShipState + `',
	         zip          = '` + bookreq.ShipperDetails.ShipZipCode + `',
	         contact      = '` + bookreq.ShipperDetails.ShipContactName + `',
	         email          = '` + bookreq.ShipperDetails.ShipEmail + `',
	         bl          = '` + bookreq.BlNumber + `',
	         time          = '` + bookreq.ShipperDetails.ShipEarliestTime + `',
	         time2          = '` + bookreq.ShipperDetails.ShipLatestTime + `',
	         date          = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	         date2          = '` + bookreq.ShipperDetails.ShipLatestDate + `',       
	         tel          = '` + bookreq.ShipperDetails.ShipPhone + `',	       
	         fax          = '` + bookreq.ShipperDetails.ShipFax + `',	       
	         po          = '` + bookreq.PoNumber + `',	       
	         instruct          = '` + bookreq.ShipperDetails.ShipLoadNotes + `',
	         addr1          = '` + bookreq.ShipperDetails.ShipAddress1 + `',
	         addr2          = '` + bookreq.ShipperDetails.ShipAddress2 + `',
	         name          = '` + bookreq.ShipperDetails.ShipName + `',
	         country      = '` + bookreq.ShipperDetails.ShipCountry + `'`

	_, err = tx.Query(lsstopsPick)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query83")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsDrop := `INSERT INTO lsstops
	     SET created = NOW(),
	         lastmod      = NOW(),
	         lscarr_id    = (select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1),
	         stop_order   = 2,
	         type         = 'DROP',
	         city         = '` + bookreq.ConsigneeDetails.ConsCity + `',
	         state        = '` + bookreq.ConsigneeDetails.ConsState + `',
	         zip          ='` + bookreq.ConsigneeDetails.ConsZipCode + `',
	         contact      = '` + bookreq.ConsigneeDetails.ConsContactName + `',
	         email          = '` + bookreq.ConsigneeDetails.ConsEmail + `',
	         bl          = '` + bookreq.BlNumber + `',
	         time          = '` + bookreq.ConsigneeDetails.ConsEarliestTime + `',
	         time2          = '` + bookreq.ConsigneeDetails.ConsLatestTime + `',
	         date          = '` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	         date2          = '` + bookreq.ConsigneeDetails.ConsLatestDate + `',      
	         tel          = '` + bookreq.ConsigneeDetails.ConsPhone + `',	       
	         fax          = '` + bookreq.ConsigneeDetails.ConsFax + `',	       
	         po          = '` + bookreq.PoNumber + `',	       
	         instruct          = '` + bookreq.ConsigneeDetails.ConsLoadNotes + `',
	         addr1          = '` + bookreq.ConsigneeDetails.ConsAddress1 + `',
	         addr2          = '` + bookreq.ConsigneeDetails.ConsAddress2 + `',
	         name          = '` + bookreq.ConsigneeDetails.ConsName + `',
	         country      = '` + bookreq.ConsigneeDetails.ConsCountry + `'`

	_, err = tx.Query(lsstopsDrop)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query84")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	for i, _ := range bookreq.Commodities.Commodities {
		stackable := 0
		hazmat := 0
		if bookreq.Commodities.Commodities[i].Stackable == true {
			stackable = 1
		}
		if bookreq.Commodities.Commodities[i].Hazmat == true {
			stackable = 1
		}
		// warning
		Query85 := `INSERT INTO lsitems
	       SET created    = NOW(),
	           lastmod    = NOW(),
	           lsstop_id  =(select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1)  order by id limit 1),
	           item_order ='` + strconv.Itoa(i+1) + `',
	           qty        = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Quantity) + `',
	           type       =  '` + bookreq.Commodities.Commodities[i].EquipmentType + `',
	           descrip    = '` + bookreq.Commodities.Commodities[i].Desc + `',
	           density    =  '` + bookreq.Commodities.Commodities[i].Density + `',
	           pallet_length = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `',
	           pallet_width  = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `',
	           pallet_height = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `',
	           weight     ='` + strconv.Itoa(bookreq.Commodities.Commodities[i].Weight) + `',
	           stackable  = '` + strconv.Itoa(stackable) + `',
	           hazmat     ='` + strconv.Itoa(hazmat) + `',
	           class      = '` + bookreq.Commodities.Commodities[i].Class + `',
	           nmfc       = '` + bookreq.Commodities.Commodities[i].NMFC + `'`

		_, err = tx.Query(Query85)
		if rollError != nil {
			fmt.Println("error in dbconnection CreateUserRepo")
			return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
		}
		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query85")
		}
		tx.Commit()
		tx, rollError = DBClient.Beginx()

		load_ltl_quote_approved_items := `INSERT INTO load_ltl_quote_approved_items
	         SET custa_id   = '` + strconv.Itoa(CustDetailsVar.CustaId) + `',
	             created    = NOW(),
	             zip        = '` + bookreq.ShipperDetails.ShipZipCode + `',
	             descrip    = '` + bookreq.Commodities.Commodities[i].Desc + `',
	             nmfc       = '` + bookreq.Commodities.Commodities[i].NMFC + `',
	             class      = '` + bookreq.Commodities.Commodities[i].Class + `',
	             pallet_length = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Length) + `',
	             pallet_width  = '` + strconv.Itoa(bookreq.Commodities.Commodities[i].Width) + `',
	             pallet_height ='` + strconv.Itoa(bookreq.Commodities.Commodities[i].Height) + `',
	             stackable     = '` + strconv.Itoa(stackable) + `',
	             hazmat        ='` + strconv.Itoa(hazmat) + `',
	             density_item  = 0,
	             density       = '` + bookreq.Commodities.Commodities[i].Density + `'`

		_, err = tx.Query(load_ltl_quote_approved_items)
		if rollError != nil {
			fmt.Println("error in dbconnection CreateUserRepo")
			return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
		}
		if err != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		} else {
			fmt.Println("success-Query86")
		}
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	fuelSurchage := ""
	for i, _ := range modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet {

		if modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.Type == "Charge" {
			if strings.Contains(modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].AssociatedCarrierPricesheet.PriceSheet.QuoteInformation.QuoteNumber, strconv.Itoa(bookreq.QuoteId)) {
				for j, _ := range modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge {

					if strings.Contains(modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Type, "ACCESSORIAL_FUEL") {
						fuelSurchage = modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Amount
					}
					// warning
					lsfees := `INSERT INTO lsfees
			                     SET created     = NOW(),
			                         lastmod     = NOW(),
			                         lsstop_id   =  (select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1)  order by id limit 1), fee_order = '` + strconv.Itoa(j+1) + `', qty = '1', type = '` + modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Type + `', descrip = '', cust_rate_type = 'FLAT', cust_rate = 'NULL', cust_charge = '` + modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Amount + `', carr_charge =  '` + modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Amount + `', carr_rate_type = 'FLAT', carr_rate = 'NULL', code = '` + modeRes.Response.MercuryResponseDto.PriceSheets.PriceSheet[i].Charges.Charge[j].Description + `'`
					fmt.Println(lsfees, i, j, "sjsjs")
					_, err = tx.Query(lsfees)
					if rollError != nil {
						fmt.Println("error in dbconnection CreateUserRepo")
						return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
					}
					if err != nil {
						tx.Rollback()
						logger.Error("Error while inserting values in the database:" + err.Error())
						return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
					} else {
						fmt.Println("success-Query89")
					}
				}
			}
		}
	}
	//
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	load_ltl_quote_request_Update := `UPDATE load_ltl_quote_request SET loadsh_id='` + strconv.Itoa(loadshIdVar.LoadshId) + `' WHERE id='` + strconv.Itoa(loadshIdVar.LtlQuoteId) + `'`

	_, err = tx.Query(load_ltl_quote_request_Update)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query94")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	// warning
	loadshUpdate := `UPDATE loadsh
	SET
	lastmod = NOW(),
	lastuser = 'Customer',
	created_by = 'Customer',
	status = 'ACTIVE',
	dispatcher = '` + CustDetailsVar.Salesperson + `',
	disp_id ='` + strconv.Itoa(CustDetailsVar.SalesId) + `',
	office = '` + CustDetailsVar.OfficeCode + `',
	ship_name = '` + bookreq.ShipperDetails.ShipName + `',
	ship_date = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	cons_name = '` + bookreq.ConsigneeDetails.ConsName + `',
	cons_date ='` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	carr_name = '` + CarrName + `',
	cust_total_primary = '` + CustTotal + `',
	load_length = '53',
	miles = '` + Distance + `',
	id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	WHERE id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(loadshUpdate)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query161")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()
	lscarridQuery := `select id from lscarr where loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1`

	var lscarridVar LscarrIdStruct
	lscarridResp := DBClient.Get(&lscarridVar, lscarridQuery)

	if lscarridResp != nil {
		logger.Error("uhghjfbkdois" + lscarridResp.Error())
		fmt.Println("", lscarridVar)
	}
	fmt.Println(lscarridVar, "asdfgCustDetailsVar")
	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-lscarridQuery")
	}

	lsstopsidQuery := `select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1)  order by id limit 1`

	var LsstopsidVar LsstopIdStruct
	LsstopsidResp := DBClient.Get(&LsstopsidVar, lsstopsidQuery)

	if LsstopsidResp != nil {
		logger.Error("uhghjfbkdois" + LsstopsidResp.Error())
		fmt.Println("", LsstopsidVar)
	}
	fmt.Println(LsstopsidResp, "asdfgCustDetailsVar")

	if err != nil {

		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-lsstopsidQuery")
	}

	lsstopsiddropQuery := `select id from lsstops where lscarr_id= ( select id from lscarr where loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' limit 1)  order by id limit 1`

	var lsstopsiddropVar LsstopIdStruct
	lsstopsiddropResp := DBClient.Get(&lsstopsiddropVar, lsstopsiddropQuery)

	if lsstopsiddropResp != nil {
		logger.Error("uhghjfbkdois" + lsstopsiddropResp.Error())
		fmt.Println("", lsstopsiddropVar)
	}
	fmt.Println(lsstopsiddropResp, "asdfgCustDetailsVar")

	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database 32", "500")
	} else {
		fmt.Println("success-Query32")
	}
	// warning
	lscarrUpdate := `UPDATE lscarr
	SET
	lastmod =  NOW(),
	carr_name ='` + CarrName + `',
	miles_total = '` + Distance + `',
	miles_agent = '',
	ltl_special_instruct = 'Quote #` + banyanQuoteNumber + `',
	id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'
	WHERE id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'`

	fmt.Println(lscarrUpdate)
	_, err = tx.Query(lscarrUpdate)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query162")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsUpdatePick := `UPDATE lsstops
	SET
	lastmod =  NOW(),
	name = '` + bookreq.ShipperDetails.ShipName + `',
	addr1 = '` + bookreq.ShipperDetails.ShipAddress1 + `',
	date = '` + bookreq.ShipperDetails.ShipEarliestDate + `',
	blind = '',
	showph = '',
	id = '` + strconv.Itoa(LsstopsidVar.LsstopId) + `'
	WHERE id = '` + strconv.Itoa(LsstopsidVar.LsstopId) + `'`
	fmt.Println(lsstopsUpdatePick)
	_, err = tx.Query(lsstopsUpdatePick)
	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query163")
	}
	tx.Commit()
	tx, rollError = DBClient.Beginx()

	lsstopsUpdateDrop := `UPDATE lsstops
	SET
	lastmod =  NOW(),
	name ='` + bookreq.ConsigneeDetails.ConsName + `',
	addr1 = '` + bookreq.ConsigneeDetails.ConsAddress1 + `',
	date = '` + bookreq.ConsigneeDetails.ConsEarliestDate + `',
	blind = '',
	showph = '',
	id = '` + strconv.Itoa(lscarridVar.LscarrId) + `'
	WHERE id ='` + strconv.Itoa(lscarridVar.LscarrId) + `'`

	_, err = tx.Query(lsstopsUpdateDrop)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query164")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	updateLoadsh := `UPDATE loadsh
	LEFT OUTER JOIN (
	  SELECT loadsh_id, amount, calc, agent_id
	  FROM commissions
	  WHERE loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND status <> 'DELETED'
	  AND type = 'SHARE'
	  ORDER BY id ASC
	  LIMIT 1
	) TT ON TT.loadsh_id = loadsh.id
	SET loadsh.share_amt = IFNULL(TT.amount, 0),
	  loadsh.share_calc = IFNULL(TT.calc, 'USD'),
	  loadsh.share_agent_id = IFNULL(TT.agent_id, 0)
	WHERE loadsh.id ='` + strconv.Itoa(loadshIdVar.LoadshId) + `'`

	_, err = tx.Query(updateLoadsh)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query166")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsInsertRes := `INSERT INTO commissions (type, calc, amount, office, loadsh_id, agent_id, basis)
	SELECT
	  'MAIN' AS type,
	  'PCT'  AS calc,
	  CAST(o.commission_pct AS DECIMAL(10,2) ) AS amount,
	  l.office,
	  l.id AS loadsh_id,
	  l.sales_id AS agent_id,
	  o.commission_basis AS basis
	FROM loadsh AS l
	JOIN offices AS o ON l.office = o.code
	LEFT JOIN commissions c ON l.id = c.loadsh_id AND c.type = 'MAIN'
	/* ONLY 1 active MAIN record per load */
	WHERE l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND IFNULL((SELECT 1 FROM commissions c2 WHERE c2.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `' AND c2.type = 'MAIN' AND c2.status <> 'DELETED' LIMIT 1), 0) < 1`

	_, err = tx.Query(commissionsInsertRes)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query167")
	}
	Query168 := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	JOIN offices o ON l.office = o.code
	JOIN agents a ON c.agent_id = a.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	  SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	  FROM lscarr tc
	  JOIN lsstops ts ON ts.lscarr_id = tc.id
	  JOIN lsfees tf ON tf.lsstop_id = ts.id
	  WHERE tf.excluded_from_commissions = 1
	  AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	SET
	 c.office = a.office,
	 c.basis = o.commission_basis,
	 c.date = NULLIF(
	      (CASE
	          WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'invoice_date'  THEN IF( NULLIF(l.invoice_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.invoice_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'delfinal_date' THEN IF( NULLIF(l.delfinal_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.delfinal_date), DATE(NOW())),  NULL )
	          ELSE NULL
	      END)
	  , '0000-00-00'),
	  c.total = CAST( IF( c.calc = 'PCT',
	      ROUND(LOAD_COMMISSION(
	          (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	          (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	          l.broker_liability_insurance,
	          0.00, /* charge GPS only to MAIN */
	          l.banyan_fee,
	          c.amount/100
	      ),2),
	      c.amount )
	  AS DECIMAL(12,2) ),
	  c.adjustment_payment_id = NULL,
	  c.cust_total = l.cust_total,
	  c.carr_total = l.carr_total,
	  c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	  c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND ( c.type = 'SHARE' OR c.type = 'INTERNAL_SHARE' )
	/* only update if commission date is in the future or TBD */
	AND (
	      NULLIF(c.date, '0000-00-00') IS NULL
	      OR c.date > DATE(NOW())
	  )`

	_, err = tx.Query(Query168)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query168")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsresq := `UPDATE commissions c
	JOIN loadsh l ON c.loadsh_id = l.id
	/* Fees to be removed from profit */
	LEFT JOIN (
	  SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	  FROM lscarr tc
	  JOIN lsstops ts ON ts.lscarr_id = tc.id
	  JOIN lsfees tf ON tf.lsstop_id = ts.id
	  WHERE tf.excluded_from_commissions = 1
	  AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY tc.loadsh_id
	) excluded_fees ON l.id = excluded_fees.loadsh_id
	JOIN offices o ON l.office = o.code
	/* FLAT USD SHAREs only joined to MAIN type */
	LEFT JOIN
	(
	  SELECT
	      IFNULL(SUM(amount), 0) amount,
	      loadsh_id
	  FROM commissions
	  WHERE loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND type = 'SHARE'
	  AND calc = 'USD'
	  AND status <> 'DELETED'
	) flat_shares ON l.id = flat_shares.loadsh_id
	/* PCT Based SHAREs 'shared' from MAIN type */
	LEFT JOIN
	(
	  SELECT
	      IFNULL(SUM(amount), 0) amount,
	      loadsh_id
	  FROM commissions
	  WHERE loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND type = 'SHARE'
	  AND calc = 'PCT'
	  AND status <> 'DELETED'
	) shares ON l.id = shares.loadsh_id
	SET
	  c.office = l.office,
	  c.agent_id = IF( l.office = ( SELECT office FROM agents WHERE id = l.sales_id ), l.sales_id, o.owner_id),
	  c.date = NULLIF(
	      (CASE
	          WHEN o.commission_basis = 'ship_date'     THEN IF( NULLIF(l.ship_date, '0000-00-00') IS NOT NULL,  GREATEST(DATE(l.ship_date), DATE(NOW())),  NULL )
	          WHEN o.commission_basis = 'invoice_date'  THEN l.invoice_date
	          WHEN o.commission_basis = 'delfinal_date' THEN l.delfinal_date
	          ELSE NULL
	      END)
	  , '0000-00-00'),
	  c.amount = o.commission_pct,
	  c.basis = o.commission_basis,
	  c.calc = 'PCT',
	  c.total = CAST(
	      ROUND(LOAD_COMMISSION(
	          (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	          (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	          l.broker_liability_insurance,
	          l.gps_total,
	          l.banyan_fee,
	          /* subtract percentage based SHAREs from MAIN */
	          ( o.commission_pct - IFNULL(shares.amount, 0) )/100
	      ),2)
	      /* subtract FLAT USD shares from MAIN */
	      - IFNULL(flat_shares.amount,0)
	  AS DECIMAL(12,2) ),
	  c.cust_total = l.cust_total,
	  c.carr_total = l.carr_total,
	  c.profit = LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0),
	  c.lastmod = CURRENT_TIMESTAMP()
	WHERE c.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	AND c.status <> 'DELETED'
	AND c.type = 'MAIN'
	/* only update if commission date is in the future or TBD */
	AND (
	      NULLIF(c.date, '0000-00-00') IS NULL
	      OR c.date > DATE(NOW())
	  )`

	_, err = tx.Query(commissionsresq)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query169")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	commissionsVar := `INSERT INTO commissions
	(type, status, basis, office, loadsh_id, agent_id, group_id, amount, calc, gl_code, date, total)
	(
	  SELECT * FROM (
	      SELECT
	          pivot.type,
	          'OPEN' status,
	          o.commission_basis,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              l.office,
	              IFNULL(share_agent.office, a.office)
	          ) office,
	          l.id loadsh_id,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              o.owner_id,
	              IFNULL(share.agent_id, ca.sales_id)
	          ) agent_id,
	          cg.id group_id,
	          IF(pivot.type = 'OFFICE_INCENTIVE',
	              cg.office_amount,
	              cg.sales_amount * IFNULL(share.amount / all_shares.total, 1.00)
	          ) amount,
	          'PCT' calc,
	          cg.gl_code,
	          DATE(main.date) date,
	          LOAD_COMMISSION(
	              (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	              (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	              l.broker_liability_insurance,
	              0, /* charge GPS only to MAIN */
	              l.banyan_fee,
	              IF(pivot.type = 'OFFICE_INCENTIVE',
	                  cg.office_amount,
	                  cg.sales_amount * IFNULL(share.amount/all_shares.total, 1.00)
	              )/100
	          ) AS total
	      FROM commission_incentives ci
	      JOIN commission_groups cg ON ci.id = cg.incentive_id
	      JOIN loadsh l ON (
	      (cg.load_field_name = 'office' AND l.office = cg.load_field_value)
	          OR
	      (cg.load_field_name = 'custa_id' AND l.custa_id = cg.load_field_value)
	      )
	      JOIN offices o ON l.office = o.code
	      /* Fees to be removed from profit */
	      LEFT JOIN (
	          SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
	          FROM lsfees tf
	          LEFT JOIN lsstops ts ON tf.lsstop_id = ts.id
	          LEFT JOIN lscarr tc ON ts.lscarr_id = tc.id
	          WHERE tf.excluded_from_commissions = 1
	          AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	          GROUP BY tc.loadsh_id
	      ) excluded_fees ON l.id = excluded_fees.loadsh_id
	      JOIN cust_agent ca ON l.custa_id = ca.id
	      -- join to MAIN in order to get commissions date of MAIN commission
	      JOIN (
	          SELECT date, loadsh_id
	          FROM commissions
	          WHERE type = 'MAIN'
	          AND NULLIF(date, '0000-00-00') IS NOT NULL
	          AND NULLIF(post_id, 0) IS NULL
	          AND status <> 'DELETED'
	          AND loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	      ) main ON l.id = main.loadsh_id
	      JOIN agents a ON ca.sales_id = a.id
	      JOIN (
	          SELECT 'OFFICE_INCENTIVE' AS type
				UNION ALL
				SELECT 'SALES_INCENTIVE'
			) pivot
	      -- get sum of all PCT Load Default Share amounts, so that we can figure out proportionate split of SALES_INCENTIVE percentage
	      LEFT JOIN (
	          SELECT loadsh.id as load_id, SUM(all_shares.amount) as total
	          FROM loadsh
	          INNER JOIN commission_groups as all_shares
	              ON IFNULL(all_shares.incentive_id, 0) = 0
	              AND all_shares.load_field_name = 'custa_id'
	              AND all_shares.load_field_value = loadsh.custa_id
	              AND all_shares.active = 1
	              AND all_shares.internal = 0
	              AND all_shares.calc = 'PCT'
	          WHERE loadsh.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	          ) all_shares ON all_shares.load_id = l.id
	      -- for sales incentives only, join to any PCT Load Default Share for this customer;
	      --  if none exist, larger SELECT will still return one SALES_INCENTIVE record because this is an outer join;
	      --  but if any _do_ exist, larger select will return one SALES_INCENTIVE record _per_ Load Default Share
	      LEFT JOIN commission_groups as share
	          ON IFNULL(share.incentive_id, 0) = 0
	          AND pivot.type = 'SALES_INCENTIVE'
	          AND share.load_field_name = cg.load_field_name
	          AND share.load_field_value = cg.load_field_value
	          AND share.active = 1
	          AND share.internal = 0
	          AND share.calc = 'PCT'
	      -- join to agent record of share agents (so we can get office)
	      LEFT JOIN agents as share_agent
	          ON share_agent.id = share.agent_id
	      -- join to any existing SALES_INCENTIVE for this load
	      LEFT JOIN commissions c ON c.group_id = cg.id
	          AND c.type = 'SALES_INCENTIVE'
	          AND pivot.type = 'SALES_INCENTIVE'
	          AND c.loadsh_id = l.id
	      -- join to any existing OFFICE_INCENTIVE for this load
	      LEFT JOIN commissions c2 ON c2.group_id = cg.id
	          AND c2.type = 'OFFICE_INCENTIVE'
	          AND pivot.type = 'OFFICE_INCENTIVE'
	          AND c2.loadsh_id = l.id
	      WHERE ci.active > 0
	      -- MAIN commission on load has a date value set
	      AND NULLIF(main.date, '0000-00-00') IS NOT NULL
	      -- commission date is after effective date (if set)
	      AND ( NULLIF(cg.effective_date, '0000-00-00') IS NULL
	          OR cg.effective_date <= (CASE
	              WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	              ELSE main.date
	          END)
	      )
	      -- commission date is before expiry date (if set)
	      AND ( NULLIF(cg.expiry_date, '0000-00-00') IS NULL
	          OR cg.expiry_date >= (CASE
	              WHEN cg.basis = 'ship_date' THEN NULLIF(l.ship_date, '0000-00-00')
	              ELSE main.date
	          END)
	      )
	      -- SALES_INCENTIVE and/or OFFICE_INCENTIVE are applicable but do not already exist for this load
	      AND ( (pivot.type = 'SALES_INCENTIVE' AND ci.sales_share > 0 AND cg.sales_amount <> 0 AND c.id IS NULL)
	          OR (pivot.type = 'OFFICE_INCENTIVE' AND ci.office_share > 0 AND cg.office_amount <> 0 AND c2.id IS NULL)
	      )
	      AND cg.active > 0
	      AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  ) temp_new_incentives
	  WHERE temp_new_incentives.amount <> 0
	)`

	_, err = tx.Query(commissionsVar)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query170")
	}
	commissionsrespq := `INSERT INTO commissions
	(
	  type,
	  status,
	  basis,
	  office,
	  loadsh_id,
	  agent_id,
	  group_id,
	  billing_adjustment_id,
	  parent_commission_id,
	  note,
	  date,
	  total,
	  cust_total,
	  carr_total,
	  profit,
	  amount,
	  calc,
	  gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	      CASE
	          WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	          WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	          WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	          ELSE 'LOAD_EDIT'
	      END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	      NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                  ROUND(LOAD_COMMISSION(
	                      (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                      (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                      l.broker_liability_insurance,
	                      IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                      l.banyan_fee,
	                      c.amount/100
	                  ),2)
	          END
	          /* ... and subtract this original commission total... */
	          - c.total
	          /* ... and subtract any other applicable commissions totals: */
	          - IFNULL((
	              CASE
	                  /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                  WHEN c.type = 'MAIN'
	                  THEN (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                          AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                          AND status <> 'DELETED'
	                  )
	                  /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                          with this commission's ID as parent commission ID. */
	                  ELSE (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          parent_commission_id = c.id
	                          AND status <> 'DELETED'
	                  )
	              END
	          ), 0)
	      AS DECIMAL(12,2) ) AS total
	      /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.cust_total
	             - c.cust_total
	             - IFNULL((
	              SELECT SUM(cust_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as cust_total
	      /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.carr_total
	             - c.carr_total
	             - IFNULL((
	              SELECT SUM(carr_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as carr_total
	      /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	             - c.profit
	             - IFNULL((
	              SELECT SUM(profit)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	  /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	  WHERE c.status <> 'DELETED'
	  /* only insert if previous commission date is NOT in the future or TBD */
	  AND NOT (
	          NULLIF(c.date, '0000-00-00') IS NULL
	          OR c.date > DATE(NOW())
	      )
	  /* Original commission records of specified type only */
	  AND c.type = 'INTERNAL_SHARE'
	  AND c.parent_commission_id IS NULL
	  AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsrespq)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query171")
	}

	commissionsrespQuery := `INSERT INTO commissions
	(
	  type,
	  status,
	  basis,
	  office,
	  loadsh_id,
	  agent_id,
	  group_id,
	  billing_adjustment_id,
	  parent_commission_id,
	  note,
	  date,
	  total,
	  cust_total,
	  carr_total,
	  profit,
	  amount,
	  calc,
	  gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	      CASE
	          WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	          WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	          WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	          ELSE 'LOAD_EDIT'
	      END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	      NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                  ROUND(LOAD_COMMISSION(
	                      (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                      (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                      l.broker_liability_insurance,
	                      IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                      l.banyan_fee,
	                      c.amount/100
	                  ),2)
	          END
	          /* ... and subtract this original commission total... */
	          - c.total
	          /* ... and subtract any other applicable commissions totals: */
	          - IFNULL((
	              CASE
	                  /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                  WHEN c.type = 'MAIN'
	                  THEN (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                          AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                          AND status <> 'DELETED'
	                  )
	                  /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                          with this commission's ID as parent commission ID. */
	                  ELSE (
	                      SELECT SUM(total)
	                      FROM commissions
	                      WHERE
	                          parent_commission_id = c.id
	                          AND status <> 'DELETED'
	                  )
	              END
	          ), 0)
	      AS DECIMAL(12,2) ) AS total
	      /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.cust_total
	             - c.cust_total
	             - IFNULL((
	              SELECT SUM(cust_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as cust_total
	      /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (l.carr_total
	             - c.carr_total
	             - IFNULL((
	              SELECT SUM(carr_total)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as carr_total
	      /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	      , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	          THEN 0.00
	        ELSE
	          (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	             - c.profit
	             - IFNULL((
	              SELECT SUM(profit)
	              FROM commissions
	              WHERE
	                  parent_commission_id = c.id
	                  AND status <> 'DELETED'
	              ), 0.00)
	          )
	        END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	  /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	  WHERE c.status <> 'DELETED'
	  /* only insert if previous commission date is NOT in the future or TBD */
	  AND NOT (
	          NULLIF(c.date, '0000-00-00') IS NULL
	          OR c.date > DATE(NOW())
	      )
	  /* Original commission records of specified type only */
	  AND c.type = 'SHARE'
	  AND c.parent_commission_id IS NULL
	  AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id =  '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsrespQuery)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query172")
	}

	commissionsrespQuery2 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'OFFICE_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsrespQuery2)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query173")
	}

	commissionsrespQuery3 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'SALES_INCENTIVE'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsrespQuery3)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query174")
	}
	commissionsrespQuery4 := `INSERT INTO commissions
	(
	 type,
	 status,
	 basis,
	 office,
	 loadsh_id,
	 agent_id,
	 group_id,
	 billing_adjustment_id,
	 parent_commission_id,
	 note,
	 date,
	 total,
	 cust_total,
	 carr_total,
	 profit,
	 amount,
	 calc,
	 gl_code
	)
	(
	SELECT *
	FROM (
		SELECT
	     CASE
	         WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE') THEN c.type
	         WHEN c.type = 'INTERNAL_SHARE' THEN 'INTERNAL_BILLING_ADJUSTMENT'
	         WHEN '0' > 0 THEN 'BILLING_ADJUSTMENT'
	         ELSE 'LOAD_EDIT'
	     END AS type,
			'OPEN' AS status,
			'adjustment' AS basis,
			c.office AS office,
			l.id AS loadsh_id,
			c.agent_id AS agent_id,
			/* DEFAULT SHARE identifier */
			NULLIF(c.group_id, 0) AS 'group_id',
	     NULLIF('0', 0) AS billing_adjustment_id,
			c.id AS parent_commission_id,
		    CASE WHEN c.deleted= 1 THEN CONCAT('Share #', c.id, ' Deleted After Commission Date')
		        ELSE ''
		        END AS note,
			/* Use today for Commission Date */
			DATE(NOW()) AS date,
			CAST(
			    /* Get full correct commission amount for this share as of now... */
			    CASE
			        WHEN c.deleted = 1
			            /* Already finalized share has now been deleted; need to back out for $0.00 total impact */
			            THEN 0.00
			        ELSE
			            /* Normal still active commission. */
	                 ROUND(LOAD_COMMISSION(
	                     (l.cust_total - IFNULL(excluded_fees.cust_charge,0)),
	                     (l.carr_total - IFNULL(excluded_fees.carr_charge,0)),
	                     l.broker_liability_insurance,
	                     IF(c.type = 'MAIN', l.gps_total, 0.00), /* charge GPS only to MAIN */
	                     l.banyan_fee,
	                     c.amount/100
	                 ),2)
	         END
	         /* ... and subtract this original commission total... */
	         - c.total
	         /* ... and subtract any other applicable commissions totals: */
	         - IFNULL((
	             CASE
	                 /* for MAIN, that means all existing SHARE, LOAD_EDIT, and BILLING_ADJUSTMENT totals for this load. */
	                 WHEN c.type = 'MAIN'
	                 THEN (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	                         AND type IN ('SHARE', 'LOAD_EDIT', 'BILLING_ADJUSTMENT')
	                         AND status <> 'DELETED'
	                 )
	                 /* for SHARE and INTERNAL_SHARE, that means any existing commissions totals of any type
	                         with this commission's ID as parent commission ID. */
	                 ELSE (
	                     SELECT SUM(total)
	                     FROM commissions
	                     WHERE
	                         parent_commission_id = c.id
	                         AND status <> 'DELETED'
	                 )
	             END
	         ), 0)
	     AS DECIMAL(12,2) ) AS total
	     /* CUST TOTAL - only show change in cust total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.cust_total
	            - c.cust_total
	            - IFNULL((
	             SELECT SUM(cust_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as cust_total
	     /* CARR TOTAL - only show change in carr total, which is new total minus original MAIN/SHARE total, minus any other diffs on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (l.carr_total
	            - c.carr_total
	            - IFNULL((
	             SELECT SUM(carr_total)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as carr_total
	     /* PROFIT - only show change in profit, which is new profit minus original MAIN/SHARE profit, minus any other diffs to profit on earlier adjustments */
	     , CASE WHEN c.type IN ('OFFICE_INCENTIVE', 'SALES_INCENTIVE')
	         THEN 0.00
	       ELSE
	         (LOAD_PROFIT(l.cust_total, l.carr_total, l.broker_liability_insurance, l.gps_total, l.banyan_fee, 0)
	            - c.profit
	            - IFNULL((
	             SELECT SUM(profit)
	             FROM commissions
	             WHERE
	                 parent_commission_id = c.id
	                 AND status <> 'DELETED'
	             ), 0.00)
	         )
	       END as profit
			, c.amount
			, c.calc  AS calc,
			c.gl_code AS 'gl_code'
		FROM loadsh l
		JOIN commissions c ON l.id = c.loadsh_id
	 /* Fees excluded from commissions impact */
		LEFT JOIN (
			SELECT tc.loadsh_id, IFNULL(SUM(tf.cust_charge),0) cust_charge, IFNULL(SUM(tf.carr_charge),0) carr_charge
			FROM lsfees tf
			JOIN lsstops ts ON tf.lsstop_id = ts.id
			JOIN lscarr tc ON ts.lscarr_id = tc.id
			WHERE tf.excluded_from_commissions = 1
			AND tc.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
			GROUP BY tc.loadsh_id
		) excluded_fees ON l.id = excluded_fees.loadsh_id
	 WHERE c.status <> 'DELETED'
	 /* only insert if previous commission date is NOT in the future or TBD */
	 AND NOT (
	         NULLIF(c.date, '0000-00-00') IS NULL
	         OR c.date > DATE(NOW())
	     )
	 /* Original commission records of specified type only */
	 AND c.type = 'MAIN'
	 AND c.parent_commission_id IS NULL
	 AND IFNULL(c.date, '0000-00-00') > '2021-10-26'
		/* PECENTAGE BASED COMMISSIONS ONLY  */
		AND c.calc = 'PCT'
		AND l.id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	 GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	_, err = tx.Query(commissionsrespQuery4)

	if err != nil {
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query175")
	}

	Query177 := `UPDATE commission_incentives ci
	JOIN commission_groups cg ON ci.id = cg.incentive_id
	JOIN cust_agent ca ON cg.load_field_value = ca.id AND cg.load_field_name = 'custa_id'
	JOIN cust_master cm ON ca.custm_id = cm.id
	JOIN (
	  SELECT
	      cust_master.id as custm_id,
	      IF(commission_groups.basis = 'ship_date',
	          loadsh.ship_date,
	          CURDATE()
	      ) as eff_date
	  FROM commissions
	  JOIN commission_groups ON commissions.group_id = commission_groups.id
	  JOIN cust_agent ON commission_groups.load_field_value = cust_agent.id AND commission_groups.load_field_name = 'custa_id'
	  JOIN cust_master ON cust_agent.custm_id = cust_master.id
	  JOIN loadsh ON commissions.loadsh_id = loadsh.id
	  WHERE commissions.loadsh_id = '` + strconv.Itoa(loadshIdVar.LoadshId) + `'
	  AND NULLIF(commission_groups.incentive_id, 0) IS NOT NULL
	  AND commission_groups.auto_expiry_months IS NOT NULL
	  AND NULLIF(commission_groups.effective_date, '0000-00-00') IS NULL
	  AND NULLIF(commission_groups.expiry_date, '0000-00-00') IS NULL
	  GROUP BY cust_master.id
	) TT ON cm.id = TT.custm_id
	SET
		cg.effective_date = TT.eff_date,
		cg.expiry_date = DATE_ADD(TT.eff_date, INTERVAL cg.auto_expiry_months MONTH)
	WHERE ci.active > 0
	AND cg.auto_expiry_months IS NOT NULL
	AND NULLIF(cg.effective_date, '0000-00-00') IS NULL
	AND NULLIF(cg.expiry_date, '0000-00-00') IS NULL`
	fmt.Println(Query177)
	_, err = tx.Query(Query177)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error Query177", "500")
	} else {
		fmt.Println("success-Query177")
	}

	tx.Commit()
	tx, rollError = DBClient.Beginx()

	loadsh_post := `INSERT INTO loadsh_post (carrier_email, comments, created_at, deleted_at, emailed_at, id, is_matched, last_error, last_error_payload, last_job_params, last_operation, last_payload, load_board_id, load_id, posted_at, posted_id, posted_load, updated_at)
	VALUES (NULL, '', '2022-08-02 07:28:03', NULL, NULL, '0', '0', NULL, NULL, NULL, NULL, NULL, '1','` + strconv.Itoa(loadshIdVar.LoadshId) + `', NULL, NULL, NULL, NULL)`

	_, err = tx.Query(loadsh_post)

	if err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	} else {
		fmt.Println("success-Query180")
	}
	tx.Commit()
	DBClient.Close()
	//	if err != nil {
	//		logger.Error("Error while inserting values in the database:" + err.Error())
	//		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	//	}

	banyanquoteid, _ := strconv.Atoi(banyanQuoteId)
	BookResponse := dto.BookLTLResponseDTO{
		LoadNumber:  loadshIdVar.LoadshId,
		QuoteId:     banyanquoteid,
		QuoteNumber: banyanQuoteNumber,
		AgentEmail:  CustDetailsVar.AgentEmail,
		CustEmail:   CustDetailsVar.CustAgentEmail,
		CustFax:     CustDetailsVar.CustAgentFax,
		CustPhone:   CustDetailsVar.CustAgentPhone,
		Contact:     CustDetailsVar.Contact,
		PriceDetails: dto.PriceDetails{
			Scac:               SCAC,
			Service:            ServiceId,
			CarrierName:        CarrName,
			TransitTime:        TransitTime,
			FlatPrice:          CustTotal,
			FuelSurchargePrice: fuelSurchage,
		},
		TotalPrice: CustTotal,
	}

	fmt.Println(BookResponse, "wdndwjx")
	return &BookResponse, nil

}

func NewBookRepositoryDb() RepositoryDb {
	return RepositoryDb{}
}
