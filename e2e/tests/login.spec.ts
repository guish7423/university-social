import { test, expect } from '@playwright/test'

test.describe('Login', () => {
  test('dev login form is visible', async ({ page }) => {
    await page.goto('/login')
    await expect(page.locator('text=开发者登录')).toBeVisible()
    await expect(page.locator('input[placeholder*="用户名"]')).toBeVisible()
  })

  test('dev login works and redirects to home', async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[placeholder*="用户名"]', 'testuser')
    await page.click('button:has-text("开发者登录")')
    await page.waitForURL('/')
    await expect(page.locator('.main-layout')).toBeVisible()
  })
})
