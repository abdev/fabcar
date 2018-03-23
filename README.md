Tendermit implementation of hyperledger fabric fabcar chaincode:
https://github.com/hyperledger/fabric-samples/tree/release/chaincode/fabcar

The state data is persistently stored in a LevelDb database. 

Implemented Operations: **createCar, queryAllCars, changeCarOwner**

# Transactions

To test:  
tx is json base64 encode, the response value in the result is also json base64 encoded

Can use https://www.base64decode.net/ to encode/decode

### 1. CreateCar car1
decoded tx value:  
**_{"operation": "createCar", "data": {"ID": "car1", "Make": "Peugeot", "Model": "205", "Colour": "red", "Owner": "Jane"}_**  
equivalent base64 encoded:
**_eyJvcGVyYXRpb24iOiAi3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjEiLCAiTWFrZSI6ICJQZXVnZW90IiwgIk1vZGVsIjogIjIwNSIsICJDb2xvdXIiOiAicmVkIiwgIk93bmVyIjogIkphbmUifX0=_**

`curl 'http://localhost:46657/broadcast_tx_commit?tx="eyJvcGVyYXRpb24iOiAi3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjEiLCAiTWFrZSI6ICJQZXVnZW90IiwgIk1vZGVsIjogIjIwNSIsICJDb2xvdXIiOiAicmVkIiwgIk93bmVyIjogIkphbmUifX0="'`

### 2. QueryAllCars: 
`curl localhost:46657/abci_query?path="allCars"'`

Sample response:  
{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "response": {
      "value": "eyJDYXI6Y2FyMSI6eyJpZCI6ImNhcjEiLCJtYWtlIjoiUGV1Z2VvdCIsIm1vZGVsIjoiMjA1IiwiY29sb3VyIjoicmVkIiwib3duZXIiOiJKYW5lIn19"
    }
  }
}    
Decoded response value:  
**_{"Car:car1":{"id":"car1","make":"Peugeot","model":"205","colour":"red","owner":"Jane"}}_**

### 3. CreateCar car3
decoded tx value:  
**_{"operation": "createCar", "data": {"ID": "car3", "Make": "Toyota", "Model": "Prius", "Colour": "blue", "Owner": "Tomoko"}}_**

based64 encoded:  
**_eyJvcGVyYXRpb24iOiAiY3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjMiLCAiTWFrZSI6ICJUb3lvdGEiLCAiTW9kZWwiOiAiUHJpdXMiLCAiQ29sb3VyIjogImJsdWUiLCAiT3duZXIiOiAiVG9tb2tvIn19_**

`curl 'http://localhost:46657/broadcast_tx_commit?tx="eyJvcGVyYXRpb24iOiAiY3JlYXRlQ2FyIiwgImRhdGEiOiB7IklEIjogImNhcjMiLCAiTWFrZSI6ICJUb3lvdGEiLCAiTW9kZWwiOiAiUHJpdXMiLCAiQ29sb3VyIjgImJsdWUiLCAiT3duZXIiOiAiVG9tb2tvIn19"'`

### 4. ChangeCarOwner car1
decoded tx value:  
_**{"operation": "changeCarOwner", "data": {"asset_id": "car1", "new_owner": "Robert"}}**_
base64 encoded:  
_**eyJvcGVyYXRpb24iOiAi2hhbmdlQ2FyT3duZXIiLCAiZGF0YSI6IHsiYXNzZXRfaWQiOiAiY2FyMSIsICJuZXdfb3duZXIiOiAiUm9iZXJ0In19**_

`curl 'http://localhost:46657/broadcast_tx_commit?tx="eyJvcGVyYXRpb24iOiAi2hhbmdlQ2FyT3duZXIiLCAiZGF0YSI6IHsiYXNzZXRfaWQiOiAiY2FyMSIsICJuZXdfb3duZXIiOiAiUm9iZXJ0In19"'`

### 5. QueryAllCars: 
`curl localhost:46657/abci_query?path="allCars"'`

Sample response:  
{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "response": {
      "value": "eyJDYXI6Y2FyMSI6eyJpZCI6ImNhcjEiLCJtYWtlIjoiUGV1Z2VvdCIsIm1vZGVsIjoiMjA1IiwiY29sb3VyIjoicmVkIiwib3duZXIiOiJSb2JlcnQifSwiQ2FyOmNhcjMiOnsiaWQiOiJjYXIzIiwibWFrZSI6IlRveW90YSIsIm1vZGVsIjoiUHJpdXMiLCJjb2xvdXIiOiJibHVlIiwib3duZXIiOiJUb21va28ifX0="
    }
  }
}  
Decoded response value:  
**_{"Car:car1":{"id":"car1","make":"Peugeot","model":"205","colour":"red","owner":"Robert"},"Car:car3":{"id":"car3","make":"Toyota","model":"Prius","colour":"blue","owner":"Tomoko"}}_**