import { test, expect } from '@playwright/test';
import { asAdmin, dropTestUsers, addTestUser } from './page/admin';
test.describe.configure({ mode: 'serial' });

test.beforeEach(async ({ page }) => {
  await asAdmin(page, async () => {
    await dropTestUsers(page, 'user');
  });
});

test('reset flow', async ({ page }) => {
  let password;
  await asAdmin(page, async () => {
    password = await addTestUser(page, 'user-reset');
  });

  // Login as the new user.
  await page.getByLabel('Username').fill('user-reset');
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);

  // Reset the password.
  await page.goto('http://ponglehub.com.localhost/auth/user/reset');

  await expect(page).toHaveTitle('Reset Password');

  const newPassword = 'Password2?';

  await page.getByLabel('Password', { exact: true }).fill(newPassword);
  await page.getByLabel('Confirm Password').fill(newPassword);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Password Reset');

  // Login with the new password.
  await page.goto('http://ponglehub.com.localhost/auth/login');

  await expect(page).toHaveTitle('Login');

  await page.getByLabel('Username').fill('user-reset');
  await page.getByLabel('Password').fill(newPassword);
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);
});
