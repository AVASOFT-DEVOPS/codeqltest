package domain

import (
	"golang/dto"
	"golang/errs"
	"net"

	"golang.org/x/crypto/ssh"
)

type LoadRepository interface {
	GetBTMSLoadResultRepo(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]LoadSearchResDTO, int, *errs.AppErrorvalidation)
	GetCPDBLoadResultRepo(LoadSearchReqDTO *dto.LoadSearchReqDTO) ([]LoadSearchResDTO, int, *errs.AppErrorvalidation)
	GetBTMSLoadDocsRepo(LoadDocsReqDTO *dto.LoadDocumentsReqDTO) ([]LoadBTMSDBResDTO, *errs.AppErrorvalidation)
	GetCPDBLoadDocsRepo(LoadDocsReqDTO *dto.LoadDocumentsReqDTO) ([]LoadCPDBResDTO, *errs.AppErrorvalidation)
}

type LoadSearchResDTO struct {
	ShipDateStart  string  `db:"shipDateStart"`
	ShipDateEnd    string  `db:"shipDateEnd"`
	LoadStatus     string  `db:"loadStatus"`
	LoadMethod     string  `db:"loadMethod"`
	LoadId         string  `db:"loadId"`
	PickupNumber   *string `db:"pickupNumber"`
	DeliveryNumber *string `db:"deliveryNumber"`
	CustTotal      string  `db:"custTotal"`
	ShipperCity    string  `db:"shipperCity"`
	ConsigneeCity  string  `db:"consigneeCity"`
	ShipperState   string  `db:"shipperState"`
	ConsigneeState string  `db:"consigneeState"`
	ProNumber      *string `db:"proNumber"`
	PoNumber       *string `db:"poNumber"`
	ShipBlNumber   *string `db:"shipBlNumber"`
	InvoiceDate    *string `db:"invoiceDate"`
	LastModified   string  `db:"lastModified"`
	LoadOrigin     string  `db:"loadOrigin"`
}

func (s LoadSearchResDTO) ToLoadDto() dto.LoadSearchResDTO {

	return dto.LoadSearchResDTO{
		ShipDateStart:  s.ShipDateStart,
		ShipDateEnd:    s.ShipDateEnd,
		LoadStatus:     s.LoadStatus,
		LoadMethod:     s.LoadMethod,
		LoadId:         s.LoadId,
		PickupNumber:   s.PickupNumber,
		DeliveryNumber: s.DeliveryNumber,
		CustTotal:      s.CustTotal,
		ShipperCity:    s.ShipperCity,
		ConsigneeCity:  s.ConsigneeCity,
		ShipperState:   s.ShipperState,
		ConsigneeState: s.ConsigneeState,
		ProNumber:      s.PoNumber,
		PoNumber:       s.PoNumber,
		ShipBlNumber:   s.ShipBlNumber,
		InvoiceDate:    s.InvoiceDate,
		LastModified:   s.LastModified,
		LoadOrigin:     s.LoadOrigin,
	}
}

type LoadCountRes struct {
	LoadCount int `db:"loadCount"`
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

type LoadBTMSDBResDTO struct {
	LoadStatus string `db:"loadStatus"`
	LoadMethod string `db:"loadMethod"`
	LoadId     string `db:"loadId"`
}

type LoadCPDBResDTO struct {
	Id      int    `db:"id"`
	TypeId  int    `db:"typeId"`
	LoadId  string `db:"loadId"`
	DocName string `db:"docName"`
	DocUrl  string `db:"docUrl"`
}
