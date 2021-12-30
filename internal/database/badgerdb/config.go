package badgerdb

var (
	databasePath string = "./tmp"
	valuePath    string = ""
)

func DatabasePath() string {
	return databasePath
}

func ValuePath() string {
	return valuePath
}

func SetDatabasePath(newPath string) {
	databasePath = newPath
}

func SetValuePath(newPath string) {
	valuePath = newPath
}
