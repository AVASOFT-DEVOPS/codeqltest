package domain

import (
	"golang/dto"
	"golang/errs"
)

type LoadDetailsRepository interface {
	GetOTMSLoadDetailsRepo(loadDetailsVar dto.GetLoadDetailsReqDTO) ([]LoadShipConsRefDTO, []LoadCommoditiesDTO, []LoadCPDBResDTO, *errs.AppErrorvalidation)
	GetBTMSLoadDetailsRepo(loadDetailsVar dto.GetLoadDetailsReqDTO) ([]LoadShipConsRefDTO, []LoadCommoditiesDTO, []EventTrackingUpdatesTL, []LocationBreadCrumbs, []EventTrackingUpdatesELTL, *errs.AppErrorvalidation)
}

type LoadShipConsRefDTO struct {
	LoadId                 string  `db:"loadId"`
	LoadStatus             string  `db:"loadStatus"`
	LoadMethod             string  `db:"loadMethod"`
	TrackingEnabled        *int    `db:"trackingEnabled"`
	ShipperName            string  `db:"shipperName"`
	ShipperAddressLine1    string  `db:"shipperAddressLine1"`
	ShipperAddressLine2    *string `db:"shipperAddressLine2"`
	ShipperCity            string  `db:"shipperCity"`
	ShipperState           string  `db:"shipperState"`
	ShipperZip             string  `db:"shipperZip"`
	EarliestShipmentsDate  *string `db:"earliestShipmentsDate"`
	EarliestShipmentsTime  *string `db:"earliestShipmentsTime"`
	LatestShipmentsDate    *string `db:"latestShipmentsDate"`
	LatestShipmentsTime    *string `db:"latestShipmentsTime"`
	ShipperDriverinDate    *string `db:"shipperDriverinDate"`
	ShipperDriverOutDate   *string `db:"shipperDriverOutDate"`
	ShipperDriverinTime    *string `db:"shipperDriverinTime"`
	ShipperDriverOutTime   *string `db:"shipperDriverOutTime"`
	ConsigneeName          string  `db:"consigneeName"`
	ConsigneeAddressLine1  string  `db:"consigneeAddressLine1"`
	ConsigneeAddressLine2  *string `db:"consigneeAddressLine2"`
	ConsigneeCity          string  `db:"consigneeCity"`
	ConsigneeState         string  `db:"consigneeState"`
	ConsigneeZip           string  `db:"consigneeZip"`
	EarliestConsigneeDate  *string `db:"earliestConsigneeDate"`
	EarliestConsigneeTime  *string `db:"earliestConsigneeTime"`
	LatestConsigneeDate    *string `db:"latestConsigneeDate"`
	LatestConsigneeTime    *string `db:"latestConsigneeTime"`
	ConsigneeDriverinDate  *string `db:"consigneeDriverinDate"`
	ConsigneeDriverOutDate *string `db:"consigneeDriverOutDate"`
	ConsigneeDriverinTime  *string `db:"consigneeDriverinTime"`
	ConsigneeDriverOutTime *string `db:"consigneeDriverOutTime"`
	PoNumber               *string `db:"poNumber"`
	ShipBlNumber           *string `db:"shipBlNumber"`
	ProNumber              *string `db:"proNumber"`
	ShipperNumber          *string `db:"shipperNumber"`
	PickupNumber           *string `db:"pickupNumber"`
	DeliveryNumber         *string `db:"deliveryNumber"`
	LoadOrigin             string  `db:"loadOrigin"`
}

type LoadCommoditiesDTO struct {
	ItemId               int      `db:"itemId"`
	Hazmat               *string  `db:"hazmat"`
	ItemQuantity         *float32 `db:"itemQuantity"`
	UnitOfMeasure        *string  `db:"unitOfMeasure"`
	Weight               *float32 `db:"weight"`
	ItemDescription      *string  `db:"itemDescription"`
	ItemValue            *int     `db:"itemValue"`
	Class                *string  `db:"class"`
	Nmfc                 *string  `db:"nmfc"`
	PalletLength         *float32 `db:"palletLength"`
	PalletWidth          *float32 `db:"palletWidth"`
	PalletHeight         *float32 `db:"palletHeight"`
	ItemDensity          *float32 `db:"itemDensity"`
	HazmatContact        *string  `db:"hazmatContact"`
	HazmatGroupName      *string  `db:"hazmatGroupName"`
	HazmatPackagingGroup *string  `db:"hazmatPackagingGroup"`
	HazmatUNNAnumber     *string  `db:"hazmatUNNAnumber"`
	HazmatClass          *string  `db:"hazmatClass"`
	HazmatPlacard        *string  `db:"hazmatPlacard"`
	HazmatFlashTemp      *int     `db:"hazmatFlashTemp"`
	HazmatFlashType      *string  `db:"hazmatFlashType"`
	HazmatUom            *string  `db:"hazmatUom"`
	HazmatCertHolderName *string  `db:"hazmatCertHolderName"`
	HazmatContactName    *string  `db:"hazmatContactName"`
	HazmatPhoneNumber    *string  `db:"hazmatPhoneNumber"`
}

func (s LoadCommoditiesDTO) ToLoadCommodityDto() dto.LoadCommoditiesDTO {

	return dto.LoadCommoditiesDTO{
		ItemId:               s.ItemId,
		Hazmat:               s.Hazmat,
		ItemQuantity:         s.ItemQuantity,
		UnitOfMeasure:        s.UnitOfMeasure,
		Weight:               s.Weight,
		ItemDescription:      s.ItemDescription,
		ItemValue:            s.ItemValue,
		Class:                s.Class,
		Nmfc:                 s.Nmfc,
		PalletLength:         s.PalletLength,
		PalletWidth:          s.PalletWidth,
		PalletHeight:         s.PalletHeight,
		ItemDensity:          s.ItemDensity,
		HazmatContact:        s.HazmatContact,
		HazmatGroupName:      s.HazmatGroupName,
		HazmatPackagingGroup: s.HazmatPackagingGroup,
		HazmatUNNAnumber:     s.HazmatUNNAnumber,
		HazmatClass:          s.HazmatClass,
		HazmatPlacard:        s.HazmatPlacard,
		HazmatFlashTemp:      s.HazmatFlashTemp,
		HazmatFlashType:      s.HazmatFlashType,
		HazmatUom:            s.HazmatUom,
		HazmatCertHolderName: s.HazmatCertHolderName,
		HazmatContactName:    s.HazmatContactName,
		HazmatPhoneNumber:    s.HazmatPhoneNumber,
	}
}

type EventTrackingUpdatesTL struct {
	EventId       int    `db:"eventId"`
	LsstopId      int    `db:"lsstopId"`
	EventDateTime string `db:"eventDateTime"`
	EventStatus   string `db:"eventStatus"`
	City          string `db:"city"`
	State         string `db:"state"`
	Country       string `db:"country"`
}

//

func (s EventTrackingUpdatesTL) ToLoadEventsTlDto() dto.EventTrackingUpdates {
	return dto.EventTrackingUpdates{
		EventId:       s.EventId,
		LsstopId:      s.LsstopId,
		EventDateTime: s.EventDateTime,
		EventStatus:   s.EventStatus,
		City:          s.City,
		State:         s.State,
		Country:       s.Country,
	}
}

type LocationBreadCrumbs struct {
	LocationId          int     `db:"locationId"`
	LocationUpdatedDate string  `db:"locationUpdatedDate"`
	ThirdParty          string  `db:"thirdParty"`
	City                string  `db:"city"`
	State               string  `db:"state"`
	Zip                 string  `db:"zip"`
	Country             string  `db:"country"`
	Latitude            float64 `db:"latitude"`
	Longitude           float64 `db:"longitude"`
	DriverMobile        string  `db:"driverMobile"`
}

func (s LocationBreadCrumbs) ToLoadLocationDto() dto.LocationBreadCrumbs {
	return dto.LocationBreadCrumbs{
		LocationId:          s.LocationId,
		LocationUpdatedDate: s.LocationUpdatedDate,
		ThirdParty:          s.ThirdParty,
		City:                s.City,
		State:               s.State,
		Zip:                 s.Zip,
		Country:             s.Country,
		Latitude:            s.Latitude,
		Longitude:           s.Longitude,
		DriverMobile:        s.DriverMobile,
	}
}

type EventTrackingUpdatesELTL struct {
	LogGroup       *int    `db:"logGroup"`
	CustCarr       *string `db:"custCarr"`
	SentRecordDate *string `db:"sentRecordDate"`
	TradingPartner *string `db:"tradingPartner"`
	StatusDateTime *string `db:"statusDateTime"`
	User           *string `db:"user"`
	InOut          *string `db:"inOut"`
	EdiType        *string `db:"ediType"`
	EdiStatus      *string `db:"ediStatus"`
	ReasonCode     *string `db:"reasonCode"`
	LoadStatus     *string `db:"loadStatus"`
	TransmitStatus *string `db:"transmitStatus"`
	SourceTable    *string `db:"sourceTable"`
	SourceId       *string `db:"sourceId"`
}

func (s EventTrackingUpdatesELTL) ToLoadEventsELTLDto() dto.EventTrackingUpdates {
	return dto.EventTrackingUpdates{
		LogGroup:       s.LogGroup,
		SentRecordDate: s.SentRecordDate,
		TradingPartner: s.TradingPartner,
		StatusDateTime: s.StatusDateTime,
		User:           s.User,
		EdiType:        s.EdiType,
		EdiStatus:      s.EdiStatus,
		ReasonCode:     s.ReasonCode,
		SourceTable:    s.SourceTable,
	}
}
