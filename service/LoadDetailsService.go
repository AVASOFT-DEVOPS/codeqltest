package Service

import (
	"crypto/sha1"
	"encoding/hex"
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

	"github.com/kataras/jwt"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=../mock/service/mockLoadDetailsService.go -package=Service golang/service LoadDetailsService
type LoadDetailsService interface {
	GetLoadDetailsService(dto.GetLoadDetailsReqDTO) (*dto.LoadDetailsResDTO, *errs.AppErrorvalidation)
}

type DefaultLoadDetailsService struct {
	repo domain.LoadDetailsRepository
}

//AP_PC_03
// The below code is a GetLoadDetails API service layer which is responsible for calling the repository layer and
// returning the response to the handler.
func (r DefaultLoadDetailsService) GetLoadDetailsService(loadDetailsVar dto.GetLoadDetailsReqDTO) (*dto.LoadDetailsResDTO, *errs.AppErrorvalidation) {

	// AP_PC_04
	//ValidateLoadRequest function is called and the loadRequest is passed as a parameter
	ValidateResponseVar := ValidateLoadRequest(loadDetailsVar)

	if len(ValidateResponseVar) != 0 {
		return nil, errs.ValidateResponse(ValidateResponseVar, 400, "")
	}

	if strings.ToUpper(loadDetailsVar.LoadOrigin) == "OTHERTMS" {

		//AP_PC_14, AP_PC_17
		//GetOTMSLoadDetailsInfo function is called and the loadDetailsVar is passed as a parameter
		otmsLoadDetails, othertmsErr := r.GetOTMSLoadDetailsInfo(loadDetailsVar)

		if othertmsErr != nil {
			return nil, errs.ValidateResponse(nil, othertmsErr.Code, othertmsErr.Message)
		}

		return otmsLoadDetails, nil

	} else if strings.ToUpper(loadDetailsVar.LoadOrigin) == "BTMS" {

		//AP_PC_06, AP_PC_10
		//GetBTMSLoadDetailsInfo function is called and the loadDetailsVar is passed as a parameter
		btmsLoadDetails, btmsErr := r.GetBTMSLoadDetailsInfo(loadDetailsVar)

		if btmsErr != nil {
			return nil, errs.ValidateResponse(nil, btmsErr.Code, btmsErr.Message)
		}

		return btmsLoadDetails, nil

	}

	return nil, nil
}

//AP_PC_08, AP_PC_09
//GetOTMSLoadDetailsInfo function is declared in here inside this function repository function is called and the loadDetailsVar is sent as a parameter
func (r DefaultLoadDetailsService) GetOTMSLoadDetailsInfo(loadDetailsVar dto.GetLoadDetailsReqDTO) (*dto.LoadDetailsResDTO, *errs.AppErrorvalidation) {

	shipConsRefInfo, commoditiesInfo, othertmsDocs, loadErr := r.repo.GetOTMSLoadDetailsRepo(loadDetailsVar)

	if loadErr != nil {
		return nil, errs.ValidateResponse(nil, loadErr.Code, loadErr.Message)
	}

	if len(shipConsRefInfo) != 0 {

		// shipper, consignee and reference details are destructured in here
		if shipConsRefInfo[0].EarliestShipmentsDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].EarliestShipmentsDate, "T")
			*shipConsRefInfo[0].EarliestShipmentsDate = dateSplit[0]
		}
		if shipConsRefInfo[0].LatestShipmentsDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].LatestShipmentsDate, "T")
			*shipConsRefInfo[0].LatestShipmentsDate = dateSplit[0]
		}
		if shipConsRefInfo[0].ShipperDriverinDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ShipperDriverinDate, "T")
			*shipConsRefInfo[0].ShipperDriverinDate = dateSplit[0]
		}
		if shipConsRefInfo[0].ShipperDriverOutDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ShipperDriverOutDate, "T")
			*shipConsRefInfo[0].ShipperDriverOutDate = dateSplit[0]
		}

		if shipConsRefInfo[0].EarliestShipmentsTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].EarliestShipmentsTime, "T")
			*shipConsRefInfo[0].EarliestShipmentsTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}
		if shipConsRefInfo[0].LatestShipmentsTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].LatestShipmentsTime, "T")
			*shipConsRefInfo[0].LatestShipmentsTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}

		if shipConsRefInfo[0].ShipperDriverinTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ShipperDriverinTime, "T")
			*shipConsRefInfo[0].ShipperDriverinTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}
		if shipConsRefInfo[0].ShipperDriverOutTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ShipperDriverOutTime, "T")
			*shipConsRefInfo[0].ShipperDriverOutTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}

		if shipConsRefInfo[0].EarliestConsigneeDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].EarliestConsigneeDate, "T")
			*shipConsRefInfo[0].EarliestConsigneeDate = dateSplit[0]
		}
		if shipConsRefInfo[0].LatestConsigneeDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].LatestConsigneeDate, "T")
			*shipConsRefInfo[0].LatestConsigneeDate = dateSplit[0]
		}
		if shipConsRefInfo[0].ConsigneeDriverinDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ConsigneeDriverinDate, "T")
			*shipConsRefInfo[0].ConsigneeDriverinDate = dateSplit[0]
		}
		if shipConsRefInfo[0].ConsigneeDriverOutDate != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ConsigneeDriverOutDate, "T")
			*shipConsRefInfo[0].ConsigneeDriverOutDate = dateSplit[0]
		}

		if shipConsRefInfo[0].EarliestConsigneeTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].EarliestConsigneeTime, "T")
			*shipConsRefInfo[0].EarliestConsigneeTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}
		if shipConsRefInfo[0].LatestConsigneeTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].LatestConsigneeTime, "T")
			*shipConsRefInfo[0].LatestConsigneeTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}

		if shipConsRefInfo[0].ConsigneeDriverinTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ConsigneeDriverinTime, "T")
			*shipConsRefInfo[0].ConsigneeDriverinTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}
		if shipConsRefInfo[0].ConsigneeDriverOutTime != nil {
			dateSplit := strings.Split(*shipConsRefInfo[0].ConsigneeDriverOutTime, "T")
			*shipConsRefInfo[0].ConsigneeDriverOutTime = strings.ReplaceAll(dateSplit[1], "Z", "")
		}
		LoadShipConsRefRes := dto.LoadShipConsRefDTO{
			LoadId:                 shipConsRefInfo[0].LoadId,
			LoadStatus:             shipConsRefInfo[0].LoadStatus,
			LoadMethod:             shipConsRefInfo[0].LoadMethod,
			ShipperName:            shipConsRefInfo[0].ShipperName,
			ShipperAddressLine1:    shipConsRefInfo[0].ShipperAddressLine1,
			ShipperAddressLine2:    shipConsRefInfo[0].ShipperAddressLine2,
			ShipperCity:            shipConsRefInfo[0].ShipperCity,
			ShipperState:           shipConsRefInfo[0].ShipperState,
			ShipperZip:             shipConsRefInfo[0].ShipperZip,
			EarliestShipmentsDate:  shipConsRefInfo[0].EarliestShipmentsDate,
			EarliestShipmentsTime:  shipConsRefInfo[0].EarliestShipmentsTime,
			LatestShipmentsDate:    shipConsRefInfo[0].LatestShipmentsDate,
			LatestShipmentsTime:    shipConsRefInfo[0].LatestShipmentsTime,
			ShipperDriverinDate:    shipConsRefInfo[0].ShipperDriverinDate,
			ShipperDriverinTime:    shipConsRefInfo[0].ShipperDriverinTime,
			ShipperDriverOutDate:   shipConsRefInfo[0].ShipperDriverOutDate,
			ShipperDriverOutTime:   shipConsRefInfo[0].ShipperDriverOutTime,
			ConsigneeName:          shipConsRefInfo[0].ConsigneeName,
			ConsigneeAddressLine1:  shipConsRefInfo[0].ConsigneeAddressLine1,
			ConsigneeAddressLine2:  shipConsRefInfo[0].ConsigneeAddressLine2,
			ConsigneeCity:          shipConsRefInfo[0].ConsigneeCity,
			ConsigneeState:         shipConsRefInfo[0].ConsigneeState,
			ConsigneeZip:           shipConsRefInfo[0].ConsigneeZip,
			EarliestConsigneeDate:  shipConsRefInfo[0].EarliestConsigneeDate,
			EarliestConsigneeTime:  shipConsRefInfo[0].EarliestConsigneeTime,
			LatestConsigneeDate:    shipConsRefInfo[0].LatestConsigneeDate,
			LatestConsigneeTime:    shipConsRefInfo[0].LatestConsigneeTime,
			ConsigneeDriverinDate:  shipConsRefInfo[0].ConsigneeDriverinDate,
			ConsigneeDriverinTime:  shipConsRefInfo[0].ConsigneeDriverinTime,
			ConsigneeDriverOutDate: shipConsRefInfo[0].ConsigneeDriverOutDate,
			ConsigneeDriverOutTime: shipConsRefInfo[0].ConsigneeDriverOutTime,
			PoNumber:               shipConsRefInfo[0].PoNumber,
			ShipBlNumber:           shipConsRefInfo[0].ShipBlNumber,
			ProNumber:              shipConsRefInfo[0].ProNumber,
			ShipperNumber:          shipConsRefInfo[0].ShipperNumber,
			PickupNumber:           shipConsRefInfo[0].PickupNumber,
			DeliveryNumber:         shipConsRefInfo[0].DeliveryNumber,
			LoadOrigin:             shipConsRefInfo[0].LoadOrigin,
		}

		//commodityResponse DTO conversion is done here
		commodityResponse := make([]dto.LoadCommoditiesDTO, 0)
		for _, c := range commoditiesInfo {
			commodityResponse = append(commodityResponse, c.ToLoadCommodityDto())
		}

		//load document Urls are destructured in here
		proofOfDelivery := 1
		billOfLading := 2
		customerConfirmation := 3
		lumperReceipt := 4
		customerReceipt := 5
		scaleWeight := 6
		accessorialApproval := 7

		docArray := make([]dto.LoadDocuments, 0)

		for j := 0; j < len(othertmsDocs); j++ {
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

		loadDetailsResult := dto.LoadDetailsResDTO{
			LoadShipConsRef: LoadShipConsRefRes,
			LoadCommodities: commodityResponse,
			LoadDocuments:   docArray,
		}

		return &loadDetailsResult, nil
	}
	return nil, nil
}

//AP_PC_16
//GetBTMSLoadDetailsInfo function is declared in here inside this function repository function is called and the loadDetailsVar is sent as a parameter
func (r DefaultLoadDetailsService) GetBTMSLoadDetailsInfo(loadDetailsVar dto.GetLoadDetailsReqDTO) (*dto.LoadDetailsResDTO, *errs.AppErrorvalidation) {

	shipConsRefInfo, commoditiesInfo, eventsTl, locationUpdates, eventsELTL, loadErr := r.repo.GetBTMSLoadDetailsRepo(loadDetailsVar)

	if loadErr != nil {
		return nil, errs.ValidateResponse(nil, loadErr.Code, loadErr.Message)
	}

	if len(shipConsRefInfo) != 0 {

		// shipper, consignee and reference details are destructured in here
		var LoadShipConsRefRes dto.LoadShipConsRefDTO
		LoadShipConsRefRes = dto.LoadShipConsRefDTO{
			LoadId:                 shipConsRefInfo[0].LoadId,
			LoadStatus:             shipConsRefInfo[0].LoadStatus,
			LoadMethod:             shipConsRefInfo[0].LoadMethod,
			TrackingEnabled:        shipConsRefInfo[0].TrackingEnabled,
			ShipperName:            shipConsRefInfo[0].ShipperName,
			ShipperAddressLine1:    shipConsRefInfo[0].ShipperAddressLine1,
			ShipperAddressLine2:    shipConsRefInfo[0].ShipperAddressLine2,
			ShipperCity:            shipConsRefInfo[0].ShipperCity,
			ShipperState:           shipConsRefInfo[0].ShipperState,
			ShipperZip:             shipConsRefInfo[0].ShipperZip,
			EarliestShipmentsDate:  shipConsRefInfo[0].EarliestShipmentsDate,
			EarliestShipmentsTime:  shipConsRefInfo[0].EarliestShipmentsTime,
			LatestShipmentsDate:    shipConsRefInfo[0].LatestShipmentsDate,
			LatestShipmentsTime:    shipConsRefInfo[0].LatestShipmentsTime,
			ShipperDriverinDate:    shipConsRefInfo[0].ShipperDriverinDate,
			ShipperDriverinTime:    shipConsRefInfo[0].ShipperDriverinTime,
			ShipperDriverOutDate:   shipConsRefInfo[0].ShipperDriverOutDate,
			ShipperDriverOutTime:   shipConsRefInfo[0].ShipperDriverOutTime,
			ConsigneeName:          shipConsRefInfo[0].ConsigneeName,
			ConsigneeAddressLine1:  shipConsRefInfo[0].ConsigneeAddressLine1,
			ConsigneeAddressLine2:  shipConsRefInfo[0].ConsigneeAddressLine2,
			ConsigneeCity:          shipConsRefInfo[0].ConsigneeCity,
			ConsigneeState:         shipConsRefInfo[0].ConsigneeState,
			ConsigneeZip:           shipConsRefInfo[0].ConsigneeZip,
			EarliestConsigneeDate:  shipConsRefInfo[0].EarliestConsigneeDate,
			EarliestConsigneeTime:  shipConsRefInfo[0].EarliestConsigneeTime,
			LatestConsigneeDate:    shipConsRefInfo[0].LatestConsigneeDate,
			LatestConsigneeTime:    shipConsRefInfo[0].LatestConsigneeTime,
			ConsigneeDriverinDate:  shipConsRefInfo[0].ConsigneeDriverinDate,
			ConsigneeDriverinTime:  shipConsRefInfo[0].ConsigneeDriverinTime,
			ConsigneeDriverOutDate: shipConsRefInfo[0].ConsigneeDriverOutDate,
			ConsigneeDriverOutTime: shipConsRefInfo[0].ConsigneeDriverOutTime,
			PoNumber:               shipConsRefInfo[0].PoNumber,
			ShipBlNumber:           shipConsRefInfo[0].ShipBlNumber,
			ProNumber:              shipConsRefInfo[0].ProNumber,
			ShipperNumber:          shipConsRefInfo[0].ShipperNumber,
			PickupNumber:           shipConsRefInfo[0].PickupNumber,
			DeliveryNumber:         shipConsRefInfo[0].DeliveryNumber,
			LoadOrigin:             shipConsRefInfo[0].LoadOrigin,
		}

		commodityResponse := make([]dto.LoadCommoditiesDTO, 0)
		for _, c := range commoditiesInfo {
			commodityResponse = append(commodityResponse, c.ToLoadCommodityDto())
		}

		//Fetching BTMS Load Documents in here
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

		//imaging API call is done here for the BTMS loadId that passed in request
		url := constant.IMAGING_END_POINT + `/api/v1/load/`
		url += loadDetailsVar.LoadId + authorizationToken + usedIn
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

		proofOfDelivery := "1"
		billOfLading := "2"
		customerConfirmation := "3"
		lumperReceipt := "4"
		customerReceipt := "5"
		scaleWeight := "6"
		accessorialApproval := "7"

		var docArray = make([]dto.LoadDocuments, 0)
		// loadId, _ := strconv.Atoi(loadDetailsVar.LoadId)
		bol := 0
		fmt.Print(bol)
		if shipConsRefInfo[0].LoadStatus == "POSTED" || shipConsRefInfo[0].LoadStatus == "INVOICED" {

			key := constant.InvoiceKey

			invoiceUrl := constant.BROKERAGE_URL + `/fats/invoice.php?to_print=`

			getSignature := "&doc_type=iv&sig="

			h := sha1.New()

			h.Write([]byte(key + loadDetailsVar.LoadId))

			getSignature += hex.EncodeToString(h.Sum(nil))

			invoiceUrl += loadDetailsVar.LoadId + getSignature

			docArray = append(docArray, dto.LoadDocuments{
				Id:      "",
				LoadId:  loadDetailsVar.LoadId,
				Path:    "",
				DocName: "Invoice",
				TypeId:  "",
				DocUrl:  invoiceUrl,
			})

		}
		if len(imagingAPIResponse[0].Documents) != 0 {
			for i := 0; i < len(imagingAPIResponse[0].Documents); i++ {

				if imagingAPIResponse[0].Documents[i].TypeId != nil {

					//POD Document
					if *imagingAPIResponse[0].Documents[i].TypeId == proofOfDelivery {
						docArray = append(docArray, dto.LoadDocuments{
							Id:      imagingAPIResponse[0].Documents[i].Id,
							LoadId:  loadDetailsVar.LoadId,
							Path:    imagingAPIResponse[0].Documents[i].Path,
							DocName: "Proof of Delivery",
							TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
							DocUrl:  imagingAPIResponse[0].Documents[i].Body,
						})
					}

					//BOL Document
					if *imagingAPIResponse[0].Documents[i].TypeId == billOfLading {
						bol = 1
						docArray = append(docArray, dto.LoadDocuments{
							Id:      imagingAPIResponse[0].Documents[i].Id,
							LoadId:  loadDetailsVar.LoadId,
							Path:    imagingAPIResponse[0].Documents[i].Path,
							DocName: "Bill of Lading",
							TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
							DocUrl:  imagingAPIResponse[0].Documents[i].Body,
						})
					}
					if slices.Contains(imagingAPIResponse[0].Documents[i].Categories, "customer") {
						switch *imagingAPIResponse[0].Documents[i].TypeId {
						case customerConfirmation:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[0].Documents[i].Id,
								LoadId:  loadDetailsVar.LoadId,
								Path:    imagingAPIResponse[0].Documents[i].Path,
								DocName: "Customer Confirmation",
								TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
								DocUrl:  imagingAPIResponse[0].Documents[i].Body,
							})
						case lumperReceipt:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[0].Documents[i].Id,
								LoadId:  loadDetailsVar.LoadId,
								Path:    imagingAPIResponse[0].Documents[i].Path,
								DocName: "Lumper Receipt",
								TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
								DocUrl:  imagingAPIResponse[0].Documents[i].Body,
							})
						case customerReceipt:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[0].Documents[i].Id,
								LoadId:  loadDetailsVar.LoadId,
								Path:    imagingAPIResponse[0].Documents[i].Path,
								DocName: "Customer Receipt",
								TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
								DocUrl:  imagingAPIResponse[0].Documents[i].Body,
							})
						case scaleWeight:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[0].Documents[i].Id,
								LoadId:  loadDetailsVar.LoadId,
								Path:    imagingAPIResponse[0].Documents[i].Path,
								DocName: "Scale/Weight",
								TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
								DocUrl:  imagingAPIResponse[0].Documents[i].Body,
							})
						case accessorialApproval:
							docArray = append(docArray, dto.LoadDocuments{
								Id:      imagingAPIResponse[0].Documents[i].Id,
								LoadId:  loadDetailsVar.LoadId,
								Path:    imagingAPIResponse[0].Documents[i].Path,
								DocName: "Accessorial Approval",
								TypeId:  *imagingAPIResponse[0].Documents[i].TypeId,
								DocUrl:  imagingAPIResponse[0].Documents[i].Body,
							})
						}
					}

				}
			}

		}

		if shipConsRefInfo[0].LoadMethod == "ELTL" || shipConsRefInfo[0].LoadMethod == "LTL" {

			baseUrl := constant.CUSTOMER_PORTAL_END_POINT + `/loads/`
			if bol == 0 {
				docArray = append(docArray, dto.LoadDocuments{
					Id:      "",
					LoadId:  loadDetailsVar.LoadId,
					Path:    "",
					DocName: "Bill of Lading",
					TypeId:  "2",
					DocUrl:  baseUrl + loadDetailsVar.LoadId + "/bol",
				})
			}

			docArray = append(docArray, dto.LoadDocuments{
				Id:      "",
				LoadId:  loadDetailsVar.LoadId,
				Path:    "",
				DocName: "Shipping label",
				TypeId:  "",
				DocUrl:  baseUrl + loadDetailsVar.LoadId + "/shipping_label",
			})

		}
		eventTrackingRes := make([]dto.EventTrackingUpdates, 0)
		if eventsTl != nil {
			for _, c := range eventsTl {
				eventTrackingRes = append(eventTrackingRes, c.ToLoadEventsTlDto())
			}
		} else if eventsELTL != nil {
			for _, c := range eventsELTL {
				eventTrackingRes = append(eventTrackingRes, c.ToLoadEventsELTLDto())
			}
		}

		locationUpdatesRes := make([]dto.LocationBreadCrumbs, 0)
		if locationUpdates != nil {
			for _, c := range locationUpdates {
				locationUpdatesRes = append(locationUpdatesRes, c.ToLoadLocationDto())
			}
		}

		loadDetailsResult := dto.LoadDetailsResDTO{
			LoadShipConsRef: LoadShipConsRefRes,
			LoadCommodities: commodityResponse,
			LoadDocuments:   docArray,
			LoadEvents:      eventTrackingRes,
			LocationUpdates: locationUpdatesRes,
		}
		return &loadDetailsResult, nil
	}
	return nil, nil
}

// AP_PC_05
//ValidateLoadRequest function is declared in here inside this function all the request validations are handled
func ValidateLoadRequest(loadDetailsVar dto.GetLoadDetailsReqDTO) []dto.ValidateResDTO {

	errorArray := make([]dto.ValidateResDTO, 0)

	re := regexp.MustCompile(`^\d+$`)

	if loadDetailsVar.LoadOrigin != "" && strings.ToUpper(loadDetailsVar.LoadOrigin) == "BTMS" {
		if !re.MatchString(loadDetailsVar.LoadId) {
			errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD001", Message: "Invalid parameter,loadId"})
		}
	}
	if loadDetailsVar.LoadOrigin == "" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD002", Message: "Missing Required parameter,load origin"})
	} else if strings.ToUpper(loadDetailsVar.LoadOrigin) != "BTMS" && strings.ToUpper(loadDetailsVar.LoadOrigin) != "OTHERTMS" {
		errorArray = append(errorArray, dto.ValidateResDTO{Code: "GLD003", Message: `Load origin should be a 'BTMS' or 'OtherTMS'`})
	}
	return errorArray
}

func NewLoadDetailsService(repository domain.LoadDetailsRepository) DefaultLoadDetailsService {
	return DefaultLoadDetailsService{repository}
}
