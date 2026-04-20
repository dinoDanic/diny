package auth

func Logout(s *Session) {
	if s == nil {
		return
	}
	s.Token = ""
}
