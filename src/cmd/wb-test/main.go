package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Handler struct {
	total  int
	mu     *sync.Mutex
	client *http.Client
}

func (h *Handler) Total() int {
	return h.total
}

func (h *Handler) FindInUrl(ctx context.Context, u *url.URL) error {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errors.New("bad status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	count := bytes.Count(body, []byte("Go"))
	log.Printf("For %s found: %d\n", u.String(), count)
	h.mu.Lock()
	h.total += count
	h.mu.Unlock()

	return nil
}

func main() {
	var (
		n      = 0
		max    = 5
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	)

	handler := &Handler{client: client, mu: &sync.Mutex{}}

	eg, ctx := errgroup.WithContext(context.Background())
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}

		u, err := url.Parse(scanner.Text())
		if err != nil {
			log.Fatal(err)
			continue
		}
		if n == max {
			break
		}

		eg.Go(func() error { return handler.FindInUrl(ctx, u) })

		n++

	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Total found: %d", handler.Total())
}
