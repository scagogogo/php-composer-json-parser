package config

import (
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// 确保返回的Config不是nil
	if config == nil {
		t.Errorf("DefaultConfig() returned nil")
		return
	}

	// 验证默认值
	if config.ProcessTimeout != 300 {
		t.Errorf("DefaultConfig().ProcessTimeout = %v, want %v", config.ProcessTimeout, 300)
	}

	if config.UseIncludePath != false {
		t.Errorf("DefaultConfig().UseIncludePath = %v, want %v", config.UseIncludePath, false)
	}

	if config.PreferredInstall != "dist" {
		t.Errorf("DefaultConfig().PreferredInstall = %v, want %v", config.PreferredInstall, "dist")
	}

	if config.StoreAuths != false {
		t.Errorf("DefaultConfig().StoreAuths = %v, want %v", config.StoreAuths, false)
	}

	expectedGithubProtocols := []string{"https", "ssh", "git"}
	if !reflect.DeepEqual(config.GithubProtocols, expectedGithubProtocols) {
		t.Errorf("DefaultConfig().GithubProtocols = %v, want %v", config.GithubProtocols, expectedGithubProtocols)
	}

	expectedGitlabProtocols := []string{"https", "ssh"}
	if !reflect.DeepEqual(config.GitlabProtocols, expectedGitlabProtocols) {
		t.Errorf("DefaultConfig().GitlabProtocols = %v, want %v", config.GitlabProtocols, expectedGitlabProtocols)
	}

	if config.SecureHttp != true {
		t.Errorf("DefaultConfig().SecureHttp = %v, want %v", config.SecureHttp, true)
	}

	if config.VendorDir != "vendor" {
		t.Errorf("DefaultConfig().VendorDir = %v, want %v", config.VendorDir, "vendor")
	}

	if config.BinDir != "vendor/bin" {
		t.Errorf("DefaultConfig().BinDir = %v, want %v", config.BinDir, "vendor/bin")
	}

	if config.OptimizeAutoloader != false {
		t.Errorf("DefaultConfig().OptimizeAutoloader = %v, want %v", config.OptimizeAutoloader, false)
	}

	// 验证未设置的字段为零值
	emptyValues := []struct {
		name string
		got  interface{}
	}{
		{"GithubOauth", config.GithubOauth},
		{"GitlabOauth", config.GitlabOauth},
		{"GitlabToken", config.GitlabToken},
		{"Bitbucket", config.Bitbucket},
		{"HttpBasic", config.HttpBasic},
		{"Platform", config.Platform},
		{"CaFile", config.CaFile},
		{"CaPath", config.CaPath},
		{"DataDir", config.DataDir},
		{"CacheDir", config.CacheDir},
		{"CacheFilesDir", config.CacheFilesDir},
		{"CacheRepoDir", config.CacheRepoDir},
		{"CacheVcsDir", config.CacheVcsDir},
		{"BinCompat", config.BinCompat},
	}

	for _, ev := range emptyValues {
		if ev.got != nil && !reflect.ValueOf(ev.got).IsZero() {
			t.Errorf("DefaultConfig().%s should be zero value, got %v", ev.name, ev.got)
		}
	}

	// 验证布尔类型的零值
	boolZeroValues := []struct {
		name string
		got  bool
	}{
		{"Disable", config.Disable},
		{"Discard", config.Discard},
		{"PrependAutoloader", config.PrependAutoloader},
		{"ClassmapAuthoritative", config.ClassmapAuthoritative},
		{"AplusADev", config.AplusADev},
		{"UseSslVerify", config.UseSslVerify},
		{"UseGithubApi", config.UseGithubApi},
		{"NotifyOnInstall", config.NotifyOnInstall},
		{"DiscardPatches", config.DiscardPatches},
	}

	for _, ev := range boolZeroValues {
		if ev.got != false {
			t.Errorf("DefaultConfig().%s should be false, got %v", ev.name, ev.got)
		}
	}
}
