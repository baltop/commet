import { test, expect } from '@playwright/test';

test.describe('Authentication - Registration', () => {
  test('should show password strength indicator', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');

    const passwordInput = page.locator('input[name="password"]');

    // Type a password that meets minimum requirements (6+ chars) to show indicator
    await passwordInput.fill('abc123');

    // Wait for Alpine.js to update - strength indicator should appear
    // The strength bar container should be visible
    await expect(page.locator('.strength-bar').first()).toBeVisible({ timeout: 5000 });

    // Type stronger password with uppercase, number, special char
    await passwordInput.fill('Abc12345!@#');

    // Verify stronger indicator shows (강함 or 매우 강함)
    await expect(page.getByText(/강함|매우 강함/)).toBeVisible({ timeout: 5000 });
  });
});
