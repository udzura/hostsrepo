package hostsrepo

import (
	"bufio"
	"fmt"
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

// XXX: re-invent of https://golang.org/src/net/hosts.go
func FindAllHosts(hostsFile string) (Hosts, error) {
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
			fmt.Fprintf(os.Stdin, "Error hosts entry: %v\n", ent)
		} else if strings.HasPrefix(ent[0], "#") {
			fmt.Fprintf(os.Stdin, "Entry is commented out: %v\n", ent)
		} else {
			ip := ent[0]
			for _, hn := range ent[1:] {
				if strings.HasPrefix(hn, "#") {
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
		hosts, err := FindAllHosts("")
		if err != nil {
			c.JSON(503, err)
		} else {
			c.JSON(200, hosts)
		}
	})
	return r.Run("0.0.0.0:10808")
}
