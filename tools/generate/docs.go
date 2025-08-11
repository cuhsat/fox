package main

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/hiforensics/fox/cmd"
	"github.com/hiforensics/fox/internal/fox"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: docs.go DIR")
	}

	mp := path.Join(os.Args[1], "man")
	md := path.Join(os.Args[1], "markdown")

	cobra.CheckErr(doc.GenMarkdownTree(cmd.Fox, md))
	cobra.CheckErr(doc.GenManTree(cmd.Fox, &doc.GenManHeader{
		Source:  fox.Product,
		Title:   "Fox",
		Section: "1",
	}, mp))
}
