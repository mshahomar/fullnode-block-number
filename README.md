# FullNode Block Number

This tool mainly to query hosted FullNode for latest block number for ETH, BSC and TRX, which will compare the result with the latest block from etherscan.io, bscscan.com, and trongrid.io.
  

## Usage

The assumption here is that the hosted FullNode requires authentication, and this is read from .env file. So ensure the following is available in your env

```env
  TRONONE=<YOUR PRIMARY FULLNODE ENDPOINT>
  TRONTWO=<YOUR SECONDARY FULLNODE ENDPOINT>
  TRONUSER=<YOUR TRX FULLNODE LOGIN DETAILS>
  TRONPASS=<YOUR TRX FULLNODE LOGIN DETAILS>
  BNBURL=<YOUR BSC FULLNODE ENDPOINT>
  BNBUSER=<YOUR BSC FULLNODE LOGIN DETAILS>
  BNBPASS=<YOUR BSC FULLNODE LOGIN DETAILS>
  ETHONE=<YOUR PRIMARY ETH FULLNODE ENDPOINT>
  ETHUSER=<YOUR ETH FULLNODE LOGIN DETAILS>
  ETHPASS=<YOUR ETH FULLNODE LOGIN DETAILS>
  TRONGRID=https://api.trongrid.io/<>
  BSCSCAN=https://bscscan.com/
  ETHERSCAN=https://etherscan.io
```
 
