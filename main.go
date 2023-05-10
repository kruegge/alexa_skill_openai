package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type AlexaRequest struct {
	Version string `json:"version"`
	Request struct {
		Type   string `json:"type"`
		Intent struct {
			Name  string `json:"name"`
			Slots struct {
				Question struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"question"`
			} `json:"slots"`
		} `json:"intent"`
		Directive struct {
			Payload struct {
				Message string `json:"message"`
			} `json:"payload"`
		} `json:"directive"`
	} `json:"request"`
}

type AlexaResponse struct {
	Version  string   `json:"version"`
	Response Response `json:"response"`
}

type Response struct {
	OutputSpeech     OutputSpeech `json:"outputSpeech"`
	Repromt          *Repromt     `json:"reprompt, omitempty"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}

type Repromt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func generateResponse(responseMessage string, repromtMessage string, withRepromt bool, endSession bool) AlexaResponse {
	alexaResponse := AlexaResponse{
		Version: "1.0",
		Response: Response{
			OutputSpeech: OutputSpeech{
				Type: "PlainText",
				Text: responseMessage,
			},
			ShouldEndSession: endSession,
		},
	}

	if withRepromt {
		alexaResponse.Response.Repromt = &Repromt{
			OutputSpeech: OutputSpeech{
				Type: "PlainText",
				Text: repromtMessage,
			},
		}
	}

	return sendResponse(alexaResponse)
}

func handleRequest(ctx context.Context, req interface{}) (AlexaResponse, error) {
	rawData, _ := json.Marshal(req)
	log.Println(string(rawData))

	var request AlexaRequest
	_ = json.Unmarshal(rawData, &request)

	switch request.Request.Type {
	case "LaunchRequest":
		return generateResponse("Du sprichst mit GPT3. Was möchtest du wissen?", "Was möchtest du wissen?", true, false), nil
	case "IntentRequest":
		if request.Request.Intent.Name == "AMAZON.StopIntent" {
			return generateResponse("Bis bald!", "", false, true), nil
		}

		if request.Request.Intent.Slots.Question.Value == "" {
			return generateResponse("Ich habe dich leider nicht richtig verstanden. Kannst du es bitte wiederholen?", "Was möchtest du wissen?", true, false), nil
		}

		return generateResponse(handleOpenAIRequest(request.Request.Intent.Slots.Question.Value), "Noch was?", true, false), nil
	case "SessionEndedRequest":
		return generateResponse("Bis bald!", "", false, true), nil
	}

	return generateResponse("Es ist ein Fehler aufgetreten. Bis bald!", "", false, true), nil
}

func handleOpenAIRequest(requestMessage string) string {
	// Create a new OpenAI API client with your API key
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Generate a short answer using the OpenAI GPT-3 API
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: requestMessage,
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	// Get the generated answer from the OpenAI API response
	answer := resp.Choices[0].Message.Content

	return answer
}

func sendResponse(response AlexaResponse) AlexaResponse {
	res, _ := json.Marshal(response)
	log.Println(string(res))

	return response
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lambda.Start(handleRequest)
}
