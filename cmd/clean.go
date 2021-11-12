 package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/isutton/helm-janitor/cmd/clean"
	"github.com/isutton/helm-janitor/pkg/helmjanitor/flags"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var cleanCmd *cobra.Command

var settings *clean.Settings

func init() {
	settings = clean.NewSettings()
	dryRunFlag := "dry-run"

	cleanCmd = &cobra.Command{
		Use:   "clean [flags] release-name",
		Short: "Removes unused artifacts of previous releases",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			settings.ReleaseName = args[0]

			dryRun, ok, _ := flags.GetBoolPtr(cmd.Flags(), dryRunFlag)
			if ok {
				settings.DryRun = dryRun
			}

			clientset, err := buildClientset(settings)
			if err != nil {
				return err
			}
			
			nsSecretsInterface := clientset.CoreV1().Secrets(settings.Namespace)

			return handleClean(
				cmd.OutOrStderr(), nsSecretsInterface, settings)
		},
	}

	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().Bool(dryRunFlag, false, "do not perform any destructive action")
}

func buildClientset(settings *clean.Settings) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", settings.KubeConfig)
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
	w io.Writer,
	secretInterface typedv1.SecretInterface,
	s *clean.Settings,
) error {
	ctx := context.Background()
	
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
	for _, e := range secretList.Items {
		fmt.Fprintf(w, "- %s\n", e.Name)
		if settings.DryRun {
			fmt.Fprintf(w, "  Would delete secret %q\n", e.Name)
		} else {
			deleteOptions := metav1.DeleteOptions{}
			if settings.DryRun {
				deleteOptions.DryRun = []string{"All"}
			}
			err := secretInterface.Delete(ctx, e.Name, deleteOptions)
			if err != nil {
				fmt.Fprintf(w, "  Error: %s\n", err)
			} else {
				fmt.Fprintf(w, "  Secret %q has been successfully deleted", e.Name)
			}
		}
	}

	return nil
}
