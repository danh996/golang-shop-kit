package registry

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "8500"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalied uri")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")

	//Singleton mode
	builderInstance = &consulBuilder{}
)

func Init() {
	fmt.Printf("calling consul init\n")
	resolver.Register(CacheBuilder())
}

type consulBuilder struct {
}

type consulResolver struct {
	address              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	Ch                   chan int
}

func NewBuilder() resolver.Builder {
	return &consulBuilder{}
}

func CacheBuilder() resolver.Builder {
	return builderInstance
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))
	if err != nil {
		fmt.Println("parse err")
		return nil, err
	}
	fmt.Println(fmt.Sprintf("consul service ==> host:%s, port%s, name:%s", host, port, name))
	cr := &consulResolver{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		Ch:                   make(chan int, 0),
	}
	go cr.watcher()
	return cr, nil

}

func (cr *consulResolver) watcher() {
	fmt.Printf("calling [%s] consul watcher\n", cr.name)
	config := api.DefaultConfig()
	config.Address = cr.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}
	t := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		fmt.Println("defer done")
	}()
	for {
		select {
		case <-t.C:
			//fmt.Println("Timing")
		case <-cr.Ch:
			//fmt.Println("ch call")
		}
		//api added lastIndex consul api is not compatible with query with lastIndex
		services, _, err := client.Health().Service(cr.name, "", true, &api.QueryOptions{})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		newAddrs := make([]resolver.Address, 0)
		for _, service := range services {
			addr := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
			newAddrs = append(newAddrs, resolver.Address{
				Addr: addr,
				//type: cannot be grpclb, grpclb will delete the last link address when processing the link, no need to set it. See => balancer_conn_wrappers => updateClientConnState
				ServerName: service.Service.Service,
			})
		}
		//cr.cc.NewAddress(newAddrs)
		//cr.cc.NewServiceConfig(cr.name)
		cr.cc.UpdateState(resolver.State{Addresses: newAddrs})
	}

}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	cr.Ch <- 1
}

func (cr *consulResolver) Close() {
}

func parseTarget(target string) (host, port, name string, err error) {

	if target == "" {
		return "", "", "", errMissingAddr
	}

	if !regexConsul.MatchString(target) {
		return "", "", "", errAddrMisMatch
	}

	groups := regexConsul.FindStringSubmatch(target)
	host = groups[1]
	port = groups[2]
	name = groups[3]
	if port == "" {
		port = defaultPort
	}
	return host, port, name, nil
}
