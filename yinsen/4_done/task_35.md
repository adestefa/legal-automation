# TASK 35

NAME: Return to use user.json file for user management
PRIORITY: HIGH
STATUS: QUEUE

## Issue Description
Replace the current hardcoded user credentials system with a user.json file-based approach to make it easier to admin the site. This will allow us to add users and assign roles to them from future UI.

## Current State
The current user management system in main.go uses hardcoded credentials:
- `admin`, `kmallon`, `demo` (all with password "password")
- Simple cookie-based session management
- No persistent user storage
- No role-based access control

## Requirements
1. Create a user.json file structure for storing users
2. Replace hardcoded credentials with dynamic loading from user.json
3. Implement user management functions (load, save, validate)
4. Add role-based access control foundation
5. Maintain backward compatibility with existing session management
6. Prepare for future admin UI implementation

## User.json Structure
```json
{
  "users": [
    {
      "id": "1",
      "username": "admin",
      "password": "password", 
      "email": "admin@mallon-law.com",
      "role": "admin",
      "active": true,
      "created": "2025-06-13T00:00:00Z",
      "lastLogin": null
    },
    {
      "id": "2", 
      "username": "kmallon",
      "password": "password",
      "email": "kmallon@mallon-law.com", 
      "role": "lawyer",
      "active": true,
      "created": "2025-06-13T00:00:00Z",
      "lastLogin": null
    }
  ]
}
```

## Implementation Tasks
- [ ] Create user.json file with initial users
- [ ] Create User struct and UserManager service
- [ ] Implement JSON file loading and saving functions
- [ ] Replace hardcoded auth logic with user.json lookup
- [ ] Add role-based access control framework
- [ ] Update session management to store user roles
- [ ] Add user management helper functions
- [ ] Test with existing functionality

## Acceptance Criteria
- [ ] User authentication works from user.json file
- [ ] Existing users (admin, kmallon, demo) can still log in
- [ ] Roles are properly assigned and stored in sessions
- [ ] Easy to add new users by editing user.json
- [ ] Backward compatible with current session system
- [ ] Foundation ready for future admin UI

## Benefits
- Easy user administration without code changes
- Role-based access control foundation
- Scalable user management approach
- Preparation for admin UI development
- Better security through externalized credentials