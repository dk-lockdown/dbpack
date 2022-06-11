package http

import (
	"context"
	"strings"

	"github.com/cectc/dbpack/pkg/log"
)

var (
	svc = &ProxyService{
		transactionInfos: make(map[string]*TransactionInfo, 0),
		tccResourceInfos: make(map[string]*TccResourceInfo, 0),
	}
)

type ProxyService struct {
	transactionInfos map[string]*TransactionInfo
	tccResourceInfos map[string]*TccResourceInfo
}

func (svc *ProxyService) RegisterTransactionInfo(ctx context.Context, req *TransactionInfo) (*Response, error) {
	svc.transactionInfos[strings.ToLower(req.RequestPath)] = req
	log.Debugf("proxy %s, will create global transaction, put xid into request header, if enabled distributed transaction.", req.RequestPath)
	return &Response{Success: true}, nil
}

func (svc *ProxyService) RegisterTccResourceInfo(ctx context.Context, req *TccResourceInfo) (*Response, error) {
	svc.tccResourceInfos[strings.ToLower(req.PrepareRequestPath)] = req
	log.Debugf("proxy %s, will register branch transaction, if enabled distributed transaction.", req.PrepareRequestPath)
	return &Response{Success: true}, nil
}

func AddTransactionInfo(req *TransactionInfo) {
	svc.transactionInfos[strings.ToLower(req.RequestPath)] = req
}

func AddTccResourceInfo(req *TccResourceInfo) {
	svc.tccResourceInfos[strings.ToLower(req.PrepareRequestPath)] = req
}

func GetTransactionInfos() map[string]*TransactionInfo {
	return svc.transactionInfos
}

func GetTccResourceInfo() map[string]*TccResourceInfo {
	return svc.tccResourceInfos
}
