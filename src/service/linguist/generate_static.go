// +build ignore

package main
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

func main() {
	if err := bake(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bake() error {
	f, err := os.Create("static.go")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "%v\n\npackage linguist\n\n", warning)
	fmt.Fprintf(w, "var files = map[string]string{\n")
	for i := 1; i < len(os.Args); i++ {
		fn := os.Args[i]
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%q: ", fn)
		if utf8.Valid(b) {
			fmt.Fprintf(w, "`%s`", sanitize(b))
		} else {
			fmt.Fprintf(w, "%q", b)
		}
		fmt.Fprintf(w, ",\n\n")
	}
	fmt.Fprintln(w, "}")
	if err := w.Flush(); err != nil {
		return err
	}
	return f.Close()
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

const warning = "// DO NOT EDIT ** This file was generated with the bake tool ** DO NOT EDIT //"
