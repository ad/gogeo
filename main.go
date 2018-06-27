package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"net/http"
)

type Geodata struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Longitude   float64 `json:"lon"`
	Latitude    float64 `json:"lat"`
}

var port = flag.String("port", "9001", "Port to listen on")

var data, _ = Asset("GeoLite2-City.mmdb")
var db, _ = geoip2.FromBytes(data)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHandler)

	log.Printf("listening on port %s", *port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+*port, mux))

	defer db.Close()
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.FormValue("ip")

	if len(ip) > 0 {

		record, err := db.City(net.ParseIP(ip))
		if err != nil {
			log.Fatal(err)
		}
		geodata := Geodata{City: record.City.Names["en"], Country: record.Country.Names["en"], CountryCode: record.Country.IsoCode, Longitude: record.Location.Longitude, Latitude: record.Location.Latitude}
		js, _ := json.Marshal(geodata)

		fmt.Println(string(js))

		w.Write(js)
	}
}
