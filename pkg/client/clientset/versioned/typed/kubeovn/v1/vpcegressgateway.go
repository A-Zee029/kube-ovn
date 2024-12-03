/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"

	v1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	scheme "github.com/kubeovn/kube-ovn/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// VpcEgressGatewaysGetter has a method to return a VpcEgressGatewayInterface.
// A group's client should implement this interface.
type VpcEgressGatewaysGetter interface {
	VpcEgressGateways(namespace string) VpcEgressGatewayInterface
}

// VpcEgressGatewayInterface has methods to work with VpcEgressGateway resources.
type VpcEgressGatewayInterface interface {
	Create(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.CreateOptions) (*v1.VpcEgressGateway, error)
	Update(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.UpdateOptions) (*v1.VpcEgressGateway, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, vpcEgressGateway *v1.VpcEgressGateway, opts metav1.UpdateOptions) (*v1.VpcEgressGateway, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.VpcEgressGateway, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.VpcEgressGatewayList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VpcEgressGateway, err error)
	VpcEgressGatewayExpansion
}

// vpcEgressGateways implements VpcEgressGatewayInterface
type vpcEgressGateways struct {
	*gentype.ClientWithList[*v1.VpcEgressGateway, *v1.VpcEgressGatewayList]
}

// newVpcEgressGateways returns a VpcEgressGateways
func newVpcEgressGateways(c *KubeovnV1Client, namespace string) *vpcEgressGateways {
	return &vpcEgressGateways{
		gentype.NewClientWithList[*v1.VpcEgressGateway, *v1.VpcEgressGatewayList](
			"vpc-egress-gateways",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1.VpcEgressGateway { return &v1.VpcEgressGateway{} },
			func() *v1.VpcEgressGatewayList { return &v1.VpcEgressGatewayList{} }),
	}
}
