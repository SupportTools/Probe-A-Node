package checks

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/miekg/dns"
	"github.com/supporttools/Probe-A-Node/modules/cli"
	"github.com/supporttools/Probe-A-Node/modules/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var log = logging.SetupLogging()

func CheckKubeApi(clientset *kubernetes.Clientset, nodeName string) (string, error) {
	node, err := clientset.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return node.Status.NodeInfo.KubeletVersion, nil
}

func CheckDns(dnsEndpoint string, dnsServer string) error {
	conf := &dns.ClientConfig{Servers: []string{dnsServer}, Port: "53"}
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dnsEndpoint), dns.TypeA)
	r, _, err := c.Exchange(m, net.JoinHostPort(conf.Servers[0], conf.Port))
	if r == nil {
		return err
	}
	if r.Rcode != dns.RcodeSuccess {
		return err
	}
	return nil
}

func CheckOverlayNetwork(clientset *kubernetes.Clientset, nodeName string) (string, int, error) {
	daemonsetName := "Probe-A-Node"
	namespace := cli.Settings().Namespace
	daemonset, err := clientset.AppsV1().DaemonSets(namespace).Get(context.Background(), daemonsetName, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("error getting DaemonSet %s: %v", daemonsetName, err)
		return msg, -1, err
	}

	if daemonset.Status.NumberReady == 0 {
		msg := fmt.Sprintf("no Pods are ready for DaemonSet %s", daemonsetName)
		return msg, -1, err
	}

	podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: metav1.FormatLabelSelector(daemonset.Spec.Selector)})
	if err != nil {
		log.Debug("error getting Pods for DaemonSet %s: %v", daemonsetName, err)
		msg := fmt.Sprintf("error getting Pods for DaemonSet %s: %v", daemonsetName, err)
		return msg, -1, err
	}

	healthCheckPort := cli.Settings().HealthCheckPort
	unhealthyPods := 0
	for _, pod := range podList.Items {
		url := fmt.Sprintf("http://%s:%s/healthz", pod.Status.PodIP, healthCheckPort)

		resp, err := http.Get(url)
		if err != nil {
			log.Warnf("Pod %s health check failed with error %v", pod.Name, err)
			unhealthyPods++
		}

		if resp.StatusCode != http.StatusOK {
			log.Warnf("Pod %s health check failed with status code %d", pod.Name, resp.StatusCode)
			unhealthyPods++
		}

		log.Infof("Pod %s is healthy", pod.Name)
		resp.Body.Close()
	}

	if unhealthyPods > 0 {
		msg := fmt.Sprintf("unhealthy pods: %d", unhealthyPods)
		return msg, unhealthyPods, nil
	}
	msg := "all pods are healthy"
	return msg, unhealthyPods, nil
}
