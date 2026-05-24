import { test, expect } from '@playwright/test'

test.describe('Home', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[placeholder*="昵称"]').fill('testuser')
    await page.getByRole('button', { name: '进入校园社' }).click()
    await expect(page.locator('.layout')).toBeVisible({ timeout: 15000 })
  })

  test('home page loads with content sections', async ({ page }) => {
    // Trending section should exist
    await expect(page.locator('.home-page')).toBeVisible({ timeout: 10000 })
    const sectionTitles = await page.locator('h2').allTextContents()
    expect(sectionTitles.length).toBeGreaterThan(0)
  })

  test('search button navigates to search page', async ({ page }) => {
    await expect(page.locator('.home-page')).toBeVisible({ timeout: 10000 })
    await page.getByText('搜索校园社').click()
    await expect(page.locator('.search-page, .square-page')).toBeVisible({ timeout: 5000 })
  })
})

