const puppeteer = require('puppeteer')

async function scrapeWebsite(url) {
  // Launch the browser 
  const browser = await puppeteer.launch({
    headless: false,
  })
  const page = await browser.newPage()

    // navigate to the URL 
    await page.goto(url, {
      waitUntil: 'networkidle2'
    }) 

    await page.waitForSelector('.css-1uwjmhu')
    
    const cardsScrape = await page.evaluate(() => {
      const icons = Array.from(document.querySelectorAll('.css-1uwjmhu'))
      const data = icons.map((el) => ({
        image: el.querySelector('img').getAttribute('src'),
        title: el.querySelector('p').innerText
      }))

      return data
    })

    console.log(cardsScrape)

    // close the browser
    await browser.close()

}

scrapeWebsite('https://www.zehrs.ca/')