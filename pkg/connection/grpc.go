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
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"

	"github.com/napptive/nerrors/pkg/nerrors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// GetConnection creates a connection with a gRPC server.
func GetConnection(cfg *Config) (*grpc.ClientConn, error) {
	if cfg.UseTLS {
		return GetTLSConnection(cfg, cfg.GetEffectiveAddress())
	}
	return GetNonTLSConnection(cfg, cfg.GetEffectiveAddress())
}

// GetTLSConnection returns a TLS wrapped connection with the playground server.
func GetTLSConnection(cfg *Config, address string) (*grpc.ClientConn, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.SkipCertValidation,
	}
	if cfg.ClientCA != "" {
		cp := x509.NewCertPool()
		decoded, err := base64.StdEncoding.DecodeString(cfg.ClientCA)
		if err != nil {
			return nil, nerrors.NewInternalErrorFrom(err, "error decoding CA")
		}
		if !cp.AppendCertsFromPEM(decoded) {
			return nil, nerrors.NewInternalError("Error appending CA")
		}
		// add the CA as valid one
		tlsConfig.RootCAs = cp
	}
	tlsCredentials := credentials.NewTLS(tlsConfig)
	return grpc.Dial(address, grpc.WithTransportCredentials(tlsCredentials))
}

// GetNonTLSConnection returns a plain connection with the playground server.
func GetNonTLSConnection(_ *Config, address string) (*grpc.ClientConn, error) {
	log.Warn().Str("address", address).Msg("using insecure connection")
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
