// package Service

// import (
// 	"golang/domain"
// 	"golang/errs"
// )

// type CreateService interface {
// 	GetCreateDetail(int) (int, *errs.AppErrorvalidation)
// }

// type DefaultCreateService struct {
// 	repo domain.CreateRepository
// }

// func (r DefaultCreateService) GetCreateDetail(i int) (int, *errs.AppErrorvalidation) {
// 	return 0, nil
// }

// func NewCreateService(repository domain.CreateRepository) DefaultCreateService {
// 	return DefaultCreateService{repository}
// }

//ps_1.4.1 package service is initialized
package Service

import (
	"encoding/json"
	"fmt"
	"golang/constant"
	"golang/domain"
	"golang/dto"
	"golang/errs"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	// "net/http"
	// "regexp"
	// "strconv"
	// "time"
)

//type Service interface {
//	GetAllstudent(string) ([]dto.StudentResponse, *errs.AppError)
//}

type DefaultTLService struct {
	repo domain.TLRepository
}
//ps_1.4.2 Interface for the service is created
//go:generate mockgen -destination=../mock/service/mockCtLtlService.go -package=Service golang/service TLService
type TLService interface {
	CreateTLService(TlRequest dto.CreateTLReq) (*dto.TLresponse, *errs.CUserError)
}

//ps_1.4.3 GetUserSummaryService is defined where the repo files are called
// The above code is to get the user summary details.

//ps_1.4.3 createTl service function is  defined here to create a truckLoad
func (s DefaultTLService) CreateTLService(TlRequest dto.CreateTLReq) (*dto.TLresponse, *errs.CUserError){
   userIdRegex :=regexp.MustCompile("[0-9]+")
   var errorarray []dto.ValidateResDTO

	if TlRequest.UserId == ""{
		errorarray=append(errorarray,dto.ValidateResDTO{
			Code:"LT001" ,
			Message:"Missing parameter,UserId",
		})
	}
	if TlRequest.UserId != "" && !userIdRegex.MatchString(TlRequest.UserId){
		errorarray=append(errorarray,dto.ValidateResDTO{
			Code:"LT001" ,
			Message:"Invalid parameter,UserId",
		})
	}
	if TlRequest.CustomerLoginId ==""{
		errorarray=append(errorarray,dto.ValidateResDTO{
			Code:"LT001" ,
			Message:"Missing parameter,CustomerLoginId",
		})
	}
	if TlRequest.CustomerLoginId !="" && !userIdRegex.MatchString(TlRequest.CustomerLoginId){
		errorarray=append(errorarray,dto.ValidateResDTO{
			Code:"LT001" ,
			Message:"Invalid parameter,CustomerLoginId",
		})
	}
	if len(errorarray)!=0{
		return nil,errs.UserRequestValidation(errorarray)
	}
	var verifyReq=dto.VerifyGCIDRequest{
		CustomerLoginId:TlRequest.CustomerLoginId,
		UserId: TlRequest.UserId,
	}

	fmt.Println(verifyReq,"verifyReq\n")
	//ps_1.4.4 validation functionis called where the permission and GCID is checked
	verifyGcidres,error :=s.repo.VerifyGCIDPer(verifyReq);
	fmt.Println(error,"error\n")
	if error !=nil{
      return nil,error
	}
	fmt.Println(verifyGcidres.Permissioncheck,"verifyGcidres.GCIDpermissioncheck\n")
	//ps_1.8.1 The permssion returned from the repo is checked
	if !verifyGcidres.Permissioncheck{
		//ps_1.9.1 if the permission value is false then error response is returned
		return nil,errs.UserAuthentication("Access Denied – You don’t have permission to access.")
	}
	//request Validation
	//ps_1.13.1,ps_1.14.1 If User has permission the request is validated
	//ps_1.15.1 ValidateTLrequest is called to validate the response
	ValidateReqres:=ValidateTLRequest(TlRequest);
	if len(ValidateReqres) != 0 {
		//ps_1.16.1 length of the errorArray returned is checked
		for _,errA:=range ValidateReqres{
			errorarray=append(errorarray,errA)}
		}
		if len(errorarray)!=0{
		return nil, errs.UserRequestValidation(errorarray)
		}

		//ps_1.19.1 validations passes the createTL repo file is called
	TLres,TLerror:=s.repo.CreateTLRepo(TlRequest)
	//ps_1.20.2 The response returned from the repo is checked for the error 
	if TLerror !=nil{
		fmt.Println(TLerror,"Entered TL error in service")
		return nil,TLerror
	}
	loadID,_:=strconv.Atoi(TLres.TLId)
	 var totalweight int
	 for _,Com:= range TlRequest.CommoditiesInfo{
		totalweight +=Com.Weight
	 }

	  MailReq := dto.EmailRequestdto{
		LoadId:           loadID,
		EarliestDate:     TlRequest.ShipperInfo.ShipEarliestDate,
		LatestDate:       TlRequest.ConsigneeInfo.ConsigEarliestDate,
		ShipperAddres:    TlRequest.ShipperInfo.ShipAdd1 + " , " + "  " + TlRequest.ShipperInfo.ShipCity + " , " + "  " + TlRequest.ShipperInfo.ShipState + " , " + "  " + TlRequest.ShipperInfo.ShipZipcode,
		ConsigneeAddress: TlRequest.ConsigneeInfo.ConsigAdd1 + "  , " + "  " + TlRequest.ConsigneeInfo.ConsigCity + " , " + "  " + TlRequest.ConsigneeInfo.ConsigState + " , " + "  " + TlRequest.ConsigneeInfo.ConsigZipcode,
		Commodity:        TlRequest.CommoditiesInfo[0].Descrip,
		LTL:              "",
		Weight:           totalweight,
		Comments:         "",
		Contact:          TLres.AgentDetails[0].Contact,
		Phone:            TLres.AgentDetails[0].CustAgentPhone,
		Fax:              TLres.AgentDetails[0].CustAgentFax,
		Email:            TLres.AgentDetails[0].CustAgentEmail,
		Sendermail:       "kirubasri.s@avasoft.com",
	 }
	 fmt.Println("MailReq",MailReq)

	 mailerr:=SendMail(MailReq,TlRequest.Token)
	 fmt.Println("mailerr",mailerr)
	 if mailerr!=nil{
		return nil,mailerr
	 }
	//ps_1.22.1 If no error occurs the responseDto is returned
	return TLres,nil
}

func SendMail(MailReq dto.EmailRequestdto,token string) (*errs.CUserError) {

	url := constant.BASE_URL + "/postmark/bookmail"
    method := "POST"

	marshalstring, _ := json.Marshal(MailReq)
    payload := strings.NewReader(string(marshalstring))
    fmt.Println()
    client := &http.Client{}

    req, err1 := http.NewRequest(method, url, payload)



    if err1 != nil {
        fmt.Println(err1, "schch")
        return errs.UserNewUnexpectedError("")
    }

    req.Header.Add("Authorization", "Bearer "+token)
    req.Header.Add("Content-Type", "application/json")

    res, err2 := client.Do(req)
    if err2 != nil {
        fmt.Println(err2, "dhdhj")
        return errs.UserNewUnexpectedError("")
    }

    defer res.Body.Close()

    body, err3 := ioutil.ReadAll(res.Body)
    if err3 != nil {
        fmt.Println(err3, "schchj")
        return  errs.UserNewUnexpectedError("")
    }

    fmt.Println(string(body))

    var PostmarkResponse dto.PostmarkResponse
    json.Unmarshal(body, &PostmarkResponse)
    fmt.Println(PostmarkResponse, "ddjsk")
    log.Println("ending of postmark", time.Now())
	return nil
}

func ValidateTLRequest(loadRequest dto.CreateTLReq) ([]dto.ValidateResDTO) {
	// GCID            string            `json:"GCID" validate:"required,alphanum"`
	// CustomerId        string          `json:"Customer" validate:"required,numeric"`
	// CustomerLoginId string            `json:"CustomerLoginId" validate:"required,numeric"`
	// UserId          string            `json:"UserId" validate:"required,numeric"`
	// ShipmentInfo    ShipmentInfo      `json:"ShipmentInfo,omitempty"`
	// ShipperInfo     ShipperInfo       `json:"ShipperInfo,omitempty"`
	// ConsigneeInfo   ConsigneeInfo     `json:"ConsigneeInfo,omitempty"`
	// CommoditiesInfo []CommoditiesInfo `json:"CommoditiesInfo,omitempty"`

	fmt.Println("Enterred the validateRequest for createTL\n")

	validate := validator.New()

	errCon:= validate.Struct(loadRequest.ConsigneeInfo)
	fmt.Println(errCon,"errCon\n")
	errShip:=validate.Struct(loadRequest.ShipperInfo)
	fmt.Println(errShip,"errShip\n")
	errShipment:=validate.Struct(loadRequest.ShipmentInfo)
	fmt.Println(errShipment,"errShipment\n")
	errComm:=validate.Struct(loadRequest.CommoditiesInfo)
	fmt.Println(errComm,"errComm\n")
	Datere := regexp.MustCompile(`[0-9]+-(0?[1-9]|[1][0-2])-(0[1-9]|[12][0-9]|3[01])`)
	Timere:=regexp.MustCompile(`(0[0-9]|1[0-9]|2[0-3]):(0[0-9]|[1-5][0-9])`)
	textre:=regexp.MustCompile(`^[a-zA-Z ]+$`)
	PhFaxre:=regexp.MustCompile(`^[0-9]{10}$`)
	PhFaxare1:=regexp.MustCompile(`^[0-9]{3}-[0-9]{3}-[0-9]{4}$`)
	PhFaxare2:=regexp.MustCompile(`^[(][0-9]{3}[)]\s[0-9]{3}-[0-9]{4}$`)
	ZipCodeCanre:=regexp.MustCompile(`^[ABCEGHJ-NPRSTVXYabceghj-nprstvxy]\d[ABCEGHJ-NPRSTV-Zabceghj-nprstv-z][ ]?\d[ABCEGHJ-NPRSTV-Zabceghj-nprstv-z]\d$`)
	ZipCodeAmere:=regexp.MustCompile(`^[0-9]{5}$`)
	layout := "2006-01-02"
	tempnowdate:=time.Now().Format("2017-09-07")
	nowdate,_:=time.Parse(layout, tempnowdate)


	errorArray := make([]dto.ValidateResDTO, 0)
	
	//Consignee DTO Validation
	// ConsigName         string `json:"ConsigName" validate:"Required,alpha"`
	// ConsigAdd1         string `json:"ConsigAdd1" validate:"Required,alphanum,"`
	// ConsigAdd2         string `json:"ConsigAdd2" validate:"omitempty"`
	// ConsigCity         string `json:"ConsigCity" validate:"Required,alpha"`
	// ConsigState        string `json:"ConsigState" validate:"Required,alpha"`
	// ConsigZipcode      int `json:"ConsigZipcode" validate:"Required,numeric"`
	// ConsigCountry      string `json:"ConsigCountry" validate:"Required,alpha"`
	// ConsigContactName  string `json:"ConsigContactName" validate:"omitempty"`
	// ConsigEmail        string `json:"ConsigEmail" validate:"omitempty"`
	// ConsigPhoneNumber  int `json:"ConsigPhoneNumber" validate:"omitempty"`
	// ConsigFax          int `json:"ConsigFax" validate:"omitempty"`
	// ConsigLoadNotes    string `json:"ConsigLoadNotes" validate:"omitempty"`
	// ConsigEarliestDate string `json:"ConsigEarliestDate" validate:"Required"`
	// ConsigEarliestTime string `json:"ConsigEarliestTime" validate:"omitempty"`
	// ConsigLatestDate   string `json:"ConsigLatestDate" validate:"omitempty"`
	// ConsigLatestTime   string `json:"ConsigLatestTime" validate:"omitempty"`
	ConsigEarliestDate, _ := time.Parse(layout, loadRequest.ConsigneeInfo.ConsigEarliestDate)
	ConsigLatestDate, _ := time.Parse(layout, loadRequest.ConsigneeInfo.ConsigLatestDate)
	ShipEarliestDate, _ := time.Parse(layout, loadRequest.ShipperInfo.ShipEarliestDate)
	ShipLatestDate, _ := time.Parse(layout, loadRequest.ShipperInfo.ShipLatestDate)
	

	if ConsigEarliestDate!=ShipEarliestDate && !ConsigEarliestDate.After(ShipEarliestDate) {
		errorArray=append(errorArray,dto.ValidateResDTO{
			Code:"TL001",
			Message:"Invalid Parameter,Consignee EarliestDate should be Greater than the Shipper EarliestDate",
		})
	}
	if !textre.MatchString(loadRequest.ConsigneeInfo.ConsigCity){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee City"})
	}
	if !textre.MatchString(loadRequest.ConsigneeInfo.ConsigCountry){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Country"})
	}
	if len(loadRequest.ConsigneeInfo.ConsigName) > 255{
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Name"})
	}


	if loadRequest.ConsigneeInfo.ConsigPhoneNumber!="" && len(loadRequest.ConsigneeInfo.ConsigPhoneNumber) <=10 && !PhFaxre.MatchString(loadRequest.ConsigneeInfo.ConsigPhoneNumber){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Phone Number"})
	}else if len(loadRequest.ConsigneeInfo.ConsigPhoneNumber) > 10 &&  len(loadRequest.ConsigneeInfo.ConsigPhoneNumber) < 12 && !PhFaxare1.MatchString(loadRequest.ConsigneeInfo.ConsigPhoneNumber) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Phone Number"})
	}else if len(loadRequest.ConsigneeInfo.ConsigPhoneNumber) > 12 && !PhFaxare2.MatchString(loadRequest.ConsigneeInfo.ConsigPhoneNumber)  {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Phone Number"})
	}


	if loadRequest.ConsigneeInfo.ConsigFax!="" && len(loadRequest.ConsigneeInfo.ConsigFax) <= 10 && !PhFaxre.MatchString(loadRequest.ConsigneeInfo.ConsigFax){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Fax Number"})
	}else if len(loadRequest.ConsigneeInfo.ConsigFax) > 10 &&  len(loadRequest.ConsigneeInfo.ConsigFax) < 12 && !PhFaxare1.MatchString(loadRequest.ShipperInfo.ShipFax) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Fax Number"})
	}else if len(loadRequest.ConsigneeInfo.ConsigFax) > 12 && !PhFaxare2.MatchString(loadRequest.ConsigneeInfo.ConsigFax)  {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Fax Number"})
	}



	if errCon != nil {
		if strings.Contains(errCon.Error(), "'ConsigName' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee Name"})
		}
		if strings.Contains(errCon.Error(), "'ConsigAdd1' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee Address"})
		}
		if strings.Contains(errCon.Error(), "'ConsigCity' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee City"})
		}
		if strings.Contains(errCon.Error(), "'ConsigState' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee State"})
		}
		if strings.Contains(errCon.Error(), "'ConsigState' failed on the 'alpha' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee State"})
		}
		if strings.Contains(errCon.Error(), "'ConsigZipcode' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee Zipcode"})
		}
		if strings.Contains(errCon.Error(), "'ConsigCountry' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee Country"})
		}
		// if strings.Contains(errCon.Error(), "'ConsigCountry' failed on the 'alpha' tag") {
		// 	errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Country"})
		// }
		if strings.Contains(errCon.Error(), "'ConsigEarliestDate' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Consignee EarliestDate"})
		}
	}

	if  len(loadRequest.ConsigneeInfo.ConsigZipcode)== 5 && !ZipCodeAmere.MatchString(loadRequest.ConsigneeInfo.ConsigZipcode) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Zipcode"})
	}else if len(loadRequest.ConsigneeInfo.ConsigZipcode)== 7 && !ZipCodeCanre.MatchString(loadRequest.ConsigneeInfo.ConsigZipcode){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Consignee Zipcode"})
	}


	if !Datere.MatchString(loadRequest.ConsigneeInfo.ConsigEarliestDate) {
		//yyyy-MM-DD
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Consignee EarliestDate"})
	}
	if loadRequest.ConsigneeInfo.ConsigEarliestTime!="" && !Timere.MatchString(loadRequest.ConsigneeInfo.ConsigEarliestTime) {
		//24hr
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Consignee EarliestTime"})
	}

	if ConsigEarliestDate.Before(nowdate) {
		errorArray=append(errorArray,dto.ValidateResDTO{
			Code:"TL001",
			Message:"Invalid Parameter,Consignee EarliestDate cannot have past dates",
		})
	}

	if loadRequest.ConsigneeInfo.ConsigLatestDate !="" && !Datere.MatchString(loadRequest.ConsigneeInfo.ConsigLatestDate) {
		//yyyy-MM-DD
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Consignee LatestDate"})
	}
	if loadRequest.ConsigneeInfo.ConsigLatestDate !="" && ConsigLatestDate.Before(nowdate) {
		errorArray=append(errorArray,dto.ValidateResDTO{
			Code:"TL001",
			Message:"Invalid Parameter,Consignee LatestDate",
		})
	}
	if loadRequest.ConsigneeInfo.ConsigLatestTime !="" && !Timere.MatchString(loadRequest.ConsigneeInfo.ConsigLatestTime) {
		//24hr
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Consignee LatestTime"})
	}
	//Consignee End
	//Shipper DTO Validation
	// ShipName           string `json:"ShipName" validate:"Required,alpha"`
	// ShipAdd1           string `json:"ShipAdd1" validate:"Required,alphanum"`
	// ShipAdd2           string `json:"ShipAdd2" validate:"omitempty"`
	// ShipCity           string `json:"ShipCity" validate:"Required,alpha"`
	// ShipState          string `json:"ShipState" validate:"Required,alpha"`
	// Zipcode            int `json:"ShipZipcode" validate:"Required,numeric"`
	// Country            string `json:"ShipCountry" validate:"Required,alpha"`
	// ShipContactName    string `json:"ShipContactName" validate:"omitempty"`
	// ShipEmail          string `json:"ShipEmail" validate:"omitempty"`
	// ShipPhoneNumber    int `json:"ShipPhoneNumber" validate:"omitempty"`
	// ShipFax            int `json:"ShipFax" validate:"omitempty"`
	// ShipLoadNotes      string `json:"ShipLoadNotes" validate:"omitempty"`
	// ShipEarliestDate   string `json:"ShipEarliestDate" validate:"Required"`
	// ShipEarliestTime   string `json:"ShipEarliestTime" validate:"omitempty"`
	// ShipLatestDate     string `json:"ShipLatestDate" validate:"omitempty"`
	// ShipLatestTime     string `json:"ShipLatestTime" validate:"omitempty"`

	if !textre.MatchString(loadRequest.ShipperInfo.ShipCity){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper City"})
	}
	if !textre.MatchString(loadRequest.ShipperInfo.Country){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Country"})
	}
	if len(loadRequest.ShipperInfo.ShipName) > 255{
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Name"})
	}


	if loadRequest.ShipperInfo.ShipPhoneNumber!="" && len(loadRequest.ShipperInfo.ShipPhoneNumber) <=10 && !PhFaxre.MatchString(loadRequest.ShipperInfo.ShipPhoneNumber){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Phone Number"})
	}else if len(loadRequest.ShipperInfo.ShipPhoneNumber) > 10 &&  len(loadRequest.ShipperInfo.ShipPhoneNumber) < 12 && !PhFaxare1.MatchString(loadRequest.ShipperInfo.ShipPhoneNumber) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Phone Number"})
	}else if len(loadRequest.ShipperInfo.ShipPhoneNumber) > 12 && !PhFaxare2.MatchString(loadRequest.ShipperInfo.ShipPhoneNumber)  {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Phone Number"})
	}


	if loadRequest.ShipperInfo.ShipFax!="" && len(loadRequest.ShipperInfo.ShipFax) <=10 && !PhFaxre.MatchString(loadRequest.ShipperInfo.ShipFax){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Fax Number"})
	}else if len(loadRequest.ShipperInfo.ShipFax) > 10 &&  len(loadRequest.ShipperInfo.ShipFax) < 12 && !PhFaxare1.MatchString(loadRequest.ShipperInfo.ShipFax) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Fax Number"})
	}else if len(loadRequest.ShipperInfo.ShipFax) > 12 && !PhFaxare2.MatchString(loadRequest.ShipperInfo.ShipFax)  {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Fax Number"})
	}


	if errShip != nil {
		if strings.Contains(errShip.Error(), "'ShipName' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper Name"})
		}

		if strings.Contains(errShip.Error(), "'ShipAdd1' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper Address"})
		}

		if strings.Contains(errShip.Error(), "'ShipCity' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper City"})
		}
		// if strings.Contains(errShip.Error(), "'ShipCity' failed on the 'alpha' tag") {
		// 	errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper City"})
		// }

		if strings.Contains(errShip.Error(), "'ShipState' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper State"})
		}
		if strings.Contains(errShip.Error(), "'ShipState' failed on the 'alpha' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper State"})
		}

		if strings.Contains(errShip.Error(), "'ShipZipcode' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper Zipcode"})
		}
		//regex validation to be added

		if strings.Contains(errShip.Error(), "'ShipCountry' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper Country"})
		}
		// if strings.Contains(errShip.Error(), "'ShipCountry' failed on the 'alpha' tag") {
		// 	errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Country"})
		// }
		if strings.Contains(errShip.Error(), "'ShipEarliestDate' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Shipper EarliestDate"})
		}
		
	}

	
	if  len(loadRequest.ShipperInfo.ShipZipcode)== 5 && !ZipCodeAmere.MatchString(loadRequest.ShipperInfo.ShipZipcode) {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Zipcode"})
	}else if len(loadRequest.ShipperInfo.ShipZipcode)== 7 && !ZipCodeCanre.MatchString(loadRequest.ShipperInfo.ShipZipcode){
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Shipper Zipcode"})
	}

	if !Datere.MatchString(loadRequest.ShipperInfo.ShipEarliestDate) {
		//yyyy-MM-DD
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Shipper EarliestDate"})
	}
	if ShipEarliestDate.Before(nowdate) {
		errorArray=append(errorArray,dto.ValidateResDTO{
			Code:"TL001",
			Message:"Invalid Parameter,Shipper EarliestDate cannot have past dates",
		})
	}

	if  loadRequest.ShipperInfo.ShipEarliestTime!="" && !Timere.MatchString(loadRequest.ShipperInfo.ShipEarliestTime) {
		//24hr
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Shipper EarliestTime"})
	}

	if loadRequest.ShipperInfo.ShipLatestDate !="" && !Datere.MatchString(loadRequest.ConsigneeInfo.ConsigLatestDate) {
		//yyyy-MM-DD
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Shipper LatestDate"})
	}
	if loadRequest.ShipperInfo.ShipLatestDate !="" && ShipLatestDate.Before(nowdate) {
		errorArray=append(errorArray,dto.ValidateResDTO{
			Code:"TL001",
			Message:"Invalid Parameter,Shipper LatestDate",
		})
	}
	if loadRequest.ShipperInfo.ShipLatestTime !="" && !Timere.MatchString(loadRequest.ShipperInfo.ShipLatestTime) {
		//24hr
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid Parameter,Shipper LatestTime"})
	}
	//Shipment DTO Validation
	// EquipType          string `json:"EquipType" validate:"Required,alpha"`
	// PONumber           int `json:"PONumber,omitempty" validate:"omitempty"`
	// BLNumber           int `json:"BLNumber,omitempty" validate:"omitempty"`
	// ShippingNumber     int `json:"ShippingNumber,omitempty" validate:"omitempty"`
	// ReferrenceNumber   int `json:"ReferrenceNumber,omitempty" validate:"omitempty"`
    if errShipment != nil {
		if strings.Contains(errShipment.Error(), "'EquipType' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, EquipType"})
		}
		if strings.Contains(errShipment.Error(), "'EquipType' failed on the 'alpha' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, EquipType"})
		}
	}
	//Commodity DTO Validation
	// BolsplIns          string `json:"BolsplIns" validate:"omitempty"`
	// Qty                int `json:"Qty" validate:"Required,numeric"`
	// UOM                string `json:"UOM" validate:"Requiredomitempty"`
	// Weight             int `json:"Weight" validate:"Required,numeric"`
	// Value              int `json:"Value" validate:"omitempty"`
	// Descrip            string `json:"Descrip" validate:"Required,alpha"`
	// Hazmat             bool  `json:"Hazmat"`
	// HazmatInfo         HazmatInfo `json:"hazmatInfo"`
	for comindex := 0; comindex < len(loadRequest.CommoditiesInfo); comindex++ {
		err2 := validate.Struct(loadRequest.CommoditiesInfo[comindex])
		fmt.Println("err2", err2)
		errComm := err2
	if errComm != nil {
		if strings.Contains(errComm.Error(), "'Qty' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Qty for commodity "})
		}
		if strings.Contains(errComm.Error(), "'Qty' failed on the 'numeric' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Qty"})
		}
		if strings.Contains(errComm.Error(), "'Weight' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, Weight"})
		}
		if strings.Contains(errComm.Error(), "'Weight' failed on the 'numeric' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, Weight"})
		}
		if strings.Contains(errComm.Error(), "'UOM' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, UOM"})
		}
		if strings.Contains(errComm.Error(), "'UOM' failed on the 'alpha' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter, UOM"})
		}
		break;
	}}

	return errorArray
}


// If the element hasn't occurred yet, add it to the result and mark it as occurred.
func unique(arr []int) []int {
	occurred := map[int]bool{}
	result := []int{}
	for e := range arr {
		if occurred[arr[e]] != true {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

// It returns true if the string str is in the slice s, and false otherwise
func isElementExist(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func NewTLService(repository domain.TLRepository) DefaultTLService {
	return DefaultTLService{repository}
}


// func (s DefaultService) GetAllstudent(status string) ([]dto.StudentResponse, *errs.AppError) {
// 	fmt.Println("entered service")

// 	Student, err := s.repo.FindAll(status)
// 	if err != nil {
// 		return nil, err
// 	}
// 	response := make([]dto.StudentResponse, 0)
// 	for _, c := range Student {
// 		response = append(response, c.ToSDto())
// 	}
// 	return response, err

// }
