var chai = require('chai')
var expect = chai.expect
const request = require('request')

const rootToken = "11e53801-8741-53cf-22fc-3152de53ea02"

describe('Integration test for user registration', () => {
  it('should return uuid of the user when valid user credentials are provided', (done) => {
    const data = {
      "username": "",
      "passphrase": "",
      "mnemonic": ""
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/register',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.equal(200)
      expect(vaultResponse.data).to.have.property('uuid')
      done();
    })
  })

  it('should return error when invalid mnemonic is provided', (done) => {
    const data = {
      "username": "",
      "passphrase": "",
      "mnemonic": "somethinginvalid"
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/register',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.equal(417)
      expect(vaultResponse).to.have.property('errors')
      done();
    })
  })

})

describe('Integration test for generating bitcoin signature', () => {
  let uuid
  before((done) => {
    const data = {
      "username": "",
      "passphrase": "",
      "mnemonic": "similar supreme jealous custom please fitness mosquito report movie valley hip hub slush foam deer"
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/register',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      uuid = vaultResponse.data.uuid;
      done();
    })
  })

  it('should return a valid signature when valid uuid, path, payload and cointype are provided', (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: 0
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: 91234
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    const actualSignature = "01000000018ad40444bde5096a2dc58f5aa278ab4c2a41eb52975857f96fb50cd732c8b481000000006a473044022064deb4f6bd3d283368e0eba6ac00f19a7412d01c8c9ff729bd30630bb0c0592502200134f0badc1796dcc7df13932f18a66454ac8558cafe3cfc82ca6f0b0200fd9b012103cd11c3e23a78a041c004ca3410575b688147ddecdf3e5931e0dda23192c8dcc7ffffffff0162640100000000001976a914c8e90996c7c6080ee06284600c684ed904d14c5c88ac00000000"

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.equal(200)
      expect(vaultResponse.data.signature).to.equal(actualSignature)
      done();
    })
  })

  it("should return error when uuid provided is not there in the database", (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: 0
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: 91234
      }]
    }

    const data = {
      "uuid": "some-other-uuid",
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("UUID does not exists")
      done();
    })

  })

  it("should return error when invalid payee address is provided", (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: 0
      }],
      outputs: [{
        address: "123456", //invalid address format
        amount: 91234
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Invalid payee address format")
      done();
    })
  })

  it("should return error when invalid UTXO id is provided", (done) => {
    const payload = {
      inputs: [{
        txhash: "f639hfa", //invalid format
        vout: 0
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: 91234
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Invalid UTXO hash")
      done();
    })
  })

  it('should return error when negative Vout is provided', (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: -1 //invalid
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: 91234
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include('Unable to decode payload')
      done();
    })
  })

  it('should return error when negative amount is provided', (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: 0
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: -234 //invalid
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/0'/0'/0/0",
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include('Invalid payee amount')
      done();
    })
  })

  it('should return error when invalid path is provided', (done) => {
    const payload = {
      inputs: [{
        txhash: "81b4c832d70cb56ff957589752eb412a4cab78a25a8fc52d6a09e5bd4404d48a",
        vout: 0
      }],
      outputs: [{
        address: "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa",
        amount: 91234
      }]
    }

    const data = {
      "uuid": uuid,
      "path": "0//0", //invalid path
      "coinType": 0,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      done();
    })
  })

})

describe('Integration test for generating ethereum signature', () => {
  let uuid
  before((done) => {
    const data = {
      "username": "",
      "passphrase": "",
      "mnemonic": "similar supreme jealous custom please fitness mosquito report movie valley hip hub slush foam deer"
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/register',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      uuid = vaultResponse.data.uuid;
      done();
    })
  })

  it("should return a valid signature when valid uuid, payload, path and cointype is provided", (done) => {
    const payload = {
      nonce: 0,
      value: 100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    const actualSignature = "0xf86780058252089495273d64876408e0eda01a45775efc2df6d1cfac88016345785d8a00008025a046300c200def6f3e96135aa279c80cd2858e02a4ac8e03e93ad01276dc443f7ca00f6515899a2293ca29747c776f874dd2c71c03bcfe6a1801d87cd4ac8e7558af"

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.equal(200)
      expect(vaultResponse.data.signature).to.equal(actualSignature)
      done();
    })
  })

  it("should return error when invalid to-address format is provided", (done) => {
    const payload = {
      nonce: 0,
      value: 100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc", //invalid
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Invalid payload data")
      done();
    })
  })

  it("should return error when nonce is negative", (done) => {
    const payload = {
      nonce: -1, //invalid
      value: 100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Unable to decode payload")
      done();
    })
  })

  it("should return error when value to send is negative", (done) => {
    const payload = {
      nonce: 0, //invalid
      value: -100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Unable to decode payload")
      done();
    })
  })

  it("should return error when gas parameters are negative", (done) => {
    const payload = {
      nonce: 0,
      value: 100000000000000000,
      gasLimit: -21000, //invalid
      gasPrice: -5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Unable to decode payload")
      done();
    })
  })

  it("should return error when chainId is negative", (done) => {
    const payload = {
      nonce: 0,
      value: 100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: -1 //invalid
    }

    const data = {
      "uuid": uuid,
      "path": "m/44'/60'/0'/0/0",
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      expect(vaultResponse.errors[0]).to.include("Invalid payload data")
      done();
    })
  })

  it("should return error when invalid path is provided", (done) => {
    const payload = {
      nonce: 0,
      value: 100000000000000000,
      gasLimit: 21000,
      gasPrice: 5,
      to: "0x95273d64876408E0eDa01a45775Efc2Df6d1CfaC",
      data: "",
      chainId: 1
    }

    const data = {
      "uuid": uuid,
      "path": "/44'/0'/0'/0/0", //invalid path
      "coinType": 60,
      "payload": JSON.stringify(payload)
    }

    request.post({
      headers: {
        'X-Vault-Token': rootToken
      },
      uri: 'http://127.0.0.1:8200/v1/api/signature',
      body: JSON.stringify(data),
      method: 'POST'
    }, function (err, response, body) {
      let vaultResponse = JSON.parse(body);
      expect(response.statusCode).to.not.equal(200)
      done();
    })
  })

})

