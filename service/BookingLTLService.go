package Service

import (
	"encoding/json"
	"fmt"
	"golang/constant"
	"golang/domain"
	dto "golang/dto"
	"golang/errs"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockgen -destination=../mock/service/mockbLtlService.go -package=Service golang/service BookService
type BookService interface {
	BookLTLService(bookReq *dto.BookLTLRequestDTO) (*dto.BookLTLResponseDto, *errs.AppErrorvalidation)
}

//Connection between service and repository will be done by this struct
type DefaultBookingService struct {
	repo domain.BookRepository
}

//BGLD_ps_1.3-1.5
//This is main service which will trigger all sub function to create a load by validating the given data.
func (r DefaultBookingService) BookLTLService(bookReq *dto.BookLTLRequestDTO) (*dto.BookLTLResponseDto, *errs.AppErrorvalidation) {

	//BGLD_ps_1.5-1.9
	//The given request will be validated by invoking the ValidateBookRequest

	validateResponseVar := ValidateBookRequest(bookReq)

	// Error will be caught from ValidateBookRequest() here...
	if len(validateResponseVar) != 0 {
		return nil, errs.ValidateResponse(validateResponseVar, 400, "")
	}
	fmt.Println(validateResponseVar)

	//UserId and LoginId Struct will be formed to pass the request to VerifyAuthDataRepo()
	verifyAuthReq := dto.AuthDataReq{
		UserId:  bookReq.UserId,
		LoginId: bookReq.SunteckLoginId,
	}
	//BGLD_ps_1.3-1.7
	// This Function will be used to verify the UserId, Login Id and whether they have permission to access bookLtl.

	verifyAuthResp, errVerAuth := r.repo.VerifyAuthDataRepo(&verifyAuthReq)

	//Error will be caught here from VerifyAuthDataRepo()

	if verifyAuthResp == nil {

		return nil, errs.ValidateResponse(nil, 401, "Unauthorised User")

	}
	if errVerAuth != nil {
		return nil, errVerAuth
	}

	//BGLD_ps_1.3-1.10
	// The execute function will be triggered here

	banyanResponse, errBadrequest, errInternal := r.execute(bookReq)

	// Error from execute function will be caught here
	if len(errBadrequest) != 0 {
		return nil, errs.ValidateResponse(errBadrequest, 400, "")
	}
	if errInternal != nil {
		return nil, errs.ValidateResponse(nil, errInternal.Code, errInternal.Message)
	}

	//The Response of the book ltl will be formed here
	banyanRes := dto.BookLTLResponseDto{
		LoadNumber:  banyanResponse.LoadNumber,
		QuoteId:     banyanResponse.QuoteId,
		QuoteNumber: banyanResponse.QuoteNumber,
		AgentEmail:  banyanResponse.AgentEmail,
		PriceDetails: dto.PriceDetails{
			Scac:               banyanResponse.PriceDetails.Scac,
			Service:            banyanResponse.PriceDetails.Service,
			CarrierName:        banyanResponse.PriceDetails.CarrierName,
			CarrierNotes:       banyanResponse.PriceDetails.CarrierNotes,
			TransitTime:        banyanResponse.PriceDetails.TransitTime,
			FlatPrice:          banyanResponse.PriceDetails.FlatPrice,
			FuelSurchargePrice: banyanResponse.PriceDetails.FuelSurchargePrice,
		},
		TotalPrice: banyanResponse.TotalPrice,
	}

	//Postmark External api call for email starts here
	//BGLD_PS_1.11-1.14,1-19-1.22
	url := constant.BASE_URL + "/postmark/bookmail"
	method := "POST"
	var request dto.EmailRequestdto
	totalWeight := 0

	for i, _ := range bookReq.Commodities.Commodities {
		totalWeight += bookReq.Commodities.Commodities[i].Weight
		request = dto.EmailRequestdto{
			LoadId:           banyanResponse.LoadNumber,
			EarliestDate:     bookReq.ShipperDetails.ShipEarliestDate,
			LatestDate:       bookReq.ConsigneeDetails.ConsEarliestDate,
			ShipperAddres:    bookReq.ShipperDetails.ShipAddress1 + " , " + "  " + bookReq.ShipperDetails.ShipCity + " , " + "  " + bookReq.ShipperDetails.ShipState + " , " + "  " + bookReq.ShipperDetails.ShipZipCode,
			ConsigneeAddress: bookReq.ConsigneeDetails.ConsAddress1 + "  , " + "  " + bookReq.ConsigneeDetails.ConsCity + " , " + "  " + bookReq.ConsigneeDetails.ConsState + " , " + "  " + bookReq.ConsigneeDetails.ConsZipCode,
			Commodity:        bookReq.Commodities.Commodities[0].Desc,
			Equipment:        bookReq.Commodities.Commodities[0].EquipmentType,
			Length:           bookReq.Commodities.Commodities[0].Length,
			Rate:             banyanResponse.TotalPrice,
			LTL:              "YES",
			Weight:           totalWeight,
			Comments:         bookReq.ConsigneeDetails.ConsLoadNotes,
			Contact:          banyanResponse.Contact,
			Phone:            banyanResponse.CustPhone,
			Fax:              banyanResponse.CustFax,
			Email:            banyanResponse.CustEmail,
			Sendermail:       "cp-non-prod-service@modetransportation.com",
		}
	}

	marshalstring, _ := json.Marshal(request)
	payload := strings.NewReader(string(marshalstring))
	fmt.Println()
	client := &http.Client{}
	req, err1 := http.NewRequest(method, url, payload)

	if err1 != nil {
		fmt.Println(err1, "schch")
		return nil, nil
	}
	req.Header.Add("Authorization", "Bearer "+banyanResponse.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println(err2, "dhdhj")
		return nil, nil
	}
	defer res.Body.Close()

	body, err3 := ioutil.ReadAll(res.Body)
	if err3 != nil {
		fmt.Println(err3, "schchj")
		return nil, nil
	}
	fmt.Println(string(body))

	var PostmarkResponse dto.PostmarkResponse

	json.Unmarshal(body, &PostmarkResponse)

	fmt.Println(PostmarkResponse, "ddjsk")

	log.Println("ending of postmark", time.Now())

	return &banyanRes, nil

}

//BGLD_ps_1.5-1.9
//The given book request will be validated and error will be thrown
func ValidateBookRequest(bookReq *dto.BookLTLRequestDTO) []dto.ValidateResDTO {

	//instance of the playground validator in created here
	validate := validator.New()

	err := validate.Struct(bookReq)

	fmt.Println(err)
	errorArray := make([]dto.ValidateResDTO, 0)

	fmt.Printf("bookReq: %+v", bookReq)

	// All the error for validator will be caught here and  thrown to the bookLTL main Service
	if err != nil {

		if strings.Contains(err.Error(), "'LoadId' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL002", Message: "Missing Required Parameter,LoadId"})
		}
		if strings.Contains(err.Error(), "'QuoteId' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL003", Message: "Missing Required Parameter,QuoteId"})
		}

		if strings.Contains(err.Error(), "'ShipName' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL005", Message: "Missing Required Parameter,ShipName"})
		}
		if strings.Contains(err.Error(), "'ShipAddress1' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL006", Message: "Missing Required Parameter,ShipAddress1"})
		}
		if strings.Contains(err.Error(), "'ShipCity' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL007", Message: "Missing Required Parameter,ShipCity"})
		}
		if strings.Contains(err.Error(), "'ShipState' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL008", Message: "Missing Required Parameter,ShipState"})
		}
		if strings.Contains(err.Error(), "'ShipCountry' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL008", Message: "Missing Required Parameter,ShipCountry"})
		}
		if strings.Contains(err.Error(), "'ShipZipCode' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL009", Message: "Missing Required Parameter,ShipZipCode"})
		}
		if strings.Contains(err.Error(), "'ShipEarliestDate' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL010", Message: "Missing Required Parameter,ShipEarliestDate"})
		}

		if strings.Contains(err.Error(), "'ConsName' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL011", Message: "Missing Required Parameter,ConsNam"})
		}

		if strings.Contains(err.Error(), "'ConsAddress1' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL012", Message: "Missing Required Parameter,ConsAddress1"})
		}

		if strings.Contains(err.Error(), "'ConsCity' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL013", Message: "Missing Required Parameter,ConsCity"})
		}
		if strings.Contains(err.Error(), "'ConsState' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL014", Message: "Missing Required Parameter,ConsState"})
		}

		if strings.Contains(err.Error(), "'ConsZipCode' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL015", Message: "Missing Required Parameter,ConsZipCode"})
		}

		if strings.Contains(err.Error(), "'ConsCountry' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL016", Message: "Missing Required Parameter,ConsCountry"})
		}

		if strings.Contains(err.Error(), "'ConsEarliestDate' failed on the 'required' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL017", Message: "Missing Required Parameter,ConsEarliestDate"})
		}

		if len(bookReq.ShipperDetails.ShipName) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL032", Message: "Invalid Required Parameter,ShipName"})
		}
		if len(bookReq.ShipperDetails.ShipAddress1) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL033", Message: "Invalid Required Parameter,ShipAddress1"})
		}
		if len(bookReq.ShipperDetails.ShipAddress2) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL034", Message: "Invalid Required Parameter,ShipAddress2"})
		}
		if len(bookReq.ShipperDetails.ShipContactName) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL035", Message: "Invalid Required Parameter,ShipContactName"})
		}

		if len(bookReq.ShipperDetails.ShipLoadNotes) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL036", Message: "Invalid Required Parameter,ShipLoadNotes"})
		}

		if len(bookReq.ConsigneeDetails.ConsName) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL042", Message: "Invalid Parameter,ConsName"})
		}
		if len(bookReq.ConsigneeDetails.ConsAddress1) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL043", Message: "Invalid Parameter,ConsAddress1"})
		}
		if len(bookReq.ConsigneeDetails.ConsAddress2) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL044", Message: "Invalid Parameter,ConsAddress1"})
		}
		if len(bookReq.ConsigneeDetails.ConsContactName) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL045", Message: "Invalid Parameter,ConsContactName"})
		}

		if len(bookReq.ConsigneeDetails.ConsLoadNotes) > 255 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL046", Message: "Invalid Parameter,ConsLoadNotes"})
		}

		//userId required
		if bookReq.UserId == "" {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL029", Message: "Missing Required Headers,UserId"})
		}
		if bookReq.SunteckLoginId == "" {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL030", Message: "Missing Required Headers,LoginId"})
		}

		if len(bookReq.Commodities.Commodities) == 0 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL018", Message: "Missing Required Parameter,Desc in Commodities"})
		}
		for index, _ := range bookReq.Commodities.Commodities {
			var CubicFt float64
			var PCF float64
			if bookReq.Commodities.Commodities[index].Length != 0 && bookReq.Commodities.Commodities[index].Width != 0 && bookReq.Commodities.Commodities[index].Height != 0 && bookReq.Commodities.Commodities[index].Quantity != 0 {

				lengthFloat := float64(bookReq.Commodities.Commodities[index].Length)
				WidthFloat := float64(bookReq.Commodities.Commodities[index].Width)
				HeightFloat := float64(bookReq.Commodities.Commodities[index].Height)
				QuantityFloat := float64(bookReq.Commodities.Commodities[index].Quantity)
				weightFloat := float64(bookReq.Commodities.Commodities[index].Weight)
				CubicFt = math.Floor(float64((lengthFloat*WidthFloat*HeightFloat)/1728) * QuantityFloat)
				PCF = weightFloat / CubicFt
				PCF = (math.Ceil(PCF*10) / 10)

			}

			if len(bookReq.Commodities.Commodities[index].Desc) > 255 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL053", Message: "Invalid Parameter,Desc in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if len(bookReq.Commodities.Commodities[index].NMFC) > 255 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL054", Message: "Invalid Parameter,NMFC in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			//if len(bookReq.Commodities[index].Density) > 255 {
			//	errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL054", Message: "Invalid Parameter,NMFC in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			//}
			if bookReq.Commodities.Commodities[index].Desc == "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL018", Message: "Missing Required Parameter,Desc in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].NMFC == "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL019", Message: "Missing Required Parameter,NMFC in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Class == "0" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL020", Message: "Missing Required Parameter,Class in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Class == "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL020", Message: "Missing Required Parameter,Class in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Length == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL021", Message: "Missing Required Parameter,Length in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Width == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL021", Message: "Missing Required Parameter,Width in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Weight == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL022", Message: "Missing Required Parameter,Weight in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Height == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL023", Message: "Missing Required Parameter,Height in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.Commodities[index].Quantity == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL024", Message: "Missing Required Parameter,Quantity in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			if bookReq.Commodities.LinearFt != "" && bookReq.Commodities.LinearFt > "12" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL026", Message: "Invalid Parameter,LinearFt in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}

			if bookReq.Commodities.Commodities[index].EquipmentType == "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL027", Message: "Missing Required Parameter,EquipmentType in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}

			if strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "BAGS" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "BUNDLE" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "CARTON" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "CRATES" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "DRUMS" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "PALLETS" && strings.ToUpper(bookReq.Commodities.Commodities[index].EquipmentType) != "ROLLS" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL027", Message: "Invalid parameter,EquipmentType in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			//if bookReq.Commodities.Commodities[index].Density > "0" && bookReq.Commodities.Commodities[index].Density != strconv.FormatFloat(PCF, 'g', 5, 64) {
			//	errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL026", Message: "Invalid Parameter,Density in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			//}

			if bookReq.Commodities.Commodities[index].CubicFt != strconv.FormatFloat(CubicFt, 'g', 5, 64) && bookReq.Commodities.Commodities[index].CubicFt != "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL025", Message: "Invalid Parameter,CubicFt in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
			}
			//res1 := strings.Trim(bookReq.ShipperDetails.ShipPhone, "-() ")
			res1 := strings.ReplaceAll(bookReq.ShipperDetails.ShipPhone, "-", "")
			res2 := strings.ReplaceAll(res1, "(", "")
			res3 := strings.ReplaceAll(res2, ")", "")
			res4 := strings.ReplaceAll(res3, " ", "")
			checkPhoneNumber, _ := regexp.MatchString(`^[0-9]{10}$`, res4)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipPhone != "" && !checkPhoneNumber {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL057", Message: "Invalid Parameter,ShipPhone"})

			}
			//res2 := strings.Trim(bookReq.ShipperDetails.ShipFax, "-() ")
			resp1 := strings.ReplaceAll(bookReq.ShipperDetails.ShipFax, "-", "")
			resp2 := strings.ReplaceAll(resp1, "(", "")
			resp3 := strings.ReplaceAll(resp2, ")", "")
			resp4 := strings.ReplaceAll(resp3, " ", "")
			checkFaxNumber, _ := regexp.MatchString(`^[0-9]{10}$`, resp4)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipFax != "" && !checkFaxNumber {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL058", Message: "Invalid Parameter,ShipFax"})

			}

			rsp1 := strings.ReplaceAll(bookReq.ConsigneeDetails.ConsPhone, "-", "")
			rsp2 := strings.ReplaceAll(rsp1, "(", "")
			rsp3 := strings.ReplaceAll(rsp2, ")", "")
			rsp4 := strings.ReplaceAll(rsp3, " ", "")
			//res3 := strings.Trim(bookReq.ConsigneeDetails.ConsPhone, "-() ")
			checkConsPhoneNumber, _ := regexp.MatchString(`^[0-9]{10}$`, rsp4)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsPhone != "" && !checkConsPhoneNumber {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL058", Message: "Invalid Parameter,ConsPhone"})

			}

			Rsp1 := strings.ReplaceAll(bookReq.ConsigneeDetails.ConsFax, "-", "")
			Rsp2 := strings.ReplaceAll(Rsp1, "(", "")
			Rsp3 := strings.ReplaceAll(Rsp2, ")", "")
			Rsp4 := strings.ReplaceAll(Rsp3, " ", "")
			//Res4 := strings.Trim(bookReq.ConsigneeDetails.ConsPhone, "-() ")
			checkConsFaxNumber, _ := regexp.MatchString(`^[0-9]{10}$`, Rsp4)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsFax != "" && !checkConsFaxNumber {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL059", Message: "Invalid Parameter,ConsFax"})

			}

			checkShipEmailType, _ := regexp.MatchString(`^[_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*(\+[a-zA-Z0-9-]+)?@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`, bookReq.ShipperDetails.ShipEmail)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipEmail != "" && !checkShipEmailType {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL049", Message: "Invalid Parameter,ShipEmail"})

			}

			checkConsEmailType, _ := regexp.MatchString(`^[_a-zA-Z0-9-]+(\.[_a-zA-Z0-9-]+)*(\+[a-zA-Z0-9-]+)?@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`, bookReq.ConsigneeDetails.ConsEmail)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsEmail != "" && !checkConsEmailType {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL050", Message: "Invalid Parameter,ConsEmail"})

			}

			checkShipEarliestDateFormat, _ := regexp.MatchString(`^\d{4}\-(0\d|1[0-2])\-(0[0-9]|1\d|2\d|3[01])$`, bookReq.ShipperDetails.ShipEarliestDate)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipEarliestDate != "" && !checkShipEarliestDateFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL037", Message: "Invalid Parameter,ShipEarliestDate"})
			}

			checkShipLatestDateFormat, _ := regexp.MatchString(`^\d{4}\-(0\d|1[0-2])\-(0[0-9]|1\d|2\d|3[01])$`, bookReq.ShipperDetails.ShipLatestDate)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipLatestDate != "" && !checkShipLatestDateFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL038", Message: "Invalid Parameter,ShipLatestDate"})
			}

			checkConsEarliestDateFormat, _ := regexp.MatchString(`^\d{4}\-(0\d|1[0-2])\-(0[0-9]|1\d|2\d|3[01])$`, bookReq.ConsigneeDetails.ConsEarliestDate)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsEarliestDate != "" && !checkConsEarliestDateFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL051", Message: "Invalid Parameter,ConsEarliestDate"})
			}

			checkConsShipLatestDateFormat, _ := regexp.MatchString(`^\d{4}\-(0\d|1[0-2])\-(0[0-9]|1\d|2\d|3[01])$`, bookReq.ConsigneeDetails.ConsLatestDate)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsLatestDate != "" && !checkConsShipLatestDateFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL052", Message: "Invalid Parameter,ConsLatestDate"})
			}

			checkShipEarliestTimeFormat, _ := regexp.MatchString(`([01]?[0-9]|2[0-3]):[0-5][0-9]`, bookReq.ShipperDetails.ShipEarliestTime)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipEarliestTime != "" && !checkShipEarliestTimeFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL053", Message: "Invalid Parameter,ShipEarliestTime"})
			}

			checkTimeFormat, _ := regexp.MatchString(`([01]?[0-9]|2[0-3]):[0-5][0-9]`, bookReq.ShipperDetails.ShipLatestTime)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ShipperDetails.ShipLatestTime != "" && !checkTimeFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL054", Message: "Invalid Parameter,ShipLatestTime"})
			}

			checkConsTimeFormat, _ := regexp.MatchString(`([01]?[0-9]|2[0-3]):[0-5][0-9]`, bookReq.ConsigneeDetails.ConsEarliestTime)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsEarliestTime != "" && !checkConsTimeFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL055", Message: "Invalid Parameter,ConsEarliestTime"})
			}

			checkConslatestTimeFormat, _ := regexp.MatchString(`([01]?[0-9]|2[0-3]):[0-5][0-9]`, bookReq.ConsigneeDetails.ConsLatestTime)
			//errorsArr := make([]dto.ValidateResponse, 0)
			if bookReq.ConsigneeDetails.ConsLatestTime != "" && !checkConslatestTimeFormat {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL056", Message: "Invalid Parameter,ConsLatestTime"})
			}

			currentTime := time.Now()
			CurrentDateS := currentTime.String()
			//currenDate := currentTime.Format("2017-09-07")
			tempArray := strings.Split(CurrentDateS, " ")
			fmt.Println(currentTime, "ghjkjhgfghjjhgfgh")
			fmt.Println(CurrentDateS, "fghjk")
			fmt.Println(tempArray[0], "sdfghjk")

			if bookReq.ShipperDetails.ShipEarliestDate < tempArray[0] {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL057", Message: "ShipEarliestDate should be greater or equal to the present date"})

			}

			if bookReq.ConsigneeDetails.ConsEarliestDate < tempArray[0] {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL059", Message: "ConsEarliestDate should be greater or equal to the present date"})
			}

			if bookReq.ShipperDetails.ShipEarliestDate > bookReq.ShipperDetails.ShipLatestDate && bookReq.ShipperDetails.ShipLatestDate != "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL061", Message: "ShipLatestDate should be greater or equal to the  ShipEarliestDate"})
			}
			if bookReq.ConsigneeDetails.ConsEarliestDate > bookReq.ConsigneeDetails.ConsLatestDate && bookReq.ConsigneeDetails.ConsLatestDate != "" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL062", Message: "ConsLatestDate should be greater or equal to the the ConsEarliestDate"})
			}
			if bookReq.ShipperDetails.ShipEarliestDate > bookReq.ConsigneeDetails.ConsEarliestDate {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BLTL063", Message: "ConsEarliestDate should be greater than the ShipEarliestDate."})
			}

		}

	}

	return errorArray
}

//BGLD_ps_1.3-1.10
//This function will be triggered from bookLtlService()
//This function will be used to check the rating engine and compare response from banyan or quote to create a valid load
func (r DefaultBookingService) execute(bookReq *dto.BookLTLRequestDTO) (*dto.BookLTLResponseDTO, []dto.ValidateResDTO, *errs.AppErrorvalidation) {

	//variable to stack all the error messages
	errorArray := make([]dto.ValidateResDTO, 0)

	//BGLD_ps_1.3-1.11
	// This function will get all the customer details like custmId and CustaId and Office code etc... from db
	log.Println("Starting of GetCustDataRepo", time.Now())

	custDataResp, errCust := r.repo.GetCustDataRepo(bookReq.SunteckLoginId)

	log.Println("ending of GetCustDataRepo", time.Now())
	//500 error will be caught here
	if errCust != nil {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: errCust.Code, Message: errCust.Message})
		return nil, nil, errs.ValidateResponse(nil, 500, errCust.Message)
	}
	//error will be caught here if there is no data in database
	if custDataResp == nil {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL001", Message: "No customer data found for the requested loginId"})
		return nil, errorArray, nil
	}

	fmt.Println("custDataResp", custDataResp)

	officeCode := custDataResp.OfficeCode
	loadId := bookReq.LoadId

	bookReq.CustaId = custDataResp.CustaId
	bookReq.CustmId = custDataResp.CustmId
	fmt.Println(bookReq.CustaId, bookReq.CustmId, officeCode)

	tempQuoteId := 0
	log.Println("Starting of cognito", time.Now())
	log.Println(constant.ACCESS_TOKEN_GRANTTYPE)
	log.Println(constant.ACCESS_TOKEN_CLIENTID)
	log.Println(constant.ACCESS_TOKEN_CLIENTSECRET)
	log.Println(constant.ACCESS_TOKEN_URL)

	//BGLD_ps_1.3-1.11
	//Cognito External api call flow starts here

	url := constant.ACCESS_TOKEN_URL + "/oauth2/token"
	method := "POST"

	payload := strings.NewReader("grant_type=" + constant.ACCESS_TOKEN_GRANTTYPE + "&client_id=" + constant.ACCESS_TOKEN_CLIENTID + "&client_secret=" + constant.ACCESS_TOKEN_CLIENTSECRET)

	client := &http.Client{}
	request, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)

	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Cookie", "XSRF-TOKEN=0364cc51-b71a-4b17-b5fe-90a39cb86154")

	res, err := client.Do(request)

	fmt.Println(res, "res")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)

	}
	var accessTokenResponse dto.AccessTokenResponse
	json.Unmarshal(body, &accessTokenResponse)

	//BGLD_ps_1.3-1.13
	//The rating engine whether it is mode or banyan will be known by CheckRatingEngine(bookReq.CustmId)

	ratingEngine, ratingEngineErr := r.repo.CheckRatingEngine(bookReq.CustmId)

	//error will be caught here
	if ratingEngineErr != nil {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: ratingEngineErr.Code, Message: ratingEngineErr.Message})
		return nil, nil, errs.ValidateResponse(nil, 500, ratingEngineErr.Message)

	}

	fmt.Println(ratingEngine)
	fmt.Println(ratingEngineErr)
	if ratingEngine.RatingEngine == "Mode" {

		//BGLD_ps_1.3-1.16
		//The QuoteDetails() will be called when Rating engine is "Mode"
		//In this function it gets the xml response of quotedetails to compare the rquest with xml Response for valid load creation

		QuoteDetailsRes, QuoteDetailsErr := r.repo.QuoteDetails(loadId)
		fmt.Println(QuoteDetailsRes)

		//Error will be caught here
		if QuoteDetailsErr != nil {
			fmt.Println(QuoteDetailsErr.Errors, QuoteDetailsErr.Message)
			errorArray = append(errorArray, dto.ValidateResDTO{Code: QuoteDetailsErr.Code, Message: QuoteDetailsErr.Message})
			return nil, errorArray, nil
		}

		var response dto.QuoteDetailsDto

		byteArray := []byte(QuoteDetailsRes.XmlTraffic)

		json.Unmarshal(byteArray, &response)

		fmt.Println(response, "db json response")

		//Looping statement will be introduced to get the details of quote which is selected by user
		//After getting Quoting details the details of shipper details , consignee Details and commodity details will be compared with request
		//Error will be thrown if any of the data is unmatched with request

		log.Println("Starting  of QuoteDetails comparision logic", time.Now())
		for index, _ := range response.Response.MercuryResponseDto.PriceSheets.PriceSheet {

			if response.Response.MercuryResponseDto.PriceSheets.PriceSheet[index].AssociatedCarrierPricesheet.PriceSheet.Type != "" {
				fmt.Println(index, "for iteration")
				str := response.Response.MercuryResponseDto.PriceSheets.PriceSheet[index].AssociatedCarrierPricesheet.PriceSheet.QuoteInformation.QuoteNumber

				tempArray := strings.Split(str, ";")

				fmt.Println(tempArray, "sample string")
				fmt.Println(len(tempArray), "sample string")

				if len(tempArray) > 1 {
					fmt.Println(tempArray, "shdshshs")
					fmt.Println(tempArray[1], "resultsgsgsg")
					if tempArray[1] == strconv.Itoa(bookReq.QuoteId) {
						tempQuoteId = bookReq.QuoteId
					}
				}
			}

		}
		if tempQuoteId == 0 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,QuoteId"})
			return nil, errorArray, nil
		}

		if response.Request.Events[0].Zip != bookReq.ShipperDetails.ShipZipCode {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipZipCode"})
			return nil, errorArray, nil
		}

		if strings.ToUpper(response.Request.Events[0].City) != strings.ToUpper(bookReq.ShipperDetails.ShipCity) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipCity"})
			return nil, errorArray, nil
		}

		if strings.Contains(strings.ToUpper(bookReq.ShipperDetails.ShipCountry), strings.ToUpper(response.Request.Events[0].Country)) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipCountry"})
			return nil, errorArray, nil
		}

		if strings.ToUpper(response.Request.Events[0].State) != strings.ToUpper(bookReq.ShipperDetails.ShipState) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipState"})
			return nil, errorArray, nil
		}

		if response.Request.Events[1].Zip != bookReq.ConsigneeDetails.ConsZipCode {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsZipCode"})
			return nil, errorArray, nil
		}

		if strings.ToUpper(response.Request.Events[1].City) != strings.ToUpper(bookReq.ConsigneeDetails.ConsCity) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsCity"})
			return nil, errorArray, nil
		}

		if strings.Contains(strings.ToUpper(bookReq.ConsigneeDetails.ConsCountry), strings.ToUpper(response.Request.Events[1].Country)) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsCountry"})
			return nil, errorArray, nil
		}

		if strings.ToUpper(response.Request.Events[1].State) != strings.ToUpper(bookReq.ConsigneeDetails.ConsState) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsState"})
			return nil, errorArray, nil
		}

		if len(response.Request.Items) != len(bookReq.Commodities.Commodities) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,commodities"})
			return nil, errorArray, nil
		}

		for index, element := range response.Request.Items {
			fmt.Println(element)
			if response.Request.Items[index].Quantity != strconv.Itoa(bookReq.Commodities.Commodities[index].Quantity) || response.Request.Items[index].Weight != strconv.Itoa(bookReq.Commodities.Commodities[index].Weight) || response.Request.Items[index].Class != bookReq.Commodities.Commodities[index].Class || response.Request.Items[index].Length != strconv.Itoa(bookReq.Commodities.Commodities[index].Length) || response.Request.Items[index].Width != strconv.Itoa(bookReq.Commodities.Commodities[index].Width) || response.Request.Items[index].Height != strconv.Itoa(bookReq.Commodities.Commodities[index].Height) ||
				response.Request.Items[index].Type != bookReq.Commodities.Commodities[index].EquipmentType {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Commodities Details"})
				return nil, errorArray, nil
			}
		}
		log.Println("Ending of QuoteDetails logic comparision", time.Now())
		//BGLD_ps_1.20-1.27
		//Valid load for mode rating engine will be created with InsertBookDetailsModeRepo()

		InsertBookDetailsRes, InsertBookDetailsErr := r.repo.InsertBookDetailsModeRepo(&response, bookReq)

		if InsertBookDetailsErr != nil {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: InsertBookDetailsErr.Code, Message: InsertBookDetailsErr.Message})
			return nil, nil, errs.ValidateResponse(nil, 500, InsertBookDetailsErr.Message)
		}

		//Response will be formed here and passed to bookLtlService function
		responseDTO := dto.BookLTLResponseDTO{
			LoadNumber:  InsertBookDetailsRes.LoadNumber,
			QuoteId:     InsertBookDetailsRes.QuoteId,
			QuoteNumber: InsertBookDetailsRes.QuoteNumber,
			AccessToken: accessTokenResponse.AccessToken,
			AgentEmail:  InsertBookDetailsRes.AgentEmail,
			CustEmail:   InsertBookDetailsRes.CustEmail,
			CustFax:     InsertBookDetailsRes.CustFax,
			CustPhone:   InsertBookDetailsRes.CustPhone,
			Contact:     InsertBookDetailsRes.Contact,
			PriceDetails: dto.PriceDetails{
				Scac:               InsertBookDetailsRes.PriceDetails.Scac,
				Service:            InsertBookDetailsRes.PriceDetails.Service,
				CarrierName:        InsertBookDetailsRes.PriceDetails.CarrierName,
				CarrierNotes:       InsertBookDetailsRes.PriceDetails.CarrierNotes,
				TransitTime:        InsertBookDetailsRes.PriceDetails.TransitTime,
				FlatPrice:          InsertBookDetailsRes.PriceDetails.FlatPrice,
				FuelSurchargePrice: InsertBookDetailsRes.PriceDetails.FuelSurchargePrice,
			},
			TotalPrice: InsertBookDetailsRes.TotalPrice,
		}
		return &responseDTO, nil, nil
	} else {
		log.Println("starting of banyan", time.Now())

		//BanyanGetLoadDetails
		//BGLD_ps_1.3-1.24
		//Banyan external api call starts here
		BanyanUrl := constant.BANYAN_URL + "/getloaddetails"
		BanyanMethod := "POST"

		BanyanPayload := strings.NewReader(`{` + "" + `	"AuthData" : {` + "" + `		"CustaId"    :"` + bookReq.CustaId + `",` + "" + `
		"CustmId"   :"` + bookReq.CustmId + `",` + "" + `
		"OfficeCode" :"` + officeCode + `"` + "" + `
	} ,` + "" + `
	"LoadId":  ` + strconv.Itoa(bookReq.LoadId) + "" + `
}`)

		log.Println(BanyanUrl)
		log.Println(bookReq.CustaId)
		log.Println(BanyanPayload)
		log.Println(officeCode)

		BanyanClient := &http.Client{}
		Banyanreq, err2 := http.NewRequest(BanyanMethod, BanyanUrl, BanyanPayload)

		if err2 != nil {
			fmt.Println(err2)
			log.Println("error 1", err2)

		}

		Banyanreq.Header.Add("Authorization", "Bearer "+accessTokenResponse.AccessToken)

		Banyanreq.Header.Add("Content-Type", "application/json")
		fmt.Println(Banyanreq.Header, "Banyanreq.Header")
		log.Println(Banyanreq.Header)
		res2, err := BanyanClient.Do(Banyanreq)

		if err != nil {
			log.Println("error 2", err)

			fmt.Println(err)

		}
		defer res2.Body.Close()

		log.Println("check 1")

		body2, err := ioutil.ReadAll(res2.Body)
		if err != nil {
			log.Println("error 3", err)
			fmt.Println(err)

		}
		fmt.Println(string(body2), "12345")

		var banResponse *dto.BanyanGetLoadDetailsDto
		var Error dto.Error
		json.Unmarshal(body2, &Error)
		json.Unmarshal(body2, &banResponse)

		fmt.Println("banResponse123234", Error)

		fmt.Println("banResponse123", banResponse)

		if banResponse.Success == false {
			if len(Error.Errors) == 0 {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "BGLD009", Message: "Load " + strconv.Itoa(bookReq.LoadId) + " does not exist for this user."})
				return nil, errorArray, nil
			}
		}

		for index, element := range Error.Errors {
			fmt.Println(element)

			if Error.Errors[index].Message == "CustaId does not exist" {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: Error.Errors[index].Code, Message: Error.Errors[index].Message})
				return nil, errorArray, nil
			}

			//if Error.Errors[index].Message == "CustmId does not exist" {
			//	errorArray = append(errorArray, dto.ValidateResDTO{Code: Error.Errors[index].Code, Message: Error.Errors[index].Message})
			//	return nil, errorArray
			//}

			if Error.Errors[index].Message == "Load "+strconv.Itoa(bookReq.LoadId)+" does not exist for this user." {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: Error.Errors[index].Code, Message: "Load " + strconv.Itoa(bookReq.LoadId) + " does not exist for this user."})
				return nil, errorArray, nil
			}
			if Error.Errors[index].Message == "HTTP 400: Not authorized." {

				errorArray = append(errorArray, dto.ValidateResDTO{Code: Error.Errors[index].Code, Message: "Load " + strconv.Itoa(bookReq.LoadId) + " does not exist for this user."})
				return nil, errorArray, nil
			}
		}
		var banyanQuoteId = 0
		log.Println("ending of banyan", time.Now())

		//Looping statement will be introduced to get the details of quote which is selected by user
		//After getting Quoting details the details of shipper details , consignee Details and commodity details will be compared with request
		//Error will be thrown if any of the data is unmatched with request

		log.Println("Starting of banyan comparision logic", time.Now())

		for i, e := range banResponse.LoadDetails[0].Quotes {
			fmt.Println(e)
			{
				if bookReq.QuoteId == banResponse.LoadDetails[0].Quotes[i].QuoteID {
					banyanQuoteId = banResponse.LoadDetails[0].Quotes[i].QuoteID

				}
				if banResponse.LoadDetails[0].Shipper.AddressInfo.Zipcode != bookReq.ShipperDetails.ShipZipCode {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipZipCode"})
					return nil, errorArray, nil
				}

				if strings.ToUpper(banResponse.LoadDetails[0].Shipper.AddressInfo.City) != strings.ToUpper(bookReq.ShipperDetails.ShipCity) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipCity"})
					return nil, errorArray, nil
				}

				if strings.Contains(strings.ToUpper(bookReq.ShipperDetails.ShipCountry), strings.ToUpper(banResponse.LoadDetails[0].Shipper.AddressInfo.CountryCode)) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipCountry"})
					return nil, errorArray, nil
				}

				if strings.ToUpper(banResponse.LoadDetails[0].Shipper.AddressInfo.State) != strings.ToUpper(bookReq.ShipperDetails.ShipState) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ShipState"})
					return nil, errorArray, nil
				}

				if banResponse.LoadDetails[0].Consignee.AddressInfo.Zipcode != bookReq.ConsigneeDetails.ConsZipCode {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsZipCode"})
					return nil, errorArray, nil
				}

				if strings.ToUpper(banResponse.LoadDetails[0].Consignee.AddressInfo.City) != strings.ToUpper(bookReq.ConsigneeDetails.ConsCity) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsCity"})
					return nil, errorArray, nil
				}

				if strings.Contains(strings.ToUpper(bookReq.ConsigneeDetails.ConsCountry), strings.ToUpper(banResponse.LoadDetails[0].Shipper.AddressInfo.CountryCode)) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsCountry"})
					return nil, errorArray, nil
				}
				if strings.ToUpper(banResponse.LoadDetails[0].Consignee.AddressInfo.State) != strings.ToUpper(bookReq.ConsigneeDetails.ConsState) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,ConsState"})
					return nil, errorArray, nil
				}
				if len(banResponse.LoadDetails[0].Products) != len(bookReq.Commodities.Commodities) {
					errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details"})
					return nil, errorArray, nil
				}

				for index, element := range banResponse.LoadDetails[0].Products {
					fmt.Println(element)
					if banResponse.LoadDetails[0].Products[index].Quantity != bookReq.Commodities.Commodities[index].Quantity {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Quantity in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}
					if banResponse.LoadDetails[0].Products[index].Weight != bookReq.Commodities.Commodities[index].Weight {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Weight in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}

					//if banResponse.LoadDetails[0].Products[index].Class != bookReq.Commodities.Commodities[index].Class {
					//	errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Class in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
					//	return nil, errorArray, nil
					//}

					if banResponse.LoadDetails[0].Products[index].Length != bookReq.Commodities.Commodities[index].Length {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Length in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}

					if banResponse.LoadDetails[0].Products[index].Width != bookReq.Commodities.Commodities[index].Width {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Width in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}

					if banResponse.LoadDetails[0].Products[index].Height != bookReq.Commodities.Commodities[index].Height {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Height in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil

					}

					if strings.ToUpper(banResponse.LoadDetails[0].Products[index].Description) != strings.ToUpper(bookReq.Commodities.Commodities[index].Desc) {

						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Desc in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}
					if banResponse.LoadDetails[0].Products[index].IsHazmat != bookReq.Commodities.Commodities[index].Hazmat {
						errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,Hazmat in Commodities" + "[" + strconv.Itoa(index+1) + "]"})
						return nil, errorArray, nil
					}
				}

			}
		}
		if banyanQuoteId == 0 {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "BL028", Message: "Invalid Load Details,QuoteId"})
			return nil, errorArray, nil
		}
		//for i, _ := range banResponse.LoadDetails[0].ShipperAccessorials.AppointmentRequired {
		//
		//}
		var AccessorialsArr []string
		if banResponse.LoadDetails[0].ShipperAccessorials.AppointmentRequired == true {
			AccessorialsArr = append(AccessorialsArr, "AppointmentRequired")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.InsidePickup == true {
			AccessorialsArr = append(AccessorialsArr, "InsideDelivery")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.LiftgatePickup == true {
			AccessorialsArr = append(AccessorialsArr, "LiftGate (PICK)")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.NonBusinessHourPickup == true {
			AccessorialsArr = append(AccessorialsArr, "NonBusinessHourPickup")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.MarkingTagging == true {
			AccessorialsArr = append(AccessorialsArr, "MarkingTagging")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.NYCMetro == true {
			AccessorialsArr = append(AccessorialsArr, "NYCMetro")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.PalletJack == true {
			AccessorialsArr = append(AccessorialsArr, "PalletJack")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.ResidentialPickup == true {
			AccessorialsArr = append(AccessorialsArr, "ResidentialPickup")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.SortSegregate == true {
			AccessorialsArr = append(AccessorialsArr, "SortSegregate (PICK)")
		}
		if banResponse.LoadDetails[0].ShipperAccessorials.TradeShowPickup == true {
			AccessorialsArr = append(AccessorialsArr, "TradeShowPickup")
		}

		//////
		if banResponse.LoadDetails[0].ConsigneeAccessorials.AppointmentRequired == true {
			AccessorialsArr = append(AccessorialsArr, "AppointmentRequired (DROP)")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.InsideDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "InsideDelivery (DROP)")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.SortSegregate == true {
			AccessorialsArr = append(AccessorialsArr, "SortSegregate (DROP)")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.PalletJack == true {
			AccessorialsArr = append(AccessorialsArr, "PalletJack")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.ResidentialDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "InsideDelivery (DROP)")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.LiftgateDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "LiftgateDelivery")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.MarkingTagging == true {
			AccessorialsArr = append(AccessorialsArr, "MarkingTagging")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.TradeShowDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "TradeShowDelivery")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.NYCMetro == true {
			AccessorialsArr = append(AccessorialsArr, "NYCMetro")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.DeliveryNotification == true {
			AccessorialsArr = append(AccessorialsArr, "DeliveryNotification")
		}
		if banResponse.LoadDetails[0].ConsigneeAccessorials.TwoHourSpecialDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "TwoHourSpecialDelivery")
		}

		if banResponse.LoadDetails[0].ConsigneeAccessorials.NonBusinessHourDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "NonBusinessHourDelivery")
		}
		//////////////

		if banResponse.LoadDetails[0].LoadAccessorials.Guaranteed == true {
			AccessorialsArr = append(AccessorialsArr, "Guaranteed")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.TimeDefinite == true {
			AccessorialsArr = append(AccessorialsArr, "InsideDelivery (DROP)")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.Expedited == true {
			AccessorialsArr = append(AccessorialsArr, "SortSegregate (DROP)")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.HolidayPickup == true {
			AccessorialsArr = append(AccessorialsArr, "PalletJack")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.HolidayDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "InsideDelivery (DROP)")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.WeightDetermination == true {
			AccessorialsArr = append(AccessorialsArr, "LiftgateDelivery")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.BlindShipment == true {
			AccessorialsArr = append(AccessorialsArr, "MarkingTagging")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.BlanketService == true {
			AccessorialsArr = append(AccessorialsArr, "TradeShowDelivery")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.ProtectFromFreezing == true {
			AccessorialsArr = append(AccessorialsArr, "NYCMetro")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.SingleShipment == true {
			AccessorialsArr = append(AccessorialsArr, "DeliveryNotification")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.CustomsInBond == true {
			AccessorialsArr = append(AccessorialsArr, "TwoHourSpecialDelivery")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.OverDimension == true {
			AccessorialsArr = append(AccessorialsArr, "OverDimension")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.Stackable == true {
			AccessorialsArr = append(AccessorialsArr, "Stackable")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.Turnkey == true {
			AccessorialsArr = append(AccessorialsArr, "Turnkey")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.FoodGradeProducts == true {
			AccessorialsArr = append(AccessorialsArr, "FoodGradeProducts")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.TSA == true {
			AccessorialsArr = append(AccessorialsArr, "TSA")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.Bulkhead == true {
			AccessorialsArr = append(AccessorialsArr, "Bulkhead")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.SignatureRequired == true {
			AccessorialsArr = append(AccessorialsArr, "SignatureRequired")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.BlanketServiceChilled == true {
			AccessorialsArr = append(AccessorialsArr, "BlanketServiceChilled")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.BlanketServiceFrozen == true {
			AccessorialsArr = append(AccessorialsArr, "BlanketServiceFrozen")
		}

		if banResponse.LoadDetails[0].LoadAccessorials.SaturdayDelivery == true {
			AccessorialsArr = append(AccessorialsArr, "SaturdayDelivery")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.SecondMan == true {
			AccessorialsArr = append(AccessorialsArr, "SecondMan")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.ReturnReceipt == true {
			AccessorialsArr = append(AccessorialsArr, "ReturnReceipt")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.ShipmentHold == true {
			AccessorialsArr = append(AccessorialsArr, "ShipmentHold")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.ProactiveResponse == true {
			AccessorialsArr = append(AccessorialsArr, "ProactiveResponse")
		}
		if banResponse.LoadDetails[0].LoadAccessorials.ShipperRelease == true {
			AccessorialsArr = append(AccessorialsArr, "ShipperRelease")
		}
		fmt.Println("AccessorialsArr", AccessorialsArr)
		SelectedAccessorials := strings.Join(AccessorialsArr, " ")
		log.Println("Ending of banyan comparision logic", time.Now())

		//BGLD_ps_1.3-1.8
		//Valid load for banyan rating engine will be created with InsertBookDetailsModeRepo()

		log.Println("starting of InsertBookDetailsRepo", time.Now())

		InsertBookDetailsRes, InsertBookDetailsErr := r.repo.InsertBookDetailsRepo(banResponse, bookReq, SelectedAccessorials)

		log.Println("ending of InsertBookDetailsRepo", time.Now())

		if InsertBookDetailsErr != nil {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: InsertBookDetailsErr.Code, Message: InsertBookDetailsErr.Message})
			return nil, nil, errs.ValidateResponse(nil, 500, InsertBookDetailsErr.Message)
		}

		fmt.Println(InsertBookDetailsRes)
		//Response will be formed here and passed to bookLtlService function

		responseDTO := dto.BookLTLResponseDTO{
			LoadNumber:  InsertBookDetailsRes.LoadNumber,
			QuoteId:     InsertBookDetailsRes.QuoteId,
			QuoteNumber: InsertBookDetailsRes.QuoteNumber,
			AccessToken: accessTokenResponse.AccessToken,
			AgentEmail:  InsertBookDetailsRes.AgentEmail,
			CustEmail:   InsertBookDetailsRes.CustEmail,
			CustFax:     InsertBookDetailsRes.CustFax,
			CustPhone:   InsertBookDetailsRes.CustPhone,
			Contact:     InsertBookDetailsRes.Contact,
			PriceDetails: dto.PriceDetails{
				Scac:               InsertBookDetailsRes.PriceDetails.Scac,
				Service:            InsertBookDetailsRes.PriceDetails.Service,
				CarrierName:        InsertBookDetailsRes.PriceDetails.CarrierName,
				CarrierNotes:       InsertBookDetailsRes.PriceDetails.CarrierNotes,
				TransitTime:        InsertBookDetailsRes.PriceDetails.TransitTime,
				FlatPrice:          InsertBookDetailsRes.PriceDetails.FlatPrice,
				FuelSurchargePrice: InsertBookDetailsRes.PriceDetails.FuelSurchargePrice,
			},
			TotalPrice: InsertBookDetailsRes.TotalPrice,
		}
		return &responseDTO, nil, nil
	}

	return nil, nil, nil
}

func NewBookService(repository domain.BookRepository) DefaultBookingService {
	return DefaultBookingService{repository}
}
