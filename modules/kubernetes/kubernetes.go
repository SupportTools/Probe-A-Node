package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Node struct {
	Name      string
	Status    bool
	PublicIp  string
	PrivateIp string
}

func SetupKubernetes() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return client
}

func GetNodeInfo(clientset *kubernetes.Clientset, nodeName string) (Node, error) {

	node, err := clientset.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return Node{}, err
	}
	var status bool
	if string(node.Status.Conditions[0].Status) == "True" {
		status = true
	} else {
		status = false
	}
	return Node{
		Name:      node.Name,
		Status:    status,
		PublicIp:  node.Status.Addresses[0].Address,
		PrivateIp: node.Status.Addresses[1].Address,
	}, nil
}
