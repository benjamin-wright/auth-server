import { test, expect } from '@playwright/test';
import { dropTestUsers } from './page/admin';
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await dropTestUsers(page);
});

test('protect self', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/login');

  // Click the get started link.
  await page.getByLabel('Username').fill("admin");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');

  const adminRow = page.getByTestId("data-row").filter({ has: page.getByText('admin', { exact: true }) });
  await expect(adminRow.locator('input[type="checkbox"]')).toBeDisabled();
  await expect(adminRow.locator('button:has-text("Delete")')).toBeDisabled();
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

test('promote to admin flow', async ({ page }) => {
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
  await page.getByLabel('Username').fill('test-admin-user');
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invited User');
  const password = await page.getByTestId('password').innerText();

  await page.getByRole('button', { name: 'Done' }).click();
  await expect(page).toHaveTitle('Admin');

  // Click the admin checkbox for the test user.
  page.getByTestId("data-row").filter({ hasText: 'test-admin-user' }).locator('input[type="checkbox"]').check();
  await page.waitForLoadState("networkidle");

  await page.getByRole('link', { name: 'Logout' }).click();
  await expect(page).toHaveTitle('Login');

  // Login as the new user.

  await page.getByLabel('Username').fill('test-admin-user');
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');
});