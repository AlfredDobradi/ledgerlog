package server

func WithAddress(addr string) Option {
	return func(s *Service) error {
		s.Addr = addr
		return nil
	}
}
