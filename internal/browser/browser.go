package browser

import (
	"context"
	"net/http"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const (
	queryUrl = "https://api.zenhub.com/v1/graphql?query=getSprintIssues"
)

func GetRawSprintIssues(ctx context.Context, url string) ([]byte, error) {
	l := launcher.NewUserMode()

	browser := rod.New().
		ControlURL(l.MustLaunch()).
		Trace(true).
		MustConnect()

	defer browser.MustClose()

	dataCh := make(chan []byte)

	router := browser.HijackRequests()
	router.MustAdd(queryUrl, func(ctx *rod.Hijack) {
		ctx.MustLoadResponse()
		if ctx.Request.Method() != http.MethodPost {
			return
		}

		dataCh <- []byte(ctx.Response.Body())
	})
	go router.Run()

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, err
	}

	if err := (proto.NetworkSetBypassServiceWorker{Bypass: true}).Call(page); err != nil {
		return nil, err
	}

	page.MustNavigate(url)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b := <-dataCh:
		return b, nil
	}
}
