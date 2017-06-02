package queier

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
)

type Router interface {
	Start(port int)
}

type router struct {
	q Querier
	http *httprouter.Router
}

func (rte *router)futures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	futures := rte.q.Futures()
	res, err := json.Marshal(futures)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (rte *router)ordersByFirmId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
	orders := rte.q.Orders(firmId)
	res, err := json.Marshal(orders)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (rte *router)ordersByFirmIdAndId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
	order := rte.q.OrderById(firmId,ps.ByName("id"))
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (rte *router) consignations(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
	rtn := rte.q.Consignations(firmId)
	res, err := json.Marshal(rtn)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (rte *router)consignationsByFirmIdAndId(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
	order := rte.q.ConsignationById(firmId,ps.ByName("id"))
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func NewRouter(q Querier) Router {
	r := httprouter.New()

	rte:= &router{
		q: q,
		http: r,
	}
	rte.register()
	return rte
}

func (rte *router) register()  {
	rte.http.GET("/futures",rte.futures)
	rte.http.GET("/firms/:firmId/orders",rte.ordersByFirmId)
	rte.http.GET("/firms/:firmId/orders/:id",rte.ordersByFirmIdAndId)
	rte.http.GET("/firms/:firmId/consignations",rte.ordersByFirmId)
	rte.http.GET("/firms/:firmId/consignations/:id",rte.ordersByFirmIdAndId)
}

func (r *router) Start(port int)  {
	http.ListenAndServe(":"+strconv.Itoa(port),r.http)
}


