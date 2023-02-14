package app

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"golang/dto"
	"golang/errs"
	Service "golang/mock/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

//func Test_should_return_load_details_with_code_200(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockService := Service.NewMockLoadDetailsService(ctrl)
//	// custmId := []int{8280, 20631, 20631, 20629, 20631, 20629, 10788, 10788, 26746}
//	// gcId := []string{"ACBFEFFC-C2B7-4EEF-846C-E53A053F6BE7"}
//	// officeCode := []string{"FL-CIRM", "NH-WTG"}
//
//	DummyRequest := dto.GetLoadDetailsReqDTO{
//		LoadId:     "5000002",
//		LoadOrigin: "othertms",
//	}
//
//	LoadCommodities := []dto.LoadCommoditiesDTO{
//		{
//			ItemId:               0,
//			Hazmat:               nil,
//			ItemQuantity:         nil,
//			UnitOfMeasure:        nil,
//			Weight:               nil,
//			ItemDescription:      nil,
//			ItemValue:            nil,
//			Class:                nil,
//			Nmfc:                 nil,
//			PalletLength:         nil,
//			PalletWidth:          nil,
//			PalletHeight:         nil,
//			ItemDensity:          nil,
//			HazmatContact:        nil,
//			HazmatGroupName:      nil,
//			HazmatPackagingGroup: nil,
//			HazmatUNNAnumber:     nil,
//			HazmatClass:          nil,
//			HazmatPlacard:        nil,
//			HazmatFlashTemp:      nil,
//			HazmatFlashType:      nil,
//			HazmatUom:            nil,
//			HazmatCertHolderName: nil,
//			HazmatContactName:    nil,
//			HazmatPhoneNumber:    nil,
//		},
//	}
//
//	mockResponse := dto.LoadDetailsResDTO{
//		LoadShipConsRef: dto.LoadShipConsRefDTO{
//			LoadId:                 10975709,
//			LoadStatus:             "DELIVERED FINAL",
//			LoadMethod:             "TL",
//			TrackingEnabled:        nil,
//			ShipperName:            "TIRECORD USA",
//			ShipperAddressLine1:    "2011 RANDOLPH RD",
//			ShipperAddressLine2:    nil,
//			ShipperCity:            "SHELBY",
//			ShipperState:           "NC",
//			ShipperZip:             "28150",
//			EarliestShipmentsDate:  nil,
//			EarliestShipmentsTime:  nil,
//			LatestShipmentsDate:    nil,
//			LatestShipmentsTime:    nil,
//			ShipperDriverinDate:    nil,
//			ShipperDriverOutDate:   nil,
//			ShipperDriverinTime:    nil,
//			ShipperDriverOutTime:   nil,
//			ConsigneeName:          "AMERICAN YARNS BURLINGTON",
//			ConsigneeAddressLine1:  "1305 GRAHAM ST",
//			ConsigneeAddressLine2:  nil,
//			ConsigneeCity:          "BURLINGTON",
//			ConsigneeState:         "NC",
//			ConsigneeZip:           "27217",
//			EarliestConsigneeDate:  nil,
//			EarliestConsigneeTime:  nil,
//			LatestConsigneeDate:    nil,
//			LatestConsigneeTime:    nil,
//			ConsigneeDriverinDate:  nil,
//			ConsigneeDriverOutDate: nil,
//			ConsigneeDriverinTime:  nil,
//			ConsigneeDriverOutTime: nil,
//			PoNumber:               nil,
//			ShipBlNumber:           nil,
//			ProNumber:              nil,
//			ShipperNumber:          nil,
//			PickupNumber:           nil,
//			DeliveryNumber:         nil,
//			LoadOrigin:             "BTMS",
//		},
//		LoadCommodities: LoadCommodities,
//		LoadDocuments:   nil,
//		LoadEvents:      nil,
//		LocationUpdates: nil,
//	}
//
//	//data, err := json.Marshal(DummyRequest)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//Dummyloadacc := *bytes.NewReader(data)
//	mockService.EXPECT().GetLoadDetailsService(&DummyRequest).Return(&mockResponse, nil)
//	ch := LDHandlers{mockService}
//	router := mux.NewRouter()
//	router.HandleFunc("/book/load/loaddetails/5000002?loadOrigin=othertms", ch.GetLoadDetails)
//	//request, _ := http.NewRequest(http.MethodPost, "/book/load/loaddetails/{loadId}", &Dummyloadacc)
//	request, _ := http.NewRequest(http.MethodGet, "book/load/loaddetails/5000002?loadOrigin=othertms", nil)
//
//	recorder := httptest.NewRecorder()
//	router.ServeHTTP(recorder, request)
//	if recorder.Code != http.StatusOK {
//		t.Error("Failed while testing status code 200")
//	}
//}
//
//func Test_should_return_load_details_with_code_200(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockService := Service.NewMockLoadDetailsService(ctrl)
//	LoadCommodities := []dto.LoadCommoditiesDTO{
//		{
//			ItemId:               0,
//			Hazmat:               nil,
//			ItemQuantity:         nil,
//			UnitOfMeasure:        nil,
//			Weight:               nil,
//			ItemDescription:      nil,
//			ItemValue:            nil,
//			Class:                nil,
//			Nmfc:                 nil,
//			PalletLength:         nil,
//			PalletWidth:          nil,
//			PalletHeight:         nil,
//			ItemDensity:          nil,
//			HazmatContact:        nil,
//			HazmatGroupName:      nil,
//			HazmatPackagingGroup: nil,
//			HazmatUNNAnumber:     nil,
//			HazmatClass:          nil,
//			HazmatPlacard:        nil,
//			HazmatFlashTemp:      nil,
//			HazmatFlashType:      nil,
//			HazmatUom:            nil,
//			HazmatCertHolderName: nil,
//			HazmatContactName:    nil,
//			HazmatPhoneNumber:    nil,
//		},
//	}
//	DummyResponse := dto.LoadDetailsResDTO{
//		LoadShipConsRef: dto.LoadShipConsRefDTO{
//			LoadId:                 10975709,
//			LoadStatus:             "DELIVERED FINAL",
//			LoadMethod:             "TL",
//			TrackingEnabled:        nil,
//			ShipperName:            "TIRECORD USA",
//			ShipperAddressLine1:    "2011 RANDOLPH RD",
//			ShipperAddressLine2:    nil,
//			ShipperCity:            "SHELBY",
//			ShipperState:           "NC",
//			ShipperZip:             "28150",
//			EarliestShipmentsDate:  nil,
//			EarliestShipmentsTime:  nil,
//			LatestShipmentsDate:    nil,
//			LatestShipmentsTime:    nil,
//			ShipperDriverinDate:    nil,
//			ShipperDriverOutDate:   nil,
//			ShipperDriverinTime:    nil,
//			ShipperDriverOutTime:   nil,
//			ConsigneeName:          "AMERICAN YARNS BURLINGTON",
//			ConsigneeAddressLine1:  "1305 GRAHAM ST",
//			ConsigneeAddressLine2:  nil,
//			ConsigneeCity:          "BURLINGTON",
//			ConsigneeState:         "NC",
//			ConsigneeZip:           "27217",
//			EarliestConsigneeDate:  nil,
//			EarliestConsigneeTime:  nil,
//			LatestConsigneeDate:    nil,
//			LatestConsigneeTime:    nil,
//			ConsigneeDriverinDate:  nil,
//			ConsigneeDriverOutDate: nil,
//			ConsigneeDriverinTime:  nil,
//			ConsigneeDriverOutTime: nil,
//			PoNumber:               nil,
//			ShipBlNumber:           nil,
//			ProNumber:              nil,
//			ShipperNumber:          nil,
//			PickupNumber:           nil,
//			DeliveryNumber:         nil,
//			LoadOrigin:             "BTMS",
//		},
//		LoadCommodities: LoadCommodities,
//		LoadDocuments:   nil,
//		LoadEvents:      nil,
//		LocationUpdates: nil,
//	}
//	var dummy dto.GetLoadDetailsReqDTO
//	mockService.EXPECT().GetLoadDetailsService(dummy).Return(&DummyResponse, nil)
//	ch := LDHandlers{mockService}
//	router := mux.NewRouter()
//	router.HandleFunc("/book/load/loaddetails/11762432", ch.GetLoadDetails)
//	request, _ := http.NewRequest(http.MethodGet, "/book/load/loaddetails/11762432", nil)
//	recorder := httptest.NewRecorder()
//	router.ServeHTTP(recorder, request)
//	if recorder.Code != http.StatusOK {
//		t.Error("Failed while testing status code 400")
//	}
//}
//
//func Test_should_return_load_details_code_500(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockService := Service.NewMockLoadDetailsService(ctrl)
//	var DummyRequest dto.GetLoadDetailsReqDTO
//
//	data, err := json.Marshal(DummyRequest)
//	if err != nil {
//		log.Fatal(err)
//	}
//	Dummyloadacc := *bytes.NewReader(data)
//	mockService.EXPECT().GetLoadDetailsService(DummyRequest).Return(nil, errs.ValidateResponse(nil, http.StatusInternalServerError, "Unexpected database error"))
//	ch := LDHandlers{mockService}
//	router := mux.NewRouter()
//	router.HandleFunc("/book/load/loaddetails/11762432", ch.GetLoadDetails)
//	request, _ := http.NewRequest(http.MethodPost, "/book/load/loaddetails/11762432", &Dummyloadacc)
//	recorder := httptest.NewRecorder()
//	router.ServeHTTP(recorder, request)
//	if recorder.Code != http.StatusInternalServerError {
//		t.Error("Failed while testing status code 500")
//	}
//}
//
//func Test_should_return_load_details_with_code_400(t *testing.T) {
//
//	DummyRequest := dto.GetLoadDetailsReqDTO{
//		LoadId:     "5QDFF000002",
//		LoadOrigin: "otherGFDStms",
//	}
//
//	appError := Service1.ValidateLoadRequest(DummyRequest)
//
//	if len(appError) == 0 {
//		t.Error("Request validation failed")
//	}
//
//}
//
//func Test_should_return_load_details_with1_code_400(t *testing.T) {
//
//	DummyRequest := dto.GetLoadDetailsReqDTO{
//		LoadId: "5000002",
//	}
//
//	appError := Service1.ValidateLoadRequest(DummyRequest)
//
//	if len(appError) == 0 {
//		t.Error("Request validation failed")
//	}
//
//}

func Test_should_return_bookLtl_with_code_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockBookService(ctrl)
	DummyRequest := dto.BookLTLRequestDTO{
		UserId:         "",
		SunteckLoginId: "",
		CustmId:        "",
		CustaId:        "",
		LoadId:         0,
		QuoteId:        0,
		
		ShipperDetails: dto.ShipperAddress{
			ShipName:         "",
			ShipAddress1:     "",
			ShipAddress2:     "",
			ShipCity:         "",
			ShipState:        "",
			ShipZipCode:      0,
			ShipCountry:      "",
			ShipContactName:  "",
			ShipEmail:        "",
			ShipPhone:        "",
			ShipFax:          "",
			ShipLoadNotes:    "",
			ShipEarliestDate: "",
			ShipLatestDate:   "",
			ShipEarliestTime: "",
			ShipLatestTime:   "",
		},
		ConsigneeDetails: dto.ConsigneeAddress{
			ConsName:         "",
			ConsAddress1:     "",
			ConsAddress2:     "",
			ConsCity:         "",
			ConsState:        "",
			ConsZipCode:      0,
			ConsCountry:      "",
			ConsContactName:  "",
			ConsEmail:        "",
			ConsPhone:        "",
			ConsFax:          "",
			ConsLoadNotes:    "",
			ConsEarliestDate: "",
			ConsLatestDate:   "",
			ConsEarliestTime: "",
			ConsLatestTime:   "",
		},
		Commodities: nil,
	}

	DummyResponse := dto.BookLTLResponseDto{
		LoadNumber:  0,
		QuoteId:     0,
		QuoteNumber: "",
		AgentEmail:  "",
		PriceDetails: dto.PriceDetails{
			Scac:               "",
			Service:            "",
			CarrierName:        "",
			CarrierNotes:       "",
			TransitTime:        "",
			FlatPrice:          "",
			FuelSurchargePrice: "",
		},
		TotalPrice: "",
	}
	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().BookLTLService(&DummyRequest).Return(&DummyResponse, nil)
	ch := BHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/book/ltl/bookltl", ch.BookLTL)
	request, _ := http.NewRequest(http.MethodPost, "/book/ltl/bookltl", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing status code 200")
	}
}

func Test_should_return_bookLtl_with_code_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockBookService(ctrl)
	DummyRequest := dto.BookLTLRequestDTO{
		UserId:         "",
		SunteckLoginId: "",
		CustmId:        "",
		CustaId:        "",
		LoadId:         0,
		QuoteId:        0,
		
		ShipperDetails: dto.ShipperAddress{
			ShipName:         "",
			ShipAddress1:     "",
			ShipAddress2:     "",
			ShipCity:         "",
			ShipState:        "",
			ShipZipCode:      0,
			ShipCountry:      "",
			ShipContactName:  "",
			ShipEmail:        "",
			ShipPhone:        "",
			ShipFax:          "",
			ShipLoadNotes:    "",
			ShipEarliestDate: "",
			ShipLatestDate:   "",
			ShipEarliestTime: "",
			ShipLatestTime:   "",
		},
		ConsigneeDetails: dto.ConsigneeAddress{
			ConsName:         "",
			ConsAddress1:     "",
			ConsAddress2:     "",
			ConsCity:         "",
			ConsState:        "",
			ConsZipCode:      0,
			ConsCountry:      "",
			ConsContactName:  "",
			ConsEmail:        "",
			ConsPhone:        "",
			ConsFax:          "",
			ConsLoadNotes:    "",
			ConsEarliestDate: "",
			ConsLatestDate:   "",
			ConsEarliestTime: "",
			ConsLatestTime:   "",
		},
		Commodities: nil,
	}

	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().BookLTLService(&DummyRequest).Return(nil, errs.ValidateResponse(nil, http.StatusBadRequest, ""))
	ch := BHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/book/ltl/bookltl", ch.BookLTL)
	request, _ := http.NewRequest(http.MethodPost, "/book/ltl/bookltl", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Error("Failed while testing status code 400")
	}
}

func Test_should_return_bookLtl_with_code_500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockBookService(ctrl)
	DummyRequest := dto.BookLTLRequestDTO{
		UserId:         "",
		SunteckLoginId: "",
		CustmId:        "",
		CustaId:        "",
		LoadId:         0,
		QuoteId:        0,
		
		ShipperDetails: dto.ShipperAddress{
			ShipName:         "",
			ShipAddress1:     "",
			ShipAddress2:     "",
			ShipCity:         "",
			ShipState:        "",
			ShipZipCode:      0,
			ShipCountry:      "",
			ShipContactName:  "",
			ShipEmail:        "",
			ShipPhone:        "",
			ShipFax:          "",
			ShipLoadNotes:    "",
			ShipEarliestDate: "",
			ShipLatestDate:   "",
			ShipEarliestTime: "",
			ShipLatestTime:   "",
		},
		ConsigneeDetails: dto.ConsigneeAddress{
			ConsName:         "",
			ConsAddress1:     "",
			ConsAddress2:     "",
			ConsCity:         "",
			ConsState:        "",
			ConsZipCode:      0,
			ConsCountry:      "",
			ConsContactName:  "",
			ConsEmail:        "",
			ConsPhone:        "",
			ConsFax:          "",
			ConsLoadNotes:    "",
			ConsEarliestDate: "",
			ConsLatestDate:   "",
			ConsEarliestTime: "",
			ConsLatestTime:   "",
		},
		Commodities: nil,
	}

	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().BookLTLService(&DummyRequest).Return(nil, errs.ValidateResponse(nil, http.StatusInternalServerError, ""))
	ch := BHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/book/ltl/bookltl", ch.BookLTL)
	request, _ := http.NewRequest(http.MethodPost, "/book/ltl/bookltl", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing status code 400")
	}
}
