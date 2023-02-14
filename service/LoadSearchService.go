package Service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"golang/constant"
	"golang/domain"
	"golang/dto"
	"golang/errs"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/jwt"
	"golang.org/x/exp/slices"
	"gopkg.in/go-playground/validator.v9"
)

//go:generate mockgen -destination=../mock/service/mockLoadSearchService.go -package=Service golang/service LoadService
type LoadService interface {
	GetLoadSearchResultService(resultVar *dto.LoadSearchReqDTO) (*dto.LoadSearchResultResDTO, *errs.AppErrorvalidation)
	GetLoadDocumentsService(loadRequest dto.LoadDocumentsReqDTO) (*dto.LoadDocumentsResDTO, *errs.AppErrorvalidation)
}

type DefaultLoadService struct {
	repo domain.LoadRepository
}

//--loadSearchResult API starts here
//AP_PC_03
// The below code is a loadSearchResult API service layer which is responsible for calling the repository layer and
// returning the response to the handler.
func (r DefaultLoadService) GetLoadSearchResultService(loadRequest *dto.LoadSearchReqDTO) (*dto.LoadSearchResultResDTO, *errs.AppErrorvalidation) {

	// AP_PC_05
	//validateRequest function is called and the loadRequest is passed as a parameter
	validateResponseVar := ValidateRequest(loadRequest)

	if len(validateResponseVar) != 0 {
		return nil, errs.ValidateResponse(validateResponseVar, 400, "")
	}

	//AP_PC_06
	//GetCPDBLoads function is called and the loadRequest is passed as a parameter
	otmsLodResponse, otmsCount, err := r.GetCPDBLoads(loadRequest)
	if err != nil {
		return nil, errs.ValidateResponse(nil, err.Code, err.Message)
	}

	//AP_PC_09
	//GetBTMSLoads function is called and the loadRequest is passed as a parameter
	btmsLodResponse, btmsCount, err := r.GetBTMSLoads(loadRequest)
	if err != nil {
		return nil, errs.ValidateResponse(nil, err.Code, err.Message)
	}

	// AP_PC_12
	//The DB responses are combined and sorted in here
	combinedResponse := append(otmsLodResponse, btmsLodResponse...)

	if (len(otmsLodResponse)) != 0 && (len(btmsLodResponse)) != 0 {

		//The combined response is sorted based on ascending or descending based on the request
		if loadRequest.SortColumn == "" || loadRequest.SortOrder == "" {
			sort.SliceStable(combinedResponse, func(i, j int) bool {
				return combinedResponse[i].LastModified > combinedResponse[j].LastModified
			})
		} else if strings.ToUpper(loadRequest.SortColumn) == "LOADID" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadId < combinedResponse[j].LoadId
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadId > combinedResponse[j].LoadId
				})
			}
		} else if strings.ToUpper(loadRequest.SortColumn) == "LOADSTATUS" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadStatus < combinedResponse[j].LoadStatus
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadStatus > combinedResponse[j].LoadStatus
				})
			}
		} else if strings.ToUpper(loadRequest.SortColumn) == "MOVETYPE" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadMethod < combinedResponse[j].LoadMethod
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].LoadMethod > combinedResponse[j].LoadMethod
				})
			}
		} else if strings.ToUpper(loadRequest.SortColumn) == "ORIGIN" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].ShipperCity < combinedResponse[j].ShipperCity
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].ShipperCity > combinedResponse[j].ShipperCity
				})
			}
		} else if strings.ToUpper(loadRequest.SortColumn) == "DESTINATION" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].ConsigneeCity < combinedResponse[j].ConsigneeCity
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].ConsigneeCity > combinedResponse[j].ConsigneeCity
				})
			}
		} else if strings.ToUpper(loadRequest.SortColumn) == "TOTAL" {
			if strings.ToUpper(loadRequest.SortOrder) == "ASCENDING" {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].CustTotal < combinedResponse[j].CustTotal
				})
			} else {
				sort.SliceStable(combinedResponse, func(i, j int) bool {
					return combinedResponse[i].CustTotal > combinedResponse[j].CustTotal
				})
			}
		}

	}
	btmsLoadRecords := 0
	otmsLoadRecords := 0
	finalLoadResponse := make([]dto.LoadSearchResDTO, 0)
	if (len(combinedResponse)) <= 10 {
		for i := 0; i < len(combinedResponse); i++ {
			if combinedResponse[i].LoadOrigin == "BTMS" {
				btmsLoadRecords++
				finalLoadResponse = append(finalLoadResponse, combinedResponse[i])
			} else {
				otmsLoadRecords++
				shipDateStartDateSplit := strings.Split(combinedResponse[i].ShipDateStart, "T")
				combinedResponse[i].ShipDateStart = shipDateStartDateSplit[0]
				shipDateEndDateSplit := strings.Split(combinedResponse[i].ShipDateEnd, "T")
				combinedResponse[i].ShipDateEnd = shipDateEndDateSplit[0]
				finalLoadResponse = append(finalLoadResponse, combinedResponse[i])
			}
		}
	} else {
		for i := 0; i < 10; i++ {
			if combinedResponse[i].LoadOrigin == "BTMS" {
				btmsLoadRecords++
				finalLoadResponse = append(finalLoadResponse, combinedResponse[i])
			} else {
				otmsLoadRecords++
				shipDateStartDateSplit := strings.Split(combinedResponse[i].ShipDateStart, "T")
				combinedResponse[i].ShipDateStart = shipDateStartDateSplit[0]
				shipDateEndDateSplit := strings.Split(combinedResponse[i].ShipDateEnd, "T")
				combinedResponse[i].ShipDateEnd = shipDateEndDateSplit[0]
				finalLoadResponse = append(finalLoadResponse, combinedResponse[i])
			}
		}
	}

	// fmt.Printf("finalResponse: %+v", finalLoadResponse)

	loadSearchResult := dto.LoadSearchResultResDTO{
		LoadResponse: finalLoadResponse,
		LoadCount: struct {
			OthertmsLoadRecords int `json:"othertmsLoadRecords"`
			BtmsLoadRecords     int `json:"btmsLoadRecords"`
			TotalLoadsCount     int `json:"totalLoadsCount"`
		}{
			OthertmsLoadRecords: otmsLoadRecords,
			BtmsLoadRecords:     btmsLoadRecords,
			TotalLoadsCount:     btmsCount + otmsCount,
		},
	}
	//AP_PC_13
	return &loadSearchResult, nil
}

//AP_PC_08
//GetCPDBLoads function is declared in here inside this function repository function is called and the load request is sent as a parameter
func (r DefaultLoadService) GetCPDBLoads(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]dto.LoadSearchResDTO, int, *errs.AppErrorvalidation) {

	otmsLoads, loadCountRes, err := r.repo.GetCPDBLoadResultRepo(LoadSearchReqDTO)

	if err != nil {
		return nil, 0, err
	}

	//DTO conversion is done here
	response := make([]dto.LoadSearchResDTO, 0)
	for _, c := range otmsLoads {
		response = append(response, c.ToLoadDto())
	}
	return response, loadCountRes, nil
}

//AP_PC_11
//GetBTMSLoads function is declared in here inside this function repository function is called and the load request is sent as a parameter
func (r DefaultLoadService) GetBTMSLoads(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]dto.LoadSearchResDTO, int, *errs.AppErrorvalidation) {

	btmsLoads, loadCountRes, err := r.repo.GetBTMSLoadResultRepo(LoadSearchReqDTO)

	if err != nil {
		return nil, 0, err
	}

	//DTO conversion is done here
	response := make([]dto.LoadSearchResDTO, 0)
	for _, c := range btmsLoads {
		response = append(response, c.ToLoadDto())
	}
	return response, loadCountRes, nil
}

// AP_PC_04
//ValidateRequest function is declared in here inside this function all the request validations are handled
func ValidateRequest(loadRequest *dto.LoadSearchReqDTO) []dto.ValidateResDTO {

	//instance of the playground validator in created here
	validate := validator.New()

	err := validate.Struct(loadRequest)

	errorArray := make([]dto.ValidateResDTO, 0)

	const layout = "2006-01-02"

	re := regexp.MustCompile(`^\d{4}\-(0\d|1[0-2])\-(0[0-9]|1\d|2\d|3[01])$`)
	IsLetter := regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

	if len(loadRequest.OfficeCode) == 0 {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Missing parameter, office code"})
	}
	if len(loadRequest.CustomerId) == 0 {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD002", Message: "Missing parameter,custmid"})
	}
	if len(loadRequest.Gcid) == 0 {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD003", Message: "Missing parameter,gcid"})
	}

	if strings.ToUpper(loadRequest.SortOrder) != "ASCENDING" && strings.ToUpper(loadRequest.SortOrder) != "DESCENDING" && loadRequest.SortOrder != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD004", Message: "Invalid parameter , sort order"})
	}
	if strings.ToUpper(loadRequest.SortColumn) != "LOADID" && strings.ToUpper(loadRequest.SortColumn) != "LOADSTATUS" &&
		strings.ToUpper(loadRequest.SortColumn) != "MOVETYPE" && strings.ToUpper(loadRequest.SortColumn) != "ORIGIN" &&
		strings.ToUpper(loadRequest.SortColumn) != "TOTAL" && strings.ToUpper(loadRequest.SortColumn) != "DESTINATION" &&
		loadRequest.SortColumn != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD005", Message: "Invalid parameter , sort column"})
	}
	if (loadRequest.SortOrder == "") && (loadRequest.SortColumn != "") {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD006", Message: "Missing parameter, sort order"})
	}
	if loadRequest.SortColumn == "" && loadRequest.SortOrder != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD007", Message: "Missing parameter, sort column"})
	}

	startDate, startDateError := time.Parse(layout, loadRequest.ShipDateStart)
	if !re.MatchString(loadRequest.ShipDateStart) && loadRequest.ShipDateStart != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD008", Message: "ship date start should be in this format `YYYY-MM-DD`"})
	} else if startDateError != nil && loadRequest.ShipDateStart != "" && startDate.String() == "0001-01-01 00:00:00 +0000 UTC" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD008", Message: "ship date start out of range"})
	}

	endDate, endDateError := time.Parse(layout, loadRequest.ShipDateEnd)
	if !re.MatchString(loadRequest.ShipDateEnd) && loadRequest.ShipDateEnd != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD009", Message: "ship date end should be in this format `YYYY-MM-DD`"})
	} else if endDateError != nil && loadRequest.ShipDateEnd != "" && endDate.String() == "0001-01-01 00:00:00 +0000 UTC" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD009", Message: "ship date end out of range"})
	}
	if loadRequest.ShipDateStart != "" && loadRequest.ShipDateEnd != "" && loadRequest.ShipDateStart > loadRequest.ShipDateEnd {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD010", Message: "End date should not be less than start date"})
	}
	if !IsLetter(loadRequest.LoadStatus) && loadRequest.LoadStatus != "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD011", Message: "Invalid parameter, load status"})
	}
	// if !IsLetter(loadRequest.LoadMethod) && loadRequest.LoadMethod != "" {
	// 	errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD012", Message: "Invalid parameter, load method"})
	// }

	if err != nil {

		// if strings.Contains(err.Error(), "'LoadId' failed on the 'numeric' tag") {
		// 	errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD013", Message: "Invalid parameter, load id"})
		// }
		if strings.Contains(err.Error(), "'OtmsLoadRecords' failed on the 'numeric' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD014", Message: "Invalid parameter, other tms Load Record"})
		}
		if strings.Contains(err.Error(), "'BtmsLoadRecords' failed on the 'numeric' tag") {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD015", Message: "Invalid parameter, btms Load Record"})
		}
	}

	return errorArray
}

//--loadSearchResult API ends here

//--GetLoadDocuments API starts here

//AP_PC_03
// The below code is a GetLoadDocuments API service layer which is responsible for calling the repository layer and
// returning the response to the handler.
func (r DefaultLoadService) GetLoadDocumentsService(LoadDocsReqDTO dto.LoadDocumentsReqDTO) (*dto.LoadDocumentsResDTO, *errs.AppErrorvalidation) {

	// AP_PC_05
	//validateRequest function is called and the LoadDocsReqDTO is passed as a parameter
	validateResponseVar := ValidateDocRequest(LoadDocsReqDTO)

	if len(validateResponseVar) != 0 {
		return nil, errs.ValidateResponse(validateResponseVar, 400, "")
	}

	//AP_PC_06
	//GetCPDBLoadDocs function is called and the LoadDocsReqDTO is passed as a parameter
	GetCPDBLoadDocs, err1 := r.GetCPDBLoadDocs(LoadDocsReqDTO)

	if err1 != nil {
		return nil, errs.ValidateResponse(nil, err1.Code, err1.Message)
	}

	GetBTMSDBLoadDocs, err := r.GetBTMSDBLoadDocs(LoadDocsReqDTO)

	if err != nil {
		return nil, errs.ValidateResponse(nil, err.Code, err.Message)
	}

	//AP_PC_12
	//The returned Document Response from the two different sources are combined as a single response in here
	combinedResponse := append(GetBTMSDBLoadDocs, GetCPDBLoadDocs...)

	loadDocumentRes := dto.LoadDocumentsResDTO{LoadDocResult: combinedResponse}
	return &loadDocumentRes, nil

}

//AP_PC_08
//GetCPDBLoadDocs function is declared in here inside this, the load document URL's for tritan and other tms loads is destructured and the response is returned
//as LoadDocsResult DTO
func (r DefaultLoadService) GetCPDBLoadDocs(LoadDocsReqDTO dto.LoadDocumentsReqDTO) ([]dto.LoadDocsResult, *errs.AppErrorvalidation) {

	othertmsDocs, err := r.repo.GetCPDBLoadDocsRepo(&LoadDocsReqDTO)

	if err != nil {
		return nil, errs.ValidateResponse(nil, err.Code, err.Message)
	}

	// fmt.Printf("othertmsDocs: %+v", othertmsDocs)

	otmsLoadDocRes := make([]dto.LoadDocsResult, 0)

	proofOfDelivery := 1
	billOfLading := 2
	customerConfirmation := 3
	lumperReceipt := 4
	customerReceipt := 5
	scaleWeight := 6
	accessorialApproval := 7

	for i := 0; i < len(LoadDocsReqDTO.OthertmsLoads); i++ {
		docArray := make([]dto.LoadDocuments, 0)
		for j := 0; j < len(othertmsDocs); j++ {
			if LoadDocsReqDTO.OthertmsLoads[i] == othertmsDocs[j].LoadId {
				switch othertmsDocs[j].TypeId {
				case proofOfDelivery:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Proof of Delivery",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case billOfLading:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Bill of Lading",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case customerConfirmation:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Customer Confirmation",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case lumperReceipt:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Lumper Receipt",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case customerReceipt:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Customer Receipt",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case scaleWeight:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Scale/Weight",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				case accessorialApproval:
					docArray = append(docArray, dto.LoadDocuments{
						Id:      strconv.Itoa(othertmsDocs[j].Id),
						LoadId:  othertmsDocs[j].LoadId,
						Path:    othertmsDocs[j].DocName,
						DocName: "Accessorial Approval",
						TypeId:  strconv.Itoa(othertmsDocs[j].TypeId),
						DocUrl:  othertmsDocs[j].DocUrl,
					})
				}
			}
		}

		otmsLoadDocRes = append(otmsLoadDocRes, dto.LoadDocsResult{
			LoadId:    LoadDocsReqDTO.OthertmsLoads[i],
			Documents: docArray,
		})
	}

	return otmsLoadDocRes, nil

}

//AP_PC_11
//GetBTMSDBLoadDocs function is declared in here inside this, the load document URL's for BTMS loads is destructured and the response is returned as LoadDocsResult DTO
func (r DefaultLoadService) GetBTMSDBLoadDocs(LoadDocsReqDTO dto.LoadDocumentsReqDTO) ([]dto.LoadDocsResult, *errs.AppErrorvalidation) {

	var imagingAPIResponse = make([]dto.ImagingAPIResDTO, 0)

	usedIn := `&used_in=BTMS`
	authorizationToken := "?auth_t="

	//jwt token for imaging API is generated in here
	var sharedKey = []byte("ABCD1234")
	claims := jwt.Map{
		"iss": "Sunteck CATS",
		"aud": "Sunteck Imaging API",
		"iat": time.Now(),
		"exp": time.Now().Add(time.Second * 3600),
	}
	token, err := jwt.Sign(jwt.HS256, sharedKey, claims, jwt.MaxAge(60*time.Minute))

	if err != nil {
		log.Printf("Generate token failure: %v", err)
		return nil, errs.ValidateResponse(nil, 500, err.Error())
	}
	tokenString := string(token)
	print(tokenString)
	authorizationToken += tokenString

	//imaging API call is done here for all the BTMS loadIds passed in request
	for i := 0; i < len(LoadDocsReqDTO.BtmsLoads); i++ {
		url := constant.IMAGING_END_POINT + `/api/v1/load/`
		url += LoadDocsReqDTO.BtmsLoads[i] + authorizationToken + usedIn
		payload := strings.NewReader(``)
		client := &http.Client{}
		req, err1 := http.NewRequest("GET", url, payload)

		if err1 != nil {
			return nil, errs.ValidateResponse(nil, 500, err1.Error())
		}

		res, clientErr := client.Do(req)

		if clientErr != nil {
			return nil, errs.ValidateResponse(nil, 500, clientErr.Error())
		}

		loadDocs, ioutilError := ioutil.ReadAll(res.Body)

		if ioutilError != nil {
			return nil, errs.ValidateResponse(nil, 500, ioutilError.Error())
		}
		var response dto.ImagingAPIResDTO

		json.Unmarshal(loadDocs, &response)

		imagingAPIResponse = append(imagingAPIResponse, response)

	}

	//AP_PC_11
	//Repository function is called in here
	btmsLoads, btmsDBError := r.repo.GetBTMSLoadDocsRepo(&LoadDocsReqDTO)

	if btmsDBError != nil {
		return nil, errs.ValidateResponse(nil, btmsDBError.Code, btmsDBError.Message)
	}

	var btmsLoadDocRes = make([]dto.LoadDocsResult, 0)

	proofOfDelivery := "1"
	billOfLading := "2"
	customerConfirmation := "3"
	lumperReceipt := "4"
	customerReceipt := "5"
	scaleWeight := "6"
	accessorialApproval := "7"

	for i := 0; i < len(imagingAPIResponse); i++ {

		var docArray = make([]dto.LoadDocuments, 0)
		bol := 0
		if len(imagingAPIResponse[i].Documents) != 0 {
			for j := 0; j < len(imagingAPIResponse[i].Documents); j++ {

				if imagingAPIResponse[i].Documents[j].TypeId != nil {

					//POD Document
					if *imagingAPIResponse[i].Documents[j].TypeId == proofOfDelivery {
						docArray = append(docArray, dto.LoadDocuments{
							Id:      imagingAPIResponse[i].Documents[j].Id,
							LoadId:  LoadDocsReqDTO.BtmsLoads[i],
							Path:    imagingAPIResponse[i].Documents[j].Path,
							DocName: "Proof of Delivery",
							TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
							DocUrl:  imagingAPIResponse[i].Documents[j].Body,
						})
					}

					//BOL Document
					if *imagingAPIResponse[i].Documents[j].TypeId == billOfLading {
						bol = 1
						docArray = append(docArray, dto.LoadDocuments{
							Id:      imagingAPIResponse[i].Documents[j].Id,
							LoadId:  LoadDocsReqDTO.BtmsLoads[i],
							Path:    imagingAPIResponse[i].Documents[j].Path,
							DocName: "Bill of Lading",
							TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
							DocUrl:  imagingAPIResponse[i].Documents[j].Body,
						})
					}
					if slices.Contains(imagingAPIResponse[i].Documents[j].Categories, "customer") {
						switch *imagingAPIResponse[i].Documents[j].TypeId {
						case customerConfirmation:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[i].Documents[j].Id,
								LoadId:  LoadDocsReqDTO.BtmsLoads[i],
								Path:    imagingAPIResponse[i].Documents[j].Path,
								DocName: "Customer Confirmation",
								TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
								DocUrl:  imagingAPIResponse[i].Documents[j].Body,
							})
						case lumperReceipt:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[i].Documents[j].Id,
								LoadId:  LoadDocsReqDTO.BtmsLoads[i],
								Path:    imagingAPIResponse[i].Documents[j].Path,
								DocName: "Lumper Receipt",
								TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
								DocUrl:  imagingAPIResponse[i].Documents[j].Body,
							})
						case customerReceipt:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[i].Documents[j].Id,
								LoadId:  LoadDocsReqDTO.BtmsLoads[i],
								Path:    imagingAPIResponse[i].Documents[j].Path,
								DocName: "Customer Receipt",
								TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
								DocUrl:  imagingAPIResponse[i].Documents[j].Body,
							})
						case scaleWeight:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[i].Documents[j].Id,
								LoadId:  LoadDocsReqDTO.BtmsLoads[i],
								Path:    imagingAPIResponse[i].Documents[j].Path,
								DocName: "Scale/Weight",
								TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
								DocUrl:  imagingAPIResponse[i].Documents[j].Body,
							})
						case accessorialApproval:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[i].Documents[j].Id,
								LoadId:  LoadDocsReqDTO.BtmsLoads[i],
								Path:    imagingAPIResponse[i].Documents[j].Path,
								DocName: "Accessorial Approval",
								TypeId:  *imagingAPIResponse[i].Documents[j].TypeId,
								DocUrl:  imagingAPIResponse[i].Documents[j].Body,
							})
						}
					}

				}
			}

		}

		// Invoice url, Shipping label and BOL URL are formed here
		for j := 0; j < len(btmsLoads); j++ {
			if LoadDocsReqDTO.BtmsLoads[i] == btmsLoads[j].LoadId {

				if btmsLoads[j].LoadStatus == "POSTED" || btmsLoads[j].LoadStatus == "INVOICED" {

					key := "7e83908c0c085bf3e9dd72ae3b537110945b421e"

					invoiceUrl := constant.BROKERAGE_URL + `/fats/invoice.php?to_print=`

					getSignature := "&doc_type=iv&sig="

					h := sha1.New()

					h.Write([]byte(key + btmsLoads[j].LoadId))

					getSignature += hex.EncodeToString(h.Sum(nil))

					invoiceUrl += btmsLoads[j].LoadId + getSignature

					docArray = append(docArray, dto.LoadDocuments{
						Id:      "",
						LoadId:  btmsLoads[j].LoadId,
						Path:    "",
						DocName: "Invoice",
						TypeId:  "",
						DocUrl:  invoiceUrl,
					})

				}

				if btmsLoads[j].LoadMethod == "ELTL" || btmsLoads[j].LoadMethod == "LTL" {

					baseUrl := constant.CUSTOMER_PORTAL_END_POINT + `/loads/`
					if bol == 0 {
						docArray = append(docArray, dto.LoadDocuments{
							Id:      "",
							LoadId:  btmsLoads[j].LoadId,
							Path:    "",
							DocName: "Bill of Lading",
							TypeId:  "2",
							DocUrl:  baseUrl + btmsLoads[j].LoadId + "/bol",
						})
					}

					docArray = append(docArray, dto.LoadDocuments{
						Id:      "",
						LoadId:  btmsLoads[j].LoadId,
						Path:    "",
						DocName: "Shipping label",
						TypeId:  "",
						DocUrl:  baseUrl + btmsLoads[j].LoadId + "/shipping_label",
					})

				}
				break
			}

		}

		btmsLoadDocRes = append(btmsLoadDocRes, dto.LoadDocsResult{
			LoadId:    LoadDocsReqDTO.BtmsLoads[i],
			Documents: docArray,
		})

	}

	return btmsLoadDocRes, nil

}

// AP_PC_04
//ValidateRequest1 function is declared in here inside this function all the request validations are handled
func ValidateDocRequest(loadDocRequest dto.LoadDocumentsReqDTO) []dto.ValidateResDTO {

	errorArray := make([]dto.ValidateResDTO, 0)
	re := regexp.MustCompile(`^\d+$`)

	if len(loadDocRequest.BtmsLoads) == 0 && len(loadDocRequest.OthertmsLoads) == 0 {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Request should have atleast one parameter"})
	}
	if len(loadDocRequest.BtmsLoads) != 0 {
		for i := 0; i < len(loadDocRequest.BtmsLoads); i++ {

			if !re.MatchString(loadDocRequest.BtmsLoads[i]) {
				errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD002", Message: "Invalid parameter, BTMS loadId"})
				break
			}
		}
	}

	return errorArray
}

//--GetLoadDocuments API ends here

func NewLoadService(repository domain.LoadRepository) DefaultLoadService {
	return DefaultLoadService{repository}
}
