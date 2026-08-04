package main

import (
	"fmt"
	"log"
	"maps"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type PeerServer struct {
	Address string
	Client  *rpc.Client
}

// when heartbeat is sent
// Go packages these fields and sends them to the other computer
type Args struct {
	GossipLive map[string]int // everything my server knows abt other servers rn
	Round      int            // sender’s logical-clock val
	Sender     string
}

type Server struct {
	live    map[string]int // Stores latest round associated w/ each known server
	lock    sync.Mutex
	Round   int
	Address string
	peers   []PeerServer
}

func (t *Server) Heartbeat(args *Args, reply *int) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	fmt.Printf("rev hb from %v\n", args.Sender)

	// update logic clock
	if args.Round > t.Round {
		t.Round = args.Round
	}

	// just heard directly from the sender
	// mark it as alive
	t.live[args.Sender] = t.Round

	// loop thru gossip only add any newer information
	for address, round := range args.GossipLive {
		if round > t.live[address] {
			t.live[address] = round
		}
	}

	return nil
}

func (t *Server) sendHeartbeat(to PeerServer) {
	t.lock.Lock()

	t.Round++

	fmt.Printf("sending hb to %v\n")
	// whose live is based on rounds
	args := &Args{
		GossipLive: maps.Clone(t.live),
		Round:      t.Round,
		Sender:     t.Address,
	}

	t.lock.Unlock()

	var reply int

	err := to.Client.Call("Server.Heartbeat", args, &reply)
	if err != nil {
		log.Println("RPC error:", err)
	}
}

func (t *Server) GenerateReport() {
	t.lock.Lock()

	log.Println("REPORT!")
	log.Println("ROUND", t.Round)
	log.Println(t.live)

	for address, round := range t.live {
		if t.Round-round <= 10 {
			fmt.Print(address)
		}
	}

	t.lock.Unlock()
}

func main() {

	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)

	my_address := "10.239.38.177:1234" //TODO: FILL IN
	server.Address = my_address
	server.Round = 0
	server.peers = make([]PeerServer, 0)
	server.live = make(map[string]int)
	peer_addresses := []string{"10.239.244.33:1234", "10.239.246.218:1234"}

	time.Sleep(10 * time.Second) // WAIT to start other servers

	for _, addr := range peer_addresses {
		if addr == my_address {
			continue
		}
		client, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		server.peers = append(server.peers, PeerServer{addr, client})
	}

	/*
		TODO: call send heartbeats to a random server every second
			- NOTE: ensure that this code is non-blocking!
		TODO: call generate report every 5 seconds
	*/

	go func() {
		for {
			server.GenerateReport()
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		server.sendHeartbeat(server.peers[rand.Intn(len(server.peers))])
		time.Sleep(1 * time.Second)
	}
}

// 10.193.25.197 shrey
//10.239.244.33. christian
// 10.239.38.177 me
// 10.239.246.218 henry
