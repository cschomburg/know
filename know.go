// Package know queries different knowledge providers and parses their result.
// Partially supported are currently Google and Wolfram Alpha.
package know

import "sync"

// Answer represents the collective result of a provider query.
type Answer struct {
	Query    string
	Question string
	Answer   string
	Media    []*Media
	Provider string
}

// Media represents a media object in the answer.
type Media struct {
	Type string
	Url  string
}

// NewAnswer constructs a new empty answer.
func NewAnswer() *Answer {
	return &Answer{
		Media: make([]*Media, 0),
	}
}

// Provider is a knowledge provider that can answer questions.
type Provider interface {
	Name() string
	Ask(question string) (*Answer, error)
}

// AskProviders asks all specified providers the same question and returns
// their answers and errors on a channel.
// The channels are closed when all providers are finished. Empty answer and
// error channels signify a successful query, but no provider could answer
// that question.
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

// Ask asks all currently supported providers the same question and returns
// their answers on a channel.
// The channels are closed when all providers are finished. Empty answer and
// error channels signify a successful query, but no provider could answer
// that question.
func Ask(question string) (<-chan *Answer, <-chan error) {
	return AskProviders(question, defaultProviders)
}

// RegisterProvider registers a new provider which is used in the default set
// of providers.
func RegisterProvider(p Provider) {
	defaultProviders = append(defaultProviders, p)
}
