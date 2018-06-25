package deployment

import (
	"github.com/ericchiang/k8s"
	"github.com/pkg/errors"
)

// Scale changes the number of app instances
func Scale(namespace string, app string, n *int32) error {
	c, err := k8s.NewInClusterClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s client")
	}

	d, err := c.ExtensionsV1Beta1().GetDeployment(ctx, app, namespace)
	if err != nil {
		return errors.Wrap(err, "failed to get deployment")
	}

	d.Spec.Replicas = n

	_, err = c.ExtensionsV1Beta1().UpdateDeployment(ctx, d)
	return errors.Wrap(err, "failed to scale deployment")
}
