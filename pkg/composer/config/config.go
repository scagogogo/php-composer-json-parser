// Package config provides functionality related to PHP Composer configuration
package config

// Config contains configuration information for Composer
type Config struct {
	ProcessTimeout        int                    `json:"process-timeout,omitempty"`
	UseIncludePath        bool                   `json:"use-include-path,omitempty"`
	PreferredInstall      string                 `json:"preferred-install,omitempty"`
	StoreAuths            bool                   `json:"store-auths,omitempty"`
	GithubProtocols       []string               `json:"github-protocols,omitempty"`
	GitlabProtocols       []string               `json:"gitlab-protocols,omitempty"`
	GithubOauth           map[string]string      `json:"github-oauth,omitempty"`
	GitlabOauth           map[string]string      `json:"gitlab-oauth,omitempty"`
	GitlabToken           map[string]string      `json:"gitlab-token,omitempty"`
	Disable               bool                   `json:"disable,omitempty"`
	SecureHttp            bool                   `json:"secure-http,omitempty"`
	Bitbucket             map[string]string      `json:"bitbucket-oauth,omitempty"`
	CaFile                string                 `json:"cafile,omitempty"`
	CaPath                string                 `json:"capath,omitempty"`
	HttpBasic             map[string]interface{} `json:"http-basic,omitempty"`
	Platform              map[string]string      `json:"platform,omitempty"`
	VendorDir             string                 `json:"vendor-dir,omitempty"`
	BinDir                string                 `json:"bin-dir,omitempty"`
	DataDir               string                 `json:"data-dir,omitempty"`
	CacheDir              string                 `json:"cache-dir,omitempty"`
	CacheFilesDir         string                 `json:"cache-files-dir,omitempty"`
	CacheRepoDir          string                 `json:"cache-repo-dir,omitempty"`
	CacheVcsDir           string                 `json:"cache-vcs-dir,omitempty"`
	CacheFileTtl          int                    `json:"cache-file-ttl,omitempty"`
	CacheFilesPerm        interface{}            `json:"cache-files-maxsize,omitempty"`
	CacheFilesMaxsize     string                 `json:"cache-files-maxsize,omitempty"`
	BinCompat             string                 `json:"bin-compat,omitempty"`
	Discard               bool                   `json:"discard-changes,omitempty"`
	AutoloadDumper        string                 `json:"autoloader-suffix,omitempty"`
	OptimizeAutoloader    bool                   `json:"optimize-autoloader,omitempty"`
	PrependAutoloader     bool                   `json:"prepend-autoloader,omitempty"`
	ClassmapAuthoritative bool                   `json:"classmap-authoritative,omitempty"`
	AplusADev             bool                   `json:"apcu-autoloader,omitempty"`
	GithubDomains         []string               `json:"github-domains,omitempty"`
	GitlabDomains         []string               `json:"gitlab-domains,omitempty"`
	UseSslVerify          bool                   `json:"use-github-api,omitempty"`
	UseGithubApi          bool                   `json:"use-github-api,omitempty"`
	NotifyOnInstall       bool                   `json:"notify-on-install,omitempty"`
	DiscardPatches        bool                   `json:"discard-patches,omitempty"`
	ArchiveFormat         string                 `json:"archive-format,omitempty"`
	ArchiveDir            string                 `json:"archive-dir,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ProcessTimeout:     300,
		UseIncludePath:     false,
		PreferredInstall:   "dist",
		StoreAuths:         false,
		GithubProtocols:    []string{"https", "ssh", "git"},
		GitlabProtocols:    []string{"https", "ssh"},
		SecureHttp:         true,
		VendorDir:          "vendor",
		BinDir:             "vendor/bin",
		OptimizeAutoloader: false,
	}
}
