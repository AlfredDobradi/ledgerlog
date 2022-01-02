package cockroach

var (
	user           string = "root"
	password       string = ""
	host           string = "127.0.0.1"
	port           string = "26256"
	database       string = "defaultdb"
	sslMode        string = "disabled"
	sslRootCert    string = ""
	cluster        string = ""
	minConnections int32  = 2
	maxConnections int32  = 4
)

func User() string {
	return user
}

func Password() string {
	return password
}

func Host() string {
	return host
}

func Port() string {
	return port
}

func Database() string {
	return database
}

func SSLMode() string {
	return sslMode
}

func SSLRootCert() string {
	return sslRootCert
}

func Cluster() string {
	return cluster
}

func MinConnections() int32 {
	return minConnections
}

func MaxConnections() int32 {
	return maxConnections
}

func SetUser(newUser string) {
	user = newUser
}

func SetPassword(newPassword string) {
	password = newPassword
}

func SetHost(newHost string) {
	host = newHost
}

func SetPort(newPort string) {
	port = newPort
}

func SetDatabase(newDatabase string) {
	database = newDatabase
}

func SetSSLMode(newSSLMode string) {
	sslMode = newSSLMode
}

func SetSSLRootCert(newSSLRootCert string) {
	sslRootCert = newSSLRootCert
}

func SetCluster(newCluster string) {
	cluster = newCluster
}

func SetMinConnections(newMin int32) {
	if newMin < 1 {
		newMin = 1
	} else if newMin > maxConnections {
		newMin = maxConnections
	}
	minConnections = newMin
}

func SetMaxConnections(newMax int32) {
	if newMax < 1 {
		newMax = 1
	} else if newMax < minConnections {
		newMax = minConnections
	}
	maxConnections = newMax
}
