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

  test('should load pie chart for traffic sources', async ({ page }) => {
    // Wait for pie chart container to load via HTMX (has 400ms delay)
    await expect(page.getByText('트래픽 소스')).toBeVisible({ timeout: 10000 });

    // Verify subtitle is also visible
    await expect(page.getByText('방문자 유입 채널 분석')).toBeVisible();
  });
});
