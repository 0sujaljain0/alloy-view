package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/a-h/templ"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// The updated structure for the final output, including identifiers
type K8sEndpoint struct {
	IP       string
	PodName  string // The name of the Pod this IP belongs to
	NodeName string // The name of the Node the Pod is running on
}

// GetEPSFromK8sSvc retrieves the IP addresses and identifiers from EndpointSlices
// associated with a specific Kubernetes Service in a given namespace.
// It now takes namespace and serviceName as arguments.
func GetEPSFromK8sSvc(namespace string, serviceName string, logger *slog.Logger) []K8sEndpoint {
	var endpoints []K8sEndpoint
	ctx := context.Background()

	// --- 1. Initialize Kubernetes Client ---
	// Load kubeconfig from the standard location for out-of-cluster execution.
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		// Fallback for non-home directory environments
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Info(fmt.Errorf("Error building kubeconfig: %v\n", err.Error()).Error())
		return nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Info(fmt.Errorf("Error creating clientset: %v\n", err.Error()).Error())
		return nil
	}

	// --- 2. Define Label Selector ---
	// EndpointSlices are linked to a Service via the standard label:
	// kubernetes.io/service-name=<service-name>
	labelSelector := fmt.Sprintf("kubernetes.io/service-name=%s", serviceName)

	// --- 3. List EndpointSlices ---
	// Access the discovery/v1 API for EndpointSlices in the specified namespace
	epsClient := clientset.DiscoveryV1().EndpointSlices(namespace)

	listOptions := metav1.ListOptions{
		LabelSelector: labelSelector,
	}

	epsList, err := epsClient.List(ctx, listOptions)
	if err != nil {
		logger.Error(fmt.Errorf("Error listing EndpointSlices for %s/%s: %v\n", namespace, serviceName, err.Error()).Error())
		return nil
	}

	// --- 4. Extract IP Addresses and Identifiers ---
	if len(epsList.Items) == 0 {
		logger.Info("No EndpointSlices found for service %s in namespace %s\n", serviceName, namespace)
		return nil
	}

	for _, eps := range epsList.Items {
		// Iterate over all endpoints within the EndpointSlice
		for _, endpoint := range eps.Endpoints {

			// Extract Pod Name (from TargetRef)
			podName := ""
			if endpoint.TargetRef != nil && endpoint.TargetRef.Kind == "Pod" {
				// The Name field of TargetRef holds the Pod name
				podName = endpoint.TargetRef.Name
			}

			// Extract Node Name
			nodeName := ""
			if endpoint.NodeName != nil {
				nodeName = *endpoint.NodeName // NodeName is a pointer to string
			}

			// An Endpoint object can have multiple IP addresses (Addresses)
			for _, ip := range endpoint.Addresses {
				endpoints = append(endpoints, K8sEndpoint{
					IP:       ip,
					PodName:  podName,
					NodeName: nodeName,
				})
			}
		}
	}

	logger.Info(fmt.Sprintf("Found %d endpoints for service %s/%s\n", len(endpoints), namespace, serviceName))
	sort.Slice(endpoints, func(i, j int) bool {
		// Use strings.Compare for standard string comparison, or strings.ToLower
		// if you need case-insensitive sorting.
		return strings.Compare(endpoints[i].PodName, endpoints[j].PodName) < 0
	})
	return endpoints
}

func CreateTemplAttrs(input ...any) (templ.Attributes, error) {
	if len(input) == 0 {
		return nil, errors.New("while creating templ.Attributes: no attributes given")
	}
	if len(input)%2 != 0 {
		return nil, errors.New("while creating templ.Attributes: odd no. of arguments")
	}

	var res templ.Attributes = templ.Attributes{}

	for i := 1; i < len(input); i += 2 {
		key, ok := input[i-1].(string)
		if !ok {
			return nil, fmt.Errorf("while creating templ.Attributes: %v is not a string", input[i])
		}
		res[key] = input[i]
	}

	return res, nil
}
