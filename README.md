# Minter Seed-phrase Auth

## About
Server for generating seed phrases and authorization using a seed phrase in the [Minter](https://minter.network) blockchain.
Server listens to port: 3999.

### API

#### New Seed-phrase

GET:
http://localhost:3999/api/v1/newMnemonic

Result:
```
{
 "status":true,
 "mnemonic":"mushroom urban cruel bone sting cash office glide impact twin finger bless",
 "err_msg":""
}
```

#### Seed-phrase authorization

GET:
http://localhost:3999/api/v1/authSeed?sp=mushroom%20urban%20cruel%20bone%20sting%20cash%20office%20glide%20impact%20twin%20finger%20bless

or

POST:
http://localhost:3999/api/v1/authSeed

Body-post:
{
  "sp":"mushroom urban cruel bone sting cash office glide impact twin finger bless"
}


Result:
```
{
 "status":true,
 "mnemonic":"mushroom urban cruel bone sting cash office glide impact twin finger bless",
 "err_msg":"",
 "address":"Mx7c09f6b08d175a5d799c876db6849ec79e5719fa",
 "priv_key":"b69361e151b563ceb83b401b3ee5bca49ef433186eb12e454a5cb5fe62b0c624"
}
```
