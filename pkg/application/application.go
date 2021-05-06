package application

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/isnlan/coral/pkg/logging"
	"github.com/mitchellh/go-homedir"
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
	logger  = logging.MustGetLogger("Application")
	cfgFile = ""
)

type Application struct {
	name    string
	cmd     *cobra.Command
	rotator logging.Rotator
}

func New(name string) *Application {
	cmd := &cobra.Command{
		Use: name,
		Run: nil,
	}
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	cobra.OnInitialize(initConfig)
	rotator := initializeLogging()

	Name = name

	return &Application{
		name:    name,
		cmd:     cmd,
		rotator: rotator,
	}
}

func (app *Application) Start(svr Server) {
	app.cmd.Run = func(cmd *cobra.Command, args []string) {
		err := svr.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		svr.Start()

		if app.rotator != nil {
			go func() {
				for {
					now := time.Now()
					next := now.Add(time.Hour * 24)
					next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
					t := time.NewTimer(next.Sub(now))
					<-t.C
					err := app.rotator.Rotate()
					if err != nil {
						log.Fatal(err)
					}
				}
			}()
		}

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for sig := range sigChan {
			fmt.Printf("get a signal %s\n", sig.String())
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				svr.Close()
				return
			case syscall.SIGHUP:
				if app.rotator != nil {
					err := app.rotator.Rotate()
					if err != nil {
						log.Fatal(err)
					}
				}
			default:
				return
			}
		}
	}
	if err := app.cmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func initializeLogging() logging.Rotator {
	loggingFile := os.Getenv("LOGGING_FILE")
	if loggingFile != "" {
		writer := &lumberjack.Logger{
			Filename:   loggingFile,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     1,    //days
			Compress:   true, // disabled by default
		}
		c := logging.NewFileConfig(writer)
		logging.Init(c)
		return writer
	} else {
		c := logging.Config{
			Format:  os.Getenv("LOGGING_FORMAT"),
			Writer:  os.Stderr,
			LogSpec: os.Getenv("LOGGING_SPEC"),
		}
		logging.Init(c)
		return nil
	}
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
