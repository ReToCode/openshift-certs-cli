// Copyright © 2018 SBB Cloud Stack Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/op/go-logging"
	"time"
)

var jsonFile string
var expiryDays int
var debug bool

var log = logging.MustGetLogger("openshift-certs-cli")

var rootCmd = &cobra.Command{
	Use:   "openshift-certs-cli",
	Short: "This cli parses 'cert-expiry-report.json' and outputs expired certs.",
	Long: `OpenShift uses SSL certificates for encrypting communication between its components. It's crucial to monitor
their expiry date and renew them as needed. The JSON file cert-expiry-report.json is generated via /usr/share/ansible/openshift-ansible/playbooks/certificate_expiry/easy-mode.yaml.`,
	Run: printExpiredCertificates,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.Flags().StringVarP(&jsonFile, "file","f", "/tmp/cert-expiry-report.json", "location of the JSON file (default is /tmp/cert-expiry-report.json)")
	rootCmd.Flags().IntVarP(&expiryDays, "expiry",  "e", 90, "number of days left before cert expires")
	rootCmd.Flags().BoolVarP(&debug, "debug",  "d", false, "print debug messages")
}

func initLogging() {
	var format = logging.MustStringFormatter(
		`%{message}`,
	)
	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logging.NewBackendFormatter(stdOutBackend, format))

	if debug {
		logging.SetLevel(logging.DEBUG, "openshift-certs-cli")
	} else {
		logging.SetLevel(logging.INFO, "openshift-certs-cli")
	}

}

func printTypeExpiry(entry []*CertEntry, server string) {
	for _, el := range entry {

		if el.DaysRemaining <= expiryDays {
			log.Infof("%d days left until %s for %s @ %s: %s", el.DaysRemaining, el.Expiry, el.Path, server, el.CertCn)
		} else {
			log.Debugf("%d days left until %s for %s @ %s: %s", el.DaysRemaining, el.Expiry, el.Path, server, el.CertCn)
		}
	}
}

func printExpiredCertificates(cmd *cobra.Command, args []string) {
	log.Debugf("Parsing JSON @ %s. Expiry is set to %d days.", jsonFile, expiryDays)

	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Errorf("Can't open the JSON file (%s).", err.Error())
		os.Exit(1)
	}

	var certExpiryReport CertExpiryReport
	err = json.Unmarshal(bytes, &certExpiryReport)
	if err != nil {
		log.Errorf("Can't unmarshall JSON file (%s).", err.Error())
		os.Exit(1)
	}

	layout := "2006-01-02 15:04:05.000000"

	for k, v := range certExpiryReport.Data {
		checked, err := time.Parse(layout, v.Meta.CheckedAtTime)
		if err != nil {
			log.Warningf("Can't parse date (%s).", err.Error())
		} else {
			log.Debugf("%s checked @ %s.", k, checked)
			daysSinceCheck := time.Now().Sub(checked).Hours()/24
			if daysSinceCheck > 1 {
				log.Errorf("JSON file was generated @ %s (%.1f days ago).", checked, daysSinceCheck)
			}
		}

		printTypeExpiry(v.Etcd, k)
		printTypeExpiry(v.Kubeconfigs, k)
		printTypeExpiry(v.OcpCerts, k)
		printTypeExpiry(v.Registry, k)
		printTypeExpiry(v.Router, k)
	}
}