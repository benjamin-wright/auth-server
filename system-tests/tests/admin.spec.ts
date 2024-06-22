import { test, expect } from '@playwright/test';
import { asAdmin, dropTestUsers, addTestUser, toggleAdmin } from './page/admin';
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await asAdmin(page, async () => {
    await dropTestUsers(page, 'admin');
  });
});

test('protect self', async ({ page }) => {
  await asAdmin(page, async () => {
    const adminRow = page.getByTestId("data-row").filter({ has: page.getByText('admin', { exact: true }) });
    await expect(adminRow.locator('input[type="checkbox"]')).toBeDisabled();
    await expect(adminRow.locator('button:has-text("Delete")')).toBeDisabled();
  });
});

test('invite flow', async ({ page }) => {
  let password;
  await asAdmin(page, async () => {
    password = await addTestUser(page, 'admin-normal-user');
  });

  // Login as the new user.
  await page.getByLabel('Username').fill('admin-normal-user');
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);
});

test('promote to admin flow', async ({ page }) => {
  let password;

  await asAdmin(page, async () => {
    password = await addTestUser(page, 'admin-admin-user');
    await toggleAdmin(page, 'admin-admin-user');
  });

  await page.getByLabel('Username').fill('admin-admin-user');
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');
});