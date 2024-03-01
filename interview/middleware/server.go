package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type response struct {
	Message string `json:"message"`
	Array   []int  `json:"array"`
}

func Returning(msg string) (string, error) {
	return msg, nil
}

func oddDouble(arr []int, wg *sync.WaitGroup, err chan error) {
	defer func() {
		close(err)
		wg.Done()
	}()
	for i := 0; i < len(arr); i++ {
		if i%2 != 0 {
			arr[i] *= 2
		}
	}
	err <- nil
}

func evenTriple(arr []int, wg *sync.WaitGroup, err chan error) {
	defer func() {
		close(err)
		wg.Done()
	}()
	for i := 0; i < len(arr); i++ {
		if i%2 == 0 {
			arr[i] *= 3
		}
	}
	err <- nil
}

func Allow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		if err := json.NewEncoder(w).Encode(response{Message: "Bad Request"}); err != nil {
			log.Fatal(err)
		}
		return
	}

	var wg sync.WaitGroup
	arr := []int{2, 3, 5, 4, 1, 8}
	errChanEven := make(chan error, 0)
	errChanOdd := make(chan error, 0)

	wg.Add(2)
	go oddDouble(arr, &wg, errChanOdd)
	go evenTriple(arr, &wg, errChanEven)

	if err := <-errChanEven; err != nil {
		log.Fatal(err)
	}

	if err := <-errChanOdd; err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	if err := json.NewEncoder(w).Encode(response{Message: "Hello There", Array: arr}); err != nil {
		log.Fatal(err)
	}
}
