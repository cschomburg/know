package know

import (
	"encoding/base64"
	"errors"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WolframProvider struct{}

func NewWolframProvider() *WolframProvider {
	return &WolframProvider{}
}

func (p *WolframProvider) Name() string {
	return "Wolfram"
}

func (p *WolframProvider) Ask(question string) (*Answer, error) {
	u, err := url.Parse("https://www.wolframalpha.com/input/")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("i", question)
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

func (p *WolframProvider) parse(body io.Reader) (*Answer, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	ans := &Answer{}
	pods := doc.Find("#answers section.pod")
	for i := 0; i < pods.Length(); i++ {
		pod := pods.Eq(i)
		title := strings.TrimSpace(pod.Find("header h2").Text())
		plain, ok := pod.Find(".sub[data-s]").Attr("data-s")
		if !ok {
			continue
		}
		b, err := base64.StdEncoding.DecodeString(plain)
		if err != nil {
			return nil, err
		}
		value := html.UnescapeString(string(b))

		switch {
		case strings.Contains(title, "Input"):
			ans.Question = value
		case strings.Contains(title, "Result"):
			ans.Answer = value
		case ans.Answer == "":
			ans.Answer = value
		}
	}
	return ans, nil
}
