package connection

import (
	"io"
	"sdt-upload-filters/pkg/utils"
)

// Orchestrator is responsible for sending commands to the pools.
type Orchestrator struct {
	pools  map[string]IPool
	queues map[IPool][]io.Reader // Each IPool has its queue of files in case the pool is full
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

func (o *Orchestrator) buildQueues(files []io.Reader) {
	for _, v := range o.pools {
		o.queues[v] = files
	}
}
