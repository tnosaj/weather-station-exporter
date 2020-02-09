package lib

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/html"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// the iframe is "https://www.bogner-lehner.eu/lwd/SITENAME_akt.php"
func getSiteMetrics(site string, ch chan<- prometheus.Metric) error {
	log.Debugf("Fetching for site - %s", site)
	url := fmt.Sprintf("https://www.bogner-lehner.eu/lwd/%s_akt.php", site)
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Failed to get page for site: %s - with error %s", site, err)
	}

	siteData := getSiteData(resp)

	log.Debugf("Found %d divs to parse", len(siteData))
	log.Debugf("Contents: %q", siteData)

	var metricsForSite []promMetric

	for key, val := range siteData {
		newMetric, merr := makeMetrics(key, val)
		// TODO: Add outdated date check here and reset metricsForSite to empty, and return error
		if merr != nil {
			log.Errorf("Metric could not be created: %s", merr)
		}
		newMetric.LabelDesc = append(newMetric.LabelDesc, "site")
		newMetric.Label = append(newMetric.Label, site)
		metricsForSite = append(metricsForSite, newMetric)
	}
	sendMetrics(metricsForSite, ch)
	return nil
}

func getSiteData(resp *http.Response) map[string]string {
	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	readText := false
	var previousId string
	var check bool
	siteData := make(map[string]string)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// End of the document, we're done
			log.Debug("Read the last token")
			return siteData
		case html.StartTagToken:
			check, previousId = isDivOfInterest(z.Token())
			if check {
				readText = true
			}
		case html.TextToken:
			if readText {
				val := string(z.Text())
				log.Debugf("Keeping: %s - %s", previousId, val)
				siteData[previousId] = val
				readText = false
				previousId = ""
			}
		}
	}
}

func isDivOfInterest(token html.Token) (bool, string) {
	var divName string
	if token.Data != "div" {
		return false, ""
	}
	for _, a := range token.Attr {
		log.Debugf("attributes: %s", a)
		divName = a.Val
	}
	switch divName {
	case
		"rahmen30",
		"rahmen31",
		"rahmen32",
		"rahmen33",
		"rahmen34",
		"rahmen35",
		"rahmen36",
		"rahmen70":
		return true, divName
	}
	return false, ""
}

func sendMetrics(metrics []promMetric, ch chan<- prometheus.Metric) {
	for _, metric := range metrics {
		log.Debugf("Pushing metric: %s, with values: %q", metric.Name, metric)
		if metric.Name != "" {
			ch <- prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					metric.Name,
					metric.Desc,
					metric.LabelDesc, nil),
				prometheus.GaugeValue,
				metric.Value,
				metric.Label...)
		}
	}
}

//
//Lufttemperatur [C] -3.5
//Luftfeuchtigkeit [%] 32
//Windgeschwindigkeit [km/h] 1.1
//Boe [km/h] 2.9
//Windrichtung [Grad] E
//Schneehoehe [cm] 139
//Temperatur Oberflaeche [C] -10.1
//
//
//the values are saved in divs like this:
//<div id="rahmen30">-2.6</div>
//<div id="rahmen31">33</div>
//<div id="rahmen32">2.9</div>
//<div id="rahmen33">5.4</div>
//<div id="rahmen34">N</div>
//<div id="rahmen35">135</div>
//<div id="rahmen36">-6.7</div>
//

func makeMetrics(id string, val string) (promMetric, error) {
	var m promMetric
	var err error
	switch id {
	case "rahmen30":
		m.Name = "weather_station_air_temperature"
		m.Desc = "Air temperature"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen31":
		m.Name = "weather_station_air_humidity"
		m.Desc = "Air humidity"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen32":
		m.Name = "weather_station_wind_speed"
		m.Desc = "Wind speed"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen33":
		m.Name = "weather_station_wind_gust_speed"
		m.Desc = "wind gust speeds"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen34":
		m.Name = "weather_station_wind_direction"
		m.Desc = "wind direction in degrees"
		m.Value, err = windToDeg(val)
	case "rahmen35":
		m.Name = "weather_station_snow_hight"
		m.Desc = "hight of snow"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen36":
		m.Name = "weather_station_surface_temperature"
		m.Desc = "temperature on the snow surface"
		m.Value, err = strconv.ParseFloat(val, 64)
	case "rahmen70":
		// TODO: check if time is up2date
		return m, nil
	case "":
		return m, nil
	default:
		err = errors.New(fmt.Sprintf("unknown div id received: %s with value: %s", id, val))
	}

	return m, err
}

// guestimations
func windToDeg(dir string) (float64, error) {
	switch dir {
	case "N":
		return float64(0), nil
	case "NE":
		return float64(45), nil
	case "E":
		return float64(90), nil
	case "SE":
		return float64(135), nil
	case "S":
		return float64(180), nil
	case "SW":
		return float64(225), nil
	case "W":
		return float64(270), nil
	case "NW":
		return float64(315), nil
	}
	return 0, errors.New("unknown wind direction")
}
