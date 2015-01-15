package main

type config struct {
	flag_debug	bool
	flag_list	string
	flag_run	string
	max_checks	int
}

func default_config() config {
	cfg := config{
		flag_debug: false,
	}
	return cfg
}
