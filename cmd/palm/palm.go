package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alextanhongpin/go-palm/internal/llms"
)

// go install cmd/palm/palm.go
// palm -prompt="what is 1+1" | glow
func main() {
	var keyPath, prompt string
	var answerOnly bool
	var dryRun bool

	flag.StringVar(&keyPath, "key", filepath.Join(os.Getenv("HOME"), ".palm"), "the path to PaLM key")
	flag.StringVar(&prompt, "prompt", "", "the prompt")
	flag.StringVar(&prompt, "p", "", "the prompt (shorthand)")
	flag.BoolVar(&answerOnly, "answer", true, "only shows the answer, not the prompt")
	flag.BoolVar(&answerOnly, "a", true, "only shows the answer, not the prompt (shorthand)")
	flag.BoolVar(&dryRun, "dry-run", false, "execute dry run")
	flag.Parse()

	b, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	b = bytes.TrimSpace(b)

	if dryRun {
		fmt.Println("# PROMPT")
		fmt.Println(prompt)
		os.Exit(0)
	}

	llm := llms.NewPalm(string(b))
	defer llm.Close()

	ctx := context.Background()

	req := llms.DefaultGenerateTextRequest()
	req.Prompt = prompt

	resp, err := llm.GenerateText(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	if !answerOnly {
		fmt.Println("# PROMPT")
		fmt.Println(prompt)
		fmt.Println()
	}
	fmt.Println("# ANSWER")
	fmt.Println(resp)
}
