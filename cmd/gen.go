/*
Copyright © 2020 Daniel Hails <daniel@hails.info>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"time"
	"math"
	"context"
	"regexp"

	"github.com/DJRHails/vanity-algos/helpers"
	"github.com/DJRHails/vanity-algos/process"

	"github.com/spf13/cobra"
	"github.com/schollz/progressbar/v3"
	"github.com/dustin/go-humanize"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

const example = "42JQEQIQNZUGC5JPJU7MMYNNIMTZ5HXKNPIS6Z2YZXRZQIN764JRO2BV44"

func init() {
	genCmd.SetArgs([]string{""})
	genCmd.Flags().IntP("count", "n", 1, "Number of matching accounts to generate")
	genCmd.Flags().IntP("timeout", "t", 0, "Length of time to generate for.")
	genCmd.Flags().BoolP("regex", "r", false, "Use a regex expression.")

	rootCmd.AddCommand(genCmd)
}

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen [pattern]",
	Aliases: []string{"g"},
	Short: "Generate a new address",
	Long: `Generates a new algorand address based on the given pattern.`,
	Example:
	`vanity-algos gen 'HAILS' # Find an address starting with HAILS
vanity-algos gen -n 3 -r '^[0-9]{3}' # Find 3 addresses starting with 3 numbers`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGen,
}

func runGen(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	count, _ := cmd.Flags().GetInt("count")
	timeout, _ := cmd.Flags().GetInt("timeout")
	concurrency, _ := cmd.Flags().GetInt("concurrency")
	regex, _ := cmd.Flags().GetBool("regex")

	// Build Matcher
	// TODO(DJRHails): Should generate for regex too
	// https://ro-che.info/articles/2018-08-01-probability-of-regex
	diff, prob50, prob99 := generateStatistics(pattern)
	bar := progressbar.Default(prob99 * int64(count), "generating")

	results := make(chan *crypto.Account, count)
	matcher := &process.Matcher{
		Bar: bar,
		Results: results,
	}

	// Set as prefix or regex
	if regex {
		matcher.Regex = regexp.MustCompile(pattern)
	} else {
		err := validPattern(pattern)
		if err != nil {
			return err
		}
		matcher.Prefix = pattern
	}

	// Initial Output
	printExample(pattern)

	fmt.Printf("Difficulty:   %s\n", humanize.Comma(int64(diff)))
	fmt.Printf("50-50 chance: %s\n", humanize.Comma(prob50))

	// Start workers, and outputer
	var ctx = context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	foundCount := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case account := <-results:
				foundCount++
				passphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)

				if err != nil {
					fmt.Printf("Error creating transaction: %s\n", err)
					foundCount--
					return
				}

				fmt.Printf("\nFound address: %s\n", account.Address)
				fmt.Printf("passphrase: %s\n", passphrase)

				if count > 0 && foundCount >= count {
					// We have met our requirements, time to leave
					cancel()
					return
				}
			}
		}
	}()

	// Create a semaphore that will count up to concurrency
	semaphore := make(chan bool, concurrency)

	for {
		select {
		case <- ctx.Done():
			return nil
		default:
			semaphore <- true
			if ctx.Err() == nil {
				matcher.Run(ctx, semaphore)
			}
		}
	}

	return nil
}

func printExample(pattern string) {
	fmt.Printf("Searching for (%s): e.g. %s%s\n",
		pattern,
		pattern,
		example[len(pattern):],
	)
}

// https://math.stackexchange.com/questions/926028/probability-of-a-substring-occurring-in-a-string
func generateStatistics(pattern string) (float64, int64, int64){
	diff := math.Pow(float64(len(helpers.Base32RuneSet)), float64(len(pattern)))
	prob50 := math.Log(0.5) / math.Log(1 - 1/diff)
	prob99 := math.Log(0.01) / math.Log(1 - 1/diff)
	return diff, int64(prob50), int64(prob99)
}

func validPattern(pattern string) error {
	if !helpers.IsBase32(pattern) {
		return fmt.Errorf("Pattern '%s' is not valid base32 ([A-Z,2-7])", pattern)
	}
	return nil
}
