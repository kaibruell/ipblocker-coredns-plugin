// Package ipblocker implements a plugin that blocks DNS queries based on IP addresses.
package ipblocker

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// IPBlocker is a plugin that blocks DNS queries based on IP addresses.
type IPBlocker struct {
	Next       plugin.Handler
	socketPath string
}

// ServeDNS implements the plugin.Handler interface.
func (ipb IPBlocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Get client IP from the question
	clientIP := w.RemoteAddr().String()
	if strings.Contains(clientIP, ":") {
		clientIP = strings.Split(clientIP, ":")[0]
	}

	// Check each question in the DNS request
	for _, question := range r.Question {
		domain := question.Name
		domain = strings.TrimSuffix(domain, ".")

		// Ask Java if domain is blocked
		isBlocked, err := ipb.isBlocked(clientIP, domain)
		if err != nil {
			fmt.Printf("[ipblocker] Error querying Java: %v\n", err)
			continue
		}

		if isBlocked {
			fmt.Printf("[ipblocker] BLOCKED: %s -> %s\n", clientIP, domain)
			// Create NXDOMAIN response
			m := new(dns.Msg)
			m.SetReply(r)
			m.SetRcode(r, dns.RcodeNameError)
			w.WriteMsg(m)
			return dns.RcodeNameError, nil
		}
	}

	// Pass through to next handler if not blocked
	return plugin.NextOrFailure(ipb.Name(), ipb.Next, ctx, w, r)
}

// isBlocked queries the Java service via UNIX socket
func (ipb IPBlocker) isBlocked(ip, domain string) (bool, error) {
	conn, err := net.Dial("unix", ipb.socketPath)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	// Send query to Java: "ip,domain"
	query := fmt.Sprintf("%s,%s\n", ip, domain)
	_, err = conn.Write([]byte(query))
	if err != nil {
		return false, err
	}

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(response)
	return response == "true", nil
}

// Name implements the Handler interface.
func (ipb IPBlocker) Name() string { return "ipblocker" }
