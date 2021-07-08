package ca

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/isnlan/coral/pkg/hlf/ca/attrmgr"

	"github.com/hyperledger/fabric/bccsp/factory"

	"github.com/cloudflare/cfssl/csr"

	"encoding/asn1"

	"github.com/cloudflare/cfssl/signer"

	"github.com/cloudflare/cfssl/config"

	"github.com/cloudflare/cfssl/log"
	cflocalsigner "github.com/cloudflare/cfssl/signer/local"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/isnlan/coral/pkg/errors"
)

type CA struct {
	csp             bccsp.BCCSP
	name            string
	csr             *CSRInfo
	parentServerURL string
	signing         *config.Signing
	enrollSigner    signer.Signer
	cert            []byte
	key             []byte
	attrMgr         *attrmgr.Mgr
}

type Config struct {
	ParentServerURL string
	Cert            []byte
	Key             []byte
}

func New(cfg *Config) (*CA, error) {
	signing := &config.Signing{
		Profiles: map[string]*config.SigningProfile{},
		Default:  config.DefaultConfig(),
	}
	ca := &CA{
		csp:             factory.GetDefault(),
		name:            "ca",
		csr:             &CSRInfo{},
		parentServerURL: cfg.ParentServerURL,
		signing:         signing,
		cert:            cfg.Cert,
		key:             cfg.Key,
		attrMgr:         attrmgr.New(),
		enrollSigner:    nil,
	}

	err := ca.validateCertAndKey(ca.cert, ca.key)
	if err != nil {
		return nil, err
	}

	ca.csr.CN, err = ca.loadCNFromEnrollmentInfo(ca.cert)
	if err != nil {
		return nil, err
	}

	//idemixConfig := &idemix.Config{
	//	IssuerPublicKeyfile:      "",
	//	IssuerSecretKeyfile:      "",
	//	RevocationPublicKeyfile:  "",
	//	RevocationPrivateKeyfile: "",
	//	RHPoolSize:               1000,
	//	NonceExpiration:          "15s",
	//	NonceSweepInterval:       "15m",
	//}
	//issuer := idemix.NewIssuer(ca.name, "/Users/snlan/go/src/github.com/isnlan/coral/test",
	//	idemixConfig, ca.csp, idemix.NewLib())
	//issuer.Init()

	//if nm, ok := issuer.(interface{ NonceManager() idemix.NonceManager }); ok && nm != nil {
	//	if ss, ok := nm.NonceManager().(interface{ StartNonceSweeper() }); ok {
	//		ss.StartNonceSweeper()
	//	}
	//}

	//ca.issuer = issuer

	ca.enrollSigner, err = BccspBackedSigner(ca.cert, ca.key, ca.signing, ca.csp)
	if err != nil {
		return nil, err
	}

	log.Debug("CA initialization successful")
	return ca, nil
}

func (ca *CA) Enroll(req signer.SignRequest, id string, pass string) ([]byte, error) {
	notBefore, notAfter, err := ca.getCACertExpiry()
	if err != nil {
		return nil, errors.New("Failed to get CA certificate information")
	}

	if !notAfter.IsZero() && req.NotAfter.After(notAfter) {
		log.Debugf("Requested expiry '%s' is after the CA certificate expiry '%s'. Will use CA cert expiry",
			req.NotAfter, notAfter)
		req.NotAfter = notAfter
	}
	// Make sure that requested expiration for enrollment certificate is not before CA certificate
	// expiration
	if !notBefore.IsZero() && req.NotBefore.Before(notBefore) {
		log.Debugf("Requested expiry '%s' is before the CA certificate expiry '%s'. Will use CA cert expiry",
			req.NotBefore, notBefore)
		req.NotBefore = notBefore
	}

	user := &User{
		Name:                      id,
		Pass:                      pass,
		Type:                      "client",
		Affiliation:               "",
		Attributes:                nil,
		State:                     0,
		MaxEnrollments:            0,
		Level:                     0,
		IncorrectPasswordAttempts: 0,
		attrs:                     nil,
	}
	err = processSignRequest(id, &req, ca, user)
	if err != nil {
		return nil, err
	}

	cert, err := ca.enrollSigner.Sign(req)
	if err != nil {
		return nil, errors.WithMessage(err, "Certificate signing failure")
	}

	//resp := &util2.EnrollmentResponseNet{
	//	Cert: util.B64Encode(cert),
	//}
	//
	//err = ca.fillCAInfo(&resp.ServerInfo)
	//if err != nil {
	//	return nil, err
	//}
	return cert, nil
}

// Process the sign request.
// Make any authorization checks needed, depending on the contents
// of the CSR (Certificate Signing Request).
// In particular, if the request is for an intermediate CA certificate,
// the caller must have the "hf.IntermediateCA" attribute.
// Check to see that CSR values do not exceed the character limit
// as specified in RFC 3280, page 103.
// Set the OU fields of the request.
func processSignRequest(id string, req *signer.SignRequest, ca *CA, user *User) error {
	// Decode and parse the request into a CSR so we can make checks
	block, _ := pem.Decode([]byte(req.Request))
	if block == nil {
		return errors.New("CSR Decode failed")
	}
	if block.Type != "NEW CERTIFICATE REQUEST" && block.Type != "CERTIFICATE REQUEST" {
		return errors.New("not a certificate or csr")
	}
	csrReq, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return err
	}
	log.Debugf("Processing sign request: id=%s, CommonName=%s, Subject=%+v", id, csrReq.Subject.CommonName, req.Subject)
	if (req.Subject != nil && req.Subject.CN != id) || csrReq.Subject.CommonName != id {
		return errors.New("The CSR subject common name must equal the enrollment ID")
	}
	isForCACert, err := isRequestForCASigningCert(csrReq, ca, req.Profile)
	if err != nil {
		return err
	}
	if isForCACert {
		// This is a request for a CA certificate, so make sure the caller
		// has the 'hf.IntermediateCA' attribute
		err := ca.attributeIsTrue(id, "hf.IntermediateCA")
		if err != nil {
			return errors.Wrap(err, "Enrolled failed")
		}
	}
	// Check the CSR input length
	err = csrInputLengthCheck(csrReq)
	if err != nil {
		return errors.Wrap(err, "CSR input validation failed")
	}

	// Set the OUs in the request appropriately.
	setRequestOUs(req, user)
	log.Debug("Finished processing sign request")
	return nil
}

func (ca *CA) attributeIsTrue(username, attrname string) error {
	return nil
}

//func (ca *CA) fillCAInfo(info *util2.CAInfoResponseNet) error {
//	caChain, err := util.ReadFile(ca.certFile)
//	if err != nil {
//		return err
//	}
//
//	info.CAName = ca.name
//	info.CAChain = util.B64Encode(caChain)
//
//	ipkBytes, err := ca.issuer.IssuerPublicKey()
//	if err != nil {
//		return err
//	}
//	rpkBytes, err := ca.issuer.RevocationPublicKey()
//	if err != nil {
//		return err
//	}
//	info.IssuerPublicKey = util.B64Encode(ipkBytes)
//	info.IssuerRevocationPublicKey = util.B64Encode(rpkBytes)
//	return nil
//}

// Check to see if this is a request for a CA signing certificate.
// This can occur if the profile or the CSR has the IsCA bit set.
// See the X.509 BasicConstraints extension (RFC 5280, 4.2.1.9).
func isRequestForCASigningCert(csrReq *x509.CertificateRequest, ca *CA, profile string) (bool, error) {
	// Check the profile to see if the IsCA bit is set
	sp := getSigningProfile(ca, profile)
	if sp == nil {
		return false, errors.Errorf("Invalid profile: '%s'", profile)
	}
	if sp.CAConstraint.IsCA {
		log.Debugf("Request is for a CA signing certificate as set in profile '%s'", profile)
		return true, nil
	}
	// Check the CSR to see if the IsCA bit is set
	for _, val := range csrReq.Extensions {
		if val.Id.Equal(basicConstraintsOID) {
			var constraints csr.BasicConstraints
			var rest []byte
			var err error
			if rest, err = asn1.Unmarshal(val.Value, &constraints); err != nil {
				return false, errors.Wrap(err, "Failed parsing CSR constraints")
			} else if len(rest) != 0 {
				return false, errors.New("Trailing data after X.509 BasicConstraints")
			}
			if constraints.IsCA {
				log.Debug("Request is for a CA signing certificate as indicated in the CSR")
				return true, nil
			}
		}
	}
	// The IsCA bit was not set
	log.Debug("Request is not for a CA signing certificate")
	return false, nil
}

func getSigningProfile(ca *CA, profile string) *config.SigningProfile {
	return ca.signing.Default
}

// Checks to make sure that character limits are not exceeded for CSR fields
func csrInputLengthCheck(req *x509.CertificateRequest) error {
	log.Debug("Checking CSR fields to make sure that they do not exceed maximum character limits")

	for _, n := range req.Subject.Names {
		value := n.Value.(string)
		switch {
		case n.Type.Equal(commonNameOID):
			if len(value) > commonNameLength {
				return errors.Errorf("The CN '%s' exceeds the maximum character limit of %d", value, commonNameLength)
			}
		case n.Type.Equal(serialNumberOID):
			if len(value) > serialNumberLength {
				return errors.Errorf("The serial number '%s' exceeds the maximum character limit of %d", value, serialNumberLength)
			}
		case n.Type.Equal(organizationalUnitOID):
			if len(value) > organizationalUnitNameLength {
				return errors.Errorf("The organizational unit name '%s' exceeds the maximum character limit of %d", value, organizationalUnitNameLength)
			}
		case n.Type.Equal(organizationOID):
			if len(value) > organizationNameLength {
				return errors.Errorf("The organization name '%s' exceeds the maximum character limit of %d", value, organizationNameLength)
			}
		case n.Type.Equal(countryOID):
			if len(value) > countryNameLength {
				return errors.Errorf("The country name '%s' exceeds the maximum character limit of %d", value, countryNameLength)
			}
		case n.Type.Equal(localityOID):
			if len(value) > localityNameLength {
				return errors.Errorf("The locality name '%s' exceeds the maximum character limit of %d", value, localityNameLength)
			}
		case n.Type.Equal(stateOID):
			if len(value) > stateOrProvinceNameLength {
				return errors.Errorf("The state name '%s' exceeds the maximum character limit of %d", value, stateOrProvinceNameLength)
			}
		}
	}

	return nil
}

// Set the OU fields of the sign request based on the identity's type and affilation.
// For example, if the type is 'peer' and the affiliation is 'a.b.c', the
// OUs become 'OU=c,OU=b,OU=a,OU=peer'.
// This is necessary because authorization decisions are made based on the OU fields,
// so we ignore any OU values specified in the enroll request and set them according
// to the type and affiliation.
func setRequestOUs(req *signer.SignRequest, user *User) {
	s := req.Subject
	if s == nil {
		s = &signer.Subject{}
	}
	names := []csr.Name{}
	// Add non-OU fields from request
	for _, name := range s.Names {
		if name.C != "" || name.L != "" || name.O != "" || name.ST != "" || name.SerialNumber != "" {
			name.OU = ""
			names = append(names, name)
		}
	}
	// Add an OU field with the type
	names = append(names, csr.Name{OU: user.Type})
	for _, aff := range strings.Split(user.Affiliation, ".") {
		names = append(names, csr.Name{OU: aff})
	}
	// Replace with new names
	s.Names = names
	req.Subject = s
}

func (ca *CA) validateCertAndKey(certPEM []byte, keyPEM []byte) error {
	log.Debug("Validating the CA certificate and key")

	cert, err := GetX509CertificateFromPEM(certPEM)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}

	if err = validateDates(cert); err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}
	if err = validateUsage(cert, ca.name); err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}
	if err = validateIsCA(cert); err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}
	if err = validateKeyType(cert); err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}
	if err = validateKeySize(cert); err != nil {
		return errors.WithMessage(err, fmt.Sprintf(certificateError+" '%s'", certPEM))
	}
	if err = validateMatchingKeys(cert, keyPEM); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Invalid certificate and/or key in files '%s' and '%s'", certPEM, keyPEM))
	}
	log.Debug("Validation of CA certificate and key successful")

	return nil
}

func (ca *CA) loadCNFromEnrollmentInfo(certPEM []byte) (string, error) {
	log.Debug("Loading CN from existing enrollment information")
	name, err := GetEnrollmentIDFromPEM(certPEM)
	if err != nil {
		return "", err
	}
	return name, nil
}

// Initialize the enrollment signer
func (ca *CA) initEnrollmentSigner() (err error) {
	log.Debug("Initializing enrollment signer")

	// If there is a config, use its signing policy. Otherwise create a default policy.
	var policy *config.Signing
	if ca.signing != nil {
		policy = ca.signing
	} else {
		policy = &config.Signing{
			Profiles: map[string]*config.SigningProfile{},
			Default:  config.DefaultConfig(),
		}
		policy.Default.CAConstraint.IsCA = true
	}

	// Make sure the policy reflects the new remote
	if ca.parentServerURL != "" {
		err = policy.OverrideRemotes(ca.parentServerURL)
		if err != nil {
			return errors.Wrap(err, "Failed initializing enrollment signer")
		}
	}

	ca.enrollSigner, err = BccspBackedSigner(ca.cert, ca.key, policy, ca.csp)
	if err != nil {
		return err
	}
	//ca.enrollSigner.SetDBAccessor(&certDBAccessor{})

	// Successful enrollment
	return nil
}

func (ca *CA) getCACertExpiry() (time.Time, time.Time, error) {
	var notAfter time.Time
	var notBefore time.Time
	signer, ok := ca.enrollSigner.(*cflocalsigner.Signer)
	if ok {
		cacert, err := signer.Certificate("", "ca")
		if err != nil {
			log.Errorf("Failed to get CA certificate for CA %s: %s", ca.name, err)
			return notBefore, notAfter, err
		} else if cacert != nil {
			notAfter = cacert.NotAfter
			notBefore = cacert.NotBefore
		}
	} else {
		log.Errorf("Not expected condition as the enrollSigner can only be cfssl/signer/local/Signer")
		return notBefore, notAfter, errors.New("Unexpected error while getting CA certificate expiration")
	}
	return notBefore, notAfter, nil
}

func canSignCRL(cert *x509.Certificate) bool {
	return cert.KeyUsage&x509.KeyUsageCRLSign != 0
}

func validateDates(cert *x509.Certificate) error {
	log.Debug("Check CA certificate for valid dates")

	notAfter := cert.NotAfter
	currentTime := time.Now().UTC()

	if currentTime.After(notAfter) {
		return errors.New("Certificate provided has expired")
	}

	notBefore := cert.NotBefore
	if currentTime.Before(notBefore) {
		return errors.New("Certificate provided not valid until later date")
	}

	return nil
}

func validateUsage(cert *x509.Certificate, caName string) error {
	log.Debug("Check CA certificate for valid usages")

	if cert.KeyUsage == 0 {
		return errors.New("No usage specified for certificate")
	}

	if cert.KeyUsage&x509.KeyUsageCertSign == 0 {
		return errors.New("The 'cert sign' key usage is required")
	}

	if !canSignCRL(cert) {
		log.Warningf("The CA certificate for the CA '%s' does not have 'crl sign' key usage, so the CA will not be able generate a CRL", caName)
	}
	return nil
}

func validateIsCA(cert *x509.Certificate) error {
	log.Debug("Check CA certificate for valid IsCA value")

	if !cert.IsCA {
		return errors.New("Certificate not configured to be used for CA")
	}

	return nil
}

func validateKeyType(cert *x509.Certificate) error {
	log.Debug("Check that key type is supported")

	switch cert.PublicKey.(type) {
	case *dsa.PublicKey:
		return errors.New("Unsupported key type: DSA")
	}

	return nil
}

func validateKeySize(cert *x509.Certificate) error {
	log.Debug("Check that key size is of appropriate length")

	switch cert.PublicKey.(type) {
	case *rsa.PublicKey:
		size := cert.PublicKey.(*rsa.PublicKey).N.BitLen()
		if size < 2048 {
			return errors.New("Key size is less than 2048 bits")
		}
	}

	return nil
}

func validateMatchingKeys(cert *x509.Certificate, keyPEM []byte) error {
	log.Debug("Check that public key and private key match")

	pubKey := cert.PublicKey
	switch pubKey.(type) {
	case *rsa.PublicKey:
		privKey, err := GetRSAPrivateKey(keyPEM)
		if err != nil {
			return err
		}

		if privKey.PublicKey.N.Cmp(pubKey.(*rsa.PublicKey).N) != 0 {
			return errors.New("Public key and private key do not match")
		}
	case *ecdsa.PublicKey:
		privKey, err := GetECPrivateKey(keyPEM)
		if err != nil {
			return err
		}

		if privKey.PublicKey.X.Cmp(pubKey.(*ecdsa.PublicKey).X) != 0 {
			return errors.New("Public key and private key do not match")
		}
	}

	return nil
}
