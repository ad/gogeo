package main

import (
	"encoding/json"
	"flag"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"strings"
	"strconv"
)

var (
	version = "dev"
	date    = "unknown"
	ipReq  chan string
	ipRes  chan IPResponse
)

// IPResponse struct
type IPResponse struct {
	body []byte
}

// Geodata struct
type Geodata struct {
	City                         string  `json:"city"`
	Country                      string  `json:"country"`
	CountryCode                  string  `json:"country_code"`
	Longitude                    float64 `json:"lon"`
	Latitude                     float64 `json:"lat"`
	AutonomousSystemNumber       uint    `json:"asn"`
	AutonomousSystemOrganization string  `json:"provider"`
}

var port = flag.String("port", "9001", "Port to listen on")

func main() {
	log.Printf("Started version %s, built at %s", version, date)

	ipReq = make(chan string, 5)
	go answerData()
	ipRes = make(chan IPResponse, 5)

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHandler)

	log.Printf("listening on port %s", *port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+*port, mux))
}

func answerData() {
	// var dbCity, _ = geoip2.Open("data/GeoLite2-City.mmdb")
	// var dbASN, _ = geoip2.Open("data/GeoLite2-ASN.mmdb")

	var dataCity, _ = Asset("GeoLite2-City.mmdb")
	var dbCity, _ = geoip2.FromBytes(dataCity)

	var dataASN, _ = Asset("GeoLite2-ASN.mmdb")
	var dbASN, _ = geoip2.FromBytes(dataASN)

	for ip := range ipReq {
		log.Println(ip)
		geodata := Geodata{}
		recordCity, _ := dbCity.City(net.ParseIP(ip))
		recordAsn, _ := dbASN.ASN(net.ParseIP(ip))

		if recordCity != nil {
			geodata.City = recordCity.City.Names["en"]
			geodata.Country = recordCity.Country.Names["en"]
			geodata.CountryCode = recordCity.Country.IsoCode
			geodata.Longitude = recordCity.Location.Longitude
			geodata.Latitude = recordCity.Location.Latitude
		}

		if recordAsn != nil {
			geodata.AutonomousSystemNumber = recordAsn.AutonomousSystemNumber
			geodata.AutonomousSystemOrganization = recordAsn.AutonomousSystemOrganization
		}

		js, _ := json.Marshal(geodata)

		ipRes <- IPResponse{js}
	}
}

// GetHandler func
func GetHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")

	if len(ip) > 0 && isIPv4(ip) {
		ipReq <- ip

		for {
			select {
			case resp := <-ipRes:
				w.Write(resp.body)
				return
			}
		}
	} else {
		w.Write([]byte(`{"status": "error"}`))
	}
}

func isIPv4(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}
	
	for _,x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
			return false
		}
		} else {
			return false
		}

	}
	return true
}
