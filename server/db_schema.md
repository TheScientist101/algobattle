## Algobattle DB Schema

#### /users
This will contain info about users by their auth uuid or whatever, will contain basic stuff, firstName, lastName, email, also an array of UUID of bots owned (shared bot ownership may not initially be supported)

#### /bots
Each document in this will contain information about the bot including the owner, the API key, maybe a name, all current equity and cash, last calculated price (bc api calls are expensive), historical stuff will be in /transactions which we can index by bot UUID…
Now, I am thinking, maybe we should have an array of historical account value for the fancy line graphs but I’m not sure how broken up the data should be, maybe daily? Hourly?
I would appreciate suggestions on this. Current idea would be to go through all bots and value them EOD.

Also, this is a separate subsection because this makes making the leaderboard easier.

#### /transactions
Contains all transactions which will have some bot UUID attribute to see who owns the transaction, the ticker for the stock purchased, and the number of shares you purchased. Note that you may not be able to recalculate the historical account value via this transactions list because of stock splits and such.

How do we manage stock splits?
Do we care about dividends (maybe that’s a later technicality)?