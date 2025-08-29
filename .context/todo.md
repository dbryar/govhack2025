# GovHack 2025 Competition Runsheet

## Competition Timeline
- **Start**: Friday 30/08/2025 @ 7:00pm
- **Development**: Saturday 31/08/2025 (full day)  
- **Submission Deadline**: Sunday 01/09/2025 @ 5:00pm

## Required Deliverables
1. ✅ Team registration in Hackerspace
2. ✅ Project page with description
3. ✅ Challenge category selection
4. ✅ Open data source URLs
5. ⏳ Working prototype/demo
6. ⏳ 3-minute video pitch
7. ⏳ Evidence repository URL (GitHub)

---

## Day 1: Friday 30/08/2025 (Evening 7pm-11pm)

### Setup & Planning (7:00pm - 8:30pm)
- [ ] **Team Registration**
  - [ ] Complete Hackerspace team setup
  - [ ] Assign team captain role
  - [ ] Select challenge category: "Identity & Digital Government Services"

- [ ] **Project Setup** 
  - [ ] Create GitHub repository for evidence
  - [ ] Initialize project structure per architecture docs
  - [ ] Set up development environment (Go, Hugo, Encore.dev)

- [ ] **Research & Data Collection**
  - [ ] Identify Australian open datasets for migration statistics
  - [ ] Locate Services Australia/Centrelink processing data
  - [ ] Find ABS census data on name origins
  - [ ] Document all data source URLs for submission

### Initial Development (8:30pm - 11:00pm)
- [ ] **Core API Structure**
  - [ ] Set up Go module with Encore.dev
  - [ ] Implement basic HTTP handlers
  - [ ] Create initial transliteration functions
  - [ ] Add Unicode to ASCII conversion

- [ ] **Hugo Site Foundation**
  - [ ] Initialize Hugo static site
  - [ ] Create basic layout and styling
  - [ ] Add API integration JavaScript
  - [ ] Implement demo form interface

---

## Day 2: Saturday 31/08/2025 (Full Development Day)

### Morning Session (8:00am - 12:00pm)
- [ ] **API Enhancement**
  - [ ] Complete name parsing logic
  - [ ] Add gender inference with confidence
  - [ ] Implement Vietnamese/Chinese/European name handling
  - [ ] Add comprehensive error handling

- [ ] **Data Integration**
  - [ ] Connect to Australian migration data APIs
  - [ ] Process and cache relevant statistics
  - [ ] Create data analysis functions
  - [ ] Generate demographic visualisations

### Afternoon Session (1:00pm - 6:00pm)  
- [ ] **Frontend Development**
  - [ ] Build interactive demo interface
  - [ ] Add data visualisation charts (migration stats)
  - [ ] Implement results display with explanations
  - [ ] Create responsive design for mobile

- [ ] **Testing & Validation**
  - [ ] Test with diverse name examples
  - [ ] Validate accuracy against known cases  
  - [ ] Performance testing and optimisation
  - [ ] Cross-browser compatibility testing

### Evening Session (7:00pm - 10:00pm)
- [ ] **Deployment & Polish**
  - [ ] Deploy API to Encore cloud
  - [ ] Deploy Hugo site to Netlify/Vercel
  - [ ] Configure CORS and security settings
  - [ ] Final UI/UX improvements

---

## Day 3: Sunday 01/09/2025 (Submission Day)

### Morning Session (8:00am - 12:00pm)
- [ ] **Video Production Setup**
  - [ ] Write video script highlighting key points:
    - [ ] Problem demonstration with real examples
    - [ ] Solution showcase with live demo
    - [ ] Data visualisation of scale/impact
    - [ ] Benefits to government and citizens
  - [ ] Gather screen recordings of working demo
  - [ ] Collect data visualisation screenshots

### Afternoon Session (12:00pm - 4:00pm)
- [ ] **Video Production**
  - [ ] Record narration and demo footage
  - [ ] Edit video to 3-minute maximum
  - [ ] Add captions and professional polish
  - [ ] Export in required format for upload
  - [ ] **Critical**: Start upload by 4:00pm for buffer time

- [ ] **Final Documentation**
  - [ ] Complete project page description
  - [ ] Finalise README with setup instructions
  - [ ] Document all open data sources used
  - [ ] Create comprehensive evidence repository

### Final Submission (4:00pm - 5:00pm)
- [ ] **Submission Checklist**
  - [ ] Confirm video upload completed successfully
  - [ ] Verify all required fields completed
  - [ ] Test demo URL is publicly accessible
  - [ ] Ensure GitHub repository is public
  - [ ] Submit before 5:00pm deadline

---

## Key Datasets to Source

### Migration & Demographics
- [ ] **data.gov.au**: Annual migration statistics by country
- [ ] **ABS Census**: Language and cultural background data
- [ ] **Department of Home Affairs**: Visa processing volumes

### Identity & Processing Issues  
- [ ] **Services Australia**: Identity verification statistics (if available)
- [ ] **AUSTRAC**: Identity fraud prevention data
- [ ] **State Government**: Service delivery efficiency metrics

### Economic Impact Data
- [ ] **Productivity Commission**: Administrative burden studies
- [ ] **ABS**: Economic impact of migration
- [ ] **Banking/Finance**: Identity verification costs

---

## Development Priorities

### Must Have (Critical Path)
1. **Working API** - Core transliteration functionality
2. **Demo Interface** - Show the solution in action  
3. **Data Visualisation** - Demonstrate problem scale
4. **Video Submission** - Required deliverable

### Should Have (Time Permitting)
1. **Advanced Language Support** - Beyond Vietnamese/Chinese
2. **Performance Optimisation** - Sub-50ms response times
3. **Comprehensive Testing** - Edge case validation
4. **Professional Styling** - Enhanced UI/UX

### Nice to Have (Stretch Goals)
1. **Real-time Statistics** - Live data feeds
2. **Advanced Analytics** - Machine learning insights
3. **Mobile App** - Native interface
4. **API Documentation** - Swagger/OpenAPI spec

---

## Risk Management

### Technical Risks & Mitigations
- **API Development Delays**: Start with simple implementation, enhance iteratively
- **Data Source Issues**: Have backup datasets and static examples ready
- **Deployment Problems**: Test deployment early, have local demo ready
- **Video Production Time**: Allocate full Sunday morning, start upload early

### Competition Risks & Mitigations  
- **Submission Platform Issues**: Complete submission by 4:30pm for buffer
- **Team Coordination**: Use shared documents and regular check-ins
- **Scope Creep**: Focus on core deliverables first, add features second

---

## Daily Success Criteria

### Friday Success: Foundation Complete
- ✅ Team registered and project page created
- ✅ Development environment set up
- ✅ Basic API and Hugo site scaffolded
- ✅ Data sources identified and documented

### Saturday Success: Core Functionality
- ✅ API handles diverse name inputs correctly
- ✅ Demo site shows clear value proposition
- ✅ Data visualisations demonstrate problem scale
- ✅ Solution deployed and publicly accessible

### Sunday Success: Professional Submission
- ✅ 3-minute video showcases solution effectively
- ✅ All GovHack requirements submitted on time
- ✅ GitHub repository provides complete evidence
- ✅ Demo impresses judges with real-world applicability

---

*Focus on delivering a working solution that solves a real problem. Polish is secondary to functionality and clear demonstration of value.*