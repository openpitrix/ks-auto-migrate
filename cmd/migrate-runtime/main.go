package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"os"
)

func main() {
	// check whether default runtime exists
	client, err := runtime.NewRuntimeManagerClient()
	if err != nil {
		logger.Error(nil, "new runtime client error: %s", err.Error())
		os.Exit(-1)
	}

	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	describeReq := &pb.DescribeRuntimeCredentialsRequest{
		RuntimeCredentialId: []string{"default"},
	}
	describeResp, err := client.DescribeRuntimeCredentials(ctxFunc(), describeReq)
	if err != nil {
		logger.Error(nil, "describe runtime credential error: %s", err.Error())
		os.Exit(-1)
	}
	if describeResp.TotalCount <= 0 || describeResp.RuntimeCredentialSet[0].RuntimeCredentialId.GetValue() != "default" {
		logger.Error(nil, "cannot find runtime: default")
		os.Exit(-1)
	}

	// if default exists,update cluster runtimeId
	clusterClient, err := cluster.NewClient()
	if err != nil {
		logger.Error(nil, "new cluster client error: %s", err.Error())
		os.Exit(-1)
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Error(nil, "cluster config client error: %s", err.Error())
		os.Exit(-1)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(nil, "client set error: %s", err.Error())
		os.Exit(-1)
	}

	namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		logger.Error(nil, "list namespaces error: %s, %s", err.Error())
		os.Exit(-1)
	}

	var runtimeIds []string
	for _, ns := range namespaces.Items {
		runtimeId, ok := ns.Annotations["openpitrix_runtime"]
		if !ok {
			continue
		}
		runtimeIds = append(runtimeIds, runtimeId)
	}

	applications, err := clusterClient.DescribeClusters(ctxFunc(), &pb.DescribeClustersRequest{
		RuntimeId: runtimeIds,
	})

	for _, application := range applications.ClusterSet {
		application.RuntimeId = &wrappers.StringValue{Value: "default"}
		_, err = clusterClient.ModifyCluster(ctxFunc(), &pb.ModifyClusterRequest{
			Cluster: application,
		})
		if err != nil {
			logger.Error(nil, "Failed to update application: %s", application.Name)
			os.Exit(-1)
		}
	}

	// delete namespace runtime
	for _, ns := range namespaces.Items {
		delete(ns.Annotations, "openpitrix_runtime")
		_, err := clientSet.CoreV1().Namespaces().Update(&ns)
		if err != nil {
			logger.Error(nil, "Failed to update namespace: %s, %s", ns.Name, err.Error())
		}
	}
}
