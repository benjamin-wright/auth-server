import { Page } from "playwright";
import { expect } from "@playwright/test";

export async function asAdmin(page: Page, action: () => Promise<void>) {
  await page.goto('http://ponglehub.com.localhost/auth/login');

  // Click the get started link.
  await page.getByLabel('Username').fill("admin");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');

  await action();

  await page.getByRole('link', { name: 'Logout' }).click();

  await expect(page).toHaveTitle('Login');
}

export async function toggleAdmin(page: Page, username: string) {
  await page.getByTestId("data-row").filter({ hasText: username }).locator('input[type="checkbox"]').check();
  await page.waitForLoadState("networkidle");
}

export async function addTestUser(page: Page, username: string): Promise<string> {
  // Click the invite button.
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invite User');

  // Fill out the form.
  await page.getByLabel('Username').fill(username);
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invited User');
  const password = await page.getByTestId('password').innerText();

  await page.getByRole('button', { name: 'Done' }).click();
  await expect(page).toHaveTitle('Admin');

  return password;
}

export async function dropTestUsers(page: Page, spec: string) {
  let running = true;
  while (running) {
    running = false;

    let rows = (await page.getByTestId("data-row").all()).map((row) =>
      row.locator("td:first-child").innerText()
    );
    let names = await Promise.all(rows);

    for (let name of names) {
      if (name.startsWith(`${spec}-`)) {
        await page
          .getByTestId("data-row")
          .filter({ hasText: name })
          .locator('button:has-text("Delete")')
          .click();
        await page.waitForLoadState("networkidle");
        running = true;
        break;
      }
    }
  }
}
