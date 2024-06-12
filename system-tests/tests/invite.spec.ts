import { test, expect } from '@playwright/test';
import { dropTestUsers } from './page/admin';
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await dropTestUsers(page);
});

test('invite flow', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/login');

  // Click the get started link.
  await page.getByLabel('Username').fill("admin");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');

  // Click the invite link.
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invite User');

  // Fill out the form.
  await page.getByLabel('Username').fill('test-invite-user');
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invited User');
  const password = await page.getByTestId('password').innerText();

  // Logout.
  await page.getByRole('link', { name: 'Logout' }).click();

  await expect(page).toHaveTitle('Login');

  // Login as the new user.
  await page.getByLabel('Username').fill('test-invite-user');
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);
});
