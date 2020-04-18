package main

import (
	"strings"
	"time"
)

func orderNationData(data []AttributeNationData) (finalResult []AttributeNationData) {
	indexFindIndonesia := 0
	for i := range data {
		if strings.ToLower(data[i].Attribute.CountryRegion) == "indonesia" {
			indexFindIndonesia = i
		}
	}

	//updating indonesia last update
	lockIndonesianLastUpdate.Lock()
	indonesianLastUpdate = data[indexFindIndonesia].Attribute.LastUpdateStr
	lockIndonesianLastUpdate.Unlock()

	finalResult = append(finalResult, data[indexFindIndonesia])
	remove(data, indexFindIndonesia)
	finalResult = append(finalResult, data...)
	return
}

func remove(slice []AttributeNationData, s int) []AttributeNationData {
	return append(slice[:s], slice[s+1:]...)
}

func generateIndonesiaTime(t time.Time) time.Time {
	wib, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), wib)
	}
	return t
}
