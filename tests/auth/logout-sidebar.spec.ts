import { test, expect } from '@playwright/test';

test.describe('Authentication - Logout', () => {
  test('should logout from sidebar', async ({ page }) => {
    // Login with valid credentials first
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인/ }).click();

    // Wait for dashboard to load
    await expect(page).toHaveURL(/\/dashboard/);

    // Click logout button in sidebar (it's a form with POST to /auth/logout)
    // The logout button contains text "로그아웃" or has an icon
    await page.locator('form[hx-post="/auth/logout"] button').first().click();

    // Wait for redirect to login page
    await expect(page).toHaveURL(/\/auth\/login/, { timeout: 10000 });
  });
});
