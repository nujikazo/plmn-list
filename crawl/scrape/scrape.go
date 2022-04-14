package scrape

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/general"
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
func Run(generalConf *general.GeneralConf, crawlerConf *config.PlmnCrawlConf) ([]Plmn, error) {
	env := crawlerConf.Plmn.Env
	var res []byte
	var err error

	switch env {
	case "remote":
		ctx := context.Background()
		plmnListURL := crawlerConf.Plmn.URL
		req, err := http.NewRequestWithContext(ctx, "GET", plmnListURL, nil)
		client := http.DefaultClient

		resp, err := client.Do(req)
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
		network := strings.TrimRight(htmlquery.InnerText(td[5]), " ")

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
