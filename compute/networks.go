package compute

import (
	"errors"
	"net/url"

	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// NetworkProperties ...
type NetworkProperties struct {
	Name string
	Type string
	Href string
}

// Network ...
type Network struct {
	API        apiv1.Network
	Properties NetworkProperties
	Network    *abiquo.Network
}

// Networks ...
type Networks map[string]*Network

func ip(n *Network, query url.Values) (l *core.Link, err error) {
	ip := n.Network.Rel("ips").Collection(query).First()
	if ip == nil {
		return nil, errors.New("IP not available on " + n.Network.Name)
	}
	return ip.Link(), nil
}

func dynamic(n *Network) (l *core.Link, err error) {
	return ip(n, url.Values{
		"available": []string{"true"},
		"free":      []string{"true"},
	})
}

func manual(n *Network) (l *core.Link, err error) {
	return ip(n, url.Values{
		"available": []string{"true"},
		"free":      []string{"true"},
		"has":       []string{n.API.IP()},
	})
}

func vip(n *Network) (l *core.Link, err error) {
	return nil, errors.New("VIP networks are not supported")
}

var getLink = map[string]func(*Network) (*core.Link, error){
	"dynamic": dynamic,
	"manual":  manual,
	"vip":     vip,
}

// NewNetwork ...
func newNetwork(apiNetwork apiv1.Network) (network *Network, err error) {
	var properties NetworkProperties
	err = apiNetwork.CloudProps().As(&properties)
	if err != nil {
		return
	}

	net := new(abiquo.Network)
	err = core.NewLink(properties.Href).SetType("vlan").Read(net)
	if err != nil {
		return
	}

	network = &Network{
		API:        apiNetwork,
		Network:    net,
		Properties: properties,
	}

	return
}

// NewNetworks ...
func newNetworks(apiNetworks apiv1.Networks) (networks Networks, err error) {
	networks = Networks{}
	for name, apiNetwork := range apiNetworks {
		network, errInt := newNetwork(apiNetwork)
		networks[name] = network
		if errInt != nil {
			err = errInt
		}
	}

	return networks, err
}

// IP ...
func (n *Network) IP() string {
	return n.API.IP()
}

// Netmask ...
func (n *Network) Netmask() string {
	return n.API.Netmask()
}

// Gateway ...
func (n *Network) Gateway() string {
	return n.API.Gateway()
}

// SetMAC ...
func (n *Network) SetMAC(mac string) {
	n.API.SetMAC(mac)
}

// CloudPropertyName ...
func (n *Network) CloudPropertyName() string {
	return n.Properties.Name
}

// CloudPropertyType ...
func (n *Network) CloudPropertyType() string {
	return n.Properties.Type
}

// AsNetworks ...
func (ns Networks) AsNetworks() apiv1.Networks {
	newNets := apiv1.Networks{}
	for name, net := range ns {
		newNets[name] = net.API
	}
	return newNets
}
