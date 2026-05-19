import { test, expect } from '@playwright/test';

const ASSETS_DIR = '../../assets';

async function takeScreenshot(page: any, name: string) {
  await page.setViewportSize({ width: 1280, height: 900 });
  await page.waitForTimeout(2000);
  await page.screenshot({ path: `${ASSETS_DIR}/${name}.png` });
}

async function registerAndLogin(page: any) {
  const email = `screenshot-${Date.now()}@test.com`;
  const password = 'testpass123';

  await page.goto('/join');
  await page.fill('input[type="text"]', 'Screenshot User');
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
  await expect(page.getByText('Account created successfully!')).toBeVisible({ timeout: 10000 });

  await page.click('button:has-text("Sign in")');
  await page.waitForTimeout(300);

  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
  await page.waitForURL('/dashboard', { timeout: 10000 });
  await expect(page.getByRole('button', { name: /Screenshot User/ })).toBeVisible({ timeout: 10000 });
}

test.describe('Screenshots', () => {
  test('01-home', async ({ page }) => {
    await page.goto('/home');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '01-home');
  });

  test('02-join-register', async ({ page }) => {
    await page.goto('/join');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '02-join-register');
  });

  test('03-join-login', async ({ page }) => {
    await page.goto('/join');
    await page.waitForTimeout(1000);
    await page.click('button:has-text("Sign in")');
    await page.waitForTimeout(500);
    await takeScreenshot(page, '03-join-login');
  });

  test('04-dashboard', async ({ page }) => {
    await registerAndLogin(page);
    await takeScreenshot(page, '04-dashboard');
  });

  test('05-users-list', async ({ page }) => {
    await registerAndLogin(page);
    await page.goto('/users');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '05-users-list');
  });

  test('06-users-new', async ({ page }) => {
    await registerAndLogin(page);
    await page.goto('/users/new');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '06-users-new');
  });

  test('07-products-list', async ({ page }) => {
    await registerAndLogin(page);
    await page.goto('/catalog/products');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '07-products-list');
  });

  test('08-products-new', async ({ page }) => {
    await registerAndLogin(page);
    await page.goto('/catalog/products/new');
    await page.waitForTimeout(2000);
    await takeScreenshot(page, '08-products-new');
  });
});
