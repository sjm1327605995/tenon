package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run scripts/api_snapshot.go <output>")
		os.Exit(2)
	}
	if err := writeSnapshot(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeSnapshot(output string) error {
	var lines []string
	for _, directory := range []string{"pkg/ui", "pkg/value", "pkg/shadcn", "pkg/decl", "pkg/router", "yoga"} {
		entries, err := os.ReadDir(directory)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}
			path := filepath.Join(directory, entry.Name())
			parsed, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.SkipObjectResolution)
			if err != nil {
				return err
			}
			for _, decl := range parsed.Decls {
				for _, line := range exportedDecls(decl) {
					lines = append(lines, parsed.Name.Name+": "+line)
				}
			}
		}
	}
	sort.Strings(lines)
	return os.WriteFile(output, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

func exportedDecls(decl ast.Decl) []string {
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		if !ast.IsExported(decl.Name.Name) {
			return nil
		}
		copy := *decl
		copy.Body = nil
		return []string{formatNode(&copy)}
	case *ast.GenDecl:
		var lines []string
		for _, spec := range decl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				if ast.IsExported(spec.Name.Name) {
					copy := *spec
					copy.Type = publicType(spec.Type)
					lines = append(lines, decl.Tok.String()+" "+formatNode(&copy))
				}
			case *ast.ValueSpec:
				copy := *spec
				copy.Names = nil
				for _, name := range spec.Names {
					if ast.IsExported(name.Name) {
						copy.Names = append(copy.Names, name)
					}
				}
				if len(copy.Names) > 0 {
					lines = append(lines, decl.Tok.String()+" "+formatNode(&copy))
				}
			}
		}
		return lines
	default:
		return nil
	}
}

func publicType(expr ast.Expr) ast.Expr {
	switch expr := expr.(type) {
	case *ast.StructType:
		copy := *expr
		fields := *expr.Fields
		fields.List = nil
		for _, field := range expr.Fields.List {
			if publicField(field) {
				fields.List = append(fields.List, field)
			}
		}
		copy.Fields = &fields
		return &copy
	default:
		return expr
	}
}

func publicField(field *ast.Field) bool {
	for _, name := range field.Names {
		if ast.IsExported(name.Name) {
			return true
		}
	}
	return len(field.Names) == 0 && publicEmbeddedType(field.Type)
}

func publicEmbeddedType(expr ast.Expr) bool {
	switch expr := expr.(type) {
	case *ast.Ident:
		return ast.IsExported(expr.Name)
	case *ast.SelectorExpr:
		return ast.IsExported(expr.Sel.Name)
	case *ast.StarExpr:
		return publicEmbeddedType(expr.X)
	case *ast.IndexExpr:
		return publicEmbeddedType(expr.X)
	default:
		return false
	}
}

func formatNode(node any) string {
	var out strings.Builder
	if err := format.Node(&out, token.NewFileSet(), node); err != nil {
		panic(err)
	}
	return strings.Join(strings.Fields(out.String()), " ")
}
