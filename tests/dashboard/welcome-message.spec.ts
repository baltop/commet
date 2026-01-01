import { test, expect } from '@playwright/test';

test.describe('Dashboard - Welcome Message', () => {
  test('should display welcome message with user name', async ({ page }) => {
    // Login first
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인|Sign in/i }).click();

    // Wait for dashboard to load
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify welcome message displays with user name
    await expect(page.getByText(/안녕하세요.*님!/)).toBeVisible();
  });
});
