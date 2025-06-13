# TASK 37

NAME: Fix Document Display - Show Only iCloud Documents
PRIORITY: HIGH
STATUS: QUEUE

## Issue Description
The document selection screen in Step 1 is showing 13 files when there are only 11 files in the iCloud folder. The extra 2 files are coming from the local test folder fallback. This causes confusion as users see documents that don't exist in their iCloud folder.

## Requirements
1. Number the documents in Step 1 for easier reference
2. Ensure only iCloud documents are displayed (11 files, not 13)
3. Remove or disable the local test folder fallback
4. Increment version to v2.5.34 after completion

## Technical Details
- The issue is in the document loading logic in ui_handlers.go
- The fallback to local test documents should be removed or made optional
- Document numbering should be added to the template

## Acceptance Criteria
- [ ] Documents are numbered in Step 1 selection screen
- [ ] Only 11 documents from iCloud folder are shown
- [ ] No test folder documents appear in the list
- [ ] Version updated to v2.5.34