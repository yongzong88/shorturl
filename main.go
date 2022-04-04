package main

import (
	"fmt"
	"net/http"
	"shorturl/redis"
	"shorturl/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := router()
	r.Run(":9000")
}

func router() *gin.Engine {
	r := gin.Default()
	r.POST("/api/v1/urls", add)
	r.GET("/:hash", visit)
	return r
}

func add(c *gin.Context) {
	var m map[string]interface{}
	err := c.Bind(&m)
	if err != nil {
		return
	}
	target, ok := m["url"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "Invaid request"})
		return
	}
	expire, ok := m["expireAt"].(string)
	if !ok {
		c.JSON(400, gin.H{"message": "Invaid request"})
		return
	}
	t, err := time.Parse(time.RFC3339, expire)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invaid request"})
		return
	}
	_expire := int(t.Sub(time.Now()).Seconds())

	// 產生一個沒用過的ID
	var textid string
	for used := true; used; used = redis.IsUsed(textid) {
		textid = utils.RandID()
	}

	err = redis.Set(textid, target, _expire)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       textid,
		"shortUrl": fmt.Sprintf("http://%s/%s", c.Request.Host, textid),
	})
}

func visit(c *gin.Context) {
	surl := c.Param("hash")
	lurl, err := redis.Get(surl)
	if err != nil {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
	c.Redirect(http.StatusPermanentRedirect, lurl)
}
