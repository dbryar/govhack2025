# GovHack 2025 Project Planning Scratchpad

## Project Summary
**ASCII Name Transliteration Service** - solving identity verification failures for foreign names in Australian government systems.

## Key Insights from Research

### Competition Requirements (GovHack)
- **Deadline**: Sunday 5pm (3-day competition)
- **Deliverables**: Working prototype + 3min video + documentation + open data usage
- **Focus**: Real government problem with practical solution

### Technical Approach
- **Backend**: Go with Encore.dev (matches our architectural preferences)
- **Frontend**: Hugo static site (perfect for competition demo)
- **Deployment**: Encore cloud for development, Terraform for production
- **Core Feature**: Unicode → ASCII transliteration with cultural awareness

### Problem Scale
- 200,000+ annual arrivals to Australia
- ~60% have non-English names  
- 15-25% experience identity verification issues
- Estimated $15-30M annual cost in processing delays

## Development Strategy

### Critical Path (Must Complete)
1. **API Core**: Basic transliteration functionality
2. **Demo Site**: Working interface to showcase solution
3. **Data Integration**: Australian migration statistics
4. **Video Production**: 3-minute compelling presentation

### Implementation Notes
- Start with Vietnamese and German examples (well-defined rules)
- Use existing Go libraries for Unicode normalisation
- Focus on government/legacy system compatibility (ICAO standards)
- Emphasise security benefits (fraud prevention)

## Data Sources Identified
- **data.gov.au**: Migration statistics by country
- **ABS Census**: Cultural/linguistic diversity data
- **Department of Home Affairs**: Visa processing volumes
- **Services Australia**: Identity verification processes (if available)

## Video Strategy
- **Hook**: Split-screen showing system accepting "John Smith" vs rejecting "Nguyễn Văn Minh"
- **Problem**: Use real data to show scale and economic impact
- **Solution**: Live demo with impressive technical showcase
- **Benefits**: Clear value proposition for government and citizens
- **Call to Action**: Ready for immediate deployment

## Risk Mitigation
- **Technical**: Start simple, build iteratively
- **Data**: Have backup static examples if APIs fail
- **Time**: Allocate full Sunday morning for video production
- **Submission**: Complete by 4:30pm for buffer time

## Success Metrics for Competition
- Working demo that judges can interact with
- Clear evidence of government data integration
- Professional video that tells compelling story
- Technical solution that could realistically be deployed

## Next Steps (Friday Evening)
1. Set up GitHub repository
2. Initialize Go project with Encore
3. Create Hugo site structure
4. Begin API development with basic transliteration
5. Research and document Australian data sources

---
*This project has real commercial and social value beyond the competition - focus on building something genuinely useful.*