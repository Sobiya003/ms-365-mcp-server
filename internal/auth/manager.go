package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Selected bool   `json:"selected"`
}

type Manager struct {
	accounts []Account
	path     string
}

func NewManager() (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".ms365-mcp", "accounts.json")
	m := &Manager{path: path}
	_ = m.load()
	return m, nil
}

func (m *Manager) LoginDeviceCode() (Account, error) {
	id := fmt.Sprintf("acct-%d", time.Now().Unix())
	acct := Account{ID: id, Username: "user@example.com", Name: "Microsoft User", Selected: true}
	for i := range m.accounts {
		m.accounts[i].Selected = false
	}
	m.accounts = append(m.accounts, acct)
	return acct, m.save()
}

func (m *Manager) HasAccount() bool { return len(m.accounts) > 0 }

func (m *Manager) ListAccounts() []Account { return append([]Account(nil), m.accounts...) }

func (m *Manager) SelectAccount(id string) bool {
	found := false
	for i := range m.accounts {
		m.accounts[i].Selected = m.accounts[i].ID == id
		if m.accounts[i].Selected {
			found = true
		}
	}
	if found {
		_ = m.save()
	}
	return found
}

func (m *Manager) RemoveAccount(id string) bool {
	idx := -1
	for i := range m.accounts {
		if m.accounts[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	m.accounts = append(m.accounts[:idx], m.accounts[idx+1:]...)
	if len(m.accounts) > 0 {
		m.accounts[0].Selected = true
	}
	_ = m.save()
	return true
}

func (m *Manager) Logout() error {
	m.accounts = nil
	return m.save()
}

func (m *Manager) load() error {
	b, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &m.accounts)
}

func (m *Manager) save() error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m.accounts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.path, b, 0o600)
}
