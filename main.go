package main

import (
	"fmt"
	"net/http"
	"reflect"
	"shorturl/redis"
	"shorturl/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/v1/urls", add)
	r.GET("/:hash", visit)
	r.Run(":9000")

}

func add(c *gin.Context) {
	var m map[string]interface{}
	err := c.Bind(&m)
	if err != nil {
		return
	}
	fmt.Println(m)
	target := m["url"].(string)
	expire := m["expireAt"].(string)
	fmt.Println(target, expire)
	fmt.Println(reflect.TypeOf(target), reflect.TypeOf(expire))
	t, _ := time.Parse(time.RFC3339, expire)
	fmt.Println(t)
	_expire := int(t.Sub(time.Now()).Seconds())
	fmt.Println(_expire)

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
