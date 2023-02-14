//ps_1.5.1 Domain package is initialized and the interfaces used in the repo is defined here
package domain

import (
	"golang/dto"
	"golang/errs"
)

// type TLRepository interface {
// 	GetUser(id string) ([]GetUserres, *errs.AppError)
// 	GetUserSummaryRepo(userReq dto.UserReq) ([]UserSummaryRes, *errs.AppError)
// }

type TLRepository interface {
	VerifyGCIDPer(dto.VerifyGCIDRequest)(*dto.VerifyGCIDresponse,*errs.CUserError)
	CreateTLRepo(TlRequest dto.CreateTLReq)(*dto.TLresponse,*errs.CUserError)
}



type VerifyUserPerRes struct {
	GCID string `db:"global_customer_id"`
}

type VerifyLoginCustRes struct {
	GCID string `db:"user_id"`
}

type verifyAddressRes struct{
	Id int `db:"id"`
}


type CustomerDetails struct {	
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

 type  loashshinsert struct{
	Loashid int `db:"id"`
 }

 type lscarrinsert struct{
	Lscarrid int `db:"id"`
 }

 type lsstopsinsert struct{
	Lsstopsid int `db:"id"`
 }
 func (s CustomerDetails) ToCustomerDetailsDto()(dto.CustomerDetailsDto){
	return dto.CustomerDetailsDto{
		Name                :s.Name,
        Salesperson     :s.Salesperson,
        AccountManager  :s.AccountManager,
        AccountManagerId:s.AccountManagerId,
        CustaId         :s.CustaId,
        CustmId         :s.CustmId,
        SalesId         :s.SalesId,
        AgentPhoneNumber:s.AgentPhoneNumber,
        AgentName       :s.AgentName,
        AgentEmail      :s.AgentEmail,
        OfficeCode      :s.OfficeCode,
        DaxBalance      :s.DaxBalance,
        BtmsBalance     :s.BtmsBalance,
        Contact         :s.Contact,
        CustAgentPhone  :s.AgentPhoneNumber,
        CustAgentFax    :s.CustAgentFax,
        CustAgentEmail  :s.CustAgentEmail,
	}
 } 


