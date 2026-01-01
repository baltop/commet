import { test, expect } from '@playwright/test';

test.describe('Authentication - Registration', () => {
  test('should show password match indicator', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');

    const passwordInput = page.locator('input[name="password"]');
    const confirmPasswordInput = page.locator('input[name="confirm_password"]');

    // Type password
    await passwordInput.fill('Test1234!');

    // Type mismatched confirm password
    await confirmPasswordInput.fill('Different');

    // Verify mismatch message is shown
    await expect(page.getByText('비밀번호가 일치하지 않습니다')).toBeVisible();

    // Clear and type matching confirm password
    await confirmPasswordInput.fill('Test1234!');

    // Verify mismatch message is hidden (match indicator shown)
    await expect(page.getByText('비밀번호가 일치하지 않습니다')).not.toBeVisible();
  });
});
