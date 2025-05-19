// /
// / Challenge #1: Echo
// / https://fly.io/dist-sys/1/
// /
package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func UniqueId() (uuid.UUID, error) {
	if id, err := uuid.NewRandom(); err != nil {
		return uuid.Nil, err
	} else {
		return id, nil
	}
}

func main() {
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	uuid.EnableRandPool()

	// Instantiate node type
	n := maelstrom.NewNode()
	// logger.Info("Node created")

	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"

		id, err := UniqueId()
		if err != nil {
			return err
		}

		body["id"] = id.String()

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
