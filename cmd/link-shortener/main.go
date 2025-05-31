package main

// func parseConfig() (config.Params, error) {
// 	cfg := config.Params{}

// 	if err := env.Parse(cfg); err != nil {
// 		return config.Params{}, fmt.Errorf(
// 			"failed to parse config for environment variables: %w", err,
// 		)
// 	}

// 	return cfg, nil
// }

// func main() {
// 	cfg, err := parseConfig()
// 	if err != nil {
// 		log.Default().Fatal(err)
// 	}

// 	linkShortener, err := instance.New(cfg)
// 	if err != nil {
// 		log.Default().Fatal(err)
// 	}

// 	if err := linkShortener.Run(context.Background()); err != nil {
// 		os.Exit(1)
// 	}
// }
