package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api/v1"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/yamlutil"
	"os"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//get sa
	sa, err := clientset.CoreV1().ServiceAccounts("openpitrix-system").Get("default", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	if len(sa.Secrets) == 0 {
		logger.Error(nil, "length of sa.secrets is zero")
		os.Exit(-1)
	}
	secretName := sa.Secrets[0].Name

	//get secret name from sa
	st, err := clientset.CoreV1().Secrets("openpitrix-system").Get(secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error(nil, "get secret error %s", err.Error())
		os.Exit(-1)
	}

	//get ca from secret
	crt, ok := st.Data["ca.crt"]
	if !ok {
		logger.Error(nil, "has no key named ca.crt")
		os.Exit(-1)
	}

	token, ok := st.Data["token"]

	if !ok {
		logger.Error(nil, "has no key named token")
		os.Exit(-1)
	}

	kconfig := v1.Config{}
	kconfig.Kind = "Config"
	kconfig.Preferences = v1.Preferences{}

	namedCluster := v1.NamedCluster{}
	namedCluster.Name = "cluster.local"
	cluster := v1.Cluster{}
	cluster.Server = "https://lb.kubesphere.local:6443"
	cluster.CertificateAuthorityData = crt
	namedCluster.Cluster = cluster
	kconfig.Clusters = []v1.NamedCluster{namedCluster}

	namedUser := v1.NamedAuthInfo{}
	namedUser.Name = sa.Name
	namedUser.AuthInfo = v1.AuthInfo{Token: string(token)}
	kconfig.AuthInfos = []v1.NamedAuthInfo{namedUser}

	namedContext := v1.NamedContext{}
	namedContext.Name = "default@cluster.local"
	namedContext.Context = v1.Context{
		Cluster:  "cluster.local",
		AuthInfo: "default",
	}
	kconfig.Contexts = []v1.NamedContext{namedContext}

	kconfig.CurrentContext = "default@cluster.local"

	kk, err := yamlutil.Encode(kconfig)
	if err != nil {
		logger.Error(nil, "encode yaml error: %s", err.Error())
		os.Exit(-1)
	}

	ctxFunc := func() (ctx context.Context) {
		ctx = context.Background()
		ctx = ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
		return
	}

	client, err := runtime.NewRuntimeManagerClient()
	if err != nil {
		logger.Error(nil, "new runtime client error %s", err.Error())
		os.Exit(-1)
	}

	credentialReq := &pb.CreateRuntimeCredentialRequest{
		Name:                     pbutil.ToProtoString("kubeconfig-default"),
		Provider:                 pbutil.ToProtoString(constants.ProviderKubernetes),
		Description:              pbutil.ToProtoString("kubeconfig"),
		RuntimeUrl:               pbutil.ToProtoString("kubesphere"),
		RuntimeCredentialContent: pbutil.ToProtoString(string(kk)),
		RuntimeCredentialId:      pbutil.ToProtoString("default"),
	}
	_, err = client.CreateRuntimeCredential(ctxFunc(), credentialReq)
	if err != nil {
		logger.Error(nil, "create runtime credential error %s", err.Error())
		os.Exit(-1)
	}

	runtimeReq := &pb.CreateRuntimeRequest{
		Name:                pbutil.ToProtoString("default"),
		RuntimeCredentialId: pbutil.ToProtoString("default"),
		Provider:            pbutil.ToProtoString(constants.ProviderKubernetes),
		Zone:                pbutil.ToProtoString("default"),
		RuntimeId:           pbutil.ToProtoString("default"),
	}

	_, err = client.CreateRuntime(ctxFunc(), runtimeReq)
	if err != nil {
		logger.Error(nil, "create runtime error %s", err.Error())
		os.Exit(-1)
	}

}
