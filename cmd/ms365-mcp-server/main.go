package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/softeria/ms-365-mcp-server/internal/auth"
	"github.com/softeria/ms-365-mcp-server/internal/cli"
	"github.com/softeria/ms-365-mcp-server/internal/server"
)

const version = "0.0.0-go"

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("parse args: %v", err)
	}

	if opts.ListPresets {
		for _, p := range cli.AvailablePresets() {
			fmt.Println(p)
		}
		return
	}

	mgr, err := auth.NewManager()
	if err != nil {
		log.Fatalf("init auth manager: %v", err)
	}

	switch {
	case opts.Login:
		acct, err := mgr.LoginDeviceCode()
		if err != nil {
			log.Fatal(err)
		}
		emit(map[string]any{"message": "login completed", "account": acct})
		return
	case opts.VerifyLogin:
		emit(map[string]any{"loggedIn": mgr.HasAccount()})
		return
	case opts.Logout:
		_ = mgr.Logout()
		emit(map[string]any{"message": "logged out"})
		return
	case opts.ListAccounts:
		emit(map[string]any{"accounts": mgr.ListAccounts()})
		return
	case opts.SelectAccount != "":
		if !mgr.SelectAccount(opts.SelectAccount) {
			log.Fatalf("account not found: %s", opts.SelectAccount)
		}
		emit(map[string]any{"message": "account selected", "id": opts.SelectAccount})
		return
	case opts.RemoveAccount != "":
		if !mgr.RemoveAccount(opts.RemoveAccount) {
			log.Fatalf("account not found: %s", opts.RemoveAccount)
		}
		emit(map[string]any{"message": "account removed", "id": opts.RemoveAccount})
		return
	}

	srv := server.New(mgr, opts, version)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

func emit(v any) {
	b, _ := json.Marshal(v)
	fmt.Println(string(b))
}
