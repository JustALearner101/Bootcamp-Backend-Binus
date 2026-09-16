package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	accountsDir  = "accounts"
	depositsDir  = "deposits"
	historyDir   = "history"
	interestRate = 0.01 // 1% per minute
)

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func getAccountPath(name string) string {
	return filepath.Join(accountsDir, name+".txt")
}

func getDepositPath(name string) string {
	return filepath.Join(depositsDir, name+".txt")
}

func getHistoryPath(name string) string {
	return filepath.Join(historyDir, name+".txt")
}

//CreateAccount ya create account lah ya
func CreateAccount(name string, amount float64) error {
	if err := ensureDir(accountsDir); err != nil {
		return fmt.Errorf("failed to create accounts directory: %w", err)
	}
	if err := ensureDir(historyDir); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	accountPath := getAccountPath(name)
	if _, err := os.Stat(accountPath); err == nil {
		return fmt.Errorf("account %s already exists", name)
	}

	balance := fmt.Sprintf("%.2f", amount)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	content := fmt.Sprintf("Balance: %s\nCreated: %s\nTransactions:\n", balance, timestamp)
	if err := os.WriteFile(accountPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write account file: %w", err)
	}

	historyEntry := fmt.Sprintf("[%s] Account created with balance %.2f\n", timestamp, amount)
	return appendToFile(getHistoryPath(name), historyEntry)
}

// Transfer antar akun
func Transfer(from, to string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	fromBalance, err := GetBalance(from)
	if err != nil {
		return fmt.Errorf("sender account error: %w", err)
	}
	if fromBalance < amount {
		return fmt.Errorf("insufficient funds: %s has %.2f, needs %.2f", from, fromBalance, amount)
	}

	toBalance, err := GetBalance(to)
	if err != nil {
		return fmt.Errorf("receiver account error: %w", err)
	}

	newFromBalance := fromBalance - amount
	newToBalance := toBalance + amount
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err := updateAccountBalance(from, newFromBalance, timestamp, fmt.Sprintf("TRANSFER OUT: -%.2f to %s", amount, to)); err != nil {
		return fmt.Errorf("failed to update sender: %w", err)
	}

	if err := updateAccountBalance(to, newToBalance, timestamp, fmt.Sprintf("TRANSFER IN: +%.2f from %s", amount, from)); err != nil {
		return fmt.Errorf("failed to update receiver: %w", err)
	}

	return nil
}

// GetBalance baca balance dari file account.txt
func GetBalance(name string) (float64, error) {
	data, err := os.ReadFile(getAccountPath(name))
	if err != nil {
		return 0, fmt.Errorf("account %s not found", name)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Balance:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				balance, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err != nil {
					return 0, fmt.Errorf("invalid balance format: %w", err)
				}
				return balance, nil
			}
		}
	}
	return 0, fmt.Errorf("balance not found in account file")
}

//buat update account balance
func updateAccountBalance(name string, newBalance float64, timestamp, transaction string) error {
	data, err := os.ReadFile(getAccountPath(name))
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "Balance:") {
			newLines = append(newLines, fmt.Sprintf("Balance: %.2f", newBalance))
		} else if line == "Transactions:" {
			newLines = append(newLines, line)
			newLines = append(newLines, fmt.Sprintf("  [%s] %s", timestamp, transaction))
		} else {
			newLines = append(newLines, line)
		}
	}

	content := strings.Join(newLines, "\n")
	if err := os.WriteFile(getAccountPath(name), []byte(content), 0644); err != nil {
		return err
	}

	historyEntry := fmt.Sprintf("[%s] %s\n", timestamp, transaction)
	return appendToFile(getHistoryPath(name), historyEntry)
}
