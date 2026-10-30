package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"

	rest "github.com/roundinternetmoney/ethereal-rest"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pk := cmp.Or(os.Getenv("MERIDIAN_PK"), os.Getenv("ETHEREAL_PK"))
	if pk == "" {
		log.Fatal("MERIDIAN_PK is required (hex private key, with or without 0x prefix; legacy ETHEREAL_PK is also accepted)")
	}

	client, err := rest.NewClient(ctx, pk, rest.Testnet)
	if err != nil {
		log.Fatalf("failed to init meridian client: %v", err)
	}

	positions, err := client.GetPosition(ctx)
	if err != nil {
		log.Fatalf("failed to get positions: %v", err)
	}

	if len(positions) == 0 {
		fmt.Println("no open positions")
		return
	}

	for _, p := range positions {
		fmt.Println(p)
	}
}
