package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

const boardRows = 6
const boardCols = 7

var gameBoard [][]int
var lastPlayer int

// Mutex?
// for connect 4, mutex is used to ensure that only one player can make a move at a time, preventing race conditions and
// ensuring the integri y of game state
var mutex sync.Mutex

// helper func for making the board initially
func makeBoard() {
	gameBoard = make([][]int, boardRows) // Allocates the outer slice

	for i := range gameBoard {
		gameBoard[i] = make([]int, boardCols) // Allocates each inner row

		for j := 0; j < boardCols; j++ {
			gameBoard[i][j] = 9 // 9 means the spot is empty
		}
	}
}

// helper func so that the output of the board matches what the starter code has
func intToStringBoard() string {
	currentBoard := ""

	for i := 0; i < boardRows; i++ {
		for j := 0; j < boardCols; j++ {
			if gameBoard[i][j] == 9 {
				currentBoard += "."
			} else if gameBoard[i][j] == 1 {
				currentBoard += "B"
			} else if gameBoard[i][j] == 0 {
				currentBoard += "W"
			} else {
				currentBoard += "invalid play"
			}
		}
		currentBoard += "\n"
	}
	return currentBoard
}

type Move struct {
	Color int
	Col   int
}

type Board struct {
	BoardString string
}

type ConnectGame int

// note to self: * for Go pointer
func (t *ConnectGame) Move(args *Move, reply *int) error {
	for i := boardRows - 1; i > -1; i-- {
		if gameBoard[i][args.Col] == 9 {
			gameBoard[i][args.Col] = args.Color
			break
		}
	}
	// unsure of how to respond to client probably has to do something with reply
	*reply = 1
	return nil
}

func (t *ConnectGame) Get(args *int, reply *Board) error {
	reply.BoardString = intToStringBoard()
	return nil
}

func main() {
	makeBoard() // first make the board 6x7 of 9s

	cg := new(ConnectGame)
	rpc.Register(cg)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Println("Serving on PORT 1234")
	http.Serve(l, nil)
}
