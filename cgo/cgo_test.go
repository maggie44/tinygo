package cgo

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestCGo(t *testing.T) {
	var cflags = []string{"--target=armv6m-none-eabi"}

	for _, name := range []string{"basic"} {
		name := name // avoid a race condition
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Read the AST in memory.
			path := filepath.Join("testdata", name+".go")
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				t.Fatal("could not parse Go source file:", err)
			}

			// Process the AST with CGo.
			cgoAST, errs := Process([]*ast.File{f}, "testdata", fset, cflags)
			for _, err := range errs {
				t.Errorf("error during CGo processing: %v", err)
			}

			// Store the (formatted) output in a buffer. Format it, so it
			// becomes easier to read (and will hopefully change less with CGo
			// changes).
			buf := &bytes.Buffer{}
			err = format.Node(buf, fset, cgoAST)
			if err != nil {
				t.Errorf("could not write out CGo AST: %v", err)
			}
			actual := strings.Replace(string(buf.Bytes()), "\r\n", "\n", -1)

			// Read the file with the expected output, to compare against.
			outfile := filepath.Join("testdata", name+".out.go")
			expectedBytes, err := ioutil.ReadFile(outfile)
			if err != nil {
				t.Fatalf("could not read expected output: %v", err)
			}
			expected := strings.Replace(string(expectedBytes), "\r\n", "\n", -1)

			// Check whether the output is as expected.
			if expected != actual {
				// It is not. Test failed.
				t.Errorf("output did not match:\n%s", string(actual))
			}
		})
	}
}
