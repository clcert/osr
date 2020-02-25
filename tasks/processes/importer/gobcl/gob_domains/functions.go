package asns

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func addDomain(link *url.URL, saver savers.Saver, ctx *tasks.Context) error {
	domain := link.Hostname()
	domain = strings.TrimPrefix(domain, "www.") // canonizing domains
	splitDomain := strings.Split(domain, ".")
	tld := splitDomain[len(splitDomain)-1]
	var name, subdomain string
	if len(splitDomain) > 1 {
		name = splitDomain[len(splitDomain)-2]
	}
	if len(splitDomain) > 2 {
		subdomain = splitDomain[len(splitDomain)-3]
	}
	err := saver.Save(&models.Domain{
		TaskID:           ctx.GetTaskID(),
		SourceID:         models.CLCERT,
		Subdomain:        subdomain,
		Name:             name,
		TLD:              tld,
		RegistrationDate: time.Now(),
	})
	if err != nil {
		return err
	}
	err = saver.Save(&models.DomainToCategory{
		TaskID:             ctx.GetTaskID(),
		SourceID:           models.CLCERT,
		DomainSubdomain:    subdomain,
		DomainName:         name,
		DomainTLD:          tld,
		DomainCategorySlug: "gov",
	})
	if err != nil {
		return err
	}
	if err != nil {
		ctx.Log.Info("Could not insert domain: %s...", err)
	}
	return nil
}

func getSelectorURLs(doc *goquery.Document, rootURL *url.URL, selector string) (sites []*url.URL) {
	sites = make([]*url.URL, 0)
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			fullURL, err := rootURL.Parse(link)
			if err != nil {
				return
			}
			sites = append(sites, fullURL)
		}
	})
	return
}

func addDomainsFromURL(url *url.URL, selector string, saver savers.Saver, ctx *tasks.Context) error {
	request, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer request.Body.Close()
	doc, err := goquery.NewDocumentFromReader(request.Body)
	if err != nil {
		return err
	}
	addDomainsFromSelector(doc, selector, saver, ctx)
	return nil
}


func addDomainsFromSelector(doc *goquery.Document, selector string, saver savers.Saver, ctx *tasks.Context) {
	// Getting public services links first
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			url, err := url.Parse(link)
			if err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"url": url,
				}).Error("error parsing domain: %s", err)
			}
			if err := addDomain(url, saver, ctx); err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"url": url,
				}).Error("error parsing domain: %s", err)
			}
		}
	})
}
