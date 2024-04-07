package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	opts := createRootOptions()
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru [options] <subcommand> <runtime service name> <execute command> -- [command arguments & options]`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.parse()
		},
	}
	opts.set(cmd.PersistentFlags())
	cmd.AddCommand(
		NewConvertCommand(opts),
	)
	return cmd
}
