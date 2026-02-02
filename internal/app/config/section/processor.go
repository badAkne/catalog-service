package section

type (
	Processor struct {
		WebServer ProcessorWebServer `split_words:"true"`
	}

	ProcessorWebServer struct {
		ListenPort string `default:"8080" split_words:"true"`
	}
)
