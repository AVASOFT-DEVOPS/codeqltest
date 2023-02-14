package dto

type BookLTLResponseDTO struct {
	LoadNumber   int    `json:"loadNumber"`
	QuoteId      int    `json:"quoteId" `
	QuoteNumber  string `json:"quoteNumber"`
	AccessToken  string `json:"accessToken"`
	AgentEmail   string `json:"AgentEmail"`
	CustEmail    string `json:"CustEmail"`
	CustFax      string `json:"CustFax"`
	CustPhone    string `json:"CustPhone"`
	Contact      string `json:"Contact"`
	PriceDetails PriceDetails
	TotalPrice   string `json:"totalPrice" `
}

type BookLTLResponseDto struct {
	LoadNumber   int    `json:"loadNumber"`
	QuoteId      int    `json:"quoteId" `
	QuoteNumber  string `json:"quoteNumber"`
	AgentEmail   string `json:"AgentEmail"`
	PriceDetails PriceDetails
	TotalPrice   string `json:"totalPrice" `
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

//bookLTLResponseDTO{
//loadNumber number
//quoteNumber string
//priceDetails {
//scac string
//service string
//carrierName string
//carrierNotes string
//transitTime number
//flatPrice number
//fuelSurchargePrice number
//otherCharges:[{
//name: string,
//code: string,
//price: string
//}]
//accessorialsPrice [{
//name string
//code string
//price number
//}]
//totalPrice number
//}
//billOfLadingURL string
//}
//}
