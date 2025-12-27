package state

import "fmt"

type AlloyNode interface {
	GetNodeName() string
	GetEndpoint() string
	String() string
}

type K8sAlloyNode struct {
	PodName  string
	NodeName string
	IP       string
}

func (n *K8sAlloyNode) String() string {
	return fmt.Sprintf("[ (kubernetes) {%s|%s|%s} ]", n.IP, n.PodName, n.NodeName)
}

func (n *K8sAlloyNode) GetNodeName() string { return n.PodName }

func (n *K8sAlloyNode) GetEndpoint() string { return fmt.Sprintf("%s:12345", n.IP) }
