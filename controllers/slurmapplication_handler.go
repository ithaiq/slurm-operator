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

func (s *SlurmApplicationHandler) GetPodByNamespacedName(ctx context.Context, name, nameSpace string) (*corev1.Pod, error) {
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

func (s *SlurmApplicationHandler) GetSvcByNamespacedName(ctx context.Context, name, nameSpace string) (*corev1.Service, error) {
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

func (s *SlurmApplicationHandler) GenerateSubSlurmPod(app *slurmoperatorv1beta1.SlurmApplication, cs *slurmoperatorv1beta1.CommonSpec, name string) (pod *corev1.Pod, err error) {
	pod = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
			Labels: map[string]string{
				SlurmApplicationLabelKey: name,
			},
		},
	}
	for k, v := range cs.Labels {
		pod.Labels[k] = v
	}
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
				Resources:    cs.Resources,
				VolumeMounts: cs.VolumeMounts,
				Env: []corev1.EnvVar{
					{
						Name:  "SLURM_CPUS_ON_NODE",
						Value: cs.Resources.Limits.Cpu().String(),
					},
					{
						Name:  "SLURM_NODENAME",
						Value: pod.Name,
					},
				},
				Command: []string{"/bin/sh", "-c"},
				Args: []string{
					GetPrePodRunCommand(
						fmt.Sprintf(SedMasterName, GetSlurmResourceName(app.Name, SlurmMaster)),
						fmt.Sprintf(SedNodeName, GetSlurmResourceName(app.Name, SlurmNode)),
						SedNLogFile,
						SedMLogFile,
						SedMLogLevel,
						SedNLogLevel,
						SedSLogFile,
						SedSLogLevel,
						CmdRun,
					),
				},
			},
		},
		Hostname:     name,
		Volumes:      app.Spec.Volumes,
		NodeSelector: cs.NodeSelector,
	}
	if err = controllerutil.SetControllerReference(app, pod, s.scheme); err != nil {
		return nil, fmt.Errorf("setting pod controller reference error : %s", err)
	}
	return pod, nil
}

func (s *SlurmApplicationHandler) GenerateJupyterPod(app *slurmoperatorv1beta1.SlurmApplication, name string) (pod *corev1.Pod, err error) {
	pod = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: app.Namespace,
			Labels: map[string]string{
				SlurmApplicationLabelKey: name,
			},
		},
	}
	for k, v := range app.Spec.Jupyter.Labels {
		pod.Labels[k] = v
	}
	pod.Spec = corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Image: app.Spec.Jupyter.Image,
				Name:  name,
				Ports: []corev1.ContainerPort{
					{
						Name:          "8888p",
						ContainerPort: 8888,
					},
				},
				Command: []string{"/bin/sh", "-c"},
				Args: []string{
					GetPrePodRunCommand(
						fmt.Sprintf(SedMasterName, GetSlurmResourceName(app.Name, SlurmMaster)),
						fmt.Sprintf(SedNodeName, GetSlurmResourceName(app.Name, SlurmNode)),
						SedNLogFile,
						SedMLogFile,
						SedSLogFile,
						SedMLogLevel,
						SedNLogLevel,
						SedSLogLevel,
						CmdRun,
					),
				},
				Resources:    app.Spec.Jupyter.Resources,
				VolumeMounts: app.Spec.Jupyter.VolumeMounts,
			},
		},
		Hostname:     name,
		Volumes:      app.Spec.Volumes,
		NodeSelector: app.Spec.Jupyter.NodeSelector,
	}
	if err = controllerutil.SetControllerReference(app, pod, s.scheme); err != nil {
		return nil, fmt.Errorf("setting pod controller reference error : %s", err)
	}
	return pod, nil
}

func (s *SlurmApplicationHandler) GenerateSubSlurmSvc(app *slurmoperatorv1beta1.SlurmApplication, name string) (svc *corev1.Service, err error) {
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
			SlurmApplicationLabelKey: name,
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
	if err = controllerutil.SetControllerReference(app, svc, s.scheme); err != nil {
		return nil, fmt.Errorf("setting svc controller reference error : %s", err)
	}
	return svc, nil
}

func (s *SlurmApplicationHandler) GenerateJupyterSvc(app *slurmoperatorv1beta1.SlurmApplication, name string) (svc *corev1.Service, err error) {
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
			SlurmApplicationLabelKey: name,
		},
		Ports: []corev1.ServicePort{
			{
				Name:       "8888p",
				Port:       8888,
				TargetPort: intstr.FromInt(8888),
			},
		},
	}
	if err = controllerutil.SetControllerReference(app, svc, s.scheme); err != nil {
		return nil, fmt.Errorf("setting svc controller reference error : %s", err)
	}
	return svc, nil
}
