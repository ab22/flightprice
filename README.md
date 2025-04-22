# Flight Price

Go microservice that retrieves flight prices from local mocked APIs.


## Quickstart

```sh
# Create our .env file.
cp .env.example .env

# This will build our backend and launch it with docker compose.
make build && make up

# Perform a quick test
curl -i http://localhost:8080/ping

# or make use of helper recipes already on backend/Makefile.
make ping

# perform a request to /flights/search without a valid JWT token. This should
# show an http.StatusNotFound (404) instead of a forbidden. We show a 404 and not a
# 403 for security purposes.
make flights-invalid

# create a sample JWT token by calling our API.
make login

# Copy the JWT token into our `.env` file inside of the API_TOKEN variable. Now
# we should be able to hit our flights endpoint.
make flights

# Test the websocket server. The current recipe sends a path parameter
# so that it receives an update every 3 seconds. This can be changed to
# /subscribe/x where x is the number of seconds on which you will receive the
# next message.
make ws
```

## Configuring HTTP/TLS

### Generating a Certificate from a Trusted Source

We would need to generate a certificate from a trusted Certificate Authority such
as DigiCert, GoDaddy, Let's Encrypt, ZeroSSL and many others.

We should also make sure that the certificate is configured with the specific TLS version
that we want to support such as TLS v1.2 or v1.3, as well as disable older versions which
are less secure. As a rule of thumb, we need to ensure that the certificate is taking all
hostnames.

### Deploying and Loading a Certificate

Once we have a certificate, we need to securely deploy it. There are many ways on securing
a certificate and it also depends on our infrastucture stack, so, one example would be
using AWS' Secrets Manager service which can securely store secrets.

Once we have it deployed, we would need to load it on our golang application by using AWS'
Golang SDK which provides a secretmanager package. Here's an example on how we should load
our secret:

```go
import "github.com/aws/aws-sdk-go/service/secretsmanager"

func getCertFromAWS() (cert, key []byte) {
	svc := secretsmanager.New(awsSession)
	secret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String("jobsity/flightsservice/cert"),
	})

	// TODO: Handle error correctly.

	// parseCert should extract the TLS Cert and Key which are later used when crreating
	// our http.Server.
	return parseCert(secret.SecretString)
}

func main() {
	certPEM, keyPEM := getCertFromAWS()
	cert, _ := tls.X509KeyPair(certPEM, keyPEM)

	server.TLSConfig.Certificates = []tls.Certificate{cert}
}
````

### Configuring our Golang app to handle secure traffic

Once we have created, deployed and configured our server to use the certificate, we need
to make a few changes on our Go code so that we correctly handle HTTPS traffic.

First, we need to change how we initialize our `http.Server{}` by changing a few
configuration options:
1. Swap from usage of `server.ListenAndServe` to `server.ListenAndServeTLS` and
1. Change our port from `:80` to `:443`.
1. Specify which cipher suites we support.
1. Set HTTP Strict Transport Security header (HSTS).
1. Redirect all traffic from port 80 to 443
1. (Optional) Specify curve preferences for the ECDH key exchange process. We could use
   these options to optimize for performance or security.

An example configuration would look like:

```go
server := &http.Server{
    Addr:    ":443",
    Handler: router,
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
        CurvePreferences: []tls.CurveID{
        		// Prefer X25519, then P256 and finally P384.
            tls.X25519,
            tls.CurveP256,
            tls.CurveP384,
        },
        PreferServerCipherSuites: true,
        Certificates: []tls.Certificate{cert},
    },
}

log.Fatal(server.ListenAndServeTLS("fullchain.pem", "privkey.pem"))
```

### Other considerations

1. A Certificate Rotation process must be defined to avoid serving expired certificates.
  This can lead to user's browers not being able to validate the authenticity of the
  suites and therefore the users/clients will not be able to securely navigate or consume
  our site/api.
1. Certificates with shorter lifetimes are preferred because they reduce risks, enforces
  automation and aligns with modern security practices.
