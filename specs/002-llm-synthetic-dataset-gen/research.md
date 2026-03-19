# Research: LLM Synthetic Dataset Generation Tool

## Phase 0 Research Findings

### Technology Decisions

#### 1. CLI Framework: Cobra

**Decision**: Use Cobra (github.com/spf13/cobra)

**Rationale**:
- Industry standard for Go CLI (43.5k stars, used by Kubernetes, Hugo, GitHub CLI)
- Built-in support for subcommands, flags, auto-completion
- Well-maintained with excellent documentation
- Seamless integration with Viper for config file support

**Alternatives evaluated**:
- urfave/cli: Good alternative but Cobra has larger ecosystem
- kingpin: Less feature-rich, not as widely used

**Installation**:
```bash
go get -u github.com/spf13/cobra@latest
go install github.com/spf13/cobra-cli@latest
```

---

#### 2. LLM API Client: go-openai

**Decision**: Use go-openai (github.com/sashabaranov/go-openai)

**Rationale**:
- Most popular Go OpenAI client (2822 importers)
- Supports OpenAI-compatible APIs: OpenAI, Azure OpenAI, Anthropic, local Ollama
- Provides full Chat Completions API with streaming support
- Actively maintained with comprehensive documentation
- Configurable base URL for custom endpoints (any OpenAI-compatible API)

**Alternatives evaluated**:
- gpt3: Less maintained, limited features
- openai-go: Less popular, fewer features

**Installation**:
```bash
go get github.com/sashabaranov/go-openai
```

---

#### 3. Configuration: Viper

**Decision**: Use Viper (github.com/spf13/viper)

**Rationale**:
- Works seamlessly with Cobra
- Native TOML support
- Environment variable overrides
- Config file watching

**Alternatives evaluated**:
- standard library encodingTOML: Simpler but no CLI integration
- koanf: Good but extra integration work needed

---

### Implementation Strategy

#### Abstraction Layer for LLM Providers

The spec requires supporting "any OpenAI-compatible API". The go-openai library supports:
- OpenAI (api.openai.com)
- Azure OpenAI
- Anthropic (via custom base URL)
- Ollama (local)
- Any custom OpenAI-compatible endpoint

**Configuration approach**: Use Viper to read TOML config with fields:
- `api_base_url`: Custom endpoint URL
- `api_key`: Authentication token
- `model`: Model name (e.g., gpt-4, gpt-3.5-turbo, llama2)

---

### Key Findings

1. **Cobra + Viper**: Standard combination for Go CLI tools (used by Hugo)
2. **go-openai**: Most feature-complete for OpenAI-compatible APIs
3. **Streaming**: For large datasets, implement streaming to handle 10k records
4. **Rate limiting**: Need to handle rate limits for batch generation
5. **Token budgets**: Plan for ~100 tokens per record generation

---

### Open Questions (Resolved)

| Question | Resolution |
|----------|------------|
| How to support multiple LLM providers? | go-openai supports any OpenAI-compatible API via base URL config |
| Offline/local LLM support? | Ollama works with go-openai via custom base URL |
| TOML config support? | Viper provides native TOML support |

---

## Summary

- **CLI**: Cobra + Cobra CLI generator
- **LLM Client**: go-openai (supports OpenAI-compatible: OpenAI, Azure, Anthropic, Ollama)
- **Config**: Viper (TOML support)
- **Testing**: Go testing + testify

All technology choices align with Constitution principles: minimal dependencies, mainstream languages, modular architecture.
