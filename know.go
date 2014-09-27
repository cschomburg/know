package know

import "sync"

type Answer struct {
	Query    string
	Question string
	Answer   string
	Media    []*Media
	Provider string
}

type Media struct {
	Type string
	Url  string
}

func NewAnswer() *Answer {
	return &Answer{
		Media: make([]*Media, 0),
	}
}

type Provider interface {
	Name() string
	Ask(question string) (*Answer, error)
}

func AskProviders(question string, providers []Provider) (<-chan *Answer, <-chan error) {
	n := len(providers)
	answers := make(chan *Answer, n)
	errors := make(chan error, n)

	wait := sync.WaitGroup{}
	wait.Add(n)
	for _, p := range providers {
		go func(p Provider) {
			ans, err := p.Ask(question)
			if err != nil {
				errors <- err
			} else if ans != nil && ans.Answer != "" {
				ans.Query = question
				ans.Provider = p.Name()
				answers <- ans
			}
			wait.Done()
		}(p)
	}
	go func() {
		wait.Wait()
		close(answers)
		close(errors)
	}()

	return answers, errors
}

var (
	Google  Provider = NewGoogleProvider()
	Wolfram Provider = NewWolframProvider()

	defaultProviders = []Provider{
		Google,
		Wolfram,
	}
)

func Ask(question string) (<-chan *Answer, <-chan error) {
	return AskProviders(question, defaultProviders)
}

func RegisterProvider(p Provider) {
	defaultProviders = append(defaultProviders, p)
}
