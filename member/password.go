package member

//PasswordProvider  member password provider interface
type PasswordProvider interface {
	VerifyPassword(uid string, password string) (bool, error)
	//PasswordChangeable return password changeable
	PasswordChangeable() bool
	//UpdatePassword update user password
	//Return any error if raised
	UpdatePassword(uid string, password string) error
}

//ServicePassword Member password module.
type ServicePassword struct {
	service *Service
}

//UpdatePassword update user password
//Return any error if raised
func (s *ServicePassword) UpdatePassword(uid string, password string) error {
	return s.service.PasswordProvider.UpdatePassword(uid, password)
}

//VerifyPassword Verify user password.
//Return verify result and any error if raised
func (s *ServicePassword) VerifyPassword(uid string, password string) (bool, error) {
	result, err := s.service.PasswordProvider.VerifyPassword(uid, password)
	if !result || err != nil {
		return result, err
	}
	if s.service.StatusProvider != nil {
		statusStore := NewStatusStore()
		err := s.service.Status().Load(statusStore, uid)
		if err != nil {
			return false, err
		}
		if !IsAvaliable(statusStore.Get(uid)) {
			return false, ErrUserBanned
		}
	}
	return true, nil
}

//PasswordChangeable return password changeable
func (s *ServicePassword) PasswordChangeable() bool {
	return s.service.PasswordProvider.PasswordChangeable()
}
