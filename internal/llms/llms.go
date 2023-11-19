package llms

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/alextanhongpin/go-palm/internal/tools"
	"google.golang.org/api/option"
)

type PalmLLM struct {
	client *gl.TextClient
}

func NewPalm(palmKey string) *PalmLLM {
	ctx := context.Background()
	client, err := gl.NewTextRESTClient(ctx, option.WithAPIKey(palmKey))
	if err != nil {
		panic(err)
	}

	return &PalmLLM{
		client: client,
	}
}

type GenerateTextRequest struct {
	Prompt          string
	Temperature     float64
	CandidateCount  int64
	MaxOutputTokens int64
	StopSequences   []string
}

func DefaultGenerateTextRequest() GenerateTextRequest {
	return GenerateTextRequest{
		Temperature:     0.0,
		CandidateCount:  1,
		MaxOutputTokens: 4096,
	}
}

func (l *PalmLLM) Close() error {
	return l.client.Close()
}

type tool interface {
	Eval(prompt string) (string, error)
	Name() string
	Description() string
	Tag() string
}

// GenerateText generates text from the prompt.
func (l *PalmLLM) GenerateText(ctx context.Context, req GenerateTextRequest, tools ...tool) (string, error) {
	return l.generateText(ctx, req, tools...)
}

func (l *PalmLLM) buildRequest(ctx context.Context, req GenerateTextRequest) (*pb.GenerateTextRequest, error) {
	temperature := float32(req.Temperature)
	count := int32(req.CandidateCount)
	maxOutputTokens := int32(req.MaxOutputTokens)

	pbreq := &pb.GenerateTextRequest{
		Model: "models/text-bison-001",
		Prompt: &pb.TextPrompt{
			Text: req.Prompt,
		},
		Temperature:     &temperature,
		CandidateCount:  &count,
		MaxOutputTokens: &maxOutputTokens,
	}

	return pbreq, nil
}

func (l *PalmLLM) generateText(ctx context.Context, req GenerateTextRequest, ts ...tool) (string, error) {
	pbreq, err := l.buildRequest(ctx, req)
	if err != nil {
		return "", err
	}

	if hasTools(pbreq.GetPrompt().GetText()) && len(ts) == 0 {
		return "", errors.New("no tools specified")
	}

	divider := strings.Repeat("=", 40)

	toolsPrompt := make([]string, len(ts))

	for i, tool := range ts {
		toolsPrompt[i] = tool.Description()
		endTag := fmt.Sprintf("</%s>", tool.Tag())
		pbreq.StopSequences = append(pbreq.StopSequences, endTag)
	}
	toolsPrompt = append(toolsPrompt, divider)
	toolsPrompt = append([]string{divider}, toolsPrompt...)
	toolPrompt := strings.Join(toolsPrompt, divider)

	prompt := tools.Template(pbreq.GetPrompt().GetText()).Format(map[string]string{
		"Tools": toolPrompt,
	})

	maxLoop := 10
	var result []string
	for {
		pbreq.Prompt.Text = strings.Join(append([]string{prompt}, result...), " ")
		resp, err := l.client.GenerateText(ctx, pbreq)
		if err != nil {
			return "", err
		}
		if len(resp.Candidates) == 0 {
			log.Println("no more output")
			return strings.Join(result, " "), nil
		}
		output := resp.Candidates[0].Output

		maxLoop--
		if maxLoop < 0 {
			return strings.Join(result, " "), nil
		}

		var hasChanged bool
		for _, tool := range ts {
			extendedPrompt, err := tool.Eval(output)
			if err != nil {
				return "", err
			}

			if output == extendedPrompt {
				continue
			}
			hasChanged = true

			result = append(result, extendedPrompt)
			break
		}
		if !hasChanged {
			result = append(result, output)
		}
	}
}

func hasTools(text string) bool {
	pattern := `\{\{\s*.Tools\s*\}\}`
	exists, err := regexp.MatchString(pattern, text)
	if err != nil {
		panic(err)
	}

	return exists
}
