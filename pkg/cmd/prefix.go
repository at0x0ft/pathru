package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/domain"
	"github.com/spf13/cobra"
	"strings"
)

func NewPrefixCommand(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "prefix",
		Short: "Only printing docker compose prefix command",
		Long:  `Usage: pathru [options] prefix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			prjOpts := opts.getProjectOptions()
			options := domain.Unmarshal(prjOpts)
			fmt.Printf("%v\n", strings.Join(options, " "))
			return nil
		},
	}
}
