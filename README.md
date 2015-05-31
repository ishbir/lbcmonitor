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

The app is intended to be a quick and dirty hack for getting the job done. You
will most likely need to modify output fields and/or some logic based on how you
want to see things.