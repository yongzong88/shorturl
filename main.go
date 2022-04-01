package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	// "shorturl/encdec"
	// "shorturl/redis"
)

func main() {
	r := gin.Default()
	r.POST("/", add)
	r.GET("/:hash", visit)
	r.RUN(":9000")

}

func add(c *gin.Context) {
	target := c.PostForm("target")
	expire := c.PostForm("expire")
	_expire, err := strconv.Atoi(expire)

}

func visit(c *gin.Context) {

}
