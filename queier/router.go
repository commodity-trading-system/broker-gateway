package queier

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"fmt"
	"strconv"
	"broker-gateway/entities"
	"github.com/satori/go.uuid"
)

type Router interface {
	Start(port int)
}

type router struct {
	q Querier
	http *httprouter.Router
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



func (rte *router)futures(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	futures := rte.q.Futures()
	echo(w, futures)
}
//
//func (rte *router)ordersByFirmId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
//	orders := rte.q.Orders(firmId)
//	res, err := json.Marshal(orders)
//	if err != nil {
//		fmt.Println(err)
//	}
//	w.Write(res)
//}
//
//func (rte *router)ordersByFirmIdAndId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
//	order := rte.q.OrderById(firmId,ps.ByName("id"))
//	res, err := json.Marshal(order)
//	if err != nil {
//		fmt.Println(err)
//	}
//	w.Write(res)
//}
//
//func (rte *router) consignations(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
//	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
//	rtn := rte.q.Consignations(firmId)
//	res, err := json.Marshal(rtn)
//	if err != nil {
//		fmt.Println(err)
//	}
//	w.Write(res)
//}
//
//func (rte *router)consignationsByFirmIdAndId(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
//	firmId,_ := strconv.Atoi(ps.ByName("firmId"))
//	order := rte.q.ConsignationById(firmId,ps.ByName("id"))
//	res, err := json.Marshal(order)
//	if err != nil {
//		fmt.Println(err)
//	}
//	w.Write(res)
//}
//
//func (rte *router)quotationsByFutureId(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
//	futureId,_ := strconv.Atoi(ps.ByName("id"))
//	rtn := rte.q.Quotations(futureId)
//	res, err := json.Marshal(rtn)
//	if err != nil {
//		fmt.Println(err)
//	}
//	w.Write(res)
//}

func echo(w http.ResponseWriter, entity interface{})  {
	if entity == nil {
		w.Write([]byte("{}"))
		return
	}
	res, err := json.Marshal(entity)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Write(res)
}


func (rte *router) register()  {
	
	// For firm
	rte.http.GET("/futures",func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		futures := rte.q.Futures()
		echo(w, futures)
	})
	rte.http.GET("/futures/:id", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		id,_ := strconv.Atoi(p.ByName("id"))
		limit,_ := strconv.Atoi(h.URL.Query().Get("limit"))
		offset,_ := strconv.Atoi(h.URL.Query().Get("offset"))
		echo(w,rte.q.Quotations(id,limit,offset))
	})
	rte.http.GET("/firms/:firmId/orders", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		id,_ := strconv.Atoi(p.ByName("firmId"))
		limit,_ := strconv.Atoi(h.URL.Query().Get("limit"))
		offset,_ := strconv.Atoi(h.URL.Query().Get("offset"))
		res := rte.q.Orders(id,limit,offset)
		for i:=0;i<len(res) ;i++  {
			res[i] = entities.TransformForFirm(res[i],id)
		}
		echo(w,res)
	})
	rte.http.GET("/firms/:firmId/orders/:id",func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		firmid,_ := strconv.Atoi(p.ByName("firmId"))
		res := rte.q.OrderById(firmid, p.ByName("id"))
		if uuid.Equal(res.ID,uuid.FromBytesOrNil(nil)) {
			echo(w, nil)
		} else {
			echo(w,res)
		}
	})
	rte.http.GET("/firms/:firmId/consignations",func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		id,_ := strconv.Atoi(p.ByName("firmId"))
		limit,_ := strconv.Atoi(h.URL.Query().Get("limit"))
		offset,_ := strconv.Atoi(h.URL.Query().Get("offset"))
		res := rte.q.Consignations(id,limit,offset)
		echo(w,res)
	})
	rte.http.GET("/firms/:firmId/consignations/:id",func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		firmid,_ := strconv.Atoi(p.ByName("firmId"))
		res := rte.q.ConsignationById(firmid, p.ByName("id"))
		if uuid.Equal(res.ID,uuid.FromBytesOrNil(nil)) {
			echo(w, nil)
		} else {
			echo(w,res)
		}
	})
	
	// For admin
	rte.http.GET("/admin/orders", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		limit,_ := strconv.Atoi(h.URL.Query().Get("limit"))
		offset,_ := strconv.Atoi(h.URL.Query().Get("offset"))
		echo(w,rte.q.AllOrders(limit,offset))
	})

	rte.http.GET("/admin/consignations", func(w http.ResponseWriter, h *http.Request, p httprouter.Params)  {
		limit,_ := strconv.Atoi(h.URL.Query().Get("limit"))
		offset,_ := strconv.Atoi(h.URL.Query().Get("offset"))
		echo(w,rte.q.AllConsignations(limit,offset))
	})
	
	
	rte.http.PUT("/admin/commissions/:id", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		firmId,_ :=strconv.Atoi(h.URL.Query().Get("firmId"))
		futureId,_ := strconv.Atoi(h.URL.Query().Get("futureId"))
		orderType,_ := strconv.Atoi(h.URL.Query().Get("orderType"))
		percent,_ := strconv.Atoi(h.URL.Query().Get("percent"))
		id,_ :=  strconv.Atoi(p.ByName("id"))
		comm := entities.Commission{
			ID: id,
			FirmId: firmId,
			FutureId: futureId,
			OrderType: orderType,
			CommissionPercent: percent,
		}
		rte.q.Save(&comm)
		echo(w, comm)
	})

	rte.http.GET("/admin/commissions", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		var res []entities.Commission
		rte.q.Query().Find(&res)
		echo(w, res)
	})
	
	rte.http.POST("/admin/commissions", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		firmId,_ :=strconv.Atoi(h.URL.Query().Get("firmId"))
		futureId,_ := strconv.Atoi(h.URL.Query().Get("futureId"))
		orderType,_ := strconv.Atoi(h.URL.Query().Get("orderType"))
		id,_ :=  strconv.Atoi(p.ByName("id"))
		percent,_ := strconv.Atoi(h.URL.Query().Get("percent"))
		comm := entities.Commission{
			ID: id,
			FirmId: firmId,
			FutureId: futureId,
			OrderType: orderType,
			CommissionPercent: percent,
		}
		rte.q.Save(&comm)
		echo(w, comm)
	})

	rte.http.GET("/admin/firmFutures", func(w http.ResponseWriter, h *http.Request, p httprouter.Params) {
		var res []entities.FirmFuture
		rte.q.Query().Find(res)
		echo(w, res)
	})


}

func (r *router) Start(port int)  {
	http.ListenAndServe(":"+strconv.Itoa(port),r.http)
}


