package resource

import (
	daasv1 "redis-operator/pkg/apis/daas/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewStatefulsetForCR(cr *daasv1.Redis) *appsv1.Deployment {
	labels := map[string]string{"name": cr.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: cr.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: NewContainer(cr),
				},
			},
			Selector: selector,
		},
	}
}
