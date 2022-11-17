package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sdt-upload-filters/pkg/connection"
	"sdt-upload-filters/pkg/utils"
)

const (
	FlagFiles  = "files"
	FlagFilesP = "f"
	FlagUrls   = "urls"
	FlagUrlsP  = "U"
	FlagUsers  = "users"
	FlagUsersP = "u"
	FlagPass   = "passwords"
	FlagPassP  = "p"
)

var (
	ErrUrlsNotEnoughCreds = errors.New(fmt.Sprintf("--%s and --%s must have a length of one, or the same length as --%s", FlagUsers, FlagPass, FlagUrls))
)

func rootCmd() *cobra.Command {
	var usernames []string
	var passwords []string
	var ips []string
	var files []string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload data to an FTP server. Featuring Handler chains, decorators and connection pools!",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.MarkFlagFilename(FlagFiles)
			if err != nil {
				return err
			}

			// To avoid extra work of handling this adequately... this be okay enough to do, honestly
			if len(usernames) == 1 {
				usernames = utils.Repeated(usernames[0], len(ips))
			}
			if len(passwords) == 1 {
				passwords = utils.Repeated(passwords[0], len(ips))
			}

			if len(usernames) != len(ips) {
				return ErrUrlsNotEnoughCreds
			}
			if len(passwords) != len(ips) {
				return ErrUrlsNotEnoughCreds
			}

			log.Printf("%+v\n", usernames)
			log.Printf("%+v\n", passwords)
			log.Printf("%+v\n", ips)
			log.Printf("%+v\n", files)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			orchestrator, err := connection.NewOrchestrator(ips, usernames, passwords)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", orchestrator)
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&usernames, FlagUsers, FlagUsersP, []string{os.Getenv("FTPUSER")}, "A list of comma separated usernames, to use for connecting to each URL. Must be only one, or the same length as `--files` flag")
	cmd.Flags().StringSliceVarP(&passwords, FlagPass, FlagPassP, []string{os.Getenv("PASS")}, "A list of comma separated passwords, to use for connecting to each URL. Must be only one, or the same length as `--files` flag")
	cmd.Flags().StringSliceVarP(&ips, FlagUrls, FlagUrlsP, []string{"82.79.159.78:21"}, "A list of comma separated URLs you want to upload the files to")
	cmd.Flags().StringSliceVarP(&files, FlagFiles, FlagFilesP, []string{"file.txt"}, "A list of comma separated files you want to upload in parallel, if the available connections permit it.\nIf not enough connections are available in the pool, we will split the load sequentially")

	cmd.Flags().SortFlags = false

	return cmd
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
