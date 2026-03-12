package cli

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type Options struct {
	Verbose                   bool
	Login                     bool
	Logout                    bool
	VerifyLogin               bool
	ListAccounts              bool
	SelectAccount             string
	RemoveAccount             string
	ReadOnly                  bool
	HTTP                      string
	EnableAuthTools           bool
	EnabledTools              string
	Preset                    string
	ListPresets               bool
	OrgMode                   bool
	Toon                      bool
	Discovery                 bool
	Cloud                     string
	EnableDynamicRegistration bool
}

func Parse(args []string) (Options, error) {
	fs := flag.NewFlagSet("ms-365-mcp-server", flag.ContinueOnError)
	var o Options

	fs.BoolVar(&o.Verbose, "v", false, "Enable verbose logging")
	fs.BoolVar(&o.Login, "login", false, "Login using a simulated device flow")
	fs.BoolVar(&o.Logout, "logout", false, "Logout and clear cached account")
	fs.BoolVar(&o.VerifyLogin, "verify-login", false, "Verify login")
	fs.BoolVar(&o.ListAccounts, "list-accounts", false, "List cached accounts")
	fs.StringVar(&o.SelectAccount, "select-account", "", "Select account ID")
	fs.StringVar(&o.RemoveAccount, "remove-account", "", "Remove account ID")
	fs.BoolVar(&o.ReadOnly, "read-only", false, "Disable write operations")
	fs.StringVar(&o.HTTP, "http", "", "Start HTTP mode on [host:]port")
	fs.BoolVar(&o.EnableAuthTools, "enable-auth-tools", false, "Enable auth tools in HTTP mode")
	fs.StringVar(&o.EnabledTools, "enabled-tools", "", "Tool regex filter")
	fs.StringVar(&o.Preset, "preset", "", "Tool preset(s)")
	fs.BoolVar(&o.ListPresets, "list-presets", false, "List presets and exit")
	fs.BoolVar(&o.OrgMode, "org-mode", false, "Enable organization/work mode")
	fs.BoolVar(&o.Toon, "toon", false, "Enable TOON output")
	fs.BoolVar(&o.Discovery, "discovery", false, "Enable discovery mode")
	fs.StringVar(&o.Cloud, "cloud", "global", "Cloud type (global|china)")
	fs.BoolVar(&o.EnableDynamicRegistration, "enable-dynamic-registration", false, "Enable DCR endpoint")

	if err := fs.Parse(args); err != nil {
		return o, err
	}
	if o.Cloud != "global" && o.Cloud != "china" {
		return o, errors.New("cloud must be one of: global, china")
	}
	if o.Preset != "" {
		for _, p := range strings.Split(o.Preset, ",") {
			if !isPreset(strings.TrimSpace(p)) {
				return o, fmt.Errorf("invalid preset: %s", p)
			}
		}
	}

	return o, nil
}

func AvailablePresets() []string {
	return []string{"mail", "calendar", "files", "personal", "work", "excel", "contacts", "tasks", "onenote", "search", "users", "all"}
}

func isPreset(s string) bool {
	for _, p := range AvailablePresets() {
		if p == s {
			return true
		}
	}
	return false
}
