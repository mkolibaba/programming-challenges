package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	var rootCmd = &cobra.Command{Use: "sudoku-solver-go"}
	rootCmd.AddCommand(NewApproxCmd())
	rootCmd.AddCommand(NewSolverCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
