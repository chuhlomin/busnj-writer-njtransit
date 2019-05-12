package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chuhlomin/busnj-writer-njtransit/pkg/redis"
	njt "github.com/chuhlomin/njtransit"

	"github.com/caarlos0/env"
)

type config struct {
	BusdataUsername       string        `env:"BUSDATA_USERNAME,required"`
	BusdataPassword       string        `env:"BUSDATA_PASSWORD,required"`
	BusdataUpdateInterval time.Duration `env:"BUSDATA_UPDATE_INTERVAL" envDefault:"5S"`
	RedisNetwork          string        `env:"REDIS_NETWORK" envDefault:"tcp"`
	RedisAddr             string        `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisSize             int           `env:"REDIS_SIZE" envDefault:"10"`
}

type busVehicleDataMessage struct {
	VehicleID            string `json:"vehicleID"`
	Route                string `json:"route"`                // 1
	RunID                string `json:"runID"`                // 21
	TripBlock            string `json:"tripBlock"`            // 001HL064
	PatternID            string `json:"patternID"`            // 264
	Destination          string `json:"destination"`          // 1 NEWARK-IVY HILL VIA RIVER TERM
	Longitude            string `json:"longitude"`            // -74.24513778686523
	Latitude             string `json:"latitude"`             // 40.73779029846192
	GPSTimestmp          string `json:"GPStimestmp"`          // 25-Apr-2019 12:15:12 AM
	LastModified         string `json:"lastModified"`         // 25-Apr-2019 12:16:10 AM
	AsInternalTripNumber string `json:"asInternalTripNumber"` // 13734490
}

func check(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	log.Println("Starting...")

	c := config{}
	err := env.Parse(&c)
	check(err, "Faield to parse environment variables")

	busData := njt.NewBusDataClient(
		c.BusdataUsername,
		c.BusdataPassword,
		njt.BusDataProdURL,
	)

	redis, err := redis.NewClient(
		c.RedisNetwork,
		c.RedisAddr,
		c.RedisSize,
	)
	check(err, "Failed to create Redis client")

	run(busData, redis, c.BusdataUpdateInterval)
}

func run(busData *njt.BusDataClient, redis *redis.Client, updateInterval time.Duration) error {
	rc := make(chan njt.BusVehicleDataRow)
	ec := make(chan error)
	go busData.GetBusVehicleDataStream(rc, ec, updateInterval, true)

	for {
		select {
		case row := <-rc:
			message := busVehicleDataMessage{
				VehicleID:            row.VehicleID,
				Route:                row.Route,
				RunID:                row.RunID,
				TripBlock:            row.TripBlock,
				PatternID:            row.PatternID,
				Destination:          strings.TrimSpace(row.Destination),
				Longitude:            row.Longitude,
				Latitude:             row.Latitude,
				GPSTimestmp:          row.GPSTimestmp,
				LastModified:         row.LastModified,
				AsInternalTripNumber: row.AsInternalTripNumber,
			}

			response, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to marshal busVehicleDataMessage: %v", err)
				continue
			}

			err = redis.SaveBusVehicleDataMessage(row.VehicleID, response)
			if err != nil {
				log.Printf("Failed to save message to Redis: %v", err)
				continue
			}

		case err := <-ec: // errors in the library
			log.Printf("njtransit-go library returned error: %v", err)
		}
	}
}
