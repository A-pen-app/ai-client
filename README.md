# AI Client Library

A Go library that provides unified AI client interface for OCR and article processing services, supporting multiple AI providers (OpenAI and Google Gemini).

## Features

### OCR Service
- 🔍 Extract names from ID cards, licenses, and certificates
- 📋 Extract detailed information based on profession type
- 🏥 Support for multiple medical professions (Doctor, Nurse, Pharmacist)
- 📨 Message queue integration for data pipeline

### Article Service
- 🏷️ Extract structured tags from job postings
- ✨ Polish and format article content with AI
- 👔 Profession-specific prompts (Doctor, Nurse, Pharmacist)
- 📝 Smart formatting and tone enhancement

### AI Client
- 🔌 Unified interface for multiple AI providers
- 🤖 OpenAI GPT-4o support
- 🌟 Google Gemini API support
- 🖼️ Vision API support for image analysis
- 🛡️ Comprehensive error handling
- 🔧 Flexible configuration options

## Installation

```bash
go get github.com/A-pen-app/ai-client
```

## Requirements

- Go 1.24.0 or higher
- OpenAI API key (for OpenAI client)
- Google Cloud Project with Gemini API access (for Gemini client)
- Message queue (implements `mq.MQ` interface) - for OCR service only

## Usage

### 1. Initialize AI Client

#### Option A: OpenAI Client

```go
import (
    "github.com/A-pen-app/ai-client/client/openai"
    openaiSDK "github.com/openai/openai-go/v2"
)

// Create OpenAI client
aiClient, err := openai.NewClient("your-openai-api-key", openaiSDK.ChatModelGPT4o)
if err != nil {
    log.Fatal(err)
}
```

#### Option B: Gemini Client

```go
import (
    "github.com/A-pen-app/ai-client/client/gemini"
)

// Create Gemini client
aiClient, err := gemini.NewClient("your-project-id", "us-central1", "gemini-2.5-flash")
if err != nil {
    log.Fatal(err)
}
```

### 2. Article Service

#### Extract Tags from Job Posting

```go
import (
    "context"
    "github.com/A-pen-app/ai-client/models"
    "github.com/A-pen-app/ai-client/store"
)

// Create article store
articleStore := store.NewArticleStore(aiClient, &store.ArticleConfig{
    MaxToken: 2048,
})

ctx := context.Background()
content := "誠徵全職主治醫師，地點：台北市南港區..."

// Extract tags (for doctors)
tags, err := articleStore.ExtractTags(ctx, content, models.PlatformTypeApen)
if err != nil {
    log.Fatal(err)
}

fmt.Println(tags) // Returns JSON with 工作類別, 需求科別, 需求職級, 職缺地點
```

#### Polish Article Content

```go
// Polish article content
polished, err := articleStore.Polish(ctx, content, models.PlatformTypeApen)
if err != nil {
    log.Fatal(err)
}

fmt.Println(polished) // Returns formatted and polished content
```

### 3. OCR Service

#### Initialize OCR Store

```go
import (
    "github.com/A-pen-app/ai-client/store"
)

// Create OCR store (requires OpenAI client and message queue)
config := &store.OpenAIConfig{
    MaxToken:    2048,
    Model:       openaiSDK.ChatModelGPT4o,
    Topic:       models.OCRTopicProd,
    MessageType: models.OCRMessageTypeIdentifyOCR,
}

ocrStore := store.NewOpenAIStore(mq, openaiClient, config)
```

#### Scan Name Only

```go
ctx := context.Background()
imageURL := "https://example.com/id-card.jpg"

name, err := ocrStore.ScanName(ctx, imageURL)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", name)
```

#### Scan Complete Information

```go
ctx := context.Background()
userID := "user-123"
imageURL := "https://example.com/doctor-license.jpg"

// For doctor
ocrInfo, err := ocrStore.ScanRawInfo(ctx, userID, imageURL, models.PlatformTypeApen)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", *ocrInfo.Name)
fmt.Printf("Birthday: %s\n", *ocrInfo.Birthday)
fmt.Printf("Position: %s\n", *ocrInfo.Position)
fmt.Printf("Department: %s\n", *ocrInfo.Department)
fmt.Printf("Facility: %s\n", *ocrInfo.Facility)
```

## Supported Platform Types

| Platform Type | Description | Extracted Fields |
|--------------|-------------|------------------|
| `PlatformTypeApen` | Doctor | Name, Birthday, Position, Department, Facility, Valid Date, Specialty Valid Date |
| `PlatformTypeNurse` | Nurse | Name, Birthday, Department, Facility, Valid Date |
| `PlatformTypePhar` | Pharmacist | Name, Birthday, Facility, Valid Date |

## API Reference

### AI Client Interface

```go
type AIClient interface {
    Generate(ctx context.Context, message models.AIChatMessage, opts models.AIClientOptions) (string, error)
}
```

#### OpenAI Client

```go
func NewClient(apiKey string, model openai.ChatModel) (store.AIClient, error)
```

**Parameters:**
- `apiKey`: OpenAI API key
- `model`: Default model to use (e.g., `openai.ChatModelGPT4o`)

#### Gemini Client

```go
func NewClient(projectID string, location string, model string) (store.AIClient, error)
```

**Parameters:**
- `projectID`: GCP project ID
- `location`: GCP region (e.g., "us-central1")
- `model`: Model name (e.g., "gemini-2.5-flash", default: "gemini-2.5-flash")

### Article Service

#### `NewArticleStore`

```go
func NewArticleStore(aiClient AIClient, config *ArticleConfig) Article
```

**Parameters:**
- `aiClient`: AI client instance (OpenAI or Gemini)
- `config`: Configuration (optional, defaults: MaxToken=2048)

#### `ExtractTags`

Extracts structured tags from job posting content.

```go
func (s *articleStore) ExtractTags(
    ctx context.Context,
    content string,
    professionType models.PlatformType,
) (string, error)
```

**Parameters:**
- `ctx`: Context for request cancellation
- `content`: Job posting content
- `professionType`: Type of profession (PlatformTypeApen/PlatformTypeNurse/PlatformTypePhar)

**Returns:** JSON string with extracted tags

**Example Output (for Doctor):**
```json
{
  "工作類別": ["全職"],
  "需求科別": ["家庭醫學科", "內科"],
  "需求職級": ["主治醫師"],
  "職缺地點": ["臺北市", "南港區"]
}
```

#### `Polish`

Polishes and formats article content with AI.

```go
func (s *articleStore) Polish(
    ctx context.Context,
    content string,
    professionType models.PlatformType,
) (string, error)
```

**Parameters:**
- `ctx`: Context for request cancellation
- `content`: Original article content
- `professionType`: Type of profession (determines prompt style)

**Returns:** Polished and formatted content

### OCR Service

#### `NewOpenAIStore`

Creates a new OCR store instance (requires OpenAI client).

```go
func NewOpenAIStore(
    mq mq.MQ,
    client *openai.Client,
    config *OpenAIConfig,
) OCR
```

**Parameters:**
- `mq`: Message queue for publishing OCR results
- `client`: OpenAI client instance
- `config`: Configuration (optional, uses default if nil)

#### `ScanName`

Extracts only the name from an image.

```go
func (os *ocrStore) ScanName(
    ctx context.Context,
    link string,
) (string, error)
```

**Parameters:**
- `ctx`: Context for request cancellation
- `link`: URL of the image to scan

**Returns:** Extracted name and error (if any)

#### `ScanRawInfo`

Extracts comprehensive information based on profession type.

```go
func (os *ocrStore) ScanRawInfo(
    ctx context.Context,
    userID string,
    link string,
    platformType models.PlatformType,
) (*models.OCRRawInfo, error)
```

**Parameters:**
- `ctx`: Context for request cancellation
- `userID`: User identifier for tracking
- `link`: URL of the image to scan
- `platformType`: Type of profession (PlatformTypeApen/PlatformTypeNurse/PlatformTypePhar)

**Returns:** Extracted OCR information and error (if any)

## Data Models

### `AIChatMessage`

```go
type AIChatMessage struct {
    SystemPrompt string
    Text         string
    ImageUrls    []string
}
```

### `AIClientOptions`

```go
type AIClientOptions struct {
    MaxTokens      int64
    Model          string
    ResponseFormat ResponseFormat  // "json" or "text"
}
```

### `PlatformType`

Profession types for OCR and article processing:

```go
const (
    PlatformTypeApen  PlatformType = "apen"   // Doctor
    PlatformTypeNurse PlatformType = "nurse"  // Nurse
    PlatformTypePhar  PlatformType = "phar"   // Pharmacist
)
```

### `OCRRawInfo`

```go
type OCRRawInfo struct {
    IdentifyURL        *string `json:"identify_url,omitempty"`
    Name               *string `json:"name"`
    Birthday           *string `json:"birthday"`
    Position           *string `json:"position,omitempty"`          // Doctor only
    Department         *string `json:"department,omitempty"`
    Facility           *string `json:"facility,omitempty"`
    ValidDate          *string `json:"valid_date,omitempty"`
    SpecialtyValidDate *string `json:"specialty_valid_date,omitempty"` // Doctor only
}
```

### `ArticleConfig`

```go
type ArticleConfig struct {
    MaxToken int64  // Maximum tokens for response (default: 2048)
}
```

### `OpenAIConfig`

```go
type OpenAIConfig struct {
    MaxToken    int64                  // Maximum tokens for response
    Model       openai.ChatModel       // OpenAI model to use
    Topic       models.OCRTopic        // Message queue topic
    MessageType models.OCRMessageType  // Message type identifier
}
```

**Default Values:**
- `MaxToken`: 1024
- `Model`: GPT-4o
- `Topic`: wanderer-dev
- `MessageType`: identify_ocr

## Position Types (for Doctors)

- `PGY` - Post-Graduate Year (不分科醫師)
- `Resident` - Resident Doctor (住院醫師)
- `VS` - Visiting Staff / Attending Physician (主治醫師)

## Error Handling

The library provides detailed error messages for common issues:

- `"openai client is not initialized"` - Client not properly configured
- `"empty response choices from OCR"` - No response from OpenAI API
- `"empty response content from OCR"` - Empty content in API response
- JSON unmarshal errors for invalid response format

## Message Queue Integration

When `ScanRawInfo` is called, the result is automatically published to the configured message queue topic as an `OCREventMessage`:

```go
type OCREventMessage struct {
    UserID    string     `json:"user_id"`
    Payload   OCRRawInfo `json:"payload"`
    CreatedAt time.Time  `json:"created_at"`
    Type      string     `json:"type"`
    Source    string     `json:"source"`
}
```

## Examples

### Example 1: Article Processing with OpenAI

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/A-pen-app/ai-client/client/openai"
    "github.com/A-pen-app/ai-client/models"
    "github.com/A-pen-app/ai-client/store"
    openaiSDK "github.com/openai/openai-go/v2"
)

func main() {
    ctx := context.Background()

    // Initialize OpenAI client
    aiClient, err := openai.NewClient("your-api-key", openaiSDK.ChatModelGPT4o)
    if err != nil {
        log.Fatal(err)
    }

    // Create article store
    articleStore := store.NewArticleStore(aiClient, nil)

    // Job posting content
    content := `
    誠徵全職家庭醫學科主治醫師
    工作地點：台北市南港區
    需具備家醫科專科證書
    待遇優渥，歡迎聯繫
    `

    // Extract tags
    tags, err := articleStore.ExtractTags(ctx, content, models.PlatformTypeApen)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Extracted Tags:", tags)

    // Polish content
    polished, err := articleStore.Polish(ctx, content, models.PlatformTypeApen)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Polished Content:", polished)
}
```

### Example 2: Article Processing with Gemini

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/A-pen-app/ai-client/client/gemini"
    "github.com/A-pen-app/ai-client/models"
    "github.com/A-pen-app/ai-client/store"
)

func main() {
    ctx := context.Background()

    // Initialize Gemini client
    aiClient, err := gemini.NewClient("your-project-id", "us-central1", "gemini-2.5-flash")
    if err != nil {
        log.Fatal(err)
    }

    // Create article store
    articleStore := store.NewArticleStore(aiClient, &store.ArticleConfig{
        MaxToken: 4096,
    })

    // Process nurse job posting
    content := "徵護理師，全職或兼職，台中市工作"

    tags, err := articleStore.ExtractTags(ctx, content, models.PlatformTypeNurse)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tags:", tags)
}
```

### Example 3: OCR with OpenAI

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/A-pen-app/ai-client/models"
    "github.com/A-pen-app/ai-client/store"
    "github.com/openai/openai-go/v2"
)

func main() {
    ctx := context.Background()

    // Initialize OpenAI client (original SDK)
    client := openai.NewClient(
        openai.WithAPIKey("your-api-key"),
    )

    // Initialize your message queue
    // mq := ... (your MQ implementation)

    // Create OCR store
    ocrStore := store.NewOpenAIStore(mq, &client, nil)

    // Scan doctor's license
    imageURL := "https://example.com/doctor-license.jpg"
    result, err := ocrStore.ScanRawInfo(
        ctx,
        "user-123",
        imageURL,
        models.PlatformTypeApen,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    if result.Name != nil {
        fmt.Printf("Name: %s\n", *result.Name)
    }
    if result.Position != nil {
        fmt.Printf("Position: %s\n", *result.Position)
    }
    if result.Department != nil {
        fmt.Printf("Department: %s\n", *result.Department)
    }
}
```

### Example 4: Direct AI Client Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/A-pen-app/ai-client/client/openai"
    "github.com/A-pen-app/ai-client/models"
    openaiSDK "github.com/openai/openai-go/v2"
)

func main() {
    ctx := context.Background()

    // Create AI client
    aiClient, err := openai.NewClient("your-api-key", openaiSDK.ChatModelGPT4o)
    if err != nil {
        log.Fatal(err)
    }

    // Prepare message
    message := models.AIChatMessage{
        SystemPrompt: "You are a helpful assistant.",
        Text:         "Explain what is Go programming language.",
        ImageUrls:    []string{},
    }

    // Generate response
    response, err := aiClient.Generate(ctx, message, models.AIClientOptions{
        MaxTokens:      1024,
        ResponseFormat: models.ResponseFormatText,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response)
}
```

## Architecture

```
ai-client/
├── client/              # AI provider implementations
│   ├── openai/         # OpenAI GPT-4o client
│   └── gemini/         # Google Gemini API client
├── models/             # Data models and prompts
│   ├── ai_client.go    # Common AI client models
│   ├── article.go      # Article processing prompts
│   └── ocr.go          # OCR models and constants
├── store/              # Service layer
│   ├── store.go        # Interface definitions
│   ├── article.go      # Article processing service
│   └── ocr.go          # OCR service
└── util/               # Utility functions
```

## Error Handling

The library provides detailed error messages for common issues:

### AI Client Errors
- `"openai API key cannot be empty"` - Missing API key
- `"openai client is not initialized"` - Client not properly configured
- `"gemini client is not initialized"` - Gemini client not initialized
- `"failed to create Gemini client"` - GCP authentication or configuration issue

### Service Errors
- `"AI client is not initialized"` - AI client not provided to service
- `"empty response content from AI client"` - Empty response from AI API
- `"empty response choices from OpenAI"` - No response choices in API response
- `"empty response from Gemini"` - No candidates in Gemini response

### OCR Specific Errors
- JSON unmarshal errors for invalid response format
- Image download errors (for Gemini with image URLs)

## Dependencies

### Core Dependencies
- [openai-go](https://github.com/openai/openai-go) - OpenAI API client (v2.7.1+)
- [google.golang.org/genai](https://pkg.go.dev/google.golang.org/genai) - Google Gemini API client (v1.36.0+)

### Internal Dependencies
- [A-pen-app/mq](https://github.com/A-pen-app/mq) - Message queue abstraction (v2.0.5+)
- [A-pen-app/logging](https://github.com/A-pen-app/logging) - Logging utilities (v0.4.0+)

### Other Dependencies
- [tidwall/sjson](https://github.com/tidwall/sjson) - JSON manipulation (v1.2.5+)

## License

Copyright © 2025 A-pen

## Contributing

This is a private library. For issues or questions, please contact the development team.

