package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"shorturl/encdec"
	"shorturl/redis"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/", add)
	r.GET("/:hash", visit)
	r.Run(":9000")

}

func add(c *gin.Context) {
	target := c.PostForm("target")
	expire := c.PostForm("expire")
	_expire, _ := strconv.Atoi(expire)

	// 產生一個沒用過的ID
	var id uint64
	for used := true; used; used = redis.IsUsed(encdec.Encode(id)) {
		id = uint64(238328) + uint64(rand.Intn(134567890))
	}

	textid := encdec.Encode(id) // ID轉為短網址

	err := redis.Set(textid, target, _expire)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "OK",
		"surl":    fmt.Sprintf("http://%s/%s", c.Request.Host, textid),
	})
}

func visit(c *gin.Context) {
	surl := c.Param("hash")
	lurl, err := redis.Get(surl)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "ERROR",
		})
	}
	c.Redirect(http.StatusPermanentRedirect, lurl)
}
