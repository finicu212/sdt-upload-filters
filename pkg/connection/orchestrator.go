package connection

import (
	"sdt-upload-filters/pkg/utils"
)

type Orchestrator struct {
	pools map[string]IPool
}

func NewOrchestrator(ips []string, usernames []string, passwords []string) (o *Orchestrator, err error) {
	o = new(Orchestrator)
	o.pools = make(map[string]IPool)
	for idx, ip := range ips {
		url, port, err := utils.SplitHostPort(ip)
		if err != nil {
			return nil, err
		}
		o.pools[ip] = NewPool(usernames[idx], passwords[idx], url, port)
	}
	return o, nil
}
