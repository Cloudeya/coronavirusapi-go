package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultAPIURL    = "https://api.covid19api.dev"
	defaultAPITimout = time.Second * 10
	defaultSleepTime = time.Second * 60
)

// Wrapper hold the data necessary for calling the endpoints
type Wrapper struct {
	token             string
	client            *http.Client
	logger            *log.Logger
	apiURL            string
	apiTimout         time.Duration
	sleepBetweenRetry time.Duration
}

// GlobalDailyReport the Global Daily Reports
type GlobalDailyReport struct {
	Code    int      `json:"Code"`
	Message string   `json:"Message"`
	Reports []Report `json:"Document"`
}

// Report holds the cases data for each month (the reports inside document)
type Report struct {
	ID                int     `json:"id,omitempty"`
	ProvinceState     string  `json:"province_state,omitempty"`
	CountryRegion     string  `json:"country_region,omitempty"`
	LastUpdate        string  `json:"last_update,omitempty"`
	Confirmed         int     `json:"confirmed,omitempty"`
	Deaths            int     `json:"deaths,omitempty"`
	Recovered         int     `json:"recovered,omitempty"`
	Active            int     `json:"active,omitempty"`
	FIPS              string  `json:"fips,omitempty"`
	Admin2            string  `json:"admin2,omitempty"`
	CaseFatalityRatio float64 `json:"case_fatality_ratio,omitempty"`
	CombinedKey       string  `json:"combined_key,omitempty"`
	IncidenceRate     float64 `json:"incidence_rate,omitempty"`
}

// TimeSeriesType for time_series_[type]_* -> deaths, confirmed or recovered
type TimeSeriesType string

const (
	// Deaths for time_series_deaths_*
	Deaths TimeSeriesType = "deaths"
	// Confirmed for time_series_confirmed_*
	Confirmed TimeSeriesType = "confirmed"
	// Recovered for time_series_recovered_*
	Recovered TimeSeriesType = "recovered"
)

// TimeSeriesCountry for country time_series_*_[location], US or Global
type TimeSeriesCountry string

const (
	// US for time_series_*_us
	US TimeSeriesCountry = "us"
	// Global for time_series_*_global
	Global TimeSeriesCountry = "global"
)

// TimeSeriesSummary Time Series Summary
type TimeSeriesSummary struct {
	Code    int
	Message string
	Reports []TimeSeriesReport
}

// TimeSeriesReport holds the time series data
type TimeSeriesReport struct {
	ID            int
	UID           int
	ISO2          string
	ISO3          string
	Code3         int
	FIPS          int
	Admin2        string
	CombinedKey   string
	ProvinceState string
	CountryRegion string
	Latitude      float64
	Longitude     float64
	Population    int
	Data          map[string]int
}

// NewWrapper creates a new wrapper with the given token
func NewWrapper(token string) *Wrapper {
	client := &http.Client{Timeout: defaultAPITimout}
	return &Wrapper{
		token:             token,
		client:            client,
		logger:            log.New(os.Stdout, "", 0),
		apiURL:            defaultAPIURL,
		apiTimout:         defaultAPITimout,
		sleepBetweenRetry: defaultSleepTime,
	}
}

// NewWrapperWithCredentials creates a new wrapper for the API with the given credentials username and password, \nThis function will request a token and use it for the upcoming requests.\nthe credentials are not stored.
func NewWrapperWithCredentials(username string, password string) (*Wrapper, error) {
	client := &http.Client{Timeout: defaultAPITimout}
	wrapper := &Wrapper{
		client:            client,
		logger:            &log.Logger{},
		apiURL:            defaultAPIURL,
		apiTimout:         defaultAPITimout,
		sleepBetweenRetry: defaultSleepTime,
	}
	token, err := wrapper.GetTokenUsingCredentials(username, password)
	if err != nil {
		return nil, err
	}
	wrapper.token = token
	return wrapper, nil
}

// SetLogger changes the logger to log to, by default the logger logs to stdout
func (wrapper *Wrapper) SetLogger(logger *log.Logger) {
	wrapper.logger = logger
}

// SetTimeout changes the apiTimout to be used within HTTP requests to the API
func (wrapper *Wrapper) SetTimeout(apiTimout time.Duration) {
	wrapper.client = &http.Client{Timeout: apiTimout}
	wrapper.apiTimout = apiTimout
}

// SetAPIUrl changes the apiURL to use in the requests
func (wrapper *Wrapper) SetAPIUrl(apiURL string) {
	wrapper.apiURL = apiURL
}

// SetTimeSleepBetweenRetry changes the time to sleep between retry to call an again when the API response is TooManyRequests
func (wrapper *Wrapper) SetTimeSleepBetweenRetry(sleepDuration time.Duration) {
	wrapper.sleepBetweenRetry = sleepDuration
}

// GetTokenUsingCredentials requests a token from the API using the username and password
func (wrapper *Wrapper) GetTokenUsingCredentials(username string, password string) (string, error) {
	wrapper.logger.Printf("Start: getting token for %s\n", username)
	url := fmt.Sprintf("%s/token", wrapper.apiURL)
	payload := strings.NewReader(fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", username, password))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get token for username %s, err: %s\n", username, err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := wrapper.client.Do(req)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get token for username %s, err: %s\n", username, err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get token for username %s, err: %s\n", username, err)
		return "", err
	}
	wrapper.logger.Printf("Finished: getting token for %s\n", username)
	return string(body), nil
}

// GetReportsAt return the list of cases in the specified month-year i.e api.GetRecordsAt("sep", "2020")
func (wrapper *Wrapper) GetReportsAt(date time.Time) (*GlobalDailyReport, error) {
	wrapper.logger.Printf("Start: Getting Global Daily Reports at %s\n", date)
	bearer := "Bearer " + wrapper.token
	url := fmt.Sprintf("%s/%s%d", wrapper.apiURL, strings.ToLower(date.Month().String()[:3]), date.Year())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get report at %s, err: %s\n", date, err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	wrapper.logger.Printf("Requesting HTTP GET %s\n", url)
	res, err := wrapper.client.Do(req)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get report at %s, err: %s\n", date, err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 429 {
		// Too many requests, wait for 60s and then retry
		wrapper.logger.Printf("Too many requests, wait for %fs and then retry\n", wrapper.sleepBetweenRetry.Seconds())
		time.Sleep(wrapper.sleepBetweenRetry)
		return wrapper.GetReportsAt(date)
	} else if res.StatusCode != 200 {
		return nil, errors.New("Get response wasn't OK")
	}
	decoder := json.NewDecoder(res.Body)
	globalDailyReport := GlobalDailyReport{}
	err = decoder.Decode(&globalDailyReport)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get report at %s, err: %s\n", date, err)
		return nil, err
	}
	return &globalDailyReport, nil
}

// GetTimeSeriesFor return the TimeSeries for TimeSeries Type, and TimeSeries country (i.e  deaths and US)
func (wrapper *Wrapper) GetTimeSeriesFor(timeSeriesType TimeSeriesType, country TimeSeriesCountry) (*TimeSeriesSummary, error) {
	wrapper.logger.Printf("Start: Getting TimeSeries Reports for %s_%s", timeSeriesType, country)
	bearer := "Bearer " + wrapper.token
	url := fmt.Sprintf("%s/time_series_%s_%s", wrapper.apiURL, timeSeriesType, country)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get Time Series Reports for %s_%s, err: %s\n", timeSeriesType, country, err)
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	wrapper.logger.Printf("Requesting HTTP GET %s\n", url)
	res, err := wrapper.client.Do(req)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get Time Series Reports for %s_%s, in wrapper.client.Do , err: %s\n", timeSeriesType, country, err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 429 {
		// Too many requests, wait sometime and then retry
		wrapper.logger.Printf("Too many requests, wait for %fs and then retry\n", wrapper.sleepBetweenRetry.Seconds())
		time.Sleep(wrapper.sleepBetweenRetry)
		return wrapper.GetTimeSeriesFor(timeSeriesType, country)
	} else if res.StatusCode != 200 {
		return nil, errors.New("Get response wasn't OK")
	}
	decoder := json.NewDecoder(res.Body)
	var rawData map[string]interface{}
	err = decoder.Decode(&rawData)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get Time Series Reports for %s_%s,in decoder.Decode, err: %s\n", timeSeriesType, country, err)
		return nil, err
	}
	timeSeriesSummary := TimeSeriesSummary{}
	for k, v := range rawData {
		switch k {
		case "Code":
			timeSeriesSummary.Code, _ = v.(int)
		case "Message":
			timeSeriesSummary.Message, _ = v.(string)

		case "Document":
			elements, _ := v.([]interface{})
			timeSeriesSummary.Reports = []TimeSeriesReport{}
			for i := 0; i < len(elements); i++ {
				timeSeriesReport := TimeSeriesReport{}
				timeSeriesReport.Data = make(map[string]int)
				mappedElements, _ := elements[i].(map[string]interface{})
				for k2, v2 := range mappedElements {
					switch k2 {
					case "id":
						value, _ := v2.(float64)
						timeSeriesReport.ID = int(value)
					case "uid":
						value, _ := v2.(float64)
						timeSeriesReport.UID = int(value)
					case "province_state":
						timeSeriesReport.ProvinceState, _ = v2.(string)
					case "country_region":
						timeSeriesReport.CountryRegion, _ = v2.(string)
					case "latitude":
						timeSeriesReport.Latitude, _ = v2.(float64)
					case "longitude":
						timeSeriesReport.Longitude = v2.(float64)
					case "iso2":
						timeSeriesReport.ISO2, _ = v2.(string)
					case "iso3":
						timeSeriesReport.ISO3, _ = v2.(string)
					case "code3":
						value, _ := v2.(float64)
						timeSeriesReport.Code3 = int(value)
					case "fips":
						value, _ := v2.(float64)
						timeSeriesReport.FIPS = int(value)
					case "admin2":
						timeSeriesReport.Admin2, _ = v2.(string)
					case "combined_key":
						timeSeriesReport.CombinedKey, _ = v2.(string)
					case "population":
						value, _ := v2.(float64)
						timeSeriesReport.Population = int(value)
					default:
						value, _ := v2.(float64)
						timeSeriesReport.Data[k2] = int(value)
					}
				}
				timeSeriesSummary.Reports = append(timeSeriesSummary.Reports, timeSeriesReport)
			}
		default:
			wrapper.logger.Printf("Unknow key found %s", k)
		}
	}
	return &timeSeriesSummary, nil
}

// GetTimeSeriesConfirmedGlobal return the timeseries for cofirmed cases globaly
func (wrapper *Wrapper) GetTimeSeriesConfirmedGlobal() (*TimeSeriesSummary, error) {
	return wrapper.GetTimeSeriesFor(Confirmed, Global)
}

// GetTimeSeriesConfirmedUS return the timeseries for cofirmed cases in US
func (wrapper *Wrapper) GetTimeSeriesConfirmedUS() (*TimeSeriesSummary, error) {
	return wrapper.GetTimeSeriesFor(Confirmed, US)
}

// GetTimeSeriesDeathsGlobal return the timeseries for Deaths cases globaly
func (wrapper *Wrapper) GetTimeSeriesDeathsGlobal() (*TimeSeriesSummary, error) {
	return wrapper.GetTimeSeriesFor(Deaths, Global)
}

// GetTimeSeriesDeathsUS return the timeseries for Deaths cases in US
func (wrapper *Wrapper) GetTimeSeriesDeathsUS() (*TimeSeriesSummary, error) {
	return wrapper.GetTimeSeriesFor(Deaths, US)
}

// GetTimeSeriesRecoveredGlobal return the timeseries for recovered cases globaly
func (wrapper *Wrapper) GetTimeSeriesRecoveredGlobal() (*TimeSeriesSummary, error) {
	return wrapper.GetTimeSeriesFor(Recovered, Global)
}
