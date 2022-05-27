package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Datum struct {
	Timestamp     uint64
	CPU           uint64
	Concurrencies uint64
}

func IngestData() []Datum {
	var Data []Datum
	current := time.Now().Add(-5 * time.Minute)
	end := time.Now()
	for end.After(current) {
		rand.Seed(time.Now().UnixNano())
		datum := Datum{
			Timestamp:     uint64(current.Unix()),
			CPU:           uint64(rand.Intn(100)),
			Concurrencies: uint64(rand.Intn(500000)),
		}
		Data = append(Data, datum)
		current = current.Add(60 * time.Second)
		fmt.Println(current.Unix())
	}
	return Data
}

func Process() error {
	Data := IngestData()
	marshalledData, err := json.Marshal(Data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("test.json", marshalledData, 0644)
	file, err := ioutil.TempFile("", "data")
	file.Write(marshalledData)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := Process(); err != nil {
		log.Fatal(err)
	}
}
