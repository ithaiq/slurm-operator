package controllers

import (
	"context"
	"fmt"
	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type SlurmApplicationHandler struct {
	scheme *runtime.Scheme
	client client.Client
}

func NewSlurmApplicationHandler(scheme *runtime.Scheme, client client.Client) *SlurmApplicationHandler {
	return &SlurmApplicationHandler{
		scheme: scheme,
		client: client,
	}
}

func (s *SlurmApplicationHandler) GenerateSlurmResource(app *slurmoperatorv1beta1.SlurmApplication, cs *slurmoperatorv1beta1.CommonSpec, pod *corev1.Pod, svc *corev1.Service, name string) (*corev1.Pod, *corev1.Service, error) {
	pod, err := s.generateSlurmPod(app, cs, pod, name)
	if err != nil {
		return nil, nil, err
	}
	svc, err = s.generateSlurmSvc(app, svc, name)
	if err != nil {
		return nil, nil, err
	}
	return pod, svc, nil
}

func (s *SlurmApplicationHandler) GenerateJupyterResource(app *slurmoperatorv1beta1.SlurmApplication, pod *corev1.Pod, svc *corev1.Service, name string) (*corev1.Pod, *corev1.Service, error) {
	pod, err := s.generateJupyterPod(app, pod, name)
	if err != nil {
		return nil, nil, err
	}
	svc, err = s.generateJupyterSvc(app, svc, name)
	if err != nil {
		return nil, nil, err
	}
	return pod, svc, nil
}

func (s *SlurmApplicationHandler) GetSlurmApplicationResource(ctx context.Context, name, nameSpace string) (*corev1.Pod, *corev1.Service, error) {
	pod, err := s.getPodByNamespacedName(ctx, name, nameSpace)
	if err != nil {
		return nil, nil, err
	}
	svc, err := s.getSvcByNamespacedName(ctx, name, nameSpace)
	if err != nil {
		return nil, nil, err
	}
	return pod, svc, nil
}

func (s *SlurmApplicationHandler) generateSlurmPod(app *slurmoperatorv1beta1.SlurmApplication, cs *slurmoperatorv1beta1.CommonSpec, srcPod *corev1.Pod, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
		},
	}
	if srcPod != nil {
		pod = srcPod.DeepCopy()
	}
	pod.Labels = map[string]string{
		SlurmApplicationLabelKey: pod.Name,
	}
	for k, v := range cs.Labels {
		pod.Labels[k] = v
	}
	if len(pod.Spec.Containers) == 0 {
		pod.Spec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: cs.Image,
					Name:  pod.Name,
					Ports: []corev1.ContainerPort{
						{
							Name:          "6817p",
							ContainerPort: 6817,
						},
						{
							Name:          "6818p",
							ContainerPort: 6818,
						},
						{
							Name:          "6819p",
							ContainerPort: 6819,
						},
					},
					Resources: cs.Resources,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      pod.Name,
							MountPath: "/home/admin",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "SLURM_CPUS_ON_NODE",
							Value: "4",
						},
						{
							Name:  "SLURM_NODENAME",
							Value: pod.Name,
						},
					},
				},
			},
			Hostname:           pod.Name,
			ServiceAccountName: "default",
			Volumes: []corev1.Volume{
				{
					Name: pod.Name,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{Path: "/root/test"},
					},
				},
			},
			NodeSelector: cs.NodeSelector,
		}
	} else {
		pod.Spec.Containers[0].Image = cs.Image
	}
	if err := controllerutil.SetControllerReference(app, pod, s.scheme); err != nil {
		return nil, fmt.Errorf("setting pod controller reference error : %s", err)
	}
	return pod, nil
}

func (s *SlurmApplicationHandler) generateJupyterPod(app *slurmoperatorv1beta1.SlurmApplication, srcPod *corev1.Pod, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
		},
	}
	if srcPod != nil {
		pod = srcPod.DeepCopy()
	}
	pod.Labels = map[string]string{
		SlurmApplicationLabelKey: pod.Name,
	}
	for k, v := range app.Spec.Jupyter.Labels {
		pod.Labels[k] = v
	}
	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers[0].Image = app.Spec.Jupyter.Image
	} else {
		pod.Spec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: app.Spec.Jupyter.Image,
					Name:  pod.Name,
					Ports: []corev1.ContainerPort{
						{
							Name:          "8888p",
							ContainerPort: 8888,
						},
					},
					Resources: app.Spec.Jupyter.Resources,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      pod.Name,
							MountPath: "/home/admin",
						},
					},
				},
			},
			Hostname: pod.Name,
			Volumes: []corev1.Volume{
				{
					Name: pod.Name,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{Path: "/root/test"},
					},
				},
			},
			NodeSelector: app.Spec.Jupyter.NodeSelector,
		}
	}
	if err := controllerutil.SetControllerReference(app, pod, s.scheme); err != nil {
		return nil, fmt.Errorf("setting pod controller reference error : %s", err)
	}
	return pod, nil
}

func (s *SlurmApplicationHandler) generateSlurmSvc(app *slurmoperatorv1beta1.SlurmApplication, svc *corev1.Service, name string) (*corev1.Service, error) {
	if svc != nil {
		return svc, nil
	}
	svc = &corev1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
		},
	}
	svc.Namespace = app.Namespace
	svc.Name = name
	svc.Labels = app.Labels
	svc.Spec = corev1.ServiceSpec{
		Type: corev1.ServiceTypeClusterIP,
		Selector: map[string]string{
			SlurmApplicationLabelKey: svc.Name,
		},
		Ports: []corev1.ServicePort{
			{
				Name:       "6817p",
				Port:       6817,
				TargetPort: intstr.FromInt(6817),
			},
			{
				Name:       "6818p",
				Port:       6818,
				TargetPort: intstr.FromInt(6818),
			},
			{
				Name:       "6819p",
				Port:       6819,
				TargetPort: intstr.FromInt(6819),
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, svc, s.scheme); err != nil {
		return nil, fmt.Errorf("setting svc controller reference error : %s", err)
	}
	return svc, nil
}

func (s *SlurmApplicationHandler) generateJupyterSvc(app *slurmoperatorv1beta1.SlurmApplication, svc *corev1.Service, name string) (*corev1.Service, error) {
	if svc != nil {
		return svc, nil
	}
	svc = &corev1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
		},
	}
	svc.Namespace = app.Namespace
	svc.Name = name
	svc.Labels = app.Labels
	svc.Spec = corev1.ServiceSpec{
		Type: corev1.ServiceTypeNodePort,
		Selector: map[string]string{
			SlurmApplicationLabelKey: svc.Name,
		},
		Ports: []corev1.ServicePort{
			{
				Name:       "8888p",
				Port:       8888,
				NodePort:   28888,
				TargetPort: intstr.FromInt(8888),
			},
		},
	}
	if err := controllerutil.SetControllerReference(app, svc, s.scheme); err != nil {
		return nil, fmt.Errorf("setting svc controller reference error : %s", err)
	}
	return svc, nil
}

func (s *SlurmApplicationHandler) getPodByNamespacedName(ctx context.Context, name, nameSpace string) (*corev1.Pod, error) {
	key := client.ObjectKey{
		Name:      name,
		Namespace: nameSpace,
	}
	pod := &corev1.Pod{}
	err := s.client.Get(ctx, key, pod)
	if err != nil {
		pod = nil
	}
	return pod, client.IgnoreNotFound(err)
}

func (s *SlurmApplicationHandler) getSvcByNamespacedName(ctx context.Context, name, nameSpace string) (*corev1.Service, error) {
	key := client.ObjectKey{
		Name:      name,
		Namespace: nameSpace,
	}
	svc := &corev1.Service{}
	err := s.client.Get(ctx, key, svc)
	if err != nil {
		svc = nil
	}
	return svc, client.IgnoreNotFound(err)
}
