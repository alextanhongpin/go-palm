package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/alextanhongpin/go-palm/internal/llms"
)

var prompt = `
	You are an experienced technical writer that is able to explain complicated terms in simple words.
	Given the following documentation/tech article:

	1. make it more readable
	2. improve the points

	Document:

	"""
	%s
	"""
`

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	text := string(b)

	llm := llms.NewPalm(os.Getenv("PALM_KEY"))
	defer llm.Close()

	ctx := context.Background()

	req := llms.DefaultGenerateTextRequest()
	req.Prompt = fmt.Sprintf(prompt, text)

	resp, err := llm.GenerateText(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
