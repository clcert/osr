package asns

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"net/url"
)

const (
	publicServiceSelector = "div.service-item > p > a"
	regionSelector        = "div.region-item > p > a"
	ministrySelector      = "div.ministry-item > p > a"
	buttonSelector        = "div.profile_buttons > a"
	communeSelector       = ".commune > article > a.text-red"
)

func Execute(ctx *tasks.Context) error {
	source := ctx.Sources[0]
	saver := ctx.Savers[0]
	blacklist := getBlacklist(ctx)
	for {
		page := source.Next()
		if page == nil {
			break
		}

		// "Instituciones" root page
		body, err := page.Open()
		if err != nil {
			return err
		}
		defer page.Close()
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return err
		}
		ctx.Log.Info("Importing public services urls...")
		addDomainsFromSelector(doc, publicServiceSelector, blacklist, saver, ctx)
		ctx.Log.Info("Done importing public services urls...")

		// parsing process url
		rootURL, err := url.Parse(page.Path()) // It should be an URL
		if err != nil {
			return err
		}
		ctx.Log.Info("Importing ministry urls...")
		ministrySites := getSelectorURLs(doc, rootURL, ministrySelector)
		for _, site := range ministrySites {
			if err := addDomainsFromURL(site, buttonSelector, blacklist, saver, ctx); err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"url": site.String(),
				}).Error("cannot add domains from ministry site: %s", err)
			}
			ctx.Log.WithFields(logrus.Fields{
				"url": site.String(),
			}).Info("done with this URL!")
		}
		ctx.Log.Info("Done importing ministry urls...")

		ctx.Log.Info("Importing region urls...")
		regionSites := getSelectorURLs(doc, rootURL, regionSelector)
		for _, site := range regionSites {
			if err := addDomainsFromURL(site, buttonSelector, blacklist, saver, ctx); err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"url": site.String(),
				}).Error("cannot add local government domains from ministry site: %s", err)
			}
			ctx.Log.WithFields(logrus.Fields{
				"url": site.String(),
			}).Info("done with the local government URL!")
			if err := addDomainsFromURL(site, communeSelector, blacklist, saver, ctx); err != nil {
				ctx.Log.WithFields(logrus.Fields{
					"url": site.String(),
				}).Error("cannot add domains from region site: %s", err)
			}
			ctx.Log.WithFields(logrus.Fields{
				"url": site.String(),
			}).Info("done with commune URLs!")
		}
		ctx.Log.Info("Done importing region urls...")

	}
	return nil
}
