# Coronavirus API Wrapper
*A Golang wrapper to work with the [Coronavirus API](https://www.covid19api.dev/)*

# How to use

## Create a new wrapper using API token

To create a new wrapper just call the following function:
```go
func NewWrapper(token string) *Wrapper
```
it will return the new wrapper ready to use to request the API.

## Create a new wrapper using username and password
To create a new wrapper username and password just call the function: 
```go
func NewWrapperWithCredentials(username string, password string)  (*Wrapper, error)
```
it will request the token using the credentials and store the token for further request. The username and password is not stored.

## Request a new token using username and password
if you are only intersted in the token, you can call the function:
```go
func (wrapper *Wrapper) GetTokenUsingCredentials(username string, password string) (string, error)
```
it will return the token if the call succeded, or an error otherwise.
## Available API call in the wrapper
Here is a list of available methods in the `Wrapper`:

### Get Global Daily Reports
`GetReportsAt` gets the Global Daily Reports at a given date.
```go
func (wrapper *Wrapper) GetReportsAt(date time.Time) (*GlobalDailyReport, error) 
```

#### Example

```go
package main

import (
	wrapper "cloudeya/coronavirusapi-go"
	"fmt"
	"time"
)

func main() {
	apiWrapper := wrapper.NewWrapper("Your_Token")
	// Get report for sep2020
	date, _ := time.Parse("2006-Jan", "2020-Sep")
	reports, err := apiWrapper.GetReportsAt(date)
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("Got %d reports for date %s\n", len(reports.Reports), date)
    fmt.Printf("Code: %d\nMessage: %s\nfirst report: %+v\n\nsecond report: %+v\n", reports.Code, reports.Message, reports.Reports[0], reports.Reports[1])
}
```

### Time Series Summary
The follwing methods get the time series for either Deaths, recovered or confirmed cases in the US or Globaly.

#### Get Time Series For a case type and country
You make a custom request using the follwing methode, timeSeriesType can either be "deaths", "recovered" or "confirmed" and country either "us" or "global".
```go
func (wrapper *Wrapper) GetTimeSeriesFor(timeSeriesType TimeSeriesType, country TimeSeriesCountry) (*TimeSeriesSummary, error) 
```

#### Confirmed-Global Time Series 
```go
func (wrapper *Wrapper) GetTimeSeriesConfirmedGlobal() (*TimeSeriesSummary, error)
```

#### Example

```go
package main

import wrapper "cloudeya/coronavirusapi-go"

func main() {
	apiWrapper := wrapper.NewWrapper("token here")
    deathsGlobal, err := apiWrapper.GetTimeSeriesConfirmedGlobal()
}
```
#### Confirmed-US Time Series 

```go
func (wrapper *Wrapper) GetTimeSeriesConfirmedUS() (*TimeSeriesSummary, error)
```

#### Example

```go
package main

import wrapper "cloudeya/coronavirusapi-go"

func main() {
	apiWrapper := wrapper.NewWrapper("token here")
    deathsGlobal, err := apiWrapper.GetTimeSeriesConfirmedUS()
}
```

#### Global-Deaths Time Series 

```go
func (wrapper *Wrapper) GetTimeSeriesDeathsGlobal() (*TimeSeriesSummary, error)
```

#### US-Deaths Time Series 

```go
func (wrapper *Wrapper) GetTimeSeriesDeathsUS() (*TimeSeriesSummary, error)
```


#### Global-Recovered Time Series 

```go
func (wrapper *Wrapper) GetTimeSeriesRecoveredGlobal() (*TimeSeriesSummary, error)
```

### Available structures

#### Global Daily Reports

```go
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
```


#### Time Series Summary

```go
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
 ```

`Data` contains as keys the date and the values is the number of cases.
Example:

```
{
    "may312020":221,
    "jun012020":233,
    "jun022020":239,
    "jun032020":239,
    "jun042020":241,
    "jun052020":248,
}
 ```

### Changing wrapper defaults

#### API URL
 By default the wrapper uses as a URL: `"https://covid19.cloudeya.org"` to change it all you need to do is call the function:

 ```go
func (wrapper *Wrapper) SetTimeout(apiTimout time.Duration)
 ```

#### HTTP Timeout when calling the API
By default the timeout is set to: 10 seconds, to change it call the function:
 ```go
func (wrapper *Wrapper) SetTimeout(apiTimout time.Duration)
  ```

#### Wrapper sleep time between retry
By default the wrapper sleeps for 60 seconds if the first call to the API is TooManyRequests, to change this duration, call this method:
 ```go
func (wrapper *Wrapper) SetTimeSleepBetweenRetry(sleepDuration time.Duration) 
  ```