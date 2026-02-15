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
		// Bool-like flags
		applyBool := func(name string, target *bool, val string) {
			if !cmd.Flags().Changed(name) {
				return
			}
			if val == "true" {
				*target = true
			} else if val == "false" {
				*target = false
			} else {
				fmt.Printf("invalid value for %s: %s\n", name, val)
				os.Exit(1)
			}
			changed = true
		}

		applyBool("tics", &cfg.Tics, flagTics)
		applyBool("stretch", &cfg.Stretch, flagStretch)
		applyBool("inline", &cfg.Inline, flagInline)
		applyBool("colorize", &cfg.Style.Colorize, flagColorize)
		applyBool("hours12", &cfg.Hours12, flagHours12)
		applyBool("live", &cfg.Live, flagLive)

		if cmd.Flags().Changed("sorting") {
			if !core.IsValidSortingMode(flagSorting) {
				fmt.Println("invalid sorting mode:", flagSorting)
				os.Exit(1)
			}
			cfg.Sorting = flagSorting
			changed = true
		}
		if cmd.Flags().Changed("sort-local-top") {
			if flagSortLocal == "true" {
				cfg.SortLocalTop = true
			} else if flagSortLocal == "false" {
				cfg.SortLocalTop = false
			} else {
				fmt.Println("invalid value for sort-local-top:", flagSortLocal)
				os.Exit(1)
			}
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
	// Reuse the same flags as root to keep UX consistent
	configCmd.Flags().StringVarP(&flagTimezones, "timezones", "t", "", "timezones to display, comma-separated")
	configCmd.Flags().StringVarP(&flagSymbols, "symbols", "s", "", "symbols to use for time blocks")
	configCmd.Flags().StringVar(&flagTics, "tics", "", "true/false")
	configCmd.Flags().StringVar(&flagStretch, "stretch", "", "true/false")
	configCmd.Flags().StringVar(&flagInline, "inline", "", "true/false")
	configCmd.Flags().StringVar(&flagColorize, "colorize", "", "true/false")
	configCmd.Flags().StringVar(&flagHours12, "hours12", "", "true/false")
	configCmd.Flags().StringVar(&flagLive, "live", "", "true/false")
	configCmd.Flags().StringVar(&flagSorting, "sorting", "", "sorting mode")
	configCmd.Flags().StringVar(&flagSortLocal, "sort-local-top", "", "true/false")
}
