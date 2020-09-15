package wrapper

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/errors/fmt"
)

const (
	apiURL    = "https://covid19.cloudeya.org"
	apiTimout = time.Second * 10
)

// Wrapper hold the data necessary for calling the endpoints
type Wrapper struct {
	token  string
	client *http.Client
	logger *log.Logger
}

// GlobalDailyReport the Global Daily Reports
type GlobalDailyReport struct {
	Code    int      `json:"Code"`
	Message string   `json:"Message"`
	Reports []Report `json:"Document"`
}

// Report holds the cases data for each month (the reports inside document)
type Report struct {
	ID            int32  `json:"id"`
	ProvinceState string `json:"province_state"`
	CountryRegion string `json:"country_region"`
	LastUpdate    string `json:"last_update"`
	Confirmed     int32  `json:"confirmed"`
	Deaths        int32  `json:"deaths"`
	Recovered     int32  `json:"recovered"`
}

// TimeSeriesSummary Time Series Summary
type TimeSeriesSummary struct {
	Code     int32
	Message  string
	Document []TimeSeriesReport
}

// TimeSeriesReport holds the time series data
type TimeSeriesReport struct {
	ID            int32
	ProvinceState string
	CountryRegion string
	Latitude      float32
	Longitude     float32
	Data          map[string]int32
}

// NewWrapper creates a new wrapper with the given token
func NewWrapper(token string) *Wrapper {
	client := &http.Client{Timeout: apiTimout}
	return &Wrapper{
		token:  token,
		client: client,
		logger: log.New(os.Stdout, "", 0),
	}
}

// NewWrapperWithCredentials creates a new wrapper with the given token and logger to log to
func NewWrapperWithCredentials(username string, password string) (*Wrapper, error) {
	client := &http.Client{Timeout: apiTimout}
	wrapper := &Wrapper{
		client: client,
		logger: &log.Logger{},
	}
	token, err := wrapper.GetTokenUsingCredentials(username, password)
	if err != nil {
		return nil, err
	}
	wrapper.token = token
	return wrapper, nil
}

// SetLogger changes the logger
func (wrapper *Wrapper) SetLogger(logger *log.Logger) {
	wrapper.logger = logger
}

// GetTokenUsingCredentials requests a token from the API using the username and password
func (wrapper *Wrapper) GetTokenUsingCredentials(username string, password string) (string, error) {
	wrapper.logger.Printf("Start: getting token for %s\n", username)
	url := fmt.Sprintf("%s/token", apiURL)
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
	url := fmt.Sprintf("%s/%s%d", apiURL, strings.ToLower(date.Month().String()[:3]), date.Year())
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
	decoder := json.NewDecoder(res.Body)
	globalDailyReport := GlobalDailyReport{}
	err = decoder.Decode(&globalDailyReport)
	if err != nil {
		wrapper.logger.Printf("Error while trying to get report at %s, err: %s\n", date, err)
		return nil, err
	}
	return &globalDailyReport, nil
}
