package main

import (
	"github.com/cenkalti/rpc2"
	"log"
	"net"
	"time"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:5000"
)

func main() {
	main2()

	log.Println("wait exit...")
	select {}
}

func main2() {
	type Args struct{ A, B int }
	type Reply int

	// server
	lis, err := net.Listen(network, addr)
	if err != nil {
		log.Println(err)
		return
	}

	var ccc *rpc2.Client

	srv := rpc2.NewServer()
	srv.SetTimeout(3)
	srv.Handle("timeout", func(client *rpc2.Client, args *Args, reply *Reply) error {
		ccc = client

		*reply = Reply(args.A + args.B)

		time.Sleep(time.Second * 5)

		return nil
	})
	srv.Handle("add", func(client *rpc2.Client, args *Args, reply *Reply) error {
		*reply = Reply(args.A + args.B)

		return nil
	})
	go srv.Accept(lis)

	// client
	conn, err := net.Dial(network, addr)
	if err != nil {
		log.Println(err)
		return
	}

	clt := rpc2.NewClient(conn)
	clt.Handle("timeout", func(client *rpc2.Client, args *Args, reply *Reply) error {
		*reply = Reply(args.A * args.B)

		time.Sleep(time.Second * 10)
		return nil
	})
	clt.Handle("mult", func(client *rpc2.Client, args *Args, reply *Reply) error {
		*reply = Reply(args.A * args.B)
		return nil
	})
	clt.SetTimeout(3)
	go clt.Run()
	defer clt.Close()

	// Test Call.
	var rep Reply
	err = clt.Call("timeout", Args{1, 2}, &rep)
	if err != nil {
		log.Println(err)
	}
	log.Println(rep)

	// Test Call.
	var rep2 Reply
	err = ccc.Call("timeout", Args{1, 2}, &rep2)
	if err != nil {
		log.Println(err)
	}
	log.Println(rep2)
}
