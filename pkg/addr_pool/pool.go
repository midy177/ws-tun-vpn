package addr_pool

import "sync"

type AddressPool struct {
	pool map[string]bool
	mask string
	mu   sync.Mutex
}

// NewAddressPool Create a new address pool with a list of addresses and a mask.
func NewAddressPool(addrs []string, mask string) *AddressPool {
	pool := make(map[string]bool)
	for _, addr := range addrs {
		// address pool determines whether it is in use, default not in use.
		pool[addr] = false
	}
	return &AddressPool{
		pool: pool,
		mask: mask,
	}
}

// GetAddressFromPool Get an available address randomly
func (a *AddressPool) GetAddressFromPool() (string, string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for addr, inUse := range a.pool {
		if !inUse {
			a.pool[addr] = true
			return addr, a.mask
		}
	}
	return "", a.mask
}

// PutAddressToPool Put an address back to the pool
func (a *AddressPool) PutAddressToPool(str string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.pool[str] = false
}

// GetMask get mask
func (a *AddressPool) GetMask() string {
	return a.mask
}
