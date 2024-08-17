const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  
  // Go to the Amazon page
  await page.goto('https://www.amazon.com/s?k=gaming+chairs&_encoding=UTF8', { waitUntil: 'networkidle2' });

  const chairData = await page.evaluate(() => {
    // Array to hold the extracted data
    let results = [];

    // Select all the product elements
    let items = document.querySelectorAll('.s-result-item[data-component-type="s-search-result"]');

    // Loop over each item and extract the necessary data
    items.forEach((item) => {
      // Extract the title
      let title = item.querySelector('h2 a span') ? item.querySelector('h2 a span').innerText : null;
      
      // Extract the image URL
      let imageUrl = item.querySelector('.s-image') ? item.querySelector('.s-image').src : null;
      
      // Extract the description (if available)
      let description = item.querySelector('.a-row.a-size-base.a-color-secondary') 
        ? item.querySelector('.a-row.a-size-base.a-color-secondary').innerText 
        : 'No description available';

      // Push the extracted data into the results array
      if (title && imageUrl) {
        results.push({ title, imageUrl, description });
      }
    });

    return results;
  });

  console.log(chairData);

  await browser.close();
})();