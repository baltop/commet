import { test, expect } from '@playwright/test';

test.describe('Authentication - Login', () => {
  test('should show error when submitting empty form', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');

    // The form has HTML5 required validation, so we need to bypass it
    // or check that the browser prevents submission

    // Try to submit and check for HTML5 validation
    const submitButton = page.getByRole('button', { name: /로그인/ });
    await submitButton.click();

    // Check that email input shows validation (required field)
    const emailInput = page.locator('input[name="email"]');
    const isInvalid = await emailInput.evaluate((el: HTMLInputElement) => !el.validity.valid);
    expect(isInvalid).toBe(true);
  });
});
