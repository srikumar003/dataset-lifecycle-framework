package e2e

import (
	goctx "context"
	"testing"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/IBM/dataset-lifecycle-framework/dataset-operator/pkg/apis"
	datasetv1alpha1 "github.com/IBM/dataset-lifecycle-framework/dataset-operator/pkg/apis/com/v1alpha1/"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestDataset(t *testing.T) {
	dataset := &datasetv1alpha1.Dataset{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, dataset)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to the framework: %v", err)
	}

	t.Run("dataset-tests", func(t *testing.T) {
		t.Run("CreateDataset", DatasetCreate)
	})
}

func DatasetCreate(t *testing.T) {

	t.Parallel()

	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&framework.CleanupOptions{testContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialise cluster resources: %v", err)
	}
	t.Log("Initialised cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}

	f := framework.Global

	err := e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "dataset-operator", 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	if err := datasetCreateTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}

}

func datasetCreateTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	
	namespace, err := ctx.Namespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}

	exampleDataset := &datasetv1alphav1.Dataset{
			ObjectMeta: metav1.ObjectMeta{
				Name: "example-dataset",
				Namespace: namespace,
			}
			Spec: datasetv1alpha1.DatasetSpec {
				Local = map[string]string{
					"type":"COS",
					"accessKeyID": "cos-key",
					"secretAccessKey": "sekrit-key",
					"endpoint": "http://cos.endpoint",
					"bucket": "test-bucket"
				}
			}
	}

	err := f.Client.Create(goctx.TODO(), exampleDataset, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	err := f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-dataset", Namespace: namespace}, exampleDataset)
	if err != nil {
		return err
	}

	return nil
}
