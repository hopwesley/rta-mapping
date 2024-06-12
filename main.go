package main

import (
	"encoding/json"
	"fmt"
	"github.com/hopwesley/fdlimit"
	"github.com/hopwesley/rta-mapping/common"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	ConfigFIleName = "config.json"
)

var (
	param = &startParam{}
)

type startParam struct {
	version bool
	config  string
}

var rootCmd = &cobra.Command{
	Use: "rtaHint",

	Short: "rtaHint",

	Long: `usage description::TODO::`,

	Run: mainRun,
}

func init() {
	flags := rootCmd.Flags()
	flags.BoolVarP(&param.version, "version",
		"v", false, "rta-map -v")
	flags.StringVarP(&param.config, "conf",
		"c", ConfigFIleName, "rta-map -c config.json")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func mainRun(_ *cobra.Command, _ []string) {

	if param.version {
		fmt.Println("\n==================================================")
		fmt.Printf("Version:\t%s\n", common.Version)
		fmt.Printf("Build:\t\t%s\n", common.BuildTime)
		fmt.Printf("Commit:\t\t%s\n", common.Commit)
		fmt.Println("==================================================")
		return
	}

	if err := fdlimit.MaxIt(); err != nil {
		panic(err)
	}
	cfg := initConfig(param.config)

	if err := InitIDMap(cfg.MysqlCfg); err != nil {
		panic(err)
	}
	common.PrintMemUsage()

	if err := InitRtaMap(cfg.RedisCfg); err != nil {
		panic(err)
	}

	common.PrintMemUsage()

	var srv = NewHttpService()
	go func() {
		srv.Start()
	}()

	waitShutdownSignal()
}

func initConfig(filName string) *Config {
	cf := new(Config)

	bts, err := os.ReadFile(filName)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(bts, &cf); err != nil {
		panic(err)
	}

	_sysConfig = cf
	return cf
}

func waitShutdownSignal() {
	var pidFile = os.Args[0] + ".pid"
	pid := strconv.Itoa(os.Getpid())
	fmt.Printf("\n>>>>>>>>>>service start at pid(%s)<<<<<<<<<<\n", pid)
	if err := os.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		fmt.Print("failed to write running pid", err)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2)
	sig := <-sigCh
	fmt.Printf("\n>>>>>>>>>>service finished(%s)<<<<<<<<<<\n", sig)
}
