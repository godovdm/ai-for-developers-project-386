## CRITICAL RULES — ОБЯЗАТЕЛЬНО

- NEVER удалять или переписывать рабочие тесты без явного запроса
- NEVER удалять файлы без подтверждения
- ALWAYS запускать тесты после любого изменения кода
- ALWAYS делать git checkpoint перед крупными рефакторингами
- Одна задача за раз. НЕ делать несколько изменений одновременно
- Если не уверен — СПРОСИ, не угадывай

## Working Style

- ДЛЯ ЗАДАЧ СОЗДАНИЯ merge request ИСПОЛЬЗУЙ merge-requester без планировщика
- ЕСЛИ ЗАДАЧА НА РАЗРАБОТКУ: Сначала ПЛАН, потом код
- Маленькие дифы: один файл → тесты → следующий файл
- Используй субагентов для исследования кодовой базы
- Сообщения коммитов оформляй по Conventional Commits: `type(scope): описание` (scope необязателен).

## Agents

- Use `planner` agent для планирования
- Use `tester` agent после изменений кода
- Use `code-reviewer` agent перед коммитами

## Agent skills

### Issue tracker

Track issues and specs in GitHub Issues for
`godovdm/ai-for-developers-project-386`.
Read `docs/agents/issue-tracker.md` before tracker operations.

### Triage labels

Use the five default triage labels.
Read `docs/agents/triage-labels.md` when assigning triage roles.

### Domain docs

Use a single-context layout: root `GLOSSARY.md` and `docs/adr/`.
Read `docs/agents/domain.md` before exploring domain concepts or decisions.
