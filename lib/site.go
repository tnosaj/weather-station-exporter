package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

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

// the iframe is "https://www.bogner-lehner.eu/lwd/SITENAME_akt.php"
func (c HttpClient) getSiteMetrics(site string, ch chan<- prometheus.Metric) error {
	log.Debugf("Fetching for site - %s", site)
	return nil
}
