package main

import (
	"flag"
	"log/slog"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" {
		slog.Error("You must specify the paths to the source and target files")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		slog.Error("Copy error:", "error", err)
		os.Exit(1)
	}
}
