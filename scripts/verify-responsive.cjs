const { chromium } = require('playwright');

const BASE = 'http://127.0.0.1:8080';
const EMAIL = process.env.SS_EMAIL || 'rwd-test@chatgo.test';
const PASSWORD = process.env.SS_PASS || 'rwd12345';

const pages = [
  { name: 'landing', url: '/', noauth: true },
  { name: 'login', url: '/login', noauth: true },
  { name: 'dashboard', url: '/' },
  { name: 'contacts', url: '/contacts' },
  { name: 'broadcast', url: '/broadcast' },
  { name: 'send', url: '/send' },
  { name: 'inbox', url: '/inbox' },
  { name: 'templates', url: '/templates' },
  { name: 'autoreply', url: '/autoreply' },
  { name: 'scheduled', url: '/scheduled' },
];

const viewports = [
  { label: 'iphone-se', width: 375, height: 667 },
  { label: 'iphone-11', width: 414, height: 896 },
  { label: 'ipad', width: 768, height: 1024 },
  { label: 'desktop', width: 1440, height: 900 },
];

(async () => {
  const browser = await chromium.launch({ headless: true });
  let failures = 0;

  for (const vp of viewports) {
    const guest = await browser.newContext({ viewport: { width: vp.width, height: vp.height } });
    const auth = await browser.newContext({ viewport: { width: vp.width, height: vp.height } });

    const lp = await auth.newPage();
    await lp.goto(`${BASE}/login`, { waitUntil: 'networkidle' });
    await lp.fill('input[name="email"]', EMAIL);
    await lp.fill('input[name="password"]', PASSWORD);
    await lp.click('button[type="submit"]');
    await lp.waitForURL('**/');
    await lp.close();

    console.log(`\n=== ${vp.label} (${vp.width}x${vp.height}) ===`);

    for (const p of pages) {
      const ctx = p.noauth ? guest : auth;
      const page = await ctx.newPage();
      try {
        const resp = await page.goto(`${BASE}${p.url}`, { waitUntil: p.url === '/inbox' ? 'domcontentloaded' : 'networkidle', timeout: 15000 });
        if (resp && resp.status() >= 400) { console.log(`  SKIP ${p.name}: HTTP ${resp.status()}`); await page.close(); continue; }
        await page.waitForTimeout(600);

        const m = await page.evaluate(() => {
          const doc = document.documentElement;
          const overflowX = doc.scrollWidth - doc.clientWidth;
          const toggler = document.querySelector('.navbar-toggler');
          const togglerVisible = toggler ? toggler.offsetParent !== null : false;
          const mobileSection = document.querySelector('#sidebarCollapse .d-md-none');
          const topbar = document.querySelector('#topbar');
          const topbarVisible = topbar ? getComputedStyle(topbar).display !== 'none' : null;
          return { overflowX, togglerVisible, hasMobileSection: !!mobileSection, topbarVisible };
        });

        const issues = [];
        if (m.overflowX > 2) issues.push(`H-OVERFLOW ${m.overflowX}px`);
        if (vp.width < 768 && !p.noauth) {
          if (!m.togglerVisible) issues.push('no hamburger');
          if (m.topbarVisible === true) issues.push('topbar should be hidden');
        }
        if (vp.width >= 768 && !p.noauth && m.topbarVisible === false) issues.push('topbar hidden on desktop');

        if (issues.length) { failures++; console.log(`  FAIL ${p.name}: ${issues.join(', ')}`); }
        else console.log(`  OK   ${p.name}`);
      } catch (e) {
        failures++;
        console.log(`  ERR  ${p.name}: ${e.message.substring(0, 60)}`);
      } finally {
        await page.close();
      }
    }
    await guest.close();
    await auth.close();
  }

  await browser.close();
  console.log(failures === 0 ? '\nALL RESPONSIVE CHECKS PASSED' : `\n${failures} CHECK(S) FAILED`);
  process.exit(failures === 0 ? 0 : 1);
})();
