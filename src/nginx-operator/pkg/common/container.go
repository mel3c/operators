package common

import (
	daasv1 "nginx-operator/pkg/apis/daas/v1"

	corev1 "k8s.io/api/core/v1"
)

func NewContainer(cr *daasv1.Nginx) []corev1.Container {
	containerPorts := []corev1.ContainerPort{}
	for _, svcPort := range cr.Spec.Ports {
		cport := corev1.ContainerPort{ContainerPort: svcPort.TargetPort.IntVal}
		containerPorts = append(containerPorts, cport)
	}

	return []corev1.Container{
		{
			Name:            cr.Name,
			Image:           cr.Spec.Image,
			Resources:       cr.Spec.Resources,
			Ports:           containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Env:             cr.Spec.Envs,
		},
	}
}
