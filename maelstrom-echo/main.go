// /
// / Challenge #1: Echo
// / https://fly.io/dist-sys/1/
// /
package main

import (
	"encoding/json"
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

	// Register a handler for echo
	n.Handle(
		"echo", // Handle messages with type "echo"
		func(msg maelstrom.Message) error {
			// Input is essentially {'src': <>, 'dest': <>, 'body': {...}}
			var body map[string]any

			if err := json.Unmarshal(msg.Body, &body); err != nil {
				return err
			}

			body["type"] = "echo_ok"

			// Send the same message back with the new body
			return n.Reply(msg, body)
		})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
