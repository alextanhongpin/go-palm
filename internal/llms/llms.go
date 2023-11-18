package llms

import (
	"context"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
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

func (l *PalmLLM) GenerateText(ctx context.Context, req GenerateTextRequest) (string, error) {
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

	resp, err := l.client.GenerateText(ctx, pbreq)
	if err != nil {
		return "", err
	}

	return resp.Candidates[0].Output, nil
}
