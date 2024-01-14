package internal

import (
	"github.com/maniSHarma7575/loadbalancer/internal/controller"
)

func StartProxyServer() {
	// This should not be hardcoded. Remove it in next version
	configs := map[string]interface{}{
		"Backends": []map[string]interface{}{
			{"Host": "localhost", "Port": 8085, "IsHealthy": true, "HealthStatusUrl": "http://localhost:8085/health"},
			{"Host": "localhost", "Port": 8086, "IsHealthy": true, "HealthStatusUrl": "http://localhost:8086/health"},
			{"Host": "localhost", "Port": 8087, "IsHealthy": true, "HealthStatusUrl": "http://localhost:8087/health"},
		},
		"Strategy": "consistent_hash",
	}

	lb := controller.InitLB(configs)

	lb.Run()
}
