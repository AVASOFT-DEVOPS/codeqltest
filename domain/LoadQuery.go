package domain

//Postgres Query to get the shipper consignee and reference numbers from CP DB are written below
var LoadShipConsRefQuery = `select ls.source_load_shipment_id as "loadId", lsd.description as "loadMethod", lst.status as "loadStatus", ls.shipper_name as "shipperName",
ls.shipper_address_line_1 as "shipperAddressLine1", ls.shipper_address_line_2 as "shipperAddressLine2", ls.shipper_city as "shipperCity",
ls.shipper_state as "shipperState", ls.shipper_zip as "shipperZip", ls.consignee_name as "consigneeName", ls.consignee_address_line_1 as "consigneeAddressLine1",
ls.consignee_address_line_2 as "consigneeAddressLine2", ls.consignee_city as "consigneeCity", ls.consignee_state
as "consigneeState", ls.consignee_zip as "consigneeZip", lssPICK.earliest_date as "earliestShipmentsDate", lssPICK.earliest_time as "earliestShipmentsTime",
lssPICK.latest_date as "latestShipmentsDate", lssPICK.latest_time as "latestShipmentsTime", lssPICK.driver_in_date as "shipperDriverinDate",
lssPICK.driver_out_date as "shipperDriverOutDate", lssPICK.driver_in_time as "shipperDriverinTime", lssPICK.driver_out_time as "shipperDriverOutTime",
 lssDROP.earliest_date as "earliestConsigneeDate", lssDROP.earliest_time as "earliestConsigneeTime",
lssDROP.latest_date as "latestConsigneeDate", lssDROP.latest_time as "latestConsigneeTime", lssDROP.driver_in_date as "consigneeDriverinDate",
lssDROP.driver_out_date as "consigneeDriverOutDate", lssDROP.driver_in_time as "consigneeDriverinTime", lssDROP.driver_out_time as "consigneeDriverOutTime",
 ls.customer_purchase_order as "poNumber", ls.shipper_bill_of_ladding as "shipBlNumber", lsc.pro_number as "proNumber",
 ls.customer_shipment_id as "shipperNumber", lssPICK.po_number as "pickupNumber", lssDROP.po_number as "deliveryNumber", 'otherTMS' as "loadOrigin"
from cpdb.load_shipment ls
inner join cpdb.load_status lst on ls.load_status_id = lst.load_status_id
inner join cpdb.load_method lsd on ls.load_method_id = lsd.load_method_id
inner join cpdb.load_carrier lsc on lsc.load_shipment_id = ls.load_shipment_id
left join cpdb.load_shipment_stops lssPICK on lssPICK.load_carrier_id = lsc.load_carrier_id and lssPICK.load_shipment_stops_type_id = 1
left join cpdb.load_shipment_stops lssDROP on lssDROP.load_carrier_id = lsc.load_carrier_id and lssDROP.load_shipment_stops_type_id = 2
where ls.source_load_shipment_id = `

//Query to get the load commodity details are writtened below
var LoadCommodityQuery = `select li.load_item_id as "itemId",li.hazmat, li.quantity as "itemQuantity", uom.unit_of_measure_desc as "unitOfMeasure", li.weight, li.description as "itemDescription",
li.value as "itemValue",  li.class, li.nmfc, li.pallet_length as "palletLength", li.pallet_width as "palletWidth", li.pallet_height as "palletHeight",
li.density as "itemDensity", li.hazmat_ship_name as "hazmatContact", li.hazmat_group_name as "hazmatGroupName", li.hazmat_package_group as "hazmatPackagingGroup",
li.hazmat_un_na_number as "hazmatUNNAnumber", li.hazmat_class as "hazmatClass", li.hazmat_placard as "hazmatPlacard",
li.hazmat_flash_temp as "hazmatFlashTemp", li.hazmat_flash_type as "hazmatFlashType", hfu.hazmat_flash_uom_desc as "hazmatUom",
li.hazmat_cert_holder as "hazmatCertHolderName", li.hazmat_contact_name as "hazmatContactName", li.hazmat_contact_telephone as "hazmatPhoneNumber"
from cpdb.load_items li
left join cpdb.unit_of_measure uom on uom.unit_of_measure_id = li.unit_of_measure_id
left join cpdb.hazmat_flash_uom hfu on hfu.hazmat_flash_uom_id = li.hazmat_flash_uom_id
inner join cpdb.load_shipment_stops lss on lss.load_shipment_stops_id = li.load_shipment_stops_id
inner join cpdb.load_carrier lc on lc.load_carrier_id = lss.load_carrier_id
inner join cpdb.load_shipment ls on ls.load_shipment_id = lc.load_shipment_id
where ls.source_load_shipment_id = `

//select Query to get the document URL
var LoadDocsQuery = `select ld.document_id as "id", ld.document_type_id as "typeId", ls.source_load_shipment_id as "loadId" , ld.document_name as "docName", ld.document_url as "docUrl"
from cpdb.document ld
inner join cpdb.document_type dt on dt.document_type_id = ld.document_type_id
inner join cpdb.load_shipment ls on ls.load_shipment_id = ld.load_shipment_id 
WHERE ls.source_load_shipment_id= `

var BtmsLoadShipConsRefQuery = `SELECT cast(loadsh.id as char) as loadId, loadsh.load_method as "loadMethod", loadsh.status as "loadStatus", loadsh.show_tracking_events_on_mysuntecktts as "trackingEnabled", lsstops.name as "shipperName",
lsstops.addr1 as "shipperAddressLine1",lsstops.addr2 as "shipperAddressLine2", lsstops.city as "shipperCity",
lsstops.state as "shipperState", lsstops.zip as "shipperZip", last_stop.name as "consigneeName", last_stop.addr1  as "consigneeAddressLine1",
last_stop.addr2 as "consigneeAddressLine2", last_stop.city as "consigneeCity", last_stop.state as "consigneeState", last_stop.zip as "consigneeZip", 
lsstops.date as "earliestShipmentsDate", lsstops.time as "earliestShipmentsTime", cast(lsstops.date2 as char) as "latestShipmentsDate", lsstops.time2 as "latestShipmentsTime", 
cast(lsstops.datein as char) as "shipperDriverinDate", cast(lsstops.dateout as char)as "shipperDriverOutDate", lsstops.timein as "shipperDriverinTime", lsstops.timeout as "shipperDriverOutTime", 
cast(last_stop.date as char) as "earliestConsigneeDate", last_stop.time as "earliestConsigneeTime", cast(last_stop.date2 as char) as "latestConsigneeDate", last_stop.time2 as "latestConsigneeTime", 
cast(last_stop.datein as char) as "consigneeDriverinDate", cast(last_stop.dateout as char) as "consigneeDriverOutDate", last_stop.timein as "consigneeDriverinTime", 
last_stop.timeout as "consigneeDriverOutTime", loadsh.cust_po as "poNumber", loadsh.ship_bl as "shipBlNumber", first_carr.pro_number as "proNumber", 
loadsh.cust_shipid as "shipperNumber", lsstops.po as "pickupNumber", last_stop.po as "deliveryNumber", 'BTMS' as "loadOrigin"
	FROM loadsh
	LEFT JOIN lscarr first_carr ON first_carr.loadsh_id=loadsh.id AND first_carr.carr_order=1
	LEFT JOIN lscarr last_carr ON last_carr.loadsh_id=loadsh.id AND last_carr.carr_order=(SELECT MAX(carr_order) FROM lscarr WHERE loadsh_id=loadsh.id AND type='carrier')
	LEFT JOIN lsstops ON lsstops.lscarr_id=first_carr.id AND lsstops.stop_order=1
	LEFT JOIN lsstops last_stop ON last_stop.lscarr_id=last_carr.id AND last_stop.stop_order=(SELECT MAX(stop_order) FROM lsstops WHERE lscarr_id=last_carr.id) 
    where loadsh.id =`

var BtmsLoadCommodityQuery = `SELECT li.id as "itemId",li.hazmat, li.qty as "itemQuantity", li.type as "unitOfMeasure", li.weight, li.descrip as "itemDescription",
li.value as "itemValue",  li.class, li.nmfc, li.pallet_length as "palletLength", li.pallet_width as "palletWidth", li.pallet_height as "palletHeight",
li.density as "itemDensity", li.hazmat_ship_name as "hazmatContact", li.hazmat_group_name as "hazmatGroupName", li.hazmat_pkg_group as "hazmatPackagingGroup",
li.hazmat_un_na_num as "hazmatUNNAnumber", li.hazmat_class as "hazmatClass", li.hazmat_placard as "hazmatPlacard",
li.hazmat_flash_temp as "hazmatFlashTemp", li.hazmat_flash_type as "hazmatFlashType", li.hazmat_flash_uom as "hazmatUom",
li.hazmat_cert_holder as "hazmatCertHolderName", li.hazmat_contact_name as "hazmatContactName", li.hazmat_contact_tel as "hazmatPhoneNumber"
FROM lsitems li
inner JOIN lsstops ON lsstops.id = li.lsstop_id
inner JOIN lscarr ON lscarr.id = lsstops.lscarr_id and lsstops.stop_order = 1
WHERE lscarr.loadsh_id = `

var LoadEventsUpdateQueryTL = `SELECT esu.id as "eventId", esu.lsstop_id as "lsstopId", esu.event_date as "eventDateTime", esu.status as "eventStatus", esu.city, esu.state, esu.country
FROM edispatch_status_updates as esu 
WHERE  ((latitude <> 0 AND longitude <> 0) OR (city <> '' AND state <> '')) 
AND (third_party IN ('BTMS USER','BANYAN') OR (third_party='FOURKITES' AND fk_type IN ('STOP_ARRIVAL','STOP_DEPARTURE')) OR (third_party = 'MACROPOINT' AND fk_type = 'TRACKING_UPDATE')) AND 
loadsh_id = `

var LocationBreadCrumbsQuery = `SELECT esu.id as "locationId", esu.created as "locationUpdatedDate", esu.third_party as "thirdParty", esu.city, esu.state, esu.country, esu.latitude,
 esu.longitude, esu.carr_cell as "driverMobile"
FROM edispatch_status_updates esu
WHERE fk_type = 'LOCATION_UPDATE' AND (latitude <> 0 AND longitude <> 0) 
AND loadsh_id = `

var LoadEventsUpdateQueryELTL = `SELECT
NULL as "logGroup",
'Cust' as "custCarr",
edi_204_in.created as "sentRecordDate",
COALESCE(cust_master.name, loadsh.cust_name) as "tradingPartner",
NULL as "statusDateTime",
edi_204_in.acdc_user as user,
'in' as "inOut",
'204' as "ediType",
edi_204_in.edi_purpose as "ediStatus",
NULL as "reasonCode",
NULL as "loadStatus",
'OK - Auto' as "transmitStatus",
'edi_204_in' as "sourceTable",
edi_204_in.id as "sourceId"
FROM
edi_204_in
INNER JOIN loadsh ON
loadsh.id = edi_204_in.loadsh_id
LEFT OUTER JOIN cust_master ON
cust_master.id = edi_204_in.custm_id_built
WHERE
edi_204_in.loadsh_id = ?
-- incoming customer EDI 204 cancelations; original EDI not easily accessible in logs, so leave source_table/id blank
UNION
SELECT
NULL as "logGroup",
'Cust' as "custCarr",
edi_204_in.cancelled as "sentRecordDate",
COALESCE(cust_master.name, loadsh.cust_name) as "tradingPartner",
NULL as "statusDateTime",
'SYSTEM' as user,
'in' as "inOut",
'204' as "ediType",
'Cancel' as "ediStatus",
NULL as "reasonCode",
NULL as "loadStatus",
'OK - Auto' as "transmitStatus",
NULL as "sourceTable",
NULL as "sourceId"
FROM
edi_204_in
INNER JOIN loadsh ON
loadsh.id = edi_204_in.loadsh_id
LEFT OUTER JOIN cust_master ON
cust_master.id = edi_204_in.custm_id_built
WHERE
edi_204_in.cancelled>'0000-00-00'
AND edi_204_in.loadsh_id = ?
-- incoming EDI 214s
UNION
SELECT
edi_status_msg_loop.lsstop_id as "logGroup",
CASE
	WHEN edi_status_msg.carr_id IS NOT NULL THEN 'Carr'
	WHEN edi_status_msg.custm_id IS NOT NULL THEN 'Cust'
	ELSE ''
END as "custCarr",
edi_status_msg.created as "sentRecordDate",
CASE
	WHEN edi_status_msg.carr_id IS NOT NULL THEN carriers.name
	WHEN edi_status_msg.custm_id IS NOT NULL THEN cust_master.name
	ELSE ''
END as "tradingPartner",
edi_status_msg_loop.event_date as status_datetime,
'SYSTEM' as user,
'in' as "inOut",
edi_status_msg.edi_type as "ediType",
edi_status_msg_loop.code as "ediStatus",
CONCAT(
	edi_status_msg_loop.delay_code,
	CASE WHEN edi_reason_codes.code IS NOT NULL THEN CONCAT(' - ', edi_reason_codes.description)
	ELSE '' END
) as "reasonCode",
CASE
	WHEN edi_status_msg_loop.nochange = 1
	AND edi_status_msg.edi_type <> '322' THEN NULL
	ELSE edi_status_msg_loop.fats_status_change
END as "loadStatus",
'OK - Auto' as "transmitStatus",
'edi_status_msg' as "sourceTable",
edi_status_msg.id as "sourceId"
FROM
edi_status_msg_loop
INNER JOIN edi_status_msg ON
edi_status_msg_loop.edi_214_in_id = edi_status_msg.id
LEFT OUTER JOIN carriers ON
carriers.id = edi_status_msg.carr_id
LEFT OUTER JOIN cust_master ON
cust_master.id = edi_status_msg.custm_id
LEFT OUTER JOIN edi_reason_codes ON
edi_reason_codes.code = edi_status_msg_loop.delay_code
WHERE
edi_status_msg.loadsh_id = ?
AND edi_status_msg_loop.fats_status_change <> 'LOCATION UPDATE'
-- per FD-3146, show on the Tracking tab instead
-- incoming webservice status updates (Banyan/Bluegrace and Interop)
UNION
SELECT
edispatch_status_updates.lsstop_id as "logGroup",
'Carr' as "custCarr",
edispatch_status_updates.created as "sentRecordDate",
IFNULL(
COALESCE(carriers.name, LC.carr_name, LC2.carr_name),
edispatch_status_updates.third_party
) as "tradingPartner",
IF(edispatch_status_updates.event_date = '0000-00-00 00:00:00',
NULL,
edispatch_status_updates.event_date) as "statusDateTime",
COALESCE(agents.name, edispatch_status_updates.third_party) as "user",
'in' as "inOut",
IF(edispatch_status_updates.status IN ('Accepted', 'Rejected', 'Declined'), '990', 'API') as "ediType",
CONCAT(
IFNULL(edispatch_status_updates.fk_type, ''),
IF(IFNULL(edispatch_status_updates.fk_type, '')= '' OR IFNULL(edispatch_status_updates.status, '')= '', '', ': '),
IFNULL(edispatch_status_updates.status, '')
) as "ediStatus",
edispatch_status_updates.message as "reasonCode",
IF(edispatch_status_updates.fats_status_change = '*INOUT*',
'',
edispatch_status_updates.fats_status_change) as "loadStatus",
'OK - Auto' as "transmitStatus",
'edi_log' as "sourceTable",
NULL as "sourceId"
FROM
edispatch_status_updates
INNER JOIN loadsh ON
loadsh.id = edispatch_status_updates.loadsh_id
LEFT OUTER JOIN carriers ON
carriers.id = edispatch_status_updates.carr_id
LEFT OUTER JOIN lscarr LC ON
LC.id = edispatch_status_updates.lscarr_id
LEFT OUTER JOIN lscarr LC2 ON
LC2.loadsh_id = loadsh.id
AND loadsh.carr_count = 1
LEFT OUTER JOIN agents ON
agents.id = edispatch_status_updates.user_id
WHERE
edispatch_status_updates.loadsh_id = ?
AND IFNULL(edispatch_status_updates.fk_type, '')<> 'LOCATION_UPDATE'
AND IFNULL(edispatch_status_updates.fats_status_change, '')<> 'LOCATION UPDATE'
-- per FD-3173 and FD-3289, show on the Tracking tab instead
AND edispatch_status_updates.edi_type NOT IN ('824')
AND edispatch_status_updates.third_party <> 'BTMS USER'
-- all outgoing EDI, whether to customers or carriers, plus any incoming carrier EDI that is NOT 214
UNION
SELECT
CASE
	WHEN edi_214_location.id THEN NULL /* FD-3111 - move log entries to Load section */
	WHEN load_status_codes.id THEN
	-- EDI 214; determine whether load-level or stop-level
	CASE
		WHEN IFNULL(load_status_codes.load_status, '') <> '' THEN
		-- actual status string was logged (newer entry); determine load-vs-stop level based on that (preferred)
		CASE
			WHEN load_status_codes.load_status NOT IN ('PICK-UP APPT', 'AT ORIGIN', 'PICKED-UP', 'DELIVERY APPT', 'AT DESTINATION', 'DELIVERED', 'DELIVERED FINAL') THEN NULL
			ELSE load_status_codes.stop_id
		END
		ELSE
		-- status string not logged (historical entry); guess at load-vs-stop based on logged status CODE and current status code mapping
		CASE
			WHEN load_status_codes.status_code NOT IN ('') THEN NULL
			ELSE load_status_codes.stop_id
		END
	END
	ELSE NULL
END as "logGroup",
CASE
	WHEN carriers.id IS NULL THEN 'Cust'
	ELSE 'Carr'
END as "custCarr",
edi_log.sent as "sentRecordDate",
COALESCE(carriers.name,
		CONCAT(
		CASE WHEN edi_log.cust_order IS NULL THEN '' ELSE CONCAT('(', edi_log.cust_order, ') ') END,
		cust_master.name
		),
		loadsh.cust_name) as "tradingPartner",
load_status_codes.event_date as "statusDateTime",
IF(edi_log.user IN ('fats_minute', 'edi_dl_daemon'), 'SYSTEM', edi_log.user) as user,
CASE
	WHEN edi_log.outgoing = 1 THEN 'out'
	ELSE 'in'
END as "inOut",
edi_log.edi_type as "ediType",
CASE
	edi_log.edi_type
	WHEN '204' THEN edi_log.tender_type
	WHEN '990' THEN 'Accepted'
	WHEN '214' THEN load_status_codes.status_code
	WHEN '210' THEN 'Invoiced'
	WHEN '997' THEN 'Acknowledged'
	WHEN '824' THEN edispatch_status_updates.status
	WHEN '404' THEN edi_log.tender_type
	ELSE 'unknown'
END as "ediStatus",
CASE
	edi_log.edi_type
	WHEN '824' THEN edispatch_status_updates.message
	ELSE CONCAT(
			load_status_codes.reason_code,
			CASE WHEN edi_reason_codes.code IS NOT NULL THEN CONCAT(' - ', edi_reason_codes.description)
			ELSE '' END
	)
END as "reasonCode",
load_status_codes.load_status as "loadStatus",
CONCAT (
	CASE
	WHEN edi_log.result IN ('SUCCESS', 'FTP SUCCESS') THEN 'OK'
	ELSE edi_log.result
END,
	CASE
	WHEN edi_log.auto_send = 1 THEN ' - Auto'
	ELSE ''
END
) as "transmitStatus",
'edi_log' as "sourceTable",
edi_log.id as "sourceId"
FROM
edi_log
INNER JOIN loadsh ON
loadsh.id = edi_log.load_id
LEFT OUTER JOIN cust_master ON
cust_master.id = edi_log.custm_id
LEFT OUTER JOIN carriers ON
carriers.id = edi_log.carr_id
LEFT OUTER JOIN load_status_codes ON
load_status_codes.edi_log_id = edi_log.id
LEFT OUTER JOIN edi_214_location ON
edi_214_location.edi_log_id = edi_log.id
LEFT OUTER JOIN edi_reason_codes ON
edi_reason_codes.code = load_status_codes.reason_code
LEFT OUTER JOIN edispatch_status_updates ON
edispatch_status_updates.id = edi_log.edispatch_status_update_id
WHERE
(edi_log.outgoing = 1
	OR edi_log.edi_type NOT IN ('990', '214', '322'))
-- select all outgoing EDI messages, but only non-990s/214s incoming
AND edi_log.message <> '0'
-- no 'EDI message was blank.' results
AND loadsh.id = ?
AND edi_log.vendor NOT IN ('Banyan', 'BlueGrace') `
