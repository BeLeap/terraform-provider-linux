package util

import (
	"context"
	"testing"
	"time"

	"github.com/melbahja/goph"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/ssh"
)

func GetLinuxContext(t *testing.T) LinuxContext {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "testcontainers/sshd:1.1.0",
		ExposedPorts: []string{"22/tcp"},
		Entrypoint:   []string{"sh", "-c", "echo 'PermitRootLogin yes'>> /etc/ssh/sshd_config && /usr/sbin/sshd && /usr/bin/tail -f /dev/null"},
		WaitingFor:   wait.ForListeningPort("22/tcp").WithStartupTimeout(10 * time.Second),
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

	auth := goph.Password("root")
	sshClient, err := goph.NewConn(&goph.Config{
		User:     "root",
		Addr:     ip,
		Port:     uint(mappedPort.Int()),
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Fatalf("Failed to connect test environment: %v", err)
		t.FailNow()
	}

	return LinuxContext{
		Ctx: ctx,
		ProviderData: &LinuxProviderData{
			SshClient: sshClient,
		},
	}
}
