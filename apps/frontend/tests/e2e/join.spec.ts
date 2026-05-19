import { test, expect } from '@playwright/test'
import { randomUUID } from 'crypto'

test.describe('Join Page - Cadastro e Login', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/join')

    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        console.log(`[Browser Error] ${msg.text()}`)
      }
    })

    page.on('requestfailed', (request) => {
      console.log(`[Request Failed] ${request.url()} - ${request.failure()?.errorText}`)
    })
  })

  function getToasts(page: ReturnType<typeof test>) {
    return page.locator('[data-sonner-toast]')
  }

  test('deve exibir formulario de cadastro por padrao', async ({ page }) => {
    await page.screenshot({ path: 'tests/e2e/screenshots/01-register-form.png' })

    await expect(page.getByRole('heading', { name: /create account|criar conta/i })).toBeVisible()
    await expect(page.getByLabel(/name|nome/i)).toBeVisible()
    await expect(page.getByLabel(/email/i)).toBeVisible()
    await expect(page.getByLabel(/password|senha/i)).toBeVisible()
    await expect(page.getByRole('button', { name: /^create account$|^criar conta$/i })).toBeVisible()
    await expect(page.getByRole('button', { name: /^sign in$|^entrar$/i })).not.toBeVisible()
  })

  test('deve alternar para formulario de login', async ({ page }) => {
    await page.getByRole('button', { name: /already have|ja tem/i }).click()

    await page.screenshot({ path: 'tests/e2e/screenshots/02-login-form.png' })

    await expect(page.getByRole('heading', { name: /sign in|entrar/i })).toBeVisible()
    await expect(page.getByLabel(/name|nome/i)).not.toBeVisible()
    await expect(page.getByLabel(/email/i)).toBeVisible()
    await expect(page.getByLabel(/password|senha/i)).toBeVisible()
    await expect(page.getByRole('button', { name: /^sign in$|^entrar$/i })).toBeVisible()
  })

  test('deve cadastrar com sucesso e exibir toast', async ({ page }) => {
    const uniqueEmail = `test-${randomUUID()}@example.com`

    await page.getByLabel(/name|nome/i).fill('Test User')
    await page.getByLabel(/email/i).fill(uniqueEmail)
    await page.getByLabel(/password|senha/i).fill('password123')

    await page.getByRole('button', { name: /^create account$|^criar conta$/i }).click()

    await page.waitForTimeout(3000)

    await page.screenshot({ path: 'tests/e2e/screenshots/03-register-success.png' })

    const toasts = getToasts(page)
    const toastCount = await toasts.count()
    expect(toastCount).toBeGreaterThan(0)

    const toastText = await toasts.first().textContent()
    expect(toastText).toMatch(/success|sucesso|created|realizado/i)

    await expect(page).toHaveURL('/join')
  })

  test('deve exibir toast de erro para email duplicado', async ({ page }) => {
    const uniqueEmail = `dup-${randomUUID()}@example.com`

    await page.getByLabel(/name|nome/i).fill('First User')
    await page.getByLabel(/email/i).fill(uniqueEmail)
    await page.getByLabel(/password|senha/i).fill('password123')
    await page.getByRole('button', { name: /^create account$|^criar conta$/i }).click()

    await page.waitForTimeout(2000)

    await page.getByLabel(/name|nome/i).fill('Second User')
    await page.getByLabel(/email/i).fill(uniqueEmail)
    await page.getByLabel(/password|senha/i).fill('password456')
    await page.getByRole('button', { name: /^create account$|^criar conta$/i }).click()

    await page.waitForTimeout(3000)

    await page.screenshot({ path: 'tests/e2e/screenshots/04-duplicate-email.png' })

    const toasts = getToasts(page)
    const errorToasts = toasts.filter({ hasText: /email|duplicate/i })
    await expect(errorToasts.first()).toBeVisible({ timeout: 5000 })

    await expect(page).toHaveURL('/join')
  })

  test('deve exibir multiplos toasters para campos invalidos', async ({ page }) => {
    await page.evaluate(() => {
      document.querySelectorAll('input[required]').forEach((el) => {
        el.removeAttribute('required')
      })
    })

    await page.getByLabel(/name|nome/i).fill('')
    await page.getByLabel(/email/i).fill('')
    await page.getByLabel(/password|senha/i).fill('123')

    await page.getByRole('button', { name: /^create account$|^criar conta$/i }).click()

    await page.waitForTimeout(3000)

    await page.screenshot({ path: 'tests/e2e/screenshots/05-multiple-errors.png' })

    const toasts = getToasts(page)
    await expect(toasts).toHaveCount(3)

    await expect(page).toHaveURL('/join')
  })

  test('deve exibir toast info ao submeter login', async ({ page }) => {
    await page.getByRole('button', { name: /already have|ja tem/i }).click()

    await page.getByLabel(/email/i).fill('user@example.com')
    await page.getByLabel(/password|senha/i).fill('password123')
    await page.getByRole('button', { name: /^sign in$|^entrar$/i }).click()

    const toasts = getToasts(page)
    await expect(toasts.first()).toBeVisible({ timeout: 5000 })

    const toastText = await toasts.first().textContent()
    expect(toastText).toMatch(/coming soon|em breve/i)

    await page.screenshot({ path: 'tests/e2e/screenshots/06-login-coming-soon.png' })

    await expect(page).toHaveURL('/join')
  })
})
