package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetDeployment returns a Deployment with the given name.
func (c *Client) GetDeployment(ctx context.Context, namespace, name string) (*v1.Deployment, error) {
	return c.client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetDeployments returns a DeploymentList.
func (c *Client) GetDeployments(ctx context.Context, namespace string, label ...string) (*v1.DeploymentList, error) {
	return c.client.AppsV1().Deployments(namespace).List(ctx, listOptions(label))
}

// CreateDeployment creates a new Deployment.
func (c *Client) CreateDeployment(ctx context.Context, deploy *v1.Deployment) (*v1.Deployment, error) {
	if len(deploy.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.AppsV1().Deployments(deploy.Namespace).Create(ctx, deploy, metav1.CreateOptions{})
}

// DeleteDeployment deletes a Deployment.
func (c *Client) DeleteDeployment(ctx context.Context, namespace, name string) error {
	return c.client.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetDaemonSet returns a DaemonSet with given name.
func (c *Client) GetDaemonSet(ctx context.Context, namespace, name string) (*v1.DaemonSet, error) {
	return c.client.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetDaemonSets returns a DaemonSetList.
func (c *Client) GetDaemonSets(ctx context.Context, namespace string, label ...string) (*v1.DaemonSetList, error) {
	return c.client.AppsV1().DaemonSets(namespace).List(ctx, listOptions(label))
}

// CreateDaemonSet creates a new DaemonSet.
func (c *Client) CreateDaemonSet(ctx context.Context, dsData *v1.DaemonSet) (*v1.DaemonSet, error) {
	if len(dsData.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.AppsV1().DaemonSets(dsData.Namespace).Create(ctx, dsData, metav1.CreateOptions{})
}

// DeleteDaemonSet deletes a DaemonSet.
func (c *Client) DeleteDaemonSet(ctx context.Context, namespace, name string) error {
	return c.client.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetStatefulSet returns a StatefulSet with given name.
func (c *Client) GetStatefulSet(ctx context.Context, namespace, name string) (*v1.StatefulSet, error) {
	return c.client.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetStatefulSets returns a StatefulSetList.
func (c *Client) GetStatefulSets(ctx context.Context, namespace string, label ...string) (*v1.StatefulSetList, error) {
	return c.client.AppsV1().StatefulSets(namespace).List(ctx, listOptions(label))
}

// CreateStatefulSet creates a new StatefulSet.
func (c *Client) CreateStatefulSet(ctx context.Context, statefulSet *v1.StatefulSet) (*v1.StatefulSet, error) {
	if len(statefulSet.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.AppsV1().StatefulSets(statefulSet.Namespace).Create(ctx, statefulSet, metav1.CreateOptions{})
}

// DeleteStatefulSet deletes a StatefulSet.
func (c *Client) DeleteStatefulSet(ctx context.Context, namespace, name string) error {
	return c.client.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetJob returns a Job with given name.
func (c *Client) GetJob(ctx context.Context, namespace, name string) (*batchv1.Job, error) {
	return c.client.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetJobs returns a JobList.
func (c *Client) GetJobs(ctx context.Context, namespace string, label ...string) (*batchv1.JobList, error) {
	return c.client.BatchV1().Jobs(namespace).List(ctx, listOptions(label))
}

// CreateJob creates a new Job.
func (c *Client) CreateJob(ctx context.Context, jobData *batchv1.Job) (*batchv1.Job, error) {
	if len(jobData.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.BatchV1().Jobs(jobData.Namespace).Create(ctx, jobData, metav1.CreateOptions{})
}

// DeleteJob deletes a Job.
func (c *Client) DeleteJob(ctx context.Context, namespace, name string) error {
	return c.client.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// GetCronJob returns a GetCronJob with given name.
func (c *Client) GetCronJob(ctx context.Context, namespace, name string) (*v1beta1.CronJob, error) {
	return c.client.BatchV1beta1().CronJobs(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetCronJobs returns a CronJobList.
func (c *Client) GetCronJobs(ctx context.Context, namespace string, label ...string) (*v1beta1.CronJobList, error) {
	return c.client.BatchV1beta1().CronJobs(namespace).List(ctx, listOptions(label))
}

// CreateCronJob creates a new CronJob.
func (c *Client) CreateCronJob(ctx context.Context, cronjob *v1beta1.CronJob) (*v1beta1.CronJob, error) {
	if len(cronjob.Namespace) == 0 {
		return nil, ErrorMissingNamespace
	}
	return c.client.BatchV1beta1().CronJobs(cronjob.Namespace).Create(ctx, cronjob, metav1.CreateOptions{})
}

// DeleteCronJob deletes a CronJob.
func (c *Client) DeleteCronJob(ctx context.Context, namespace, name string) error {
	return c.client.BatchV1beta1().CronJobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
