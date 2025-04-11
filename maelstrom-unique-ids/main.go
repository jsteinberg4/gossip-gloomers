// /
// / Challenge #1: Echo
// / https://fly.io/dist-sys/1/
// /
package main

import (
	"log"
	"log/slog"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	// Instantiate node type
	n := maelstrom.NewNode()
	logger.Info("Node created")

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
