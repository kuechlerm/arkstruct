package cmd

import (
	"arkstruct/generate"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate arktype types and RPC-Client from Go structs",
	Long: `Generate arktype types and RPC-Client from Go structs.
	Example:

	arkstruct generate -i /path/to/folder -o output.ts
	`,
	Run: func(cmd *cobra.Command, args []string) {
		in, _ := cmd.Flags().GetString("input")
		out, _ := cmd.Flags().GetString("output")

		err := generate.Generate(in, out)
		if err != nil {
			cmd.PrintErrf("Error generating types: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input", "i", "", "Folder with Go files containing structs")
	generateCmd.Flags().StringP("output", "o", "", "Output TypeScript file for generated types")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
