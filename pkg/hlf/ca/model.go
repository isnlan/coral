/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ca

import (
	"encoding/asn1"

	"github.com/cloudflare/cfssl/csr"
)

const (
	commonNameLength             = 64
	serialNumberLength           = 64
	countryNameLength            = 2
	localityNameLength           = 128
	stateOrProvinceNameLength    = 128
	organizationNameLength       = 64
	organizationalUnitNameLength = 64
)

var (
	// The X.509 BasicConstraints object identifier (RFC 5280, 4.2.1.9)
	basicConstraintsOID   = asn1.ObjectIdentifier{2, 5, 29, 19}
	commonNameOID         = asn1.ObjectIdentifier{2, 5, 4, 3}
	serialNumberOID       = asn1.ObjectIdentifier{2, 5, 4, 5}
	countryOID            = asn1.ObjectIdentifier{2, 5, 4, 6}
	localityOID           = asn1.ObjectIdentifier{2, 5, 4, 7}
	stateOID              = asn1.ObjectIdentifier{2, 5, 4, 8}
	organizationOID       = asn1.ObjectIdentifier{2, 5, 4, 10}
	organizationalUnitOID = asn1.ObjectIdentifier{2, 5, 4, 11}
)

const certificateError = "Invalid certificate in file"

// CSRInfo is Certificate Signing Request (CSR) Information
type CSRInfo struct {
	CN           string        `json:"CN"`
	Names        []csr.Name    `json:"names,omitempty"`
	Hosts        []string      `json:"hosts,omitempty"`
	KeyRequest   *KeyRequest   `json:"key,omitempty"`
	CA           *csr.CAConfig `json:"ca,omitempty" hide:"true"`
	SerialNumber string        `json:"serial_number,omitempty"`
}

// KeyRequest encapsulates size and algorithm for the key to be generated.
// If ReuseKey is set, reenrollment requests will reuse the existing private
// key.
type KeyRequest struct {
	Algo     string `json:"algo" yaml:"algo" help:"Specify key algorithm"`
	Size     int    `json:"size" yaml:"size" help:"Specify key size"`
	ReuseKey bool   `json:"reusekey" yaml:"reusekey" help:"Reuse existing key during reenrollment"`
}

// Attribute is a name and value pair
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	ECert bool   `json:"ecert,omitempty"`
}

// GetName returns the name of the attribute
func (a *Attribute) GetName() string {
	return a.Name
}

// GetValue returns the value of the attribute
func (a *Attribute) GetValue() string {
	return a.Value
}

// AttributeRequest is a request for an attribute.
// This implements the certmgr/AttributeRequest interface.
type AttributeRequest struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional,omitempty"`
}

// GetName returns the name of an attribute being requested
func (ar *AttributeRequest) GetName() string {
	return ar.Name
}

// IsRequired returns true if the attribute being requested is required
func (ar *AttributeRequest) IsRequired() bool {
	return !ar.Optional
}

type User struct {
	Name                      string
	Pass                      string `mask:"password"`
	Type                      string
	Affiliation               string
	Attributes                []Attribute
	State                     int
	MaxEnrollments            int
	Level                     int
	IncorrectPasswordAttempts int
	attrs                     map[string]Attribute
}
