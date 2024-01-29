package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/ssh"
)

func GetLinuxContextForTest(t *testing.T) LinuxContext {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "testcontainers/sshd:1.1.0",
		ExposedPorts: []string{"22/tcp"},
		Entrypoint:   []string{"sh", "-c", "echo 'PermitRootLogin yes'>> /etc/ssh/sshd_config && /usr/sbin/sshd && /usr/bin/tail -f /dev/null"},
		WaitingFor:   wait.ForListeningPort("22/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to create test environment: %v", err)
		t.FailNow()
	}

	ip, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get test environment ip: %v", err)
		t.FailNow()
	}

	mappedPort, err := container.MappedPort(ctx, "22")
	if err != nil {
		t.Fatalf("Failed to get test environment mapped port: %v", err)
		t.FailNow()
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("root"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, mappedPort.Int()), config)
	if err != nil {
		t.Fatalf("Failed to connect test environment: %v", err)
		t.FailNow()
	}

	session, err := conn.NewSession()
	if err != nil {
		t.Fatalf("Failed to create session for test environment: %v", err)
		t.FailNow()
	}

	return LinuxContext{
		Ctx: ctx,
		ProviderData: &LinuxProviderData{
			SshSession: session,
		},
	}
}
