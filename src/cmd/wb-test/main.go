package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func init() {
	log.SetFlags(0)
}

func main() {

	var (
		count int
		total int
		uri   string
		mux   sync.Mutex
	)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// Get url from stdin.
		uri = scanner.Text()

		// Check if user finished typing.
		if uri == "" {
			break
		}
		// uri is valid
		_, err := url.ParseRequestURI(uri)
		if err != nil {
			fmt.Println("Invalid URL, please try again.")
			// log.Fatalln(err)
			continue
		}

		go func() {
			// Do http GET request.
			resp, err := http.Get(uri)
			if err != nil {
				log.Fatalln(err)
			}
			// Close response body.
			defer resp.Body.Close()
			// Read response body.
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			mux.Lock()
			// Only one goroutine at a time can access the count and total vars.
			count = bytes.Count(body, []byte("Go"))
			total += count
			mux.Unlock()
			log.Printf("Count for %s %d", uri, count)
		}()

	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Total: %v", total)
}
