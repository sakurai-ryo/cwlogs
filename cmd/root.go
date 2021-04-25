package cmd

import (
	"bufio"
	"context"
	"cwlogs/aws"
	"errors"
	"fmt"
	"io"
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

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	// cobra.OnInitialize(initConfig)
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&profile, "profile", "p", "profile(default is $HOME/.aws/credentials)")
	pf.StringVar(&region, "region", "r", "target log group region")
	pf.StringVar(&prefix, "prefix", "x", "log group prefix")
	cobra.MarkFlagRequired(pf, "profile")
	cobra.MarkFlagRequired(pf, "region")
	cobra.MarkFlagRequired(pf, "prefix")
}

func do() error {
	c := aws.NewCW(region)
	ctx := context.Background()
	names, err := aws.ListLogGroup(ctx, c, prefix)
	if err != nil {
		return err
	}
	t := usePrompt(names)

	var cmd *exec.Cmd
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	errChan := make(chan error)
	if err := startCmd(sigCtx, cmd, t, errChan); err != nil {
		return err
	}

	for {
		select {
		case <-sigCtx.Done():
			if ctx.Err() != nil {
				return ctx.Err()
			} else {
				return err
			}
		case err := <-errChan:
			return err
		}
	}
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

func startCmd(ctx context.Context, cmd *exec.Cmd, name string, errChan chan error) error {
	cmd = exec.CommandContext(ctx, "cw", "tail", "-f", name)
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go cwReader(outReader)
	go cwErrorHandler(errReader, errChan)
	return nil
}

func cwReader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
}

func cwErrorHandler(e io.ReadCloser, errChan chan error) {
	defer e.Close()
	scanner := bufio.NewScanner(e)
	for scanner.Scan() {
		errChan <- errors.New(scanner.Text())
	}
}
