package main

import (
	"fmt"
	"os"
	"strconv"

	"banking/services"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "create_account":
		if len(os.Args) != 4 {
			fmt.Println("Usage: create_account <name> <amount>")
			return
		}
		name := os.Args[2]
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Printf("Invalid amount: %s\n", os.Args[3])
			return
		}
		if err := services.CreateAccount(name, amount); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Account created for %s with balance %.2f\n", name, amount)

	case "transfer":
		if len(os.Args) != 5 {
			fmt.Println("Usage: transfer <from> <to> <amount>")
			return
		}
		from := os.Args[2]
		to := os.Args[3]
		amount, err := strconv.ParseFloat(os.Args[4], 64)
		if err != nil {
			fmt.Printf("Invalid amount: %s\n", os.Args[4])
			return
		}
		if err := services.Transfer(from, to, amount); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Transferred %.2f from %s to %s\n", amount, from, to)

	case "add_deposit":
		if len(os.Args) != 4 {
			fmt.Println("Usage: add_deposit <name> <amount>")
			return
		}
		name := os.Args[2]
		amount, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Printf("Invalid amount: %s\n", os.Args[3])
			return
		}
		if err := services.AddDeposit(name, amount); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Deposit of %.2f added to %s's savings\n", amount, name)

	case "accrue_interest":
		if err := services.AccrueInterest(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Interest accrued on all deposits")

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("BLU CLI Banking System")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create_account <name> <amount>   Create a new account")
	fmt.Println("  transfer <from> <to> <amount>    Transfer funds between accounts")
	fmt.Println("  add_deposit <name> <amount>      Add a savings deposit")
	fmt.Println("  accrue_interest                  Compound interest on all deposits")
}
