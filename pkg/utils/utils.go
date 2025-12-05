package utils

type K8sEndpoint struct {
	IP string
}

func GetEPSFromK8sSvc() []K8sEndpoint {
	endpoints := make([]K8sEndpoint, 3, 6)
	endpoints = append(endpoints, K8sEndpoint{
		IP: "10.2.12.12",
	})
	endpoints = append(endpoints, K8sEndpoint{
		IP: "10.2.12.14",
	})
	endpoints = append(endpoints, K8sEndpoint{
		IP: "10.4.12.12",
	})

	return endpoints
}
