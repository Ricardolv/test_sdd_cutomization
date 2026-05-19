import { test, expect } from '@playwright/test'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090'

test.describe('Login flow', () => {
  test('registers a user, logs in, and sees dashboard', async ({ page }) => {
    const email = `e2e-${Date.now()}@test.com`
    const password = 'testpass123'

    await page.goto('/join')

    await page.fill('input[type="text"]', 'E2E User')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await expect(page.getByText('Account created successfully!')).toBeVisible()

    await page.click('button:has-text("Sign in")')

    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await page.waitForURL('/dashboard')
    await expect(page).toHaveURL(/\/dashboard/)

    await expect(page.getByRole('button', { name: /E2E User/ })).toBeVisible()
  })

  test('shows error on wrong password', async ({ page }) => {
    await page.goto('/join')
    await page.click('button:has-text("Sign in")')

    await page.fill('input[type="email"]', 'wrong@test.com')
    await page.fill('input[type="password"]', 'wrongpass')
    await page.click('button[type="submit"]')

    await expect(page.getByText(/Invalid email or password/)).toBeVisible()
  })

  test('redirects authenticated user from /join to /dashboard', async ({ page }) => {
    const email = `e2e-redirect-${Date.now()}@test.com`
    const password = 'testpass123'

    await page.goto('/join')
    await page.fill('input[type="text"]', 'Redirect User')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await expect(page.getByText('Account created successfully!')).toBeVisible()

    await page.click('button:has-text("Sign in")')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await page.waitForURL('/dashboard')

    await page.goto('/join')
    await page.waitForURL('/dashboard')
  })

  test('logout clears session and redirects to /join', async ({ page }) => {
    const email = `e2e-logout-${Date.now()}@test.com`
    const password = 'testpass123'

    await page.goto('/join')
    await page.fill('input[type="text"]', 'Logout User')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await expect(page.getByText('Account created successfully!')).toBeVisible()

    await page.click('button:has-text("Sign in")')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await page.waitForURL('/dashboard')

    await page.getByRole('button', { name: /Logout User/ }).click()
    await page.getByRole('button', { name: 'Logout', exact: true }).click()
    await page.waitForURL('/join')

    await page.goto('/dashboard')
    await page.waitForURL('/join')
  })

  test('persists session after page reload', async ({ page }) => {
    const email = `e2e-persist-${Date.now()}@test.com`
    const password = 'testpass123'

    await page.goto('/join')
    await page.fill('input[type="text"]', 'Persist User')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await expect(page.getByText('Account created successfully!')).toBeVisible()

    await page.click('button:has-text("Sign in")')
    await page.fill('input[type="email"]', email)
    await page.fill('input[type="password"]', password)
    await page.click('button[type="submit"]')

    await page.waitForURL('/dashboard')
    await expect(page.getByRole('button', { name: /Persist User/ })).toBeVisible()

    await page.reload()
    await expect(page.getByRole('button', { name: /Persist User/ })).toBeVisible()
    await expect(page).toHaveURL(/\/dashboard/)
  })
})
