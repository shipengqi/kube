package kube

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

const (
	// EnvVarKubeConfig The KUBECONFIG environment variable holds a list of kubeconfig files.
	// For Linux and Mac, the list is colon-delimited. For Windows, the list is semicolon-delimited.
	EnvVarKubeConfig          = "KUBECONFIG"
	DefaultKubeHomeDir        = ".kube"
	DefaultKubeConfigFileName = "config"
)

var (
	defaultYamlDelimiter = []byte("---")
	defaultConfigDir     = filepath.Join(homedir.HomeDir(), DefaultKubeHomeDir)
	defaultHomeFile      = filepath.Join(defaultConfigDir, DefaultKubeConfigFileName)
)

// RetrievesDefaultKubeConfig returns a available kubeconfig file.
func RetrievesDefaultKubeConfig() string {
	if env, ok := os.LookupEnv(EnvVarKubeConfig); ok && len(env) > 0 {
		list := filepath.SplitList(env)
		d := deduplicate(list)
		if len(d) > 0 {
			return d[0]
		}
	} else {
		return defaultHomeFile
	}
	return ""
}

// LoadDefaultKubeConfig starts by running the clientcmdapi.MigrationRules and then
// takes the loading rules and returns a clientcmdapi.Config object.
func LoadDefaultKubeConfig() (*clientcmdapi.Config, error) {
	dc := clientcmd.NewDefaultClientConfigLoadingRules()
	return dc.Load()
}

// GetObjects returns the list of objects parsed from the given files.
func GetObjects(files ...string) ([]runtime.Object, error) {
	var objs []runtime.Object
	for _, f := range files {
		fBytes, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		subs, err := getObjects(fBytes)
		if err != nil {
			return nil, err
		}
		if len(subs) > 0 {
			objs = append(objs, subs...)
		}
	}
	return objs, nil
}

func getObjects(content []byte) ([]runtime.Object, error) {
	objs := make([]runtime.Object, 0)

	delimited := bytes.Split(content, defaultYamlDelimiter)
	for _, del := range delimited {
		if len(del) == 0 {
			continue
		}
		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode(del, nil, nil)
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
		// switch o := obj.(type) {
		// case *corev1.Pod:
		// default:
		// }
	}
	return objs, nil
}

func applyObject(helper *resource.Helper, namespace, name string, obj runtime.Object) error {
	if _, err := helper.Get(namespace, name); err != nil {
		_, err = helper.Create(namespace, false, obj)
		if err != nil {
			return err
		}
	} else {
		_, err = helper.Replace(namespace, name, true, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteObject(helper *resource.Helper, namespace, name string) error {
	if _, err := helper.Get(namespace, name); err == nil {
		_, err = helper.Delete(namespace, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func retrievesMetaFromObject(obj runtime.Object) (namespace, name string, err error) {
	name, err = meta.NewAccessor().Name(obj)
	if err != nil {
		return
	}
	namespace, err = meta.NewAccessor().Namespace(obj)
	if err != nil {
		return
	}
	return
}

// deduplicate removes any duplicated values and returns a new slice.
func deduplicate(s []string) []string {
	encountered := map[string]bool{}
	ret := make([]string, 0)
	for i := range s {
		if len(s[i]) == 0 {
			continue
		}
		if encountered[s[i]] {
			continue
		}
		encountered[s[i]] = true
		ret = append(ret, s[i])
	}
	return ret
}

// isset reports whether the given string pointer has a value.
func isset(s *string) bool {
	return s != nil && len(*s) != 0
}

func listOptions(label []string) metav1.ListOptions {
	opts := metav1.ListOptions{}
	if len(label) > 0 {
		opts.LabelSelector = strings.Join(label, ",")
	}
	return opts
}
