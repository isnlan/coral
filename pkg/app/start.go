package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/snlansky/coral/pkg/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdRoot = "SVR"

var Name string

type Server interface {
	Initialize() error
	Start()
	Close()
}

var (
	logger  = logging.MustGetLogger("App")
	cfgFile = ""
)

type App struct {
	name string
	cmd  *cobra.Command
}

func New(name string) *App {
	cmd := &cobra.Command{
		Use: name,
		Run: nil,
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	cobra.OnInitialize(initConfig)
	initializeLogging()

	Name = name

	return &App{
		name: name,
		cmd:  cmd,
	}
}

func (app *App) Start(svr Server) {
	app.cmd.Run = func(cmd *cobra.Command, args []string) {
		err := svr.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		svr.Start()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for sig := range sigChan {
			fmt.Printf("get a signal %s", sig.String())
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				svr.Close()
				return
			case syscall.SIGHUP:
				//logger.Rotate(false)
			default:
				return
			}
		}
	}
	if err := app.cmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func initializeLogging() {
	loggingSpec := os.Getenv("LOGGING_SPEC")
	loggingFormat := os.Getenv("LOGGING_FORMAT")
	logging.Init(logging.Config{
		Format:  loggingFormat,
		Writer:  os.Stderr,
		LogSpec: loggingSpec,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix(cmdRoot)
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else if os.Getenv("CONFIG_FILE") != "" {
		viper.SetConfigFile(os.Getenv("CONFIG_FILE"))
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logger.Fatal(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName("config-dev")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Using config file:%s\n", viper.ConfigFileUsed())
}

func InitConfig(config interface{}) error {
	err := viper.Unmarshal(config)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(buf))
	return nil
}
