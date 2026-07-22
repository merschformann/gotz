package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/merschformann/gotz/core"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit persistent configuration",
	Long: `View or edit the persistent configuration stored in $XDG_CONFIG_HOME/gotz/config.json.

Boolean flags: use --flag to enable, --flag=false to disable the setting.
For example: gotz config --live       enables live mode permanently
             gotz config --live=false disables live mode permanently`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := core.Load()
		if err != nil {
			fmt.Println("error loading configuration:", err)
			if ok, inErr := core.AskUser("Reset configuration?"); inErr != nil {
				fmt.Println("error asking user:", inErr)
				os.Exit(1)
			} else if ok {
				cfg, inErr = core.SaveDefault()
				if inErr != nil {
					fmt.Println("error resetting configuration:", inErr)
					os.Exit(1)
				}
			} else {
				os.Exit(0)
			}
		}

		changed := false

		if cmd.Flags().Changed("timezones") {
			tzs, err := core.ParseTimezones(flagTimezones)
			if err != nil {
				fmt.Println("invalid timezones:", err)
				os.Exit(1)
			}
			cfg.Timezones = tzs
			changed = true
		}
		if cmd.Flags().Changed("symbols") {
			if err := core.CheckSymbolMode(flagSymbols); err != nil {
				fmt.Println("invalid symbols:", err)
				os.Exit(1)
			}
			cfg.Style.Symbols = flagSymbols
			changed = true
		}
		if cmd.Flags().Changed("tics") {
			cfg.Tics = flagTics
			changed = true
		}
		if cmd.Flags().Changed("stretch") {
			cfg.Stretch = flagStretch
			changed = true
		}
		if cmd.Flags().Changed("inline") {
			cfg.Inline = flagInline
			changed = true
		}
		if cmd.Flags().Changed("colorize") {
			cfg.Style.Colorize = flagColorize
			changed = true
		}
		if cmd.Flags().Changed("hours12") {
			cfg.Hours12 = flagHours12
			changed = true
		}
		if cmd.Flags().Changed("live") {
			cfg.Live = flagLive
			changed = true
		}
		if cmd.Flags().Changed("sorting") {
			if !core.IsValidSortingMode(flagSorting) {
				fmt.Println("invalid sorting mode:", flagSorting)
				os.Exit(1)
			}
			cfg.Sorting = flagSorting
			changed = true
		}
		if cmd.Flags().Changed("sort-local-top") {
			cfg.SortLocalTop = flagSortLocal
			changed = true
		}

		if changed {
			if err := cfg.Save(); err != nil {
				fmt.Println("error saving config:", err)
				os.Exit(1)
			}
			fmt.Println("configuration saved")
			return
		}

		// If no changes, print config
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(data))
	},
}

func init() {
	configCmd.Flags().StringVarP(&flagTimezones, "timezones", "t", "", "timezones to display, comma-separated")
	configCmd.Flags().StringVarP(&flagSymbols, "symbols", "s", "", "symbols to use for time blocks")
	configCmd.Flags().BoolVar(&flagTics, "tics", false, "whether to use local time tics on the time axis")
	configCmd.Flags().BoolVar(&flagStretch, "stretch", false, "whether to stretch across the terminal width")
	configCmd.Flags().BoolVar(&flagInline, "inline", false, "whether to display time info and bars in one line")
	configCmd.Flags().BoolVar(&flagColorize, "colorize", false, "whether to colorize the symbols")
	configCmd.Flags().BoolVar(&flagHours12, "hours12", false, "whether to use 12-hour clock")
	configCmd.Flags().BoolVar(&flagLive, "live", false, "whether to display time live (quit via 'q' or 'Ctrl+C')")
	configCmd.Flags().StringVar(&flagSorting, "sorting", "", "indicates how to sort the timezones")
	configCmd.Flags().BoolVar(&flagSortLocal, "sort-local-top", false, "whether to keep the local timezone at the top")
}
