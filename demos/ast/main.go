package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	src := `
        package main
        
        import "fmt"
        
        var x int = 42

        func hello() {
            fmt.Println("Hello, World!")
        }
        
        func add(a int, b int) int {
            return a + b
        }
    `

	// 创建一个文件集，表示 Go 源代码的文件集
	fset := token.NewFileSet()

	// 解析源代码，得到 AST
	node, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 使用 ast.Inspect 递归遍历 AST 节点
	ast.Inspect(node, func(n ast.Node) bool {
		// 提取函数声明
		if fn, ok := n.(*ast.FuncDecl); ok {
			fmt.Printf("Function name: %s\n", fn.Name.Name)
			fmt.Println("Parameters:")
			for _, param := range fn.Type.Params.List {
				for _, name := range param.Names {
					fmt.Printf("\tName: %s, Type: %s\n", name.Name, param.Type)
				}
			}
			if fn.Type.Results != nil {
				fmt.Println("Return types:")
				for _, result := range fn.Type.Results.List {
					fmt.Printf("\tType: %s\n", result.Type)
				}
			}
		}

		// 提取变量声明
		if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if vspec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vspec.Names {
						fmt.Printf("Variable name: %s, Type: %s\n", name.Name, vspec.Type)
					}
				}
			}
		}

		return true
	})
}
