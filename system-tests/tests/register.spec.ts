import { test, expect } from '@playwright/test';

test.beforeAll(async () => {
  const res = await fetch("http://localhost:3000/users");
  console.log(res);
});

test('Has correct title', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/register');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle('Register');
});

const FAIL_TEST = new Error("shouldn't have reached here");

test('login flow', async ({ page }) => {
  await page.goto('http://ponglehub.localhost/auth/register');

  // Click the get started link.
  await page.getByLabel('Username').fill("benofbenton3");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByLabel('Confirm Password').fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle(/Welcome to nginx/);
});
