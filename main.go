package main

import (
	"context"
	"fmt"
	"os"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := gl.NewTextRESTClient(ctx, option.WithAPIKey(os.Getenv("PALM_KEY")))
	if err != nil {
		panic(err)
	}

	defer client.Close()
	req := &pb.GenerateTextRequest{
		Model: "models/text-bison-001",
		Prompt: &pb.TextPrompt{
			Text: "What is the world's largest island that's not a continent?",
		},
	}

	resp, err := client.GenerateText(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Candidates[0].Output)
}
