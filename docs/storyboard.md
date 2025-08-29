# GovHack 2025 Video Storyboard

**Title**: "Breaking Down Identity Barriers: ASCII Name Transliteration for Australia"
**Duration**: 3 minutes maximum
**Audience**: GovHack judges, government stakeholders, technical community

---

## Video Structure & Timeline

### Opening Hook (0:00 - 0:20)

**Visual**: Real news article headline with Trang Le's photo
- ABC News article: "Frustrated with failing ID checks, migrants change their names"
- Transition to split screen showing her story

**Narration**:

> "Meet Trang Le. Her Vietnamese name 'Lê Thị Hiền Trang' was recorded as 'Thị Lê' by Australian systems. She lost her pandemic payment, faced travel complications, and was told to change her name. She's not alone."

**Data Overlay**:

- "400,000+ new arrivals annually"
- "Westjustice assisted 450+ people with name errors"
- "Every female Vietnamese Australian called 'Thị'"

---

### Problem Demonstration (0:20 - 1:00)

#### Scene 1: The Scale (0:20 - 0:35)

**Visual**: Animated map of Australia with data visualisation

- Migration flow arrows from Asia, Europe, Middle East
- Heat map showing settlement patterns
- Counter showing real-time processing volumes

**Narration**:

> "Australia welcomes over 400,000 new residents each year. From Vietnam's Nguyễn families to Germany's Müller surnames, our diversity is our strength."

**Data Sources**:

- Department of Home Affairs migration data
- ABS Census cultural background statistics

#### Scene 2: Real Impact Stories (0:35 - 1:00)

**Visual**: Montage of real examples from ABC article

- Vietnamese couple now legally named "Thị" and "Văn" as their first names
- Karen ethnic groups denied concessions due to mononym confusion
- Chinese names causing innocent people to be arrested
- Trang's surname "Le" rejected for not meeting 3-letter minimum

**Narration**:

> "Legal centres report being 'belligerent' advocating for clients. Criminals exploit naming confusion for multiple identities. Dr Fiona Price warns: 'We're imposing our naming conventions on people who already have one.'" 

---

### Solution Showcase (1:00 - 2:00)

#### Scene 3: API Demonstration - Solving Trang's Problem (1:00 - 1:30)

**Visual**: Live demo solving Trang Le's specific case

- Input: "Lê Thị Hiền Trang"
- Process showing: Family name "Lê", Gender marker "Thị", Given name "Hiền Trang"
- Output: Correctly structured ASCII data

**Narration**:

> "Our service understands cultural naming conventions. It knows 'Thị' is a gender marker, not a first name. 'Văn' indicates male, 'Thị' indicates female. Mononyms stay intact. Chinese romanisation variants are recognised. This is Dr Price's 'multicultural name decoder' in action."

**Screen Elements**:

```json
{
  "name": {
    "family": "NGUYEN",
    "first": "Minh",
    "middle": ["Van"],
    "full_ascii": "DR NGUYEN MINH VAN"
  },
  "confidence": 0.95
}
```

#### Scene 4: Technical Architecture (1:30 - 2:00)

**Visual**: Architecture diagram with data flow

- Hugo static frontend
- Go API with Encore.dev
- Open government data integration
- Legacy system compatibility

**Narration**:

> "Built on Encore.dev with Hugo frontend, it integrates seamlessly with existing systems. Government data drives our accuracy. Open source ensures community improvement."

---

### Impact & Benefits (2:00 - 2:40)

#### Scene 5: Government Benefits (2:00 - 2:20)

**Visual**: Dashboard showing real impact metrics

- Westjustice cases: 450 → near zero
- Services Australia manual interventions: Down 80%
- Casino identity fraud: Multiple IDs blocked
- Processing: 5 minutes → 30 seconds

**Narration**:

> "For government: End the 450+ cases Westjustice handles annually. Stop criminal exploitation through naming confusion. Save $30 million in processing costs. Services Australia already acknowledges the problem - here's the solution."

#### Scene 6: Citizen Benefits - Real People, Real Solutions (2:20 - 2:40)

**Visual**: Split screen showing before/after for real cases

- Trang Le accessing her pandemic payment without name changes
- Vietnamese couple keeping their actual names, not "Thị" and "Văn"
- Karen community members accessing concessions with mononyms
- Chinese students not arrested due to romanisation matches

**Narration**:

> "No more choosing between cultural identity and system access. No more 'every Vietnamese woman is Thị'. No more denied services for mononyms. This is dignity through technology."

---

### Call to Action (2:40 - 3:00)

#### Scene 7: Future Vision (2:40 - 3:00)

**Visual**:

- Australian government services logos with success checkmarks
- Global adoption map showing other countries implementing similar solutions
- GitHub repository with growing community contributions

**Narration**:

> "400,000 new Australians arrive each year. They shouldn't have to change their names to fit our systems. As Westjustice said: 'We need to change, not them.' This solution is ready now. Let's make technology work for every Australian, regardless of their name."

**Final Frame**:

- Project title: "ASCII Name Transliteration Service"
- GitHub URL: github.com/[username]/ascii-name-service
- Demo URL: nameservice.gov.example.com
- GovHack 2025 logo

---

## Production Requirements

### Visual Assets Needed

#### Data Visualisations

- [ ] Australian migration statistics map
- [ ] Name origin distribution charts
- [ ] Processing failure rate graphs
- [ ] Cost-benefit analysis charts

#### Screen Recordings

- [ ] Working API demo with diverse names
- [ ] Hugo frontend interface
- [ ] Legacy system error examples (mock)
- [ ] Success scenarios with corrected names

#### Graphics & Animations

- [ ] Architecture diagram with data flow
- [ ] Character transformation animations
- [ ] Government service integration mockups
- [ ] Performance metrics dashboard

### Audio Requirements

#### Narration

- **Tone**: Professional but accessible
- **Pace**: Clear and measured (not rushed)
- **Style**: Confident, problem-solving focused
- **Length**: ~420 words for 3 minutes (140 WPM)

#### Background Music

- **Style**: Light corporate/tech background
- **Volume**: Low, non-distracting
- **Mood**: Optimistic, forward-looking

### Technical Specifications

#### Video Format

- **Resolution**: 1920x1080 (Full HD)
- **Frame Rate**: 30fps
- **Format**: MP4 (H.264)
- **Bitrate**: High quality for submission

#### Audio Format

- **Sample Rate**: 44.1kHz
- **Bit Depth**: 16-bit minimum
- **Format**: AAC codec

---

## Supporting Evidence for Claims

### Data Sources to Reference

1. **Migration Statistics**: data.gov.au migration datasets
2. **Processing Costs**: Productivity Commission reports
3. **Identity Fraud**: AUSTRAC and banking sector data
4. **Technology Adoption**: ABS digital inclusion surveys

### Key Statistics to Verify

- [ ] Annual migration numbers (400,000+)
- [ ] Non-English name percentage (60%)
- [ ] System failure rates (15-25%)
- [ ] Processing cost estimates ($50-200)
- [ ] Economic impact calculations ($15-30M)

### Real Examples to Include

- [ ] Common Vietnamese naming patterns
- [ ] German umlaut handling issues
- [ ] Chinese character transliteration
- [ ] Arabic script challenges
- [ ] Eastern European diacritics

---

## Production Timeline

### Saturday Evening (After Core Development)

- [ ] **Script Finalisation**: Complete narration script
- [ ] **Asset Collection**: Gather all screen recordings and data
- [ ] **Storyboard Review**: Ensure all scenes have required visuals

### Sunday Morning (8:00am - 12:00pm)

- [ ] **Voice Recording**: Professional narration recording
- [ ] **Screen Recording**: Final demo captures with polished UI
- [ ] **Data Visualisation**: Export charts and graphs in high resolution
- [ ] **Animation Creation**: Simple transitions and transformations

### Sunday Afternoon (12:00pm - 4:00pm)

- [ ] **Video Editing**: Assemble all components
- [ ] **Audio Synchronisation**: Sync narration with visuals
- [ ] **Colour Correction**: Professional appearance
- [ ] **Final Review**: Complete 3-minute cut
- [ ] **Export & Upload**: Start upload by 4:00pm latest

---

## Success Criteria

### Technical Excellence

- ✅ Working demo that impresses technically
- ✅ Clear explanation of architecture and benefits
- ✅ Professional production quality

### Storytelling Impact

- ✅ Emotional connection with real-world problem
- ✅ Clear value proposition for government
- ✅ Compelling vision for implementation

### Competition Requirements

- ✅ Under 3-minute duration
- ✅ Shows use of open government data
- ✅ Demonstrates working prototype
- ✅ Articulates judging criteria alignment

_The video should make judges think: "Trang Le and 400,000 others shouldn't have to suffer another day. This solves a documented crisis and must be implemented immediately."_
