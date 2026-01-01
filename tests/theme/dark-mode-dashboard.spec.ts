import { test, expect } from '@playwright/test';

test.describe('Dark Mode', () => {
  test('should toggle dark mode on dashboard', async ({ page }) => {
    // Login first
    await page.goto('/auth/login');
    await page.evaluate(() => localStorage.removeItem('darkMode'));
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인/ }).click();

    // Wait for dashboard to load
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify page starts in light mode
    await expect(page.locator('html')).not.toHaveClass(/dark/);

    // The dashboard has a dark mode toggle button in the header with title attribute
    // Find the button by its title or by the SVG moon icon
    const darkModeButton = page.locator('button[title*="다크 모드"]').or(
      page.locator('header button').filter({ has: page.locator('svg path[d*="M20.354"]') })
    );

    await darkModeButton.first().click();

    // Verify dashboard switches to dark mode
    await expect(page.locator('html')).toHaveClass(/dark/, { timeout: 5000 });
  });
});
