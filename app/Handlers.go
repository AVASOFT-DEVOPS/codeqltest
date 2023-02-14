package app

import (
	"encoding/json"
	"fmt"
	"golang/dto"
	"golang/errs"
	Service "golang/service"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type LHandlers struct {
	LService Service.LoadService
}

type LDHandlers struct {
	LDService Service.LoadDetailsService
}

type BHandlers struct {
	BService Service.BookService
}

type TLHandlers struct {
	TLService Service.TLService
}

type Success struct {
	Success bool
	ErrorID int
	Message string
}

// AP_PC_02
// A handler function is declared below inside this functions the service layer is called and Request body is sent as a parameter.
// The Response returned from the service layer is readed in here
func (h *LHandlers) GetLoadSearchResult(w http.ResponseWriter, r *http.Request) {

	var searchResultVar *dto.LoadSearchReqDTO

	err := json.NewDecoder(r.Body).Decode(&searchResultVar)

	// fmt.Printf("request%+v\n", searchResultVar)

	if err != nil {
		errormessage := errs.ErrorResponse{
			Errors: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "GLD001",
				Message: strings.ReplaceAll(err.Error(), "json: ", ""),
			},
		}
		writeResponse(w, http.StatusBadRequest, errormessage)
	} else {
		loadSearchResultServiceResVar, errors := h.LService.GetLoadSearchResultService(searchResultVar)
		//AP_PC_14
		if errors != nil {
			if errors.Code != 500 {
				writeResponse(w, errors.Code, errors.Errors)
			} else {
				errormessage := errs.ErrorResponse{
					Errors: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "500",
						Message: errors.Message,
					},
				}
				writeResponse(w, errors.Code, errormessage)
			}
		} else {
			if len(loadSearchResultServiceResVar.LoadResponse) == 0 {
				NoContent(w)
			} else {
				writeResponse(w, http.StatusOK, loadSearchResultServiceResVar)
			}
		}
	}
}

//AP_PC_02
// A handler function is declared below inside this functions the service layer is called and Request body is sent as a parameter.
// The Response returned from the service layer is readed in here
func (h *LHandlers) GetLoadDocuments(w http.ResponseWriter, r *http.Request) {

	var loadDocumentsVar dto.LoadDocumentsReqDTO

	err := json.NewDecoder(r.Body).Decode(&loadDocumentsVar)
	//fmt.Printf("loadDocumentsVar : %+v\n", loadDocumentsVar)

	if err != nil {
		errormessage := errs.ErrorResponse{
			Errors: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "GLD001",
				Message: strings.ReplaceAll(err.Error(), "json: ", ""),
			},
		}
		writeResponse(w, http.StatusBadRequest, errormessage)
	} else {
		loadDocumentsVarServiceResVar, errors := h.LService.GetLoadDocumentsService(loadDocumentsVar)

		if errors != nil {

			if errors.Code != 500 {

				writeResponse(w, errors.Code, errors.Errors)
			} else {
				errormessage := errs.ErrorResponse{
					Errors: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "500",
						Message: errors.Message,
					},
				}
				writeResponse(w, errors.Code, errormessage)
			}
		} else {
			//if len(loadDocumentsVarServiceResVar) == 0 {
			//	NoContent(w)
			//} else {
			writeResponse(w, http.StatusOK, loadDocumentsVarServiceResVar)
			//}
		}
	}
}

//AP_PC_02
// A handler function is declared below inside this functions the service layer is called and Request body is sent as a parameter.
// The Response returned from the service layer is readed in here
func (h *LDHandlers) GetLoadDetails(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	loadId := vars["loadId"]
	loadOrigin := r.FormValue("loadOrigin")

	var loadDetailsVar dto.GetLoadDetailsReqDTO

	loadDetailsVar.LoadId = loadId
	loadDetailsVar.LoadOrigin = loadOrigin

	loadDetailsResVar, errors := h.LDService.GetLoadDetailsService(loadDetailsVar)
	if errors != nil {
		if errors.Code != 500 {
			writeResponse(w, errors.Code, errors.Errors)
		} else {
			errormessage := errs.ErrorResponse{
				Errors: struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				}{
					Code:    "500",
					Message: errors.Message,
				},
			}
			writeResponse(w, errors.Code, errormessage)
		}
	} else {
		if loadDetailsResVar == nil {
			NoContent(w)
		} else {
			writeResponse(w, http.StatusOK, loadDetailsResVar)
		}
	}
}

// BGLD_PS_1.2
// A handler function is declared below inside this functions the service layer is called and Request body is sent as a parameter.
// The Response returned from the service layer is readed in here of bookLTLService

func (h *BHandlers) BookLTL(w http.ResponseWriter, req *http.Request) {
	log.Println("Starting of handler", time.Now())
	var BookReq *dto.BookLTLRequestDTO
	errReq := json.NewDecoder(req.Body).Decode(&BookReq)
	fmt.Println(errReq)
	userId := req.Header.Get("userId")
	loginId := req.Header.Get("loginId")

	BookReq.UserId = userId
	BookReq.SunteckLoginId = loginId
	if errReq != nil {
		errormessage := errs.ErrorResponse{
			Errors: struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}{
				Code:    "BLTL001",
				Message: errReq.Error(),
			},
		}
		writeResponse(w, http.StatusBadRequest, errormessage)
		fmt.Println(errormessage, "errrmsg")
		fmt.Println(http.StatusBadRequest, "statttttt")
	} else {
		bookResponse, errors := h.BService.BookLTLService(BookReq)

		if errors != nil {
			if errors.Message == "Unauthorised User" {
				errormessage := errs.ErrorResponse{
					Errors: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "BL001",
						Message: errors.Message,
					},
				}
				writeResponse(w, errors.Code, errormessage)
			} else if errors.Code != 500 {
				fmt.Println("errors.Code", errors.Code, errors.Errors)
				writeResponse(w, errors.Code, errors.Errors)
			} else {
				errormessage := errs.ErrorResponse{
					Errors: struct {
						Code    string `json:"code"`
						Message string `json:"message"`
					}{
						Code:    "500",
						Message: errors.Message,
					},
				}
				writeResponse(w, errors.Code, errormessage)
			}

		} else {
			//if len(loadDocumentsVarServiceResVar) == 0 {
			//	NoContent(w)
			//} else {
			log.Println("Ending of handler", time.Now())
			writeResponse(w, http.StatusOK, bookResponse)
			//}
		}
	}
}


//ps_1.2.1 To create a TL to calling createTlService is defined here
//ps_1.2.3
func (h *TLHandlers) CreateTL(w http.ResponseWriter, r *http.Request) {

	// /ps_1.3.1 a varible is declared to store the request
	var TlRequest dto.CreateTLReq
	//ps_1.3.2 Using the json Decoder the request is assigned to the variable created
	err := json.NewDecoder(r.Body).Decode(&TlRequest)

	reqToken := r.Header.Get("Authorization")
	log.Println(" token authorization : ", reqToken)
	splitToken := strings.Split(reqToken, "Bearer ")
	log.Println(" splitToken authorization : ", splitToken)
	Token := splitToken[1]
	log.Println(" token authorization after split : ", Token)
	TlRequest.Token=Token;
	fmt.Println(TlRequest, "request\n")
	if err != nil {
		fmt.Println(err, "request error")
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		//ps_1.3.3 Call to the createTLService is Made
		TLRes, error := h.TLService.CreateTLService(TlRequest)
		// /ps_1.3.4 A if statement is used to check for the error
		if error != nil {
			//ps_1.3.5 to check for the error response
			//ps_1.10.1 error response is returned here
			writeResponse(w, error.Code, error.AsUserMessage())
		} else {
			writeResponse(w, http.StatusOK, TLRes.AsUserRespopnse())
		}
		//fmt.Println(error)
		//writeResponse(w, http.StatusOK, user)
	}
}
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func NoContent(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusNoContent)
	// send the headers with a 204 response code.
}
