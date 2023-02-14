package dto

import "time"

type BookLTLRequestDTO struct {
	UserId           string
	SunteckLoginId   string
	CustmId          string
	CustaId          string
	LoadId           int    `json:"loadId" validate:"required"`
	QuoteId          int    `json:"quoteId" validate:"required"`
	PoNumber         string `json:"poNumber"`
	BlNumber         string `json:"blNumber"`
	ShippingNumber   string `json:"shippingNumber" `
	ShipperDetails   ShipperAddress
	ConsigneeDetails ConsigneeAddress
	Commodities      Commodities `json:"Commodities" validate:"omitempty,numeric"`
}

type Commodities struct {
	Commodities []Commodity `json:"Commodities"`
	LinearFt    string      `json:"LinearFt" validate:"required,alpha"`
}

type ShipperAddress struct {
	ShipName         string `json:"ShipName" validate:"required,alpha"`
	ShipAddress1     string `json:"ShipAddress1" validate:"required,alpha"`
	ShipAddress2     string `json:"ShipAddress2" validate:"required,alpha"`
	ShipCity         string `json:"ShipCity" validate:"required,alpha"`
	ShipState        string `json:"ShipState" validate:"required,alpha"`
	ShipZipCode      string `json:"ShipZipCode" validate:"required,alpha"`
	ShipCountry      string `json:"ShipCountry" validate:"required,alpha"`
	ShipContactName  string `json:"ShipContactName" validate:"omitempty,alpha"`
	ShipEmail        string `json:"ShipEmail" validate:"omitempty,alpha"`
	ShipPhone        string `json:"ShipPhone" validate:"omitempty,numeric"`
	ShipFax          string `json:"ShipFax" validate:"omitempty,numeric"`
	ShipLoadNotes    string `json:"ShipLoadNotes" validate:"omitempty,alpha"`
	ShipEarliestDate string `json:"ShipEarliestDate" validate:"required,alpha"`
	ShipLatestDate   string `json:"ShipLatestDate" validate:"omitempty,alpha"`
	ShipEarliestTime string `json:"ShipEarliestTime" validate:"omitempty,alpha"`
	ShipLatestTime   string `json:"ShipLatestTime" validate:"omitempty,alpha"`
}
type ConsigneeAddress struct {
	ConsName         string `json:"ConsName" validate:"required,alpha"`
	ConsAddress1     string `json:"ConsAddress1" validate:"required,alpha"`
	ConsAddress2     string `json:"ConsAddress2" validate:"omitempty,alpha"`
	ConsCity         string `json:"ConsCity" validate:"required,alpha"`
	ConsState        string `json:"ConsState" validate:"required,alpha"`
	ConsZipCode      string `json:"ConsZipCode" validate:"required,alpha"`
	ConsCountry      string `json:"ConsCountry" validate:"required,alpha"`
	ConsContactName  string `json:"ConsContactName" validate:"omitempty,alpha"`
	ConsEmail        string `json:"ConsEmail" validate:"omitempty,alpha"`
	ConsPhone        string `json:"ConsPhone" validate:"omitempty,numeric"`
	ConsFax          string `json:"ConsFax" validate:"omitempty,numeric"`
	ConsLoadNotes    string `json:"ConsLoadNotes" validate:"omitempty,alpha"`
	ConsEarliestDate string `json:"ConsEarliestDate" validate:"required,alpha"`
	ConsLatestDate   string `json:"ConsLatestDate" validate:"omitempty,alpha"`
	ConsEarliestTime string `json:"ConsEarliestTime" validate:"omitempty,alpha"`
	ConsLatestTime   string `json:"ConsLatestTime" validate:"omitempty,alpha"`
}

type Commodity struct {
	Desc          string `json:"Desc" validate:"required,alpha"`
	NMFC          string `json:"NMFC" validate:"required,alpha"`
	Class         string `json:"Class" validate:"required",alpha`
	Stackable     bool   `json:"Stackable" `
	Hazmat        bool   `json:"Hazmat" `
	Length        int    `json:"Length" validate:"required,numeric"`
	Width         int    `json:"Width" validate:"required,numeric"`
	Height        int    `json:"Height" validate:"required,numeric"`
	Weight        int    `json:"Weight" validate:"required,numeric"`
	Quantity      int    `json:"Quantity" validate:"required,numeric"`
	CubicFt       string `json:"CubicFt" validate:"required,numeric"`
	Density       string `json:"density"`
	EquipmentType string `json:"EquipmentType" validate:"required,alpha"`
}

type ValidateDTO struct {
	Code    string `json:"code" validate:"omitempty,alpha"`
	Message string `json:"message" validate:"omitempty,alpha"`
}
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AuthDataReq struct {
	UserId  string
	LoginId string
}

type VerifyAuthRes struct {
	UserId             string
	LoginId            string
	CustmId            string
	CustomerPermission string
}

type CustDataRes struct {
	CustmId      string
	CustaId      string
	OfficeCode   string
	RatingEngine string
}
type EmailRequestdto struct {
	LoadId           int    `json:"LoadId"`
	EarliestDate     string `json:"EarliestDate"`
	LatestDate       string `json:"LatestDate"`
	ShipperAddres    string `json:"ShipperAddres"`
	ConsigneeAddress string `json:"ConsigneeAddress"`
	Commodity        string `json:"Commodity"`
	Equipment        string `json:"Equipment"`
	Length           int    `json:"Length"`
	Rate             string `json:"Rate"`
	LTL              string `json:"LTL"`
	Weight           int    `json:"Weight"`
	Comments         string `json:"Comments"`
	Contact          string `json:"Contact"`
	Phone            string `json:"Phone"`
	Fax              string `json:"Fax"`
	Email            string `json:"Email"`
	Sendermail       string `json:"Sendermail"`
}
type PostmarkResponse struct {
	To          string    `json:"To,omitempty"`
	SubmittedAt time.Time `json:"SubmittedAt,omitempty"`
	MessageID   string    `json:"MessageID,omitempty"`
	ErrorCode   int       `json:"ErrorCode"`
	Message     string    `json:"Message"`
}
