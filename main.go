package main

import "github.com/maniSHarma7575/loadbalancer/internal"

func main() {
	lb := internal.InitLB()

	lb.Run()
}
