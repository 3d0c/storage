package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/apiserver"
	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/log"
)

var (
	globalCtx context.Context
	globalWG  *sync.WaitGroup
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Storage Proxy API Server",
	Long:  `runs Storage Proxy API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		log.InitLogger(config.Proxy().Logger)
		log.TheLogger().Debug("proxy component",
			zap.String("config", fmt.Sprintf("%#v", config.Proxy())))

		runProcesses()
		globalWG.Wait()
	},
}

func init() {
	var (
		cancel func()
	)

	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	globalCtx, cancel = context.WithCancel(context.Background())
	globalWG = &sync.WaitGroup{}

	globalWG.Add(1)
	go signalHandler(cancel)
}

func signalHandler(fn func()) {
	defer globalWG.Done()
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.TheLogger().Info("stop execution", zap.String("signal", sig.String()))
	fn()
	close(sigs)
}

func runProcesses() {
	var (
		apiSrv *apiserver.APIHTTPServer
		err    error
	)

	globalWG.Add(1)
	defer globalWG.Done()

	if apiSrv, err = apiserver.NewAPIHTTPServer(); err != nil {
		log.TheLogger().Fatal("error initializing API server", zap.Error(err))
	}

	apiSrv.Run(globalCtx)
}
