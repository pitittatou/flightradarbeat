package beater

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/pitittatou/flightradarbeat/config"
)

// flightradarbeat configuration.
type flightradarbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of flightradarbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &flightradarbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *flightradarbeat) eventFeeder(flights <-chan Flight, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case f := <-flights:
			event := beat.Event{
				Timestamp: time.Unix(f.Timestamp, 0),
				Fields: common.MapStr{
					"type":                     "flightradarbeat",
					"id":                       f.Id,
					"latitude":                 f.Latitude,
					"longitude":                f.Longitude,
					"heading":                  f.Heading,
					"altitude":                 f.Altitude,
					"on_ground":                f.OnGround,
					"ground_speed":             f.GroundSpeed,
					"vertical_speed":           f.VerticalSpeed,
					"number":                   f.Number,
					"icao":                     f.ICAO,
					"aircraft_code":            f.AircraftCode,
					"aircraft_type":            f.AircraftType,
					"squawk":                   f.Squawk,
					"call_sign":                f.CallSign,
					"registration":             f.Registration,
					"origin_airport_iata":      f.OriginAirportIATA,
					"destination_airport_iata": f.DestinationAirportIATA,
					"airline_iata":             f.AirlineIATA,
					"airline_icao":             f.AirlineICAO,
				},
			}
			bt.client.Publish(event)
		case <-stop:
			return
		}
	}
}

// Run starts flightradarbeat.
func (bt *flightradarbeat) Run(b *beat.Beat) error {
	logp.Info("Flightradarbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	flightChannel := make(chan Flight)
	stopChannel := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go bt.eventFeeder(flightChannel, stopChannel, &wg)

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			close(stopChannel)
			wg.Wait()
			return nil
		case <-ticker.C:
		}

		err := getFlights(flightChannel)
		if err != nil {
			logp.Err("Error fetching flights: %v", err)
		} else {
			logp.Info("Fetched new flights infos")
		}
	}
}

// Stop stops flightradarbeat.
func (bt *flightradarbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
