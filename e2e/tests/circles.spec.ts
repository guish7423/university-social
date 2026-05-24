import { test, expect } from '@playwright/test'

test.describe('Circles', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[placeholder*="昵称"]').fill('testuser')
    await page.getByRole('button', { name: '进入校园社' }).click()
    await expect(page.locator('.layout')).toBeVisible({ timeout: 15000 })
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
    await expect(page.locator('.detail-page')).toBeVisible({ timeout: 5000 })
  })
})
