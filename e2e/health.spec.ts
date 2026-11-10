import { expect, test } from "@playwright/test";

test("shows the guest page served by the application", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("link", { name: "Calendar" })).toBeVisible();
  await expect(page.getByRole("heading", { name: "Выберите время для встречи" })).toBeVisible();
});
