package broker

import "errors"

const (
	SASLPlaintextSecurityProtocol = "sasl_plaintext"
	SASLSSLSecurityProtocol       = "sasl_ssl"

	SASLMechanism = "PLAIN"
)

// SASL the based authentication with broker.
// While there are multiple SASL authentication methods the current implementation is limited
// to plaintext authentication.
type SASL struct {
	// Username is the authentication identity to present for plaintext authentication.
	Username string `yaml:"username"`

	// Password for plaintext authentication.
	Password string `yaml:"password"`
}

// SSL represents ssl connecting to the broker.
type SSL struct {
	// CALocation a file or directory path to CA certificate for verifying the broker's key.
	CALocation string `yaml:"ca-location"`
}

// NetworkConfig specifies the configuration needed to create a connection to the broker
type NetworkConfig struct {
	// BrokerAddresses initial list of brokers as a list of broker host in host:port format.
	BrokerAddresses []string `yaml:"brokers"`

	*SASL `yaml:"sasl"`

	*SSL `yaml:"ssl,omitempty"`
}

func (c *NetworkConfig) Init() error {
	if len(c.BrokerAddresses) == 0 {
		return errors.New("broker addresses cannot be empty")
	}

	if c.SASL == nil {
		return errors.New("sasl configuration cannot be empty")
	}

	if len(c.SASL.Username) == 0 {
		return errors.New("sasl username cannot be empty")
	}

	if len(c.SASL.Password) == 0 {
		return errors.New("sasl password cannot be empty")
	}

	if c.SSL != nil {
		if len(c.SSL.CALocation) == 0 {
			return errors.New("ssl ca certificate location cannot be empty")
		}
	}

	return nil
}
