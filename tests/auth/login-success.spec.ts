import { test, expect } from '@playwright/test';

test.describe('Authentication - Login', () => {
  // Note: This test requires a valid test user to exist in the database
  // You may need to run database seeds or create the user before running this test
  test('should successfully login with valid credentials', async ({ page }) => {
    // Navigate to login page
    await page.goto('/auth/login');

    // Enter valid test user credentials
    // Using default seeded user credentials from the application
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');

    // Click submit button
    await page.getByRole('button', { name: /로그인|Sign in/i }).click();

    // Wait for redirect to dashboard
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify dashboard page displays user's name in welcome message
    await expect(page.getByText(/안녕하세요.*님!/)).toBeVisible();
  });
});
