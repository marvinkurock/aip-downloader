package aip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"time"
)

type ManifestJson struct {
  Name string `json:"name"`
  Abbreviation string `json:"abbreviation"`
  Version float64 `json:"version"`
  OrganizationName string `json:"organizationName"`
  EffectiveDate string `json:"effectiveDate"`
}

type AirportConfig struct {
	M    string `json:"m"`
	N    string `json:"n"`
	Icao string
}

func SetManifestJson(){
  manifestFile, err := os.Open("manifest.json")
  if err != nil {
    panic(err)
  }
  defer manifestFile.Close()
  bytesValue, _ := ioutil.ReadAll(manifestFile)

  var manifest ManifestJson
  json.Unmarshal(bytesValue, &manifest)

  now := time.Now().UTC()
  manifest.EffectiveDate = now.Format("20060102T15:04:05Z")
  manifest.Version += 0.1
  manifest.Version = math.Round(manifest.Version * 100) / 100

  newManifest, _ := json.Marshal(manifest)
  os.WriteFile("manifest.json", newManifest, 0644)
  os.WriteFile("charts/manifest.json", newManifest, 0644)
}

func ParseConfigJson() []AirportConfig {
	jsonFile, err := os.Open("airports.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var airports []AirportConfig
	var icaoAirports []AirportConfig
	json.Unmarshal(byteValue, &airports)
	r, _ := regexp.Compile("[A-Z]{4}$")
	for _, ap := range airports {
		ap.Icao = r.FindString(ap.N)
		if ap.Icao != "" {
			fmt.Println(ap.N, ": ", ap.Icao)
      icaoAirports = append(icaoAirports, ap)
		}
	}
  return icaoAirports
}
