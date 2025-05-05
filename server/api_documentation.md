# AlgoBattle API Documentation

This document provides detailed information on how to make requests to the AlgoBattle trading platform API. The API allows you to manage portfolios, execute trades, and retrieve stock data.

## Authentication

All API endpoints require authentication using an API key. The API key should be provided in the `Authorization` header of each request.

Example:
```
Authorization: your_api_key_here
```

## Endpoints

### Portfolio Management

#### Get Portfolio

Retrieves the authenticated user's portfolio including cash balance, holdings, and transaction history.

- **URL**: `/portfolio`
- **Method**: `GET`
- **Authentication**: Required

**Example Request:**
```http
GET http://localhost:8080/portfolio
Authorization: your_api_key_here
```

**Example Response:**
```json
{
  "type": "portfolio",
  "payload": {
    "accountValue": 10500.25,
    "historicalAccountValue": [
      {
        "date": "2023-01-01T00:00:00Z",
        "value": 10000.00
      },
      {
        "date": "2023-01-02T00:00:00Z",
        "value": 10500.25
      }
    ],
    "cash": 5000.25,
    "holdings": {
      "AAPL": {
        "numShares": 10,
        "purchaseValue": 150.00
      },
      "GOOG": {
        "numShares": 5,
        "purchaseValue": 1000.00
      }
    },
    "transactions": [
      {
        "time": "2023-01-01T12:00:00Z",
        "numShares": 10,
        "unitCost": 150.00,
        "ticker": "AAPL",
        "action": "buy"
      },
      {
        "time": "2023-01-01T13:00:00Z",
        "numShares": 5,
        "unitCost": 1000.00,
        "ticker": "GOOG",
        "action": "buy"
      }
    ]
  }
}
```

### Stock Data

#### Add Ticker

Adds one or more stock tickers to the watchlist for price monitoring and data collection.

- **URL**: `/add_ticker`
- **Method**: `GET`
- **Authentication**: Required
- **Query Parameters**:
  - `ticker` (array of strings): Ticker symbols to add (can specify multiple)

**Example Request:**
```http
GET http://localhost:8080/add_ticker?ticker=AAPL&ticker=GOOG
Authorization: your_api_key_here
```

**Example Response:**
```json
{
  "type": "result",
  "payload": {
    "payload": "successfully added tickers: [AAPL GOOG]",
    "success": true
  }
}
```

#### Get Daily Stock Data

Retrieves daily historical stock data for all tickers in the watchlist.

- **URL**: `/daily_stock_data`
- **Method**: `GET`
- **Authentication**: Required

**Example Request:**
```http
GET http://localhost:8080/daily_stock_data
Authorization: your_api_key_here
```

**Example Response:**
```json
{
  "type": "daily_stock_data",
  "payload": {
    "tickers": {
      "AAPL": {
        "dataStart": "2000-01-01T00:00:00Z",
        "dataEnd": "2023-01-01T00:00:00Z"
      },
      "GOOG": {
        "dataStart": "2004-08-19T00:00:00Z",
        "dataEnd": "2023-01-01T00:00:00Z"
      }
    },
    "rows": [
      {
        "date": "2023-01-01T00:00:00Z",
        "data": {
          "AAPL": {
            "open": 150.00,
            "high": 155.00,
            "low": 149.00,
            "close": 152.00,
            "volume": 1000000,
            "adjClose": 152.00,
            "adjHigh": 155.00,
            "adjLow": 149.00,
            "adjOpen": 150.00,
            "adjVolume": 1000000,
            "divCash": 0,
            "splitFactor": 1
          },
          "GOOG": {
            "open": 1000.00,
            "high": 1010.00,
            "low": 990.00,
            "close": 1005.00,
            "volume": 500000,
            "adjClose": 1005.00,
            "adjHigh": 1010.00,
            "adjLow": 990.00,
            "adjOpen": 1000.00,
            "adjVolume": 500000,
            "divCash": 0,
            "splitFactor": 1
          }
        }
      }
    ]
  }
}
```

#### Get Live Stock Data

Retrieves the latest stock prices for all tickers in the watchlist.

- **URL**: `/live_stock_data`
- **Method**: `GET`
- **Authentication**: Required

**Example Request:**
```http
GET http://localhost:8080/live_stock_data
Authorization: your_api_key_here
```

**Example Response:**
```json
{
  "type": "live_stock_data",
  "payload": {
    "AAPL": 152.35,
    "GOOG": 1007.25
  }
}
```

### Transactions

#### Execute Transaction

Processes a buy or sell transaction for a specified ticker and number of shares.

- **URL**: `/transact`
- **Method**: `POST`
- **Authentication**: Required
- **Content-Type**: `application/json`
- **Request Body**:
  - `action` (string): "buy" or "sell"
  - `numShares` (number): Number of shares to buy or sell
  - `ticker` (string): Stock ticker symbol

**Example Request:**
```http
POST http://localhost:8080/transact
Authorization: your_api_key_here
Content-Type: application/json

{
  "action": "buy",
  "numShares": 10,
  "ticker": "AAPL"
}
```

**Example Response:**
```json
{
  "type": "result",
  "payload": {
    "payload": "successfully executed transaction",
    "success": true
  }
}
```

## Error Handling

All API endpoints return appropriate HTTP status codes and error messages in case of failure:

- `200 OK`: Request successful
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Authentication failed or insufficient permissions
- `500 Internal Server Error`: Server-side error

Error responses follow the same format as success responses, but with `success` set to `false` and an error message in the `payload` field.

Example error response:
```json
{
  "type": "result",
  "payload": {
    "payload": "error: not enough cash to buy 10 shares of AAPL",
    "success": false
  }
}
```