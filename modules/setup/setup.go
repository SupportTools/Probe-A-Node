package setup

import (
	"context"

	"github.com/supporttools/Probe-A-Node/modules/cli"
	"github.com/supporttools/Probe-A-Node/modules/logging"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var log = logging.SetupLogging()

func CreateOrUpdateDaemonSet(clientset *kubernetes.Clientset) (string, error) {
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "overlay-test-pod",
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "overlay-test-pod",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "overlay-test-pod",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "overlay-test-pod",
							Image: "rancherlabs/swiss-army-knife:latest",
							Command: []string{
								"sh",
								"-c",
								"tail -f /dev/null",
							},
						},
					},
				},
			},
		},
	}

	namespace := cli.Settings().Namespace
	_, err := clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), daemonset, metav1.CreateOptions{})
	if err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return "", err
		}
		// DaemonSet already exists, let's update it
		log.Debugf("Updating daemonset %q\n", daemonset.Name)
		result, updateErr := clientset.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonset, metav1.UpdateOptions{})
		if updateErr != nil {
			return "", updateErr
		}
		log.Debugln("Updated daemonset %q.\n", result.GetObjectMeta().GetName())
		return result.GetObjectMeta().GetName(), nil
	}
	log.Debugln("Created daemonset %q.\n", daemonset.Name)
	return daemonset.Name, nil
}
