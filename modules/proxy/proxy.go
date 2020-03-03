package proxy

import (
	"errors"
	"fmt"
	mods "github.com/hellgate75/go-deploy-modules/modules"
	"github.com/hellgate75/go-deploy/modules/meta"
)

type Module interface {
	GetComponent() (meta.Converter, error)
}

type module struct {
	module string
	stub   meta.ProxyStub
}

func (m *module) GetComponent() (meta.Converter, error) {
	return m.stub.Discover(m.module)
}

type Proxy interface {
	DiscoverModule(name string) (Module, error)
}

type proxy struct {
	modules map[string]meta.ProxyStub
}

func (p *proxy) DiscoverModule(name string) (Module, error) {
	fmt.Println(fmt.Sprintf("module map: %v", p.modules))
	for k, s := range p.modules {
		fmt.Println(fmt.Sprintf("module map entry: %s", k))
		if k == name {
			return &module{
				module: k,
				stub:   s,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to discover module: %s", name))
}

func getModules() map[string]meta.ProxyStub {
	return mods.GetModulesMap()
}

func NewProxy() Proxy {
	return &proxy{
		modules: getModules(),
	}
}
