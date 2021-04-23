package cmd

import (
	"context"
	"cwlogs/aws"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

var profile string
var region string
var prefix string

var rootCmd = &cobra.Command{
	Use:   "cwlogs",
	Short: "cwlogs tail Cloud Watch Logs",
	Long:  ``,
	// PreRunE: //TODO あとでバリデーションを追加する
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := do(); err != nil {
			return err
		}
		return nil
	},
}

func do() error {
	c := aws.NewCW(region)
	ctx := context.Background()
	names, err := aws.ListLogGroup(ctx, c, prefix)
	if err != nil {
		return err
	}
	t := usePrompt(names)
	// fmt.Sprintf("tail %s", t)

	var cmd *exec.Cmd
	go func() error {
		cmd = exec.Command("cw", "tail", t, "-f")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case <-quit:
			log.Print("Ctr+c Clicked")
			if err := cmd.Process.Kill(); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func usePrompt(names []string) string {
	var s []prompt.Suggest
	for _, n := range names {
		s = append(s, prompt.Suggest{Text: n})
	}
	fmt.Println("Please select logGroup.")
	t := prompt.Input(">> ", completer(s))
	fmt.Println("Display Logs: ", t)
	return t
}

func completer(s []prompt.Suggest) prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "p", "profile(default is $HOME/.aws/credentials)")
	rootCmd.PersistentFlags().StringVar(&region, "region", "r", "target log group region")
	rootCmd.PersistentFlags().StringVar(&prefix, "prefix", "x", "log group prefix")
}

// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		cobra.CheckErr(err)

// 		// Search config in home directory with name ".cwlogs" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".cwlogs")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
// 	}
// }
