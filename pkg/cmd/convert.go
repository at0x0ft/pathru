package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/domain"
	"github.com/spf13/cobra"
	"strings"
)

func NewConvertCommand(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "convert",
		Short: "Only printing converted exec arguments",
		Long:  `Usage: pathru [options] convert <runtime service name> <execute command> -- [command arguments & options]`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf(
					"arguments must be given more than 1 [actual = \"%v\"]",
					args,
				)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			prjOpts, baseService := opts.getProjectOptions(), opts.getBaseService()
			convertedArgs, err := domain.Convert(
				prjOpts,
				baseService,
				args,
			)
			if err != nil {
				return err
			}

			fmt.Printf("%v\n", strings.Join(convertedArgs, " "))
			return nil
		},
	}
}
