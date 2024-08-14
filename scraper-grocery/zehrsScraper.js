const puppeteer = require('puppeteer')

async function scrapeWebsite(url) {
  // Launch the browser 
  const browser = await puppeteer.launch()
  const page = await browser.newPage()

  try {
    // navigate to the URL 
    await page.goto(url, { timeout: 600000 }) 

    // scrape data 
    const allGroceries = await page.evaluate(() => {
      const siteContent = document.querySelector('.site-content')
      if (!siteContent) {
        console.error('No element with class .site-content found')
        return []
      }

      const groceries = Array.from(siteContent.querySelectorAll('ul'))
      return groceries.map(grocery => {
        const priceElement = grocery.querySelector('span')
        const price = priceElement ? priceElement.innerText : 'N/A'
        return { price }
      })
    })

    console.log(allGroceries)
  } catch (error) {
    console.error('Error during navigation:', error)
  } finally {
    // close the browser
    await browser.close()
  }
}

scrapeWebsite('https://www.zehrs.ca/food/fruits-vegetables/c/28000?navid=flyout-L2-fruits-vegetables&icid=gr_fruits-vegetables-shop-categories_tile_3_hp')