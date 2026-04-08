package section

import "time"

type (
	Processor struct {
		WebServer ProcessorWebServer `split_words:"true"`
		Grpc      ProcessorGrpc
		Gateway   ProcessorGateway
	}

	ProcessorWebServer struct {
		ListenPort        uint32        `default:"8080" split_words:"true"`
		ReadHeaderTimeout time.Duration `default:"3s" split_words:"true"`
	}

	ProcessorGrpc struct {
		Host       string `default:"localhost" split_words:"true"`
		ListenPort uint32 `default:"50052" split_words:"true"`
	}

	ProcessorGateway struct {
		Host         string `default:"localhost" split_words:"true"`
		ListenPort   uint32 `default:"8081" split_words:"true"`
		GrpcEndpoint string `default:"localhost:8081" split_words:"true"`
	}
)
