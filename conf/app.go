package conf

var rootDir string

func SetAppRootDir(dir string) {
	rootDir = dir
}

func GetAppRootDir() string {
	return rootDir
}
