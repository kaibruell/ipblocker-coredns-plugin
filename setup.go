package ipblocker

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

const (
	socketPath = "/tmp/ipblocker.sock"
	javaJar    = "/coredns-launcher.jar"
)

var javaProcess *exec.Cmd

func init() {
	fmt.Println("IPBlocker plugin initializing...")
	plugin.Register("ipblocker", setup)
}

func setup(c *caddy.Controller) error {
	if c.Next() {
		// Parse configuration here if needed in the future
	}

	// Start Java service
	if err := startJavaService(); err != nil {
		fmt.Printf("[ipblocker] Error starting Java service: %v\n", err)
		// Don't fail - continue without blocking
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return IPBlocker{
			Next:       next,
			socketPath: socketPath,
		}
	})

	return nil
}

func startJavaService() error {
	fmt.Println("[ipblocker] Starting Java service...")

	// Check if JAR exists
	if _, err := os.Stat(javaJar); err != nil {
		return fmt.Errorf("Java JAR not found at %s: %w", javaJar, err)
	}

	javaProcess = exec.Command("java", "-jar", javaJar)
	javaProcess.Stdout = os.Stdout
	javaProcess.Stderr = os.Stderr

	if err := javaProcess.Start(); err != nil {
		return fmt.Errorf("failed to start Java process: %w", err)
	}

	fmt.Printf("[ipblocker] Java process started (PID: %d)\n", javaProcess.Process.Pid)

	// Give Java time to start the socket server
	time.Sleep(2 * time.Second)

	return nil
}
