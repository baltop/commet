import { test, expect } from '@playwright/test';

test.describe('Responsive Design', () => {
  test('should show mobile menu on small screens', async ({ page }) => {
    // Login first
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인/ }).click();

    // Wait for dashboard
    await expect(page).toHaveURL(/\/dashboard/);

    // Set viewport to mobile size
    await page.setViewportSize({ width: 375, height: 667 });

    // Verify hamburger menu button is visible (the button with @click="sidebarOpen = true")
    const hamburgerButton = page.locator('button').filter({ has: page.locator('svg path[d*="M4 6h16"]') });
    await expect(hamburgerButton).toBeVisible();

    // Click hamburger menu button to open mobile sidebar
    await hamburgerButton.click();

    // Verify mobile sidebar is visible
    await expect(page.locator('aside.fixed')).toBeVisible();
  });
});
