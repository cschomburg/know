stark
=====

[![API Documentation](http://img.shields.io/badge/api-GoDoc-blue.svg?style=flat-square)](http://godoc.org/github.com/xconstruct/know)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://opensource.org/licenses/MIT)

Package know queries different knowledge providers and parses their result.
Partially supported are currently Google and Wolfram Alpha.

## Installation

`go get github.com/xconstruct/know`

## Example

```go
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
```
