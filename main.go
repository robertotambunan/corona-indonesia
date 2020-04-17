package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

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
	http.HandleFunc("/", templateHandler)
	log.Println("running in port :", port)
	http.ListenAndServe(":"+port, nil)
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, AllDataCache)
}

type (
	// AllData : United all of the data
	AllData struct {
		Nations []AttributeNationData
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
)

const (
	dataURLNation = "https://api.kawalcorona.com/"
)

func runCron() {
	nationData := getWorldCoronaData()
	AllDataCache.Nations = nationData

	c := cron.New()
	c.AddFunc("@every 10m", func() {
		log.Println("Cron is Running every 10m")
		nationData := getWorldCoronaData()
		if len(nationData) > 0 {
			lock.Lock()
			AllDataCache.Nations = nationData
			lock.Unlock()
		}
	})
	// Start cron with one scheduled job
	c.Start()
}

func getWorldCoronaData() (result []AttributeNationData) {
	req, err := http.NewRequest("GET", dataURLNation, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("NewDecoder", err)
	}

	for i := range result {
		tm := time.Unix(result[i].Attribute.LastUpdate, 0)
		tempArr := strings.Split(tm.String(), " ")
		if len(tempArr) >= 3 {
			result[i].Attribute.LastUpdateStr = tempArr[1]
		}
	}
	return
}
