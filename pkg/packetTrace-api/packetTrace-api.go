package packet_trace_api

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CapturePodRequest struct {
	PodName   string `json:"podName"`
	Namespace string `json:"Namespace"`
	Seconds   string `json:"Seconds`
}

func GetContainerIdAndNode(req CapturePodRequest, clientset kubernetes.Interface) (string, string) {
	api := clientset.CoreV1()
	var ns, label, field string
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	pods, err := api.Pods(ns).List(listOptions)
	if err != nil {
		log.Panicln(err)
		return "", ""
	}
	for _, pod := range pods.Items {
		if pod.Name == req.PodName {
			return pod.Status.ContainerStatuses[0].ContainerID, pod.Spec.NodeName
		}
	}

	return "", ""
}

func GetNodeIp(nodeName string, clientset kubernetes.Interface) string {
	api := clientset.CoreV1()
	var label, field string
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	nodes, err := api.Nodes().List(listOptions)
	if err != nil {
		return ""
	}
	for _, node := range nodes.Items {
		if node.Name == nodeName {
			return node.Status.Addresses[0].Address
		}
	}

	return ""
}
