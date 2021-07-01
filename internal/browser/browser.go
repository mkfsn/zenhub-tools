package browser

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const (
	bin      = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	queryUrl = "https://api.zenhub.com/v1/graphql?query=getSprintIssues"
)

var (
	homeDir, _  = os.UserHomeDir()
	userDataDir = homeDir + "/Library/Application Support/Google/Chrome"
)

func GetRawSprintIssues(ctx context.Context, url string, profileDir string) ([]byte, error) {
	l := launcher.New().
		Headless(true).
		Leakless(true).
		UserDataDir(userDataDir).
		ProfileDir(profileDir).
		Bin(bin)

	browser := rod.New().
		Timeout(time.Second * 15).
		ControlURL(l.MustLaunch()).
		Trace(true).
		MustConnect()

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
	wait := page.MustWaitRequestIdle()
	wait()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case b := <-dataCh:
		return b, nil
	}
}
