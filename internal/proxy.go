package internal

import (
	"github.com/maniSHarma7575/loadbalancer/internal/controller"
)

func StartProxyServer() {
	// This should not be hardcoded. Remove it in next version
	configs := map[string]interface{}{
		"Backends": []map[string]interface{}{
			{"Host": "localhost", "Port": 8085, "IsHealthy": true},
			{"Host": "localhost", "Port": 8086, "IsHealthy": true},
			{"Host": "localhost", "Port": 8087, "IsHealthy": true},
		},
		"Strategy": "consistent_hash",
	}

	lb := controller.InitLB(configs)

	lb.Run()
}
