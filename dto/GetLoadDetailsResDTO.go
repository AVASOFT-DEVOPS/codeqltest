package dto

type LoadDetailsResDTO struct {
	LoadShipConsRef LoadShipConsRefDTO     `json:"loadShipConsRef"`
	LoadCommodities []LoadCommoditiesDTO   `json:"loadCommodities"`
	LoadDocuments   []LoadDocuments        `json:"loadDocuments"`
	LoadEvents      []EventTrackingUpdates `json:"loadEvents"`
	LocationUpdates []LocationBreadCrumbs  `json:"locationUpdates"`
}

type LoadShipConsRefDTO struct {
	LoadId                 string  `json:"loadId"`
	LoadStatus             string  `json:"loadStatus"`
	LoadMethod             string  `json:"loadMethod"`
	TrackingEnabled        *int    `json:"trackingEnabled"`
	ShipperName            string  `json:"shipperName"`
	ShipperAddressLine1    string  `json:"shipperAddressLine1"`
	ShipperAddressLine2    *string `json:"shipperAddressLine2"`
	ShipperCity            string  `json:"shipperCity"`
	ShipperState           string  `json:"shipperState"`
	ShipperZip             string  `json:"shipperZip"`
	EarliestShipmentsDate  *string `json:"earliestShipmentsDate"`
	EarliestShipmentsTime  *string `json:"earliestShipmentsTime"`
	LatestShipmentsDate    *string `json:"latestShipmentsDate"`
	LatestShipmentsTime    *string `json:"latestShipmentsTime"`
	ShipperDriverinDate    *string `json:"shipperDriverinDate"`
	ShipperDriverOutDate   *string `json:"shipperDriverOutDate"`
	ShipperDriverinTime    *string `json:"shipperDriverinTime"`
	ShipperDriverOutTime   *string `json:"shipperDriverOutTime"`
	ConsigneeName          string  `json:"consigneeName"`
	ConsigneeAddressLine1  string  `json:"consigneeAddressLine1"`
	ConsigneeAddressLine2  *string `json:"consigneeAddressLine2"`
	ConsigneeCity          string  `json:"consigneeCity"`
	ConsigneeState         string  `json:"consigneeState"`
	ConsigneeZip           string  `json:"consigneeZip"`
	EarliestConsigneeDate  *string `json:"earliestConsigneeDate"`
	EarliestConsigneeTime  *string `json:"earliestConsigneeTime"`
	LatestConsigneeDate    *string `json:"latestConsigneeDate"`
	LatestConsigneeTime    *string `json:"latestConsigneeTime"`
	ConsigneeDriverinDate  *string `json:"consigneeDriverinDate"`
	ConsigneeDriverOutDate *string `json:"consigneeDriverOutDate"`
	ConsigneeDriverinTime  *string `json:"consigneeDriverinTime"`
	ConsigneeDriverOutTime *string `json:"consigneeDriverOutTime"`
	PoNumber               *string `json:"poNumber"`
	ShipBlNumber           *string `json:"shipBlNumber"`
	ProNumber              *string `json:"proNumber"`
	ShipperNumber          *string `json:"shipperNumber"`
	PickupNumber           *string `json:"pickupNumber"`
	DeliveryNumber         *string `json:"deliveryNumber"`
	LoadOrigin             string  `json:"loadOrigin"`
}

type LoadCommoditiesDTO struct {
	ItemId               int      `json:"itemId"`
	Hazmat               *string  `json:"hazmat"`
	ItemQuantity         *float32 `json:"itemQuantity"`
	UnitOfMeasure        *string  `json:"unitOfMeasure"`
	Weight               *float32 `json:"weight"`
	ItemDescription      *string  `json:"itemDescription"`
	ItemValue            *int     `json:"itemValue"`
	Class                *string  `json:"class"`
	Nmfc                 *string  `json:"nmfc"`
	PalletLength         *float32 `json:"palletLength"`
	PalletWidth          *float32 `json:"palletWidth"`
	PalletHeight         *float32 `json:"palletHeight"`
	ItemDensity          *float32 `json:"itemDensity"`
	HazmatContact        *string  `json:"hazmatContact"`
	HazmatGroupName      *string  `json:"hazmatGroupName"`
	HazmatPackagingGroup *string  `json:"hazmatPackagingGroup"`
	HazmatUNNAnumber     *string  `json:"hazmatUNNAnumber"`
	HazmatClass          *string  `json:"hazmatClass"`
	HazmatPlacard        *string  `json:"hazmatPlacard"`
	HazmatFlashTemp      *int     `json:"hazmatFlashTemp"`
	HazmatFlashType      *string  `json:"hazmatFlashType"`
	HazmatUom            *string  `json:"hazmatUom"`
	HazmatCertHolderName *string  `json:"hazmatCertHolderName"`
	HazmatContactName    *string  `json:"hazmatContactName"`
	HazmatPhoneNumber    *string  `json:"hazmatPhoneNumber"`
}

type EventTrackingUpdates struct {
	EventId        int     `json:"eventId"`
	LsstopId       int     `json:"lsstopId"`
	EventDateTime  string  `json:"eventDateTime"`
	EventStatus    string  `json:"eventStatus"`
	City           string  `json:"city"`
	State          string  `json:"state"`
	Country        string  `json:"country"`
	LogGroup       *int    `json:"logGroup"`
	SentRecordDate *string `json:"sentRecordDate"`
	TradingPartner *string `json:"tradingPartner"`
	StatusDateTime *string `json:"statusDateTime"`
	User           *string `json:"user"`
	EdiType        *string `json:"ediType"`
	EdiStatus      *string `json:"ediStatus"`
	ReasonCode     *string `json:"reasonCode"`
	SourceTable    *string `json:"sourceTable"`
}

type LocationBreadCrumbs struct {
	LocationId          int     `json:"locationId"`
	LocationUpdatedDate string  `json:"locationUpdatedDate"`
	ThirdParty          string  `json:"thirdParty"`
	City                string  `json:"city"`
	State               string  `json:"state"`
	Zip                 string  `json:"zip"`
	Country             string  `json:"country"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	DriverMobile        string  `json:"driverMobile"`
}
