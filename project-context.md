# Japanese Vocabulary Learning App - Project Context

## Project Overview
A mobile application focused on Japanese vocabulary learning through daily word exposure and home screen widget integration. The app aims to improve vocabulary retention through passive reinforcement via widget visibility throughout the day.

## Problem Statement
- **Duolingo Issue**: Repetition without deep explanation leads to poor retention; users with 880+ day streaks still struggle to remember vocabulary
- **Anki Issue**: Requires actively opening the app to study, easy to skip or forget
- **Market Gap**: Existing solutions lack passive, ambient learning through consistent visual exposure

## Core Solution
- Level-appropriate vocabulary delivery based on JLPT levels (N5-N1)
- Rich explanations for deep understanding
- **Home screen widget** for passive, repeated exposure (20-30+ daily glances)
- Ambient learning: see the word throughout the day without opening the app

## User Flow
1. **Initial Setup**: User takes placement test to determine starting level (N5/N4/N3/N2/N1)
2. **Daily Learning Cycle**:
   - Widget displays: vocab + short meaning (minimal cognitive load)
   - User taps widget → opens app with detailed information
   - In app: full explanation, nuance, example sentences, usage notes
   - After learning (opened once), widget does reinforcement throughout the day
3. **Skip Functionality**: User can mark word as "known" → immediately shows next word
4. **Progress**: Track which words are learned, skipped, or need review

## Technical Architecture

### Backend (Go)
- User authentication and management
- Vocabulary database (organized by JLPT level)
- Progress tracking and analytics
- Level placement test logic
- API endpoints for:
  - User profile and progress
  - Daily word fetching
  - Skip/mark as known
  - Vocabulary data by level
  - Placement test results

### iOS App (Swift/SwiftUI)
- Main application UI
- WidgetKit integration for home screen widget
- Local caching (widget works offline)
- Deep linking from widget to app
- StoreKit for subscription management
- Core features:
  - Placement test interface
  - Word detail view (full explanation)
  - Progress tracking UI
  - Settings (level selection, skip management)
  - Widget configuration

### Data Structure

#### Vocabulary Entry
```
{
  "id": "uuid",
  "word": "諦める",
  "reading": "あきらめる",
  "short_meaning": "to give up",
  "detailed_explanation": "...",
  "example_sentences": [...],
  "jlpt_level": "N4",
  "index": 45,
  "usage_notes": "..."
}
```

#### User Progress
```
{
  "user_id": "uuid",
  "current_level": "N4",
  "current_vocab_index": 45,
  "known_vocab_ids": [...],
  "skipped_vocab_ids": [...],
  "review_queue": [...],
  "streak_days": 15,
  "last_updated": "2026-01-12T10:30:00Z"
}
```

## Feature Breakdown

### MVP Features (Phase 1)
1. ✅ Placement test to determine starting level
2. ✅ Daily word with short + detailed explanation
3. ✅ Home screen widget (start with one size: small or medium)
4. ✅ Skip/mark as known functionality (instant next word)
5. ✅ Basic progress tracking
6. ✅ Vocabulary organized by JLPT level (N5-N1)
7. ✅ Widget auto-updates at midnight
8. ✅ Deep link from widget to app

### Phase 2 Features (Post-MVP)
- Review algorithm: 80% new words + 20% review from previously learned
- Advanced analytics (skip patterns, retention metrics)
- Multiple widget sizes (small, medium, large)
- Customization options (widget theme, font size)
- Offline mode improvements
- Export progress data

### Future Considerations
- Android version
- Additional languages beyond Japanese
- Social features (study groups, streaks sharing)
- AI-generated personalized example sentences
- Integration with Japanese media (manga, anime, games)

## Vocabulary Rotation Logic

### Phase 1: Initial Introduction (Days 1-200 for N4)
- Sequential display: vocab index 0 → 199
- Widget updates daily at midnight (00:00)
- User can skip ahead at any time → next word appears immediately
- Track skipped words for future review

### Phase 2: Review + New Content (Days 201+)
- **Option A**: If user stays at N4 level
  - Mix of new cycle (80%) + random review (20%)
  - Shuffle order for second cycle through N4
- **Option B**: If user progresses to N3
  - Begin N3 vocabulary sequence
  - Occasional N4 review words mixed in

### Skip Ahead Implementation
- **Behavior**: Unlimited skips allowed
- **Tracking**: Log all skip events with timestamp
- **Data collected**:
  - `vocab_id`, `user_id`, `skipped_at`, `reason: "already_known"`
- **Use cases for data**:
  - Identify words that are too easy for a level
  - Detect if user should be bumped up a level (15+ skips in one session)
  - Words to include in future review cycles

## Widget Implementation Details

### iOS WidgetKit Specifics
- Uses timeline entries for updates
- Supports small, medium, large sizes (start with small/medium)
- Deep links into app when tapped
- Can display offline (cached data)
- Updates:
  - Scheduled: daily at midnight
  - Manual: when user skips to next word
  - Background: when app is opened

### Widget Display
**Small Widget**:
```
┌─────────────────┐
│  諦める          │
│  (あきらめる)    │
│                 │
│  to give up     │
└─────────────────┘
```

**Medium Widget**:
```
┌───────────────────────────────┐
│  諦める (あきらめる)           │
│  to give up; to abandon       │
│                               │
│  Tap for details              │
│  Day 45/200 • N4              │
└───────────────────────────────┘
```

## Monetization Strategy

### Subscription Model
- **Free Tier**:
  - One level only (e.g., N5 or user's tested level)
  - 50 words
  - Basic widget
  - Ads (non-intrusive)
  
- **Premium Subscription** ($2.99 - $4.99/month):
  - All JLPT levels (N5-N1)
  - Unlimited vocabulary
  - No ads
  - Advanced analytics
  - Custom widget themes
  - Review algorithm access

### Pricing Strategy
- Start lower ($2.99) to build user base
- iOS users more likely to convert to paid
- Annual option: $29.99/year (2 months free)

## Development Timeline (Solo Developer)

### Week 1-2: Backend Foundation
- Set up Go backend structure
- Database schema design
- User authentication API
- Vocabulary data model
- Basic CRUD endpoints

### Week 3-4: Vocabulary Data & Placement Test
- Source/create N4 vocabulary dataset (200 words)
- Design placement test algorithm
- Implement placement test API
- Create test questions/logic

### Week 5-7: iOS App Basics
- Learn SwiftUI fundamentals (if new to Swift)
- Set up iOS project structure
- User authentication flow
- Main app navigation
- Word detail view UI

### Week 8: Widget Development
- Learn WidgetKit
- Implement basic widget (one size)
- Timeline provider for daily updates
- Deep linking to app
- Local caching for offline

### Week 9: Core Features Integration
- Skip/mark as known functionality
- Progress tracking
- Widget update triggers
- API integration (fetch daily word)

### Week 10-11: Polish & Testing
- UI/UX refinements
- Bug fixes
- TestFlight beta testing
- Performance optimization
- Error handling

### Week 12: Launch Prep
- StoreKit subscription setup
- App Store assets (screenshots, description)
- Privacy policy & terms
- Final testing
- Soft launch to friends/small group

## Success Metrics

### User Engagement
- Daily active users (DAU)
- Widget tap rate
- Words learned per user per week
- Skip rate per word (quality indicator)
- Retention: Day 7, Day 30, Day 90

### Learning Effectiveness
- Self-reported vocabulary retention
- Time to complete a JLPT level
- Review test scores (Phase 2)
- User testimonials

### Business Metrics
- Free to paid conversion rate (target: 5-10%)
- Monthly recurring revenue (MRR)
- Churn rate
- Lifetime value (LTV)

## Risks & Mitigations

### Technical Risks
- **Risk**: Widget doesn't update reliably on iOS
  - **Mitigation**: Test timeline provider thoroughly; provide manual refresh option
- **Risk**: App crashes or performance issues
  - **Mitigation**: Extensive testing; crash analytics (Firebase Crashlytics)

### Product Risks
- **Risk**: Users don't find passive widget learning effective
  - **Mitigation**: A/B test different widget formats; collect feedback early
- **Risk**: Vocabulary quality issues (wrong meanings, poor examples)
  - **Mitigation**: Source from reputable databases; user feedback system

### Market Risks
- **Risk**: Crowded market, hard to stand out
  - **Mitigation**: Focus on unique widget UX; target Japanese learners specifically
- **Risk**: Low conversion to paid
  - **Mitigation**: Ensure free tier is useful but limited; clear value prop for premium

## Open Questions

1. **Vocabulary Source**: Where to get high-quality N4 vocabulary with example sentences?
   - Options: JLPT official lists, WaniKani API, custom curation
   
2. **Placement Test Design**: How many questions? What format?
   - Suggestion: 20 questions, multiple choice, mix of vocab + grammar
   
3. **Widget Design**: Minimalist vs. feature-rich?
   - Lean toward minimalist for MVP (just word + short meaning)
   
4. **Subscription Timing**: When to prompt for upgrade?
   - After 7 days of free use? After hitting 50-word limit?

5. **Review Algorithm**: Simple random or spaced repetition (SM2/Anki-style)?
   - Start simple (random 20%), add SRS in Phase 2

## Next Steps

1. Validate concept with potential users (Japanese learners at N4-N3 level)
2. Set up development environment (Go backend, iOS project)
3. Source/create initial N4 vocabulary dataset (200 words minimum)
4. Design database schema and API structure
5. Create wireframes for key screens (placement test, word detail, widget)
6. Begin backend development (authentication, vocab API)
7. Parallel: Learn SwiftUI and WidgetKit fundamentals

## Resources & References

### Technical
- [WidgetKit Documentation](https://developer.apple.com/documentation/widgetkit)
- [SwiftUI Tutorials](https://developer.apple.com/tutorials/swiftui)
- [Go Backend Best Practices](https://github.com/golang-standards/project-layout)

### Vocabulary Sources
- JLPT official word lists
- [Jisho.org](https://jisho.org) - Japanese dictionary
- [WaniKani](https://wanikani.com) - kanji/vocab learning
- [Bunpro](https://bunpro.jp) - grammar points

### Design Inspiration
- Duolingo widget
- Anki mobile
- Drops language app
- Todoist widget design

---

## Personal Context (Developer)
- **Developer**: Erwin
- **Current Level**: Japanese N4
- **Goal**: Build personal learning tool first, potential startup/side project later
- **Tech Stack**: Go backend (comfortable), iOS/Swift (learning)
- **Platform**: iOS first (using iPhone)
- **Timeline**: Solo development, aiming for 12-week MVP

## Core Philosophy
"The best learning happens when information is present but not intrusive. The widget makes vocabulary visible 20-30 times per day through normal phone use, creating effortless repetition that sticks."