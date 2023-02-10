package node

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/3d0c/storage/pkg/config"
	"github.com/3d0c/storage/pkg/log"
	"github.com/3d0c/storage/pkg/node"
)

var (
	globalCtx context.Context
	globalWG  *sync.WaitGroup
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Storage Node API Server",
	Long:  `runs Storage Node API Server`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		log.InitLogger(config.Node().Logger)
		log.TheLogger().Debug("node component",
			zap.String("config", fmt.Sprintf("%#v", config.Node())))

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
		apiSrv *node.APIHTTPServer
		err    error
	)

	globalWG.Add(1)
	defer globalWG.Done()

	if apiSrv, err = node.NewAPIHTTPServer(); err != nil {
		log.TheLogger().Fatal("error initializing API server", zap.Error(err))
	}

	apiSrv.Run(globalCtx)
}
