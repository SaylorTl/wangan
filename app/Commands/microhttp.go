package Commands

import (
	"github.com/spf13/cobra"
	"log"
	"wangxin2.0/app/Constant"
	"wangxin2.0/routes"
)

func init() {
	rootCmd.AddCommand(microHttpCmd)
}

var microHttpCmd = &cobra.Command{
	Use:   "microhttp",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		initRouter := routes.InitMicroRouter()
		initErr := initRouter.Run(":8444")
		if initErr != nil {
			log.Fatal(Constant.INIT_ROUTE_ERROR+", err = ", initErr)
			return
		}
	},
}
