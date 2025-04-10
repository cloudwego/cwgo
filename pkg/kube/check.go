package kube

import (
	"errors"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"strings"
)

// check validates the configuration arguments for Kubernetes deployment.
func check(c *config.KubeArgument) error {
	// Validate template type; set default if necessary.
	if c.Template == consts.Kube {
		c.Template = tpl.KubeDir
	}

	// Validate port number is within the valid range (1024-65535).
	if c.Port > 0 && (c.Port < 1024 || c.Port > 65535) {
		return errors.New("port must be between 1024 and 65535")
	}

	// Validate image name is provided.
	if c.Image == "" {
		return errors.New("image must be provided")
	}

	// Validate image pull policy (valid values: Always, Never, IfNotPresent).
	if c.ImagePullPolicy != "" && !isValidImagePullPolicy(c.ImagePullPolicy) {
		return errors.New("invalid imagePullPolicy. Valid values are Always, Never, IfNotPresent")
	}

	// Validate CPU and memory limits and requests; set defaults if not provided.
	if c.LimitCpu <= 0 {
		c.LimitCpu = 1000
	}
	if c.LimitMem <= 0 {
		c.LimitMem = 1024
	}
	if c.RequestCpu <= 0 {
		c.RequestCpu = 500
	}
	if c.RequestMem <= 0 {
		c.RequestMem = 512
	}

	// Validate replica count; set defaults if not provided.
	if c.Replicas <= 0 {
		c.Replicas = 3
	}
	if c.MaxReplicas <= 0 {
		c.MaxReplicas = 10
	}
	if c.MinReplicas <= 0 {
		c.MinReplicas = 3
	}
	// Validate that MinReplicas is less than or equal to MaxReplicas.
	if c.MinReplicas > c.MaxReplicas {
		return errors.New("minReplicas must be less than or equal to maxReplicas")
	}

	// Validate deployment name is provided.
	if c.Name == "" {
		return errors.New("name must be provided")
	}

	// Validate Kubernetes namespace is provided.
	if c.Namespace == "" {
		return errors.New("namespace must be provided")
	}

	// Validate nodePort is within the valid range (30000-32767).
	if c.NodePort > 0 && (c.NodePort < 30000 || c.NodePort > 32767) {
		return errors.New("nodePort must be between 30000 and 32767")
	}

	// Validate targetPort; set it to the value of port if not provided.
	if c.TargetPort <= 0 {
		c.TargetPort = c.Port
	}

	// Validate output YAML file name is provided.
	if c.Output == "" {
		return errors.New("output file (o) must be provided")
	}

	// Validate secret field is provided if specified.
	if len(c.Secret) > 0 && c.Secret == "" {
		return errors.New("secret must be provided if specified")
	}

	// Validate service account field is provided if specified.
	if len(c.ServiceAccount) > 0 && c.ServiceAccount == "" {
		return errors.New("serviceAccount must be provided if specified")
	}

	// Validate template path exists.
	if !strings.HasSuffix(c.Template, consts.SuffixGit) {
		isExist, err := utils.PathExist(c.Template)
		if err != nil {
			return err
		}
		if !isExist {
			return errors.New("DockerFile template not exist")
		}
	}

	return nil
}

// isValidImagePullPolicy checks if the provided image pull policy is valid.
func isValidImagePullPolicy(policy string) bool {
	return policy == "Always" || policy == "Never" || policy == "IfNotPresent"
}
