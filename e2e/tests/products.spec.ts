import { test, expect } from '@playwright/test'

test.describe('Products', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[placeholder*="昵称"]').fill('testuser')
    await page.getByRole('button', { name: '进入校园社' }).click()
    await expect(page.locator('.layout')).toBeVisible({ timeout: 15000 })
  })

  test('product list loads with tabs', async ({ page }) => {
    await page.goto('/#/products')
    await expect(page.locator('.products-page')).toBeVisible({ timeout: 10000 })
    const tabs = page.locator('.tab-btn')
    await expect(tabs.first()).toBeVisible()
  })

  test('category tab switches', async ({ page }) => {
    await page.goto('/#/products')
    await expect(page.locator('.products-page')).toBeVisible({ timeout: 10000 })
    const tabs = page.locator('.tab-btn')
    if (await tabs.count() > 1) {
      await tabs.nth(1).click()
      await page.waitForTimeout(300)
    }
  })
})