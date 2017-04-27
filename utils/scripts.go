package utils

import "os/exec"

func AssignNodeScript(node string) ([]byte, error) {
	args := []string{"ADD", node, "Request Via Outage API"}
	cmd := exec.Command(OutageScript, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}

func DeassignNodeScript(node string) ([]byte, error) {
	args := []string{"REMOVE", node, "Request Via Outage API"}
	cmd := exec.Command(OutageScript, args...)
	out, err := cmd.Output()
	if err != nil {
		return out, err
	}
	return out, nil
}
