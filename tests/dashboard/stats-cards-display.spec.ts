import { test, expect } from '@playwright/test';

test.describe('Dashboard - Statistics Cards', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/auth/login');
    await page.locator('input[name="email"]').fill('test@example.com');
    await page.locator('input[name="password"]').fill('password123');
    await page.getByRole('button', { name: /로그인/ }).click();
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should display all four statistics cards', async ({ page }) => {
    // Verify Total Users card is visible
    await expect(page.getByText('총 사용자')).toBeVisible();

    // Verify Total Revenue card is visible
    await expect(page.getByText('총 매출')).toBeVisible();

    // Verify Total Orders card is visible
    await expect(page.getByText('총 주문')).toBeVisible();

    // Verify Conversion Rate card is visible
    await expect(page.getByText('전환율')).toBeVisible();
  });
});
