package resource

import (
	daasv1 "nginx-operator/pkg/apis/daas/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewSlaveServiceForCR(cr *daasv1.Nginx) *corev1.Service {
	selector := map[string]string{
		"name": cr.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Ports:    cr.Spec.Ports,
			Selector: selector,
		},
	}
}
