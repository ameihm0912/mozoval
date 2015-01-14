package main

type config struct {
	flag_debug	bool
	flag_list	string
	flag_run	string
}

func default_config() config {
	cfg := config{
		flag_debug: false,
	}
	return cfg
}
