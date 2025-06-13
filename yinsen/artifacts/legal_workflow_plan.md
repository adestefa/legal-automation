# Legal Document Processing Workflow - Implementation Plan

## Project Overview
Extend the Sensei server with endpoints to automate legal complaint form generation by extracting relevant data from iCloud documents (PDFs and meeting notes) and populating a complaint template.

## Architecture Decision
Direct Go implementation with new Sensei server endpoints rather than n8n workflow engine to maintain:
- Single binary deployment
- Minimal dependencies 
- Fast execution
- Direct control over document processing

## Required Endpoints

### 1. `/api/legal/connect-icloud` (POST)
- Handle iCloud OAuth authentication
- Store auth tokens securely
- Return connection status

### 2. `/api/legal/documents` (GET)
- List available documents from connected iCloud account
- Filter for PDFs and relevant file types
- Return document metadata (name, size, modified date)

### 3. `/api/legal/extract` (POST)
- Accept array of selected document IDs
- Process PDFs using `unidoc/unipdf` library
- Extract text and chunk for embedding generation
- Generate embeddings via Ollama API
- Store vectors in Zilliz Cloud with metadata
- Use semantic search for relevant legal precedents
- Use LLM to identify relevant legal data
- Return structured data in JSON format

### 4. `/api/legal/search` (GET)
- Semantic search across stored legal documents
- Query Zilliz Cloud for relevant case precedents
- Return ranked results with similarity scores

### 5. `/api/legal/preview` (GET)
- Generate markdown summary from extracted data
- Include relevant precedents from semantic search
- Present formatted preview to user
- Allow for manual edits if needed

### 6. `/api/legal/generate-complaint` (POST)
- Accept complaint form template
- Insert extracted data into specific form sections
- Include relevant precedents from vector search
- Save completed form (PDF/DOCX format)
- Return file path/download link

### 7. `/api/legal/summary` (GET)
- Generate summary of changes made to complaint form
- Compare original template vs populated form
- Return concise change log for lawyer review

## Technical Stack

### Core Dependencies
```go
import (
    "github.com/unidoc/unipdf/v3"    // PDF processing
    "github.com/milvus-io/milvus-sdk-go/v2"  // Zilliz Cloud integration
    // iCloud API client (research available options)
    // Ollama API client (HTTP calls)
)
```

### Vector Database: Zilliz Cloud
- **Free Tier**: 1M vectors, 2 collections, perfect for prototyping
- **Free Credits**: $100 initial + $100 marketplace bonus = $200 total
- **Credit Extension**: 1 year expiration with payment method added
- **API Integration**: Direct REST API calls from Go
- **Semantic Search**: Sub-10ms query latency
- **Managed Service**: No infrastructure overhead

### Frontend Implementation
- HTML/HTMX interface with Tailwind CSS
- Real-time progress indicators during document processing
- Markdown preview with edit capabilities
- One-click complaint generation
- Responsive design for desktop/tablet use

### Key Features
1. **iCloud Integration**: Secure OAuth connection to access user documents
2. **Smart Extraction**: LLM-powered identification of relevant legal data
3. **Semantic Search**: Vector search for similar cases and precedents via Zilliz Cloud
4. **Template Mapping**: Intelligent insertion of data into complaint form sections
5. **Review Workflow**: Markdown preview before final generation
6. **Change Tracking**: Summary of all modifications for lawyer review

## Implementation Priority
1. Basic PDF text extraction and chunking for embeddings
2. Zilliz Cloud integration and vector storage
3. Ollama integration for embeddings and LLM processing
4. iCloud API integration
5. Semantic search and precedent matching
6. Complaint form template engine
7. HTMX frontend interface
8. Change tracking and summary generation

## Estimated Scope
- Core functionality: 500-1000 lines of Go code
- Vector operations: Additional 200-300 lines for Zilliz integration
- Frontend interface: Clean HTML/HTMX with Tailwind styling
- Testing with sample legal documents
- Documentation for lawyer workflow

## Success Criteria
- Lawyer can connect to iCloud and select relevant documents
- System accurately extracts case-relevant information
- Semantic search returns relevant legal precedents
- Complaint form is properly populated with minimal manual intervention
- Generated summary allows for quick review of changes
- End-to-end process completes in under 2 minutes

## Infrastructure Strategy: Zilliz Cloud vs Self-Hosted Milvus

### Phase 1: Zilliz Cloud (Recommended Start)
**Benefits:**
- Free tier: 1M vectors, 2 collections (perfect for early legal docs)
- $200 in credits ($100 initial + $100 marketplace bonus)
- Credit extension to 1 year with payment method
- Zero infrastructure management
- Built-in monitoring and backups
- Sub-10ms search latency

**Cost Analysis:**
- Free tier sufficient for initial 1000+ legal documents
- Paid plans start around $100/month for production scale
- No infrastructure costs or maintenance

### Phase 2: Self-Hosted Milvus (Migration Path)
**Benefits:**
- Complete control and data sovereignty
- Lower long-term costs at scale
- Custom configuration and optimization
- Runs perfectly on Linode infrastructure

**Considerations:**
- Requires Docker or standalone deployment
- Manual backup and monitoring setup
- Additional Linode instances for high availability

### Recommended Approach
Start with Zilliz Cloud for rapid development and validation, then migrate to self-hosted Milvus once you hit scale or want full control. The Milvus SDK supports both seamlessly.

## Notes
- Maintain Sensei server's single binary deployment model
- Ensure secure handling of sensitive legal documents
- Design for speed and simplicity over feature complexity
- Consider legal compliance requirements for document handling