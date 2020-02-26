package asns

import (
	"github.com/clcert/osr/savers"
	"github.com/clcert/osr/tasks"
	"gopkg.in/tchap/go-patricia.v2/patricia"
	"io"
	"net/url"
)

type findContext struct {
	trie      *patricia.Trie
	blacklist map[string]struct{}
	saver     savers.Saver
	taskCtx   *tasks.Context
	stack     []*stackItem
}

type stackItem struct {
	io.Closer
	io.Reader
	url *url.URL
}
