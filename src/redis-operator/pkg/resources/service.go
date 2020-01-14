package resource

import (
	"fmt"
	daasv1 "redis-operator/pkg/apis/daas/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewServiceForCR(cr *daasv1.Redis) *corev1.Service {
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
			Ports:    getServicePorts(cr),
			Selector: selector,
		},
	}
}

func getServicePorts(cr *daasv1.Redis) []corev1.ServicePort {
	return []corev1.ServicePort{corev1.ServicePort{
		Name:     fmt.Sprintf("tcp-%s-master-service", cr.Name),
		Port:     6379,
		Protocol: corev1.ProtocolTCP,
	}}
}
