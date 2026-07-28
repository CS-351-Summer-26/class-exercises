package main

import (
	"fmt"
	"log"
	"net/rpc"
	// "time"
)

type Args struct {
	Key      string
	Value    string
	ClientID string
}

type Reply struct {
	Success bool
	Value   string
}

// Asks the server for exclusive control over one key
func Lock(client *rpc.Client, client_id string, key string) {
	args := &Args{Key: key, ClientID: client_id}
	var reply Reply
	reply.Success = false

	for reply.Success != true {
		err := client.Call("KVServer.Lock", args, &reply)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Tells the server that you are finished using that key
func Unlock(client *rpc.Client, client_id string, key string) {
	args := &Args{Key: key, ClientID: client_id}
	var reply Reply
	err := client.Call("KVServer.Unlock", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

// Asks the server for the current value of a key
func Get(client *rpc.Client, key string) string {
	args := &Args{Key: key}
	var reply Reply
	err := client.Call("KVServer.Get", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	return reply.Value
}

// Sends a request to the server to change a value
func Put(client *rpc.Client, key string, value string) {
	args := &Args{Key: key, Value: value}
	var reply Reply
	err := client.Call("KVServer.Put", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
}

// This function should get/put values into the key value store
func run_txn(client *rpc.Client, client_id string) {
	Put(client, "x", "10") // have partner put 20
	// time.Sleep() = time.Millisecond()
	Put(client, "y", "20") // and 40
}

// This function should get values from the keyvalue store server
// This should check if you get serializable results and print an error if the results are unserializable.
// inspect the server and determine whether it contains a valid or invalid combination
func run_read_txn(client *rpc.Client, client_id string) {
	x := Get(client, "x")
	y := Get(client, "y")

	// checks for the two mixed combinations
	// first: x came from my transaction, y came from partner’s transaction
	// second: x came from partner’s transaction, y came frommy transaction
	if (x != "10" && y == "20") || (x != "20" && y == "40") {
		fmt.Println("ERROR: unserializable result:", "x =", x, "y =", y)
	}

	// if (x == "10" && y == "20") || (x == "20" && y == "40") {
	// 	fmt.Println("serializable!!", "x: ", x, "y: ", y)
	// }

	fmt.Println(x, y)
}

// (x != "30" && y == "60") ||

func main() {
	lockserver_address := "10.239.51.175:1234"
	client_id := "anna"

	client, err := rpc.DialHTTP("tcp", lockserver_address)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	for range 1000 {
		run_txn(client, client_id)
		run_read_txn(client, client_id)
	}
}
