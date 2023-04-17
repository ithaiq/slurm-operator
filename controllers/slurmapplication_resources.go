package controllers

import (
	slurmoperatorv1beta1 "github.com/aistudio/slurm-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

type SlurmApplicationResources struct {
	app  *slurmoperatorv1beta1.SlurmApplication
	pods map[string]*corev1.Pod
	svcs map[string]*corev1.Service
}

func NewSlurmApplicationResources(app *slurmoperatorv1beta1.SlurmApplication) *SlurmApplicationResources {
	return &SlurmApplicationResources{
		app:  app,
		pods: make(map[string]*corev1.Pod),
		svcs: make(map[string]*corev1.Service),
	}
}

func (s *SlurmApplicationResources) GetPods() map[string]*corev1.Pod {
	return s.pods
}

func (s *SlurmApplicationResources) SetPods(pod *corev1.Pod) {
	if pod == nil {
		return
	}
	s.pods[pod.Name] = pod
}

func (s *SlurmApplicationResources) GetSvcs() map[string]*corev1.Service {
	return s.svcs
}

func (s *SlurmApplicationResources) SetSvcs(svc *corev1.Service) {
	if svc == nil {
		return
	}
	s.svcs[svc.Name] = svc
}
