package cockroach

var (
	user        string = "root"
	password    string = ""
	host        string = "127.0.0.1"
	port        string = "26256"
	database    string = "defaultdb"
	sslMode     string = "disabled"
	sslRootCert string = ""
	cluster     string = ""
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
