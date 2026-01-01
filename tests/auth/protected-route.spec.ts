import { test, expect } from '@playwright/test';

test.describe('Authentication - Protected Routes', () => {
  test('should redirect unauthenticated user to login', async ({ page }) => {
    // Clear any existing auth cookies
    await page.context().clearCookies();

    // Navigate directly to dashboard (protected route)
    await page.goto('/dashboard');

    // Verify user is redirected to login page
    await expect(page).toHaveURL(/\/auth\/login/);
  });
});
