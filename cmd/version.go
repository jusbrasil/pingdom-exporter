// Copyright 2019 Veepee.
// Copyright 2016 Giant Swarm GmbH.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run:   versionRun,
	}

	version   string
	goVersion string
	gitCommit string
	osArch    string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Version:\t%v\n", version)
	fmt.Printf("Go version:\t%v\n", goVersion)
	fmt.Printf("Git commit:\t%v\n", gitCommit)
	fmt.Printf("OS/Arch:\t%v\n", osArch)
}
