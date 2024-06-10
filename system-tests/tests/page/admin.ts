import { Page } from 'playwright';
import { expect } from '@playwright/test';

export async function dropTestUsers(page: Page) {
    await page.goto('http://ponglehub.localhost/auth/login');

    await expect(page).toHaveTitle('Login')
    await page.getByLabel('Username').fill("admin");
    await page.getByLabel('Password').fill("Password1!");
    await page.getByRole('button').click();
    
    await expect(page).toHaveTitle('Admin')
  
    let running = true;
    while (running) {
      running = false;
  
      let rows = (await page.getByTestId('data-row').all()).map(row => row.locator('td:first-child').innerText());
      let names = await Promise.all(rows);
  
      for (let name of names) {
        if (name.startsWith('test-')) {
          await page.getByTestId('data-row').filter({ hasText: name }).locator('button:has-text("Delete")').click();
          await page.waitForLoadState('networkidle');
          running = true;
          break;
        }
      }
    }
}