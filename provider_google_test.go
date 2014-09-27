package know

import "testing"

func TestGoogleAlbum(t *testing.T) {
	google := GoogleProvider{t.Logf}
	ans, err := google.Ask("imaginaerum album release date")
	if err != nil {
		t.Fatal(err)
	}
	if ans == nil {
		t.Fatal("no answer found")
	}
	if ans.Question != "Imaginaerum, Release date" {
		t.Error("wrong question:", ans.Question)
	}
	if ans.Answer != "November 30, 2011" {
		t.Error("wrong answer:", ans.Answer)
	}
	t.Log(ans)
}

func TestGoogleFacts(t *testing.T) {
	google := GoogleProvider{t.Logf}
	ans, err := google.Ask("tuomas holopainen birth place")
	if err != nil {
		t.Fatal(err)
	}
	if ans == nil {
		t.Fatal("no answer found")
	}
	if ans.Question != "Tuomas Holopainen, Place of birth" {
		t.Error("wrong question:", ans.Question)
	}
	if ans.Answer != "Kitee, Finland" {
		t.Error("wrong answer:", ans.Answer)
	}
	t.Log(ans)
}

func TestGoogleCalc(t *testing.T) {
	google := GoogleProvider{t.Logf}
	ans, err := google.Ask("sin(3/2*pi)")
	if err != nil {
		t.Fatal(err)
	}
	if ans == nil {
		t.Fatal("no answer found")
	}
	if ans.Question != "sin((3 / 2) * pi radians) =" {
		t.Error("wrong question:", ans.Question)
	}
	if ans.Answer != "-1" {
		t.Error("wrong answer:", ans.Answer)
	}
	t.Log(ans)
}
