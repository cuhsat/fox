package rag

import (
    "context"
    "fmt"
    "runtime"
    "strconv"
    "strings"

    "github.com/philippgille/chromem-go"

    "github.com/hiforensics/fox/internal/pkg/sys"
    "github.com/hiforensics/fox/internal/pkg/types/heap"
)

type RAG struct {
    db *chromem.DB // in-memory database
}

func New() *RAG {
    return &RAG{
        db: chromem.NewDB(),
    }
}

func (rag *RAG) Embed(h *heap.Heap) *chromem.Collection {
    fn := chromem.NewEmbeddingFuncOllama("nomic-embed-text", "")

    col, err := rag.db.GetOrCreateCollection("fox", nil, fn)

    if err != nil {
        sys.Error(err)
        return nil
    }

    var docs []chromem.Document

    for i, str := range *h.FMap() {
        docs = append(docs, chromem.Document{
            ID:       strconv.Itoa(i),
            Metadata: map[string]string{"path": h.Base},
            Content:  fmt.Sprintf("line %d: %s", str.Nr, str.Str),
        })
    }

    err = col.AddDocuments(context.Background(), docs, runtime.NumCPU())

    if err != nil {
        sys.Error(err)
        return nil
    }

    return col
}

func (rag *RAG) Query(query string, col *chromem.Collection) string {
    res, err := col.Query(context.Background(), query, col.Count(), nil, nil)

    if err != nil {
        sys.Error(err)
        return ""
    }

    var sb strings.Builder

    for _, r := range res {
        sb.WriteString(r.Content)
        sb.WriteRune('\n')
    }

    return sb.String()
}
