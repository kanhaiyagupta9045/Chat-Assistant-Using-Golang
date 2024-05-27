package chatassistant

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

func ChatAssistant(filename string, filecontent []byte, assistantName, assistantInstructions, content string) (string, error) {
	if len(bytes.TrimSpace(filecontent)) == 0 {
		return "", fmt.Errorf("file content is empty or contains only whitespace")
	}
	client := openai.NewClient(os.Getenv("API_KEY"))
	openAIFile, err := client.CreateFileBytes(context.Background(), openai.FileBytesRequest{
		Name:    filename,
		Bytes:   filecontent,
		Purpose: openai.PurposeAssistants,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create file bytes: %w", err)
	}

	assistant, err := client.CreateAssistant(context.Background(), openai.AssistantRequest{
		Name:         &(assistantName),
		Model:        openai.GPT3Dot5Turbo1106,
		Instructions: &assistantInstructions,
		Tools:        []openai.AssistantTool{{Type: openai.AssistantToolTypeRetrieval}},
		FileIDs:      []string{openAIFile.ID},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create assistant: %w", err)
	}

	thread, err := client.CreateThread(context.Background(), openai.ThreadRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to create thread: %w", err)
	}

	_, err = client.CreateMessage(context.Background(), thread.ID, openai.MessageRequest{
		Role:    string(openai.ThreadMessageRoleUser),
		Content: content,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create message: %w", err)
	}

	run, err := client.CreateRun(context.Background(), thread.ID, openai.RunRequest{
		AssistantID: assistant.ID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create run: %w", err)
	}

	for run.Status != openai.RunStatusCompleted {
		time.Sleep(5 * time.Second)
		run, err = client.RetrieveRun(context.Background(), thread.ID, run.ID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve run: %w", err)
		}
	}

	msgs, err := client.ListMessage(context.Background(), thread.ID, nil, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to list messages: %w", err)
	}

	if len(msgs.Messages[0].Content) == 0 {
		return "", fmt.Errorf("no content found in the first message")
	}

	return msgs.Messages[0].Content[0].Text.Value, nil
}
