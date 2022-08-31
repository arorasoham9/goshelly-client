package basic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

func GoGenCert(){
	
	key, err := rsa.GenerateKey(rand.Reader, 4096)
    if err != nil {
		fmt.Println("Could not generate private key: ",err)
        os.Exit(1)
    }
    keyBytes := x509.MarshalPKCS1PrivateKey(key)
    // PEM encoding of private key
    keyPEM := pem.EncodeToMemory(
        &pem.Block{
            Type:  "RSA PRIVATE KEY",
            Bytes: keyBytes,
        },
    )
    // fmt.Println(string(keyPEM))

    notBefore := time.Now()
    notAfter := notBefore.Add(365*24*10*time.Hour)

    //Create certificate templet
    template := x509.Certificate{
        SerialNumber:          big.NewInt(0),
        Subject:               pkix.Name{CommonName: "localhost"},
        SignatureAlgorithm:    x509.SHA256WithRSA,
        NotBefore:             notBefore,
        NotAfter:              notAfter,
        BasicConstraintsValid: true,
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
    }
    //Create certificate using template
    derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
    if err != nil {
		fmt.Println("Could not generate SSL Certificates: ",err)
        os.Exit(1)
    }
    //pem encoding of certificate
    certPem := string(pem.EncodeToMemory(
        &pem.Block{
            Type:  "CERTIFICATE",
            Bytes: derBytes,
        },
    ))
    // fmt.Println(certPem)
	err = os.MkdirAll("./goshelly_certs/", os.ModePerm)
	f,err :=os.Create("./goshelly_certs/client.key")
	if err != nil {
        fmt.Println("Could not open key.", err)
		os.Exit(1)
    }
	f.Close()
	f,err = os.Create("./goshelly_certs/client.pem")
	if err != nil {
        fmt.Println("Could not open pem file.", err)
		os.Exit(1)
    }
	f.Close()
	writeBytesToFile("./goshelly_certs/client.key", keyPEM)
	writeBytesToFile("./goshelly_certs/client.pem", []byte(certPem))
}
func writeBytesToFile(path string, b []byte){
	err := ioutil.WriteFile(path, b, 0644)
    if err != nil {
        fmt.Println("Could not write key or pem to file.", err)
		os.Exit(1)
    }
}