package crawl

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/nujikazo/plmn-list/crawl/config"
)

type Plmn struct {
	MCC         string
	MNC         string
	ISO         string
	Country     string
	CountryCode string
	Network     string
}

// Run
func Run(generalConf *config.GeneralConf, crawlerConf *config.PlmnCrawlConf) ([]Plmn, error) {
	env := generalConf.Env
	var res []byte
	var err error

	switch env {
	case "remote":
		plmnListURL := crawlerConf.Plmn.URL
		resp, err := http.Get(plmnListURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	case "local":
		res, err = ioutil.ReadFile(crawlerConf.Plmn.LocalFile)
		if err != nil {
			return nil, err
		}
	}

	doc, err := htmlquery.Parse(bytes.NewReader(res))
	if err != nil {
		return nil, err
	}

	tr := htmlquery.Find(doc, crawlerConf.Plmn.Path.Tr)

	var plmns []Plmn
	for _, t := range tr {
		td := htmlquery.Find(t, crawlerConf.Plmn.Path.Td)
		mcc := htmlquery.InnerText(td[0])
		mnc := htmlquery.InnerText(td[1])
		iso := htmlquery.InnerText(td[2])
		country := htmlquery.InnerText(td[3])
		countryCode := htmlquery.InnerText(td[4])
		network := htmlquery.InnerText(td[5])

		p := Plmn{
			MCC:         mcc,
			MNC:         mnc,
			ISO:         iso,
			Country:     country,
			CountryCode: countryCode,
			Network:     network,
		}

		plmns = append(plmns, p)

	}

	return plmns, nil
}
