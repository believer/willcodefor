package utils

import (
	"bytes"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/anchor"
)

type ASTTransformer struct{}

// Custom AST transformer for own purposes. Probably all sorts of
// bugs in here due to missing data, but it works for now.
func (g *ASTTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		switch v := n.(type) {
		case *ast.Link:
			link := v.Destination

			// If the link is an external link, add the target and rel attributes
			if len(link) > 0 && bytes.HasPrefix(link, []byte("http")) {
				v.SetAttributeString("target", []byte("_blank"))
				v.SetAttributeString("rel", []byte("noopener noreferrer"))
			}
		}

		return ast.WalkContinue, nil
	})
}

func MarkdownToHTML(input []byte) bytes.Buffer {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(
				util.Prioritized(&ASTTransformer{}, 1000),
			),
		),
		goldmark.WithExtensions(
			&anchor.Extender{
				Attributer: anchor.Attributes{
					"class": "!text-gray-400 dark:!text-gray-500 no-underline",
				},
				Texter: anchor.Text("#"),
			},
			extension.Strikethrough,
			extension.Typographer,
			extension.NewFootnote(
				extension.WithFootnoteBacklinkClass([]byte("font-mono no-underline")),
			),
			extension.Table,
			highlighting.NewHighlighting(
				highlighting.WithStyle("base16-snazzy"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	if err := md.Convert(input, &buf); err != nil {
		panic(err)
	}

	return buf
}
