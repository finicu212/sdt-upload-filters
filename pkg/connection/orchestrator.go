package connection

import (
	"sdt-upload-filters/pkg/file"
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

func (o *Orchestrator) addToQueue(files []file.FileDetails) {
	for _, p := range o.pools {
		p.Q().AddToQueue(files)
	}
}

func (o *Orchestrator) Handle(files []file.FileDetails) error {
	o.addToQueue(files)
	return o.handleQueue()
}

// handleQueue feeds files to the connection pools that have available connections.
// If a connection pool is currently filled, skip it.
// Once we iterate through all the connection pools, wait 5s, and start this process once again with only non-empty queues
func (o *Orchestrator) handleQueue() error {
	q := o.dropPoolsWithEmptyQueues()
	for ; len(q) > 0; q = o.dropPoolsWithEmptyQueues() {
		for _, p := range o.pools {
			newConn, err := p.GetConnection()
			if err != nil {
				if err == ErrConnectionLimit {
					continue
				}
				return err
			}
			nextFile := p.Q().PopQueue()
			err = newConn.Store("TODO", nextFile.DataReader)
			if err != nil {
				return err
			}
			time.Sleep(5 * time.Second)
		}
	}
	return nil
}

func (o *Orchestrator) dropPoolsWithEmptyQueues() map[string]IPool {
	for k, p := range o.pools {
		if len(p.Q().GetQueue()) == 0 {
			delete(o.pools, k)
		}
	}
	return o.pools
}
