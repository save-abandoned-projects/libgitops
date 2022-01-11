//go:build ignore
// +build ignore

/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_resource_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	"fmt"

	api "API_DIR"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage"
	"github.com/weaveworks/libgitops/pkg/storage/filterer"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceClient is an interface for accessing Resource-specific API objects
type ResourceClient interface {
	// New returns a new Resource
	New() *api.Resource
	// Get returns the Resource matching given UID from the storage
	Get(runtime.UID) (*api.Resource, error)
	// Set saves the given Resource into persistent storage
	Set(*api.Resource) error
	// Patch performs a strategic merge patch on the object with
	// the given UID, using the byte-encoded patch given
	Patch(runtime.UID, []byte) error
	// Find returns the Resource matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filterer.BaseFilter) (*api.Resource, error)
	// FindAll returns multiple Resources matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filterer.BaseFilter) ([]*api.Resource, error)
	// Delete deletes the Resource with the given UID from the storage
	Delete(uid runtime.UID) error
	// List returns a list of all Resources available
	List() ([]*api.Resource, error)
}

// Resources returns the ResourceClient for the Client object
func (c *Client) Resources() ResourceClient {
	if c.resourceClient == nil {
		c.resourceClient = newResourceClient(c.storage, c.gv)
	}

	return c.resourceClient
}

// resourceClient is a struct implementing the ResourceClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type resourceClient struct {
	storage  storage.Storage
	filterer *filterer.Filterer
	gvk      schema.GroupVersionKind
}

// newResourceClient builds the resourceClient struct using the storage implementation and a new Filterer
func newResourceClient(s storage.Storage, gv schema.GroupVersion) ResourceClient {
	return &resourceClient{
		storage:  s,
		filterer: filterer.NewFilterer(s),
		gvk:      gv.WithKind(api.KindResource.Title()),
	}
}

// New returns a new Object of its kind
func (c *resourceClient) New() *api.Resource {
	log.Tracef("Client.New; GVK: %v", c.gvk)
	obj, err := c.storage.New(c.gvk)
	if err != nil {
		panic(fmt.Sprintf("Client.New must not return an error: %v", err))
	}
	return obj.(*api.Resource)
}

// Find returns a single Resource based on the given Filter
func (c *resourceClient) Find(filter filterer.BaseFilter) (*api.Resource, error) {
	log.Tracef("Client.Find; GVK: %v", c.gvk)
	object, err := c.filterer.Find(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	return object.(*api.Resource), nil
}

// FindAll returns multiple Resources based on the given Filter
func (c *resourceClient) FindAll(filter filterer.BaseFilter) ([]*api.Resource, error) {
	log.Tracef("Client.FindAll; GVK: %v", c.gvk)
	matches, err := c.filterer.FindAll(c.gvk, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Resource, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Resource))
	}

	return results, nil
}

// Get returns the Resource matching given UID from the storage
func (c *resourceClient) Get(uid runtime.UID) (*api.Resource, error) {
	log.Tracef("Client.Get; UID: %q, GVK: %v", uid, c.gvk)
	object, err := c.storage.Get(c.gvk, uid)
	if err != nil {
		return nil, err
	}

	return object.(*api.Resource), nil
}

// Set saves the given Resource into the persistent storage
func (c *resourceClient) Set(resource *api.Resource) error {
	log.Tracef("Client.Set; UID: %q, GVK: %v", resource.GetUID(), c.gvk)
	return c.storage.Set(c.gvk, resource)
}

// Patch performs a strategic merge patch on the object with
// the given UID, using the byte-encoded patch given
func (c *resourceClient) Patch(uid runtime.UID, patch []byte) error {
	return c.storage.Patch(c.gvk, uid, patch)
}

// Delete deletes the Resource from the storage
func (c *resourceClient) Delete(uid runtime.UID) error {
	log.Tracef("Client.Delete; UID: %q, GVK: %v", uid, c.gvk)
	return c.storage.Delete(c.gvk, uid)
}

// List returns a list of all Resources available
func (c *resourceClient) List() ([]*api.Resource, error) {
	log.Tracef("Client.List; GVK: %v", c.gvk)
	list, err := c.storage.List(c.gvk)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Resource, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Resource))
	}

	return results, nil
}
