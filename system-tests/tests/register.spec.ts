import { test, expect } from '@playwright/test';
import { dropTestUsers } from './page/admin';
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await dropTestUsers(page);
});

test('Has correct title', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/register');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle('Register');
});

test('login flow', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/register');

  // Click the get started link.
  await page.getByLabel('Username').fill("test-user");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByLabel('Confirm Password').fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);
});
