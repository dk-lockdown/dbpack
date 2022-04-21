package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // register mysql

	"github.com/dbpack/samples/product_svc/dao"
)

var dsn = `dksl:123456@tcp(127.0.0.1:13307)/product?timeout=60s&readTimeout=60s&writeTimeout=60s&parseTime=true&loc=Local&charset=utf8mb4,utf8`

func main() {
	r := gin.Default()

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(4 * time.Hour)

	if err != nil {
		panic(err)
	}
	d := &dao.Dao{
		DB: sqlDB,
	}

	r.POST("/allocateInventory", func(c *gin.Context) {
		xid := c.GetHeader("xid")

		type req struct {
			Req []*dao.AllocateInventoryReq
		}
		var q req
		if err := c.ShouldBindJSON(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := d.AllocateInventory(c, xid, q.Req)

		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "fail",
			})
		} else {
			c.JSON(200, gin.H{
				"success": true,
				"message": "success",
			})
		}
	})

	r.Run(":3002")
}
