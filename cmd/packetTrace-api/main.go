package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	packet_trace_api "github.com/ormazz/K8SpacketTrace/pkg/packetTrace-api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	port := os.Getenv("AGENTPORT")
	var err error
	var config *rest.Config
	if port == "" {
		port = "8888"
	}
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		log.Printf("no vaild kubecofig")
		panic(err)
	}
	config.TLSClientConfig.Insecure = false
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.POST("/capturePod", func(c *gin.Context) {
		var capturePodRequest packet_trace_api.CapturePodRequest
		if c.Bind(&capturePodRequest) == nil {
			containerId, containerNode := packet_trace_api.GetContainerIdAndNode(capturePodRequest, clientset)
			log.Println(containerId)
			log.Println(containerNode)
			agentIP := packet_trace_api.GetDeamonsetIp(containerNode, clientset)
			log.Println(agentIP)
			url := "http://" + agentIP + ":" + port + "/packetcapture"
			log.Printf(url)
			containerId = containerId[9:]
			var jsonStr = []byte(`{"containerId": "` + containerId + `","seconds": "2"}`)
			log.Printf(`{"containerId": "` + containerId + `","seconds": "2"}`)
			//url = "http://127.0.0.1:8888/packetcapture"
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
			if err != nil {
				log.Printf(err.Error())
				c.JSON(500, gin.H{"message": "error"})
				return
			}
			defer resp.Body.Close()

			_, err = io.Copy(c.Writer, resp.Body)
			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
		} else {
			c.JSON(500, gin.H{"message": "tddddddddump"})
		}
	})
	r.Run(":5000")

}
