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

package connection

import (
	"context"
	"time"

	"github.com/napptive/go-utils/pkg/printer"
	"google.golang.org/grpc/metadata"
)

// AgentHeader with the key name for the agent payload.
const AgentHeader = "agent"

// VersionHeader with the key name for the version payload.
const VersionHeader = "version"

// ContextTimeout with the default timeout for Napptive playground operations.
const ContextTimeout = 5 * time.Minute

// ContextHelper structure to facilitate the generation of secure contexts.
type ContextHelper struct {
	// Version of the application sending the request.
	Version string
	// Agent sending the request.
	Agent string

	printer.ResultPrinter
}

// NewContextHelper creates a ContextHelper with a given configuration.
func NewContextHelper(version string, agent string, printer printer.ResultPrinter) *ContextHelper {
	return &ContextHelper{
		Version:       version,
		Agent:         agent,
		ResultPrinter: printer,
	}
}

// GetContext returns a valid gRPC context with the appropriate authorization header.
func (ch *ContextHelper) GetContext() (context.Context, context.CancelFunc) {
	md := metadata.New(map[string]string{AgentHeader: ch.Agent, VersionHeader: ch.Version})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return context.WithTimeout(ctx, ContextTimeout)
}
