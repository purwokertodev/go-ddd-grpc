'use strict';

const grpc = require('grpc');
const employeeProto = grpc.load('employee.proto');
const path = require('path');
const fs = require('fs');

// with insecure
//const client = new employeeProto.api.EmployeeService('127.0.0.1:8080', grpc.credentials.createInsecure());

function createEmployee(client, name, age, address, salary){
  let em = {name: name, age: age, address: address, salary: salary}
  client.createEmployee(em, (err, res) => {
    if (err) {
      console.log(err);
    } else {
      console.log(res);
    }
  });
}

function getEmployee(client, id){
  let call = client.getEmployee({key: id});
  call.on('data', (em) => {
    console.log(em)
  });

  call.on('error', (err) => {
    console.log(err)
  });
}

function readKey(){
  return new Promise((f, r) => {
    let file = path.join(__dirname, '../../cert/', 'server.crt');
    fs.readFile(file, 'utf-8', (err, data) => {
      if(err){
        r(err);
      }
      f(data);
    });
  });
};

readKey().then(key => {
  const sslCreds = grpc.credentials.createSsl(Buffer.from(key, 'utf-8'));
  const client = new employeeProto.api.EmployeeService('127.0.0.1:8080', sslCreds.toString());
  getEmployee(client, '9e10c3fc-f3dd-4dca-b7e9-8f1f1b038dcc');
});

//createEmployee("Andi", 17, "Banjarnegara", 8000.0);
