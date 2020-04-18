package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getIndonesiaCoronaData() (result []AttributeIndonesianData) {
	req, err := http.NewRequest("GET", dataURLIndonesia, nil)
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
		lockIndonesianLastUpdate.RLock()
		result[i].Attribute.LastUpdateStr = indonesianLastUpdate
		lockIndonesianLastUpdate.RUnlock()
	}

	return
}
