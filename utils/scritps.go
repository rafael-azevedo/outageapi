package utils

import "os/exec"

func AssignNodeScript(node string) ([]byte, error) {
	nodeName := "node_name=" + node
	args := []string{"ADD", nodeName, "Request Via Outage API"}
	cmd := exec.Command(OutageScript, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}

func DeassignNodeScript(node string) ([]byte, error) {
	nodeName := "node_name=" + node
	args := []string{"REMOVE", nodeName, "Request Via Outage API"}
	cmd := exec.Command(OutageScript, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
