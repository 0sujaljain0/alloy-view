package state

type AlloyNode interface {
	GetNodeName() string
}

type K8sAlloyNode struct {
	PodName  string
	NodeName string
	IP       string
}

func (n *K8sAlloyNode) GetNodeName() string { return n.PodName }
