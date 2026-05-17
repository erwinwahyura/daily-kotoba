# Japanese Learning Data Resources for Kotoba

## 1. Open Data Sources

### JLPT Official Resources
| Resource | URL | Description | License |
|----------|-----|-------------|---------|
| JEES Official | https://www.jees.or.jp/ | Japan Educational Exchanges and Services | Proprietary |
| JLPT Sample Questions | https://www.jlpt.jp/samples/sampleindex.html | Official sample tests | Free use |

### Kanji Datasets
| Resource | URL | Kanji Count | Features |
|----------|-----|-------------|----------|
| **KanjiVG** | https://kanjivg.tagaini.net/ | 6,355 | SVG stroke order paths |
| **Kanji Daishi** | https://www.kanjidatabase.com/ | 13,000+ | Readings, meanings, radicals |
| **Kanjidamage** | https://kanjidamage.com/ | 1,700 | Mnemonics, ordering |
| **WaniKani** | https://www.wanikani.com/ | 2,000+ | Radicals, mnemonics |
| **Jisho.org API** | https://jisho.org/ | 200,000+ words | Dictionary, example sentences |

### Vocabulary Datasets
| Resource | URL | Word Count | Features |
|----------|-----|------------|----------|
| **Core 6000** | http://core6000.neocities.org/ | 6,000 | Frequency-based, audio |
| **Tatoeba** | https://tatoeba.org/ | 1M+ sentences | Example sentences, translations |
| **JMdict** | http://www.edrdg.org/jmdict/j_jmdict.html | 180,000+ | Multilingual dictionary |
| **JLPT Vocab Lists** | https://jlptstudy.net/ | 10,000+ | Organized by level |
| **iKnow.jp** | https://iknow.jp/ | 6,000 | SRS optimized, audio |

### Grammar Resources
| Resource | URL | Patterns | Features |
|----------|-----|----------|----------|
| **Tae Kim's Guide** | https://guidetojapanese.org/ | 100+ | Free, comprehensive |
| **Bunpro** | https://bunpro.jp/ | 1,000+ | SRS, example sentences |
| **Imabi** | https://www.imabi.net/ | 400+ | Detailed explanations |
| **Japanese Grammar Guide** | https://www.kanshudo.com/ | 1,500+ | JLPT focused |
| **JGram** | http://www.jgram.org/ | 500+ | Community grammar database |

## 2. Audio Resources

### Text-to-Speech (TTS)
| Service | Quality | Cost | Notes |
|---------|---------|------|-------|
| **ElevenLabs** | ⭐⭐⭐⭐⭐ | $5/1000 chars | Best quality, natural |
| **Google Cloud TTS** | ⭐⭐⭐⭐ | $4/1M chars | Good, affordable |
| **Amazon Polly** | ⭐⭐⭐⭐ | $4/1M chars | Neural voices |
| **Azure TTS** | ⭐⭐⭐⭐ | $16/1M chars | High quality |
| **VOICEVOX** | ⭐⭐⭐ | Free | Open source, anime-style |
| **OpenAI TTS** | ⭐⭐⭐⭐ | $15/1M chars | Whisper quality |

### Native Audio
| Resource | URL | Content | License |
|----------|-----|---------|---------|
| **Forvo** | https://forvo.com/ | Native pronunciations | CC BY-NC-SA |
| **JapanesePod101** | https://www.japanesepod101.com/ | Lesson audio | Proprietary |
| **NHK World** | https://www3.nhk.or.jp/nhkworld/ | News audio | Free |
| **Shadowing Japanese** | Book/CD | Conversation audio | Commercial |

## 3. Conversation & Dialogue Resources

### Real-life Dialogues
| Resource | URL | Scenarios | Format |
|----------|-----|-----------|--------|
| **NHK Easy Japanese** | https://www.nhk.or.jp/lesson/english/ | Daily life | Audio + Text |
| **Japanese Conversation** | https://www.youtube.com/c/JapaneseConversation/ | Various | YouTube |
| **Real Japanese** | https://reajapanese.com/ | Natural speech | Audio |
| **Shadowing Japan** | Book series | 50+ scenarios | Audio + Text |

### AI Training Data
| Resource | Size | Quality | Use Case |
|----------|------|---------|----------|
| **OpenSubtitles** | 10M+ lines | Medium | Casual conversation |
| **Tatoeba** | 1M+ sentences | High | Example sentences |
| **Japanese Text Initiative** | 100M+ chars | High | Literary/formal |
| **BCCWJ** | 100M+ words | High | Balanced corpus |

## 4. Recommended Data Structure

### Kanji Schema
```json
{
  "id": "kanji-001",
  "character": "日",
  "jlpt_level": "N5",
  "stroke_count": 4,
  "meanings": ["day", "sun", "Japan"],
  "readings": {
    "onyomi": ["にち", "じつ"],
    "kunyomi": ["ひ", "か"]
  },
  "radical": {
    "character": "日",
    "meaning": "sun",
    "position": "whole"
  },
  "components": ["日"],
  "strokes": [
    {
      "number": 1,
      "type": "vertical",
      "path_svg": "M 50 20 L 50 80",
      "direction": {"start": [50, 20], "end": [50, 80]}
    }
  ],
  "mnemonics": {
    "english": "The sun is a vertical line with a horizontal line through it",
    "japanese": "太陽の形を描いた漢字"
  },
  "example_words": [
    {"word": "日本", "reading": "にほん", "meaning": "Japan"}
  ],
  "frequency_rank": 1
}
```

### Vocabulary Schema
```json
{
  "id": "vocab-001",
  "word": "日本",
  "reading": "にほん",
  "meanings": ["Japan"],
  "jlpt_level": "N5",
  "word_type": "noun",
  "tags": ["country", "common"],
  "audio_url": "/audio/nihon.mp3",
  "example_sentences": [
    {
      "japanese": "日本に行きたいです。",
      "reading": "にほんにいきたいです。",
      "english": "I want to go to Japan.",
      "audio_url": "/audio/ex-001.mp3"
    }
  ],
  "frequency_rank": 1,
  "srs_data": {
    "interval_days": 1,
    "ease_factor": 2.5,
    "review_count": 0
  }
}
```

### Grammar Schema
```json
{
  "id": "gram-001",
  "pattern": "～たいです",
  "meaning": "Want to do ~",
  "jlpt_level": "N5",
  "category": "desire",
  "explanation": "Expresses the speaker's desire to do something",
  "formation": "Verb (stem) + たいです",
  "examples": [
    {
      "japanese": "日本に行きたいです。",
      "reading": "にほんにいきたいです。",
      "english": "I want to go to Japan."
    }
  ],
  "related_patterns": ["～たくないです", "～たかったです"],
  "common_mistakes": [
    {
      "incorrect": "食べるたいです",
      "correct": "食べたいです",
      "explanation": "Use stem form, not dictionary form"
    }
  ]
}
```

## 5. Implementation Priorities

### Phase 1: Foundation (Current)
- [x] N5 Kanji (11 characters)
- [x] Basic API structure
- [x] 5 conversation scenarios
- [ ] Expand to full N5 (103 kanji)
- [ ] Add N5 vocabulary (~800 words)
- [ ] Add N5 grammar (~50 patterns)

### Phase 2: Content Expansion
- [ ] N4 Kanji (181 characters)
- [ ] N4 Vocabulary (~1,500 words)
- [ ] N4 Grammar (~100 patterns)
- [ ] 20+ conversation scenarios
- [ ] Audio integration (TTS)

### Phase 3: Advanced Content
- [ ] N3 Kanji (368 characters)
- [ ] N3 Vocabulary (~3,000 words)
- [ ] N3 Grammar (~150 patterns)
- [ ] 50+ conversation scenarios
- [ ] Native audio samples

### Phase 4: Complete Coverage
- [ ] N2 Content (368 kanji, ~6,000 words)
- [ ] N1 Content (1,026 kanji, ~10,000 words)
- [ ] 100+ conversation scenarios
- [ ] Dialect variations
- [ ] Business keigo track

## 6. Data Import Scripts

### Kanji Import (KanjiVG)
```bash
# Download KanjiVG
curl -L https://github.com/KanjiVG/kanjivg/releases/download/r2022.03.21/kanjivg-20220321-main.zip -o kanjivg.zip
unzip kanjivg.zip

# Parse and import
python3 scripts/import_kanjivg.py --input kanjivg/ --output data/kanji.json
```

### Vocabulary Import (JMdict)
```bash
# Download JMdict
curl -L http://ftp.monash.edu.au/pub/nihongo/JMdict_e.gz | gunzip > JMdict_e

# Parse and filter by JLPT level
python3 scripts/import_jmdict.py --input JMdict_e --level N5 --output data/vocab-n5.json
```

### Example Sentences (Tatoeba)
```bash
# Download Tatoeba
curl -L https://downloads.tatoeba.org/exports/sentences.tar.bz2 | tar -xj
curl -L https://downloads.tatoeba.org/exports/links.tar.bz2 | tar -xj

# Filter Japanese-English pairs
python3 scripts/import_tatoeba.py --sentences sentences.csv --links links.csv --output data/sentences.json
```

## 7. Quality Assurance

### Data Validation
- All kanji have stroke order data
- All vocabulary has readings and meanings
- All grammar has at least 3 examples
- All audio files are < 500KB
- All content is properly licensed

### Community Contribution
- Accept user corrections
- Flag inappropriate content
- Suggest new examples
- Rate naturalness of AI responses

## 8. Licensing Notes

### Open Data Licenses
- **CC BY-SA**: ShareAlike - modifications must use same license
- **CC BY-NC-SA**: NonCommercial + ShareAlike
- **EDRDG License**: Free for educational use (JMdict)

### Commercial Use Considerations
- KanjiVG: CC BY-SA 3.0 (attribution required)
- JMdict: Free for use (EDRDG license)
- Tatoeba: CC BY 2.0 FR (attribution required)
- Audio: Check individual sources

---

*Last updated: 2026-04-22*