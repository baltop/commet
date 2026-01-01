import { test, expect } from '@playwright/test';

test.describe('Authentication - Registration', () => {
  test('should show validation error when submitting empty form', async ({ page }) => {
    // Navigate to registration page
    await page.goto('/auth/register');

    // Check the terms checkbox to enable submit
    await page.locator('input[type="checkbox"]').check();

    // Try to submit the form
    await page.getByRole('button', { name: /계정 만들기/ }).click();

    // The form has HTML5 required validation
    // Check that name input shows validation (required field)
    const nameInput = page.locator('input[name="name"]');
    const isInvalid = await nameInput.evaluate((el: HTMLInputElement) => !el.validity.valid);
    expect(isInvalid).toBe(true);
  });
});
