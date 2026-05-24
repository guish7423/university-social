import { test, expect } from '@playwright/test'

test.describe('Circles', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[placeholder*="用户名"]', 'testuser')
    await page.click('button:has-text("开发者登录")')
    await page.waitForURL('/')
  })

  test('circle list loads', async ({ page }) => {
    await page.goto('/circles')
    await expect(page.locator('.circle-card').first()).toBeVisible({ timeout: 10000 })
  })

  test('circle detail page opens', async ({ page }) => {
    await page.goto('/circles')
    const firstCircle = page.locator('.circle-card').first()
    await expect(firstCircle).toBeVisible({ timeout: 10000 })
    await firstCircle.click()
    await expect(page.locator('.circle-detail')).toBeVisible({ timeout: 5000 })
  })
})
