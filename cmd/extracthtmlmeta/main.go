package main

import (
	"encoding/json"
	"log"
	"os"

	"git.rappet.de/index/htmlmeta"
)

func main() {
	meta, err := htmlmeta.CreatePageMeta(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(meta)
}
