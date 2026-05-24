import { test, expect } from '@playwright/test'

test.describe('Explore', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.locator('input[placeholder*="昵称"]').fill('testuser')
    await page.getByRole('button', { name: '进入校园社' }).click()
    await expect(page.locator('.layout')).toBeVisible({ timeout: 15000 })
  })

  test('friends page loads', async ({ page }) => {
    await page.goto('/#/friends')
    await expect(page.locator('.friends-page')).toBeVisible({ timeout: 10000 })
  })

  test('notifications page loads', async ({ page }) => {
    await page.goto('/#/notifications')
    await expect(page.locator('.notif-page')).toBeVisible({ timeout: 10000 })
  })

  test('search page loads', async ({ page }) => {
    await page.goto('/#/search')
    await expect(page.locator('.search-page')).toBeVisible({ timeout: 10000 })
  })

  test('found items page loads', async ({ page }) => {
    await page.goto('/#/found')
    await expect(page.locator('.found-page')).toBeVisible({ timeout: 10000 })
  })

  test('favorites page loads', async ({ page }) => {
    await page.goto('/#/favorites')
    await expect(page.locator('.favorite-page')).toBeVisible({ timeout: 10000 })
  })

  test('whispers page loads', async ({ page }) => {
    await page.goto('/#/whispers')
    await expect(page.locator('.whisper-page')).toBeVisible({ timeout: 10000 })
  })

  test('activities page loads', async ({ page }) => {
    await page.goto('/#/activities')
    await expect(page.locator('.activities-page')).toBeVisible({ timeout: 10000 })
  })

  test('university page loads', async ({ page }) => {
    await page.goto('/#/university')
    await expect(page.locator('.university-page')).toBeVisible({ timeout: 10000 })
  })

  test('account security page loads', async ({ page }) => {
    await page.goto('/#/account/security')
    await expect(page.locator('.security-page')).toBeVisible({ timeout: 10000 })
  })

  test('chat page loads', async ({ page }) => {
    await page.goto('/#/chat')
    await expect(page.locator('.chat-page')).toBeVisible({ timeout: 10000 })
  })

  test('search performs query', async ({ page }) => {
    await page.goto('/#/search')
    await expect(page.locator('.search-page')).toBeVisible({ timeout: 10000 })
    // Type into search input and verify results appear
    const searchInput = page.locator('.search-page input[type="text"], .search-page input:not([type="hidden"])').first()
    await searchInput.fill('test')
    await searchInput.press('Enter')
    await page.waitForTimeout(500)
    // Should show search results or empty state
    const hasContent = await page.locator('.stagger-item, .empty-state').first().isVisible().catch(() => false)
    expect(hasContent).toBe(true)
  })
})
