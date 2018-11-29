package internal

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"time"
)

const resolveName = "consul"

var c *api.Client

func init() {
	var err error
	conf := api.DefaultConfig()
	conf.Address = "localhost:8500"
	c, err = api.NewClient(conf)
	if err != nil {
		grpclog.Fatal(err)
	}
}

func RegisterService(name string, address string, port int) {
	grpclog.Infoln(c.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    name,
		Address: address,
		Port:    port,
	}))
}

func NewResolveBuilder(host, port string) *resolveBuilder {
	return &resolveBuilder{host, port}
}

type resolveBuilder struct {
	host string
	port string
}

func (rb *resolveBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {

	cr := &consulResolve{
		cc:     cc,
		t:      time.NewTimer(0),
		rn:     make(chan struct{}, 1),
		host:   rb.host,
		port:   rb.port,
		target: target,
	}

	go cr.watch()

	return cr, nil
}

func (*resolveBuilder) Scheme() string {
	return resolveName
}

type consulResolve struct {
	cc     resolver.ClientConn
	t      *time.Timer
	rn     chan struct{}
	stop   chan struct{}
	target resolver.Target

	host                 string
	port                 string
	disableServiceConfig bool
}

func (cr *consulResolve) ResolveNow(resolver.ResolveNowOption) {
	select {
	case cr.rn <- struct{}{}:
	default:
	}
}

func (cr *consulResolve) Close() {
	close(cr.stop)
}

func (cr *consulResolve) watch() {
	for {
		select {
		case <-cr.stop:
			return
		case <-cr.rn:
		case <-cr.t.C:
			as, _, err := c.Agent().Service(cr.target.Endpoint, nil)
			if err != nil {
				grpclog.Errorln(err)
				return
			}
			cr.cc.NewAddress([]resolver.Address{
				{
					Addr:       fmt.Sprintf("%s:%d", as.Address, as.Port),
					Type:       resolver.Backend,
					ServerName: as.Address,
				},
			})
			cr.t.Reset(5 * time.Second)
		}
	}
}
