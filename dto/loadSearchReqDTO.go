package dto

//import _ "github.com/go-playground/validator/v10"

type LoadSearchReqDTO struct {
	CustomerId       []int    `json:"customerId" `
	Gcid             []string `json:"gcid" `
	OfficeCode       []string `json:"officeCode" `
	ShipDateStart    string   `json:"shipDateStart"`
	ShipDateEnd      string   `json:"shipDateEnd"`
	LoadStatus       string   `json:"loadStatus" validate:"omitempty,alpha"`
	LoadMethod       string   `json:"loadMethod" validate:"omitempty,alpha"`
	LoadId           string   `json:"loadId" validate:"omitempty,numeric"`
	PoNumber         string   `json:"poNumber"`
	CustomerblNumber string   `json:"customerblNumber"`
	PickupNumber     string   `json:"pickupNumber"`
	DeliveryNumber   string   `json:"deliveryNumber"`
	SortOrder        string   `json:"sortOrder"`
	SortColumn       string   `json:"sortColumn"`
	OtmsLoadRecords  string   `json:"otmsLoadRecords" validate:"omitempty,numeric" `
	BtmsLoadRecords  string   `json:"btmsLoadRecords" validate:"omitempty,numeric"`
}

type ValidateResDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
