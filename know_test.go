package know

import (
	"testing"
	"time"
)

type TestProvider struct {
	Delay  time.Duration
	MyName string
	T      *testing.T
}

func (p *TestProvider) Name() string {
	return p.MyName
}

func (p *TestProvider) Ask(question string) (*Answer, error) {
	time.Sleep(p.Delay)

	ans := &Answer{Answer: p.MyName}
	if question == "all" || question == p.MyName {
		return ans, nil
	}
	return nil, nil
}

type TestQuestion struct {
	Q       string
	Answers []string
}

func TestMultiple(t *testing.T) {
	provs := []Provider{
		&TestProvider{30 * time.Millisecond, "alice", t},
		&TestProvider{10 * time.Millisecond, "bob", t},
	}

	tests := []TestQuestion{
		{"none", []string{}},
		{"alice", []string{"alice"}},
		{"bob", []string{"bob"}},
		{"all", []string{"bob", "alice"}},
	}

	for _, test := range tests {
		answers, errors := AskProviders(test.Q, provs)
		for err := range errors {
			t.Error(err)
		}

		i := 0
		for ans := range answers {
			if i < len(test.Answers) {
				if ans.Answer != test.Answers[i] {
					t.Errorf("%s: expected answer '%s', got '%s'", test.Q,
						test.Answers[i], ans.Answer)
				} else {
					t.Logf("%s: correct answer '%s'", test.Q, ans.Answer)
				}
			} else {
				t.Errorf("%s: expected no answer, got '%s'", test.Q, ans.Answer)
			}
			i++
		}
		if i < len(test.Answers) {
			t.Logf("%s: %d missing answers", test.Q, len(test.Answers)-i)
		}
	}
}
