package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// User represents a user in the system
type User struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Active    bool       `json:"active"`
	Created   time.Time  `json:"created"`
	LastLogin *time.Time `json:"lastLogin"`
}

// UserData represents the structure of the users.json file
type UserData struct {
	Users []User `json:"users"`
}

// UserService manages user data from JSON file
type UserService struct {
	filePath string
	users    map[string]*User // username -> User mapping
}

// NewUserService creates a new UserService instance
func NewUserService(configPath string) (*UserService, error) {
	service := &UserService{
		filePath: filepath.Join(configPath, "users.json"),
		users:    make(map[string]*User),
	}
	
	err := service.LoadUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to load users: %v", err)
	}
	
	return service, nil
}

// LoadUsers loads user data from the JSON file
func (s *UserService) LoadUsers() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read users file: %v", err)
	}
	
	var userData UserData
	err = json.Unmarshal(data, &userData)
	if err != nil {
		return fmt.Errorf("failed to parse users JSON: %v", err)
	}
	
	// Build username index
	s.users = make(map[string]*User)
	for i := range userData.Users {
		user := &userData.Users[i]
		s.users[user.Username] = user
	}
	
	fmt.Printf("[UserService] Loaded %d users from %s\n", len(s.users), s.filePath)
	return nil
}

// SaveUsers saves user data to the JSON file
func (s *UserService) SaveUsers() error {
	// Convert map back to slice
	var users []User
	for _, user := range s.users {
		users = append(users, *user)
	}
	
	userData := UserData{Users: users}
	
	data, err := json.MarshalIndent(userData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users JSON: %v", err)
	}
	
	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write users file: %v", err)
	}
	
	fmt.Printf("[UserService] Saved %d users to %s\n", len(s.users), s.filePath)
	return nil
}

// ValidateUser validates username and password
func (s *UserService) ValidateUser(username, password string) (*User, bool) {
	user, exists := s.users[username]
	if !exists {
		return nil, false
	}
	
	if !user.Active {
		return nil, false
	}
	
	if user.Password != password {
		return nil, false
	}
	
	// Update last login time
	now := time.Now()
	user.LastLogin = &now
	
	// Save updated user data
	s.SaveUsers()
	
	return user, true
}

// GetUser returns a user by username
func (s *UserService) GetUser(username string) (*User, bool) {
	user, exists := s.users[username]
	return user, exists
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() []*User {
	var users []*User
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// AddUser adds a new user
func (s *UserService) AddUser(user *User) error {
	if _, exists := s.users[user.Username]; exists {
		return fmt.Errorf("user already exists: %s", user.Username)
	}
	
	user.Created = time.Now()
	s.users[user.Username] = user
	
	return s.SaveUsers()
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *User) error {
	if _, exists := s.users[user.Username]; !exists {
		return fmt.Errorf("user not found: %s", user.Username)
	}
	
	s.users[user.Username] = user
	return s.SaveUsers()
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(username string) error {
	if _, exists := s.users[username]; !exists {
		return fmt.Errorf("user not found: %s", username)
	}
	
	delete(s.users, username)
	return s.SaveUsers()
}

// GetUserRoles returns available user roles
func (s *UserService) GetUserRoles() []string {
	return []string{"admin", "lawyer", "user"}
}