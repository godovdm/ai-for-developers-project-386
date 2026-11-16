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
- Каждый коммит реализации и исправлений связывай с GitHub issue: указывай `#номер` в заголовке или теле коммита.
- Перед завершением изменений запускай `mise --locked run check` и добивайся успешного результата.
- Не редактируй вручную сгенерированные `mise.lock` и `.mise/locks/`.
- Не отключай правила проверок ради прохождения: узкое исключение с причиной требует решения человека.

## Запуск и проверки

Команды выполняются из корня проекта, если не указано иное:

1. Подготовка: `mise install`, `mise exec -- pnpm install --frozen-lockfile`, `mise exec -- pnpm build`.
2. Сервер: `(cd backend && mise exec -- go run ./cmd/server)`; для разработки клиента отдельно: `mise exec -- pnpm dev`.
3. Установи браузеры: `mise exec -- pnpm exec playwright install chromium firefox webkit`.
4. Тесты и линтеры: `mise --locked run check` (как в CI). Отдельные проверки: `mise --locked run check:backend`, `check:frontend`, `check:api`, `check:e2e` и `check:docs`; для каждой задачи повторяй полный префикс `mise --locked run`.

TypeSpec в `api/` — источник API. После его изменения выполни генерацию в порядке из раздела README «Для генерации и проверки артефактов»: OpenAPI `api/generated/openapi.yaml`, Go API `backend/internal/api/generated/api.gen.go`, TypeScript SDK `frontend/src/api/generated/`, затем патч SDK и проверку генерации. Не редактируй эти сгенерированные артефакты вручную.

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
