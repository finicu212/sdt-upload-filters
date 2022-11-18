package connection

import (
	"io"
	"sdt-upload-filters/pkg/utils"
	"time"
)

// Orchestrator is responsible for sending commands to the pools.
type Orchestrator struct {
	pools map[string]IPool // [IP -> Pool]
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

func (o *Orchestrator) AddToQueue(files []io.Reader) {
	for _, p := range o.pools {
		p.AddToQueue(files)
	}
}

// HandleQueue feeds files to the connection pools that have available connections.
// If a connection pool is currently filled, skip it.
// Once we iterate through all the connection pools, wait 5s, and start this process once again with only non-empty queues
func (o *Orchestrator) HandleQueue() error {
	q := o.DropPoolsWithEmptyQueues()
	for ; len(q) > 0; q = o.DropPoolsWithEmptyQueues() {
		for _, p := range o.pools {
			newConn, err := p.GetConnection()
			if err != nil {
				if err == ErrConnectionLimit {
					continue
				}
				return err
			}
			nextFile := p.PopQueue()
			err = newConn.Store("TODO", nextFile)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func (o *Orchestrator) DropPoolsWithEmptyQueues() map[string]IPool {
	for k, p := range o.pools {
		if len(p.GetQueue()) == 0 {
			delete(o.pools, k)
		}
	}
	return o.pools
}
