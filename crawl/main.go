package crawl

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/db"
	"github.com/nujikazo/plmn-list/database/models"
)

const plmnListURL = "https://www.mcc-mnc.com/"

// Run
func Run(conf *config.GeneralConf) error {
	ctx := context.Background()

	resp, err := http.Get(plmnListURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return err
	}

	tr := htmlquery.Find(doc, `//div[@id="main"]/div[@class="content"]/table[@id="mncmccTable"]/tbody/tr`)

	queries, err := db.New(conf)
	if err != nil {
		return err
	}

	for _, t := range tr {
		td := htmlquery.Find(t, `//td`)
		mcc := htmlquery.InnerText(td[0])
		mnc := htmlquery.InnerText(td[1])
		iso := htmlquery.InnerText(td[2])
		country := htmlquery.InnerText(td[3])
		countryCode := htmlquery.InnerText(td[4])
		network := htmlquery.InnerText(td[5])

		_, err := queries.UpsertPlmn(
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
		)
		if err != nil {
			return err
		}
	}

	return nil
}
