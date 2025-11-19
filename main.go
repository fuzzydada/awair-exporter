package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML configuration file
type Config struct {
	Hosts      []string `yaml:"hosts"`
	ListenPort int      `yaml:"listen_port"`
}

// AwairData represents the structure of the JSON response from the Awair Local API
type AwairData struct {
	Timestamp      time.Time `json:"timestamp"`
	Score          float64   `json:"score"`
	DewPoint       float64   `json:"dew_point"`
	Temp           float64   `json:"temp"`
	Humid          float64   `json:"humid"`
	AbsHumid       float64   `json:"abs_humid"`
	CO2            float64   `json:"co2"`
	CO2Est         float64   `json:"co2_est"`
	CO2EstBaseline float64   `json:"co2_est_baseline"`
	VOC            float64   `json:"voc"`
	VOCBaseline    float64   `json:"voc_baseline"`
	VOCH2Raw       float64   `json:"voc_h2_raw"`
	VOCEthanolRaw  float64   `json:"voc_ethanol_raw"`
	PM25           float64   `json:"pm25"`
	PM10Est        float64   `json:"pm10_est"`
}

// AwairCollector implements the prometheus.Collector interface
type AwairCollector struct {
	hosts          []string
	httpClient     *http.Client
	score          *prometheus.Desc
	dewPoint       *prometheus.Desc
	temperature    *prometheus.Desc
	humidity       *prometheus.Desc
	absHumidity    *prometheus.Desc
	co2            *prometheus.Desc
	co2Est         *prometheus.Desc
	co2EstBaseline *prometheus.Desc
	voc            *prometheus.Desc
	vocBaseline    *prometheus.Desc
	vocH2Raw       *prometheus.Desc
	vocEthanolRaw  *prometheus.Desc
	pm25           *prometheus.Desc
	pm10Est        *prometheus.Desc
}

// NewAwairCollector creates a new AwairCollector
func NewAwairCollector(hosts []string) *AwairCollector {
	return &AwairCollector{
		hosts: hosts,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		score: prometheus.NewDesc("awair_score",
			"Awair score (0-100)",
			[]string{"device"}, nil,
		),
		dewPoint: prometheus.NewDesc("awair_dew_point_celsius",
			"Dew point in Celsius",
			[]string{"device"}, nil,
		),
		temperature: prometheus.NewDesc("awair_temperature_celsius",
			"Temperature in Celsius",
			[]string{"device"}, nil,
		),
		humidity: prometheus.NewDesc("awair_humidity_percent",
			"Relative humidity percentage (0-100)",
			[]string{"device"}, nil,
		),
		absHumidity: prometheus.NewDesc("awair_absolute_humidity_g_m3",
			"Absolute humidity in grams per cubic meter",
			[]string{"device"}, nil,
		),
		co2: prometheus.NewDesc("awair_co2_ppm",
			"Carbon Dioxide in parts per million",
			[]string{"device"}, nil,
		),
		co2Est: prometheus.NewDesc("awair_co2_estimated_ppm",
			"Estimated Carbon Dioxide in parts per million",
			[]string{"device"}, nil,
		),
		co2EstBaseline: prometheus.NewDesc("awair_co2_estimated_baseline",
			"CO2 sensor baseline value for estimation algorithm",
			[]string{"device"}, nil,
		),
		voc: prometheus.NewDesc("awair_voc_ppb",
			"Volatile Organic Compounds in parts per billion",
			[]string{"device"}, nil,
		),
		vocBaseline: prometheus.NewDesc("awair_voc_baseline",
			"VOC sensor baseline value for estimation algorithm",
			[]string{"device"}, nil,
		),
		vocH2Raw: prometheus.NewDesc("awair_voc_h2_raw",
			"Raw H2 sensor value for VOC calculation",
			[]string{"device"}, nil,
		),
		vocEthanolRaw: prometheus.NewDesc("awair_voc_ethanol_raw",
			"Raw ethanol sensor value for VOC calculation",
			[]string{"device"}, nil,
		),
		pm25: prometheus.NewDesc("awair_pm25_ug_m3",
			"Particulate Matter (2.5 microns) in micrograms per cubic meter",
			[]string{"device"}, nil,
		),
		pm10Est: prometheus.NewDesc("awair_pm10_estimated_ug_m3",
			"Estimated Particulate Matter (10 microns) in micrograms per cubic meter",
			[]string{"device"}, nil,
		),
	}
}

// Describe implements prometheus.Collector
func (collector *AwairCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.score
	ch <- collector.dewPoint
	ch <- collector.temperature
	ch <- collector.humidity
	ch <- collector.absHumidity
	ch <- collector.co2
	ch <- collector.co2Est
	ch <- collector.co2EstBaseline
	ch <- collector.voc
	ch <- collector.vocBaseline
	ch <- collector.vocH2Raw
	ch <- collector.vocEthanolRaw
	ch <- collector.pm25
	ch <- collector.pm10Est
}

// Collect implements prometheus.Collector
func (collector *AwairCollector) Collect(ch chan<- prometheus.Metric) {
	for _, host := range collector.hosts {
		data, err := collector.fetchData(host)
		if err != nil {
			log.Printf("Error fetching data from %s: %v", host, err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(collector.score, prometheus.GaugeValue, data.Score, host)
		ch <- prometheus.MustNewConstMetric(collector.dewPoint, prometheus.GaugeValue, data.DewPoint, host)
		ch <- prometheus.MustNewConstMetric(collector.temperature, prometheus.GaugeValue, data.Temp, host)
		ch <- prometheus.MustNewConstMetric(collector.humidity, prometheus.GaugeValue, data.Humid, host)
		ch <- prometheus.MustNewConstMetric(collector.absHumidity, prometheus.GaugeValue, data.AbsHumid, host)
		ch <- prometheus.MustNewConstMetric(collector.co2, prometheus.GaugeValue, data.CO2, host)
		ch <- prometheus.MustNewConstMetric(collector.co2Est, prometheus.GaugeValue, data.CO2Est, host)
		ch <- prometheus.MustNewConstMetric(collector.co2EstBaseline, prometheus.GaugeValue, data.CO2EstBaseline, host)
		ch <- prometheus.MustNewConstMetric(collector.voc, prometheus.GaugeValue, data.VOC, host)
		ch <- prometheus.MustNewConstMetric(collector.vocBaseline, prometheus.GaugeValue, data.VOCBaseline, host)
		ch <- prometheus.MustNewConstMetric(collector.vocH2Raw, prometheus.GaugeValue, data.VOCH2Raw, host)
		ch <- prometheus.MustNewConstMetric(collector.vocEthanolRaw, prometheus.GaugeValue, data.VOCEthanolRaw, host)
		ch <- prometheus.MustNewConstMetric(collector.pm25, prometheus.GaugeValue, data.PM25, host)
		ch <- prometheus.MustNewConstMetric(collector.pm10Est, prometheus.GaugeValue, data.PM10Est, host)
	}
}

// fetchData fetches and parses air quality data from a single Awair device
func (collector *AwairCollector) fetchData(host string) (*AwairData, error) {
	url := "http://" + host + "/air-data/latest"
	resp, err := collector.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data AwairData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// loadConfig loads configuration from file or environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{
		ListenPort: 9101, // Default port
	}
	
	// Attempt to load from config file first
	configFile := "/config/config.yml"
	if _, err := os.Stat(configFile); err == nil {
		log.Printf("Loading configuration from %s", configFile)
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(yamlFile, cfg)
		if err != nil {
			return nil, err
		}
	}

	// Override with environment variables if they are set
	if awairHosts := os.Getenv("AWAIR_HOSTS"); awairHosts != "" {
		cfg.Hosts = strings.Split(awairHosts, ",")
	}
	if listenPort := os.Getenv("LISTEN_PORT"); listenPort != "" {
		port, err := strconv.Atoi(listenPort)
		if err != nil {
			return nil, err
		}
		cfg.ListenPort = port
	}

	if len(cfg.Hosts) == 0 {
		log.Fatal("Configuration error: No Awair hosts specified. Please set AWAIR_HOSTS or define hosts in /config/config.yml")
	}

	return cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	collector := NewAwairCollector(cfg.Hosts)
	prometheus.MustRegister(collector)

	listenAddr := ":" + strconv.Itoa(cfg.ListenPort)
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting server on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
