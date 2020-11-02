package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	. "saas-controller.cloud-native-coding.dev/controllers"
	"saas-controller.cloud-native-coding.dev/saasclient"
)

var podYAML = `
initContainers:
  - name: context
    image: busybox
    args:
      - /bin/sh
      - -c
      - cp -Lr /context/. /workspace
    volumeMounts:
      - name: context
        mountPath: /context
      - name: build
        mountPath: /workspace
containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:latest
    args:
      [
        "--dockerfile=/workspace/Dockerfile",
        "--context=/workspace",
        "--cache=true",
        "--destination=$IMAGE_TAG",
      ]
    volumeMounts:
      - name: kaniko-secret
        mountPath: /secret
      - name: build
        mountPath: /workspace
    env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/kaniko-secret.json
restartPolicy: Never
volumes:
  - name: kaniko-secret
    secret:
      secretName: c-n-d-e-kaniko-secret
  - name: context
    configMap:
      name: c-n-d-e-dev-env-build-k8s-go
  - name: build
    emptyDir: {}
`

var _ = Describe("Controllers", func() {

	Describe("Testing K8sService", func() {
		var (
			builder *saasclient.Builder
		)

		BeforeEach(func() {
			builder = &saasclient.Builder{
				Name:  "blahblup",
				Value: podYAML,
			}
		})

		Context("checking YAML/JSON to PodSpec mashaling", func() {
			It("should work with this Pod YAML", func() {
				zapLog, _ := zap.NewDevelopment()
				log := zapLog.Sugar()

				k8sService := NewK8sService(cfg, log)

				b, _ := k8sService.CreateBuilder(builder, "")
				//Expect(err).ToNot(HaveOccurred())
				Expect(b).ToNot(BeNil())
			})
		})
	})
})
