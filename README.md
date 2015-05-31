LBCMonitor
==========

LBCMonitor is a simple Go app for getting your money's worth of Bitcoins if
you're trading on LocalBitcoins.

# How does it work?

Change URL (firstPage) in lbcmonitor.go to whatever page of LBC that you want to
monitor appended with .json. Also change the variable selling to true/false
appropriately. Examples:

```
firstPage="https://localbitcoins.com/buy-bitcoins-online/US/united-states/cash-deposit/.json"
selling=false
```
```
firstPage="https://localbitcoins.com/sell-bitcoins-online/US/united-states/cash-deposit/.json"
selling=true
```
```
firstPage="https://localbitcoins.com/sell-bitcoins-online/sepa-eu-bank-transfer/.json"
selling=true
```

Compile the app and run it as `lbcmonitor --help` to see help arguments. There
are two modes of operation. One is to track the price of Bitcoin and report the
most favourable values (prices higher than threshold specified by --xbtprice if
selling, and lower if buying). The second is to automatically calculate the rate
from the given values of fiat and XBT and report good deals.

The app is intended to be a quick and dirty hack for getting the job done. You
will most likely need to modify output fields and/or some logic based on how you
want to see things.

# License

Do whatever the hell you want.
