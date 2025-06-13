# Defect 1b: UI State Restoration for Navigation

## Overview
Update HTML templates to read and display saved user selections from session data, ensuring the UI reflects previously entered information when users navigate back to earlier steps.

## Problem Statement
After implementing session infrastructure (Defect 1a), the templates need to be updated to read session data and pre-populate form fields, show selected options, and display previously entered information when users return to a step.

## Scope
- Template modifications to show preserved state
- Form field pre-population
- Visual indicators for selected items

## Technical Requirements
1. **Step 1 Template Updates**
   - Pre-select case folder from session
   - Display case details from session data
   - Show folder as "selected" in UI
   - Preserve scroll position if possible

2. **Step 2 Template Updates**
   - Highlight previously selected template
   - Show template as "active" with visual indicators
   - Pre-populate any template-specific fields
   - Maintain template selection state

3. **Step 3 Template Updates**
   - Pre-populate all form fields from session
   - Restore any document generation state
   - Show progress indicators if applicable
   - Preserve any generated document links

4. **Template Data Structure**
   ```go
   type TemplateData struct {
       // Existing fields
       CaseFolders   []CaseFolder
       Templates     []Template
       
       // New session-based fields
       SelectedFolder    string
       CaseDetails      map[string]string
       SelectedTemplate string
       FormData         map[string]string
       IsReturningUser  bool
   }
   ```

## Implementation Tasks
- [ ] Update Step 1 template (_step1_case_selection.gohtml)
- [ ] Update Step 2 template (_step2_template_selection.gohtml)
- [ ] Update Step 3 template (_step3_document_generation.gohtml)
- [ ] Add session data to template context in handlers
- [ ] Create helper functions for template state restoration
- [ ] Add visual indicators for restored state

## Success Criteria
- Form fields pre-populate with session data
- Previously selected items are visually highlighted
- Users can see they're returning to a previous step
- All user input is preserved and displayed correctly
- Templates gracefully handle missing session data

## Dependencies
- Defect 1a (Session Infrastructure) must be completed first

## Testing Plan
1. Manual testing of each step with session data
2. Test navigation flow: Step 1 → Step 2 → Step 1 (verify state)
3. Test edge cases (missing session data, corrupted data)
4. Cross-browser testing for form restoration

## Estimated Time
2-3 hours

## Priority
High - Critical for user experience

## Related Files
- dev/templates/_step1_case_selection.gohtml
- dev/templates/_step2_template_selection.gohtml
- dev/templates/_step3_document_generation.gohtml
- dev/handlers/ui_handlers.go (template data preparation)