import { test, expect } from '@playwright/test';

test.describe('Dashboard - Charts', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인|Sign in/i }).click();
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should load line chart for monthly sales', async ({ page }) => {
    // Wait for line chart container to load via HTMX
    // The chart title should be visible after HTMX loads the partial
    await expect(page.getByText('월별 매출 추이')).toBeVisible({ timeout: 10000 });

    // Verify chart canvas is rendered
    await expect(page.locator('canvas').first()).toBeVisible();
  });
});
