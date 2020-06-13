package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/robfig/cron"
)

var (
	// AllDataCache : caching data
	AllDataCache AllData
	lock         = sync.RWMutex{}
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	runCron()
	http.HandleFunc("/world", templateWorldHandler)
	http.HandleFunc("/", templateIndonesiaHandler)
	log.Println("running in port :", port)
	http.ListenAndServe(":"+port, nil)
}

func templateWorldHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index-world.html"))
	t.Execute(w, AllDataCache)
}

func templateIndonesiaHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index-indonesia.html"))
	t.Execute(w, AllDataCache)
}

type (
	// AllData : United all of the data
	AllData struct {
		Nations   []AttributeNationData
		Indonesia []AttributeIndonesianData
	}

	// NationData : Present Model of Nation
	NationData struct {
		CountryRegion string `json:"Country_Region"`
		Confirmed     int    `json:"Confirmed"`
		Deaths        int    `json:"Deaths"`
		Recovered     int    `json:"Recovered"`
		LastUpdate    int64  `json:"Last_Update"`
		LastUpdateStr string
	}

	// AttributeNationData : API struct for data
	AttributeNationData struct {
		Attribute NationData `json:"attributes"`
	}

	// IndonesianData : Present model of indonesian data
	IndonesianData struct {
		Provinsi      string `json:"Provinsi"`
		Confirmed     int    `json:"Kasus_Posi"`
		Recovered     int    `json:"Kasus_Semb"`
		Deaths        int    `json:"Kasus_Meni"`
		LastUpdateStr string
	}

	// AttributeIndonesianData : API struct for data
	AttributeIndonesianData struct {
		Attribute IndonesianData `json:"attributes"`
	}
)

func runCron() {
	nationData := getWorldCoronaData()
	AllDataCache.Nations = nationData
	indonesianData := getIndonesiaCoronaData()
	AllDataCache.Indonesia = indonesianData
	c := cron.New()
	c.AddFunc("@every 10m", func() {
		log.Println("Cron is Running every 10m")
		nationData := getWorldCoronaData()
		indonesianData := getIndonesiaCoronaData()
		if len(nationData) > 0 {
			lock.Lock()
			AllDataCache.Nations = nationData
			AllDataCache.Indonesia = indonesianData
			lock.Unlock()
		}
	})
	// Start cron with one scheduled job
	c.Start()
}
