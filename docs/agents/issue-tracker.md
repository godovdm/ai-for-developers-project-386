# Issue tracker: GitHub

Issues and specs live in GitHub Issues for
`godovdm/ai-for-developers-project-386`. Use the `gh` CLI.

## Operations

- Verify the target: `gh repo view --json nameWithOwner`.
- Create: `gh issue create --title "..." --body-file <path>`.
  Write multiline bodies to a temporary file.
- Read: `gh issue view <number> --json number,title,body,labels,comments`.
- List: `gh issue list --state open --json number,title,labels`.
  Add label and state filters as needed.
- Comment: `gh issue comment <number> --body-file <path>`.
- Label: `gh issue edit <number> --add-label "..."`.
- Remove a label: `gh issue edit <number> --remove-label "..."`.
- Close: post the explanation as a comment, then use
  `gh issue close <number>`.

Infer the repository from the clone's remote. Before creating the first
issue in a session, verify the target. If it differs from the repository
above, ask the user before publishing.

## Relationships

Use native GitHub sub-issues and issue dependencies when available.
Consult the installed CLI's help for supported commands.

If native relationships are unavailable, put `Part of #<parent>` and
`Blocked by: #<number>, #<number>` at the top of the child issue body.
Maintain a task list of children in the parent.
An issue is unblocked when all its blockers are closed.

## Skill instructions

“Publish to the issue tracker” means create a GitHub issue.
“Fetch the relevant ticket” means read the referenced GitHub issue.

## Pull requests as a triage surface

**PRs as a request surface: no.**

## Wayfinding

The map is an issue labelled `wayfinder:map`.
Link child tickets as sub-issues, or use the relationship fallback above.
Use `wayfinder:research`, `wayfinder:prototype`, `wayfinder:grilling`,
or `wayfinder:task` for child ticket types.

The frontier is the first open child in map order with no open blocker
and no assignee. Claim it with `gh issue edit <number> --add-assignee @me`.
To resolve it, comment with the result, close the issue, and add the
result and its link to the map's Decisions-so-far section.
