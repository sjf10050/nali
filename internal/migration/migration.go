package migration

// Run applies pending config/data directory migrations. It must be called at
// startup after constant.InitPaths, since the migrations read the resolved
// config/data directories.
func Run() {
	migration2v6()
	migration2v7()
}
