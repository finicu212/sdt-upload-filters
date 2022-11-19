package connection

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
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

func (o *Orchestrator) Handle(ctx context.Context, files []file.FileDetails) error {
	for _, p := range o.pools {
		q := p.Q()
		q.AddToQueue(files)
	}
	return o.handleQueue(ctx)
}

// handleQueue feeds files to the connection pools that have available connections.
// If a connection pool is currently filled, skip it.
// Once we iterate through all the connection pools, wait 5s, and start this process once again with only non-empty queues
func (o *Orchestrator) handleQueue(ctx context.Context) error {
	fmt.Println("handleQueue")
	errs, ctx := errgroup.WithContext(ctx)
	q := o.dropPoolsWithEmptyQueues()
	for ; len(q) > 0; q = o.dropPoolsWithEmptyQueues() {
		for k, p := range o.pools {
			fmt.Printf("IP %s\n", k)
			newConn, err := p.GetConnection(ctx)
			if err != nil {
				if err == ErrConnectionLimit {
					fmt.Printf("Skipping pool for IP %s as it is currently busy\n", k)
					continue
				}
				return err
			}
			q := p.Q()
			nextFile := q.PopQueue()
			errs.Go(func() error {
				return newConn.Store(nextFile.RemotePath, nextFile.DataReader)
			})
			fmt.Printf("Ran Store(%s)\n", nextFile.RemotePath)
			if err != nil {
				return err
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errs.Wait()
}

func (o *Orchestrator) dropPoolsWithEmptyQueues() map[string]IPool {
	for k, p := range o.pools {
		q := p.Q()
		if len(q.GetQueue()) == 0 {
			fmt.Printf("Dropped empty queue over IP %s\n", k)
			delete(o.pools, k)
		}
	}
	return o.pools
}
