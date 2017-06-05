package executor

import (
	"github.com/shopspring/decimal"
	"broker-gateway/entities"
	"time"
)

type Inspector interface {
	GetCommission(futureId, firmId, orderType int, amount decimal.Decimal) decimal.Decimal
	InspectSetting()
	InspectFirmFuture()
	ValidateFutureId(firmId, futureId int) bool
}




type inspector struct {
	db DB
	//  setting[futureId][firmId][type]
	commissionSetting map[int]map[int]map[int]int
	firmFutures map[int][]int
}

func NewInspector(db DB) Inspector {
	return &inspector{
		db:db,
		commissionSetting: map[int]map[int]map[int]int{},
	}
}

func (in *inspector) GetCommission(futureId, firmId, orderType int, amount decimal.Decimal) decimal.Decimal {
	firmAndOrderType, exist := in.commissionSetting[futureId]
	if ! exist {
		return decimal.Zero
	}
	drTy, exist := firmAndOrderType[firmId]
	if ! exist {
		return decimal.Zero
	}

	percent, exist := drTy[orderType]
	if ! exist {
		return decimal.Zero
	}
	return amount.Mul(decimal.New(int64(percent),-2))
}

func (in *inspector) ValidateFutureId(firmId, futureId int) bool  {
	futures, exist := in.firmFutures[firmId]
	if ! exist {
		return true
	} else {
		for i:=0; i<len(futures) ;i++  {
			if futures[i] == futureId {
				return true
			}
		}
	}
	return false
}


func (in *inspector) InspectSetting()  {
	for true {
		var commissions []entities.Commission
		in.db.Query().Find(&commissions)
		for i:=0; i< len(commissions) ;i++  {

			firmAndOrderType, exist := in.commissionSetting[commissions[i].FutureId]
			if ! exist {
				in.commissionSetting[commissions[i].FutureId] = map[int]map[int]int{}
				firmAndOrderType = in.commissionSetting[commissions[i].FutureId]
			}
			drTy, exist := firmAndOrderType[commissions[i].FirmId]
			if ! exist {
				firmAndOrderType[commissions[i].FirmId] = map[int]int{}
				drTy = firmAndOrderType[commissions[i].FirmId]
			}

			//percent, _ := drTy[commissions[i].OrderType]
			drTy[commissions[i].OrderType] = commissions[i].CommissionPercent

		}
		time.Sleep(time.Second * 10)
	}

}


func (in *inspector) InspectFirmFuture()  {
	for true {
		in.firmFutures = map[int][]int{}
		var entis []entities.FirmFuture
		in.db.Query().Find(&entities)
		for i:=0; i< len(entis) ;i++  {

			_, exist := in.commissionSetting[entis[i].FirmId]
			if ! exist {
				in.firmFutures[entis[i].FirmId] = []int{entis[i].FutureId}
				continue
			} else {
				in.firmFutures[entis[i].FirmId] = append(in.firmFutures[entis[i].FirmId],entis[i].FutureId)
			}

		}
		time.Sleep(time.Second * 10)
	}

}
