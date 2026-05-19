import { test, expect } from '@playwright/test'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090'

async function createTestUser(email: string, password: string, name: string) {
  const res = await fetch(`${API_URL}/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, password }),
  })
  if (!res.ok && res.status !== 201) {
    throw new Error(`Failed to create test user: ${res.status}`)
  }
}

async function login(page: any, email: string, password: string) {
  await page.goto('/join')
  await page.click('button:has-text("Sign in")')
  await page.fill('input[type="email"]', email)
  await page.fill('input[type="password"]', password)
  await page.click('button[type="submit"]')
  await page.waitForURL('/dashboard')
}

test.describe('User CRUD flow', () => {
  test('creates, edits, lists, and deletes a user', async ({ page }) => {
    const email = `crud-${Date.now()}@test.com`
    const password = 'testpass123'
    const name = 'CRUD Test User'

    await createTestUser(email, password, 'Admin User')
    await login(page, email, password)

    await page.getByRole('link', { name: 'Users' }).click()
    await page.waitForURL('/users')

    await page.getByRole('button', { name: /New User/ }).click()
    await page.waitForURL('/users/new')

    const newName = 'New User'
    const newEmail = `new-${Date.now()}@test.com`
    const newPassword = 'newpass123'

    await page.fill('input[placeholder="User full name"]', newName)
    await page.fill('input[placeholder="user@email.com"]', newEmail)
    await page.fill('input[placeholder="Enter password"]', newPassword)
    await page.fill('input[placeholder="Confirm your password"]', newPassword)
    await page.click('button[type="submit"]')

    await page.waitForURL('/users')
    await expect(page.getByText(newName)).toBeVisible()

    await page.getByRole('row', { name: newName }).getByRole('button').first().click()
    await page.waitForURL(/\/users\/[^/]+$/)

    const editedName = 'Edited User'
    await page.fill('input[placeholder="User full name"]', editedName)
    await page.click('button[type="submit"]')

    await page.waitForURL('/users')
    await expect(page.getByText(editedName)).toBeVisible()

    await page.getByRole('row', { name: editedName }).getByRole('button').nth(1).click()
    await expect(page.getByText(/Are you sure/)).toBeVisible()
    await page.getByRole('button', { name: 'Delete' }).click()

    await page.waitForTimeout(500)
    await expect(page.getByText(editedName)).not.toBeVisible()
  })

  test('shows error when passwords do not match', async ({ page }) => {
    const email = `crud-admin-${Date.now()}@test.com`
    const password = 'testpass123'

    await createTestUser(email, password, 'Admin')
    await login(page, email, password)

    await page.getByRole('link', { name: 'Users' }).click()
    await page.waitForURL('/users')
    await page.getByRole('button', { name: /New User/ }).click()
    await page.waitForURL('/users/new')

    await page.fill('input[placeholder="User full name"]', 'Mismatch')
    await page.fill('input[placeholder="user@email.com"]', 'mismatch@test.com')
    await page.fill('input[placeholder="Enter password"]', 'pass123')
    await page.fill('input[placeholder="Confirm your password"]', 'different123')
    await page.click('button[type="submit"]')

    await expect(page.getByText(/Passwords do not match/)).toBeVisible()
  })
})
