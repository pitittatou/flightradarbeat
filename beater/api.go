package beater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
)

type Flight struct {
	Id                     string
	Timestamp              int64
	ICAO                   string
	Latitude               float64
	Longitude              float64
	Heading                int64
	Altitude               int64
	GroundSpeed            int64
	Squawk                 string
	AircraftType           string
	AircraftCode           string
	Registration           string
	OriginAirportIATA      string
	DestinationAirportIATA string
	Number                 string
	AirlineIATA            string
	OnGround               bool
	VerticalSpeed          int64
	CallSign               string
	AirlineICAO            string
}

var HEADERS = map[string]string{
	"Accept-Encoding": "gzip",
	"Accept-Language": "fr-FR,fr;q=0.6",
	"Cache-Control":   "max-age=0",
	"Origin":          "https://www.flightradar24.com",
	"Referer":         "https://www.flightradar24.com/",
	"Sec-Fetch-Dest":  "empty",
	"Sec-Fetch-Mode":  "cors",
	"Sec-Fetch-Site":  "same-site",
	"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
}

func apiRequest(url string) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range HEADERS {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var bodyReader *bytes.Reader
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		gzipBody, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(gzipBody)
	} else {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewReader(body)
	}

	var data map[string]interface{}
	err = json.NewDecoder(bodyReader).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseFlight(id string, data []interface{}) Flight {
	flight := Flight{}

	flight.Id = id
	flight.ICAO = data[0].(string)
	flight.Latitude = data[1].(float64)
	flight.Longitude = data[2].(float64)
	flight.Heading = int64(data[3].(float64))
	flight.Altitude = int64(data[4].(float64))
	flight.GroundSpeed = int64(data[5].(float64))
	flight.Squawk = data[6].(string)
	flight.AircraftType = data[7].(string)
	flight.AircraftCode = data[8].(string)
	flight.Registration = data[9].(string)
	flight.Timestamp = int64(data[10].(float64))
	flight.OriginAirportIATA = data[11].(string)
	flight.DestinationAirportIATA = data[12].(string)
	flight.Number = data[13].(string)
	flight.AirlineIATA = func() string {
		if str, ok := data[13].(string); ok && str != "" {
			return str[:2]
		}
		return ""
	}()
	flight.OnGround = data[14].(float64) != 0.0
	flight.VerticalSpeed = int64(data[15].(float64))
	flight.CallSign = data[16].(string)
	flight.AirlineICAO = data[18].(string)

	return flight
}

func getFlights(flightChannel chan<- Flight) error {
	url := "https://data-cloud.flightradar24.com/zones/fcgi/feed.js"
	resp, err := apiRequest(url)
	if err != nil {
		return err
	}

	for key, value := range resp {
		if key != "full_count" && key != "version" {
			flightChannel <- parseFlight(key, value.([]interface{}))
		}
	}

	return nil
}
