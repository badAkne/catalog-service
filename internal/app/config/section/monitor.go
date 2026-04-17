package section

type Monitor struct {
	LogLevel          string `default:"debug" split_words:"true"`
	Environment       string `default:"development"`
	MonitorPrometheus `default:"true" split_words:"true"`
}

type MonitorPrometheus struct {
	Enabled bool
}
