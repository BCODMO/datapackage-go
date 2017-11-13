package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/frictionlessdata/datapackage-go/resource"
)

const (
	resourcePropName = "resources"
)

type resourceFactory func(map[string]interface{}) (*resource.Resource, error)

// Package represents a https://specs.frictionlessdata.io/data-package/
type Package struct {
	resources []*resource.Resource

	descriptor map[string]interface{}
	resFactory resourceFactory
}

// GetResource return the resource which the passed-in name or nil if the resource is not part of the package.
func (p *Package) GetResource(name string) *resource.Resource {
	for _, r := range p.resources {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// AddResource adds a new resource to the package, updating its descriptor accordingly.
func (p *Package) AddResource(d map[string]interface{}) error {
	if p.resFactory == nil {
		return fmt.Errorf("invalid resource factory. Did you mean resources.FromDescriptor?")
	}
	r, err := p.resFactory(d)
	if err != nil {
		return err
	}
	p.resources = append(p.resources, r)
	if p.descriptor == nil {
		p.descriptor = make(map[string]interface{})
	}
	p.descriptor[resourcePropName] = newResourcesDescriptor(p.resources)
	return nil
}

func newResourcesDescriptor(resources []*resource.Resource) []interface{} {
	descRes := make([]interface{}, len(resources))
	for i := range resources {
		descRes[i] = resources[i].Descriptor
	}
	return descRes
}

//RemoveResource removes the resource from the package, updating its descriptor accordingly.
func (p *Package) RemoveResource(name string) {
	index := -1
	for i := range p.resources {
		if p.resources[i].Name == name {
			index = i
			break
		}
	}
	if index != -1 {
		p.resources = append(p.resources[:index], p.resources[:index+1]...)
	}
	p.descriptor[resourcePropName] = newResourcesDescriptor(p.resources)
}

// Valid returns true if the passed-in descriptor is valid or false, otherwise.
func Valid(descriptor map[string]interface{}) bool {
	return valid(descriptor, resource.New)
}

func valid(descriptor map[string]interface{}, resFactory resourceFactory) bool {
	_, err := fromDescriptor(descriptor, resFactory)
	return err == nil
}

func fromDescriptor(descriptor map[string]interface{}, resFactory resourceFactory) (*Package, error) {
	r, ok := descriptor[resourcePropName]
	if !ok {
		return nil, fmt.Errorf("resources property is required, with at least one resource")
	}
	rSlice, ok := r.([]interface{})
	if !ok || len(rSlice) == 0 {
		return nil, fmt.Errorf("resources property is required, with at least one resource")
	}
	resources := make([]*resource.Resource, len(rSlice))
	for pos, rInt := range rSlice {
		rDesc, ok := rInt.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("resources must be a json object. got:%v", rInt)
		}
		r, err := resFactory(rDesc)
		if err != nil {
			return nil, err
		}
		resources[pos] = r
	}
	return &Package{
		resources:  resources,
		resFactory: resFactory,
		descriptor: descriptor,
	}, nil
}

// FromDescriptor creates a data package from a json descriptor.
func FromDescriptor(descriptor map[string]interface{}) (*Package, error) {
	return fromDescriptor(descriptor, resource.New)
}

func fromReader(r io.Reader, resFactory resourceFactory) (*Package, error) {
	b, err := ioutil.ReadAll(bufio.NewReader(r))
	if err != nil {
		return nil, err
	}
	var descriptor map[string]interface{}
	if err := json.Unmarshal(b, &descriptor); err != nil {
		return nil, err
	}
	return fromDescriptor(descriptor, resFactory)
}

// FromReader validates and returns a data package from an io.Reader.
func FromReader(r io.Reader) (*Package, error) {
	return fromReader(r, resource.New)
}
