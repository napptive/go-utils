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
	"fmt"

	"github.com/napptive/go-utils/pkg/validation"
	"github.com/rs/zerolog/log"
)

// Config contains the configuration elements related to the connection with gRPC services.
type Config struct {
	// Name of the connection as user information.
	Name string
	// ServerAddress with the dns/IP of the target gRPC server.
	ServerAddress string
	// ServerPort with the port of the catalog-manager gRPC server.
	ServerPort int
	// AuthEnable with a flag to indicate if the authentication is enabled or not
	AuthEnable bool
	// UseTLS indicates that a TLS connection is expected with the service.
	UseTLS bool
	// SkipCertValidation flag that enables ignoring the validation step of the certificate presented by the server.
	SkipCertValidation bool
	// ClientCA with a client trusted CA
	ClientCA string
}

// IsValid checks if the configuration options are valid.
func (cc *Config) IsValid() error {
	if err := validation.CheckNotEmpty(cc.ServerAddress, "serverAddress"); err != nil {
		return err
	}
	if err := validation.CheckPositive(cc.ServerPort, "serverPort"); err != nil {
		return err
	}

	return nil
}

// Print the configuration using the application logger.
func (cc *Config) Print() {
	log.Info().Str("name", cc.Name).Str("server", cc.ServerAddress).Int("Port", cc.ServerPort).Bool("useTLS", cc.UseTLS).Bool("skipCertValidation", cc.SkipCertValidation).Msg("Connection options")
}

// GetEffectiveAddress returns an address:port string.
func (cc *Config) GetEffectiveAddress() string {
	return fmt.Sprintf("%s:%d", cc.ServerAddress, cc.ServerPort)
}
