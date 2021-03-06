package service

import corev1 "k8s.io/api/core/v1"

// Merge merges `source` into `dest`. Both arguments will remain unchanged
// a new service will be created and returned.
// The "merging" process is arbitrary and it only handle specific attributes
func Merge(dest corev1.Service, source corev1.Service) corev1.Service {
	for k, v := range source.ObjectMeta.Annotations {
		dest.ObjectMeta.Annotations[k] = v
	}

	for k, v := range source.ObjectMeta.Labels {
		dest.ObjectMeta.Labels[k] = v
	}

	var nodePort int32 = 0
	if len(dest.Spec.Ports) > 0 {
		// Save the NodePort for later, in case this ServicePort is changed.
		nodePort = dest.Spec.Ports[0].NodePort
	}

	if len(source.Spec.Ports) > 0 {
		dest.Spec.Ports = source.Spec.Ports

		if nodePort > 0 && source.Spec.Ports[0].NodePort == 0 {
			// There *is* a nodePort defined already, and a new one is not being passed
			dest.Spec.Ports[0].NodePort = nodePort
		}
	}

	dest.Spec.Type = source.Spec.Type
	dest.Spec.LoadBalancerIP = source.Spec.LoadBalancerIP
	dest.Spec.ExternalTrafficPolicy = source.Spec.ExternalTrafficPolicy
	return dest
}
