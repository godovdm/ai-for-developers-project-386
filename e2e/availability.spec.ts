import { expect, test } from "@playwright/test";

test("guest selects an available start", async ({
  page,
  request,
}, testInfo) => {
  const id = `availability-${testInfo.project.name}`;
  const created = await request.post("/api/v1/meeting-types", {
    data: { id, name: "Выбор времени", description: "", durationMinutes: 60 },
  });
  expect(created.status()).toBe(201);
  const response = await request.get(
    `/api/v1/meeting-types/${id}/availability`,
  );
  const calendar: { days: { date: string; slots: string[] }[] } =
    await response.json();
  const day = calendar.days.find((item) => item.slots.length > 0);
  expect(day).toBeDefined();
  if (!day) return;
  await page.goto(`/meeting-types/${id}`);
  await page.getByRole("button", { name: day.date }).click();
  await page
    .getByRole("button", { name: day.slots[0].slice(11, 16), exact: true })
    .click();
  await expect(
    page.getByText(
      `Выбрано: ${day.date} ${day.slots[0].slice(11, 16)} (Москва)`,
    ),
  ).toBeVisible();
});
