package asns

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/clcert/osr/models"
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"gopkg.in/tchap/go-patricia.v2/patricia"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func addDomain(domain string, ctx *findContext) error {
	domain = strings.TrimPrefix(domain, "www.") // canonizing domains
	splitDomain := strings.Split(domain, ".")
	tld := splitDomain[len(splitDomain)-1]
	var name, subdomain string
	if len(splitDomain) > 1 {
		name = splitDomain[len(splitDomain)-2]
	}
	if _, ok := ctx.blacklist[name]; ok {
		return nil
	}
	if len(splitDomain) > 2 {
		subdomain = strings.Join(splitDomain[:len(splitDomain)-2], ".")
	}
	if len(tld) > 0 && len(domain) > 0 {
		err := ctx.saver.Save(&models.Domain{
			TaskID:           ctx.taskCtx.GetTaskID(),
			SourceID:         models.CLCERT,
			Subdomain:        subdomain,
			Name:             name,
			TLD:              tld,
			RegistrationDate: time.Now(),
		})
		if err != nil {
			ctx.taskCtx.Log.WithFields(logrus.Fields{
				"domain": domain,
			}).Info("Could not insert domain: %s...", err)
			return err
		}
		err = ctx.saver.Save(&models.DomainToCategory{
			TaskID:             ctx.taskCtx.GetTaskID(),
			SourceID:           models.CLCERT,
			DomainSubdomain:    subdomain,
			DomainName:         name,
			DomainTLD:          tld,
			DomainCategorySlug: "gov",
		})
		if err != nil {
			ctx.taskCtx.Log.WithFields(logrus.Fields{
				"domain": domain,
			}).Info("Could not insert domain category: %s...", err)
		}
	}
	return nil
}

func getURLs(doc *goquery.Document, rootURL *url.URL) (local map[string]*url.URL, external map[string]struct{}) {
	local = make(map[string]*url.URL, 0)
	external = make(map[string]struct{}, 0)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, exists := s.Attr("href"); exists {
			anURL, err := rootURL.Parse(link)
			if err != nil {
				return
			}
			if anURL.Hostname() == rootURL.Hostname() {
				local[anURL.String()] = anURL
			} else {
				external[anURL.Hostname()] = struct{}{}
			}
		}
	})
	return
}

func checkLocalURLs(localURLs map[string]*url.URL, ctx *findContext) {
	for urlStr, url := range localURLs {
		localPath := patricia.Prefix(urlStr)
		if ctx.trie.Match(localPath) {
			continue
		}
		// URL is not in our history, we add it and request it
		ctx.trie.Set(localPath, struct{}{})
		request, err := http.Get(url.String())
		if err != nil {
			ctx.taskCtx.Log.WithFields(logrus.Fields{
				"url": url.String(),
			}).Error("error reading url: %v", err)
			continue
		}
		ctx.stack = append(ctx.stack, &stackItem{
			Closer: request.Body,
			Reader: request.Body,
			url:    url,
		})
	}
}

func getBlacklist(ctx *tasks.Context) map[string]struct{} {
	blacklistMap := make(map[string]struct{})
	blacklist := strings.Split(ctx.Params.Get("blacklist", ""), ",")
	for _, domain := range blacklist {
		blacklistMap[domain] = struct{}{}
	}
	return blacklistMap
}

func processPage(page *stackItem, rootURL *url.URL, ctx *findContext) {
	defer page.Close()
	// Load the HTML document
	ctx.taskCtx.Log.WithFields(logrus.Fields{
		"url": page.url.String(),
	}).Info("Checking page")
	doc, err := goquery.NewDocumentFromReader(page.Reader)
	if err != nil {
		ctx.taskCtx.Log.WithFields(logrus.Fields{
			"url": page.url.String(),
		}).Errorf("error reading file: %v", err)
		return
	}

	// get URLs in HTML document
	localURLs, externalDomains := getURLs(doc, rootURL)
	ctx.taskCtx.Log.WithFields(logrus.Fields{
		"url": page.url.String(),
		"localURLs": len(localURLs),
		"externalDomains": len(externalDomains),
	}).Info("Page parsed")
	if len(externalDomains) > 0 {
		for domain, _ := range externalDomains {
			if err := addDomain(domain, ctx); err != nil {
				ctx.taskCtx.Log.WithFields(logrus.Fields{
					"url": domain,
				}).Errorf("cannot add domain: %s", err)
			}
		}
	}
	if len(localURLs) > 0 {
		checkLocalURLs(localURLs, ctx)
	}
}
