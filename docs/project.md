# GovHack 2025: What's in a Name?

## Team Details

### Hsckerspace

- Team page: [Transliterates](https://hackerspace.govhack.org/teams/12345)
- Project page: [What's In A Name?](https://hackerspace.govhack.org/projects/12345)

### Members

| Name         | Details                                | Role                            |
| ------------ | -------------------------------------- | ------------------------------- |
| Daniel Bryar | ([@dbryar](https://github.com/dbryar)) | Project Lead, Software Engineer |

### Evidence of Eligibility

- Respository: [GitHub](https://github.com/dbryar/govhack2025)
- Video: ![Presentation](https://user-images.githubusercontent.com/1234567/abcdefg.mp4)
- Demo: [Encore.dev](https://example.com)
- Data Sources:
  - [data.gov.au](https://data.gov.au)
  - [ABS](https://abs.gov.au)

## Project Overview

**Application**: ASCII Name Transliteration Service
**Challenge Category**: Identity & Digital Government Services
**Timeline**: 3 days (Friday 7pm - Sunday 5pm)

### Problem Statement

Based on the ABC News investigation (see [article](./article.md)), multicultural names create systematic barriers in Australian (and other Western) systems of Government and commerce due to cultural anomolies that are unable to be addressed by legacy systems.

**Real-World Impact Examples:**

- **Vietnamese names**: Trang Le's story - gender markers (Thị/Văn) recorded as first names, causing payment failures and travel complications
- **Single names (mononyms)**: Karen, Kareni, and Chin ethnic groups unable to access concessions, prescriptions, or banking
- **Chinese names**: Identical romanisation causing innocent people to be arrested for matching wanted suspects
- **Malaysian names**: Patronymic structures (binti/bin) incompatible with first/last name assumptions

**Documented Consequences:**

1. **Identity Verification Failures**: Westjustice "My Name" project assisted 450+ individuals with incorrect name records
2. **Administrative Burden**: Legal centres report being "belligerent" advocating for clients, requiring costly manual interventions
3. **Security Vulnerabilities**: Criminals exploit naming confusion to create multiple identities (casinos, financial systems)
4. **Economic Impact**: Productivity losses from processing delays, police checks, and system failures
5. **Cultural Erosion**: People changing names to "fit" systems rather than systems adapting to people

### Target Impact

**Addressing Documented Problems:**

- **Primary**: Eliminate the need for people like Trang Le to choose between cultural identity and system compatibility
- **Secondary**: Reduce manual interventions by legal centres and government agencies (80% reduction target)
- **Tertiary**: Prevent criminal exploitation of naming system vulnerabilities
- **Quaternary**: Enable proper cultural name preservation while maintaining ASCII system compatibility

## Technical Solution

### Core Service: ASCII Transliteration API

**Purpose**: Convert Unicode names to ASCII-compatible structured JSON for legacy systems

**Key Features**:

- Unicode to ASCII transliteration (Vietnamese, Chinese, Arabic, European)
- Structured parsing (family/given/middle/titles)
- Gender inference with confidence scoring
- Cultural naming convention awareness
- CORS-enabled for static sites

**Example Transformation**:

```
Input: "Doctor Nguyễn Văn Minh"
Output: {
  "name": {
    "family": "NGUYEN",
    "first": "Minh",
    "middle": ["Van"],
    "full_ascii": "DR NGUYEN MINH VAN"
  },
  "gender": { "value": "M", "confidence": 0.65 }
}
```

### Architecture

**Backend**: Go service with Encore.dev framework

- RESTful API with client authentication
- Rate limiting and CORS controls
- Open data integration for migration statistics

**Frontend**: Hugo static site

- Public demonstration interface
- API documentation
- Migration data visualisations

**Deployment**:

- Backend: Encore cloud (development)
- Frontend: Netlify/Vercel
- Production: Terraform export to AWS/Azure

## Data Sources & Evidence

### Australian Government Open Data

1. **Migration Statistics** (data.gov.au)

   - Annual arrivals by country of origin
   - Visa categories and processing volumes
   - Settlement patterns and demographics

2. **Centrelink/Services Australia Data**

   - Identity verification failure rates
   - Processing time statistics
   - System error categories

3. **ABS Census Data**
   - Name origin distributions
   - Linguistic diversity statistics
   - Cultural background demographics

### Problem Scale (Evidence-Based)

**From ABC Article & Supporting Data:**

- **Annual Immigration**: 400,000+ new arrivals with diverse naming conventions
- **Affected Demographics**: Vietnamese (Thị/Văn gender markers), Indonesian/Myanmar (mononyms), Chinese (romanisation variants), Malaysian (patronymic structures)
- **Documented Cases**: 450+ individuals assisted by single legal centre's "My Name" project
- **System Failures**: Identity document mismatches, payment failures, travel complications, insurance rejections
- **Processing Burden**: Legal centres, Services Australia, police, casinos, universities all affected
- **Criminal Exploitation**: Multiple identity creation, security system bypassing
- **Economic Impact**: Estimated $15-30M annually in processing costs, legal fees, and lost productivity

## GovHack Deliverables

### 1. Working Prototype ✅

- Functional API with Go backend
- Hugo demonstration site
- Live deployment with test data

### 2. Video Presentation (3 minutes max) ✅

- Problem demonstration
- Solution showcase
- Data visualisation
- Impact projection

### 3. Documentation ✅

- Technical specifications
- API documentation
- Implementation guide
- Open data references

### 4. Evidence Repository ✅

- GitHub repository with complete source
- Data analysis and visualisations
- Testing results and benchmarks
- Deployment instructions

## Value Proposition

### Government Benefits

**Direct Solutions to ABC Article Issues:**

1. **Cost Reduction**: Eliminate manual name correction processes (Westjustice handled 450+ cases from one legal centre alone)
2. **Security Enhancement**: Prevent criminal exploitation of naming inconsistencies (casino multiple identity issues)
3. **Service Delivery**: End payment failures like Trang Le's Pandemic Leave Disaster payment
4. **Cultural Competency**: Proper handling of Vietnamese Thị/Văn markers, Indonesian mononyms, Chinese romanisation
5. **Legal Compliance**: Reduce discrimination against people with non-Anglo names

### Citizen Benefits

**Addressing Real Stories:**

1. **Cultural Preservation**: People like Trang Le won't need to adopt incorrect gender markers as their "first name"
2. **Identity Consistency**: End situations where "every female Vietnamese Australian has the same name, Thị"
3. **System Access**: Karen, Kareni, and Chin communities can access concessions and banking with mononyms
4. **Travel Freedom**: Consistent name handling across Australian and international systems
5. **Dignity**: End the need for people to change their cultural names to fit inflexible systems

### Technical Benefits

1. **Legacy Integration**: Works with existing systems
2. **Standards Compliance**: ICAO/passport compatible output
3. **Scalability**: Cloud-native architecture
4. **Open Source**: Community-driven improvements

## Implementation Roadmap

### Phase 1: Core Service (3 days - GovHack)

- ✅ API development and testing
- ✅ Hugo demo site with visualisations
- ✅ Open data integration
- ✅ Video production and submission

### Phase 2: Government Pilot (3 months)

- Services Australia integration
- Real-world testing with anonymised data
- Performance optimisation
- Security hardening

### Phase 3: National Rollout (12 months)

- Multi-agency deployment
- Advanced language support
- Machine learning enhancements
- International standard compliance

## Success Metrics

### Technical KPIs

- **API Response Time**: <50ms average
- **Accuracy Rate**: >95% correct transliteration
- **Uptime**: 99.9% availability
- **Throughput**: 1000+ requests/minute

### Business KPIs

- **Processing Time Reduction**: 60% decrease
- **Error Rate Reduction**: 80% fewer manual interventions
- **Cost Savings**: $15-30M annually
- **User Satisfaction**: >90% positive feedback

## Risk Mitigation

### Technical Risks

- **Unicode Complexity**: Extensive test coverage for edge cases
- **Performance**: Caching and CDN deployment
- **Security**: Input validation and rate limiting

### Business Risks

- **Privacy Concerns**: No data retention, audit logging
- **Cultural Sensitivity**: Community consultation and feedback
- **Integration Complexity**: Gradual rollout with fallback systems

## Team & Resources

### Skills Required

- Go/TypeScript backend development
- Hugo static site generation
- Government data analysis
- Video production and editing

### Infrastructure

- Encore.dev development platform
- GitHub for version control
- Netlify/Vercel for frontend hosting
- Government open data APIs

---

_This project aims to solve a real problem affecting hundreds of thousands of Australians while demonstrating the power of open government data and modern cloud-native architecture._
