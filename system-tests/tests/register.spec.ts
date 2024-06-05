import { test, expect } from '@playwright/test';
test.describe.configure({ mode: 'serial' });

import { Client } from './test-client';
const client = new Client("http://localhost:3000");

test.beforeAll(async () => {
  const users = await client.getUsers();
  console.log(users);
  for (let i = 0; i < users.length; i++) {
    console.log(`Deleting user ${users[i].id}`)
    await client.deleteUser(users[i].id);
  }
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
