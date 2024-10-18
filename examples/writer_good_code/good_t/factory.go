package good_t

/**
func (s *ProxyServer) createProxier(...) (proxy.Provider, error) {
	var proxier proxy.Provider
	if config.Mode == proxyconfigapi.ProxyModeIPTables {
		if dualStack {
			proxier, err = iptables.NewDualStackProxier()
		} else {
			proxier, err = iptables.NewProxier()
		}
	} else if config.Mode == proxyconfigapi.ProxyModeIPVS {
		if dualStack {
			proxier, err = ipvs.NewDualStackProxier()
		} else {
			proxier, err = ipvs.NewProxier()
		}
	} else if config.Mode == proxyconfigapi.ProxyModeNFTables {
		if dualStack {
			proxier, err = nftables.NewDualStackProxier()
		} else {
			proxier, err = nftables.NewProxier()
		}

	}
	return proxier, nil
}
**/
