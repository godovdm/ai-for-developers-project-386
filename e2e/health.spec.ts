import { expect, test } from '@playwright/test';

test('shows the status returned by the server', async ({ page }) => {
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'Календарь встреч' })).toBeVisible();
  await expect(page.getByText('Сервер работает')).toBeVisible();
});
