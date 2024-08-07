package main

import (
	servicehttp "dataflow-api/pkg/http"
	"dataflow-api/pkg/repository/sale"
	sale2 "dataflow-api/pkg/service/sale"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	saleRepository := sale.NewInMemorySaleRepository()
	saleService := sale2.NewService(saleRepository)

	router := mux.NewRouter()
	servicehttp.NewSaleHandler(router, saleService)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
