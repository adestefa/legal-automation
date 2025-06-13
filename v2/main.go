package main

import (
	"log"
	"net/http"
	"time"

	"mallon-legal-v2/handlers"
	"mallon-legal-v2/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Initialize session service with 24 hour TTL
	sessionService := services.NewSessionService(24 * time.Hour)

	// Session middleware - adds session to context
	router.Use(func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil || sessionToken == "" {
			// Generate a temporary session ID for unauthenticated users
			sessionToken = "temp_session"
		}
		
		// Validate session token header if present
		headerToken := c.GetHeader("X-Session-Token")
		if headerToken != "" && headerToken != sessionToken {
			log.Printf("[WARNING] Session token mismatch: cookie=%s, header=%s", sessionToken, headerToken)
		}
		
		// Add session service and session ID to context
		c.Set("sessionService", sessionService)
		c.Set("sessionID", sessionToken)
		c.Next()
	})

	// Setup CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Serve static files
	router.Static("/static", "frontend")
	
	// Create handlers
	uiHandlers := handlers.NewUIHandlers()
	
	// Serve login page for unauthenticated users
	router.GET("/login", func(c *gin.Context) {
		c.File("frontend/login.html")
	})
	
	// Serve main application for authenticated users
	router.GET("/", func(c *gin.Context) {
		// Check if user is authenticated via cookie
		sessionToken, err := c.Cookie("session_token")
		if err != nil || sessionToken == "" {
			// No valid session cookie, redirect to login
			c.Redirect(http.StatusFound, "/login")
			return
		}
		
		// Set username in context (would come from session validation)
		c.Set("username", "Kevin Mallon")
		
		// Serve main application with Go SSR + HTMX
		uiHandlers.ShowMainPage(c)
	})

	// HTMX UI endpoints for partial page updates
	ui := router.Group("/ui")
	{
		// Step navigation
		ui.GET("/step/:step", uiHandlers.GetStep)
		
		// iCloud folder operations
		ui.GET("/icloud-folders", uiHandlers.GetICloudFolders)
		ui.GET("/case-folders", uiHandlers.GetCaseFolders)
		ui.POST("/select-parent-folder", uiHandlers.SelectParentFolder)
		ui.POST("/select-case-folder", uiHandlers.SelectCaseFolder)
		
		// iCloud setup modal
		ui.GET("/icloud-setup", uiHandlers.ShowICloudSetup)
		ui.POST("/icloud-auth", uiHandlers.HandleICloudAuth)
		
		// Document operations
		ui.GET("/load-documents", uiHandlers.LoadDocuments)
		ui.POST("/select-documents", uiHandlers.SelectDocuments)
		ui.POST("/select-template", uiHandlers.SelectTemplate)
		
		// Document preview
		ui.GET("/preview-document", uiHandlers.PreviewDocument)
		
		// Document viewer/editor
		ui.GET("/view-document", uiHandlers.ViewDocument)
		ui.GET("/edit-document", uiHandlers.EditDocument)
		ui.POST("/save-document", uiHandlers.SaveDocument)
	}

	// Initialize user service
	userService, err := services.NewUserService("config")
	if err != nil {
		log.Printf("Failed to initialize user service: %v", err)
		log.Printf("Falling back to hardcoded credentials")
		userService = nil
	}

	// Legacy API endpoints (keep for backward compatibility during transition)
	api := router.Group("/api")
	{
		// Authentication endpoints (simplified)
		api.POST("/login", func(c *gin.Context) {
			var loginRequest struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			
			if err := c.BindJSON(&loginRequest); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request format"})
				return
			}
			
			var isValid bool
			var user *services.User
			
			// Try user service first, fallback to hardcoded
			if userService != nil {
				user, isValid = userService.ValidateUser(loginRequest.Username, loginRequest.Password)
			} else {
				// Fallback to hardcoded credentials
				isValid = (loginRequest.Username == "admin" || loginRequest.Username == "kmallon" || loginRequest.Username == "demo") && 
						 loginRequest.Password == "password"
			}
			
			if isValid {
				// Set session cookie
				c.SetCookie("session_token", "demo_session_token", 3600, "/", "", false, true)
				
				// Determine display name and role
				displayName := "Kevin Mallon"
				role := "Attorney"
				
				if user != nil {
					// Use data from user service
					displayName = user.Username
					role = user.Role
					log.Printf("[LOGIN] User %s logged in with role: %s", user.Username, user.Role)
				} else {
					log.Printf("[LOGIN] User %s logged in (fallback mode)", loginRequest.Username)
				}
				
				// Return success with user info
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"user": gin.H{
						"username": loginRequest.Username,
						"displayName": displayName,
						"role": role,
					},
					"sessionToken": "demo_session_token",
				})
			} else {
				log.Printf("[LOGIN] Invalid login attempt for user: %s", loginRequest.Username)
				c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid credentials"})
			}
		})
		
		api.GET("/validate-session", func(c *gin.Context) {
			// Get session token from Authorization header or cookie
			authHeader := c.GetHeader("Authorization")
			sessionCookie, _ := c.Cookie("session_token")
			
			// Validate session token (simplified for demo)
			if authHeader != "" || sessionCookie != "" {
			// In production, verify the token with proper session management
			c.JSON(http.StatusOK, gin.H{"valid": true})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
			}
		})
		
		api.POST("/logout", func(c *gin.Context) {
			c.SetCookie("session_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out"})
		})
		
		// Legacy document endpoints (for fallback)
		api.GET("/documents", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{}) // Return empty for now
		})
		
		// Legacy templates endpoint
		api.GET("/templates", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{
				{
					"id":   "fcra-credit-card-fraud",
					"name": "FCRA Complaint - Credit Card Fraud",
					"desc": "For cases involving fraudulent credit card transactions",
				},
			})
		})
	}

	// Start the server
	log.Println("[INFO] Starting Mallon Legal Server v2.5.41 on :8080")
	log.Printf("[INFO] Features: Dynamic document processing (Task 8), document editing, Go SSR + HTMX, Enhanced Session Navigation (Defect 1C)")
	log.Printf("[INFO] Templates directory: /Users/corelogic/satori-dev/clients/proj-mallon/v2/templates")
	log.Printf("[INFO] Test iCloud directory: /Users/corelogic/satori-dev/clients/proj-mallon/test_icloud")
	log.Printf("[INFO] Session TTL: 24 hours with automatic cleanup")
	router.Run(":8080")
}