package scrape

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/nujikazo/plmn-list/config"
	cg "github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/database"
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
func Run(generalConf *config.GeneralConf, crawlerConf *cg.CrawlConf, list *[]*database.Schema) error {

	for _, v := range crawlerConf.Plmn {
		env := v.Env
		var res []byte
		var err error

		switch env {
		case "remote":
			ctx := context.Background()
			plmnListURL := v.URL
			req, err := http.NewRequestWithContext(ctx, "GET", plmnListURL, nil)
			client := http.DefaultClient

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			res, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
		case "local":
			res, err = ioutil.ReadFile(v.LocalFile)
			if err != nil {
				return err
			}
		default:
			return errors.New("Environment must be set to 'local' or 'remote")
		}

		doc, err := htmlquery.Parse(bytes.NewReader(res))
		if err != nil {
			return err
		}

		tr := htmlquery.Find(doc, v.Path.Tr)

		for _, t := range tr {
			td := htmlquery.Find(t, v.Path.Td)
			mcc := htmlquery.InnerText(td[0])
			mnc := htmlquery.InnerText(td[1])
			iso := htmlquery.InnerText(td[2])
			country := htmlquery.InnerText(td[3])
			countryCode := htmlquery.InnerText(td[4])
			network := strings.TrimRight(htmlquery.InnerText(td[5]), " ")

			p := &database.Schema{
				MCC:         mcc,
				MNC:         mnc,
				ISO:         iso,
				Country:     country,
				CountryCode: countryCode,
				Network:     network,
			}

			*list = append(*list, p)

		}
	}

	return nil
}
