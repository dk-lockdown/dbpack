package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dbpack/samples/aggregation_svc/svc"
)

func main() {
	r := gin.Default()

	r.POST("/v1/order/create", func(c *gin.Context) {
		xid := c.GetHeader("x_dbpack_xid")
		err := svc.GetSvc().CreateSo(c, xid, false)
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

	r.GET("/v1/order/create2", func(c *gin.Context) {
		xid := c.GetHeader("x_dbpack_xid")
		err := svc.GetSvc().CreateSo(c, xid, true)
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

	r.Run(":3000")
}
