package know

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoogleProvider struct {
	Debugf func(format string, v ...interface{})
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{nil}
}

func (p *GoogleProvider) Name() string {
	return "Google"
}

func (p *GoogleProvider) debug(format string, v ...interface{}) {
	if p.Debugf == nil {
		return
	}
	p.Debugf(format, v...)
}

func (p *GoogleProvider) Ask(question string) (*Answer, error) {
	u, err := url.Parse("https://www.google.com/search")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("q", question)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Unexpected status: " + resp.Status)
	}

	defer resp.Body.Close()
	ans, err := p.parse(resp.Body)
	if ans != nil && ans.Answer == "" {
		return nil, err
	}
	return ans, err
}

func (p *GoogleProvider) parse(body io.Reader) (*Answer, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	ans := &Answer{}
	if block := doc.Find(".kp-blk"); block.Length() > 0 {
		p.parseFreebase(block.Eq(0), ans)
	} else if knavi := doc.Find(".knavi"); knavi.Length() == 1 {
		p.parseKnavi(knavi, ans)
	}

	return ans, nil
}

func (p *GoogleProvider) parseFreebase(blk *goquery.Selection, ans *Answer) {
	p.debug("freebase")
	if ctxs := blk.Find(".kno-fb-ctx"); ctxs.Length() > 0 {
		for i := 0; i < ctxs.Length(); i++ {
			ctx := ctxs.Eq(i)
			if ctx.ParentsFiltered("#media_result_group").Length() > 0 {
				p.debug("media")
				if media := p.parseMedia(ctx); media != nil {
					ans.Media = append(ans.Media, media)
				}
			} else if ans.Answer == "" {
				divs := ctx.Find("div")
				ans.Answer = divs.Eq(0).Text()
				ans.Question = divs.Eq(1).Text()
			}
		}
		return
	}

}

func (p *GoogleProvider) parseKnavi(knavi *goquery.Selection, ans *Answer) {
	p.debug("knavi")
	// Calculator widget
	if cwos := knavi.Find("#cwos"); cwos.Length() == 1 {
		p.debug("calculator")
		ans.Answer = strings.TrimSpace(cwos.Text())
		if q := knavi.Find("#cwletbl"); q.Length() == 1 {
			ans.Question = strings.TrimSpace(q.Text())
		}
	}
}

func (p *GoogleProvider) parseMedia(ctx *goquery.Selection) *Media {
	media := &Media{}
	if u, ok := ctx.Find("img").Attr("src"); ok {
		media.Url = u
		return media
	}
	return nil
}
