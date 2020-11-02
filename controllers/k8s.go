package controllers

import (
	"encoding/json"
	"os"
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
	clientset "saas-controller.cloud-native-coding.dev/pkg/generated/clientset/versioned"
	"saas-controller.cloud-native-coding.dev/saasclient"
	"sigs.k8s.io/yaml"
)

const (
	cndeLabel         = "cnde-name"
	appLabel          = "cnde-controller"
	contextVolumeName = "cnde-context"
)

var (
	cndeNamespace = os.Getenv("CNDE_NS")
)

// K8sService main structure
type K8sService struct {
	clientSet        *clientset.Clientset
	k8sClientSet     *kubernetes.Clientset
	metricsClientSet *metricsv.Clientset
	log              *zap.SugaredLogger
}

// NewK8sService creates a new Kubernetes Service instance
func NewK8sService(c *rest.Config, log *zap.SugaredLogger) *K8sService {
	// creates the clientset
	k8sClientSet, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := clientset.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	metricsClientSet, err := metricsv.NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return &K8sService{clientSet, k8sClientSet, metricsClientSet, log}
}

func labelsForDevEnv(name string) map[string]string {
	return map[string]string{cndeLabel: name, "app": appLabel}
}

// CreateBuildFile create config Map with Build-File (Dockerfile)
func (s *K8sService) CreateBuildFile(bf *saasclient.BuildFile) (*corev1.ConfigMap, error) {
	name := strings.ToLower(bf.Name)
	labels := labelsForDevEnv(name)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Data: map[string]string{"Template": bf.Value},
	}
	return s.k8sClientSet.CoreV1().ConfigMaps(cndeNamespace).Create(cm)
}

// DeleteBuildFile deletes one BuildFile ConfigMap
func (s *K8sService) DeleteBuildFile(cm *corev1.ConfigMap) error {
	return s.k8sClientSet.CoreV1().ConfigMaps(cndeNamespace).Delete(cm.Name, &metav1.DeleteOptions{})
}

// DeleteBuildFileByName deletes one BuildFile ConfigMap
func (s *K8sService) DeleteBuildFileByName(name string) error {
	return s.k8sClientSet.CoreV1().ConfigMaps(cndeNamespace).Delete(strings.ToLower(name), &metav1.DeleteOptions{})
}

// GetBuildFiles gets existing K8s resource
func (s *K8sService) GetBuildFiles() (*corev1.ConfigMapList, error) {
	m := map[string]string{"app": appLabel}
	return s.k8sClientSet.CoreV1().ConfigMaps(cndeNamespace).List(metav1.ListOptions{LabelSelector: labels.Set(m).String()})
}

// CreateBuilder creates a builder instance
func (s *K8sService) CreateBuilder(b *saasclient.Builder, buildFileName string) (*v1alpha1.Builder, error) {

	pod := b.Value
	v, err := yaml.YAMLToJSON([]byte(pod))
	if err == nil { // if error != nil, we assume that is was already json
		pod = string(v)
	}

	podSpec := corev1.PodSpec{}
	if err := json.Unmarshal([]byte(pod), &podSpec); err != nil {
		return nil, err
	}

	// check if Volume context is there
	found := false
	volumes := podSpec.Volumes
	for _, vol := range volumes {
		if vol.Name == contextVolumeName {
			found = true
		}
	}

	// if no -> create one for CM defined by buildFileName
	if !found {
		podSpec.Volumes = append(volumes, corev1.Volume{
			Name: contextVolumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: strings.ToLower(buildFileName),
					},
				},
			},
		})
	}

	name := strings.ToLower(b.Name)

	labels := labelsForDevEnv(name)
	builder := &v1alpha1.Builder{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1alpha1.BuilderSpec{
			Template: podSpec,
		},
	}

	return s.clientSet.CndecontrollerV1alpha1().Builders(cndeNamespace).Create(builder)
}

// GetBuilders gets existing K8s resource
func (s *K8sService) GetBuilders() (*v1alpha1.BuilderList, error) {
	return s.clientSet.CndecontrollerV1alpha1().Builders(cndeNamespace).List(metav1.ListOptions{})
}

// DeleteBuilder deletes existing K8s resource
func (s *K8sService) DeleteBuilder(b *v1alpha1.Builder) error {
	return s.clientSet.CndecontrollerV1alpha1().Builders(cndeNamespace).Delete(b.Name, &metav1.DeleteOptions{})
}

// DeleteBuilderByName deletes existing K8s resource
func (s *K8sService) DeleteBuilderByName(name string) error {
	return s.clientSet.CndecontrollerV1alpha1().Builders(cndeNamespace).Delete(strings.ToLower(name), &metav1.DeleteOptions{})
}

// CreateDevEnv creates new Dev Env Instance
func (s *K8sService) CreateDevEnv(de *saasclient.DevEnvUser, keycloakHost string, builder string) (*v1alpha1.DevEnv, error) {
	name := strings.ToLower(de.Name)

	labels := labelsForDevEnv(name)
	devenv := &v1alpha1.DevEnv{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1alpha1.DevEnvSpec{
			// Volume settings
			DockerVolumeSize: de.ContainerVolumeSize,
			HomeVolumeSize:   de.HomeVolumeSize,
			DeleteVolumes:    de.DeleteVolume,

			// Operator environment
			UserEnvDomain: de.UserEnvDomain,
			KeycloakHost:  keycloakHost,
			UserEmail:     de.Email,

			// Definition of Container Images
			//DockerImg: "",		// using default
			DevEnvImg:    de.DevEnvImage,
			ConfigureImg: de.DevEnvImage,
			//KubeConfigImg: "", 	// using defaults
			//OauthProxyImg: "",	/ using defaults

			// DevEnv configuration
			//SSHSecret:       "",	// future
			ClusterRoleName: de.ClusterRoleName,
			RoleName:        de.RoleName,
			BuilderName:     strings.ToLower(builder),
		},
	}
	return s.clientSet.CndecontrollerV1alpha1().DevEnvs().Create(devenv)
}

// DeleteDevEnv deletes existing K8s resource
func (s *K8sService) DeleteDevEnv(de *v1alpha1.DevEnv) error {
	return s.clientSet.CndecontrollerV1alpha1().DevEnvs().Delete(de.Name, &metav1.DeleteOptions{})
}

// GetDevEnvs gets existing K8s resource
func (s *K8sService) GetDevEnvs() (*v1alpha1.DevEnvList, error) {
	return s.clientSet.CndecontrollerV1alpha1().DevEnvs().List(metav1.ListOptions{})
}

// GetPodMetrics returns Pod metrics for a specific namespace
func (s *K8sService) GetPodMetrics(ns string) (*v1beta1.PodMetricsList, error) {
	return s.metricsClientSet.MetricsV1beta1().PodMetricses(ns).List(v1.ListOptions{})
}

// GetNodeMetrics returns Pod metrics for a specific namespace
func (s *K8sService) GetNodeMetrics(nodeName string) (*v1beta1.NodeMetrics, error) {
	return s.metricsClientSet.MetricsV1beta1().NodeMetricses().Get(nodeName, v1.GetOptions{})
}
