package dto

type LoadSearchResultResDTO struct {
	LoadResponse []LoadSearchResDTO `json:"loadResponse"`
	LoadCount    struct {
		OthertmsLoadRecords int `json:"othertmsLoadRecords"`
		BtmsLoadRecords     int `json:"btmsLoadRecords"`
		TotalLoadsCount     int `json:"totalLoadsCount"`
	} `json:"loadCount"`
}

type LoadSearchResDTO struct {
	ShipDateStart  string  `json:"shipDateStart"`
	ShipDateEnd    string  `json:"shipDateEnd"`
	LoadStatus     string  `json:"loadStatus"`
	LoadMethod     string  `json:"loadMethod"`
	LoadId         string  `json:"loadID"`
	PickupNumber   *string `json:"pickupNumber"`
	DeliveryNumber *string `json:"deliveryNumber"`
	CustTotal      string  `json:"custTotal"`
	ShipperCity    string  `json:"shipperCity"`
	ConsigneeCity  string  `json:"consigneeCity"`
	ShipperState   string  `json:"shipperState"`
	ConsigneeState string  `json:"consigneeState"`
	ProNumber      *string `json:"proNumber"`
	PoNumber       *string `json:"poNumber"`
	ShipBlNumber   *string `json:"shipBlNumber"`
	InvoiceDate    *string `json:"invoiceDate"`
	LastModified   string  `json:"lastModified"`
	LoadOrigin     string  `json:"loadOrigin"`
}
