package certificate

import (
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "net"
    "net/url"
    "time"
)

type CertificateDetails struct {
    Subject        string
    Issuer         pkix.Name
    CommonName     string
    Issued         string
    Expiration     string
    DNSNames       []string
    IPAddresses    []net.IP
    EmailAddresses []string
    URIs           []*url.URL
}

func BuildCertificateDetails(certPath string) CertificateDetails {
    certFile, err := ioutil.ReadFile(certPath)
    if err != nil {
        fmt.Println(err)
    }

    pemBlock, _ := pem.Decode([]byte(certFile))
    if pemBlock == nil {
        fmt.Println(err)
    }

    cert, err := x509.ParseCertificate(pemBlock.Bytes)
    if err != nil {
        fmt.Println(err)
    }

    subject := cert.Subject

    return CertificateDetails{
        Subject:        subject.String(),
        Issuer:         cert.Issuer,
        CommonName:     subject.CommonName,
        Issued:         cert.NotBefore.Format(time.RFC3339),
        Expiration:     cert.NotAfter.Format(time.RFC3339),
        DNSNames:       cert.DNSNames,
        IPAddresses:    cert.IPAddresses,
        EmailAddresses: cert.EmailAddresses,
        URIs:           cert.URIs,
    }
}
