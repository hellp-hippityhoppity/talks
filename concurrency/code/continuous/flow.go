package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
	"github.com/nathj07/talks/concurrency/code/common"
)

// Provider is the basic unit we are working with for this example
type Provider struct {
	url  string
	name string
}

var providerChan chan *Provider

func main() {
	db := common.GetDBConnection()

	providerChan = make(chan *Provider) // unbuffered
	// fetch data
	go fetchData(db)
	// blocking calls
	useData()
	//useDataWorkerPool() // using a worker pool to interleave the tasks
}

func fetchData(db *sql.DB) {
	// bounded loop to simulate repeated fetches
	for i := 0; i <= 1; i++ {
		fmt.Println("Iteration ", i)
		fmt.Println("Fetch Data from DB")
		rows, err := db.Query("SELECT name, url FROM provider")
		if err != nil {
			log.Fatalf("Error fetching data: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			p := &Provider{}
			if err := rows.Scan(&p.name, &p.url); err != nil {
				fmt.Printf("Error in scan: %v", err)
				continue
			}
			fmt.Printf("Write to chan: %q\n", p.name)
			providerChan <- p
		}
	}
	close(providerChan) // artificial closure for the demo
}

func useData() {
	for p := range providerChan {
		func() {
			fmt.Printf("Read from chan: %q, %q\n", p.name, p.url)
			resp, err := http.Head(p.url)
			if err != nil {
				fmt.Printf("Error making head request for %q: %v\n", p.url, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Processing Data: %q\t%v\n", p.name, resp)
		}()
	}

}

// useDataWorkerPool will read until the chan is closed but it will read items concurrently and the
// work is controlled with a loop acting as a pool and a WaitGroup to ensure it all finishes before we return.
func useDataWorkerPool() {
	var wg sync.WaitGroup
	concurrencyRate := 5 // in the wild you'd use a config variable for this
	for i := 0; i <= concurrencyRate; i++ {
		fmt.Println("Worker ", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range providerChan {
				func() {
					fmt.Printf("Read from chan: %q, %q\n", p.name, p.url)
					resp, err := http.Head(p.url)
					if err != nil {
						fmt.Printf("Error making head request for %q: %v\n", p.url, err)
						return
					}
					defer resp.Body.Close()
					fmt.Printf("Processing Data: %q\t%v\n", p.name, resp)
				}()
			}
		}()
	}
	wg.Wait()
}