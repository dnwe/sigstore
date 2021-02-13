package x509pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"

	"github.com/projectrekor/signer/config"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func GenPrivKeyPEM() (*rsa.PrivateKey, error) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// TODO: The followinfg subj values should be gathered from
// a developers profile (likely something in ~/.config)
func GenerateCsr(config config.Config, keyBytes interface{}) ([]byte, error) {
	emailAddress := config.Email
    subj := pkix.Name{
        CommonName:          config.CommonName,
        Country:            []string{config.Country},
        Province:           []string{config.Province},
        Locality:           []string{config.Locality},
        Organization:       []string{config.Organization},
        OrganizationalUnit: []string{config.OrganizationalUnit},
        ExtraNames: []pkix.AttributeTypeAndValue{
            {
                Type:  oidEmailAddress,
                Value: asn1.RawValue{
                    Tag:   asn1.TagIA5String,
                    Bytes: []byte(emailAddress),
                },
            },
        },
    }

    template := x509.CertificateRequest{
        Subject:            subj,
        SignatureAlgorithm: x509.SHA256WithRSA,
    }

    csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}), nil

}
