package dao

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	allocateInventorySql = `update /*+ XID('%s') */ product.inventory set available_qty = available_qty - ?, 
		allocated_qty = allocated_qty + ? where product_sysno = ? and available_qty >= ?`
)

type Dao struct {
	*sql.DB
}

type AllocateInventoryReq struct {
	ProductSysNo int64
	Qty          int32
}

func (dao *Dao) AllocateInventory(ctx context.Context, xid string, reqs []*AllocateInventoryReq) error {
	tx, err := dao.Begin()
	if err != nil {
		return err
	}
	updateInventory := fmt.Sprintf(allocateInventorySql, xid)
	for _, req := range reqs {
		_, err := tx.Exec(updateInventory, req.Qty, req.Qty, req.ProductSysNo, req.Qty)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
