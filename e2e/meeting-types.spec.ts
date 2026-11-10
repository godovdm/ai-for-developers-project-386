import { expect, test } from "@playwright/test";

test("owner creates a meeting type and the guest sees it", async ({ page }, testInfo) => {
  const name = `Вводный звонок ${testInfo.project.name}`;

  await page.goto("/owner");
  await page.getByLabel("Название").fill(name);
  await page.getByLabel("Описание").fill("Знакомство с командой");
  await page.getByLabel("Длительность").selectOption("60");
  await page.getByRole("button", { name: "Создать тип встречи" }).click();

  await expect(page.getByRole("heading", { name })).toBeVisible();

  await page.getByRole("link", { name: "Гостю" }).click();

  await expect(page.getByRole("heading", { name: "Выберите время для встречи" })).toBeVisible();
  await expect(page.getByRole("heading", { name })).toBeVisible();
  await expect(page.getByText("60 минут").first()).toBeVisible();
});
