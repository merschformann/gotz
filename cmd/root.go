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
	flagTics      string
	flagStretch   string
	flagInline    string
	flagColorize  string
	flagHours12   string
	flagLive      string
	flagSorting   string
	flagSortLocal string
	flagTime      string
)

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gotz",
	Short: "gotz - terminal timezone plots",
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
		// Helper to parse bool-like flags
		applyBool := func(name string, target *bool, val string) error {
			if !cmd.Flags().Changed(name) {
				return nil
			}
			switch l := val; l {
			case "true", "TRUE", "True":
				*target = true
			case "false", "FALSE", "False":
				*target = false
			default:
				return fmt.Errorf("invalid value for %s: %s", name, val)
			}
			return nil
		}

		if err := applyBool("tics", &runConfig.Tics, flagTics); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := applyBool("stretch", &runConfig.Stretch, flagStretch); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := applyBool("inline", &runConfig.Inline, flagInline); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := applyBool("colorize", &runConfig.Style.Colorize, flagColorize); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := applyBool("hours12", &runConfig.Hours12, flagHours12); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := applyBool("live", &runConfig.Live, flagLive); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cmd.Flags().Changed("sorting") {
			if !core.IsValidSortingMode(flagSorting) {
				fmt.Println("invalid sorting mode:", flagSorting)
				os.Exit(1)
			}
			runConfig.Sorting = flagSorting
		}
		if cmd.Flags().Changed("sort-local-top") {
			if flagSortLocal == "true" {
				runConfig.SortLocalTop = true
			} else if flagSortLocal == "false" {
				runConfig.SortLocalTop = false
			} else {
				fmt.Println("invalid value for sort-local-top:", flagSortLocal)
				os.Exit(1)
			}
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
	// String flags
	rootCmd.Flags().StringVarP(&flagTimezones, "timezones", "t", "", "timezones to display, comma-separated")
	rootCmd.Flags().StringVarP(&flagSymbols, "symbols", "s", "", "symbols to use for time blocks")
	rootCmd.Flags().StringVar(&flagTics, "tics", "", "indicates whether to use local time tics on the time axis (true/false)")
	rootCmd.Flags().StringVar(&flagStretch, "stretch", "", "indicates whether to stretch across the terminal width (true/false)")
	rootCmd.Flags().StringVar(&flagInline, "inline", "", "indicates whether to display time info and bars in one line (true/false)")
	rootCmd.Flags().StringVar(&flagColorize, "colorize", "", "indicates whether to colorize the symbols (true/false)")
	rootCmd.Flags().StringVar(&flagHours12, "hours12", "", "indicates whether to use 12-hour clock (true/false)")
	rootCmd.Flags().StringVar(&flagLive, "live", "", "indicates whether to display time live (true/false)")
	rootCmd.Flags().StringVar(&flagSorting, "sorting", "", "indicates how to sort the timezones")
	rootCmd.Flags().StringVar(&flagSortLocal, "sort-local-top", "", "indicates whether to keep the local timezone at the top (true/false)")
	rootCmd.Flags().StringVarP(&flagTime, "time", "T", "", "time to display (e.g. 20:00 or 8pm)")

	// Add config subcommand
	rootCmd.AddCommand(configCmd)
}
