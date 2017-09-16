package main

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
	clientSet       *kubernetes.Clientset
	//restClient      *rest.RESTClient
	//clusterInformer *cache.SharedInformer
	cfg             *OperatorConfig
}

func main() {
  NewHdfsOperatorLib()
}

func NewHdfsOperatorLib() (*HdfsController, error) {
    op, err := NewHdfsOperator(OperatorConfig{})
    if err != nil {
        fmt.Println("Failed to initialize the Hdfs operator: %s", err)
        return nil, err
    }

/*
    fmt.Println("Running operator")
    go o.Run(stop, wg)
*/
    //restartPod(context.Background(),
	//go zookeeper.UpdateZookeeperCluster(ctx, controller.clientSet, newSpec, oldSpec, controller.cfg)

	fmt.Println("Hello, world.\n")
    return op, nil
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
        clientSet:  clientSet,
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
func restartPod(ctx context.Context, clientset *kubernetes.Clientset, spec *ZookeeperClusterSpec, instance int) error {
    opts := metav1.NewDeleteOptions(0)
    opts.PropagationPolicy = &[]metav1.DeletionPropagation{metav1.DeletePropagationForeground}[0]

    podname := spec.GetPodName(instance)

    log.INFO(ctx).Printf("Restarting pod %s", podname)
    startTime, err := getPodStartTimeWithRetries(ctx, clientset, podname, 5)
    if err != nil {
        return errors.Wrap(err, "couldn't get pod start time")
    }

    if err := clientset.CoreV1Client.Pods(NAMESPACE).Delete(podname, opts); err != nil {
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
*/

/*
import "context"
ctx := context.Background()
*/

