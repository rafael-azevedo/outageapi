package utils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

const (
	// Status id iterperetation
	StatusSuccess    = 0
	StatusFailed     = 1
	StatusPending    = 2
	StatusProcessing = 3

	// HPOM constants
	OVBin        = "/opt/OV/bin"
	NetworkIP    = "net_type=NETWORK_IP"
	NodeGroup    = "group_name=outage"
	OPCNode      = OVBin + "/opcnode"
	ScriptRepo   = "/var/opt/OV/SPLS_scripts/"
	OutageScript = ScriptRepo + "outage.ksh"
)

var statusText = map[int]string{
	StatusPending:    "Pending",
	StatusProcessing: "Processing",
	StatusFailed:     "Failed",
	StatusSuccess:    "Successs",
}

func StatusText(code int) string {
	return statusText[code]
}

type SubRequest struct {
	Host   string
	Status int
}

type OutageRequest struct {
	ID           int
	Action       string       `json:"action"`
	UserName     string       `json:"username"`
	Email        string       `json:"email"`
	ChangeTicket string       `json:"changeticket"`
	IP           string       `json:"ip"`
	Status       int          `json:"status"`
	ServerList   []SubRequest `json:"serverlist"`
}

func GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func CreateRequest(r *http.Request) (OutageRequest, error) {
	var or OutageRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return or, err
	}
	if err := r.Body.Close(); err != nil {
		return or, err
	}
	if err := json.Unmarshal(body, &or); err != nil {
		return or, err
	}
	or.Status = 2
	or.IP = GetIP(r)

	LID, err := LastID()
	if err != nil {
		log.Println(err)
	}

	or.ID = LID + 1

	for i := range or.ServerList {
		or.ServerList[i].Status = 2
	}
	return or, nil
}

func (or OutageRequest) LogRequest() error {
	logFile := os.Getenv("OUTAGELOGDIR") + "/" + strconv.Itoa(or.ID) + ".json"
	log.Printf("%s : %s\n", "Logging to ", logFile)
	logtext, err := toJson(or)
	if err != nil {
		log.Println()
	}
	if err := ioutil.WriteFile(logFile, logtext, 0644); err != nil {
		return err
	}
	return nil
}

func toJson(o interface{}) ([]byte, error) {
	bytes, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ReadOutageLog(id int) (OutageRequest, error) {
	logFile := os.Getenv("OUTAGELOGDIR") + "/" + strconv.Itoa(id) + ".json"
	log.Println(logFile)

	var or OutageRequest

	raw, err := ioutil.ReadFile(logFile)
	if err != nil {
		return or, err
	}

	json.Unmarshal(raw, &or)

	return or, err
}

func (or *OutageRequest) Assign() error {
	or.Status = 3
	for i, item := range or.ServerList {
		out, err := AssignNodeScript(item.Host)
		if err != nil {
			or.ServerList[i].Status = 1
		} else if err == nil {
			or.ServerList[i].Status = 0
		}
		log.Printf("%s\n%s\n", item.Host, out)
	}
	or.Status = 1
	err := or.LogRequest()
	if err != nil {
		return err
	}
	return nil
}

func (or *OutageRequest) Deassign() error {
	or.Status = 3
	for i, item := range or.ServerList {
		out, err := DeassignNodeScript(item.Host)
		if err != nil {
			or.ServerList[i].Status = 1
		} else if err == nil {
			or.ServerList[i].Status = 0
		}
		log.Printf("%s\n%s\n", item.Host, out)
	}
	or.Status = 1
	err := or.LogRequest()
	if err != nil {
		return err
	}
	return nil
}

func AssignNodeCMD(node string) ([]byte, error) {
	nodeName := "node_name=" + node
	args := []string{"-assign_node", nodeName, NodeGroup, NetworkIP}
	log.Println(OPCNode, NetworkIP, nodeName, args)
	cmd := exec.Command(OPCNode, args...)
	log.Println(cmd)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}

func DeassignNodeCMD(node string) ([]byte, error) {
	nodeName := "node_name=" + node
	args := []string{"-deassign_node", nodeName, NodeGroup, NetworkIP}
	cmd := exec.Command(OPCNode, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}

func TestAssignNode(node string) ([]byte, error) {

	switch rand.Intn(2) {
	case 0:
		return []byte("You have assigend " + node), nil
	}
	return []byte("You have  failed to assigend " + node), errors.New("You have  failed to assigend " + node)
}

func TestDeassignNode(node string) ([]byte, error) {

	switch rand.Intn(2) {
	case 0:
		return []byte("You have deassigend " + node), nil
	}
	return []byte("You have  failed to deassigend " + node), errors.New("You have  failed to deassigend " + node)
}
