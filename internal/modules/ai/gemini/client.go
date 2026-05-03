package gemini

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"ticktick-ai/internal/domain"
)

const endpointTemplate = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s"

type Client struct {
	httpClient *http.Client
	apiKey     string
	model      string
}

func NewClient(httpClient *http.Client, apiKey string, model string) *Client {
	return &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		model:      model,
	}
}

func (c *Client) ParseText(ctx context.Context, text string, timezone string) (domain.ParsedIntent, error) {
	return c.generate(ctx, []part{
		{Text: buildUserPrompt(text, timezone)},
	})
}

func (c *Client) ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (domain.ParsedIntent, error) {
	if mimeType == "" {
		mimeType = "audio/ogg"
	}

	return c.generate(ctx, []part{
		{Text: buildUserPrompt("Parse this voice message.", timezone)},
		{InlineData: &inlineData{
			MimeType: mimeType,
			Data:     base64.StdEncoding.EncodeToString(audio),
		}},
	})
}

func (c *Client) generate(ctx context.Context, parts []part) (domain.ParsedIntent, error) {
	reqBody := generateRequest{
		SystemInstruction: content{Parts: []part{{Text: systemPrompt}}},
		Contents: []content{
			{
				Role:  "user",
				Parts: parts,
			},
		},
		GenerationConfig: generationConfig{
			ResponseMimeType: "application/json",
			Temperature:      0.1,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return domain.ParsedIntent{}, err
	}

	endpoint := fmt.Sprintf(endpointTemplate, url.PathEscape(c.model), url.QueryEscape(c.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return domain.ParsedIntent{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.ParsedIntent{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ParsedIntent{}, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return domain.ParsedIntent{}, fmt.Errorf("gemini api status %d: %s", resp.StatusCode, string(respBody))
	}

	var parsed generateResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return domain.ParsedIntent{}, err
	}

	text := parsed.firstText()
	if text == "" {
		return domain.ParsedIntent{}, errors.New("gemini returned empty response")
	}

	return decodeIntent(text)
}

func buildUserPrompt(text string, timezone string) string {
	return fmt.Sprintf("Timezone: %s\nUser message: %s", timezone, text)
}

func decodeIntent(text string) (domain.ParsedIntent, error) {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var intent domain.ParsedIntent
	if err := json.Unmarshal([]byte(text), &intent); err != nil {
		return domain.ParsedIntent{}, fmt.Errorf("decode gemini intent: %w: %s", err, text)
	}

	if intent.Priority == "" {
		intent.Priority = domain.PriorityNone
	}
	return intent, nil
}

type generateRequest struct {
	SystemInstruction content          `json:"systemInstruction"`
	Contents          []content        `json:"contents"`
	GenerationConfig  generationConfig `json:"generationConfig"`
}

type generationConfig struct {
	ResponseMimeType string  `json:"responseMimeType"`
	Temperature      float64 `json:"temperature"`
}

type content struct {
	Role  string `json:"role,omitempty"`
	Parts []part `json:"parts"`
}

type part struct {
	Text       string      `json:"text,omitempty"`
	InlineData *inlineData `json:"inlineData,omitempty"`
}

type inlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type generateResponse struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

func (r generateResponse) firstText() string {
	for _, candidate := range r.Candidates {
		for _, part := range candidate.Content.Parts {
			if strings.TrimSpace(part.Text) != "" {
				return part.Text
			}
		}
	}
	return ""
}
