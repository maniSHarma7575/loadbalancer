package internal

import (
	"github.com/maniSHarma7575/loadbalancer/internal/controller"
)

func StartProxyServer() {
	lb := controller.InitLB()

	lb.Run()
}
