package adapter

type BrokerConnConfig struct {
	MEMPHIS_HOST         string
	ROOT_USER            string
	ROOT_PASSWORD        string
	CONNECTION_TOKEN     string
	CLIENT_CERT_PATH     string
	CLIENT_KEY_PATH      string
	ROOT_CA_PATH         string
	USER_PASS_BASED_AUTH bool
	DEBUG                bool
	CLOUD_ENV            bool
	DEV_ENV              bool
}
