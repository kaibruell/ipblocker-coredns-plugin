package ipblocker

import (
	"fmt"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	for i := 0; i < 10; i++ {
		fmt.Println("Hello CoreDNS Container from ipblocker plugin!")
	}
	plugin.Register("ipblocker", setup)
}

func setup(c *caddy.Controller) error {
	if c.Next() {
		// Parse configuration here if needed in the future
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return IPBlocker{Next: next}
	})

	return nil
}
