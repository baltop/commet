import { test, expect } from '@playwright/test';

test.describe('Authentication - Registration', () => {
  test('should display registration form with all required fields', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');

    // Verify name input field is visible
    await expect(page.locator('input[name="name"]')).toBeVisible();

    // Verify email input field is visible
    await expect(page.locator('input[name="email"]')).toBeVisible();

    // Verify password input field is visible
    await expect(page.locator('input[name="password"]')).toBeVisible();

    // Verify confirm password input field is visible
    await expect(page.locator('input[name="confirm_password"]')).toBeVisible();

    // Verify terms checkbox is visible
    await expect(page.locator('input[type="checkbox"]')).toBeVisible();

    // Verify submit button is visible (text is "계정 만들기")
    await expect(page.getByRole('button', { name: /계정 만들기/ })).toBeVisible();
  });
});
