package know

import "testing"

func TestWolfram(t *testing.T) {
	ans, err := Wolfram.Ask("imaginaerum album release date")
	if err != nil {
		t.Fatal(err)
	}
	if ans == nil {
		t.Fatal("no answer found")
	}
	if ans.Question != "Imaginaerum (album) | release date" {
		t.Error("wrong question:", ans.Question)
	}
	if ans.Answer != "November 30, 2011" {
		t.Error("wrong answer:", ans.Answer)
	}
	t.Log(ans)
}
