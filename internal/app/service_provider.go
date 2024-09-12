package app

type serviceProvider struct {
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

//func (s *serviceProvider) Config() config.Config {
//	if s.config == nil {
//		s.config = &config.Config{
//			App:  config.App(),
//			Grpc: config.Grpc(),
//			Http: config.Http(),
//			Pg:   config.Pg(),
//		}
//	}
//
//	return s.config
//}
