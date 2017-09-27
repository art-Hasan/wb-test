package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func init() {
	log.SetFlags(0)
}

func HandleUri(uri string, total *int, wg *sync.WaitGroup, mux *sync.Mutex) {
	// Decrement the counter when the goroutine completes.
	defer wg.Done()
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
	count := bytes.Count(body, []byte("Go"))
	mux.Lock()
	// Only one goroutine at a time can access the total variable.
	*total += count
	mux.Unlock()

	log.Printf("Count for %s %d", uri, count)
}

func ValidUrl(uri string) {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	var (
		total int
		uri   string
		wg    sync.WaitGroup
		mux   sync.Mutex
	)

	// Read stdin.
	_bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	// Convert bytes to string.
	toStr := string(_bytes)

	// Split string by \n
	splited := strings.Split(toStr, "\n")
	// Delete empty strings from a slice.
	for i := 0; i < len(splited); i++ {
		if splited[i] == "" {
			splited = append(splited[:i], splited[i+1:]...)
		}
	}

	for i := 0; i < len(splited); i++ {
		uri = splited[i]
		// uri is valid
		ValidUrl(uri)
		// Increment the WaitGroup counter.
		wg.Add(1)
		go HandleUri(uri, &total, &wg, &mux)
	}

	// Wait for complete goroutines.
	wg.Wait()
	log.Printf("Total: %v", total)
}
