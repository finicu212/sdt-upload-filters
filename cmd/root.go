package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sdt-upload-filters/pkg/connection"
)

var rootCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload data to an FTP server. Featuring Handler chains, decorators and connection pools!",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool := connection.NewPool("", 21)
		c, err := pool.GetConnection("root", "root")
		if err != nil {
			return err
		}

		filePath := "file.txt"
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		err = c.Store("file.txt", file)
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
