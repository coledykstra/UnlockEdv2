const axios = require('axios');
const dotenv = require('dotenv');
const { faker } = require('@faker-js/faker');
const fs = require('fs');
const path = require('path');
const puppeteer = require('puppeteer');

async function postDiscussion(username, password) {
  (async () => {
    // Launch the browser
    const browser = await puppeteer.launch({ headless: false });
    const page = await browser.newPage();

    // Navigate to the Kolibri login page
    await page.goto(process.env.KOLIBRI_URL);
  
    // Perform login
    // Wait for the first textbox to become visible and type username
    await page.waitForSelector('.ui-textbox-input');
    await page.type('.ui-textbox-input', username);

    // Click the next button
    await page.click('.login-btn.button._1q0nwed.raised');

    // Wait for the textbox to become visible and type password
    await page.waitForSelector('.ui-text-inputbox');
    await page.type('.ui-text-inputbox', password);

    // Click the Sign In button
    await page.click('.login-btn.button._1q0nwed.raised');

    // Wait for navigation to complete
    await page.waitForNavigation();

    // Add a 5-second delay before closing the browser
    await page.waitForTimeout(5000);
        
    // Close the browser
    await browser.close();
})();
}

async function delay(time) {
  return new Promise(resolve => setTimeout(resolve, time));
}

async function createCourse(courseName) {
  (async () => {
    // Launch the browser
    const browser = await puppeteer.launch({ headless: false });
    const page = await browser.newPage();

    // Navigate to the Kolibri login page
    console.log('goto KOLIBRI_URL');
    await page.goto(process.env.KOLIBRI_URL);
  
    // Perform login
    // Wait for the first textbox to become visible and type username
    await page.waitForSelector('.ui-textbox-input');
    console.log('typing ADMIN_USERNAME');
    await page.type('.ui-textbox-input', process.env.ADMIN_USERNAME);

    // Click the next button
    console.log('clicking next button');
    await page.click('.login-btn.button._1q0nwed.raised');

    await delay(1000);

    // await page.waitForNavigation();
    // console.log('waitForNavigation complete');

    // Wait for the textbox to become visible and type password
    console.log('waiting for type password');
    await page.waitForSelector('input[type="password"]');
    console.log('found type password');
    await page.type('input[type="password"]', process.env.ADMIN_PASSWORD);
    console.log('password entered');

    // Click the Sign In button
    await page.keyboard.press('Enter');
    // await page.waitForSelector('input[type="submit"]');
    // console.log('Click the Sign In button');
    // await page.click('input[type="submit"]');

    // Wait for navigation to complete
    // await delay(1000);

    // // Navigate to the classes page
    await page.goto(process.env.KOLIBRI_URL + '/en/facility/#/classes', {
      waitUntil: 'networkidle2'
    });

    // // Click the New Class button
    await page.click('.move-down.button._1q0nwed.raised');
    
    await page.waitForSelector('input[type="text"]');
    await page.type('input[type="text"]', courseName);

    await page.keyboard.press('Enter');

    await page.waitForNavigation();

    // Find the link containing the courseName and click it
    await page.evaluate(() => {
      const links = Array.from(document.querySelectorAll('a.link'));
      const targetLink = links.find(link => link.textContent.trim().includes(courseName));
      if (targetLink) {
        targetLink.click();
      }
    });

    // Add a 5-second delay before closing the browser
    await delay(5000);
        
    // Close the browser
    await browser.close();
})();
}

async function main() {
  // Specify the path to the .env file for the local directory
  const envFilePath = path.resolve(__dirname, '.env');

  // Load the .env file
  const result = dotenv.config({ path: envFilePath });

  if (result.error) {
    console.error('Error loading .env file:', result.error);
  } else {
    console.log('Successfully loaded .env file');
  }


  // // Read and parse the tab-delimited credentials file
  // tab_file = path.join(__dirname, process.env.TAB_FILE_PATH)
  // console.log(`Reading the file: ${tab_file}`);
  // fs.readFile(tab_file, 'utf8', (err, data) => {
  //   if (err) {
  //     console.error('Error reading the file:', err);
  //     return;
  //   }

  //   // Split the file content by lines
  //   const lines = data.trim().split('\n');

  //   // Process each line (skipping blank lines and comment lines)
  //   lines.forEach((line, index) => {
  //     line = line.trim();
  //     if (line === '' || line.startsWith('#')) {
  //       return;
  //     }

  //     // Split each line by tab to get username and password
  //     const [username, password] = line.split('\t');
  //     // console.log(`${username} ${password}`);
  //     postDiscussion(username, password);
  //   });

  //   console.log('Tab-delimited text file successfully processed');
  // });

  try {
    const courseName = 'New Course';
    await createCourse(courseName);
    console.log(`Successfully created course`);
  } catch (error) {
    console.error('Error creating course:', error);
  }
}

main();