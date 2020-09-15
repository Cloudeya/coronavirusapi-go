package main

import "time"

const (
	// APIURL is the url of the covid19 API
	APIURL = "https://covid19.cloudeya.org"
)

// Wrapper hold the data necessary for calling the endpoints
type Wrapper struct {
	token string
}

// GlobalDailyReport the Global Daily Reports
type GlobalDailyReport struct {
	Code    int      `json:"Code,int"`
	Message string   `json:"Message,string"`
	Reports []Report `json:"Document"`
}

// Report holds the cases data for each month (the reports inside document)
type Report struct {
	ID            int       `json:"id,int"`
	ProvinceState string    `json:"province_state,string"`
	CountryRegion string    `json:"country_region,string"`
	LastUpdate    time.Time `json:"last_update,string"`
	Confirmed     int       `json:"confirmed,string"`
	Deaths        int       `json:"deaths,string"`
	Recovered     int       `json:"recovered,string"`
}

// TimeSeriesSummary Time Series Summary
type TimeSeriesSummary struct {
	Code     int
	Message  string
	Document []TimeSeriesReport
}

// TimeSeriesReport holds the time series data
type TimeSeriesReport struct {
	ID            int
	ProvinceState string
	CountryRegion string
	Latitude      float64
	Longitude     float64
	Data          map[string]int
}

// NewWrapper creates a new wrapper with the given token
func NewWrapper(token string) Wrapper {
	return Wrapper{
		token,
	}
}

// GetRecordsAt return the list of cases in the specified month-year i.e api.GetRecordsAt("sep", "2020")
func (wrapper *Wrapper) GetRecordsAt(month string, year string) GlobalDailyReport {
	return GlobalDailyReport{}
}
