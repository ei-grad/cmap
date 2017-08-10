package main // import "github.com/ei-grad/cmap/cmap-gen"

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	output        = flag.String("output", "", "output file name; default srcdir/<type>_cmap.go")
	typeName      = flag.String("type", "", "value type name (required)")
	mapTypeName   = flag.String("map", "", "map type name (default: TypeMap)")
	shardTypeName = flag.String("shard", "", "shard type name (default: TypeShard)")
	keyTypeName   = flag.String("key", "string", "key type name")
	newMethodName = flag.String("new", "", "the New() method name (default: NewTypeMap)")
	packageName   = flag.String("package", "", "package name (required)")
)

type templateParams struct {
	TypeName      string
	MapTypeName   string
	ShardTypeName string
	KeyTypeName   string
	NewMethodName string
	Package       string
	Cmdline       string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("cmap-gen: ")
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	var dir string

	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		dir = filepath.Dir(args[0])
	}

	params := getParams()

	var buf bytes.Buffer

	err := tmpl.Execute(&buf, params)
	if err != nil {
		log.Fatalf("rendering template: %s", err)
	}

	// Format the output.
	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		panic(fmt.Errorf("warning: internal error: invalid Go generated: %s", err))
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_cmap.go", strings.ToLower(*typeName))
		outputName = filepath.Join(dir, baseName)
	}
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

func getParams() (params templateParams) {

	params.TypeName = *typeName
	params.Package = *packageName
	params.Cmdline = strings.Join(os.Args, " ")

	if *mapTypeName == "" {
		params.MapTypeName = *typeName + "Map"
	} else {
		params.MapTypeName = *mapTypeName
	}

	if *shardTypeName == "" {
		params.ShardTypeName = *typeName + "Shard"
	} else {
		params.ShardTypeName = *shardTypeName
	}

	params.KeyTypeName = *keyTypeName

	if *newMethodName == "" {
		params.NewMethodName = "New" + *typeName + "Map"
	} else {
		params.NewMethodName = *newMethodName
	}

	return

}
