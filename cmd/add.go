// Copyright © 2018 Yasin Bahtiyar <yasin@bahtiyar.org>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	gdns "github.com/0x1EE7/cloudns/googledns"
	"github.com/spf13/cobra"
)

var addFlags gdns.DNSRecord

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add given IPs to the domain",
	// Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Adding IPs %v to %v\n", *addFlags.Ips, *addFlags.Domain)
		dns, err := gdns.NewDNSProvider()
		if err == nil {
			for i := 0; i < retryNum; i++ {
				err = dns.MakeChange(addFlags, true)
				if !RetryOn(err) {
					break
				}
			}
		}
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	ips := addCmd.Flags().IPSliceP("ip", "i", nil, "IP address list")
	addCmd.MarkFlagRequired("ip")
	domain := addCmd.Flags().StringP("domain", "d", "", "Domain address to remove")
	addCmd.MarkFlagRequired("domain")
	addFlags = gdns.DNSRecord{Ips: ips, Domain: domain}

}
