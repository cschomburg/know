package know_test

import (
	"fmt"

	"github.com/xconstruct/know"
)

func ExampleAll() {
	answers, errs := know.Ask("What is the capital of germany?")

	// Get first answer
	ans, ok := <-answers
	if !ok {
		fmt.Println("No answer found!")
		for err := range errs {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(ans.Question, "is", ans.Answer)
	// Output:
	// Germany, Capital is Berlin
}
