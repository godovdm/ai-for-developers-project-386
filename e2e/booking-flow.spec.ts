import { expect, test } from "@playwright/test";

test("owner sees a booking made by a guest", async ({
  page,
  request,
}, testInfo) => {
  const id = `flow-${testInfo.project.name}`;
  const name = `Встреча ${testInfo.project.name}`;
  const guest = `Гость ${testInfo.project.name}`;
  const email = `${testInfo.project.name}@example.com`;
  await page.goto("/owner");
  await page.getByLabel("ID").fill(id);
  await page.getByLabel("Название").fill(name);
  await page.getByLabel("Описание").fill("Знакомство");
  await page.getByRole("button", { name: "Создать тип встречи" }).click();
  await expect(page.getByRole("heading", { name })).toBeVisible();
  await page.getByRole("link", { name: "Гостю" }).click();
  await page
    .getByRole("heading", { name })
    .locator("xpath=../..")
    .getByRole("button", { name: "Выбрать время" })
    .click();
  const response = await request.get(
    `/api/v1/meeting-types/${id}/availability`,
  );
  const calendar: { days: { date: string; slots: string[] }[] } =
    await response.json();
  const day = calendar.days.filter((item) => item.slots.length > 0)[
    ["chromium", "firefox", "webkit"].indexOf(testInfo.project.name)
  ];
  expect(day).toBeDefined();
  if (!day) return;
  await page.getByRole("button", { name: day.date }).click();
  await page
    .getByRole("button", { name: day.slots[0].slice(11, 16), exact: true })
    .click();
  await page.getByLabel("Имя").fill(guest);
  await page.getByLabel("Email").fill(email);
  await page.getByRole("button", { name: "Записаться" }).click();
  await expect(page.getByText("Запись подтверждена")).toBeVisible();
  await page.getByRole("link", { name: "Владельцу" }).click();
  const bookings = page.getByRole("region", { name: "Предстоящие записи" });
  await expect(bookings.getByText(name)).toBeVisible();
  await expect(bookings.getByText(guest)).toBeVisible();
  await expect(bookings.getByText(email)).toBeVisible();
  await expect(
    bookings.getByText(`${day.date} ${day.slots[0].slice(11, 16)} (Москва)`),
  ).toBeVisible();
});

for (const outcome of ["confirmed", "conflict"] as const) {
  test(`availability error clears after a ${outcome} booking refresh`, async ({
    page,
    request,
  }, testInfo) => {
    const id = `refresh-${outcome}-${testInfo.project.name}`;
    const created = await request.post("/api/v1/meeting-types", {
      data: {
        id,
        name: "Обновление времени",
        description: "",
        durationMinutes: 60,
      },
    });
    expect(created.status()).toBe(201);
    const response = await request.get(
      `/api/v1/meeting-types/${id}/availability`,
    );
    const calendar: { days: { date: string; slots: string[] }[] } =
      await response.json();
    const day = calendar.days.filter((item) => item.slots.length >= 2)[
      ["chromium", "firefox", "webkit"].indexOf(testInfo.project.name) +
        (outcome === "confirmed" ? 3 : 6)
    ];
    expect(day).toBeDefined();
    if (!day) return;
    const nextSlot = day.slots.find(
      (slot) => Date.parse(slot) >= Date.parse(day.slots[0]) + 60 * 60_000,
    );
    expect(nextSlot).toBeDefined();
    if (!nextSlot) return;
    await page.goto(`/meeting-types/${id}`);
    await page.getByRole("button", { name: day.date }).click();
    await page
      .getByRole("button", { name: day.slots[0].slice(11, 16), exact: true })
      .click();
    await page.getByLabel("Имя").fill("Гость");
    await page.getByLabel("Email").fill("refresh@example.com");
    await page.route(
      `**/api/v1/meeting-types/${id}/availability`,
      (route) => route.abort(),
      { times: 1 },
    );
    await page.getByRole("button", { name: "Записаться" }).click();
    const error = page.getByText("Не удалось загрузить доступное время");
    await expect(error).toBeVisible();
    await page
      .getByRole("button", { name: nextSlot.slice(11, 16), exact: true })
      .click();
    if (outcome === "conflict") {
      const reserved = await request.post("/api/v1/bookings", {
        data: {
          meetingTypeId: id,
          start: nextSlot,
          guestName: "Другой гость",
          guestEmail: "other@example.com",
        },
      });
      expect(reserved.status()).toBe(201);
    }
    await page.getByRole("button", { name: "Записаться" }).click();
    await expect(
      page.getByText(
        outcome === "confirmed"
          ? "Запись подтверждена"
          : "Это время уже занято. Выберите другое время.",
      ),
    ).toBeVisible();
    await expect(error).not.toBeVisible();
    await expect(
      page.getByRole("button", {
        name: nextSlot.slice(11, 16),
        exact: true,
      }),
    ).toHaveCount(0);
  });
}
