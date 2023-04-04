package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Number []string `yaml:"number"`
	URL    string   `yaml:"url"`
}
type request struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Env       string `json:"env"`
			Instance  string `json:"instance"`
			Job       string `json:"job"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Summary string `json:"summary"`
		} `json:"annotations"`
		StartsAt     time.Time `json:"startsAt"`
		EndsAt       time.Time `json:"endsAt"`
		GeneratorURL string    `json:"generatorURL"`
		Fingerprint  string    `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
		Job       string `json:"job"`
	} `json:"groupLabels"`
	CommonLabels struct {
		Alertname string `json:"alertname"`
		Env       string `json:"env"`
		Job       string `json:"job"`
		Severity  string `json:"severity"`
	} `json:"commonLabels"`
	CommonAnnotations struct {
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

type Payload struct {
	Message     string `json:"Message"`
	PhoneNumber string `json:"PhoneNumber"`
}

var InfoLogger *log.Logger
var WarningLogger *log.Logger
var ErrorLogger *log.Logger

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func requestreadr(w http.ResponseWriter, resp *http.Request) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorLogger.Println("Some errors in reading request body")
		log.Fatal(err)
	}
	InfoLogger.Println("Get successfylly alertmanager request")

	var req request
	err = json.Unmarshal(body, &req)
	if err != nil {
		ErrorLogger.Println("Some errors in working with json file and unmarshaling incomming request data")
		log.Fatal(err)
	}
	InfoLogger.Println("Unmarashal successfully request data")

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		ErrorLogger.Println("Some errors in reading config file include URL or Numbers")
		log.Fatal(err)
		return
	}
	InfoLogger.Println("Read successfully config file data's")

	// Create a struct to hold the YAML data
	var config Config

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		ErrorLogger.Println("Some errors in working with YAML file and unmarshaling config file data")
		return
	}
	InfoLogger.Println("Unmarshall successfully config file's data")

	sms_url := config.URL
	for _, number := range config.Number {
		date_time := ptime.Now().Format("yyyy-MM-dd HH-mm-ss")
		message := "Platform2\n\nStatus: " + req.Status + "\nSummary: " + req.Alerts[0].Annotations.Summary + "\nDate and Time: " + date_time
		data := Payload{
			PhoneNumber: number,
			Message:     message,
		}
		alertname := req.Alerts[0].Labels.Alertname
		payloadBytes, err := json.Marshal(data)
		if err != nil {
			ErrorLogger.Println("We have some error in unmarshalling json data")
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", sms_url, body)
		if err != nil {
			ErrorLogger.Println("we have some error for create a request")
		}
		req.Header.Set("Content-Type", "application/json")
		InfoLogger.Println("request was created")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ErrorLogger.Println("some errors accured in sending http request")
		}
		defer resp.Body.Close()
		InfoLogger.Println("We have sent alert ", alertname, " sms to", number, "at", date_time)
	}
}

func main() {
	http.HandleFunc("/", requestreadr)
	http.ListenAndServe(":8040", nil)
}
