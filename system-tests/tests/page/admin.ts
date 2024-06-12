import { Page } from "playwright";
import { expect } from "@playwright/test";

export async function addTestUser(page: Page, username: string): Promise<string> {
  await page.goto('http://ponglehub.localhost/auth/login');

  // Click the get started link.
  await page.getByLabel('Username').fill("admin");
  await page.getByLabel('Password', { exact: true }).fill("Password1!");
  await page.getByRole('button').click();

  await expect(page).toHaveTitle('Admin');

  // Click the invite link.
  await page.getByRole('link', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invite User');

  // Fill out the form.
  await page.getByLabel('Username').fill(username);
  await page.getByRole('button', { name: 'Invite' }).click();

  await expect(page).toHaveTitle('Invited User');
  return await page.getByTestId('password').innerText();
}

export async function dropTestUsers(page: Page) {
  await page.goto("http://ponglehub.localhost/auth/login");

  await expect(page).toHaveTitle("Login");
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("Password1!");
  await page.getByRole("button").click();

  await expect(page).toHaveTitle("Admin");

  let running = true;
  while (running) {
    running = false;

    let rows = (await page.getByTestId("data-row").all()).map((row) =>
      row.locator("td:first-child").innerText()
    );
    let names = await Promise.all(rows);

    for (let name of names) {
      if (name.startsWith("test-")) {
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
