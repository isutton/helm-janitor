package clean

import (
	"github.com/isutton/helm-janitor/pkg/helmjanitor/env"
)

type Settings struct {
	dryRun               bool
	releaseName          string
	KubeConfig           string
	helmPlugins          string
	helmPluginDir        string
	helmBin              string
	helmRegistryConfig   string
	helmRepositoryCache  string
	helmRepositoryConfig string
	helmNamespace        string
	helmKubeContext      string
}

func (s *Settings) Namespace() string {
	return s.helmNamespace
}

func (s *Settings) ReleaseName() string {
	return s.releaseName
}

func (s *Settings) SetReleaseName(n string) {
	s.releaseName = n
}

func (s *Settings) DryRun() bool {
	return s.dryRun
}

func (s *Settings) SetDryRun(b bool) {
	s.dryRun = b
}

func NewSettings() *Settings {
	return &Settings{
		dryRun:               env.BoolOr("HELM_JANITOR_CLEAN_DRY_RUN", false),
		releaseName:          "",
		KubeConfig:           env.String("KUBECONFIG"),
		helmPlugins:          env.String("HELM_PLUGINS"),
		helmPluginDir:        env.String("HELM_PLUGIN_DIR"),
		helmBin:              env.String("HELM_DEBUG"),
		helmRegistryConfig:   env.String("HELM_REGISTRY_CONFIG"),
		helmRepositoryCache:  env.String("HELM_REPOSITORY_CACHE"),
		helmRepositoryConfig: env.String("HELM_REPOSITORY_CONFIG"),
		helmNamespace:        env.String("HELM_NAMESPACE"),
		helmKubeContext:      env.String("HELM_KUBECONTEXT"),
	}
}

