package main

import (
	wrapper "cloudeya/coronavirusapi-go"
	"fmt"
	"time"
)

func main() {
	apiWrapper := wrapper.NewWrapper("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3RhcGkxIiwiaWF0IjoxNjAwMTgyODYzLCJleHAiOjE2MDAzODI4NjN9.YTf_Fx_GDKKBvST_jeVhL-YLbz6ZSuSYYQjqJyNPgQY")
	date, _ := time.Parse("2006-Jan-02", "2020-Sep-02")
	reports, err := apiWrapper.GetReportsAt(date)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Got %d reports for date %s\n", len(reports.Reports), date)
}
