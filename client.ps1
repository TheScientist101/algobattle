try {
    # Set console colors for better readability
    $host.UI.RawUI.BackgroundColor = "Black"
    $host.UI.RawUI.ForegroundColor = "White"
    Clear-Host

    Write-Host "================================" -ForegroundColor Cyan
    Write-Host "     ALGOBATTLE API CLIENT     " -ForegroundColor Green
    Write-Host "================================" -ForegroundColor Cyan
    Write-Host

    # Get API key from user
    Write-Host "Please enter your AlgoBattle API key:" -ForegroundColor Yellow
    $SECURE_API_KEY = Read-Host -AsSecureString
    $bstr = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($SECURE_API_KEY)
    $API_KEY = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($bstr)
    # $API_KEY = ConvertFrom-SecureString -SecureString $SECURE_API_KEY -AsPlainText
    Write-Host

    if (-not $API_KEY) {
        Write-Host "API key is required. Exiting." -ForegroundColor Red
        exit 1
    }

    $BASE_URL = "http://algobattle.thescientist101.hackclub.app"

    # Display available tickers
    Write-Host "NOTE: You can only interact with the following stocks:" -ForegroundColor Yellow
    Write-Host "- AAPL (Apple)"
    Write-Host "- GOOG (Google)"
    Write-Host "- JPM (JPMorgan Chase)"
    Write-Host

    # Function to make API requests
    function Invoke-AlgoBattleRequest {
        param(
            [string]$Method,
            [string]$Endpoint,
            [object]$Body = $null,
            [hashtable]$QueryParams = $null
        )

        $headers = @{
            "Authorization" = $API_KEY
        }

        $uri = $BASE_URL + $Endpoint
        
        # Append query parameters if provided
        if ($QueryParams) {
            $queryString = ""
            foreach ($param in $QueryParams.GetEnumerator()) {
                if ($queryString -eq "") {
                    $queryString = "?"
                } else {
                    $queryString += "&"
                }
                $queryString += "$($param.Key)=$($param.Value)"
            }
            $uri += $queryString
        }

        $params = @{
            Method = $Method
            Uri = $uri
            Headers = $headers
            ContentType = "application/json"
        }
        
        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json)
        }
        
        try {
            $response = Invoke-RestMethod @params
            Write-Host "RESPONSE:" -ForegroundColor Green
            Write-Host ($response | ConvertTo-Json -Depth 10)
            return $response
        } catch {
            Write-Host "ERROR: $_" -ForegroundColor Red
            # try {
            #     $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
            #     $responseBody = $reader.ReadToEnd()
            #     Write-Host $responseBody -ForegroundColor Red
            # } catch {
            #     # Ignore any errors trying to read response body
            # }
        }
    }

    # Show menu and handle user input
    while ($true) {
        Write-Host
        Write-Host "Available requests:" -ForegroundColor Cyan
        Write-Host "1. Add Tickers To Fetch Data (GET)"
        Write-Host "2. Buy Stocks (POST)"
        Write-Host "3. Sell Stocks (POST)"
        Write-Host "4. Get Portfolio (GET)"
        Write-Host "5. Get Daily Stock Data (GET)"
        Write-Host "6. Get Live Stock Data (GET)"
        Write-Host "0. Exit"
        Write-Host
        
        Write-Host "Enter your choice (0-6):"
        $choice = Read-Host
        
        switch ($choice) {
            "1" {
                Write-Host "Adding tickers to fetch data..." -ForegroundColor Green
                Write-Host "This will add AAPL, GOOG, and JPM tickers to the server."
                $params = @{
                    "ticker" = @("AAPL&ticker=GOOG&ticker=JPM")
                }
                Invoke-AlgoBattleRequest -Method "GET" -Endpoint "/add_ticker" -QueryParams $params
            }
            
            "2" {
                Write-Host "Buy stocks" -ForegroundColor Green
                Write-Host "Enter ticker (AAPL, GOOG, or JPM):"
                $ticker = Read-Host
                
                # Validate ticker
                if ($ticker -notmatch "^(AAPL|GOOG|JPM)$") {
                    Write-Host "Invalid ticker. Only AAPL, GOOG, or JPM are allowed." -ForegroundColor Yellow
                    continue
                }
                
                Write-Host "Enter number of shares to buy:"
                $shares = Read-Host
                if ($shares -notmatch "^\d+$") {
                    Write-Host "Invalid input. Please enter a valid number." -ForegroundColor Yellow
                    continue
                }
                
                $body = @{
                    action = "buy"
                    numShares = [int]$shares
                    ticker = $ticker
                }
                
                Invoke-AlgoBattleRequest -Method "POST" -Endpoint "/transact" -Body $body
            }
            
            "3" {
                Write-Host "Sell stocks" -ForegroundColor Green
                Write-Host "Enter ticker (AAPL, GOOG, or JPM):"
                $ticker = Read-Host
                
                # Validate ticker
                if ($ticker -notmatch "^(AAPL|GOOG|JPM)$") {
                    Write-Host "Invalid ticker. Only AAPL, GOOG, or JPM are allowed." -ForegroundColor Yellow
                    continue
                }
                
                Write-Host "Enter number of shares to sell:"
                $shares = Read-Host
                if ($shares -notmatch "^\d+$") {
                    Write-Host "Invalid input. Please enter a valid number." -ForegroundColor Yellow
                    continue
                }
                
                $body = @{
                    action = "sell"
                    numShares = [int]$shares
                    ticker = $ticker
                }
                
                Invoke-AlgoBattleRequest -Method "POST" -Endpoint "/transact" -Body $body
            }
            
            "4" {
                Write-Host "Getting your portfolio..." -ForegroundColor Green
                Invoke-AlgoBattleRequest -Method "GET" -Endpoint "/portfolio"
            }
            
            "5" {
                Write-Host "Getting daily stock data (BETA)..." -ForegroundColor Green
                Invoke-AlgoBattleRequest -Method "GET" -Endpoint "/daily_stock_data"
            }
            
            "6" {
                Write-Host "Getting live stock data..." -ForegroundColor Green
                Invoke-AlgoBattleRequest -Method "GET" -Endpoint "/live_stock_data"
            }
            
            "0" {
                Write-Host "Exiting..."
                exit 0
            }
            
            default {
                Write-Host "Invalid choice. Please try again." -ForegroundColor Yellow
            }
        }
        
        Write-Host
        Write-Host "Press Enter to continue..."
        Read-Host
        Clear-Host
    }
} catch {
    # Global error handler for any unexpected errors
    Write-Host "`n`n" -ForegroundColor Red
    Write-Host "=======================================" -ForegroundColor Red
    Write-Host "   An unexpected error has occurred   " -ForegroundColor Red
    Write-Host "=======================================" -ForegroundColor Red
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "`nPress Enter to exit..." -ForegroundColor Yellow
    Read-Host
    exit 1
} finally {
    # Clean up resources if needed
    # This will run whether there was an error or not
}
