package app

type App struct {
	config *config.Config
	db     *database.Database
}

func NewApp(*App, error) {
	cfg, err := config.LoadConfig()

}
