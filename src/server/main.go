package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view  int
	click int
}

func (c *counters) incView() {
	c.Lock()
	defer c.Unlock()
	c.view++
}

func (c *counters) incClick() {
	c.Lock()
	defer c.Unlock()
	c.click++
}

var (
	c         = counters{}
	statsMax  = 10
	dataStore = make(map[string]*counters)
	content   = []string{"sports", "entertainment", "business", "education"}
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func getKey() string {
	data := content[rand.Intn(len(content))]
	dt := time.Now()
	formatedTime := dt.Format("01-02-2006 15:04:05")
	key := data + ":" + formatedTime
	return key
}

func getCounter() *counters {
	key := getKey()
	counter := dataStore[key]
	if counter == nil {
		counter = &counters{}
		dataStore[key] = counter
	}
	log.Println(key)
	log.Println(counter)
	return counter
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	counter := getCounter()
	counter.incView()
	fmt.Println(counter)

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		counter.incClick()
	}
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed() {
		w.WriteHeader(429)
		return
	}
	c.Lock()
	c.view++
	c.Unlock()
}

func isAllowed() bool {
	return c.view <= statsMax
}

func uploadCounters(dt int64) error {
	for {
		counter := getCounter()
		counter.incView()
		counter.incClick()
		log.Println(counter)
		time.Sleep(time.Duration(dt) * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)
	go uploadCounters(5000)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
