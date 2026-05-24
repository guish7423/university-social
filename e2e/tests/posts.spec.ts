import { test, expect } from '@playwright/test'

test.describe('Posts', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.fill('input[placeholder*="用户名"]', 'testuser')
    await page.click('button:has-text("开发者登录")')
    await page.waitForURL('/')
  })

  test('square page loads posts', async ({ page }) => {
    await page.goto('/square')
    await expect(page.locator('.post-card').first()).toBeVisible({ timeout: 10000 })
  })

  test('post detail page opens', async ({ page }) => {
    await page.goto('/square')
    const firstPost = page.locator('.post-card').first()
    await expect(firstPost).toBeVisible({ timeout: 10000 })
    await firstPost.click()
    await expect(page.locator('.post-detail')).toBeVisible({ timeout: 5000 })
  })
})
