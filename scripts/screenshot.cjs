const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE = 'http://127.0.0.1:8080';
const OUT = path.join(__dirname, '..', 'public', 'marketing', 'screens');
const EMAIL = 'admin@chatgo.test';
const PASSWORD = 'password';

const pages = [
  { name: '01-login', url: '/login' },
  { name: '02-dashboard', url: '/' },
  { name: '03-wa-qr', url: '/wa' },
  { name: '04-send-message', url: '/send' },
  { name: '05-broadcast', url: '/broadcast' },
  { name: '06-scheduled', url: '/scheduled' },
  { name: '07-autoreply-ai', url: '/autoreply' },
  { name: '08-autoreply-rules', url: '/autoreply' },
  { name: '09-contacts', url: '/contacts' },
  { name: '10-templates', url: '/templates' },
  { name: '11-settings', url: '/settings' },
  { name: '12-inbox', url: '/inbox' },
  { name: '13-docs', url: '/docs' },
  { name: '14-admin-users', url: '/admin/users' },
  { name: '15-admin-packages', url: '/admin/packages' },
];

(async () => {
  if (!fs.existsSync(OUT)) fs.mkdirSync(OUT, { recursive: true });

  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: { width: 1440, height: 900 },
    deviceScaleFactor: 1,
  });

  // Login
  const loginPage = await context.newPage();
  await loginPage.goto(`${BASE}/login`, { waitUntil: 'networkidle' });
  await loginPage.fill('input[name="email"]', EMAIL);
  await loginPage.fill('input[name="password"]', PASSWORD);
  await loginPage.click('button[type="submit"]');
  await loginPage.waitForURL('**/');
  await loginPage.close();

  // Capture all pages
  for (const p of pages) {
    console.log(`Capturing: ${p.name}...`);
    const page = await context.newPage();
    try {
      await page.goto(`${BASE}${p.url}`, { waitUntil: 'networkidle', timeout: 15000 });
      await page.waitForTimeout(1000);
      // For autoreply, also click the Rules tab
      if (p.name === '08-autoreply-rules') {
        await page.click('.nav-tabs a[href="#tab-rules"]');
        await page.waitForTimeout(500);
      }
      if (p.name === '07-autoreply-ai') {
        // AI Config tab is default active, no need to click
        await page.waitForTimeout(500);
      }
      await page.screenshot({ path: path.join(OUT, `${p.name}.png`), fullPage: false });
      console.log(`  OK: ${p.name}.png`);
    } catch (e) {
      console.log(`  FAILED: ${p.name} - ${e.message.substring(0, 80)}`);
    } finally {
      await page.close();
    }
  }

  await browser.close();
  console.log('\nDone! Screenshots saved to:', OUT);
})();
