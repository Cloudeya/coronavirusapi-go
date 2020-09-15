package main

import (
	wrapper "cloudeya/coronavirusapi-go"
	"fmt"
	"time"
)

func main() {
	apiWrapper := wrapper.NewWrapper("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RhcGkxIiwiaWF0IjoxNjAwMTgyODYzLCJleHAiOjE2MDAzODI4NjN9.YTf_Fx_GDKKBvST_jeVhL-YLbz6ZSuSYYQjqJyNPgQY")
	// Get report for sep2020
	date, _ := time.Parse("2006-Jan", "2020-Sep")
	reports, err := apiWrapper.GetReportsAt(date)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Got %d reports for date %s\n", len(reports.Reports), date)
	fmt.Printf("Code: %d\nMessage: %s\nfirst report: %+v\n\nsecond report: %+v\n", reports.Code, reports.Message, reports.Reports[0], reports.Reports[1])
	fmt.Printf("\n-----------------\n")
	deathsGlobal, err := apiWrapper.GetTimeSeriesDeathsGlobal()
	fmt.Printf("Got %d reports for deaths globaly\n", len(deathsGlobal.Reports))
	fmt.Printf("Code: %d\nMessage: %s\nfirst report: %+v\n\nsecond report: %+v\nThird report: %+v\n", deathsGlobal.Code, deathsGlobal.Message, deathsGlobal.Reports[0], deathsGlobal.Reports[1], deathsGlobal.Reports[2])

	fmt.Printf("\n-----------------\n")
	deathsGlobal, err = apiWrapper.GetTimeSeriesDeathsUS()
	fmt.Printf("Got %d reports for deaths in US\n", len(deathsGlobal.Reports))
	fmt.Printf("Code: %d\nMessage: %s\nfirst report: %+v\n\nsecond report: %+v\nThird report: %+v\n", deathsGlobal.Code, deathsGlobal.Message, deathsGlobal.Reports[0], deathsGlobal.Reports[1], deathsGlobal.Reports[2])
}
