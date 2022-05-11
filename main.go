package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type InfoData struct {
	HostIp         string
	ClientIp       string
	ClientLocation string
}

type JsonData struct {
	Status string
}

type ResponseData struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func main() {
	fs := http.FileServer(http.Dir("assets/"))

	tmpl := template.Must(template.ParseFiles("index.html"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data, _ := json.Marshal(JsonData{
			Status: "OK",
		})
		fmt.Fprintf(w, string(data))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		userIp := r.Header.Get("X-Forwarded-For")

		if userIp == "" {
			userIp = ip
		}
		loc := getClientAddress(userIp)
		fmt.Println(r.URL)
		data := InfoData{
			HostIp:         getHostIp(),
			ClientIp:       userIp,
			ClientLocation: strings.Trim(loc, " "),
		}
		tmpl.Execute(w, data)
	})

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	server.ListenAndServe()
	fmt.Println("Listening at port 8080")
}

func getClientAddress(ip string) string {
	urlstrings := []string{"https://ipinfo.io/", ip, "/json"}
	url := strings.Join(urlstrings, "")
	response, err := http.Get(url)

	if err != nil {
		return "?"
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "?"
	}
	var responseObj ResponseData
	json.Unmarshal(responseData, &responseObj)
	loc := []string{responseObj.Country, responseObj.City}
	return strings.Join(loc, " ")
}

func getHostIp() string {
	response, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return "?"
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "?"
	}
	var responseObj ResponseData
	json.Unmarshal(responseData, &responseObj)
	return responseObj.IP
}
