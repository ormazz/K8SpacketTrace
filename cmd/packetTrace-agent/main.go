package main

import (
	"github.com/gin-gonic/gin"
	packet_trace_agent "github.com/ormazz/K8SpacketTrace/pkg/packetTrace-agent"
	"log"
	"strconv"
)

type dumpRequest struct {
	ContainerID string `json:"containerId"`
	Seconds     string `json:"seconds"`
}

func main() {

	r := gin.Default()
	r.POST("/packetcapture", func(c *gin.Context) {
		var dumpRequest dumpRequest
		if c.Bind(&dumpRequest) == nil {
			log.Println(dumpRequest.ContainerID)
			log.Println(dumpRequest.Seconds)
			c.Header("Content-Description", "File Transfer")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Content-Type", "application/octet-stream")
			seconds, err := strconv.Atoi(dumpRequest.Seconds)
			if err != nil {
				c.JSON(500, gin.H{"shittt": "fuckk"})
				return
			}
			var filename string
			filename = packet_trace_agent.CapturePackets(dumpRequest.ContainerID, seconds)
			c.Header("Content-Disposition", "attachment; filename="+filename)
			log.Printf("/tmp/" + filename)
			c.FileAttachment("/tmp/"+filename, "capture")
			log.Printf("send file")
			//remove file
		} else {
			c.JSON(500, gin.H{"message": "tddddddddump"})
		}
	})
	r.Run(":8888")

}
