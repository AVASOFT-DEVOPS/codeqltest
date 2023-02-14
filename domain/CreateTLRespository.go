package domain

import (
	// "database/sql"
	"fmt"
	"golang/database"
	"golang/dto"
	"golang/errs"
	"golang/logger"
	"strconv"
)

func (d RepositoryDb) VerifyGCIDPer(req dto.VerifyGCIDRequest) (*dto.VerifyGCIDresponse, *errs.CUserError) {
	// snippet-start:[dynamodb.go.read_item.session]
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	fmt.Println(req, "Enetred verifyGCID with request\n")
	verifyGCIDResponse := dto.VerifyGCIDresponse{
		Permissioncheck: false,
	}

	//GCID permission check
	// var VerifyUserPerResponse dto.VerifyUserDto
	//ps_1.5.2 db Connection is established
	dbclnt, errcl := database.GetCpdbClient()
	if errcl != nil {
		fmt.Println(errcl, "dbclnt err123")
		logger.Error("Error while selecting values in the database:")
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request 123.")
	}

	//ps_1.5.3 The variable to store the response is declared
	var VerifyUserPer []VerifyUserPerRes
	//ps_1.6.1 Query to determine the User permission is defined below
	verifyUserPermissionQuery := `select c.global_customer_id from cpdb.customer_permission c inner join cpdb.user_permission u on c.global_customer_id =u.global_customer_id where u.permission_id =2 and c.permission_id = 2 and u.active=1 and c.active=1 and u.user_id=` + req.UserId + `and u.mysunteck_login_id= '` + req.CustomerLoginId + `'`
	fmt.Println(verifyUserPermissionQuery, "verifyUserQuery")
	queeryerror := dbclnt.Select(&VerifyUserPer, verifyUserPermissionQuery)
	fmt.Println(queeryerror, "queeryerror")
	if queeryerror != nil {
		logger.Error("Error while scanning customer " + queeryerror.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request")
	}
	if len(VerifyUserPer) != 0 {
		// /ps_1.6.2 boolean value in the response is set based on the response returned from the db
		verifyGCIDResponse.Permissioncheck = true
	}

	return &verifyGCIDResponse, nil
}

// ps_1.20.1
func (d RepositoryDb) CreateTLRepo(TlRequest dto.CreateTLReq) (*dto.TLresponse, *errs.CUserError) {
	dbclnt, errcl := database.BTMSDBWriterConnection()
	if errcl != nil {
		fmt.Println(errcl, "dbclnt err123")
		logger.Error("Error while selecting values in the database:" + errcl.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request 123.")
	}
	fmt.Println("entered Userpermission check ")

	//AgentName.AgentId,office,custm_id name,custa_id,

	CustomerInfoQuery := `select  cust_agent.salesperson as "salesperson",cust_master.name as "name",cust_agent.id as "custaId", cust_agent.custm_id as "custmId", agents.tel as "agentPhoneNumber", agents.name as "agentName", agents.email as "agentEmail", agents.office as "officeCode", IFNULL(SUM(loadsh_balance.current_amount), 0) AS "daxBalance", IFNULL(cust_bal_temp.fats_balance, 0) AS "btmsBalance",agents.id as "salesId",cust_agent.contact as "contact",cust_agent.phone as "cusPhone",COALESCE(cust_agent.email, '')  as "cusEmail",COALESCE(cust_agent.acct_manager_id, '0')  as 'accountmanagerid',COALESCE(cust_agent.acct_manager, '') as 'accountmanager' ,cust_agent.fax as "cusFax" from agents inner join cust_agent on cust_agent.sales_id=agents.id inner join cust_master on cust_master.id=cust_agent.custm_id inner join cust_bal_temp on cust_master.id=cust_bal_temp.custm_id left JOIN loadsh_balance on cust_bal_temp.custm_id = loadsh_balance.custm_id inner join cust_mysunteck_logins on cust_mysunteck_logins.custa_id=cust_agent.id where cust_mysunteck_logins.login ='` + TlRequest.CustomerLoginId + `'`
	var CustomerDetails []CustomerDetails
	fmt.Println(CustomerInfoQuery, "CustomerInfoQuery")
	queeryerror := dbclnt.Select(&CustomerDetails, CustomerInfoQuery)
	fmt.Println(queeryerror, "queeryerror")
	if queeryerror != nil {
		logger.Error("Error while scanning customer " + queeryerror.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request")
	}
	var Customerdetail []dto.CustomerDetailsDto
	for _, custD := range CustomerDetails {
		Customerdetail = append(Customerdetail, custD.ToCustomerDetailsDto())
	}
	totalWeight := 0
	for i, _ := range TlRequest.CommoditiesInfo {
		totalWeight += TlRequest.CommoditiesInfo[i].Weight
	}
	insertloadshquery := `INSERT INTO loadsh 
	(
	 lastmod, 
	 lastuser, lastsync, 
	 commission_data_posted, opentime, 
	 openuser, never_edited, 
	 created, audited, 
	 auditor, removed_from_queue, 
	 removed_from_queue_by, created_by, 
	 created_id, carr_count, 
	 stop_count, ltl_flag, 
	 ltl_indirect, ltl_foreign_quote_id, 
	 ltl_foreign_quote_number, ltl_quote_id, 
	 ltl_quote_result_id, ltl_quote_spot_result_id, 
	 ltl_carr_contact, ltl_fees, 
	 carriers_LTL_billto_id, loadtype, 
	 status, substatus, 
	 date, flag, 
	 call_time, call_ok, 
	 dispatcher, disp_id, 
	 salesperson, sales_id, 
	 acct_manager, acct_manager_id, 
	 billing_associate, office, 
	 cust_name, custm_id, 
	 custa_id, cust_po, 
	 cust_shipid, cust_data_posted_to_dax_at, 
	 cust_ref, x_cust_key, 
	 ship_bl, ship_name, 
	 ship_city, ship_state, 
	 ship_zip, ship_intl_addr, 
	 ship_country, ship_date, 
	 cons_name, cons_city, 
	 cons_state, cons_zip, 
	 cons_intl_addr, cons_country, 
	 cons_date, carr_name, 
	 carr_tel, carr_ref, 
	 carr_id, carr_ontime, 
	 carrsat_id, x_cust_total, 
	 x_carr_total, cust_rate, 
	 carr_rate, cust_total, 
	 cust_total_primary, carr_total, 
	 gps_total, notes, 
	 invoice_notes, ci_notes, 
	 disp_notes, terms, 
	 carr_eq, carr_ltl, 
	 load_method, load_qty, 
	 load_type, load_descrip, 
	 load_weight, load_length, 
	 fats_post, post_time, 
	 natl_post, natl_post_time, 
	 natl_comm, tstop_post, 
	 tstop_post_time, tstop_comm, 
	 miles, miles_customer, 
	 confsign, delfinal_date, 
	 invoice_date, invoice_by, 
	 inv_hold_pod, inv_hold_date, 
	 finance_issues, fin_iss_date, 
	 carr_bill_recd, cv, 
	 share_amt, share_calc, 
	 share_agent_id, loadboard_shares_amt, 
	 loadboard_shares_type, disaster_relief, 
	 claim_status, claim_opened, 
	 claim_closed, qu_rate_pmi, 
	 qu_rate_cwt, broker_liability_insurance, 
	 ptbli_finalized, ptbli_finalize_time, 
	 edi_transaction_number, shipment_delayed, 
	 edi_delay_reason, carr_conf_exists, 
	 finance_message, cust_change, 
	 edi_custom1, edi_custom2, 
	 edi_custom3, edi_carrier_scac, 
	 edi_tender_accepted_at, edi_carrier_load_id, 
	 force_paper_invoice, edi_commodity_code, 
	 capreload_note, kofax_notified, 
	 ecancel_date, ebook_date, 
	 edispatch_date, edispatch_user, 
	 edispatch_load_id, edispatch_ack_date, 
	 edispatch_ack_user, edispatch_ack_acdc, 
	 edispatch_other_refs, linear_feet, 
	 ltl_volume_quote, ltl_transit_days, 
	 external_integration_id, billing_associate_id, 
	 billing_queue_visible, invoice_emailed_ts, 
	 billing_alert, billing_alert_message, 
	 intercomp_ref, fk_auto_load_status_update, 
	 fk_alert_dispatcher, customer_notification_emails, 
	 use_gps_with_customer_notification, banyan_fee, 
	 last_finance_contact_date, finance_issue_waiting_on, 
	 ar_settlement_issue, ar_settlement_notes, 
	 ap_settlement_issue, ap_settlement_notes, 
	 assigned_to, show_tracking_events_on_mysuntecktts, 
	 edi_payment_terms, office_level_code_flag_reason_code_id, 
	 office_level_code_operation_code_id, office_level_code_revenue_code_id, 
	 cust_rate_single, import_source, 
	 import_source_id, manifest_id, 
	 transitioned, doc_audit_status, 
	 doc_audit_datetime, doc_audit_user, 
	 invoice_datetime, write_off_override_date)
	VALUES 
	(
	  current_timestamp, 
	 'Customer', '0000-00-00 00:00:00', 
	 NULL, '0000-00-00 00:00:00', 
	 '', '0', 
	 current_timestamp, NULL, 
	 NULL, NULL, 
	 NULL, 'Customer', 
	 '0', '1', 
	 '2', '0', 
	 NULL, NULL, 
	 NULL, NULL, 
	 NULL, NULL, 
	 NULL, NULL, 
	 NULL, '', 
	 'ACTIVE', NULL, 
	 '0000-00-00', 'NO', 
	 '', 'NO', 
	 '` + Customerdetail[0].Salesperson + `','` + strconv.Itoa(Customerdetail[0].SalesId) + `',
	 '` + Customerdetail[0].Salesperson + `','` + strconv.Itoa(Customerdetail[0].SalesId) + `',
	  NULL,NULL,
	  NULL,'` + Customerdetail[0].OfficeCode + `',
	 '` + Customerdetail[0].Name + `','` + strconv.Itoa(Customerdetail[0].CustmId) + `',
	 '` + strconv.Itoa(Customerdetail[0].CustaId) + `','` + TlRequest.ShipmentInfo.PONumber + `',
     '` + TlRequest.ShipmentInfo.ShippingNumber + `',NULL,
     '','',
     '` + TlRequest.ShipmentInfo.BLNumber + `','` + TlRequest.ShipperInfo.ShipName + `',
     '` + TlRequest.ShipperInfo.ShipCity + `','` + TlRequest.ShipperInfo.ShipState + `',
     '` + TlRequest.ShipperInfo.ShipZipcode + `','',
     '` + TlRequest.ShipperInfo.Country + `','` + TlRequest.ShipperInfo.ShipEarliestDate + `',
     '` + TlRequest.ConsigneeInfo.ConsigName + `','` + TlRequest.ConsigneeInfo.ConsigCity + `',
     '` + TlRequest.ConsigneeInfo.ConsigState + `','` + TlRequest.ConsigneeInfo.ConsigZipcode + `',
     '','` + TlRequest.ConsigneeInfo.ConsigCountry + `',
     '` + TlRequest.ConsigneeInfo.ConsigEarliestDate + `','',
     '','',
     '0','0',
     '0','0',
     '0','0',
     '0','0',
     '0','0',
     '0','',
     '','',
     '','',
     '', '',
     'TL','1',
     '3/4 TL','',
     '` + strconv.Itoa(totalWeight) + `','53',
     'NO','0000-00-00 00:00:00',
     'NO','0000-00-00 00:00:00',
     '','NO',
     '0000-00-00 00:00:00','',
     '0','0',
     '', NULL,
     '0000-00-00','0',
     'NO','0000-00-00 00:00:00',
     NULL,NULL,
     NULL,'0',
     '0','USD',
     '0','0',
     'USD','NO',
     NULL,NULL,
     NULL,'0',
     '0','0',
     'NO','0000-00-00 00:00:00',
     NULL,NULL,
     NULL,NULL,
     NULL,'0',
     '','',
     '',NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     NULL,NULL,
     '1',NULL,
     '','',
     NULL,'0',
     '0',NULL,
     NULL,'0',
     NULL,NULL,
     '0',NULL,
     '0',NULL,
     '','0',
     NULL,NULL,
     NULL, NULL,
     '0', NULL,
     NULL,NULL,
     '0',NULL,
     NULL,NULL,
     NULL,NULL)`

	fmt.Println(insertloadshquery, "insertloadshquery")
	var loadId string
	var temploadID int64
	//  var insertloadhshid []loashshinsert

	tx, rollError :=  dbclnt.Beginx()
	if rollError != nil {
		fmt.Println(errcl, "dbclnt err123")
		logger.Error("Error while selecting values in the database:" + errcl.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request 123.")
	}
	insertedloashId, insertloadhSherr := tx.Exec(insertloadshquery)
	fmt.Println(insertloadhSherr, "insertloadhSherr")
	if insertloadhSherr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertloadhSherr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	} else {
		temploadID, _ = insertedloashId.LastInsertId()
	}

	fmt.Println(temploadID, "temploadID")
	loadId = strconv.FormatInt(temploadID, 10)
	fmt.Println(loadId, "loadId")

	insertLscarrQuery := `INSERT INTO 
	lscarr 
	 ( lastmod, 
	 created, loadsh_id,
	 import_source_id, type, 
	 carrier_mode, move_type,
	 carr_id, carr_id_LTL, 
	 factor_id, carrsat_id, 
	 carr_order, carr_name, 
	 vendor, driver_name, 
	 driver_cel, truck_num, 
	 trailer_num, container_codes, 
	 container_num, container_checksum, 
	 cust_rate, cust_rate_pmi, 
	 cust_rate_cwt, carr_rate, 
	 pay_rate, carr_total, 
	 carr_tel, equip, 
	 ref_num, trailer_length, 
	 weight_total, booked_with, 
	 miles_total, miles_total_questionable, 
	 miles_agent, cargo_value, 
	 cargo_value_opt_id, cargo_value_cust_min, 
	 cargo_value_cust_min_override, ltl_rate_requested, 
	 ltl_special_instruct, stop_missing_paperwork_notification, 
	 carrier_invoice_number, booked_with_email, 
	 driver_2_name, driver_2_cel, 
	 booked_with_phone, carr_bill_recd, 
	 hold_pay, hold_pay_reason, 
	 trailer_used, pro_number, 
	 quickpay_requested, available_truck_id, 
	 gps_tracking_charge, gps_tracking_state, 
	 external_tracking_loadid, external_tracking_type_id, 
	 temp_control, temp_min, 
	 temp_max, temp_scale, 
	 temp_type, temp_service, 
	 seal_num, post_status, 
	 carrier_invoice_amount, intermodal_notify_party, 
	 intermodal_service_type_code, intermodal_railroad_plan_code, 
	 intermodal_fc_type, intermodal_csa_custa_id, 
	 intermodal_spq, rule11, 
	 interchange_to, booking_num, 
	 eq_res_num, eq_res_exp_date, 
	 eq_res_exp_time, pick_empty_at, 
	 return_empty_to, gate_res_num, 
	 gate_res_exp, econfirm_status, 
	 target_rate, payables_waiting_on, 
	 payables_message, payables_last_updated, 
	 waybill_num, waybill_num_lock)
	
	
	VALUES
	
	 ( '0000-00-00 00:00:00',
	 current_timestamp, '` + loadId + `',
	 NULL, 'carrier',
	 'Road', NULL,
	 '0', '0',
	 NULL, '0',
	  '1', '',
	 '', '',
	 '', '',
	 '', NULL, 
	 NULL, NULL,
	 '0', '0',
	 '0', '0',
	 '0', '0',
	 '', 'BC',
	 '', '53',
	 '1', '',
	 '0', NULL,
	  NULL, '0',
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 '', NULL, 
	 '0', NULL,
	 NULL, NULL,
	 NULL, NULL,
	  NULL, 'DISABLED',
	 NULL, NULL, 
	 '0', NULL,
	 NULL, NULL,
	 NULL, '0',
	  NULL, 'IN PROCESS',
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, 'NO',
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, 'UNSENT',
	 NULL, NULL,
	 NULL, NULL,
	 '', '0')`

	var lscarrId string
	var templscarrId int64

	fmt.Println(insertLscarrQuery, "insertLscarrQuery")
	insertlscarrId, insertLscarrErr := tx.Exec(insertLscarrQuery)
	fmt.Println(insertLscarrErr, "insertLscarrErr")
	if insertLscarrErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertLscarrErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	} else {
		templscarrId, _ = insertlscarrId.LastInsertId()
	}

	fmt.Println(templscarrId, "templscarrId")
	lscarrId = strconv.FormatInt(templscarrId, 10)
	fmt.Println(lscarrId, "lscarrId")

	insertLsstopsPickQuery := 
	`INSERT INTO lsstops 
  ( lastmod, created,
	lscarr_id, stop_order,
	type, status,
	name, addr1,
	addr2, city, 
	state, zip, 
	intl_addr, country,
	contact,email, 
	bl, date,
	time, date2,
	time2,datein,
	timein,dateout,
	timeout, ontime,
	tel,fax,
	po,blind,
	showph,showemail,
	hide_on_cust_invoice,showdirect,
	interline,instruct,
    direct, notes,
	prt_cust_instruct, prt_cust_direct, 
	prt_cust_notes, prt_carr_instruct,
	prt_carr_direct,prt_carr_notes, 
	edi_214_codes, edi_location_qualifier,
	edi_location_id, sub_bol,
	verified_at,verified_by,
    geo_data,geo_candidate,
	date_extra_info,date_extra_info_other,
	delay_reason,intermodal_rail_ramp_id)
  VALUES 
  ( '0000-00-00 00:00:00', 'current_tiemstamp',
   '` + lscarrId + `', '1', 
   'PICK', '',
    '` + TlRequest.ShipperInfo.ShipName + `', '` + TlRequest.ShipperInfo.ShipAdd1 + `', 
	'` + TlRequest.ShipperInfo.ShipAdd2 + `', '` + TlRequest.ShipperInfo.ShipCity + `',
	'` + TlRequest.ShipperInfo.ShipState + `', '` + TlRequest.ShipperInfo.ShipZipcode+ `',
	'', '` + TlRequest.ShipperInfo.Country + `',  '` + TlRequest.ShipperInfo.ShipContactName + `',
	'` + TlRequest.ShipperInfo.ShipEmail + `', '',
	'` + TlRequest.ShipperInfo.ShipEarliestDate + `','` + TlRequest.ShipperInfo.ShipEarliestTime + `', 
	'` + TlRequest.ShipperInfo.ShipLatestDate + `', '` + TlRequest.ShipperInfo.ShipLatestTime + `', 
	'', '',
	'', '',
	'', '` + TlRequest.ShipperInfo.ShipPhoneNumber + `', 
	'', '',
	'','',
	'NO','NO', 
	'', '',
	'','',
	'', '',
	'', '',
	'', '',
	'', NULL,
	NULL,NULL,
	NULL,NULL,
	NULL, NULL, 
	NULL, NULL,
	NULL, NULL,
	NULL)`
	// `INSERT INTO lsstops
	// (
	//    lastmod,created,
	//    lscarr_id,stop_order,
	//    type,status,
	//    name, addr1,
	//    addr2,city,
	//    state,zip,
	//    intl_addr,country,
	//    contact, email,
	//    bl,date,
	// 	 time,date2,
	// 	 time2,datein,
	// 	 timein,dateout,
	// 	 timeout,ontime,
	// 	 tel,fax,
	// 	 po,blind,
	// 	 showph, showemail,
	// 	 hide_on_cust_invoice,showdirect,
	// 	 interline,instruct,
	// 	 direct,notes,
	// 	 prt_cust_instruct,prt_cust_direct,
	// 	 prt_cust_notes, prt_carr_instruct,
	// 	 prt_carr_direct,prt_carr_notes,
	// 	 edi_214_codes,edi_location_qualifier,
	// 	 edi_location_id,sub_bol,
	// 	 verified_at,verified_by,
	// 	 geo_data, geo_candidate,
	// 	 date_extra_info,date_extra_info_other,
	//    delay_reason,intermodal_rail_ramp_id)

	// VALUES

	//   (NULL, '0000-00-00 00:00:00',
	//   current_timestamp, '`+lscarrId+`',
	//   '1', 'PICK',
	//   '', '`+TlRequest.ShipperInfo.ShipName+`',
	//   '`+TlRequest.ShipperInfo.ShipAdd1+`', '`+TlRequest.ShipperInfo.ShipAdd2+`',
	//   '`+TlRequest.ShipperInfo.ShipCity+`', '`+TlRequest.ShipperInfo.ShipState+`',
	//   '`+string(TlRequest.ShipperInfo.Zipcode)+`', '',
	//   '`+TlRequest.ShipperInfo.Country+`', '`+TlRequest.ShipperInfo.ShipContactName+`',
	//   '`+TlRequest.ShipperInfo.ShipEmail+`', '',
	//   '`+TlRequest.ShipperInfo.ShipEarliestDate+`', '`+TlRequest.ShipperInfo.ShipEarliestTime+`',
	//   '', '`+TlRequest.ShipperInfo.ShipLatestTime+`',
	//   '', '',
	//   '', '',
	//   '706-654-3677', '',
	//   '', '',
	//   '', 'NO',
	//   'NO', '',
	//    '', '',
	//    '', '',
	//    '', '',
	//    '', '',
	//    '', '',
	//   NULL, NULL,
	//   NULL, NULL,
	//   NULL, NULL,
	//   NULL, NULL,
	//   NULL, NULL,
	//   NULL, NULL)`

	var lsStopsId string
	var templsStopsId int64
	fmt.Println(insertLsstopsPickQuery, "insertLsstopsPickQuery")
	insertlsStopsId, insertLsstopsPickErr := tx.Exec(insertLsstopsPickQuery)
	fmt.Println(insertLsstopsPickErr, "insertLsstopsPickErr")
	if insertLsstopsPickErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertLsstopsPickErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	} else {
		templsStopsId, _ = insertlsStopsId.LastInsertId()
	}

	fmt.Println(templsStopsId, "templsStopsId")
	lsStopsId = strconv.FormatInt(templsStopsId, 10)
	fmt.Println(lsStopsId, "lsStopsId")

	for comindex := 0; comindex < len(TlRequest.CommoditiesInfo); comindex++ {
		var insertLsItemsQuery string
		if !TlRequest.CommoditiesInfo[comindex].Hazmat {
			insertLsItemsQuery = `INSERT INTO lsitems (
		  lastmod,
		 created, lsstop_id,
		 item_order, deliv_lsstop_id,
		 deliv_item_order, qty,
		 type, descrip,
		 weight, length,
		 class, nmfc,
		 pallet_length, pallet_width,
		 pallet_height, stackable,
		 hazmat, density,
		 external_integration_id, value,
		 hazmat_ship_name, hazmat_descrip,
		 hazmat_group_name, hazmat_pkg_group,
		 hazmat_un_na_num, hazmat_class,
		 hazmat_placard, hazmat_flash_temp,
		 hazmat_flash_uom, hazmat_flash_type,
		 hazmat_cert_holder, hazmat_contact_name,
		 hazmat_contact_tel, intermodal_commodity_code_id)
		VALUES 
		(
		  '0000-00-00 00:00:00',
		 current_timestamp, '` + lsStopsId + `',
		 '` + strconv.Itoa(comindex+1) + `', '0', 
		 '0', '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Qty) + `', 
		 '` + TlRequest.CommoditiesInfo[comindex].UOM + `', '` + TlRequest.CommoditiesInfo[comindex].Descrip + `', 
		 '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Weight) + `', '0', 
		 '', NULL,
		 NULL, NULL,
		 NULL, '0',
		 '0', NULL,
		 NULL, '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Value) + `',
		 NULL, NULL,
		 NULL, NULL,
		 NULL, NULL,
		 NULL, NULL,
		 NULL, NULL,
		 NULL, NULL,
		 NULL, NULL)`
		} else {
			insertLsItemsQuery = `INSERT INTO lsitems (
	     lastmod,
	    created, lsstop_id,
	    item_order, deliv_lsstop_id,
	    deliv_item_order, qty,
	    type, descrip,
	    weight, length,
	    class, nmfc,
	    pallet_length, pallet_width,
	    pallet_height, stackable,
	    hazmat, density,
	    external_integration_id, value,
	    hazmat_ship_name, hazmat_descrip,
	    hazmat_group_name, hazmat_pkg_group,
	    hazmat_un_na_num, hazmat_class,
	    hazmat_placard, hazmat_flash_temp,
	    hazmat_flash_uom, hazmat_flash_type,
	    hazmat_cert_holder, hazmat_contact_name,
	    hazmat_contact_tel, intermodal_commodity_code_id)
	   VALUES 
	   (
	     '0000-00-00 00:00:00',
	    current_timestamp, '` + lsStopsId + `',
	    '` + strconv.Itoa(comindex+1) + `', '0', 
	    '0', '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Qty) + `', 
	    '` + TlRequest.CommoditiesInfo[comindex].UOM + `', '` + TlRequest.CommoditiesInfo[comindex].Descrip + `', 
	    '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Weight) + `', '0', 
	    '', NULL,
	    NULL, NULL,
	    NULL, '0',
	    '1', NULL,
	    NULL, '` + strconv.Itoa(TlRequest.CommoditiesInfo[comindex].Value) + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.ShippingName + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.Description + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.GroupName + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.PackagingGroup + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.UNNANumber + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.Class + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.PlacardType + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.FlashTemp + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.UOM + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.FlashType + `',
	     '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.CHName + `', '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.ContactName + `',
	    '` + TlRequest.CommoditiesInfo[comindex].HazmatInfo.PhoneNumber + `', NULL)`
		}
		fmt.Println(insertLsItemsQuery, "insertLsItemsQuery")
		insertLsItemsErr := tx.QueryRow(insertLsItemsQuery).Err()
		fmt.Println(insertLsItemsErr, "insertLsItemsErr")
		if insertLsItemsErr != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + insertLsItemsErr.Error())
			return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
		}
	}

		//Verify shippoints
		 var verifyPickAddressRes []verifyAddressRes
		 verifyaddressQuery :=`select id from ship_points where custa_id='`+strconv.Itoa(Customerdetail[0].CustaId)+`' and name='`+TlRequest.ShipperInfo.ShipName+`'`
		 verifyaddressQueryerr:=tx.Select(&verifyPickAddressRes,verifyaddressQuery)
		 if verifyaddressQueryerr != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + verifyaddressQueryerr.Error())
			return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
		}
        
	    if len(verifyPickAddressRes)==0{
			insertShippointsBfrQuery := `INSERT INTO ship_points (
			  lastmod,
			 created, custa_id,
			 type_pick, type_drop,
			 name, addr1,
			 addr2, city,
			 state, zip,
			 intl_addr, country,
			 contact, email,
			 tel, fax,
			 blind, showph,
			 showemail, showdirect,
			 instruct, direct,
			 edi_location_qualifier, edi_location_id,
			 verified_at, verified_by,
			 geo_data, geo_candidate)
			VALUES 
			 ( '0000-00-00 00:00:00', 
			 current_timestamp, '` + strconv.Itoa(Customerdetail[0].CustaId) + `',
			 'YES', 'YES',
			 '` + TlRequest.ShipperInfo.ShipName + `', '` + TlRequest.ShipperInfo.ShipAdd1 + `',
			 '` + TlRequest.ShipperInfo.ShipAdd2 + `', '` + TlRequest.ShipperInfo.ShipCity + `', 
			 '` + TlRequest.ShipperInfo.ShipState + `','` + TlRequest.ShipperInfo.ShipZipcode + `',
			 '', '` + TlRequest.ShipperInfo.Country + `', 
			 '` + TlRequest.ShipperInfo.ShipContactName + `', '` + TlRequest.ShipperInfo.ShipEmail + `', 
			 '` + TlRequest.ShipperInfo.ShipPhoneNumber + `', '', 
			 'NO','NO', 
			 'NO', 'NO',
			 '', '',
			 '', '', 
			 NULL, '', 
			 '', '')`

	fmt.Println(insertShippointsBfrQuery, "insertShippointsBfrQuery")
	insertShippointsBfrErr := tx.QueryRow(insertShippointsBfrQuery).Err()
	fmt.Println(insertShippointsBfrErr, "insertShippointsBfrErr")
	if insertShippointsBfrErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertShippointsBfrErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}
}

	insertLsstopsDropQuery := `INSERT INTO lsstops 
		  ( lastmod, 
			created, lscarr_id,
			stop_order, type,
			status, name,
			addr1, addr2,
			city, state,
		    zip, intl_addr,
		    country, contact, 
		    email, bl,
		    date, time,
		    date2, time2,
		 	datein, timein,
		 	dateout, timeout,
		 	ontime, tel,
		 	fax, po,
		 	blind, showph,
		 	showemail, hide_on_cust_invoice,
		 	showdirect, interline,
		    instruct,direct,
		    notes,prt_cust_instruct,
		    prt_cust_direct,prt_cust_notes,
		    prt_carr_instruct,prt_carr_direct,
		    prt_carr_notes,edi_214_codes,
		    edi_location_qualifier,edi_location_id,
		    sub_bol,verified_at,
		    verified_by,geo_data,
		    geo_candidate,date_extra_info, 
		    date_extra_info_other,delay_reason,
		    intermodal_rail_ramp_id)
		VALUES 
		
	 ( '0000-00-00 00:00:00',
	 current_timestamp, '` + lscarrId + `',
	 '2', 'DROP',
	 '', '` + TlRequest.ConsigneeInfo.ConsigName + `',
	 '` + TlRequest.ConsigneeInfo.ConsigAdd1 + `', '` + TlRequest.ConsigneeInfo.ConsigAdd2 + `',
	 '` + TlRequest.ConsigneeInfo.ConsigCity + `', '` + TlRequest.ConsigneeInfo.ConsigState + `',
	 '` + TlRequest.ConsigneeInfo.ConsigZipcode + `', '',
	 '` + TlRequest.ConsigneeInfo.ConsigCountry + `', '',
	 '', '',
	 '` + TlRequest.ConsigneeInfo.ConsigEarliestDate + `', '` + TlRequest.ConsigneeInfo.ConsigEarliestTime + `',
	 '` + TlRequest.ConsigneeInfo.ConsigLatestDate + `', '` + TlRequest.ConsigneeInfo.ConsigLatestTime + `',
	 '', '',
	 '', '',
	 '', '` + TlRequest.ConsigneeInfo.ConsigPhoneNumber + `',
	 '` + TlRequest.ConsigneeInfo.ConsigFax + `', '',
	 '', '',
	 'NO', 'NO',
	 '', '',
	 '', '',
	 '', '',
	 '', '',
	 '', '',
	 '', NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL,
	 NULL, NULL, 
	 NULL)`

	fmt.Println(insertLsstopsDropQuery, "insertLsstopsDropQuery")

	insertLsstopsDropErr := tx.QueryRow(insertLsstopsDropQuery).Err()
	fmt.Println(insertLsstopsDropErr, "insertLsstopsDropErr")
	if insertLsstopsDropErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertLsstopsDropErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}


	var verifyDropAddressRes []verifyAddressRes
	verifyDropaddressQuery :=`select id from ship_points where custa_id='`+strconv.Itoa(Customerdetail[0].CustaId)+`' and name='`+TlRequest.ConsigneeInfo.ConsigName+`'`
	verifyDropaddressQueryerr:=dbclnt.Select(&verifyDropAddressRes,verifyDropaddressQuery)
	if verifyDropaddressQueryerr != nil {
	   // tx.Rollback()
	   logger.Error("Error while inserting values in the database:" + verifyaddressQueryerr.Error())
	   return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
   }
   if len(verifyDropAddressRes)==0{
	insertShippointsaftQuery:=`INSERT INTO ship_points (
		 lastmod,
		created, custa_id,
		type_pick, type_drop,
		name, addr1,
		addr2, city,
		state, zip,
		intl_addr, country,
		contact, email,
	    tel, fax,
		blind, showph,
		showemail, showdirect,
		instruct, direct,
		edi_location_qualifier, edi_location_id,
		verified_at, verified_by,
		geo_data, geo_candidate)

	 VALUES
	   ( '0000-00-00 00:00:00',
	    current_timestamp, '` + strconv.Itoa(Customerdetail[0].CustaId) + `',
	    'YES', 'YES',
	    '` + TlRequest.ConsigneeInfo.ConsigName + `', '` + TlRequest.ConsigneeInfo.ConsigAdd1 + `',
		'', '` + TlRequest.ConsigneeInfo.ConsigCity + `',
		'` + TlRequest.ConsigneeInfo.ConsigState + `', '` + TlRequest.ConsigneeInfo.ConsigZipcode + `',
		'', '` + TlRequest.ConsigneeInfo.ConsigCountry + `',
		'', '',
		'', '',
		'NO', 'NO',
		'NO', 'NO',
		'', '', 
		'', '',
		NULL, '',
		'', '')`
    
		fmt.Println(insertShippointsaftQuery, "insertShippointsaftQuery")
		insertShippointsaftErr:=tx.QueryRow(insertShippointsaftQuery).Err()
		fmt.Println(insertShippointsaftErr, "insertShippointsaftErr")
		if insertShippointsaftErr != nil {
			tx.Rollback()
			logger.Error("Error while inserting values in the database:" + insertShippointsaftErr.Error())
			return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
		}}

	//COMMISSIONS

	insertShareCommissions1Query := `INSERT INTO commissions
     (
	  type,
	  status,
	  basis,office,
	  loadsh_id,agent_id,
	  group_id,post_id,
	  adjustment_type_id,adjustment_payment_id,
	  parent_commission_id,billing_adjustment_id,
	  note,date,
	  total,amount,
	  calc,gl_code)
     (
     SELECT
	     CASE
	  	   WHEN cg.internal=1 THEN 'INTERNAL_SHARE'
	  	   ELSE 'SHARE'
	     END                 AS 'type','OPEN'             AS 'status',
	     o.commission_basis  AS 'basis',a.office            AS 'office',
	     l.id                AS 'loadsh_id',cg.agent_id         AS 'agent_id',
	     cg.id               AS 'group_id',NULL                AS 'post_id',
	     NULL                AS 'adjustment_type_id',NULL        AS 'adjustment_payment_id',
	     NULL                AS 'parent_commission_id',NULL                AS 'billing_adjustment_id',
	     NULL                AS 'note',NULL                AS 'date',
	     NULL                AS 'total',cg.amount           AS 'amount',
	     cg.calc             AS 'calc','510110'            AS 'gl_code'
     FROM loadsh l
     JOIN offices o
	     ON o.code = l.office
     LEFT JOIN commission_groups cg
	     ON cg.load_field_name = 'custa_id'
	     AND l.custa_id = cg.load_field_value
     LEFT JOIN agents a
	     ON a.id = cg.agent_id
     WHERE l.id = '` + loadId + `'
     AND cg.id IS NOT NULL
     AND  a.id IS NOT NULL
     AND IFNULL(cg.active, 1) > 0
     AND IFNULL(cg.incentive_id, 0) < 1 -- exclude incentives
     )`
	fmt.Println(insertShareCommissions1Query, "insertShareCommissions1Query")
	insertShareCommissions1Err := tx.QueryRow(insertShareCommissions1Query).Err()
	fmt.Println(insertShareCommissions1Err, "insertShareCommissions1Err")
	if insertShareCommissions1Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertShareCommissions1Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissions2Query := `INSERT INTO commissions
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
     WHERE l.id = '` + loadId + `'
     AND cg.id IS NOT NULL
     AND  a.id IS NOT NULL
     AND IFNULL(cg.active, 1) > 0
     AND IFNULL(cg.incentive_id, 0) < 1 -- exclude incentives
     )`

	fmt.Println(insertCommissions2Query, "insertCommissions2Query")
	insertCommissions2Err := tx.QueryRow(insertCommissions2Query).Err()
	fmt.Println(insertCommissions2Err, "insertCommissions2Err")
	if insertCommissions2Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions2Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissionInternalQuery := `INSERT INTO commissions_internal
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
     WHERE l.id = '` + loadId + `'
     AND cg.id IS NOT NULL
     AND  a.id IS NOT NULL
     AND IFNULL(cg.active, 1) > 0
     AND IFNULL(cg.internal, 0) > 0 -- internal only
     )`

	fmt.Println(insertCommissionInternalQuery, "insertCommissionInternalQuery")

	insertCommissionInternalErr := tx.QueryRow(insertCommissionInternalQuery).Err()
	fmt.Println(insertCommissionInternalErr, "insertCommissionInternalErr")
	if insertCommissionInternalErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissionInternalErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	updateLoadhshQuery := `UPDATE loadsh
     LEFT OUTER JOIN (
	     SELECT loadsh_id, amount, calc, agent_id
	     FROM commissions
	     WHERE loadsh_id = '` + loadId + `'
	     AND status <> 'DELETED'
	     AND type = 'SHARE'
	     ORDER BY id ASC
	     LIMIT 1
     ) TT ON TT.loadsh_id = loadsh.id
     SET loadsh.share_amt = IFNULL(TT.amount, 0),
	     loadsh.share_calc = IFNULL(TT.calc, 'USD'),
	     loadsh.share_agent_id = IFNULL(TT.agent_id, 0)
     WHERE loadsh.id = '` + loadId + `'`

	fmt.Println(updateLoadhshQuery, "updateLoadhshQuery")
	updateLoadshErr := tx.QueryRow(updateLoadhshQuery).Err()
	fmt.Println(updateLoadshErr, "updateLoadshErr")
	if updateLoadshErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + updateLoadshErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertMainCommissions3Query := `INSERT INTO commissions (type, calc, amount, office, loadsh_id, agent_id, basis)
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
     WHERE l.id = '` + loadId + `'
     AND IFNULL((SELECT 1 FROM commissions c2 WHERE c2.loadsh_id = '` + loadId + `' AND c2.type = 'MAIN' AND c2.status <> 'DELETED' LIMIT 1), 0) < 1`

	fmt.Println(insertMainCommissions3Query, "insertMainCommissions3Query")
	insertMainCommissions3Err := tx.QueryRow(insertMainCommissions3Query).Err()
	fmt.Println(insertMainCommissions3Err, "insertMainCommissions3Err")
	if insertMainCommissions3Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertMainCommissions3Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	updateShareCommission1Query := `UPDATE commissions c
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
   	  AND tc.loadsh_id = '` + loadId + `'
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
     WHERE c.loadsh_id = '` + loadId + `'
     AND ( c.type = 'SHARE' OR c.type = 'INTERNAL_SHARE' )
     /* only update if commission date is in the future or TBD */
     AND (
		  NULLIF(c.date, '0000-00-00') IS NULL
		  OR c.date > DATE(NOW())
	  )`
	fmt.Println(updateShareCommission1Query, "updateShareCommission1Query")
	updateShareCommission1Err := tx.QueryRow(updateShareCommission1Query).Err()
	fmt.Println(updateShareCommission1Err, "updateShareCommission1Err")
	if updateShareCommission1Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + updateShareCommission1Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	updateMainCommission2Query := `UPDATE commissions c
	  JOIN loadsh l ON c.loadsh_id = l.id
	  /* Fees to be removed from profit */
	  LEFT JOIN (
		  SELECT tc.loadsh_id, SUM(tf.cust_charge) cust_charge, SUM(tf.carr_charge) carr_charge
		  FROM lscarr tc
		  JOIN lsstops ts ON ts.lscarr_id = tc.id
		  JOIN lsfees tf ON tf.lsstop_id = ts.id
		  WHERE tf.excluded_from_commissions = 1
		  AND tc.loadsh_id = '` + loadId + `'
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
		  WHERE loadsh_id = '` + loadId + `'
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
		  WHERE loadsh_id = '` + loadId + `'
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
	  WHERE c.loadsh_id = '` + loadId + `'
	  AND c.status <> 'DELETED'
	  AND c.type = 'MAIN'
	  /* only update if commission date is in the future or TBD */
	  AND (
			  NULLIF(c.date, '0000-00-00') IS NULL
			  OR c.date > DATE(NOW())
		  )`

	fmt.Println(updateMainCommission2Query, "updateMainCommission2Query")

	updateMainCommission2Err := tx.QueryRow(updateMainCommission2Query).Err()
	fmt.Println(updateMainCommission2Err, "updateMainCommission2Err")
	if updateMainCommission2Err != nil {
		  tx.Rollback()
		logger.Error("Error while inserting values in the database:" + updateMainCommission2Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissions4Query := `INSERT INTO commissions
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
				AND tc.loadsh_id = '` + loadId + `'
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
				AND loadsh_id = '` + loadId + `'
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
				WHERE loadsh.id = '` + loadId + `'
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
			AND l.id = '` + loadId + `'
		) temp_new_incentives
		WHERE temp_new_incentives.amount <> 0
	)`

	fmt.Println(insertCommissions4Query, "insertCommissions4Query")
	insertCommissions4Err := tx.QueryRow(insertCommissions4Query).Err()
	fmt.Println(insertCommissions4Err, "insertCommissions4Err")
	if insertCommissions4Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions4Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissions5Query := `INSERT INTO commissions
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
								loadsh_id = '` + loadId + `'
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
			AND tc.loadsh_id = '` + loadId + `'
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
		AND l.id = '` + loadId + `'
		GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01)`

	fmt.Println(insertCommissions5Query, "insertCommissions5Query")
	insertCommissions5Err := tx.QueryRow(insertCommissions5Query).Err()
	fmt.Println(insertCommissions5Err, "insertCommissions5Err")
	if insertCommissions5Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions5Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissions6Query := `INSERT INTO commissions
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
								loadsh_id = '` + loadId + `'
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
			AND tc.loadsh_id = '` + loadId + `'
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
		AND l.id = '` + loadId + `'
		GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	fmt.Println(insertCommissions6Query, "insertCommissions6Query")
	insertCommissions6Err := tx.QueryRow(insertCommissions6Query).Err()
	fmt.Println(insertCommissions6Err, "insertCommissions6Err")
	if insertCommissions6Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions6Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}
	insertCommissions7Query := `INSERT INTO commissions
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
								loadsh_id = '` + loadId + `'
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
			AND tc.loadsh_id = '` + loadId + `'
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
		AND l.id = '` + loadId + `'
		GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`
	insertCommissions7Err := tx.QueryRow(insertCommissions7Query).Err()
	fmt.Println(insertCommissions7Err, "insertCommissions7Err")
	if insertCommissions7Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions7Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}
	insertCommission8Query := `INSERT INTO commissions
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
								loadsh_id = '` + loadId + `'
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
			AND tc.loadsh_id = '` + loadId + `'
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
		AND l.id = '` + loadId + `'
		GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`
	insertCommissions8Err := tx.QueryRow(insertCommission8Query).Err()
	fmt.Println(insertCommissions8Err, "insertCommissions8Err")
	if insertCommissions8Err != nil {
		// tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions8Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertCommissions9Query := `INSERT INTO commissions
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
								loadsh_id = '` + loadId + `'
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
			AND tc.loadsh_id = '` + loadId + `'
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
		AND l.id = '` + loadId + `'
		GROUP BY c.id
	) adjustments
	-- only log when commission amount changes by more than 1 penny
	WHERE total NOT BETWEEN -0.01 AND 0.01
	)`

	insertCommissions9Err := tx.QueryRow(insertCommissions9Query).Err()
	fmt.Println(insertCommissions9Err, "insertCommissions9Err")
	if insertCommissions9Err != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertCommissions9Err.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	updateCommissionIncentivesQuery := `UPDATE commission_incentives ci
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
		WHERE commissions.loadsh_id = '` + loadId + `'
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

	updateCommissionIncentivesErr := tx.QueryRow(updateCommissionIncentivesQuery).Err()
	fmt.Println(updateCommissionIncentivesErr, "updateCommissionIncentivesErr")
	if updateCommissionIncentivesErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + updateCommissionIncentivesErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	insertLoashshPostQuery := `INSERT INTO loadsh_post 
	(carrier_email, comments,
    created_at, deleted_at,
	emailed_at, id,
	is_matched, last_error,
	last_error_payload, last_job_params,
	last_operation, last_payload,
	load_board_id, load_id,
	posted_at, posted_id,
	posted_load, updated_at)
	VALUES 
	(NULL, NULL, 
	current_timestamp, NULL,
	NULL, '0',
	'0', NULL,
	NULL, NULL,
	NULL, NULL,
	'1', '` + loadId + `',
	NULL, NULL,
	NULL, NULL)`

	insertLoashshPostErr := tx.QueryRow(insertLoashshPostQuery).Err()
	fmt.Println(insertLoashshPostErr, "insertLoashshPostErr")
	if insertLoashshPostErr != nil {
		tx.Rollback()
		logger.Error("Error while inserting values in the database:" + insertLoashshPostErr.Error())
		return nil, errs.UserNewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
	}

	fmt.Println(loadId, "loadId")
	//ps_1.22.1 The response is returned to the service if no error occurs
	tx.Commit()
	dbclnt.Close()
	return &dto.TLresponse{
		TLId: loadId,
		AgentDetails: Customerdetail,
	}, nil
}

func NewTLRepositoryDb() RepositoryDb {
	return RepositoryDb{}
}

//ps_1.16.1 getUserSummaryRepo is defined here with its request got from the service and response returned to the service
// func (d RepositoryDb) GetUserSummaryRepo(userReq dto.UserReq) ([]UserSummaryRes, *errs.AppError) {
// 	fmt.Println("entered Repository")

// 	//ps_1.16.2 user variable is used to store the summary response
// 	userres := make([]UserSummaryRes, 0)
// 	dbclnt, err := database.GetCpdbClient()
// 	fmt.Println("entered if repo")
// 	if err != nil {
// 		fmt.Println("error in dbconnection GetUserSummaryRepo")
// 		return nil, errs.NewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request 123.")
// 	}
// 	var errors error

// 	//ps_1.16.3 & ps_1.16.4 query to be executed is formed using the below if statements

// 	userDQuery := "Select users.user_id,users.first_name,users.last_name,users.email_id,users.phone_number,cpdb.permission.permission_id,cpdb.permission.permission_type,cpdb.user_status.status_type,cpdb.user_customer_map.global_customer_name,cpdb.user_customer_map.global_customer_id from cpdb.users inner join cpdb.user_permission on users.user_id =user_permission.user_id inner join cpdb.permission on permission.permission_id=user_permission.permission_id inner join cpdb.user_status on user_status.user_status_id=users.user_status_id inner join cpdb.user_customer_map on user_customer_map.user_id = users.user_id where  user_permission.active = 1 and user_customer_map.active = 1 and users.active=1"
// 	//users for admin include all the users under a particular coorporate
// 	fmt.Println("Admin and tableName", userReq.Admin, userReq.TableName)

// 	if userReq.Admin == "yes" {
// 		userDQuery = userDQuery + " and user_customer_map.corporate_id = " + "'" + userReq.CoorporateId + "'"
// 	} else if userReq.TableName == "Cust_mysunteck_logins" {
// 		userDQuery = userDQuery + " and users.user_id in (select user_id from cpdb.user_customer_map where global_customer_id= '" + userReq.GCID + "')"
// 	} else {
// 		userDQuery = userDQuery + " and users.created_by = '" + userReq.UserId +"'"
// 	}
// 	if userReq.FilStatus != "" {
// 		userDQuery = userDQuery + " and users.user_status_id =" + userReq.FilStatus
// 	}
// 	if userReq.FilCus != "" {
// 		userDQuery = userDQuery + " and user_customer_map.global_customer_name = '" + userReq.FilCus + "'"
// 	}
// 	if userReq.FilFromDate != "" && userReq.FilToDate == "" {
// 		userDQuery = userDQuery + " and users.created_date >= " + "'" + userReq.FilFromDate + "'"
// 	} else if userReq.FilToDate == "" && userReq.FilFromDate != "" {
// 		userDQuery = userDQuery + " and users.created_date <= " + "'" + userReq.FilToDate + "'"
// 	} else if userReq.FilToDate != "" && userReq.FilFromDate != "" {
// 		userDQuery = userDQuery + " and users.created_date >= " + "'" + userReq.FilFromDate + "'" + "and users.created_date <= " + "'" + userReq.FilToDate + "'"
// 	}
// 	if userReq.Searchparam != "" {
// 		userDQuery = userDQuery + " and (phone_number like '%" + userReq.Searchparam + "%' or email_id like '%" + userReq.Searchparam + "%' or first_name like '%" + userReq.Searchparam + "%' or last_name like '%" + userReq.Searchparam + "%')"
// 	}

// 	if userReq.SortCol == "User ID" {
// 		userReq.SortCol = "user_id"
// 		userDQuery = userDQuery + " order by users." + userReq.SortCol + " " + userReq.SortOrder
// 	} else if userReq.SortCol == "Username" {
// 		userReq.SortCol = "first_name"
// 		userDQuery = userDQuery + " order by users." + userReq.SortCol + " " + userReq.SortOrder
// 	} else if userReq.SortCol == "Email ID" {
// 		userReq.SortCol = "email_id"
// 		userDQuery = userDQuery + " order by users." + userReq.SortCol + " " + userReq.SortOrder
// 	} else if userReq.SortCol == "Phone Number" {
// 		userReq.SortCol = "phone_number"
// 		userDQuery = userDQuery + " order by users." + userReq.SortCol + " " + userReq.SortOrder
// 	} else if userReq.SortCol == "Status" {
// 		userReq.SortCol = "status_type"
// 		userDQuery = userDQuery + " order by user_status." + userReq.SortCol + " " + userReq.SortOrder
// 	} else if userReq.SortCol == "" && userReq.SortOrder == "" {
// 		userReq.SortCol = "created_date"
// 		userReq.SortOrder = "Asc"
// 		// userDQuery = userDQuery + " order by users.$3 $4"
// 		userDQuery = userDQuery + " order by users." + userReq.SortCol + " " + userReq.SortOrder
// 	}

// 	//userID, _ := strconv.Atoi(userReq.UserId)
// 	fmt.Print(userReq.UserId)
// 	//ps_1.17.1 The query is executed using the below function
// 	fmt.Println("userDQuery", userDQuery)
// 	resulterr := dbclnt.Select(&userres, userDQuery)
// 	fmt.Println("repo err", resulterr)

// 	if errors != nil {
// 		//ps_1.17.2 the if statement is used to check for the error in the db Response
// 		logger.Error("Error while selecting values in the database:" + resulterr.Error())
// 		return nil, errs.NewUnexpectedError("Oops! The server encountered a temporary error and could not complete your request.")
// 	}

// 	fmt.Println("rep userres", userres)
// 	dbclnt.Close()
// 	//ps_1.18.1 the response is returned to the service
// 	return userres, nil
// }

// sess := session.Must(session.NewSessionWithOptions(session.Options{
// 	SharedConfigState: session.SharedConfigEnable,
// }))

// // Create DynamoDB client
// svc := dynamodb.New(sess)
// // snippet-end:[dynamodb.go.read_item.session]

// // snippet-start:[dynamodb.go.read_item.call]
// // customerportal-dev-gcdb-table
// //customerportal-dev-gcdb-table-test
// tableName := "customerportal-dev-gcdb-table-test"
// GlobalCustomerId := req.GCID

// result, err := svc.GetItem(&dynamodb.GetItemInput{
// 	TableName: aws.String(tableName),
// 	Key: map[string]*dynamodb.AttributeValue{
// 		"ID": {
// 			S: aws.String(GlobalCustomerId),
// 		},
// 	},
// })
// fmt.Println(result, "GCDB Response\n")

// if err != nil {
// 	log.Fatalf("Got error calling GetItem: %s", err.Error())
// }
// // snippet-end:[dynamodb.go.read_item.call]

// // snippet-start:[dynamodb.go.read_item.unmarshall]
// if result.Item == nil {
// 	//msg := "Could not find '" + *title + "'"
// 	return nil, errs.UserAuthentication("Access Denied no item returned")
// 	//return nil, errors.New(msg)
// }

// fmt.Println("item unmarshling \n")
// // item := Item{}
// var item Item

// err = dynamodbattribute.UnmarshalMap(result.Item, &item)
// fmt.Println(result.Item, "result.Item\n")

// if err != nil {
// 	panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
// }

// fmt.Println("Found item:")
// fmt.Printf("Found item: %+v\n", item)
// fmt.Println("customer name :  ", item.CustomerName)
// //fmt.Println("Title: ", item.Title)
// //fmt.Println("Plot:  ", item.Plot)
// //fmt.Println("Rating:", item.Rating)
// // snippet-end:[dynamodb.go.read_item.unmarshall]

// if item.ID != "" {
// 	fmt.Println(item.ID, "Enterned GCID check\n")
// 	verifyGCIDResponse.GCID = true
// }
// for _, cust := range item.SourceInfo {
// 	if cust.SourceCustomerID == req.Customer {
// 		fmt.Println(item.ID, "Enterned GCID_cust check\n")
// 		verifyGCIDResponse.GCIDCusCheck = true
// 		break
// 	}
// }
