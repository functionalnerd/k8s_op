package k8s_op

import (
	"fmt"

	"net"
    "os"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"

/*
	"sync"
	"time"

	"k8s.io/client-go/tools/cache"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"

	"git.hubteam.com/HubSpot/HyperDrive/hdfs/pkg/config"
	"git.hubteam.com/HubSpot/HyperDrive/hdfs/pkg/log"
	"git.hubteam.com/HubSpot/HyperDrive/hdfs/pkg/hdfs"
	"github.com/pkg/errors"
*/
)

type OperatorConfig struct {
}

type HdfsController struct {
	ClientSet       *kubernetes.Clientset
	//restClient      *rest.RESTClient
	//clusterInformer *cache.SharedInformer
	cfg             *OperatorConfig
}

func main() {
  NewHdfsOperatorLib()
}

func NewHdfsOperatorLib() (*HdfsController) {
    op, err := NewHdfsOperator(OperatorConfig{})
    if err != nil {
        fmt.Println("Failed to initialize the Hdfs operator: %s", err)
        return nil
    }

/*
    fmt.Println("Running operator")
    go o.Run(stop, wg)
*/
    //restartPod(context.Background(),
	//go zookeeper.UpdateZookeeperCluster(ctx, controller.clientSet, newSpec, oldSpec, controller.cfg)

	fmt.Println("Hello, world.\n")
    return op
}

// NewHdfsOperator is the constructor for Hdfs operators and creates all the initial clients needed
func NewHdfsOperator(cfg OperatorConfig) (*HdfsController, error) {
    // Create an empty context
    //ctx := context.Background()
    clientConfig, err := newConfig()
    if err != nil {
        return nil, errors.Wrap(err, "Failed to retrieve kubernetes client config")
    }

    clientSet, err := kubernetes.NewForConfig(clientConfig)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to create Kubernetes client")
    }

/*
    restClient, err := NewClientForConfig(clientConfig)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to create REST client")
    }
*/

    controller := &HdfsController{
        ClientSet:  clientSet,
        //restClient: restClient,
        cfg:        &cfg,
    }

	return controller, nil
}

// New creates a new config for use in creating clients or otherwise.
func newConfig() (*rest.Config, error) {
    // try in cluster config first, it should fail quickly on lack of env vars
    cfg, err := inClusterConfig()
    if err != nil {
        cfg, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
        if err != nil {
            return nil, errors.Wrap(err, "failed to get InClusterConfig and Config from kube_config")
        }
    }
    return cfg, nil
}

// Thanks Bryan :D
func inClusterConfig() (*rest.Config, error) {
    // Work around https://github.com/kubernetes/kubernetes/issues/40973
    // See https://github.com/coreos/etcd-operator/issues/731#issuecomment-283804819
    if len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0 {
        addrs, err := net.LookupHost("kubernetes.default.svc")
        if err != nil {
            panic(err)
        }
        os.Setenv("KUBERNETES_SERVICE_HOST", addrs[0])
    }
    if len(os.Getenv("KUBERNETES_SERVICE_PORT")) == 0 {
        os.Setenv("KUBERNETES_SERVICE_PORT", "443")
    }
    return rest.InClusterConfig()

    return nil, nil
}

/*
// StatefulSetName returns the expected name of the stateful set for the given spec
func StatefulSetName(spec *ZookeeperClusterSpec, instance int) string {
    return fmt.Sprintf("%s-%d", spec.ClusterName, instance)
}*/

/*
func restartNamenode(clientset *kubernetes.Clientset, instance int) error {
    podname := getPodName(instance)
    restartPod(podname)
}

func getPodName(instance int) string {
	return fmt.Sprintf("%s-%s%d-0", hdfsClusterName, namenodeName, instance)
}

func restartPod(clientset *kubernetes.Clientset, podname string) error {
    opts := metav1.NewDeleteOptions(0)
    opts.PropagationPolicy = &[]metav1.DeletionPropagation{metav1.DeletePropagationForeground}[0]

    fmt.Println("Restarting pod %s", podname)
    startTime, err := getPodStartTimeWithRetries(ctx, clientset, podname, 5)
    if err != nil {
        return errors.Wrap(err, "couldn't get pod start time")
    }

    if err := clientset.CoreV1().Pods(NAMESPACE).Delete(podname, opts); err != nil {
        return fmt.Errorf("failed to delete pod %s", podname)
    }

    // This just makes sure we don't validate the state of the old pod which can happen before the old pod terminates
    if err := waitForNewStartTime(ctx, clientset, podname, startTime); err != nil {
        return fmt.Errorf("pod %s never terminated", podname)
    }

    if err := WaitForHealthyPod(ctx, clientset, spec, instance, podname, successfulHealthChecks); err != nil {
        return errors.Wrap(err, "failed to restart pod")
    }

    log.INFO(ctx).Printf("Pod %s started", podname)
    return nil
}

func getPodStartTime(clientset *kubernetes.Clientset, podname string) (time.Time, error) {
    pod, err := clientset.CoreV1().Pods(NAMESPACE).Get(podname, metav1.GetOptions{})
    if err != nil {
        return time.Time{}, errors.Wrap(err, "failed to get pod")
    }

    if pod.Status.StartTime == nil {
        return time.Time{}, errors.New("got nil start time for pod")
    }

    return pod.Status.StartTime.Time, nil
}

func getPodStartTimeWithRetries(ctx context.Context, clientset *kubernetes.Clientset, podname string, retries int) (time.Time, error) {
    for i := 0; i < retries; i++ {
        startTime, err := getPodStartTime(clientset, podname)
        if err != nil {
            log.WARN(ctx).Printf("Failed to get start time for pod %s: %s", podname, err)
            time.Sleep(5 * time.Second)
            continue
        }

        return startTime, nil
    }

    return time.Time{}, errors.New("failed to get pod start time within the alloted retries")
}

func waitForNewStartTime(ctx context.Context, clientset *kubernetes.Clientset, podname string, currentStartTime time.Time) error {
    timeout := time.Now().Add(60 * time.Second)
    sleepDuration := 5 * time.Second
    for {
        if time.Now().After(timeout) {
            return errors.New("failed to restart within timeout")
        }

        newStartTime, err := getPodStartTimeWithRetries(ctx, clientset, podname, 5)
        if err != nil {
            log.WARN(ctx).Printf("non-fatal error getting pod start time, will retry if we're still within the time out: %s", err)
            time.Sleep(sleepDuration)
            continue
        }

        if newStartTime.After(currentStartTime) {
            break
        }

        time.Sleep(sleepDuration)
    }

    return nil
}
*/

/*
func deletePod() {
    podname := "jumbo-test-test-busybox"
    opts := metav1.NewDeleteOptions(0)
    opts.PropagationPolicy = &[]metav1.DeletionPropagation{metav1.DeletePropagationForeground}[0]

    if err := op.ClientSet.CoreV1().Pods("hdfs").Delete(podname, opts); err != nil {
        fmt.Println("failed to delete pod %s", podname)
    }
}
*/

