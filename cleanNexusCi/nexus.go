package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func HttpGet(url string, cred Cred) []byte {
	client := &http.Client{}
	get, err := http.NewRequest(http.MethodGet, url, nil)
	get.SetBasicAuth(cred.usr, cred.pwd)
	req, err := client.Do(get)
	if err != nil {
		log.Fatal(err.Error())
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		log.Fatal(err.Error())
	}

	defer req.Body.Close()

	return body
}

func HttpDel(url string, cred Cred) (*http.Response, error) {
	log.Printf("http delete: %s", url)

	client := &http.Client{
		Timeout: 4 * time.Minute,
	}

	httReq, err := http.NewRequest(http.MethodDelete, url, nil)
	httReq.SetBasicAuth(cred.usr, cred.pwd)
	req, err := client.Do(httReq)
	return req, err
}

func deleteDistributive(url string, cred Cred, wg *sync.WaitGroup, quotaCh chan int,
	deletedCount *int64, countArchives int) {
	quotaCh <- 1
	defer wg.Done()

	req, err := HttpDel(url, cred)
	if err != nil {
		if !strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			log.Fatal(err.Error())
		} else {
			time.Sleep(time.Minute * 5)
		}
	} else {
		defer req.Body.Close()
		if req.StatusCode != 204 && req.StatusCode != 502 {
			log.Fatal(req.StatusCode, req.Status)
		}
	}

	atomic.AddInt64(deletedCount, 1)
	log.Printf("[%d/%d]", *deletedCount, countArchives)
	<-quotaCh
}

func deleteOldArchives(archives []OldArchives, cred Cred, jobQueue int) {
	quotaCh := make(chan int, jobQueue)
	wg := &sync.WaitGroup{}
	var deletedCount int64 = 0

	for archiveNumber := range archives {
		wg.Add(1)
		go deleteDistributive(archives[archiveNumber].url, cred, wg, quotaCh,
			&deletedCount, len(archives))
	}
	wg.Wait()
}
