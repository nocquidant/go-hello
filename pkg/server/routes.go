package server

// One place for all routes
func (s *server) routes() {
	s.router.HandleFunc("/health", s.handleHealth()).Methods("GET")
	s.router.HandleFunc("/hello", s.handleHello()).Methods("GET")
	s.router.HandleFunc("/gurl", s.handleGetURL()).Methods("GET")
}
