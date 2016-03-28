package hostsrepo

import (
	"github.com/gin-gonic/gin"
)

type Host struct {
	Name   string `json:"name"`
	IPAddr string `json:"ipaddr"`
}
type Hosts []*Host

func RunServer() error {
	r := gin.Default()
	r.GET("/hosts", func(c *gin.Context) {
		hosts := Hosts{
			{Name: "localhost", IPAddr: "127.0.0.1"},
		}
		c.JSON(200, hosts)
	})
	return r.Run("0.0.0.0:10808")
}
