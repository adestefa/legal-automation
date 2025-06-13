# Enhanced Document Data Flow Diagram

```mermaid
graph TD
    subgraph "Source Documents"
        A1["Atty_Notes.docx"] --> B1["Text Extraction"]
        A2["Adverse_Action_Letter_Cap_One.pdf"] --> B1
        A3["Barclays_Application_Denial_1.pdf"] --> B1
        A4["Barclays_Application_Denial_2.pdf"] --> B1
        A5["Civil Cover Sheet.pdf"] --> B1
        A6["Summons_Experian.pdf"] --> B1
        A7["Summons_TD Bank.pdf"] --> B1
        A8["Summons_Trans Union.pdf"] --> B1
        A9["Summons_Equifax.pdf"] --> B1
    end

    subgraph "Data Extraction"
        B1 --> C1["Client Information"]
        B1 --> C2["Financial Details"]
        B1 --> C3["Fraud Information"]
        B1 --> C4["Dispute History"]
        B1 --> C5["Credit Impact"]
        B1 --> C6["Court Information"]
        B1 --> C7["Defendant Details"]
        B1 --> C8["Attorney Information"]
    end

    subgraph "Enhanced ClientCase Struct"
        C1 --> D1["Basic Client Data"]
        C2 --> D2["Account Information"]
        C3 --> D3["Fraud Details"]
        C4 --> D4["Dispute History"]
        C5 --> D5["Credit Impact"]
        C6 --> D6["Court Information"]
        C7 --> D7["Defendant Objects"]
        C8 --> D8["Attorney Details"]
    end

    subgraph "Document Generation"
        D1 --> E1["Header Section"]
        D2 --> E2["Account Allegations"]
        D3 --> E3["Fraud Allegations"]
        D4 --> E4["Dispute Narrative"]
        D5 --> E5["Damages Section"]
        D6 --> E1
        D7 --> E6["Defendants Section"]
        D8 --> E7["Signature Block"]
    end

    subgraph "Final Document"
        E1 --> F["Complete Legal Complaint"]
        E2 --> F
        E3 --> F
        E4 --> F
        E5 --> F
        E6 --> F
        E7 --> F
    end

classDef highlight fill:#ffecb3,stroke:#ffab00,stroke-width:2px
classDef new fill:#e3f2fd,stroke:#2196f3,stroke-width:2px
class A5,A9,C6,C7,C8,D6,D7,D8 new
class D1,D2,D3,D4,D5,E1,E2,E3,E4,E5,E6,E7 highlight
```

This diagram illustrates the enhanced data flow with the addition of the Civil Cover Sheet and Equifax summons. The blue nodes represent new components, while the yellow highlighted nodes show updated components that will incorporate the enhanced data.

## Key Improvements in Data Flow

1. **New Document Sources**: 
   - Civil Cover Sheet.pdf provides court information, case classification, and attorney details
   - Summons_Equifax.pdf adds another defendant with proper legal entity information

2. **Enhanced Data Extraction**:
   - Court Information extraction from Civil Cover Sheet
   - More detailed Defendant extraction from all summons documents
   - Attorney information extraction for proper document signing

3. **Structured Data Model**:
   - New Court Information section in ClientCase
   - Enhanced Attorney Details section
   - Structured Defendant objects instead of simple strings
   - Proper legal entity naming and addressing

4. **Improved Document Generation**:
   - More accurate court header section 
   - Complete and properly formatted defendants list 
   - Professional signature block with all required attorney information
   - Proper claim amount and filing details

The enhanced data flow creates a more comprehensive and accurate legal document that incorporates all available information from the source documents.
