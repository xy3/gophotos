package photos

type config struct {
	SqliteDSN string
	Host      string
	BasePath  string
}

var (
	Config = config{
		SqliteDSN: "file:photos.sqlite",
		Host:      "localhost:8090",
		BasePath:  "/api/v1",
	}
)
