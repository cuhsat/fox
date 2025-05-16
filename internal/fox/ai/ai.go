//go:build ai

/*
Advanced Usage: RAG System with Ollama Embeddings
Here's a more advanced example showing how to build a simple RAG (Retrieval Augmented Generation) system using Ollama embeddings and LangChainGo:

package main

import (
    "context"
    "fmt"
    "log"
    "sort"

    "github.com/tmc/langchaingo/embeddings"
    "github.com/tmc/langchaingo/embeddings/ollama"
    "github.com/tmc/langchaingo/llms"
    ollamallm "github.com/tmc/langchaingo/llms/ollama"
)

// Document represents a text with its embedding
type Document struct {
    Content   string
    Embedding []float64
}

// DocumentWithScore represents a document with a similarity score
type DocumentWithScore struct {
    Content string
    Score   float64
}

// Calculate cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
    var dotProduct float64
    var normA float64
    var normB float64

    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }

    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func main() {
    ctx := context.Background()

    // Initialize the Ollama embeddings model with a specific model
    embeddingModel, err := ollama.New(
        ollama.WithModel("nomic-embed-text"), // Specify embedding model
    )
    if err != nil {
        log.Fatalf("Failed to initialize Ollama embeddings: %v", err)
    }

    // Initialize the Ollama LLM for text generation
    llm, err := ollamallm.New(
        ollamallm.WithModel("llama3"), // Specify LLM model
    )
    if err != nil {
        log.Fatalf("Failed to initialize Ollama LLM: %v", err)
    }

    // Sample knowledge base
    knowledgeBase := []string{
        "LangChain is a framework for developing applications powered by language models.",
        "Ollama is an open-source project that allows running LLMs locally.",
        "Embeddings are vector representations of text that capture semantic meaning.",
        "RAG stands for Retrieval Augmented Generation, a technique to enhance LLM responses with external knowledge.",
        "Go is a statically typed, compiled programming language designed at Google.",
    }

    // Create document embeddings
    fmt.Println("Creating document embeddings...")
    var documents []Document

    docEmbeddings, err := embeddingModel.EmbedDocuments(ctx, knowledgeBase)
    if err != nil {
        log.Fatalf("Failed to generate document embeddings: %v", err)
    }

    for i, content := range knowledgeBase {
        documents = append(documents, Document{
            Content:   content,
            Embedding: docEmbeddings[i],
        })
    }

    // Function to retrieve relevant documents
    retrieveRelevantDocs := func(query string, k int) ([]DocumentWithScore, error) {
        // Generate embedding for the query
        queryEmbedding, err := embeddingModel.EmbedQuery(ctx, query)
        if err != nil {
            return nil, fmt.Errorf("failed to generate query embedding: %v", err)
        }

        // Calculate similarity scores
        var docsWithScores []DocumentWithScore
        for _, doc := range documents {
            score := cosineSimilarity(queryEmbedding, doc.Embedding)
            docsWithScores = append(docsWithScores, DocumentWithScore{
                Content: doc.Content,
                Score:   score,
            })
        }

        // Sort by similarity score (descending)
        sort.Slice(docsWithScores, func(i, j int) bool {
            return docsWithScores[i].Score > docsWithScores[j].Score
        })

        // Return top k results
        if k > len(docsWithScores) {
            k = len(docsWithScores)
        }
        return docsWithScores[:k], nil
    }

    // Example query
    query := "What is LangChain and how does it relate to embeddings?"
    fmt.Printf("\nQuery: %s\n", query)

    // Retrieve relevant documents
    relevantDocs, err := retrieveRelevantDocs(query, 2)
    if err != nil {
        log.Fatalf("Failed to retrieve relevant documents: %v", err)
    }

    fmt.Println("\nRelevant documents:")
    var context string
    for i, doc := range relevantDocs {
        fmt.Printf("%d. [Score: %.4f] %s\n", i+1, doc.Score, doc.Content)
        context += doc.Content + "\n"
    }

    // Generate response using the LLM with retrieved context
    prompt := fmt.Sprintf("Based on the following information:\n\n%s\n\nAnswer this question: %s", context, query)

    response, err := llm.Call(ctx, prompt, llms.WithTemperature(0.2))
    if err != nil {
        log.Fatalf("Failed to generate response: %v", err)
    }

    fmt.Printf("\nGenerated response:\n%s\n", response)
}
*/

package ai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"

	"github.com/cuhsat/fox/internal/pkg/sys"
	"github.com/cuhsat/fox/internal/pkg/text"
	"github.com/cuhsat/fox/internal/pkg/types/heap"
	"github.com/cuhsat/fox/internal/pkg/user/history"
)

const (
	Build = true
)

const (
	Default = "mistral"
)

var (
	llm *ollama.LLM = nil
)

type Chat struct {
	sync.RWMutex

	file  *os.File              // chat file
	parts []llms.MessageContent // chat parts
	ch    chan string           // chat channel
}

func Init(model string) bool {
	var err error

	if len(model) == 0 {
		model = Default
	}

	if strings.ToLower(model) == "default" {
		model = Default
	}

	llm, err = ollama.New(ollama.WithModel(model))

	// TODO: Add embedding model with chain

	if err != nil {
		sys.Error(err)
		return false
	}

	return true
}

func NewChat() *Chat {
	return &Chat{
		file:  sys.TempFile("chat", ".txt"),
		parts: make([]llms.MessageContent, 0),
		ch:    make(chan string, 16),
	}
}

func (c *Chat) Path() string {
	return c.file.Name()
}

func (c *Chat) Close() {
	_ = c.file.Close()
}

func (c *Chat) Prompt(s string, h *heap.Heap) {
	c.write(fmt.Sprintf("%s %s\n", text.Chevron, s))
	c.human(s)

	em := make([]string, h.Lines())

	for _, str := range *h.SMap() {
		em = append(em, h.Unmap(&str))
	}

	if _, err := llm.CreateEmbedding(
		context.Background(),
		em,
	); err != nil {
		sys.Error(err)
	}

	if _, err := llm.GenerateContent(
		context.Background(),
		c.parts,
		llms.WithSeed(0),
		llms.WithTemperature(0),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			if len(chunk) > 0 {
				c.ch <- string(chunk)
			} else {
				c.ch <- "\n\n"
			}
			return nil
		}),
	); err != nil {
		sys.Error(err)
	}
}

func (c *Chat) Listen(hi *history.History) {
	var buf strings.Builder

	for s := range c.ch {
		// response start
		if buf.Len() == 0 {
			s = strings.TrimLeft(s, " ")
		}

		// response chunk
		c.write(s)
		buf.WriteString(s)

		// response end
		if s == "\n\n" {
			s = buf.String()

			c.system(s)
			hi.AddSystem(s)
			buf.Reset()
		}
	}
}

func (c *Chat) write(s string) {
	c.Lock()

	_, err := c.file.WriteString(s)

	if err != nil {
		sys.Error(err)
	}

	err = c.file.Sync()

	if err != nil {
		sys.Error(err)
	}

	c.Unlock()
}

func (c *Chat) human(s string) {
	c.history(llms.ChatMessageTypeHuman, s)
}

func (c *Chat) system(s string) {
	c.history(llms.ChatMessageTypeSystem, s)
}

func (c *Chat) history(r llms.ChatMessageType, s string) {
	c.Lock()
	c.parts = append(c.parts, llms.TextParts(r, s))
	c.Unlock()
}
