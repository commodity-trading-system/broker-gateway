package main

import (
	"broker-gateway/queier"
	"broker-gateway/executor"
	"github.com/joho/godotenv"
	"fmt"
	"strconv"
	"os"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var q queier.Querier


func futures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q.Futures()
}

func consignationsByFirmId(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {

}

func main()  {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	mysqlPort,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	config := executor.DBConfig{
		Host: os.Getenv("MYSQL_HOST"),
		Port: int(mysqlPort),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName: os.Getenv("MYSQL_DB"),
		User: os.Getenv("MYSQL_USER"),
	}
	queryPort := os.Getenv("QUERIER_PORT")
	if queryPort == "" {
		queryPort = "5002"
	}
	port,_ := strconv.Atoi(queryPort)
	q = queier.NewQuerier(config)

	router := queier.NewRouter(q)
	router.Start(port)

}