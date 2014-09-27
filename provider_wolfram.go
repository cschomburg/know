package know

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WolframProvider struct {
	ApiKey string
}

func NewWolframProvider() *WolframProvider {
	return &WolframProvider{}
}

func (p *WolframProvider) SetApiKey(key string) {
	p.ApiKey = key
}

func (p *WolframProvider) Name() string {
	if p.ApiKey != "" {
		return "Wolfram API"
	}
	return "Wolfram"
}

type apiQueryResult struct {
	Success  bool     `xml:"success,attr"`
	HasError bool     `xml:"error,attr"`
	Error    string   `xml:"error>msg"`
	Pods     []apiPod `xml:"pod"`
}

type apiPod struct {
	Title   string      `xml:"title,attr"`
	SubPods []apiSubPod `xml:"subpod"`
}

type apiSubPod struct {
	Title     string `xml:"attr"`
	Plaintext string `xml:"plaintext"`
}

func (p *WolframProvider) Ask(question string) (*Answer, error) {
	if p.ApiKey == "" {
		return p.askPublic(question)
	}

	u, err := url.Parse("https://api.wolframalpha.com/v2/query")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("input", question)
	v.Set("appid", p.ApiKey)
	u.RawQuery = v.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Unexpected status: " + resp.Status)
	}

	defer resp.Body.Close()
	dec := xml.NewDecoder(resp.Body)
	result := apiQueryResult{}
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	if result.HasError {
		return nil, errors.New(result.Error)
	}

	ans := &Answer{}
	for _, pod := range result.Pods {
		for _, s := range pod.SubPods {
			switch {
			case strings.Contains(pod.Title, "Input"):
				ans.Question = s.Plaintext
			case strings.Contains(pod.Title, "Result"):
				ans.Answer = s.Plaintext
			case ans.Answer == "":
				ans.Answer = s.Plaintext
			}
		}
	}

	return ans, nil
}

func (p *WolframProvider) askPublic(question string) (*Answer, error) {
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
