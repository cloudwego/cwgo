package kube

import "github.com/cloudwego/cwgo/config"

// KubeInfo contains all the information required to generate a Kubernetes deployment.
type KubeInfo struct {
	Name            string
	Namespace       string
	Image           string
	Secret          string
	Replicas        int
	Revisions       int
	Port            int
	TargetPort      int
	NodePort        int
	UseNodePort     bool
	RequestCpu      int
	RequestMem      int
	LimitCpu        int
	LimitMem        int
	MinReplicas     int
	MaxReplicas     int
	ServiceAccount  string
	ImagePullPolicy string
}

// FillInfo populates a KubeInfo struct with data from the KubeArgument configuration.
func FillInfo(c *config.KubeArgument) *KubeInfo {
	return &KubeInfo{
		Name:            c.Name,
		Namespace:       c.Namespace,
		Image:           c.Image,
		Secret:          c.Secret,
		Replicas:        c.Replicas,
		Revisions:       c.Revisions,
		Port:            c.Port,
		TargetPort:      c.TargetPort,
		NodePort:        c.NodePort,
		UseNodePort:     c.NodePort > 0, // If NodePort is greater than 0, UseNodePort is true
		RequestCpu:      c.RequestCpu,
		RequestMem:      c.RequestMem,
		LimitCpu:        c.LimitCpu,
		LimitMem:        c.LimitMem,
		MinReplicas:     c.MinReplicas,
		MaxReplicas:     c.MaxReplicas,
		ServiceAccount:  c.ServiceAccount,
		ImagePullPolicy: c.ImagePullPolicy,
	}
}
