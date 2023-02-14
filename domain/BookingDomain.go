package domain

import (
	"golang/dto"
	"golang/errs"
)

type BookRepository interface {
	VerifyAuthDataRepo(authDataReq *dto.AuthDataReq) (*VerifyAuthRes, *errs.AppErrorvalidation)
	GetCustDataRepo(loginId string) (*CustDataRes, *errs.AppError)
	CheckRatingEngine(CustmId string) (*CheckRatingEngineResponse, *errs.AppError)
	QuoteDetails(loadId int) (*QuoteDetailsResponse, *errs.AppError)
	InsertBookDetailsRepo(banres *dto.BanyanGetLoadDetailsDto, bookreq *dto.BookLTLRequestDTO, SelectedAccessorials string) (*dto.BookLTLResponseDTO, *errs.AppError)
	InsertBookDetailsModeRepo(modeRes *dto.QuoteDetailsDto, bookreq *dto.BookLTLRequestDTO) (*dto.BookLTLResponseDTO, *errs.AppError)
}

type CheckRatingEngineResponse struct {
	RatingEngine string `db:"rating_engine"`
}

type QuoteDetailsResponse struct {
	XmlTraffic string `db:"xml_traffic"`
}

type CustmerDetails struct {
	Name             string `db:"name"`
	Salesperson      string `db:"salesperson"`
	AccountManager   string `db:"accountmanager"`
	AccountManagerId string `db:"accountmanagerid"`
	CustaId          int    `db:"custaId"`
	CustmId          int    `db:"custmId"`
	SalesId          int    `db:"salesId"`
	AgentPhoneNumber string `db:"agentPhoneNumber"`
	AgentName        string `db:"agentName"`
	AgentEmail       string `db:"agentEmail"`
	OfficeCode       string `db:"officeCode"`
	DaxBalance       string `db:"daxBalance"`
	BtmsBalance      string `db:"btmsBalance"`
	Contact          string `db:"contact"`
	CustAgentPhone   string `db:"cusPhone"`
	CustAgentFax     string `db:"cusFax"`
	CustAgentEmail   string `db:"cusEmail"`
}

type CarrIdDetails struct {
	Id              int    `db:"id"`
	Name            string `db:"name"`
	Notes           string `db:"notes"`
	PalletLength    int    `db:"pallet_length"`
	PalletWidth     int    `db:"pallet_width"`
	PalletHeight    int    `db:"pallet_height"`
	PalletWeight    int    `db:"pallet_weight"`
	LtlFlag         int    `db:"ltl_flag"`
	BanSupports204  int    `db:"ban_supports204"`
	CarrActivated   int    `db:"carr_activated"`
	CarrId          int    `db:"carr_id"`
	CarrUsable      int    `db:"carr_usable"`
	CarrPref        int    `db:"carr_pref"`
	OfficeProhibted int    `db:"office_prohibited"`
}

type LoadShIdStruct struct {
	LoadshId   int `db:"id"`
	LtlQuoteId int `db:"ltl_quote_id"`
	CarrId     int `db:"carr_id"`
}

type LscarrIdStruct struct {
	LscarrId int `db:"id"`
}
type LsstopIdStruct struct {
	LsstopId int `db:"id"`
}

type BookLTLResponseDTO struct {
	LoadNumber   int `json:"loadNumber"`
	QuoteId      int `json:"quoteId" `
	QuoteNumber  int `json:"quoteNumber"`
	PriceDetails PriceDetails
	TotalPrice   int `json:"totalPrice" `
}

type PriceDetails struct {
	Scac               string `json:"scac"`
	Service            string `json:"service" `
	CarrierName        string `json:"carrierName" `
	CarrierNotes       string `json:"carrierNotes" `
	TransitTime        string `json:"transitTime" `
	FlatPrice          string `json:"flatPrice" `
	FuelSurchargePrice string `json:"fuelSurchargePrice" `
}

type VerifyAuthRes struct {
	UserId             string `db:"userId"`
	LoginId            string `db:"sunteckLoginId"`
	CustmId            string `db:"custmId"`
	CustomerPermission string `db:"customerPermission"`
}

type CustDataRes struct {
	CustmId      string `db:"custm_id"`
	CustaId      string `db:"custa_id"`
	OfficeCode   string `db:"office_code"`
	RatingEngine string `db:"rating_engine"`
}

func (s VerifyAuthRes) ToSDto() dto.VerifyAuthRes {

	return dto.VerifyAuthRes{
		UserId:             s.UserId,
		LoginId:            s.LoginId,
		CustmId:            s.CustmId,
		CustomerPermission: s.CustomerPermission,
	}

}

func (s CustDataRes) ToSDto() dto.CustDataRes {
	return dto.CustDataRes{
		CustmId:      s.CustmId,
		CustaId:      s.CustaId,
		OfficeCode:   s.OfficeCode,
		RatingEngine: s.RatingEngine,
	}
}
