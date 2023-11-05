package Commands

import (
	"github.com/spf13/cobra"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/routes"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("Exception1111", err)
					return
				}
			}()
			initMicroRouter := routes.InitMicroRouter()
			initMicroErr := initMicroRouter.Run(":8444")
			if initMicroErr != nil {
				log.Fatal(Constant.INIT_ROUTE_ERROR+", err = ", initMicroErr)
				return
			}
		}()
		initRouter := routes.InitRouter()
		initErr := initRouter.RunTLS(":8445", Constant.FociiCrt, Constant.FociiCerts)
		if initErr != nil {
			log.Fatal(Constant.INIT_ROUTE_ERROR+", err = ", initErr)
			return
		}

	},
}
