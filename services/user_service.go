package services

// type UserService struct {
// 	users    map[string]*models.User
// 	filePath string
// 	mu       sync.RWMutex
// }

// func NewUserService() *UserService {
// 	return &UserService{}
// }

// func (s *UserService) loadUsers() error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	// 检查文件是否存在
// 	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
// 		return nil
// 	}

// 	data, err := ioutil.ReadFile(s.filePath)
// 	if err != nil {
// 		return err
// 	}

// 	var users []*models.User
// 	if err := json.Unmarshal(data, &users); err != nil {
// 		return err
// 	}

// 	for _, user := range users {
// 		s.users[user.Username] = user
// 	}

// 	return nil
// }

// func (s *UserService) SaveUsers() error {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	users := make([]*models.User, 0, len(s.users))
// 	for _, user := range s.users {
// 		users = append(users, user)
// 	}

// 	data, err := json.MarshalIndent(users, "", "    ")
// 	if err != nil {
// 		return err
// 	}

// 	return ioutil.WriteFile(s.filePath, data, 0644)
// }

// func (s *UserService) GetUser(username string) *models.User {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	return s.users[username]
// }

// func (s *UserService) GetAllUsers() []*models.User {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()

// 	users := make([]*models.User, 0, len(s.users))
// 	for _, user := range s.users {
// 		users = append(users, user)
// 	}
// 	return users
// }

// func (s *UserService) AddUser(user *models.User) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, exists := s.users[user.Username]; exists {
// 		return ErrUserExists
// 	}

// 	s.users[user.Username] = user
// 	return s.SaveUsers()
// }

// func (s *UserService) UpdateUser(user *models.User) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, exists := s.users[user.Username]; !exists {
// 		return ErrUserNotFound
// 	}

// 	s.users[user.Username] = user
// 	return s.SaveUsers()
// }

// func (s *UserService) DeleteUser(username string) error {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, exists := s.users[username]; !exists {
// 		return ErrUserNotFound
// 	}

// 	delete(s.users, username)
// 	return s.SaveUsers()
// }
