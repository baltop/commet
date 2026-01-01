import { test, expect } from '@playwright/test';

test.describe('Authentication - Login', () => {
  test('should display login form with all fields', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');

    // Verify email input field is visible
    await expect(page.locator('input[name="email"]')).toBeVisible();

    // Verify password input field is visible
    await expect(page.locator('input[name="password"]')).toBeVisible();

    // Verify submit button is visible
    await expect(page.getByRole('button', { name: /로그인/ })).toBeVisible();

    // Verify link to registration page exists (text is "무료로 시작하기")
    await expect(page.getByRole('link', { name: /무료로 시작하기/ })).toBeVisible();
  });
});
