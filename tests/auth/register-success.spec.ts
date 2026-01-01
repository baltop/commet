import { test, expect } from '@playwright/test';

test.describe('Authentication - Registration', () => {
  test('should successfully register new user', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');

    // Generate unique email for this test
    const uniqueEmail = `testuser_${Date.now()}@test.com`;

    // Fill in registration form
    await page.locator('input[name="name"]').fill('Test User');
    await page.locator('input[name="email"]').fill(uniqueEmail);
    await page.locator('input[name="password"]').fill('TestPass123!');
    await page.locator('input[name="confirm_password"]').fill('TestPass123!');

    // Check terms checkbox
    await page.locator('input[type="checkbox"]').check();

    // Click submit button
    await page.getByRole('button', { name: /계정 만들기/ }).click();

    // Wait for redirect to login page
    await expect(page).toHaveURL(/\/auth\/login/, { timeout: 10000 });
  });
});
