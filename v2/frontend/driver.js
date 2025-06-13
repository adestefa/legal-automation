  function init() {
            // Initialize authentication
            const sessionToken = localStorage.getItem('sessionToken');
            const storedUser = localStorage.getItem('user');
            
            if (!sessionToken) {
                window.location.href = '/login';
                return;
            }
            
            // Load stored user data
            if (storedUser) {
                try {
                    this.currentUser = JSON.parse(storedUser);
                } catch (e) {
                    console.error('Error parsing stored user data:', e);
                }
            }
            
            // Validate session before making API calls
            fetch('/api/validate-session', {
                headers: { 'Authorization': sessionToken }
            })
            .then(response => response.json())
            .then(data => {
                if (!data.valid) {
                    localStorage.removeItem('sessionToken');
                    localStorage.removeItem('user');
                    window.location.href = '/login';
                    return;
                }
                
                if (data.user) {
                    this.currentUser = data.user;
                    localStorage.setItem('user', JSON.stringify(data.user));
                }
                
                loadApplicationData.call(this, sessionToken);
            })
            .catch(error => {
                console.error('Session validation error:', error);
                localStorage.removeItem('sessionToken');
                localStorage.removeItem('user');
                window.location.href = '/login';
            });
        }
        
        function loadApplicationData(token) {
            const headers = { 'Authorization': token };
            
            fetch('/api/documents', { headers })
                .then(response => response.json())
                .then(data => {
                    this.allDocuments = data;
                    console.log('Loaded documents:', this.allDocuments);
                })
                .catch(error => console.error('Error fetching documents:', error));
                
            fetch('/api/templates', { headers })
                .then(response => response.json())
                .then(data => {
                    this.allTemplates = data;
                    console.log('Loaded templates:', this.allTemplates);
                })
                .catch(error => console.error('Error fetching templates:', error));
                
            // Check for existing iCloud session (don't log 401 as error - it's expected)
            fetch('/api/icloud/validate', { headers })
            .then(response => {
                if (response.status === 401) {
                    // No iCloud session exists, this is normal on fresh login
                    console.log('[INFO] No existing iCloud session - user can set up iCloud later');
                    this.icloudSessionValid = false;
                    return null;
                }
                if (response.ok) {
                    return response.json();
                } else {
                    console.log('[INFO] iCloud session check failed with status:', response.status);
                    return null;
                }
            })
            .then(data => {
                if (data && data.valid) {
                    this.icloudSessionValid = true;
                    console.log('[SUCCESS] Existing iCloud session found for:', data.username);
                    loadICloudFolders.call(this);
                } else {
                    this.icloudSessionValid = false;
                }
            })
            .catch(error => {
                // Don't log this as an error - it's normal when no iCloud session exists
                console.log('[INFO] No existing iCloud session available');
                this.icloudSessionValid = false;
            });
        }
        
        // Document Editing Functions (Task 15)
        function enableDocumentEdit() {
            this.originalDocumentHTML = this.documentHTML;
            this.documentEditMode = true;
            console.log('Document editing enabled');
        }
        
        function saveDocumentEdits() {
            const editableElement = document.getElementById('document-preview');
            if (editableElement) {
                this.documentHTML = editableElement.innerHTML;
                this.documentEditMode = false;
                console.log('Document edits saved');
            }
        }
        
        function cancelDocumentEdits() {
            this.documentHTML = this.originalDocumentHTML;
            this.documentEditMode = false;
            this.editedSections = {};
            console.log('Document edits cancelled');
        }
        
        // Generate clean document HTML for final output (no highlights)
        function generateCleanDocumentHTML(htmlContent) {
            let cleanHTML = htmlContent;
            
            // Remove yellow highlighting classes and inline styles
            cleanHTML = cleanHTML.replace(/class="edited-content"/g, '');
            cleanHTML = cleanHTML.replace(/style="background-color:\s*#fef3c7[^"]*"/g, '');
            cleanHTML = cleanHTML.replace(/style="background-color:\s*#fef3c7[^"]*;\s*color:\s*#[^"]*"/g, '');
            
            // Remove any remaining edited-content classes
            cleanHTML = cleanHTML.replace(/edited-content/g, '');
            
            // Clean up empty class attributes
            cleanHTML = cleanHTML.replace(/class="\s*"/g, '');
            cleanHTML = cleanHTML.replace(/class='\s*'/g, '');
            
            // Clean up empty style attributes
            cleanHTML = cleanHTML.replace(/style="\s*"/g, '');
            cleanHTML = cleanHTML.replace(/style='\s*'/g, '');
            
            return cleanHTML;
        }
        
        function trackContentEdit(event) {
            // Find the nearest paragraph or editable element
            let target = event.target;
            while (target && target.parentNode && target.parentNode.id !== 'document-preview') {
                target = target.parentNode;
            }
            
            // Add yellow highlighting to edited content
            if (target && !target.classList.contains('edited-content')) {
                target.classList.add('edited-content');
                target.style.backgroundColor = '#fef3c7'; // yellow-100
                console.log('Content edited and highlighted:', target.tagName);
            }
        }
        
        function handleKeyDown(event) {
            // Handle Escape key to cancel editing
            if (event.key === 'Escape' && this.documentEditMode) {
                this.cancelDocumentEdits();
            }
            // Handle Ctrl+S to save (prevent default browser save)
            if ((event.ctrlKey || event.metaKey) && event.key === 's' && this.documentEditMode) {
                event.preventDefault();
                this.saveDocumentEdits();
            }
        }
        
        // Logout function
        function logout() {
            const sessionToken = localStorage.getItem('sessionToken');
            
            fetch('/api/logout', {
                method: 'POST',
                headers: { 'Authorization': sessionToken }
            })
            .then(response => response.json())
            .then(data => {
                console.log('Logout successful:', data.message);
            })
            .catch(error => {
                console.error('Logout error:', error);
            })
            .finally(() => {
                localStorage.removeItem('sessionToken');
                localStorage.removeItem('user');
                window.location.href = '/login';
            });
        }
        
        // iCloud Functions (scoped to Alpine.js data)
        
        function loadICloudFolders() {
            const alpineData = this;
            
            fetch('/api/icloud/folders')
            .then(response => response.json())
            .then(data => {
                if (data.folders) {
                    alpineData.icloudFolders = data.folders;
                    console.log('Loaded iCloud folders:', data.folders.length);
                }
            })
            .catch(error => console.error('Error loading iCloud folders:', error));
        }
        
        function loadICloudDocuments(folder) {
            // 'this' is expected to be the Alpine component instance via .call(this)
            const alpineContext = this;
            
            if (!folder) {
                alpineContext.icloudDocuments = [];
                return;
            }
            
            console.log('[DEBUG] Loading iCloud documents from folder:', folder);
            console.log('[DEBUG] Alpine context check - current icloudDocuments length:', alpineContext.icloudDocuments?.length || 'undefined');
            
            // Get session token for authorization
            const sessionToken = localStorage.getItem('sessionToken');
            const headers = {};
            if (sessionToken) {
                headers['Authorization'] = sessionToken;
            }

            fetch(`/api/icloud/documents?folder=${encodeURIComponent(folder)}`, { headers })
            .then(response => {
                console.log('[DEBUG] Document fetch response status:', response.status);
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                console.log('[DEBUG] Document fetch response data:', data);
                if (data && data.documents) {
                    // CRITICAL: Ensure we're updating the Alpine.js reactive data
                    alpineContext.icloudDocuments = data.documents;
                    console.log('[SUCCESS] Successfully loaded', data.documents.length, 'iCloud documents to Alpine context');
                    console.log('[DEBUG] Alpine context icloudDocuments after assignment:', alpineContext.icloudDocuments.length);
                    
                    // Debug: Log each document
                    data.documents.forEach((doc, index) => {
                        console.log(`[DEBUG] Document ${index + 1}: ${doc.name} ${doc.type} ${doc.size} bytes`);
                    });
                    
                    // Force Alpine.js to update the UI by triggering a small delay
                    setTimeout(() => {
                        console.log('[DEBUG] Post-timeout check - Alpine icloudDocuments.length =', alpineContext.icloudDocuments.length);
                    }, 100);
                } else {
                    alpineContext.icloudDocuments = []; // Clear if no documents array in response
                    console.log('[ERROR] No documents found in response for folder:', folder);
                }
                alpineContext.selectedICloudFolder = folder; // Track the folder for which documents were loaded/attempted
            })
            .catch(error => {
                console.error('[ERROR] Error loading iCloud documents for folder:', folder, error);
                alpineContext.icloudDocuments = []; // Clear on error
                
                // Try without authentication as fallback
                console.log('[FALLBACK] Retrying without authentication...');
                fetch(`/api/icloud/documents?folder=${encodeURIComponent(folder)}`)
                .then(response => {
                    if (response.ok) {
                        return response.json();
                    }
                    throw new Error(`Fallback failed with status: ${response.status}`);
                })
                .then(data => {
                    if (data && data.documents) {
                        alpineContext.icloudDocuments = data.documents;
                        console.log('[FALLBACK SUCCESS] Loaded', data.documents.length, 'documents without auth to Alpine context');
                        console.log('[FALLBACK DEBUG] Alpine context icloudDocuments after fallback:', alpineContext.icloudDocuments.length);
                        // Force UI update after fallback success
                        setTimeout(() => {
                            console.log('[FALLBACK DEBUG] Post-timeout - UI should now show', alpineContext.icloudDocuments.length, 'documents');
                        }, 100);
                    }
                })
                .catch(fallbackError => {
                    console.error('[FALLBACK ERROR]', fallbackError);
                });
            });
        }
        
        function loadCaseFolders(parentFolder) {
            if (!parentFolder) return;
            
            console.log('Loading case folders for parent:', parentFolder);
            
            fetch(`/api/icloud/case-folders?parent=${encodeURIComponent(parentFolder)}`)
            .then(response => {
                console.log('Case folders response status:', response.status);
                return response.json();
            })
            .then(data => {
                console.log('Case folders API response:', data);
                if (data.folders) {
                    // Use Alpine.store or find the Alpine component and update it directly
                    const alpineComponent = Alpine.$data(document.querySelector('[x-data]'));
                    if (alpineComponent) {
                        alpineComponent.caseFolders = data.folders;
                        console.log('Updated Alpine caseFolders:', alpineComponent.caseFolders.length);
                    } else {
                        console.error('Could not find Alpine component');
                    }
                } else {
                    console.log('No folders in response:', data);
                    const alpineComponent = Alpine.$data(document.querySelector('[x-data]'));
                    if (alpineComponent) {
                        alpineComponent.caseFolders = [];
                    }
                }
            })
            .catch(error => {
                console.error('Error loading case folders:', error);
                const alpineComponent = Alpine.$data(document.querySelector('[x-data]'));
                if (alpineComponent) {
                    alpineComponent.caseFolders = [];
                }
            });
        }
        
        // Make functions available globally for Alpine.js
        window.init = init;
        window.loadApplicationData = loadApplicationData;
        window.enableDocumentEdit = enableDocumentEdit;
        window.saveDocumentEdits = saveDocumentEdits;
        window.cancelDocumentEdits = cancelDocumentEdits;
        window.generateCleanDocumentHTML = generateCleanDocumentHTML;
        window.trackContentEdit = trackContentEdit;
        window.handleKeyDown = handleKeyDown;
        window.logout = logout;
        window.loadICloudFolders = loadICloudFolders;
        window.loadICloudDocuments = loadICloudDocuments;
        window.loadCaseFolders = loadCaseFolders;