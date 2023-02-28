/*
 *   Copyright 2023 Napptive
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *        https://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package structs

import (
	"encoding/json"
	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
)

// Convert transform a struct into another one. This method is used to convert into playground structs the internal ones
func Convert(entry interface{}, result interface{}) error {
	to, err := json.Marshal(entry)
	if err != nil {
		log.Error().Err(err).Msg("error in marshal")
		return nerrors.NewInternalError("unable to marshal entry")
	}
	err = json.Unmarshal(to, &result)
	if err != nil {
		log.Error().Err(err).Str("to", string(to)).Msg("error in unmarshal")
		return nerrors.NewInternalError("unable to unmarshal entry")
	}
	return nil
}
