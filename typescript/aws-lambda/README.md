# AWS Lambda with TypeScript

This project demonstrates a simple AWS Lambda function written in TypeScript.

## Prerequisites

- Node.js (v18 or later)
- AWS CLI configured with appropriate credentials
- AWS SAM CLI (for deployment)

## Project Structure

```
.
├── src/
│   └── index.ts        # Lambda function handler
├── dist/               # Compiled JavaScript files
├── package.json        # Project dependencies
├── tsconfig.json       # TypeScript configuration
└── README.md          # This file
```

## Setup

1. Install dependencies:

```bash
npm install
```

2. Build the project:

```bash
npm run build
```

## Development

- The main Lambda handler is in `src/index.ts`
- Build TypeScript files: `npm run build`
- Run tests: `npm test`

## Deployment

To deploy the Lambda function:

```bash
npm run deploy
```

This will build the TypeScript files and deploy using AWS SAM.

## API Endpoints

The Lambda function responds to HTTP requests with the following structure:

- Method: Any HTTP method
- Response: JSON object containing:
  - message: Welcome message
  - path: Request path
  - method: HTTP method used
  - timestamp: Current timestamp

## Error Handling

The function includes basic error handling and will return:

- 200: Successful response
- 500: Internal server error with error details
