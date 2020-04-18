package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

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
		timeMS := result[i].Attribute.LastUpdate / int64(1000)
		tm := time.Unix(timeMS, 0)
		tm = generateIndonesiaTime(tm)
		tempArr := strings.Split(tm.String(), " ")
		if len(tempArr) >= 3 {
			result[i].Attribute.LastUpdateStr = tempArr[1] + "-WIB " + tempArr[0]
		}
	}
	result = orderNationData(result)
	return
}
