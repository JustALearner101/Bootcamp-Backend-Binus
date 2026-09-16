package services

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//AddDeposit nambahin depositan ke file deposit.txt
func AddDeposit(name string, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}

	if err := ensureDir(depositsDir); err != nil {
		return fmt.Errorf("failed to create deposits directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("%.2f|%s|0\n", amount, timestamp)

	depositPath := getDepositPath(name)
	f, err := os.OpenFile(depositPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open deposit file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		return fmt.Errorf("failed to write deposit: %w", err)
	}

	return nil
}

//AccrueInterest nambahin bunga ke semua depositan
func AccrueInterest() error {
	if err := ensureDir(depositsDir); err != nil {
		return fmt.Errorf("failed to access deposits directory: %w", err)
	}

	entries, err := os.ReadDir(depositsDir)
	if err != nil {
		return fmt.Errorf("failed to read deposits directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".txt")
		if err := accrueInterestForDeposits(name); err != nil {
			return fmt.Errorf("failed to accrue interest for %s: %w", name, err)
		}
	}

	return nil
}

// accrueInterestForDeposits nambahin bunga ke depositan user tertentu
func accrueInterestForDeposits(name string) error {
	depositPath := getDepositPath(name)
	data, err := os.ReadFile(depositPath)
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil
	}

	var newLines []string
	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			newLines = append(newLines, line)
			continue
		}

		amount, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			newLines = append(newLines, line)
			continue
		}

		depositTime, err := time.Parse("2006-01-02 15:04:05", parts[1])
		if err != nil {
			newLines = append(newLines, line)
			continue
		}

		minutesElapsed := now.Sub(depositTime).Minutes()
		if minutesElapsed > 0 {
			periods := int(minutesElapsed)
			// Compound interest: A = P * (1 + r)^n
			newAmount := amount * math.Pow(1+interestRate, float64(periods))
			newLines = append(newLines, fmt.Sprintf("%.2f|%s|%d", newAmount, parts[1], periods))
		} else {
			newLines = append(newLines, line)
		}
	}

	content := strings.Join(newLines, "\n") + "\n"
	if err := os.WriteFile(depositPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to update deposit file: %w", err)
	}

	_ = filepath.Join(historyDir, name+".txt")
	historyEntry := fmt.Sprintf("[%s] Interest accrued on deposits\n", timestamp)
	return appendToFile(getHistoryPath(name), historyEntry)
}

// ini cuma buat append file, create file kalo gaada
func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
