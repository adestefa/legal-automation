# 🏯 Yinsen Project Readme File - Advanced Coding Master Agent
*Your persistent, project-aware development sensei*

**Status:** Active  
**Priority:** High  
**Client:** Mallon Law Firm  
**Deadline:** 2025-07-01  

## Project Overview

**Project**: Mallon Legal Assistant
**Agent**: Yinsen v1.1
**Last Updated**: 2025-05-20
**Maintained By**: Satori Tech Consulting

## Project Scope

The Mallon Legal Assistant is an automated tool for streamlining the legal complaint form generation process for attorneys. Key components include:

1. **Document Selection & Input**:
   - Connect to iCloud account to select source documents (PDFs, meeting notes, etc.)
   - Select a target complaint form template from previous cases
   - Extract text and relevant information from all input documents

2. **Information Extraction & Analysis**:
   - Extract client-specific information from source documents using pattern matching
   - Identify approximately 20 key data points needed for complaint form population
   - Generate structured data representing the client's case details

3. **Review & Approval**:
   - Generate a markdown "preflight review" document summarizing extracted information
   - Present this summary to the attorney for verification and approval
   - Allow manual adjustments to extracted data if needed

4. **Template Population & Generation**:
   - Insert approved information into specific sections of the complaint form
   - Maintain proper legal formatting and document structure
   - Highlight changes in the populated document for easy verification

5. **Finalization & Summary**:
   - Save the completed complaint form to the attorney's iCloud account
   - Provide a summary of all changes made to the document
   - Allow approval or rejection with option to edit further

## Success Criteria

Success will be measured by the following outcomes:

1. **Time Savings**: Reduce the time to create a complaint form from 4+ hours to under 30 minutes
2. **Accuracy**: Correctly extract and place client information in the appropriate sections
3. **Workflow Integration**: Seamless connection with attorney's existing iCloud workflow
4. **Usability**: Intuitive interface that doesn't require technical expertise
5. **Flexibility**: Accommodate different case types and complaint form templates
6. **Quality Control**: Provide clear summaries and highlighting of changes for verification

## Implementation Approach

1. **Backend Processing**:
   - Go-based API endpoints integrated with Sensei server
   - PDF processing with UniDoc/UniPDF library
   - Pattern matching with regular expressions for data extraction
   - Template mapping using JSON schema for flexible insertion

2. **Frontend Interface**:
   - HTML with Tailwind CSS for styling
   - HTMX for dynamic content updates without heavy JavaScript
   - Responsive design compatible with desktop and tablet use

3. **Integration Components**:
   - iCloud API integration for document access
   - Vector database integration (Zilliz Cloud) for semantic search
   - Markdown generation for preflight review
   - Document template engine for final output


## Current Status & Next Steps

1. **Completed**:
   - Analysis of legal documents and complaint form structure
   - Identification of key data points for extraction
   - Planning of implementation strategy and technical approach
   - Creation of prototype design documents and code examples
   - Development of template mapping schema and pattern matching rules

2. **In Progress**:
   - Creating minimal working prototype for demo to Kevin
   - Implementing Go backend for document text extraction
   - Developing HTML/HTMX frontend interface

3. **Upcoming**:
   - Meeting with Kevin to review approach and gather feedback
   - Implementation of full workflow with iCloud integration
   - Testing with sample legal documents
   - Integration with Sensei server endpoints

## Implementation Resources

Key project resources are located in:

- `/build/` - Implementation artifacts and code examples
- `/yinsen/artifacts/` - Technical planning documents and specifications
- `/legal_artifacts/` - Sample legal documents for testing

All technical design documents have been created and stored for reference during implementation.

---

*"A master's knowledge transcends individual sessions through elegant persistence."*
