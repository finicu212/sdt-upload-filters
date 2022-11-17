package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"sdt-upload-filters/pkg/connection"
)

var rootCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload data to an FTP server. Featuring Handler chains, decorators and connection pools!",
	RunE: func(cmd *cobra.Command, args []string) error {
		ips, err := cmd.Flags().GetIPSlice(FlagUrl)
		if err != nil {
			return err
		}
		var conns []connection.IConnection

		fmt.Printf("%s\n", c.GetUUID())

		//filePath := "file.txt"
		//_, err = os.Open(filePath)
		//if err != nil {
		//	return err
		//}

		//err = c.Store("file.txt", file)
		//if err != nil {
		//	return err
		//}
		return nil
	},
}

const (
	FlagFile  = "file"
	FlagFileP = "f"
	FlagUrl   = "url"
	FlagUrlP  = "u"
	FlagPort  = "port"
	FlagPortP = "p"
)

func init() {
	rootCmd.Flags().StringSliceP(FlagFile, FlagFileP, []string{"file.txt"}, "A list of comma separated files you want to upload in parallel, if the available connections permit it.\nThis is to permit integration of the connection pool pattern")
	rootCmd.Flags().IPSliceP(FlagUrl, FlagUrlP, []net.IP{net.IP("172.17.0.1:21")}, "A list of URLs you want to upload the files to")
	//rootCmd.Flags().StringP(FlagPort, FlagPortP, "21", "The port you want to connect to")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
