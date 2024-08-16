const puppeteer = require('puppeteer')
const fs = require('fs')

async function scrapeWebsite(url) {
  // Launch the browser 
  const browser = await puppeteer.launch({
    headless: false,
  })
  const page = await browser.newPage()

    // navigate to the URL 
    await page.goto(url) 
    
    const cardsScrape = await page.evaluate(() => {
      const icons = Array.from(document.querySelectorAll('.category'))
      const data = icons.map((el) => ({
        image: el.querySelector('img').getAttribute('src'),
        title: el.querySelector('h3 a').innerText
      }))

      return data
    })

    console.log(cardsScrape)

    // close the browser
    await browser.close()

    fs.writeFile('categories.json', JSON.stringify(cardsScrape), (err) => {
      if(err) throw err
      console.log("JSON successfully created");
    })

}

scrapeWebsite('https://www.saveonfoods.com/sm/planning/rsid/1982/')