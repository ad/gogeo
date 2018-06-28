package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	version = "dev"
	date    = "unknown"
	ip_req  chan string
	ip_res  chan IpResponse
)

// type asset struct {
// 	bytes []byte
// 	info  os.FileInfo
// }

// var _bindata = map[string]func() (*asset, error){}

type IpResponse struct {
	body []byte
}

type Geodata struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Longitude   float64 `json:"lon"`
	Latitude    float64 `json:"lat"`
}

var port = flag.String("port", "9001", "Port to listen on")

func main() {
	log.Printf("Started version %s, built at %s", version, date)

	ip_req = make(chan string, 5)
	go answerData()
	ip_res = make(chan IpResponse, 5)

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHandler)

	log.Printf("listening on port %s", *port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+*port, mux))
}

func answerData() {
	// var db, _ = geoip2.Open("data/GeoLite2-City.mmdb")

	// if version != "dev" {
	var data, _ = Asset("GeoLite2-City.mmdb")
	db, _ = geoip2.FromBytes(data)
	// }

	for ip := range ip_req {
		log.Println(ip)
		record, err := db.City(net.ParseIP(ip))
		if err != nil {
			log.Fatal(err)
		}

		geodata := Geodata{City: record.City.Names["en"], Country: record.Country.Names["en"], CountryCode: record.Country.IsoCode, Longitude: record.Location.Longitude, Latitude: record.Location.Latitude}
		js, _ := json.Marshal(geodata)

		ip_res <- IpResponse{js}
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")

	if len(ip) > 0 {
		ip_req <- ip

		for {
			select {
			case resp := <-ip_res:
				w.Write(resp.body)
				return
			}
		}
	}
}

// func Asset(name string) ([]byte, error) {
// 	cannonicalName := strings.Replace(name, "\\", "/", -1)
// 	if f, ok := _bindata[cannonicalName]; ok {
// 		a, err := f()
// 		if err != nil {
// 			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
// 		}
// 		return a.bytes, nil
// 	}
// 	return nil, fmt.Errorf("Asset %s not found", name)
// }
