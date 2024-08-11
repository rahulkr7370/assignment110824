# Concurrent Web Scraper

## Overview

This project is a simple web scraper application built in Go. The scraper retrieves the title and the first 100 words of content from a list of web pages. It uses Go's concurrency features, such as goroutines and channels, to fetch data from multiple web pages simultaneously. The scraper also handles timeouts and errors gracefully, ensuring that the program does not crash even if a webpage takes too long to respond or if it is unreachable.

## Features

- **Concurrent Scraping**: Fetches data from multiple web pages at the same time, making it faster and more efficient.
- **Timeout Handling**: Allows you to specify a timeout in milliseconds. If a webpage takes longer than the specified time, it returns a "Timeout" message for that page.
- **Error Handling**: Reports errors without crashing the program if a webpage cannot be retrieved.
- **Custom Output**: Displays the URL, title, and the first 100 words of content for each page.

## Installation

### Prerequisites

- **Go**: Ensure Go is installed on your machine. If not, download it from [Go's official website](https://golang.org/dl/).

### Steps

1. **Clone the Repository**:
   - Open your terminal or command prompt.
   - Run the following command to clone the repository:
     ```bash
     git clone https://github.com/rahulkr7370/assignment110824.git
     ```

2. **Navigate to the Project Directory**:
   - Change into the project directory using:
     ```bash
     cd assignment110824
     ```

3. **Run the Scraper**:
   - Execute the following command to run the scraper:
     ```bash
     go run scrapper.go
     ```
   - The scraper will fetch data from the specified web pages and print the results in the terminal.

## Usage

The scraper automatically fetches the title and content from a list of pre-defined web pages. You can modify the URLs in the `scrapper.go` file to scrape different pages.


## Running Tests

Unit tests are provided to ensure the scraper works as expected. To run these tests:

1. **Run the Tests**:
   - Use the following command:
     ```bash
     go test -v
     ```
   - Test results will be displayed in the terminal.

## Customization

You can customize the scraper by:

1. **Changing the URLs**: Edit the list of URLs in the `scrapper.go` file to scrape different web pages.
2. **Adjusting the Timeout**: Modify the timeout value in the `scrapper.go` file to set how long the scraper waits for a response.

## Troubleshooting

### Common Issues

- **Timeouts**: If you see a "Timeout" message, the scraper couldn't get a response within the specified time. Increase the timeout duration in the `scrapper.go` file.
- **Errors**: An "Error" message indicates issues with reaching the page or other network problems. Check the URLs and your internet connection.


## Conclusion

This project demonstrates how to use Go's concurrency features to build efficient applications. It's straightforward to set up and run, even for beginners. We hope you find it useful for your scraping needs and encourage you to explore the code further!

**Happy Scraping!**


