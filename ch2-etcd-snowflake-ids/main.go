package main

import (
	"ch2-etcd-snwoflake-ids/snowflake"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	etcd "go.etcd.io/etcd/client/v3"
)

// TODO: Use leases to handle node failures
func assignMachineId(client *etcd.Client) (snowflake.MachineID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	for range 5 {
		id := rand.IntN(2 >> 10)
		key := fmt.Sprintf("/machines/%d", id)

		txn := client.Txn(ctx).If(etcd.Compare(etcd.CreateRevision(key), "=", 0)).Else(etcd.OpPut(key, "", nil))

		resp, err := txn.Commit()
		if err != nil {
			return 0, err
		}

		if resp.Succeeded {
			return snowflake.MachineID(id), nil
		}
	}

	return 0, errors.New("assignMachineId: failed to select a unique ID")
}

func main() {
	// Setup etcd connection for unique machine ID
	etcdClient, err := etcd.New(etcd.Config{Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}, DialTimeout: 5 * time.Second})
	if err != nil {
		log.Fatalf("main: Failed to connect to etcd: %s", err.Error())
	}
	defer etcdClient.Close()

	nodeId, err := assignMachineId(etcdClient)
	if err != nil {
		log.Fatal(err)
	}
	counter := snowflake.NewCounter(nodeId)

	// Connect to maelstrom network
	node := maelstrom.NewNode()

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		snowflake := counter.Next()
		body["id"] = snowflake.Pack()

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
