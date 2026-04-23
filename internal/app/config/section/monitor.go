package section

type Monitor struct {
	LogLevel          string `default:"debug" split_words:"true"`
	Environment       string `default:"development"`
	MonitorPrometheus `default:"true" split_words:"true"`
	MonitorSentry     `split_word:"true"`
}

type (
	MonitorPrometheus struct {
		Enabled bool
	}

	MonitorSentry struct {
		Enabled bool `default:"false"`
		DSN     string
	}
)
