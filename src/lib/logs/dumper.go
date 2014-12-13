// Debug beautiful output
package logs

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
)

// dump Dump all variables to bytes.Buffer
func dump(idl ...interface{}) bytes.Buffer {
	var buf bytes.Buffer
	var wr io.Writer

	wr = io.MultiWriter(&buf)
	for _, field := range idl {
		fset := token.NewFileSet()
		ast.Fprint(wr, fset, field, ast.NotNilFilter)
	}
	return buf
}

// Dumper Dump all variables to STDOUT
func Dumper(idl ...interface{}) {
	ret := dump(idl...)
	fmt.Print(ret.String())
}

// DumperString Dump all variables to string
func DumperString(idl ...interface{}) string {
	ret := dump(idl...)
	return ret.String()
}

// DumperByte Dump all variables to []byte
func DumperByte(idl ...interface{}) []byte {
	ret := dump(idl...)
	return ret.Bytes()
}
