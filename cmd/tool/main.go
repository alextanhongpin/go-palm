package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alextanhongpin/go-palm/internal/llms"
	"github.com/alextanhongpin/go-palm/internal/tools"
)

func main() {
	llm := llms.NewPalm(getPalmKey())
	defer llm.Close()

	question := `I have 77 houses, each with 31 cats.
Each cat owns 14 mittens, and 6 hats.
Each mitten was knit from 141m of yarn, each hat from 55m.
How much yarn was needed to make all the items?`

	prompt := `You are an expert at solving word problems. Here's a question:

{{.Question}}

-------------------

{{.Tools}}

-------------------

Work throught it step by step, and show your work.
One step per line.

Your solution:
`

	finalPrompt := tools.Template(prompt).Format(map[string]any{
		"Question": question,
		"Tools":    "{{.Tools}}",
	})

	ctx := context.Background()
	req := llms.DefaultGenerateTextRequest()
	req.Prompt = finalPrompt
	fmt.Println(llm.GenerateText(ctx, req, tools.NewMath()))
}

func getPalmKey() string {
	keyPath := filepath.Join(os.Getenv("HOME"), ".palm")
	b, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	b = bytes.TrimSpace(b)

	return string(b)
}
