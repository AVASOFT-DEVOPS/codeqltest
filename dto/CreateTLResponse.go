package dto


type TLresponse struct {
	Code    int    `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	TLId    string`json:"TLId,omitempty"`
	AgentDetails []CustomerDetailsDto `json:"AgentDetails,omitempty"`
}

func (TL TLresponse) AsUserRespopnse() *TLresponse {
   return &TLresponse{
      TLId: TL.TLId,
   }
}

type VerifyGCIDresponse struct {
	Permissioncheck bool `json:Permissioncheck`

}

type CustomerDetailsDto struct {	
	Name             string `json:"name"`
	Salesperson      string `json:"salesperson"`
	AccountManager   string `json:"accountmanager"`
	AccountManagerId string `json:"accountmanagerid"`
	CustaId          int    `json:"custaId"`
	CustmId          int    `json:"custmId"`
	SalesId          int    `json:"salesId"`
	AgentPhoneNumber string `json:"agentPhoneNumber"`
	AgentName        string `json:"agentName"`
	AgentEmail       string `json:"agentEmail"`
	OfficeCode       string `json:"officeCode"`
	DaxBalance       string `json:"daxBalance"`
	BtmsBalance      string `json:"btmsBalance"`
	Contact          string `json:"contact"`
	CustAgentPhone   string `json:"cusPhone"`
	CustAgentFax     string `json:"cusFax"`
	CustAgentEmail   string `json:"cusEmail"`
 }
