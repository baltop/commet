import { test, expect } from '@playwright/test';

test.describe('Dark Mode', () => {
  test('should toggle dark mode on login page', async ({ page }) => {
    // Navigate and clear localStorage to start fresh
    await page.goto('/auth/login');
    await page.evaluate(() => localStorage.removeItem('darkMode'));
    await page.reload();

    // Verify page is in light mode initially (html element should not have 'dark' class)
    const htmlElement = page.locator('html');
    await expect(htmlElement).not.toHaveClass(/dark/);

    // Click dark mode toggle button (fixed top-right, has @click="darkMode = !darkMode")
    await page.locator('button.fixed').click();

    // Verify page switches to dark mode
    await expect(htmlElement).toHaveClass(/dark/);
  });
});
