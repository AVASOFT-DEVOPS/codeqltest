package main

import (
	"fmt"
	"golang/app"
	"golang/constant"
	"golang/domain"
	Service "golang/service"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// It creates a new instance of router, creates a new instance of repository, creates a new instance of service, creates a new instance of a handler, and
// then creates a new route
func main() {
	print("hellow")
	log.Println("Starting the application...")

	router := mux.NewRouter()

	loadRepository := domain.NewLoadRepositoryDb()
	loadDetailsRepository := domain.NewLoadDetailsRepositoryDb()
	bookRepository := domain.NewBookRepositoryDb()
	CreateRepository := domain.NewTLRepositoryDb()
	// AP_PC_01
	HL := app.LHandlers{Service.NewLoadService(loadRepository)}
	HLD := app.LDHandlers{Service.NewLoadDetailsService(loadDetailsRepository)}
	HB := app.BHandlers{Service.NewBookService(bookRepository)}
	HC := app.TLHandlers{Service.NewTLService(CreateRepository)}

	router.HandleFunc("/book/load/loadsearch", HL.GetLoadSearchResult).Methods(http.MethodPost)
	router.HandleFunc("/book/load/loaddocument", HL.GetLoadDocuments).Methods(http.MethodPost)
	router.HandleFunc("/book/load/loaddetails/{loadId}", HLD.GetLoadDetails).Methods(http.MethodGet)

	log.Println("Starting of application Main", time.Now())

	router.HandleFunc("/book/ltl/bookltl", HB.BookLTL).Methods(http.MethodPost)

	log.Println("ending of application Main", time.Now())

	router.HandleFunc("/book/load/creatl", HC.CreateTL).Methods(http.MethodPost)

	listenAddr := ":" + constant.PORT

	fmt.Println(constant.PORT, "port number")

	log.Printf("About to listen on %s. Go to https://127.0.0.1:%s", constant.PORT, constant.PORT)
	log.Fatal(http.ListenAndServe(listenAddr, router))

}
