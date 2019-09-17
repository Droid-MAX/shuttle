package dns

import (
	"flag"
	"net"
	"os"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/assets"
	"github.com/sirupsen/logrus"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var fileName = flag.String("geoip", os.Getenv("GEOIP_DB"), "geo ip db path")

var geoipDB *geoip2.Reader

func InitGeoIP() error {
	var err error
	geoipFileBytes, err := assets.ReadFile(*fileName)
	if err != nil {
		return errors.Errorf("reade geo file [%s] failed: %s", *fileName, err.Error())
	}
	geoipDB, err = geoip2.FromBytes(geoipFileBytes)
	if err != nil {
		return errors.Errorf("reade geo file [%s] failed: %s", *fileName, err.Error())
	}
	return nil
}

func GeoLookUp(ip net.IP) (countryCode string) {
	if geoipDB == nil {
		return
	}
	country, err := geoipDB.Country(ip)
	if err == nil && country != nil {
		logrus.WithField("country-code", country.Country.IsoCode).WithField("ip", ip.String()).
			Debug("GeoIP lookup")
		return country.Country.IsoCode
	}
	logrus.Errorf("[GeoIP] lookup [%s] country failed: %s", ip, err.Error())
	return
}

func CloseGeoDB() error {
	return geoipDB.Close()
}
