import { test, expect } from '@playwright/test';

test.describe('Dashboard - Charts', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인/ }).click();
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should load bar chart for category sales', async ({ page }) => {
    // Wait for bar chart heading to load via HTMX
    await expect(page.getByRole('heading', { name: '카테고리별 판매' })).toBeVisible({ timeout: 10000 });

    // Verify canvas elements are rendered for charts
    const canvasElements = page.locator('canvas');
    await expect(canvasElements.first()).toBeVisible({ timeout: 10000 });
  });
});
