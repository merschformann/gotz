package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/merschformann/gotz/core"
	"github.com/spf13/cobra"
)

var (
	flagTimezones string
	flagSymbols   string
	flagTics      bool
	flagStretch   bool
	flagInline    bool
	flagColorize  bool
	flagHours12   bool
	flagLive      bool
	flagSorting   string
	flagSortLocal bool
	flagTime      string
)

// Version is set by main to enable the built-in cobra --version flag.
var Version string

// Execute runs the root command.
func Execute() {
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gotz",
	Short: "gotz - terminal timezone plots",
	Long: `gotz plots timezone information in the terminal.
Use gotz config to manage your persistent defaults, override via command-line flags if needed (=false for disabling).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		config, err := core.Load()
		if err != nil {
			fmt.Println("error loading configuration:", err)
			if ok, inErr := core.AskUser("Reset configuration?"); inErr != nil {
				fmt.Println("error asking user:", inErr)
				os.Exit(1)
			} else if ok {
				config, inErr = core.SaveDefault()
				if inErr != nil {
					fmt.Println("error resetting configuration:", inErr)
					os.Exit(1)
				}
			} else {
				os.Exit(0)
			}
		}

		// Work on a copy so flags are NOT sticky by default
		runConfig := config

		// Apply ephemeral flags when provided
		if cmd.Flags().Changed("timezones") {
			tzs, err := core.ParseTimezones(flagTimezones)
			if err != nil {
				fmt.Println("invalid timezones:", err)
				os.Exit(1)
			}
			runConfig.Timezones = tzs
		}
		if cmd.Flags().Changed("symbols") {
			runConfig.Style.Symbols = flagSymbols
			if err := core.CheckSymbolMode(runConfig.Style.Symbols); err != nil {
				fmt.Println("invalid symbols:", err)
				os.Exit(1)
			}
		}
		if cmd.Flags().Changed("tics") {
			runConfig.Tics = flagTics
		}
		if cmd.Flags().Changed("stretch") {
			runConfig.Stretch = flagStretch
		}
		if cmd.Flags().Changed("inline") {
			runConfig.Inline = flagInline
		}
		if cmd.Flags().Changed("colorize") {
			runConfig.Style.Colorize = flagColorize
		}
		if cmd.Flags().Changed("hours12") {
			runConfig.Hours12 = flagHours12
		}
		if cmd.Flags().Changed("live") {
			runConfig.Live = flagLive
		}
		if cmd.Flags().Changed("sorting") {
			if !core.IsValidSortingMode(flagSorting) {
				fmt.Println("invalid sorting mode:", flagSorting)
				os.Exit(1)
			}
			runConfig.Sorting = flagSorting
		}
		if cmd.Flags().Changed("sort-local-top") {
			runConfig.SortLocalTop = flagSortLocal
		}

		// Determine requested time (flag wins, otherwise last arg if looks like time)
		var rt time.Time
		if cmd.Flags().Changed("time") {
			r, err := core.ParseRequestTime(runConfig, flagTime)
			if err != nil {
				fmt.Println("invalid time:", err)
				os.Exit(1)
			}
			rt = r
		} else if len(args) > 0 {
			last := args[len(args)-1]
			if len(last) > 0 && last[0] >= '0' && last[0] <= '9' {
				r, err := core.ParseRequestTime(runConfig, last)
				if err != nil {
					fmt.Println("invalid time:", err)
					os.Exit(1)
				}
				rt = r
			}
		}

		// Run plotting using the ephemeral config
		if err := core.Plot(runConfig, rt); err != nil {
			fmt.Println("error plotting time:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&flagTimezones, "timezones", "t", "", "timezones to display, comma-separated")
	rootCmd.Flags().StringVarP(&flagSymbols, "symbols", "s", "", "symbols to use for time blocks")
	rootCmd.Flags().BoolVar(&flagTics, "tics", false, "whether to use local time tics on the time axis")
	rootCmd.Flags().BoolVar(&flagStretch, "stretch", false, "whether to stretch across the terminal width")
	rootCmd.Flags().BoolVar(&flagInline, "inline", false, "whether to display time info and bars in one line")
	rootCmd.Flags().BoolVar(&flagColorize, "colorize", false, "whether to colorize the symbols")
	rootCmd.Flags().BoolVar(&flagHours12, "hours12", false, "whether to use 12-hour clock")
	rootCmd.Flags().BoolVar(&flagLive, "live", false, "whether to display time live (quit via 'q' or 'Ctrl+C')")
	rootCmd.Flags().StringVar(&flagSorting, "sorting", "", "indicates how to sort the timezones")
	rootCmd.Flags().BoolVar(&flagSortLocal, "sort-local-top", false, "whether to keep the local timezone at the top")
	rootCmd.Flags().StringVarP(&flagTime, "time", "T", "", "time to display (e.g. 20:00 or 8pm)")

	rootCmd.AddCommand(configCmd)
}
