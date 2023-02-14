package dto

//import _ "github.com/go-playground/validator/v10"

type VerifyGCIDRequest struct {
	UserId          string `json:"UserId"`
	CustomerLoginId string `json:"CustomerLoginId"`
}

type CreateTLReq struct {
	GCID            string            `json:"GCID" validate:"required,alphanum"`
	CustomerId      string            `json:"CustomerId" validate:"required,numeric"`
	CustomerLoginId string            `json:"CustomerLoginId" validate:"required,numeric"`
	UserId          string            `json:"UserId" validate:"required,numeric"`
	ShipmentInfo    ShipmentInfo      `json:"ShipmentInfo,omitempty"`
	ShipperInfo     ShipperInfo       `json:"ShipperInfo,omitempty"`
	ConsigneeInfo   ConsigneeInfo     `json:"ConsigneeInfo,omitempty"`
	CommoditiesInfo []CommoditiesInfo `json:"CommoditiesInfo,omitempty"`
	Token          string             `json:"Token,omitempty"`
}
// Shipper:[Address]
// Consignee:[Address]
// Commodity: [Commodity Details]

// Contact:[]
// Phone:[]
// Fax:[]
// Email: []
// type SendMailRequest struct{
// 	Shipper string
// 	Consignee string
// 	TempPass string
// 	Commodity string
// 	Contact string
// 	Phone string
// 	Fax string
// 	Email string
// 	Token string
// }


type HazmatInfo struct {
	ShippingName   string `json:"ShippingName"`
	Description    string `json:"Description"`
	GroupName      string `json:"GroupName"`
	PackagingGroup string `json:"PackagingGroup"`
	UNNANumber     string `json:"UNNANumber"`
	Class          string `json:"Class"`
	PlacardType    string `json:"PlacardType"`
	FlashTemp      string `json:"FlashTemp"`
	UOM            string `json:"UOM"`
	FlashType      string `json:"FlashType"`
	CHName         string `json:"CHName"`
	ContactName    string `json:"ContactName"`
	PhoneNumber    string `json:"PhoneNumber"`
}
type CommoditiesInfo struct {
	BolsplIns  string     `json:"BolsplIns" validate:"omitempty"`
	Qty        int        `json:"Qty" validate:"required,numeric"`
	UOM        string     `json:"UOM" validate:"required,alpha,omitempty"`
	Weight     int        `json:"Weight" validate:"required,numeric"`
	Value      int        `json:"Value" validate:"omitempty"`
	Descrip    string     `json:"Descrip" validate:"required,alpha"`
	Hazmat     bool       `json:"Hazmat"`
	HazmatInfo HazmatInfo `json:"HazmatInfo"`
}
type ShipmentInfo struct {
	EquipType        string    `json:"EquipType" validate:"required,alpha"`
	PONumber         string    `json:"PONumber,omitempty" validate:"omitempty"`
	BLNumber         string    `json:"BLNumber,omitempty" validate:"omitempty"`
	ShippingNumber   string    `json:"ShippingNumber,omitempty" validate:"omitempty"`
	ReferrenceNumber string    `json:"ReferrenceNumber,omitempty" validate:"omitempty"`
}
type ShipperInfo struct {
	ShipName         string `json:"ShipName" validate:"required"`
	ShipAdd1         string `json:"ShipAdd1" validate:"required,alphanum"`
	ShipAdd2         string `json:"ShipAdd2" validate:"omitempty"`
	ShipCity         string `json:"ShipCity" validate:"required,alpha"`
	ShipState        string `json:"ShipState" validate:"required,alpha"`
	ShipZipcode      string `json:"ShipZipcode" validate:"required,numeric"`
	Country          string `json:"ShipCountry" validate:"required,alpha"`
	ShipContactName  string `json:"ShipContactName" validate:"omitempty"`
	ShipEmail        string `json:"ShipEmail" validate:"omitempty"`
	ShipPhoneNumber  string `json:"ShipPhoneNumber" validate:"omitempty"`
	ShipFax          string `json:"ShipFax" validate:"omitempty"`
	ShipLoadNotes    string `json:"ShipLoadNotes" validate:"omitempty"`
	ShipEarliestDate string `json:"ShipEarliestDate" validate:"required"`
	ShipEarliestTime string `json:"ShipEarliestTime" validate:"omitempty"`
	ShipLatestDate   string `json:"ShipLatestDate" validate:"omitempty"`
	ShipLatestTime   string `json:"ShipLatestTime" validate:"omitempty"`
}
type ConsigneeInfo struct {
	ConsigName         string `json:"ConsigName"        validate:"required,alpha"`
	ConsigAdd1         string `json:"ConsigAdd1"        validate:"required,alphanum"`
	ConsigAdd2         string `json:"ConsigAdd2"        validate:"omitempty"`
	ConsigCity         string `json:"ConsigCity"        validate:"required,alpha"`
	ConsigState        string `json:"ConsigState"       validate:"required,alpha"`
	ConsigZipcode      string `json:"ConsigZipcode"     validate:"required,numeric"`
	ConsigCountry      string `json:"ConsigCountry"     validate:"required,alpha"`
	ConsigContactName  string `json:"ConsigContactName" validate:"omitempty"`
	ConsigEmail        string `json:"ConsigEmail"       validate:"omitempty"`
	ConsigPhoneNumber  string `json:"ConsigPhoneNumber" validate:"omitempty"`
	ConsigFax          string `json:"ConsigFax"         validate:"omitempty"`
	ConsigLoadNotes    string `json:"ConsigLoadNotes"   validate:"omitempty"`
	ConsigEarliestDate string `json:"ConsigEarliestDate" validate:"required"`
	ConsigEarliestTime string `json:"ConsigEarliestTime" validate:"omitempty"`
	ConsigLatestDate   string `json:"ConsigLatestDate"   validate:"omitempty"`
	ConsigLatestTime   string `json:"ConsigLatestTime"   validate:"omitempty"`
}
