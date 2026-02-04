# skills-x

AI Agent Skills management tool - Download and manage Claude/AI agent skills.

## Installation

```bash
npm install -g skills-x
```

Or use directly with npx:

```bash
npx skills-x list
```

## Usage

```bash
# List all available skills
skills-x list

# Download a skill to ~/.claude/skills/
skills-x init pdf
skills-x init ui-ux-pro-max

# Download all skills
skills-x init --all

# Specify custom target directory
skills-x init pdf --target ./my-skills
```

## Alternative Installation

If npm installation fails, you can install via Go:

```bash
go install github.com/castle-x/skills-x/cmd/skills-x@latest
```

## Links

- [GitHub Repository](https://github.com/castle-x/skills-x)
- [Skills Documentation](https://github.com/castle-x/skills-x#readme)

## License

MIT
