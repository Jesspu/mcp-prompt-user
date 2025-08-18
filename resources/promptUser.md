# promptUser Tool

## Usage

The `promptUser` tool should be used by an LLM or AI agent when it lacks sufficient information, is uncertain about how to proceed, or requires guidance or clarification from the user. This tool allows the agent to pause its current operation and ask the user for the necessary input.

### When to use `promptUser`:

-   **Ambiguous Instructions:** When the user's request is unclear or can be interpreted in multiple ways.
-   **Missing Information:** When the agent requires specific details from the user to complete a task (e.g., a file name, a configuration setting, or a confirmation).
-   **Seeking Guidance:** When the agent has multiple options and needs the user to choose one.
-   **Confirmation:** Before performing a potentially destructive or irreversible action.

---

## Examples

### Example 1: Ambiguous file name

**Scenario:** The user asks to "delete the report," but multiple report files exist.

**Action:**
```json
{
  "tool": "promptUser",
  "prompt": "I found multiple report files: 'report-final.pdf', 'report-draft.docx', and 'report-data.csv'. Which one should I delete?"
}
```

### Example 2: Needing confirmation for a destructive action

**Scenario:** The user asks to "reset the database."

**Action:**
```json
{
  "tool": "promptUser",
  "prompt": "Are you sure you want to reset the database? This action cannot be undone."
}
```

### Example 3: Clarification on a feature implementation

**Scenario:** The user wants to "add authentication," but the specific type isn't specified.

**Action:**
```json
{
  "tool": "promptUser",
  "prompt": "What type of authentication would you like to add? (e.g., OAuth, API Key, JWT)"
}
```
