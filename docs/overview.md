# ASCII Name Transliteration Service

Thousands of Australians with non-Anglo names face daily barriers when interacting with government services, banking, insurance, and other systems that weren't designed for multicultural naming conventions. The problem is both systemic and personal.

### The Human Cost

**Trang Le's Story** ([source: ABC News article](./article.md))

When Trang arrived from Vietnam, her full name "Lê Thị Hiền Trang" was incorrectly recorded as "Thị Lê" because Australian systems couldn't handle Vietnamese naming conventions where:

- **Lê** is the family name
- **Thị** indicates female gender
- **Hiền Trang** is her actual given name(s)

This simple misunderstanding led to:

- Missing her Pandemic Leave Disaster payment
- Additional police checks for different name versions
- Trouble with insurance applications (surname "Le" failed minimum 3-letter requirements)
- Flight booking complications requiring her to use the wrong name
- Emotional stress and administrative burden

> _"In the Vietnamese community we have a joke that every female Vietnamese Australian in Australia has the same name, Thị, because of how it's messed up."_ - Trang Le

### The Systemic Problem

**Identity Verification Failures**

Australian databases and systems make incorrect assumptions about naming structures:

- **Vietnamese names**: Gender markers (Văn/Thị) mistaken for given names
- **Indonesian/Myanmar names**: Single names (mononyms) split incorrectly across first/last name fields
- **Chinese names**: Lost tonal information causes innocent people to be arrested for matching wanted suspects
- **Malaysian names**: Patronymic structures (binti/bin) don't fit Western naming expectations

**Criminal Exploitation**

Dr. Fiona Swee-Lin Price notes that naming confusion is being exploited by criminals:

- Multiple identities created using different romanisation of Chinese characters
- Casinos struggling to identify the same person across Mandarin/Cantonese spellings
- Security vulnerabilities in identity verification systems

### The Scale

From the ABC article and supporting data:

- **200,000+ annual arrivals** to Australia
- **Majority have non-Anglo naming structures**
- **Westjustice assisted 450+ individuals** through their "My Name" project alone
- **Multiple government agencies affected**: Services Australia, immigration, police, defence

### Current "Solutions" Are Inadequate

People are being forced to:

- **Change their names** to fit inflexible systems
- **Accept incorrect documentation** that doesn't match their identity
- **Navigate multiple inconsistent records** across different agencies
- **Pay for costly legal assistance** to rectify name errors

As Westjustice legal director Joseph Nunweek states: _"We need to change. We don't need to impose a naming convention on people who already have got one."_

## The Technical Solution

This project provides a **multicultural name decoder** - exactly what Dr Price suggested - to help Australian systems correctly interpret, store, and display names from diverse cultural backgrounds while maintaining compatibility with existing infrastructure.

---

## Project Documentation

### Core Information

- **[Project Summary](./project.md)** - Complete GovHack competition project details
- **[Technical Inspiration](./article.md)** - ABC News article highlighting the problem
- **[Original Technical Specification](./idea.md)** - Detailed implementation blueprint

### Technical Documentation

- **[Architecture](./architecture.md)** - Encore.dev framework and infrastructure approach
- **[Technical Standards](./technical.md)** - Coding conventions and implementation guidelines
- **[Locale Standards](../.claude/locale.md)** - Australian English and regional formatting requirements

### Competition Materials

- **[Video Storyboard](./storyboard.md)** - 3-minute video production plan for GovHack presentation
- **[Competition Runsheet](./.context/todo.md)** - 3-day development timeline and deliverables

---

## The Vision

Every Australian, regardless of their name's cultural origin, should have seamless access to government services, banking, healthcare, and other essential systems.

This isn't just about technical compatibility - it's about **digital inclusion**, **cultural respect**, and **administrative efficiency**.

When complete, this service will:

- **Eliminate identity verification failures** for people with non-Anglo names
- **Reduce administrative burden** for government agencies
- **Improve security** by providing consistent name handling
- **Save millions in processing costs** through automation
- **Preserve cultural identity** while enabling system compatibility

_The technology exists. The need is urgent. The solution is within reach._
