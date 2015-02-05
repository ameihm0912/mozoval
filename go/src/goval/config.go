package main

type config struct {
	flagDebug bool
	flagList  string
	flagRun   string
	maxChecks int
}

func defaultConfig() config {
	return config{
		flagDebug: false,
	}
}
