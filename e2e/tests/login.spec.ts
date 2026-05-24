import { test, expect } from '@playwright/test'

test.describe('Login', () => {
  test('dev login form is visible', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('.login-card')).toBeVisible()
    await expect(page.locator('input[placeholder*="昵称"]')).toBeVisible()
  })

  test('dev login works and redirects to home', async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[placeholder*="昵称"]').fill('testuser')
    await page.getByRole('button', { name: '进入校园社' }).click()
    await expect(page.locator('.layout')).toBeVisible({ timeout: 15000 })
  })
})
