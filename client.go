package kube

import (
	"errors"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	defaultApiUriPath  = "/api"
	defaultApisUriPath = "/apis"
)

var ErrorMissingNamespace = errors.New("missing namespace")

// Client represents kubernetes Client.
type Client struct {
	client    kubernetes.Interface
	metrics   *versioned.Clientset
	cfg       *Config
	proxy     func(request *http.Request) (*url.URL, error)
	inCluster bool
}

// New returns a new Client.
func New(cfg *Config) *Client {
	return &Client{cfg: cfg}
}

// NewDefault returns a new Client with default kubeconfig.
func NewDefault() *Client {
	flags := genericclioptions.NewConfigFlags(UsePersistentConfig)
	return &Client{cfg: NewConfig(flags)}
}

// NewInCluster returns a new Client with InClusterConfig.
func NewInCluster() *Client {
	return &Client{inCluster: true}
}

// WithProxy set proxy. Proxy is the proxy func to be used for all requests made by this client.
// If Proxy is nil, http.ProxyFromEnvironment is used. If Proxy returns a nil *URL, no proxy is used.
func (c *Client) WithProxy(fn func(request *http.Request) (*url.URL, error)) *Client {
	c.proxy = fn
	return c
}

// DialMetrics returns a client to the metrics server.
func (c *Client) DialMetrics() (*versioned.Clientset, error) {
	if c.metrics != nil {
		return c.metrics, nil
	}
	cfg, err := c.RestConfig()
	if err != nil {
		return nil, err
	}
	if c.metrics, err = versioned.NewForConfig(cfg); err != nil {
		return nil, err
	}
	return c.metrics, nil
}

// Dial returns a client to the kubernetes apiserver.
func (c *Client) Dial() (kubernetes.Interface, error) {
	if c.client != nil {
		return c.client, nil
	}
	cfg, err := c.RestConfig()
	if err != nil {
		return nil, err
	}
	if c.client, err = kubernetes.NewForConfig(cfg); err != nil {
		return nil, err
	}
	return c.client, nil
}

// RestConfig returns a complete rest client config.
func (c *Client) RestConfig() (*rest.Config, error) {
	var (
		cfg *rest.Config
		err error
	)
	if c.inCluster {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = c.cfg.RESTConfig()
		if err != nil {
			return nil, err
		}
	}
	if c.proxy != nil {
		cfg.Proxy = c.proxy
	}
	cfg.Timeout = c.cfg.Timeout()

	return cfg, nil
}

// ResourceClient returns a client to the given schema.GroupVersion.
func (c *Client) ResourceClient(gv schema.GroupVersion) (rest.Interface, error) {
	cfg, err := c.RestConfig()
	if err != nil {
		return nil, err
	}
	cfg.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	cfg.GroupVersion = &gv
	if len(gv.Group) == 0 {
		cfg.APIPath = defaultApiUriPath
	} else {
		cfg.APIPath = defaultApisUriPath
	}
	return rest.RESTClientFor(cfg)
}
