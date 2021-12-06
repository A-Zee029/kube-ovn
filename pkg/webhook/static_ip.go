package webhook

import (
	"context"
	"net/http"

	ovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	deploymentGVK  = metav1.GroupVersionKind{Group: appsv1.SchemeGroupVersion.Group, Version: appsv1.SchemeGroupVersion.Version, Kind: "Deployment"}
	statefulSetGVK = metav1.GroupVersionKind{Group: appsv1.SchemeGroupVersion.Group, Version: appsv1.SchemeGroupVersion.Version, Kind: "StatefulSet"}
	daemonSetGVK   = metav1.GroupVersionKind{Group: appsv1.SchemeGroupVersion.Group, Version: appsv1.SchemeGroupVersion.Version, Kind: "DaemonSet"}
	podGVK         = metav1.GroupVersionKind{Group: corev1.SchemeGroupVersion.Group, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}
	subnetGVK      = metav1.GroupVersionKind{Group: ovnv1.SchemeGroupVersion.Group, Version: ovnv1.SchemeGroupVersion.Version, Kind: "Subnet"}
)

func (v *ValidatingHook) DeploymentCreateHook(ctx context.Context, req admission.Request) admission.Response {
	o := appsv1.Deployment{}
	if err := v.decoder.Decode(req, &o); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	// Get pod template static ips
	staticIPSAnno := o.Spec.Template.GetAnnotations()[util.IpPoolAnnotation]
	klog.V(3).Infof("%s %s@%s, ip_pool: %s", o.Kind, o.GetName(), o.GetNamespace(), staticIPSAnno)
	if staticIPSAnno == "" {
		return ctrlwebhook.Allowed("by pass")
	}
	return v.validateIp(ctx, o.Spec.Template.GetAnnotations(), o.Kind, o.GetName(), o.GetNamespace())
}

func (v *ValidatingHook) StatefulSetCreateHook(ctx context.Context, req admission.Request) admission.Response {
	o := appsv1.StatefulSet{}
	if err := v.decoder.Decode(req, &o); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	// Get pod template static ips
	staticIPSAnno := o.Spec.Template.GetAnnotations()[util.IpPoolAnnotation]
	klog.V(3).Infof("%s %s@%s, ip_pool: %s", o.Kind, o.GetName(), o.GetNamespace(), staticIPSAnno)
	if staticIPSAnno == "" {
		return ctrlwebhook.Allowed("by pass")
	}
	return v.validateIp(ctx, o.Spec.Template.GetAnnotations(), o.Kind, o.GetName(), o.GetNamespace())
}

func (v *ValidatingHook) DaemonSetCreateHook(ctx context.Context, req admission.Request) admission.Response {
	o := appsv1.DaemonSet{}
	if err := v.decoder.Decode(req, &o); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	// Get pod template static ips
	staticIPSAnno := o.Spec.Template.GetAnnotations()[util.IpPoolAnnotation]
	klog.V(3).Infof("%s %s@%s, ip_pool: %s", o.Kind, o.GetName(), o.GetNamespace(), staticIPSAnno)
	if staticIPSAnno == "" {
		return ctrlwebhook.Allowed("by pass")
	}
	return v.validateIp(ctx, o.Spec.Template.GetAnnotations(), o.Kind, o.GetName(), o.GetNamespace())
}

func (v *ValidatingHook) PodCreateHook(ctx context.Context, req admission.Request) admission.Response {
	o := corev1.Pod{}
	if err := v.decoder.Decode(req, &o); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	poolAnno := o.GetAnnotations()[util.IpPoolAnnotation]
	klog.V(3).Infof("%s %s@%s, ip_pool: %s", o.Kind, o.GetName(), o.GetNamespace(), poolAnno)
	if poolAnno != "" {
		return ctrlwebhook.Allowed("by pass")
	}
	staticIP := o.GetAnnotations()[util.IpAddressAnnotation]
	klog.V(3).Infof("%s %s@%s, ip_address: %s", o.Kind, o.GetName(), o.GetNamespace(), staticIP)
	if staticIP == "" {
		return ctrlwebhook.Allowed("by pass")
	}
	return v.validateIp(ctx, o.GetAnnotations(), o.Kind, o.GetName(), o.GetNamespace())
}

func (v *ValidatingHook) SubnetCreateHook(ctx context.Context, req admission.Request) admission.Response {
	o := ovnv1.Subnet{}
	if err := v.decoder.Decode(req, &o); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	if err := util.ValidateSubnet(o); err != nil {
		return ctrlwebhook.Denied(err.Error())
	}

	subnetList := &ovnv1.SubnetList{}
	if err := v.cache.List(ctx, subnetList); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	if err := util.ValidateCidrConflict(o, subnetList.Items); err != nil {
		return ctrlwebhook.Denied(err.Error())
	}

	return ctrlwebhook.Allowed("by pass")
}

func (v *ValidatingHook) validateIp(ctx context.Context, annotations map[string]string, kind, name, namespace string) admission.Response {
	if err := util.ValidatePodNetwork(annotations); err != nil {
		klog.Errorf("validate %s %s/%s failed: %v", kind, namespace, name, err)
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}

	ipList := &ovnv1.IPList{}
	if err := v.cache.List(ctx, ipList); err != nil {
		return ctrlwebhook.Errored(http.StatusBadRequest, err)
	}
	if err := util.ValidateIPConflict(annotations, ipList.Items); err != nil {
		return ctrlwebhook.Denied(err.Error())
	}

	return ctrlwebhook.Allowed("by pass")
}
