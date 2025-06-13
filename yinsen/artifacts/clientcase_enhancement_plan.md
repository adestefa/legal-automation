# ClientCase Struct Enhancement Plan

This document outlines the proposed updates to the ClientCase struct in main.go to incorporate data from the Civil Cover Sheet and additional summons documents.

## Current Implementation

```go
// ClientCase represents the structured data extracted from documents
type ClientCase struct {
    ClientName            string    `json:"clientName"`
    ContactInfo           string    `json:"contactInfo"`
    ResidenceLocation     string    `json:"residenceLocation"`
    FinancialInstitution  string    `json:"financialInstitution"`
    AccountOpenDate       time.Time `json:"accountOpenDate"`
    CreditLimit           string    `json:"creditLimit"`
    TravelLocation        string    `json:"travelLocation"`
    TravelStartDate       time.Time `json:"travelStartDate"`
    TravelEndDate         time.Time `json:"travelEndDate"`
    FraudAmount           string    `json:"fraudAmount"`
    FraudStartDate        time.Time `json:"fraudStartDate"`
    FraudEndDate          time.Time `json:"fraudEndDate"`
    FraudDetails          string    `json:"fraudDetails"`
    DiscoveryDate         time.Time `json:"discoveryDate"`
    DisputeCount          int       `json:"disputeCount"`
    DisputeMethods        []string  `json:"disputeMethods"`
    BankResponse          string    `json:"bankResponse"`
    PoliceReportFiled     bool      `json:"policeReportFiled"`
    PoliceReportDetails   string    `json:"policeReportDetails"`
    CreditBureauDisputes  []string  `json:"creditBureauDisputes"`
    CreditBureauDisputeDate time.Time `json:"creditBureauDisputeDate"`
    AdditionalEvidence    string    `json:"additionalEvidence"`
    CreditImpact          string    `json:"creditImpact"`
}
```

## Proposed Enhancements

```go
// Defendant represents a party being sued
type Defendant struct {
    EntityType     string `json:"entityType"`      // "Credit Bureau" or "Financial Institution"
    Name           string `json:"name"`            // Full legal entity name
    Address        string `json:"address"`         // Legal address for service
    RegisteredAgent string `json:"registeredAgent"` // Agent for service of process
    State          string `json:"state"`           // State of incorporation/registration
}

// ClientCase represents the structured data extracted from documents
type ClientCase struct {
    // Existing Fields
    ClientName            string    `json:"clientName"`
    ContactInfo           string    `json:"contactInfo"`
    ResidenceLocation     string    `json:"residenceLocation"`
    FinancialInstitution  string    `json:"financialInstitution"`
    AccountOpenDate       time.Time `json:"accountOpenDate"`
    CreditLimit           string    `json:"creditLimit"`
    TravelLocation        string    `json:"travelLocation"`
    TravelStartDate       time.Time `json:"travelStartDate"`
    TravelEndDate         time.Time `json:"travelEndDate"`
    FraudAmount           string    `json:"fraudAmount"`
    FraudStartDate        time.Time `json:"fraudStartDate"`
    FraudEndDate          time.Time `json:"fraudEndDate"`
    FraudDetails          string    `json:"fraudDetails"`
    DiscoveryDate         time.Time `json:"discoveryDate"`
    DisputeCount          int       `json:"disputeCount"`
    DisputeMethods        []string  `json:"disputeMethods"`
    BankResponse          string    `json:"bankResponse"`
    PoliceReportFiled     bool      `json:"policeReportFiled"`
    PoliceReportDetails   string    `json:"policeReportDetails"`
    CreditBureauDisputes  []string  `json:"creditBureauDisputes"`
    CreditBureauDisputeDate time.Time `json:"creditBureauDisputeDate"`
    AdditionalEvidence    string    `json:"additionalEvidence"`
    CreditImpact          string    `json:"creditImpact"`
    
    // New Fields from Civil Cover Sheet
    CourtJurisdiction     string    `json:"courtJurisdiction"`     // E.g., "EASTERN DISTRICT OF NEW YORK"
    CourtDivision         string    `json:"courtDivision"`         // E.g., "BROOKLYN"
    CaseClassification    string    `json:"caseClassification"`    // E.g., "CONSUMER CREDIT"
    CauseOfAction         string    `json:"causeOfAction"`         // E.g., "15 U.S.C. § 1681 Fair Credit Reporting Act"
    JuryDemand            bool      `json:"juryDemand"`            // Whether a jury is demanded
    CaseNumber            string    `json:"caseNumber"`            // Assigned case number if available
    
    // Enhanced Attorney Information
    AttorneyName          string    `json:"attorneyName"`          // Attorney's full name
    AttorneyBarNumber     string    `json:"attorneyBarNumber"`     // Attorney's bar number
    AttorneyFirm          string    `json:"attorneyFirm"`          // Law firm name
    AttorneyAddress       string    `json:"attorneyAddress"`       // Complete address
    AttorneyPhone         string    `json:"attorneyPhone"`         // Office phone number
    AttorneyEmail         string    `json:"attorneyEmail"`         // Email address
    
    // Enhanced Defendant Information
    Defendants            []Defendant `json:"defendants"`           // Detailed defendant information
    RelatedCases          []string    `json:"relatedCases"`        // Any related case numbers
    
    // Claim Information
    ClaimAmount           string    `json:"claimAmount"`           // Amount claimed in damages
    FilingDate            time.Time `json:"filingDate"`            // Date of filing
}
```

## Impact on Document Generation

The enhanced struct will require modifications to the `generateDocumentHTML` function to utilize the new fields in the complaint document:

1. **Court Header Section**:
   - Use `CourtJurisdiction` and `CourtDivision` in the document header
   - Include `CaseNumber` if available (or placeholder if not)

2. **Attorney Information**:
   - Update signature block with complete attorney details
   - Include bar number in attorney identifier

3. **Defendant Handling**:
   - Loop through `Defendants` array to dynamically build all defendants
   - Include legal entity names and proper formatting
   - Ensure Equifax is included in credit bureaus

4. **Document Properties**:
   - Add `CaseClassification` and `CauseOfAction` in appropriate sections
   - Include `JuryDemand` flag in the document header

## Implementation Strategy

1. **Update main.go**:
   - Add the new `Defendant` struct
   - Enhance the `ClientCase` struct with new fields

2. **Modify handleGenerateSummary function**:
   - Populate the new fields from extracted document data
   - Include default values for testing purposes

3. **Update generateDocumentHTML function**:
   - Modify the document generation to use the new fields
   - Replace hardcoded values with dynamic data from the enhanced struct

4. **Enhance document extraction**:
   - Add PDF text extraction for Civil Cover Sheet
   - Process Summons_Equifax.pdf for defendant information
   - Extract attorney and court details

## Sample Data Population

```go
// Sample data for testing
clientCase := ClientCase{
    // Existing fields...
    
    // New fields
    CourtJurisdiction: "UNITED STATES DISTRICT COURT\nEASTERN DISTRICT OF NEW YORK",
    CourtDivision:     "BROOKLYN DIVISION",
    CaseClassification: "CONSUMER CREDIT",
    CauseOfAction:     "15 U.S.C. § 1681 Fair Credit Reporting Act",
    JuryDemand:        true,
    
    AttorneyName:      "Kevin Mallon",
    AttorneyBarNumber: "KM1234",
    AttorneyFirm:      "MALLON CONSUMER LAW GROUP",
    AttorneyAddress:   "500 Fifth Avenue, Suite 1900\nNew York, NY 10110",
    AttorneyPhone:     "(212) 732-5777",
    AttorneyEmail:     "kmallon@mallonlaw.com",
    
    Defendants: []Defendant{
        {
            EntityType:     "Financial Institution",
            Name:           "TD BANK, N.A.",
            Address:        "1701 Route 70 East, Cherry Hill, NJ 08034",
            RegisteredAgent: "The Corporation Trust Company",
            State:          "Delaware",
        },
        {
            EntityType:     "Credit Bureau",
            Name:           "EXPERIAN INFORMATION SOLUTIONS, INC.",
            Address:        "475 Anton Boulevard, Costa Mesa, CA 92626",
            RegisteredAgent: "CT Corporation System",
            State:          "Ohio",
        },
        {
            EntityType:     "Credit Bureau",
            Name:           "EQUIFAX INFORMATION SERVICES LLC",
            Address:        "1550 Peachtree Street NE, Atlanta, GA 30309",
            RegisteredAgent: "Corporation Service Company",
            State:          "Georgia",
        },
        {
            EntityType:     "Credit Bureau",
            Name:           "TRANS UNION LLC",
            Address:        "555 West Adams Street, Chicago, IL 60661",
            RegisteredAgent: "Prentice Hall Corporation",
            State:          "Delaware",
        },
    },
    
    ClaimAmount: "$75,000",
    FilingDate:  parseDate("April 9, 2025"),
}
```

## Migration Considerations

1. **Backward Compatibility**:
   - Keep existing field names and structure to maintain compatibility
   - Add new fields with proper JSON tags for API consistency
   - Default values for new fields to handle existing data

2. **Frontend Updates**:
   - Update Alpine.js data bindings to include new fields
   - Enhance the document preview to display new information
   - Add Civil Cover Sheet handling to the document selection UI

3. **Testing Strategy**:
   - Test with both existing data format and enhanced data format
   - Verify document generation with all defendant combinations
   - Validate HTML output with new fields

This enhancement plan provides a comprehensive approach to incorporating data from all available documents while maintaining compatibility with the existing system.
