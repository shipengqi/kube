package kube

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/restmapper"
)

type Operation int

const (
	OperationApply Operation = iota
	OperationDelete
)

// Apply is like kubectl apply -f
// references:
// - https://github.com/kubernetes/client-go/issues/193
// - https://stackoverflow.com/questions/58783939/using-client-go-to-kubectl-apply-against-the-kubernetes-api-directly-with-mult
func (c *Client) Apply(files []string) error {
	return c.execute(OperationApply, files)
}

// Delete is like kubectl delete -f
func (c *Client) Delete(files []string) error {
	return c.execute(OperationDelete, files)
}

func (c *Client) execute(op Operation, files []string) error {
	objs, err := GetObjects(files)
	if err != nil {
		return err
	}
	// Create a REST mapper that tracks information about the available resources in the cluster.
	groupResources, err := restmapper.GetAPIGroupResources(c.client.Discovery())
	if err != nil {
		return err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	for i := range objs {
		// Get some metadata needed to make the REST request.
		gvk := objs[i].GetObjectKind().GroupVersionKind()
		gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
		mapping, err := mapper.RESTMapping(gk, gvk.Version)
		if err != nil {
			return err
		}
		namespace, name, err := retrievesMetaFromObject(objs[i])
		if err != nil {
			return err
		}
		cli, err := c.ResourceClient(mapping.GroupVersionKind.GroupVersion())
		if err != nil {
			return err
		}
		helper := resource.NewHelper(cli, mapping)
		switch op {
		case OperationApply:
			err = applyObject(helper, namespace, name, objs[i])
			if err != nil {
				return err
			}
		case OperationDelete:
			err = deleteObject(helper, namespace, name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
