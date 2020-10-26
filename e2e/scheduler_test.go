package e2e

import (
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

func testScheduler() {
	testNamespacePrefix := "scheduler-test"

	It("should be deployed topolvm-scheduler pod", func() {
		Eventually(func() error {
			result, stderr, err := kubectl("get", "-n=topolvm-system", "pods", "--selector=app.kubernetes.io/name=topolvm-scheduler", "-o=json")
			if err != nil {
				return fmt.Errorf("%v: stdout=%s, stderr=%s", err, result, stderr)
			}

			var podlist corev1.PodList
			err = json.Unmarshal(result, &podlist)
			if err != nil {
				return err
			}

			if len(podlist.Items) == 0 {
				return errors.New("pod is not found")
			}

			for _, pod := range podlist.Items {
				podReady := false
				for _, cond := range pod.Status.Conditions {
					fmt.Fprintln(GinkgoWriter, cond)
					if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
						podReady = true
						break
					}
				}
				if !podReady {
					return errors.New("topolvm-scheduler is not yet ready")
				}
			}

			return nil
		}).Should(Succeed())
	})

	It("should schedule pod if requested capacity is sufficient", func() {
		ns := testNamespacePrefix + randomString(10)
		createNamespace(ns)
		podYml := fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: testhttpd
  namespace: %s
  labels:
    app.kubernetes.io/name: testhttpd
  annotations:
    capacity.topolvm.cybozu.com/ssd: "1073741824"
spec:
  containers:
  - name: ubuntu
    image: quay.io/cybozu/ubuntu:18.04
    command: ["/usr/local/bin/pause"]
    resources:
      requests:
        topolvm.cybozu.com/capacity: 1
      limits:
        topolvm.cybozu.com/capacity: 1
`, ns)
		stdout, stderr, err := kubectlWithInput([]byte(podYml), "apply", "-f", "-")
		Expect(err).ShouldNot(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)

		Eventually(func() error {
			result, stderr, err := kubectl("get", "-n", ns, "pods/testhttpd", "-o=json")
			if err != nil {
				return fmt.Errorf("%v: stdout=%s, stderr=%s", err, result, stderr)
			}

			var pod corev1.Pod
			err = json.Unmarshal(result, &pod)
			if err != nil {
				return err
			}

			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
					return nil
				}
			}

			return fmt.Errorf("testhttpd is not yet ready: %v", pod.Status)
		}).Should(Succeed())
		kubectl("delete", "namespaces", ns)
	})

	It("should not schedule pod if requested capacity is not sufficient", func() {
		ns := testNamespacePrefix + randomString(10)
		createNamespace(ns)
		podYml := fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: testhttpd
  namespace: %s
  labels:
    app.kubernetes.io/name: testhttpd
  annotations:
    capacity.topolvm.cybozu.com/ssd: "21474836480"
spec:
  containers:
  - name: ubuntu
    image: quay.io/cybozu/ubuntu:18.04
    command: ["/usr/local/bin/pause"]
    resources:
      requests:
        topolvm.cybozu.com/capacity: 1
      limits:
        topolvm.cybozu.com/capacity: 1
`, ns)
		stdout, stderr, err := kubectlWithInput([]byte(podYml), "apply", "-f", "-")
		Expect(err).ShouldNot(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)

		Eventually(func() error {
			result, stderr, err := kubectl("get", "-n", ns, "pods/testhttpd", "-o=json")
			if err != nil {
				return fmt.Errorf("%v: stdout=%s, stderr=%s", err, result, stderr)
			}

			var pod corev1.Pod
			err = json.Unmarshal(result, &pod)
			if err != nil {
				return err
			}

			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodScheduled && cond.Status == corev1.ConditionFalse {
					return nil
				}
			}

			return errors.New("testhttpd should not be scheduled")
		}).Should(Succeed())
		kubectl("delete", "namespaces", ns)
	})
}
