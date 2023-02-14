package app

import (
	"bytes"
	"encoding/json"
	"golang/dto"
	"golang/errs"
	Service "golang/mock/service"
	Service1 "golang/service"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func Test_should_return_load_search_with_code_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockLoadService(ctrl)
	custmId := []int{8280, 20631, 20631, 20629, 20631, 20629, 10788, 10788, 26746}
	gcId := []string{"ACBFEFFC-C2B7-4EEF-846C-E53A053F6BE7"}
	officeCode := []string{"FL-CIRM", "NH-WTG"}

	DummyRequest := dto.LoadSearchReqDTO{

		CustomerId:       custmId,
		Gcid:             gcId,
		OfficeCode:       officeCode,
		LoadId:           "",
		ShipDateStart:    "",
		ShipDateEnd:      "",
		LoadStatus:       "",
		PoNumber:         "",
		CustomerblNumber: "",
		PickupNumber:     "W2098",
		DeliveryNumber:   "",
		SortOrder:        "ascending",
		SortColumn:       "destination",
		OtmsLoadRecords:  "",
		BtmsLoadRecords:  "",
	}

	Dummy := []dto.LoadSearchResDTO{
		{
			ShipDateStart:  "2019-02-11T00:00:00Z",
			ShipDateEnd:    "2019-02-11T00:00:00Z",
			LoadStatus:     "DELIVERED",
			LoadMethod:     "Unknown",
			LoadId:         "25229398",
			PickupNumber:   nil,
			DeliveryNumber: nil,
			CustTotal:      "0.0000",
			ShipperCity:    "Fort Payne",
			ConsigneeCity:  "ABBEVILLE",
			ShipperState:   "AL",
			ConsigneeState: "GA",
			ProNumber:      nil,
			PoNumber:       nil,
			ShipBlNumber:   nil,
			InvoiceDate:    nil,
			LastModified:   "2022-11-01T12:10:35.618509Z",
			LoadOrigin:     "otherTMS",
		},
	}

	mockResponse := dto.LoadSearchResultResDTO{
		LoadResponse: Dummy,
		LoadCount: struct {
			OthertmsLoadRecords int `json:"othertmsLoadRecords"`
			BtmsLoadRecords     int `json:"btmsLoadRecords"`
			TotalLoadsCount     int `json:"totalLoadsCount"`
		}{
			OthertmsLoadRecords: 10,
			BtmsLoadRecords:     125,
			TotalLoadsCount:     123456,
		},
	}
	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().GetLoadSearchResultService(&DummyRequest).Return(&mockResponse, nil)
	ch := LHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/api/loadSearch", ch.GetLoadSearchResult)
	request, _ := http.NewRequest(http.MethodPost, "/api/loadSearch", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing status code 200")
	}
}

func Test_should_return_load_search_code_500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockLoadService(ctrl)
	DummyRequest := dto.LoadSearchReqDTO{

		CustomerId:       nil,
		Gcid:             nil,
		OfficeCode:       nil,
		LoadId:           "",
		ShipDateStart:    "",
		ShipDateEnd:      "",
		LoadStatus:       "",
		PoNumber:         "",
		CustomerblNumber: "",
		PickupNumber:     "W2098",
		DeliveryNumber:   "",
		SortOrder:        "ascending",
		SortColumn:       "destination",
		OtmsLoadRecords:  "",
		BtmsLoadRecords:  "",
	}

	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().GetLoadSearchResultService(&DummyRequest).Return(nil, errs.ValidateResponse(nil, http.StatusInternalServerError, "Unexpected database error"))
	ch := LHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/api/loadSearch", ch.GetLoadSearchResult)
	request, _ := http.NewRequest(http.MethodPost, "/api/loadSearch", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing status code 500")
	}
}

func Test_should_return_load_search_with_code_400(t *testing.T) {

	DummyRequest := dto.LoadSearchReqDTO{
		CustomerId:       nil,
		Gcid:             nil,
		OfficeCode:       nil,
		LoadId:           "",
		ShipDateStart:    "",
		ShipDateEnd:      "",
		LoadStatus:       "",
		LoadMethod:       "1234",
		PoNumber:         "",
		CustomerblNumber: "",
		PickupNumber:     "",
		DeliveryNumber:   "",
		SortOrder:        "",
		SortColumn:       "destifhhnation",
		OtmsLoadRecords:  "1",
		BtmsLoadRecords:  "1",
	}

	appError := Service1.ValidateRequest(&DummyRequest)

	if len(appError) == 0 {
		t.Error("Request validation failed")
	}

}

func Test_should_return_load_search_result_with_code_400(t *testing.T) {

	DummyRequest := dto.LoadSearchReqDTO{
		Gcid:             nil,
		OfficeCode:       nil,
		LoadId:           "qwerty",
		ShipDateStart:    "qwertyt",
		ShipDateEnd:      "qwerytrew",
		LoadStatus:       "12345",
		PoNumber:         "",
		CustomerblNumber: "",
		PickupNumber:     "W2098",
		DeliveryNumber:   "",
		SortOrder:        "ascggending",
		SortColumn:       "destifhhnation",
		OtmsLoadRecords:  "sdfgh",
		BtmsLoadRecords:  "asdfgh",
	}

	appError := Service1.ValidateRequest(&DummyRequest)

	if len(appError) == 0 {
		t.Error("Request validation failed")
	}

}

//GetLoadDocuments unit Test case

func Test_should_return_load_documents_with_code_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockLoadService(ctrl)
	btmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}
	otmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}

	DummyRequest := dto.LoadDocumentsReqDTO{
		OthertmsLoads: otmsLoadIds,
		BtmsLoads:     btmsLoadIds,
	}

	loadDocuments := []dto.LoadDocuments{
		{
			Id:      "3553820",
			LoadId:  "8948589",
			Path:    "20180202102208-003.pdf",
			DocName: "Proof of Delivery",
			TypeId:  "1",
			DocUrl:  "https://sunteck-imaging-staging.s3.amazonaws.com/3553820?AWSAccessKeyId=AKIASKNF6RMTGALLYIOT&Expires=1669818924&Signature=G7hGwPvy1d%2FZZ8Ozs35bHq46zA0%3D",
		},
	}
	loadDocsResult := []dto.LoadDocsResult{
		{LoadId: "8948589",
			Documents: loadDocuments,
		},
	}
	mockResponse := dto.LoadDocumentsResDTO{
		LoadDocResult: loadDocsResult,
	}

	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().GetLoadDocumentsService(DummyRequest).Return(&mockResponse, nil)
	ch := LHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/load/loaddocument", ch.GetLoadDocuments)
	request, _ := http.NewRequest(http.MethodPost, "/load/loaddocument", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing status code 200")
	}
}

func Test_should_return_load_documents_code_500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := Service.NewMockLoadService(ctrl)
	btmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}
	otmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}

	DummyRequest := dto.LoadDocumentsReqDTO{
		OthertmsLoads: otmsLoadIds,
		BtmsLoads:     btmsLoadIds,
	}

	data, err := json.Marshal(DummyRequest)
	if err != nil {
		log.Fatal(err)
	}
	Dummyloadacc := *bytes.NewReader(data)
	mockService.EXPECT().GetLoadDocumentsService(DummyRequest).Return(nil, errs.ValidateResponse(nil, http.StatusInternalServerError, "Unexpected database error"))
	ch := LHandlers{mockService}
	router := mux.NewRouter()
	router.HandleFunc("/load/loaddocument", ch.GetLoadDocuments)
	request, _ := http.NewRequest(http.MethodPost, "/load/loaddocument", &Dummyloadacc)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing status code 500")
	}
}

func Test_should_return_load_documents_with_code_400(t *testing.T) {

	btmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}
	DummyRequest := dto.LoadDocumentsReqDTO{
		BtmsLoads: btmsLoadIds,
	}

	appError := Service1.ValidateDocRequest(DummyRequest)

	if len(appError) == 0 {
		t.Error("Request validation failed")
	}

}

func Test_should_return_load_document_with_code_400(t *testing.T) {

	otmsLoadIds := []string{"8948589", "11532004", "8904929", "11761695"}

	DummyRequest := dto.LoadDocumentsReqDTO{
		OthertmsLoads: otmsLoadIds,
	}

	appError := Service1.ValidateDocRequest(DummyRequest)

	if len(appError) == 0 {
		t.Error("Request validation failed")
	}

}
