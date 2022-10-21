/**
 * Copyright 2022 Napptive
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"os"

	"github.com/napptive/go-utils/pkg/printer"
	"github.com/spf13/cobra"
)

// CrashOnError prints the error if found and returns a non-zero value as the result of the playground CLI execution.
func CrashOnError(err error) {
	if err != nil {
		printer.PrintError(err)
		os.Exit(1)
	}
}

// CrashWithHelp shows the command help before exiting.
func CrashWithHelp(cmd *cobra.Command) {
	_ = cmd.Help()
	os.Exit(1)
}
