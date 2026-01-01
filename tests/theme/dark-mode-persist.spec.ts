import { test, expect } from '@playwright/test';

test.describe('Dark Mode', () => {
  test('should persist dark mode preference', async ({ page }) => {
    // Navigate and clear localStorage to start fresh
    await page.goto('/auth/login');
    await page.evaluate(() => localStorage.removeItem('darkMode'));
    await page.reload();

    // Enable dark mode by clicking the toggle
    await page.locator('button.fixed').click();

    // Verify dark mode is enabled
    await expect(page.locator('html')).toHaveClass(/dark/);

    // Refresh the page
    await page.reload();

    // Verify dark mode is still enabled after refresh
    await expect(page.locator('html')).toHaveClass(/dark/);

    // Verify localStorage has the dark mode preference
    const darkModeValue = await page.evaluate(() => localStorage.getItem('darkMode'));
    expect(darkModeValue).toBe('true');
  });
});
