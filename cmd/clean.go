package cmd

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/isutton/helm-janitor/pkg/helm/plugin"
	"github.com/isutton/helm-janitor/pkg/helmjanitor/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var cleanCmd *cobra.Command

var settings *plugin.Settings

func addKubeconfigFlag(flags *pflag.FlagSet) *string {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flags.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flags.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return kubeconfig
}

func init() {
	settings = plugin.NewSettings()
	dryRunFlag := "dry-run"
	var kubeconfig *string

	cleanCmd = &cobra.Command{
		Use:   "clean [flags] release-name",
		Short: "remove unused artifacts of previous failed releases",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			settings.ReleaseName = args[0]
			// kubeconfig should have a valid value here, thus no check
			settings.Kubeconfig = *kubeconfig

			dryRun, ok, _ := flags.GetBoolPtr(cmd.Flags(), dryRunFlag)
			if ok {
				settings.DryRun = dryRun
			}

			clientset, err := buildClientset(settings)
			if err != nil {
				return err
			}

			nsSecretsInterface := clientset.CoreV1().Secrets(settings.Namespace)

			ctx := context.Background()

			return handleClean(
				ctx, cmd.OutOrStderr(), nsSecretsInterface, settings)
		},
	}
	kubeconfig = addKubeconfigFlag(cleanCmd.Flags())
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().Bool(dryRunFlag, false, "do not perform any destructive action")
}

func buildClientset(settings *plugin.Settings) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", settings.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("building config from flags: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("creating clientset: %w", err)
	}
	return clientset, err
}

func handleClean(
	ctx context.Context,
	w io.Writer,
	secretInterface typedv1.SecretInterface,
	s *plugin.Settings,
) error {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"owner": "helm",
			"name":  s.ReleaseName,
		},
	}

	labelMap, err := metav1.LabelSelectorAsMap(labelSelector)
	if err != nil {
		return fmt.Errorf("converting label selector to map: %w", err)
	}

	listOptions := metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
		Limit:         100,
	}
	secretList, err := secretInterface.List(ctx, listOptions)
	if err != nil {
		return fmt.Errorf("listing secrets for release %q: %w", s.ReleaseName, err)
	}

	fmt.Fprintf(w, "Found %d secret(s) for release %q\n", len(secretList.Items), s.ReleaseName)
	for _, secret := range secretList.Items {
		fmt.Fprintf(w, "- %s\n", secret.Name)

		if !shouldDeleteSecret(secret) {
			fmt.Fprintf(w, "  Should not delete secret %q\n", secret.Name)
			continue
		}

		if settings.DryRun {
			fmt.Fprintf(w, "  Would delete secret %q\n", secret.Name)
		} else {
			deleteOptions := metav1.DeleteOptions{}
			if settings.DryRun {
				deleteOptions.DryRun = []string{"All"}
			}
			err := secretInterface.Delete(ctx, secret.Name, deleteOptions)
			if err != nil {
				fmt.Fprintf(w, "  Error: %s\n", err)
			} else {
				fmt.Fprintf(w, "  Secret %q has been successfully deleted", secret.Name)
			}
		}
	}

	return nil
}

func shouldDeleteSecret(secret corev1.Secret) bool {
	return true
}
