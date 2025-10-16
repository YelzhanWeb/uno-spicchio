package main

import (
	"fmt"

	"github.com/YelzhanWeb/uno-spicchio/pkg/hash"
)

func main() {
	passwords := []string{"admin123", "manager123", "waiter123", "cook123"}
	for _, p := range passwords {
		h, _ := hash.Hash(p)
		fmt.Printf("%s → %s\n", p, h)
	}
}
