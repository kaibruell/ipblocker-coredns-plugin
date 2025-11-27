// Package ipblocker implements a plugin that blocks DNS queries based on IP addresses.
package ipblocker

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// IPBlocker is a plugin that blocks DNS queries based on IP addresses.
type IPBlocker struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (ipb IPBlocker) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Hello World: Just pass through to the next handler for now
	return plugin.NextOrFailure(ipb.Name(), ipb.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (ipb IPBlocker) Name() string { return "ipblocker" }
