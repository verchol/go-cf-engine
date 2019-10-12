package docker

import "testing"

func TestRun(t *testing.T) {
	args := map[string]string{
		"commands": "sleep, 4",
	}
	RunDocker("docker.io/library/alpine", args)
}
