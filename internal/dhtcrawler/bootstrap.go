package dhtcrawler

import (
	"context"
	"net"
	"time"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/dht/ktable"
)

func (c *crawler) reseedBootstrapNodes(ctx context.Context) {
	interval := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			for _, strAddr := range c.bootstrapNodes {
				addr, err := net.ResolveUDPAddr("udp", strAddr)
				if err != nil {
					c.logger.Warnf("failed to resolve bootstrap node address: %s", err)
					continue
				}
        // Normalize possible IPV4-mapped IPV6 Addr to IPV4
				normalizedAddrPort := netip.AddrPortFrom(
					addr.AddrPort().Addr().Unmap(),
					uint16(addr.Port),
				)
				select {
				case <-ctx.Done():
					return
				case c.nodesForPing.In() <- ktable.NewNode(ktable.ID{}, normalizedAddrPort):
					continue
				}
			}
		}

		interval = c.reseedBootstrapNodesInterval
	}
}
