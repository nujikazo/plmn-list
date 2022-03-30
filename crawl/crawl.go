package crawl

import (
	"bytes"
	"context"
	"database/sql"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/db"
	"github.com/nujikazo/plmn-list/database/models"
)

type plmn struct {
	MCC         string
	MNC         string
	ISO         string
	Country     string
	CountryCode string
	Network     string
}

// Run
func Run(generalConf *config.GeneralConf, crawlerConf *config.PlmnCrawlConf) error {
	var res []byte
	var err error

	switch generalConf.Env {
	case "remote":
		plmnListURL := crawlerConf.Plmn.URL
		resp, err := http.Get(plmnListURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	case "local":
		res, err = ioutil.ReadFile("")
		if err != nil {
			return err
		}
	}

	doc, err := htmlquery.Parse(bytes.NewReader(res))
	if err != nil {
		return err
	}

	tr := htmlquery.Find(doc, crawlerConf.Plmn.Path.Tr)

	queries, err := db.New(generalConf)
	if err != nil {
		return err
	}

	ctx := context.Background()
	var plmns []plmn

	for _, t := range tr {
		td := htmlquery.Find(t, crawlerConf.Plmn.Path.Td)
		mcc := htmlquery.InnerText(td[0])
		mnc := htmlquery.InnerText(td[1])
		iso := htmlquery.InnerText(td[2])
		country := htmlquery.InnerText(td[3])
		countryCode := htmlquery.InnerText(td[4])
		network := htmlquery.InnerText(td[5])

		p := plmn{
			MCC:         mcc,
			MNC:         mnc,
			ISO:         iso,
			Country:     country,
			CountryCode: countryCode,
			Network:     network,
		}

		plmns = append(plmns, p)

		if _, err := queries.UpsertPlmn(
			ctx, models.UpsertPlmnParams{
				Mcc:     mcc,
				Mnc:     mnc,
				Iso:     iso,
				Country: country,
				CountryCode: sql.NullString{
					String: countryCode,
					Valid:  true,
				},
				Network:   network,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Mcc_2:     mcc,
				Mnc_2:     mnc,
				Iso_2:     iso,
				Country_2: country,
				CountryCode_2: sql.NullString{
					String: countryCode,
					Valid:  true,
				},
				Network_2:   network,
				UpdatedAt_2: time.Now(),
			},
		); err != nil {
			return err
		}
	}

	return nil
}
