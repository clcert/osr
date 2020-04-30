package gob_domains

import (
	"github.com/clcert/osr/tasks"
	"github.com/sirupsen/logrus"
	"gopkg.in/tchap/go-patricia.v2/patricia"
	"net/url"
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

		// root stackItem
		body, err := page.Open()
		if err != nil {
			return err
		}

		// parsing process url
		rootURL, err := url.Parse(page.Path()) // It should be an URL
		if err != nil {
			return err
		}

		// creating Context
		findCtx := &findContext{
			taskCtx:   ctx,
			trie:      patricia.NewTrie(),
			saver:     saver,
			blacklist: blacklist,
			stack:     make([]*stackItem, 0),
		}
		findCtx.stack = append(findCtx.stack, &stackItem{Closer: page, Reader: body, url: rootURL})
		findCtx.trie.Set(patricia.Prefix(rootURL.String()), struct{}{})
		checked := 0
		for len(findCtx.stack) > 0 {
			item := findCtx.stack[len(findCtx.stack) - 1]
			findCtx.stack = findCtx.stack[:len(findCtx.stack) - 1]
			processPage(item, rootURL, findCtx)
			checked++
		}
		findCtx.taskCtx.Log.WithFields(logrus.Fields{
			"rootURL": rootURL.String(),
		}).Infof("Checked %d pages", checked)
	}
	return nil
}
