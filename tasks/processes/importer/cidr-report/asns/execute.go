package asns

import (
	"bufio"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"regexp"
	"strconv"
)

// Regex used for scrapping lines from webpage.
const lineRegex = `<a href=".*">(AS\d+)\s*<\/a>\s*(.*), ([A-Z]{2})`

// Downloads the CIDR webpage and uploads its data to the
// ASNs database.
func Execute(args *tasks.Context) error {
	source := args.Sources[0]
	saver := args.Savers[0]
	re := regexp.MustCompile(lineRegex)

	for {
		page := source.Next()

		if page == nil {
			break
		}

		body, err := page.Open()
		if err != nil {
			return err
		}

		// Page is encoded in Windows1252 :(
		dec := transform.NewReader(body, charmap.Windows1252.NewDecoder())

		scanner := bufio.NewScanner(dec)
		args.Log.Info("Parsing ASNs...")
		for scanner.Scan() {
			result := re.FindStringSubmatch(scanner.Text())
			if len(result) >= 4 && len(result[1]) > 2 {
				asnID, err := strconv.Atoi(result[1][2:])
				if err == nil {
					err = saver.Save(&models.ASN{
						ID:            asnID,
						Name:          result[2],
						CountryAlpha2: result[3],
					})
					if err != nil {
						args.Log.Info("Could not insert asn: %s...", err)
					}
				}
			}
		}
	}
	return nil
}
