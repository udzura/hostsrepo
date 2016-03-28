package hostsrepo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Host struct {
	Name   string `json:"name"`
	IPAddr string `json:"ipaddr"`
}
type Hosts []*Host

var spaces = regexp.MustCompile("\\s+")
var privateClassA, privateClassB, privateClassC *net.IPNet

func init() {
	_, privateClassA, _ = net.ParseCIDR("192.168.0.0/16")
	_, privateClassB, _ = net.ParseCIDR("172.16.0.0/12")
	_, privateClassC, _ = net.ParseCIDR("10.0.0.0/8")
}

// XXX: re-invent of https://golang.org/src/net/hosts.go
func FindHosts(hostsFile string, domain string, privateOnly bool) (Hosts, error) {
	if hostsFile == "" {
		hostsFile = "/etc/hosts"
	}

	hosts := Hosts{}
	fp, err := os.Open(hostsFile)
	if err != nil {
		return hosts, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		ent := spaces.Split(strings.TrimSpace(scanner.Text()), -1)
		if len(ent) < 2 {
			fmt.Fprintf(os.Stdin, "Entry is empty or invalid: %v\n", ent)
		} else if strings.HasPrefix(ent[0], "#") {
			fmt.Fprintf(os.Stdin, "Entry is commented out: %v\n", ent)
		} else {
			ip := ent[0]
			if privateOnly {
				isPrivate := false
				theIP := net.ParseIP(ip)
				for _, net := range []*net.IPNet{privateClassA, privateClassB, privateClassC} {
					if net.Contains(theIP) {
						isPrivate = true
					}
				}
				if !isPrivate {
					continue
				}
			}
			for _, hn := range ent[1:] {
				if strings.HasPrefix(hn, "#") {
					break
				}
				if domain != "" && !strings.HasSuffix(hn, domain) {
					break
				}
				hosts = append(hosts, &Host{Name: hn, IPAddr: ip})
			}
		}

	}
	if err := scanner.Err(); err != nil {
		return hosts, err
	}

	return hosts, nil
}

func RunServer() error {
	r := gin.Default()
	r.GET("/hosts", func(c *gin.Context) {
		hosts, err := FindHosts("", "", false)
		if err != nil {
			c.JSON(503, err)
		} else {
			c.JSON(200, hosts)
		}
	})
	r.GET("/hosts/private", func(c *gin.Context) {
		hosts, err := FindHosts("", "", true)
		if err != nil {
			c.JSON(503, err)
		} else {
			c.JSON(200, hosts)
		}
	})
	r.GET("/hosts/domains/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		hosts, err := FindHosts("", domain, false)
		if err != nil {
			c.JSON(503, err)
		} else {
			c.JSON(200, hosts)
		}
	})
	r.GET("/hosts/private/domains/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		hosts, err := FindHosts("", domain, true)
		if err != nil {
			c.JSON(503, err)
		} else {
			c.JSON(200, hosts)
		}
	})

	return r.Run("0.0.0.0:10808")
}
