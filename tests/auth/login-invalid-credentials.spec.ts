import { test, expect } from '@playwright/test';

test.describe('Authentication - Login', () => {
  test('should show error for invalid credentials', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');

    // Enter invalid credentials
    await page.locator('input[name="email"]').fill('invalid@test.com');
    await page.locator('input[name="password"]').fill('wrongpassword');

    // Click submit button
    await page.getByRole('button', { name: /로그인|Sign in/i }).click();

    // Wait for and verify error message
    await expect(page.getByText('이메일 또는 비밀번호가 올바르지 않습니다.')).toBeVisible();

    // Verify email field value is preserved
    await expect(page.locator('input[name="email"]')).toHaveValue('invalid@test.com');
  });
});
