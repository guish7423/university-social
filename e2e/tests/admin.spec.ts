import { test, expect } from '@playwright/test'

test.describe('Admin', () => {
  test('admin login page loads', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('form')).toBeVisible()
  })
})
