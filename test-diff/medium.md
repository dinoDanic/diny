# Medium Test File

This is a medium-sized test file for git diff testing with more substantial content.

## Introduction
This file contains moderate amount of content to test how the AI handles medium-sized diffs. It should be large enough to test the system but not overwhelming.

## Features
- Multiple sections
- More detailed explanations
- Code examples
- Lists and formatting

## Code Examples

Here's some sample code:

```javascript
function generateCommitMessage(diff) {
    const prompt = buildPrompt(diff);
    return ai.process(prompt);
}

const config = {
    useEmoji: true,
    useConventional: true,
    tone: 'casual',
    length: 'normal'
};
```

## Configuration Options

### Emoji Settings
- Enable emoji prefixes: `✨ feat:` `🐛 fix:` `♻️ refactor:`
- Disable for plain text commits

### Conventional Commits
- Follow conventional commit format
- Use types: feat, fix, docs, style, refactor, test, chore
- Include scope when relevant: `feat(auth): add login system`

### Tone Options
1. **Professional** - Formal and matter-of-fact
2. **Casual** - Light but clear (default)
3. **Friendly** - Warm and approachable

### Length Settings
- **Short**: Subject only, no body
- **Normal**: Subject + optional body bullets
- **Long**: Subject + detailed body explanations

## Usage Examples

```bash
# Initialize configuration
diny init

# Generate commit message
diny commit

# Use with custom model
diny commit --model qwen2.5-coder:7b
```

## Best Practices
- Stage your changes first with `git add`
- Review generated messages before committing
- Customize configuration to match your team's style
- Use conventional commits for better changelog generation

This concludes the medium-sized test file content.