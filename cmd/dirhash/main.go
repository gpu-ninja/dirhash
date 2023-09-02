/* SPDX-License-Identifier: Apache-2.0
 *
 * Copyright 2023 Damian Peckett <damian@pecke.tt>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"golang.org/x/mod/sumdb/dirhash"
)

func main() {
	app := &cli.App{
		Name:      "dirhash",
		Usage:     "Cryptographically checksums a directory and its contents.",
		ArgsUsage: "[directory]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "Sign the hash with the provided ed25519 key.",
			},
			&cli.StringFlag{
				Name:    "key-passphrase",
				Usage:   "Optional passphrase for the provided ed25519 key.",
				EnvVars: []string{"DIRHASH_KEY_PASSPHRASE"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}

			dir := c.Args().Get(0)

			hash, err := dirhash.HashDir(dir, "", dirhash.Hash1)
			if err != nil {
				return fmt.Errorf("failed to hash directory: %s", err)
			}

			if c.IsSet("key") {
				keyData, err := os.ReadFile(c.String("key"))
				if err != nil {
					return fmt.Errorf("failed to read key: %s", err)
				}

				var signer ssh.Signer
				if c.IsSet("key-passphrase") {
					signer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(c.String("key-passphrase")))
				} else {
					signer, err = ssh.ParsePrivateKey(keyData)
				}
				if err != nil {
					return fmt.Errorf("failed to parse key: %s", err)
				}

				if signer.PublicKey().Type() != "ssh-ed25519" {
					return fmt.Errorf("only ed25519 keys are supported")
				}

				signature, err := signer.Sign(rand.Reader, []byte(hash))
				if err != nil {
					return fmt.Errorf("failed to sign hash: %s", err)
				}

				hash += fmt.Sprintf(",s1:%s", base64.StdEncoding.EncodeToString(signature.Blob))
			}

			fmt.Println(hash)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "verify",
				Usage:     "Verify a directory hash.",
				ArgsUsage: "[hash] [directory]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "Verify the hash signature with the provided ed25519 public key.",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() != 2 {
						return cli.ShowAppHelp(c)
					}

					dir := c.Args().Get(1)

					values := strings.Split(c.Args().Get(0), ",")
					hash := values[0]

					var signature *ssh.Signature
					if len(values) > 1 {
						if !c.IsSet("key") {
							return fmt.Errorf("found signature but no key provided")
						}

						signatureBlob, err := base64.StdEncoding.DecodeString(values[1][3:])
						if err != nil {
							return fmt.Errorf("failed to decode signature: %s", err)
						}

						signature = &ssh.Signature{
							Format: "ssh-ed25519",
							Blob:   signatureBlob,
						}
					}

					computedHash, err := dirhash.HashDir(dir, "", dirhash.Hash1)
					if err != nil {
						return fmt.Errorf("failed to hash directory: %s", err)
					}

					if computedHash != hash {
						return fmt.Errorf("expected hash %s, got %s", hash, computedHash)
					}

					if signature != nil {
						publicKeyData, err := os.ReadFile(c.String("key"))
						if err != nil {
							return fmt.Errorf("failed to read public key: %s", err)
						}

						publicKey, _, _, _, err := ssh.ParseAuthorizedKey(publicKeyData)
						if err != nil {
							return fmt.Errorf("failed to parse public key: %s", err)
						}

						if err := publicKey.Verify([]byte(computedHash), signature); err != nil {
							return fmt.Errorf("failed to verify signature: %s", err)
						}
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
