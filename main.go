package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
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
			filename = capturePackets(dumpRequest.ContainerID, seconds)
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
func capturePackets(containerID string, seconds int) string {
	//seconds := 1
	//id := "ddd8aed2e0ea9ddd962deac8bad86b3aeca33be3d5161613ae77f572927fa3c0"
	pid, err := getContainerPid(containerID)
	if err != nil {
		panic(err)
	}
	num := getNetworkInterfaceId(pid)
	log.Printf("got interface num %d for container %s", num, containerID)
	interfaceName, err := getInterfaceName(num)
	log.Printf("interface name for container id %s is %s", containerID, interfaceName)
	t := time.Now()
	timeString := fmt.Sprint(t.Format("2006-01-02::15:04:05"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	filename := fmt.Sprint(containerID + "-" + timeString + ".pcap")
	cmd := exec.CommandContext(ctx, "tcpdump", "-i", interfaceName, "-w", "/tmp/"+filename)
	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("tcpdump on interface %s for %d seconds", interfaceName, seconds)
		return filename
	}
	if err != nil {
		panic(err)
	}
	return ""

}

func getInterfaceName(num int) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v\n", err.Error()))
		return "", err
	}

	for _, i := range ifaces {
		if i.Index == num {
			return i.Name, nil
		}
	}
	return "", err
}

func getContainerPid(id string) (int, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}

	containerSpec, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		panic(err)
	}
	log.Printf("got pid num %d for container %s", containerSpec.State.Pid, id)
	return containerSpec.State.Pid, nil
}

func getNetworkInterfaceId(pid int) int {
	cmd := exec.Command("nsenter", "-t", strconv.Itoa(pid), "-n", "ip", "link")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	interfaces := readStuff(scanner)
	defer cmd.Wait()
	re := regexp.MustCompile(`@if\w+`)
	match := re.FindStringSubmatch(interfaces)
	log.Printf("inteface name %s", match[0])
	interfaceNumStr := strings.Replace(match[0], "@if", "", 1)

	interfaceNum, err := strconv.Atoi(interfaceNumStr)
	if err != nil {
		panic(err)
	}
	return interfaceNum
}
func readStuff(scanner *bufio.Scanner) string {
	var str string
	for scanner.Scan() {
		str += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return str

}
