# Citation Scanner

Citation Scanner is a Go-based project that uses a web scraper and ChatGPT to parse web pages, extract claims, and associate these claims with their respective sources in JSON format. It is especially useful for summarizing content and identifying sources within large web articles, and recursively discovering the quotation for said sources from their original paper/article, as well as identifying potentially erroneous circular claims being made or claims made with no associating source material. It also includes an API server for exposing these features as RESTful endpoints, and an API key generation utility for secure access. It also includes an API server for exposing these features as RESTful endpoints, and an API key generation utility for secure access.

## Features
- **Web Scraping**: Scrapes the content of web pages using the `webscraper` package.
- **OpenAI Integration**: Uses OpenAI API to analyze and extract claims and sources from the scraped content.
- **API Server**: Exposes RESTful API endpoints for parsing pages, using the `chi` router to manage routes.
- **Claims and Sources**: Returns the claims made in an article along with their corresponding sources in JSON format.
- **API Key Management**: Secure API key generation with HMAC for authentication, using a utility in the `cmd/keygen` folder.

## Requirements
- [Go](https://golang.org/dl/) 1.16 or later.
- [OpenAI API Key](https://openai.com/api/)
- [Git](https://git-scm.com/) (for cloning the repository)
- [Chi Router](https://github.com/go-chi/chi) for managing API routes

## Installation

1. **Clone the repository**
   ```sh
   git clone https://github.com/yourusername/citation-scanner.git
   cd citation-scanner
   ```

2. **Install Dependencies**
   Ensure you have Go installed. Then, run:
   ```sh
   go mod tidy
   ```

3. **Set Up Environment Variables**
   Create a `.env` file in the `configs/` directory to store your API keys and secret phrase.
   ```
   OPENAI_API_KEY=your_openai_api_key_here
   SECRET_PHRASE=your_secret_phrase_here
   ```

## Usage

### Running the API Server

To start the API server:

```sh
cd cmd/app
go run main.go
```

The API server will start on port `4145` by default, and provide the following endpoints:

- **GET /**: Basic health check endpoint.
- **POST /parse**: Accepts a JSON payload with a `url` parameter to parse claims from the provided webpage.

Example request to parse a page:
```sh
curl -X POST http://localhost:4145/parse -H "Content-Type: application/json" -d '{"url": "https://en.wikipedia.org/wiki/Go_(programming_language)"}'
```

### Generating an API Key
To generate an API key, you can use the key generation tool located under `cmd/keygen`.

```sh
cd cmd/keygen
go run main.go
```

Make sure you have the `SECRET_PHRASE` set in your `.env` file. This command will generate an API key that is required for accessing the `/parse` endpoint of the API server.

## Testing

You can test the functionality of the project using Go's built-in test tool:

```sh
go test -v ./internal/parser
```

This will run the tests located in the `parser` package and print verbose output, which can help with debugging.

## Example
The following is an example of how the project works:

1. **Scrape a Wikipedia Page**: Given a URL (e.g., `https://en.wikipedia.org/wiki/Go_(programming_language)`), the `webscraper` package will scrape the text content.
2. **Extract Claims and Sources**: The `parser` package then uses OpenAI to extract claims and their linked sources, returning a JSON response such as:

   ```json
   {
     "page": "https://en.wikipedia.org/wiki/Go_(programming_language)",
     "claims": [
       {
         "claim": "Go is a statically typed, compiled programming language.",
         "source": "[12]"
       },
       {
         "claim": "Go was designed at Google in 2009.",
         "source": "[4]"
       }
     ]
   }
   ```

## Project Configuration
- **Configs**: The project uses a `.env` file to store configuration like API keys and secret phrases. Make sure to add your OpenAI key and secret phrase in the `.env` file located under the `configs/` directory.

## API Key Security

The API server uses an HMAC-based authentication mechanism to ensure secure access. To interact with the `/parse` endpoint, you need to include a valid API key in the headers of your request:

- **Header**: `X-API-Key`
- **API Key Generation**: Use the key generation script in `cmd/keygen` to generate an API key using a secret phrase.

## Troubleshooting
1. **OpenAI Response Parsing Issues**:
   - If you receive errors related to JSON parsing (`invalid character looking for beginning of value`), make sure the prompt is correctly formatted and you are requesting only JSON-formatted responses from the API.
2. **Environment Variables**:
   - Ensure that the `.env` file is correctly loaded and contains valid values for `OPENAI_API_KEY` and `SECRET_PHRASE`.
3. **API Key Validation**:
   - Ensure that the `X-API-Key` in your request headers matches the key generated by the `cmd/keygen` tool.

