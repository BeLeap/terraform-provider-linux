package util

import (
	"context"
	"testing"

	"github.com/melbahja/goph"
)

func GetLinuxContext(t *testing.T) LinuxContext {
	auth, err := goph.UseAgent()
	if err != nil {
		t.Fatalf("Failed to create auth info: %v", err)
		t.FailNow()
	}

	sshClient, err := goph.New("root", "test-node.fox-deneb.ts", auth)
	if err != nil {
		t.Fatalf("Failed to connect test node: %v", err)
		t.FailNow()
	}

	return LinuxContext{
		Ctx: context.Background(),
		ProviderData: &LinuxProviderData{
			SshClient: sshClient,
		},
	}
}
