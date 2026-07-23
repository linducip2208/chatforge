const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const BASE = 'http://127.0.0.1:8080';
const OUT = path.join(__dirname, '..', 'web', 'screens-mobile');
const EMAIL = process.env.SS_EMAIL || 'rwd-test@chatgo.test';
const PASSWORD = process.env.SS_PASS || 'rwd12345';

const pages = [
  { name: 'm01-landing', url: '/', noauth: true },
  { name: 'm02-login', url: '/login', noauth: true },
  { name: 'm03-dashboard', url: '/' },
  { name: 'm04-sidebar-open', url: '/', menu: true },
  { name: 'm05-contacts', url: '/contacts' },
  { name: 'm06-broadcast', url: '/broadcast' },
  { name: 'm07-inbox', url: '/inbox' },
  { name: 'm08-send', url: '/send' },
];

(async () => {
  if (!fs.existsSync(OUT)) fs.mkdirSync(OUT, { recursive: true });

  const browser = await chromium.launch({ headless: true });

  const guest = await browser.newContext({
    viewport: { width: 414, height: 896 },
    deviceScaleFactor: 2,
    isMobile: true,
    hasTouch: true,
  });
  const auth = await browser.newContext({
    viewport: { width: 414, height: 896 },
    deviceScaleFactor: 2,
    isMobile: true,
    hasTouch: true,
  });

  const loginPage = await auth.newPage();
  await loginPage.goto(`${BASE}/login`, { waitUntil: 'networkidle' });
  await loginPage.fill('input[name="email"]', EMAIL);
  await loginPage.fill('input[name="password"]', PASSWORD);
  await loginPage.click('button[type="submit"]');
  await loginPage.waitForURL('**/');
  await loginPage.close();

  for (const p of pages) {
    console.log(`Capturing: ${p.name}...`);
    const ctx = p.noauth ? guest : auth;
    const page = await ctx.newPage();
    try {
      await page.goto(`${BASE}${p.url}`, { waitUntil: p.url === '/inbox' ? 'domcontentloaded' : 'networkidle', timeout: 15000 });
      await page.waitForTimeout(800);
      if (p.menu) {
        await page.click('.navbar-toggler');
        await page.waitForTimeout(600);
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
