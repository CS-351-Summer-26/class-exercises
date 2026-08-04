package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Move struct {
	Color int
	Col   int
}

type Board struct {
	BoardString string
}

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// First send one move to the server
	var move Move

	fmt.Print("Enter color (0 = White, 1 = Black): ")
	fmt.Scan(&move.Color)

	fmt.Print("Enter column (0-6): ")
	fmt.Scan(&move.Col)

	var replyMove int

	err = client.Call("ConnectGame.Move", &move, &replyMove)
	if err != nil {
		log.Fatal("move error:", err)
	}

	log.Printf("Move reply: %v", replyMove)

	// Then ask the server for the board
	var reply Board
	var args int

	err = client.Call("ConnectGame.Get", &args, &reply)
	if err != nil {
		log.Fatal("game error:", err)
	}

	log.Printf("Game:\n%v", reply.BoardString)
}
