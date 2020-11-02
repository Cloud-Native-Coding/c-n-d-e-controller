package controllers_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
	clientset "saas-controller.cloud-native-coding.dev/pkg/generated/clientset/versioned"
	"saas-controller.cloud-native-coding.dev/pkg/generated/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	testEnv       *envtest.Environment
	cfg           *rest.Config
	k8sClientSet  *kubernetes.Clientset
	cndeClientSet *clientset.Clientset
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		//CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
		CRDInstallOptions: envtest.CRDInstallOptions{
			ErrorIfPathMissing: true,
			Paths:              []string{filepath.Join("..", "config", "crd", "bases")},
		},
		AttachControlPlaneOutput: true,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClientSet, err = kubernetes.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClientSet).ToNot(BeNil())

	cndeClientSet, err = clientset.NewForConfig(cfg)
	Expect(err).ToNot(HaveOccurred())
	Expect(cndeClientSet).ToNot(BeNil())

	close(done)

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
