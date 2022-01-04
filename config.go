package kube

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	// DefaultTimeout default request timeout.
	DefaultTimeout = 10 * time.Second

	// UsePersistentConfig caches client config to avoid reloads.
	UsePersistentConfig = true
)

// Config represents a kubernetes configuration.
// references:
// - https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/
// - https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/
type Config struct {
	flags *genericclioptions.ConfigFlags
	mutex *sync.RWMutex
}

// NewConfig returns a new kube client configuration.
func NewConfig(f *genericclioptions.ConfigFlags) *Config {
	return &Config{
		flags: f,
		mutex: &sync.RWMutex{},
	}
}

// Timeout returns the request timeout if set or the default if not set.
func (c *Config) Timeout() time.Duration {
	if !isset(c.flags.Timeout) {
		return DefaultTimeout
	}
	dur, err := time.ParseDuration(*c.flags.Timeout)
	if err != nil {
		return DefaultTimeout
	}

	return dur
}

// RESTConfig returns a complete rest client config.
func (c *Config) RESTConfig() (*rest.Config, error) {
	return c.clientConfig().ClientConfig()
}

// Flags returns configuration flags.
func (c *Config) Flags() *genericclioptions.ConfigFlags {
	return c.flags
}

func (c *Config) RawConfig() (clientcmdapi.Config, error) {
	return c.clientConfig().RawConfig()
}

func (c *Config) clientConfig() clientcmd.ClientConfig {
	return c.flags.ToRawKubeConfigLoader()
}

// CurrentContextName returns the currently active config context.
func (c *Config) CurrentContextName() (string, error) {
	if isset(c.flags.Context) {
		return *c.flags.Context, nil
	}
	cfg, err := c.RawConfig()
	if err != nil {
		return "", err
	}

	return cfg.CurrentContext, nil
}

// GetContext returns a context with the given name.
func (c *Config) GetContext(n string) (*clientcmdapi.Context, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return nil, err
	}
	if ctx, ok := cfg.Contexts[n]; ok {
		return ctx, nil
	}

	return nil, fmt.Errorf("context: %s not found", n)
}

// Contexts returns all available contexts.
func (c *Config) Contexts() (map[string]*clientcmdapi.Context, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Contexts, nil
}

// UseContext changes the current context.
func (c *Config) UseContext(name string) error {
	if _, err := c.GetContext(name); err != nil {
		return fmt.Errorf("context %q not found", name)
	}
	return c.useContext(name)
}

// SetContext set k/v of the given context or add a new context if not exist.
func (c *Config) SetContext(name string, ctx *clientcmdapi.Context) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	cfg.Contexts[name] = ctx
	return c.modify(cfg)
}

// RemoveContext remove the given context from the configuration.
func (c *Config) RemoveContext(n string) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	delete(cfg.Contexts, n)

	return c.modify(cfg)
}

// ClusterNameFromContext returns the cluster associated with the given context.
func (c *Config) ClusterNameFromContext(context string) (string, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return "", err
	}

	if ctx, ok := cfg.Contexts[context]; ok {
		return ctx.Cluster, nil
	}
	return "", fmt.Errorf("cluster not found from context: %s", context)
}

// CurrentClusterName returns the active cluster name.
func (c *Config) CurrentClusterName() (string, error) {
	if isset(c.flags.ClusterName) {
		return *c.flags.ClusterName, nil
	}
	cfg, err := c.RawConfig()
	if err != nil {
		return "", err
	}
	context, err := c.CurrentContextName()
	if err != nil {
		context = cfg.CurrentContext
	}

	if ctx, ok := cfg.Contexts[context]; ok {
		return ctx.Cluster, nil
	}

	return "", errors.New("current cluster not found")
}

// ClusterNames returns all kubeconfig defined clusters.
func (c *Config) ClusterNames() ([]string, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return nil, err
	}

	cc := make([]string, 0)
	for name := range cfg.Clusters {
		cc = append(cc, name)
	}

	return cc, nil
}

// SetCluster set k/v of the given cluster or add a new cluster if not exist.
func (c *Config) SetCluster(name string, cluster *clientcmdapi.Cluster) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	cfg.Clusters[name] = cluster
	return c.modify(cfg)
}

// RemoveCluster remove the given cluster from the configuration.
func (c *Config) RemoveCluster(name string) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	delete(cfg.Clusters, name)

	return c.modify(cfg)
}

// GetCluster returns a cluster with the given name.
func (c *Config) GetCluster(name string) (*clientcmdapi.Cluster, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return nil, err
	}
	if cluster, ok := cfg.Clusters[name]; ok {
		return cluster, nil
	}

	return nil, fmt.Errorf("cluster: %s not found", name)
}

// CurrentUserName retrieves the active user name.
func (c *Config) CurrentUserName() (string, error) {
	if isset(c.flags.Impersonate) {
		return *c.flags.Impersonate, nil
	}

	if isset(c.flags.AuthInfoName) {
		return *c.flags.AuthInfoName, nil
	}

	cfg, err := c.RawConfig()
	if err != nil {
		return "", err
	}

	current := cfg.CurrentContext
	if isset(c.flags.Context) {
		current = *c.flags.Context
	}
	if ctx, ok := cfg.Contexts[current]; ok {
		return ctx.AuthInfo, nil
	}

	return "", errors.New("current user not found")
}

// SetCredential set k/v of the given credential or add a new credential if not exist.
func (c *Config) SetCredential(name string, auth *clientcmdapi.AuthInfo) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	cfg.AuthInfos[name] = auth
	return c.modify(cfg)
}

// RemoveCredential remove the given credential from the configuration.
func (c *Config) RemoveCredential(name string) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	delete(cfg.AuthInfos, name)

	return c.modify(cfg)
}

// GetCredential returns a credential with the given name.
func (c *Config) GetCredential(name string) (*clientcmdapi.AuthInfo, error) {
	cfg, err := c.RawConfig()
	if err != nil {
		return nil, err
	}
	if auth, ok := cfg.AuthInfos[name]; ok {
		return auth, nil
	}

	return nil, fmt.Errorf("credential: %s not found", name)
}

// CurrentNamespace retrieves the active namespace.
func (c *Config) CurrentNamespace() (string, error) {
	ns, _, err := c.clientConfig().Namespace()

	return ns, err
}

// ConfigAccess returns the current kubeconfig api server access configuration.
func (c *Config) ConfigAccess() (clientcmd.ConfigAccess, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.clientConfig().ConfigAccess(), nil
}

func (c *Config) useContext(name string) error {
	cfg, err := c.RawConfig()
	if err != nil {
		return err
	}
	cfg.CurrentContext = name

	return c.modify(cfg)
}

func (c *Config) modify(cfg clientcmdapi.Config) error {
	acc, err := c.ConfigAccess()
	if err != nil {
		return err
	}
	return clientcmd.ModifyConfig(acc, cfg, true)
}
