package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nicescale/pangolin"
)

func main() {
	log.SetLevel(log.DebugLevel)
	hub := pangolin.NewHub(nil)
	ech := make(chan error)
	go func(ech chan<- error) {
		ech <- hub.ListenAndServe(":9000")
	}(ech)

	client := http.Client{
		Transport: &http.Transport{
			Dial: hub.Dial,
		},
	}

	for {
		time.Sleep(time.Second * 1)
		resp, err := client.Get("http://node1/info")
		if err != nil {
			log.Error("http get: ", err)
		} else {
			d, err := httputil.DumpResponse(resp, true)
			fmt.Println(string(d), err)
		}
	}

	err := <-ech
	if err != nil {
		log.Fatal("server: ", err)
	}
}
