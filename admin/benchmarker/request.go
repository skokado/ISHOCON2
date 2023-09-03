package main

// 初期化(N秒以内)
import (
	"crypto/tls"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/http2"
)

func getInitialize() {
	log.Print("Start GET /initialize")
	finishTime := time.Now().Add(10 * time.Second)
	httpRequest("GET", "/initialize", nil)
	if time.Now().Sub(finishTime) > 0 {
		log.Print("Timeover at GET /initialize")
		os.Exit(1)
	}
}

func postVote(v Vote) bool {
	prm := url.Values{}
	prm.Add("name", v.Name)
	prm.Add("address", v.Address)
	prm.Add("mynumber", v.Mynumber)
	prm.Add("candidate", v.Candidate)
	prm.Add("keyword", v.Keyword)
	prm.Add("vote_count", v.VoteCount)
	if httpRequest("POST", "/vote", prm) == 200 {
		return true
	}
	return false
}

func getIndex() bool {
	if httpRequest("GET", "/", nil) == 200 {
		return true
	}
	return false
}

func getCandidate() bool {
	id := strconv.Itoa(getRand(1, 30))
	if httpRequest("GET", "/candidates/"+id, nil) == 200 {
		return true
	}
	return false
}

func getPoliticalParty() bool {
	set := []string{"国民元気党", "国民10人大活躍党", "夢実現党", "国民平和党"}
	party := set[getRand(0, 3)]
	if httpRequest("GET", "/political_parties/"+party, nil) == 200 {
		return true
	}
	return false
}

func getCSS() bool {
	if httpRequest("GET", "/css/bootstrap.min.css", nil) == 200 {
		return true
	}
	return false
}

var clients []http.Client

func createClients(size int) {
	clients = make([]http.Client, size)
	for i := 0; i < size; i++ {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		if err := http2.ConfigureTransport(tr); err != nil {
			log.Fatalf("Failed to configure h2 transport: %s", err)
		}
		clients[i] = http.Client{Transport: tr}
	}
	rand.Seed(time.Now().Unix())
}

func httpRequest(method string, path string, params url.Values) int {
	req, _ := http.NewRequest(method, origin+path, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := clients[rand.Intn(len(clients))]

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return 500
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func httpRequestDoc(method string, path string, params url.Values) *goquery.Document {
	req, _ := http.NewRequest(method, origin+path, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := clients[rand.Intn(len(clients))]

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	return doc
}
