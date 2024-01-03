package main

import "github.com/maniSHarma7575/loadbalancer/lib"

func main() {
	lb := lib.InitLB()

	lb.Run()
}
